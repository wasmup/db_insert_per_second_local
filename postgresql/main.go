package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	dsn := os.ExpandEnv("user=$PSQL_USERNAME password=$PSQL_PASSWORD host=localhost port=5432 dbname=$PSQL_DB sslmode=disable")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("db open failed", "error", err)
		return
	}
	defer db.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url)
              VALUES ($1, $2, $3, $4)`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0
	ctx := context.Background()

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"

		_, err := db.ExecContext(ctx, query, impressionID, adID, imageURL, clickURL)
		if err != nil {
			slog.Error("db insert failed", "error", err)
			return
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
