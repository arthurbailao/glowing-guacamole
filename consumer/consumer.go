package consumer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/arthurbailao/cint-ad/detector"
	"github.com/redis/go-redis/v9"
)

// Consumer is a consumer of data points
type Consumer struct {
	redisClient *redis.Client
	channel     string
	detector    *detector.Detector
}

// NewConsumer creates a new consumer
func NewConsumer(
	redisClient *redis.Client,
	channel string,
	detector *detector.Detector,
) *Consumer {
	return &Consumer{
		redisClient: redisClient,
		channel:     channel,
		detector:    detector,
	}
}

// Start starts the consumer
func (c *Consumer) Start(ctx context.Context) error {
	pubsub := c.redisClient.Subscribe(ctx, c.channel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	for msg := range pubsub.Channel() {
		if err := c.processMessage(msg, time.Now()); err == nil {
			continue
		}

		// TODO: If the message is not valid, we should log it and continue

	}

	return nil
}

func (c *Consumer) processMessage(msg *redis.Message, ts time.Time) error {
	value, err := strconv.ParseFloat(msg.Payload, 64)

	if err != nil {
		return fmt.Errorf("failed to parse float value: %w", err)
	}

	z, anomaly, warming := c.detector.Detect(value)

	if warming {
		return nil
	}

	timestamp := ts.Format(time.RFC3339)

	if anomaly {
		fmt.Printf(
			"[%s] Data point: %.2f | Status: ANOMALY DETECTED! | Z-score: %.2f | ALERT: Significant deviation detected\n",
			timestamp,
			value,
			z,
		)
	} else {
		fmt.Printf("[%s] Data point: %.2f | Status: OK | Z-score: %.2f\n", timestamp, value, z)
	}

	return nil
}
