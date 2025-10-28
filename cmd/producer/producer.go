package producer

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Run() {
	client, err := newRedisClient()
	if err != nil {
		panic(err)
	}

	throughputInput := os.Getenv("PRODUCER_THROUGHPUT_PER_SECOND")
	if throughputInput == "" {
		throughputInput = "5"
	}

	throughput, err := strconv.ParseInt(throughputInput, 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse PRODUCER_THROUGHPUT_PER_SECOND: %w", err))
	}

	ticker := time.NewTicker(time.Millisecond * 1000 / time.Duration(throughput))

	var x int64
	for range ticker.C {
		y := generateData(x)
		x++
		client.Publish(context.Background(), "data", strconv.FormatFloat(y, 'f', -1, 64))
	}
}

func generateData(x int64) float64 {
	y := 5 + float64(x)*0.01
	y += rand.NormFloat64()*0.15 + 0.5

	if rand.Float64() < 0.1 {
		y += float64(rand.IntN(30)) + 10.0
	}

	return y
}

func newRedisClient() (*redis.Client, error) {
	url, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return nil, errors.New("REDIS_URL not set")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_URL: %w", err)
	}

	return redis.NewClient(opts), nil
}
