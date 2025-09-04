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

const BATCH_SIZE = 1000

func main() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
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

	pipe := rdb.Pipeline()

	for time.Since(startTime) < maxDuration {
		impressionKey := fmt.Sprintf("impression:%s", uuid.New().String())
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"

		pipe.HSet(ctx, impressionKey,
			"ad_id", adID,
			"image_url", imageURL,
			"click_url", clickURL)

		count++

		if count%BATCH_SIZE == 0 {
			_, err := pipe.Exec(ctx)
			if err != nil {
				slog.Error("pipeline exec failed", "error", err)
				return
			}
			pipe = rdb.Pipeline()
		}
	}

	if count%BATCH_SIZE != 0 {
		_, err := pipe.Exec(ctx)
		if err != nil {
			slog.Error("final pipeline exec failed", "error", err)
			return
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d impressions in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
