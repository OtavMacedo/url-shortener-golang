package infra

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse redis url: %w", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("unable to reach redis: %w", err)
	}

	fmt.Println("Connected to redis")
	return client, nil
}
