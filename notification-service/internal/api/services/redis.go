package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	*redis.Client
}

func NewRedisClient() (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisService{client}, nil
}

func (r *RedisService) IsProcessed(ctx context.Context, bookingID string) (bool, error) {
	exists, err := r.Exists(ctx, fmt.Sprintf("notification.sent.%s", bookingID)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisService) MarkProcessed(ctx context.Context, bookingID string) error {
	return r.SetEX(ctx, fmt.Sprintf("notification.sent.%s", bookingID), "processed", 24*time.Hour).Err() // TTL 24h
}
