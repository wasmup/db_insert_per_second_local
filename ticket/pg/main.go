package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	_ "embed"
)

func main() {
	mode := flag.String("mode", "loadtest", "serve|seed|loadtest")
	count := flag.Int("count", 100_000, "Number of tickets to seed (for seed mode)")
	eventID := flag.Int64("event", 1, "Event ID to seed or operate on")

	dsn := flag.String("dsn", os.Getenv("PG_DSN"), "PostgreSQL DSN, e.g. postgres://user:pass@localhost:5432/db?sslmode=disable")

	requests := flag.Int("requests", 1_000, "Total requests for loadtest (for loadtest mode)")
	concurrency := flag.Int("concurrency", runtime.NumCPU(), "Concurrency for loadtest (for loadtest mode)")
	url := flag.String("url", "http://localhost:8080/sell", "URL for loadtest (for loadtest mode)")

	flag.Parse()

	if *dsn == "" {
		slog.Error(`DSN failed`)
		return
	}

	ctx := context.Background()
	var err error
	pool, err = pgxpool.New(ctx, *dsn)
	if err != nil {
		slog.Error(`pool failed`, `error`, err)
		return
	}
	defer pool.Close()

	if err := ensureSchema(ctx, pool); err != nil {
		slog.Error(`schema setup failed`, `error`, err)
		return
	}

	switch *mode {
	case "seed":
		log.Println("Seeding tickets...")
		if err := seedTickets(ctx, pool, *eventID, *count); err != nil {
			slog.Error(`seed failed`, `error`, err)
			return
		}
		slog.Info(`Seeding complete.`)

	case "loadtest":
		go serve()
		time.Sleep(1 * time.Second)

		if strings.TrimSpace(*url) == "" {
			slog.Error(`url failed`, `error`, err)
			return
		}
		fmt.Printf("Starting load test: url=%s concurrency=%d total=%d\n", *url, *concurrency, *requests)
		start := time.Now()
		success, failed := runLoadTest(*url, *concurrency, *requests, *eventID)
		elapsed := time.Since(start)
		fmt.Printf("Load test completed. Success=%d, Failed=%d, elapsed=%s, throughput=%.2f req/s\n",
			success, failed, elapsed, float64(success)/elapsed.Seconds())
		return
	}

	serve()
}

func serve() {
	slog.Info(`Starting HTTP server on :8080`)
	http.HandleFunc("/sell", sellHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error(`server failed`, `error`, err)
		return
	}
}

func seedTickets(ctx context.Context, p *pgxpool.Pool, eventID int64, count int) error {
	batch := 1000
	for done := 0; done < count; done += batch {
		n := batch
		if done+n > count {
			n = count - done
		}
		placeholders := make([]string, 0, n)
		args := make([]any, 0, 2*n)
		for i := 0; i < n; i++ {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", len(args)+1, len(args)+2))
			args = append(args, eventID, "AVAILABLE")
		}
		stmt := "INSERT INTO tickets(event_id, status) VALUES " + strings.Join(placeholders, ",")
		_, err := p.Exec(ctx, stmt, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func sellHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SellRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.EventID == 0 {
		req.EventID = 1
	}
	if req.UserID == "" {
		req.UserID = "anonymous"
	}
	if req.HoldMinutes <= 0 {
		req.HoldMinutes = 15
	}

	ctx := r.Context()
	var ticketID int64

	// Atomic operation inside a transaction
	// The CTE selects one AVAILABLE ticket for the given event, locks it, and updates it to HELD
	// Returning the ticket_id on success
	err := pool.QueryRow(ctx, query, req.EventID, req.UserID, req.HoldMinutes).Scan(&ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "sold_out"})
			return
		}
		http.Error(w, "internal_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseTicket{TicketID: ticketID})
}

func runLoadTest(url string, concurrency int, total int, eventID int64) (int64, int64) {
	var (
		success int64
		failed  int64
		wg      sync.WaitGroup
		sem     = make(chan struct{}, concurrency)
		client  = &http.Client{
			Timeout: 10 * time.Second,
		}
	)

	for i := range total {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			req := SellRequest{
				EventID:     eventID,
				UserID:      fmt.Sprintf("loadtest-%d", i),
				HoldMinutes: 15,
			}
			body, err := json.Marshal(req)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}

			reqHTTP, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			reqHTTP.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(reqHTTP)
			if err != nil {
				atomic.AddInt64(&failed, 1)
				return
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				atomic.AddInt64(&success, 1)
			} else {
				atomic.AddInt64(&failed, 1)
			}
		}(i)
	}

	wg.Wait()
	return success, failed
}

//go:embed query.sql
var query string

//go:embed create_schema.sql
var queryCreateSchema string

//go:embed create_index.sql
var queryCreateIndex string

func ensureSchema(ctx context.Context, p *pgxpool.Pool) error {
	_, err := p.Exec(ctx, queryCreateSchema)
	if err != nil {
		return err
	}

	_, err = p.Exec(ctx, queryCreateIndex)
	if err != nil {
		return err
	}

	return nil
}

type SellRequest struct {
	EventID     int64  `json:"event_id"`
	UserID      string `json:"user_id"`
	HoldMinutes int    `json:"hold_minutes"`
}

type ResponseTicket struct {
	TicketID int64 `json:"ticket_id"`
}

var pool *pgxpool.Pool
