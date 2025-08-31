package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func New(host, port, password string, db int) *ClientRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return &ClientRedis{
		client: client,
	}
}

type ClientRedis struct {
	client *redis.Client
}

func (rc *ClientRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *ClientRedis) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %v", key, err)
	}
	return val, nil
}
