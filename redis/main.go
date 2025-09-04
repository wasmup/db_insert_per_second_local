package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Addr:     os.Getenv("REDIS_ADDR"),
		// Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0, // use default DB
	})

	ctx := context.Background()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		slog.Error("redis ping failed", "error", err)
		return
	}

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0

	for time.Since(startTime) < maxDuration {
		impressionKey := fmt.Sprintf("impression:%s", uuid.New().String())
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"

		err := rdb.HSet(ctx, impressionKey,
			"ad_id", adID,
			"image_url", imageURL,
			"click_url", clickURL).Err()
		if err != nil {
			slog.Error("redis HSet failed", "error", err)
			return
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d impressions in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
