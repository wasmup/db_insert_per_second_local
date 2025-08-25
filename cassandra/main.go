package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "ks1"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer session.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url, impression_time)
              VALUES (?, ?, ?, ?, ?)`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"
		impressionTime := time.Now().UTC()

		if err := session.Query(query, impressionID, adID, imageURL, clickURL, impressionTime).Exec(); err != nil {
			log.Fatalf("Insert failed: %v", err)
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
