package consumer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/arthurbailao/cint-ad/consumer"
	"github.com/arthurbailao/cint-ad/detector"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Run runs the consumer
func Run() {
	fx.New(
		fx.Provide(
			context.Background,
			newRedisClient,
			newConsumer,
			newDetector,
		),
		fx.Invoke(runConsumer),
	).Run()
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

func newConsumer(client *redis.Client, detector *detector.Detector) *consumer.Consumer {
	channel := os.Getenv("REDIS_PUBSUB_CHANNEL")

	if channel == "" {
		channel = "data"
	}

	return consumer.NewConsumer(client, channel, detector)
}

func newDetector() (*detector.Detector, error) {
	thresholdEnv := os.Getenv("Z_SCORE_THRESHOLD")

	if thresholdEnv == "" {
		thresholdEnv = "3.0"
	}

	threshold, err := strconv.ParseFloat(thresholdEnv, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Z_SCORE_THRESHOLD: %w", err)
	}

	windowSizeEnv := os.Getenv("DETECTOR_WINDOW_SIZE")

	if windowSizeEnv == "" {
		windowSizeEnv = "50"
	}

	windowSize, err := strconv.ParseInt(windowSizeEnv, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DETECTOR_WINDOW_SIZE: %w", err)
	}

	if windowSize < 5 {
		return nil, fmt.Errorf("DETECTOR_WINDOW_SIZE must be greater than 5")
	}

	return detector.NewDetector(threshold, int(windowSize)), nil
}

func runConsumer(c *consumer.Consumer, ctx context.Context) error {
	return c.Start(ctx)
}
