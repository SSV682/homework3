package app

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"user-service/internal/config"
)

func initRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Network:  cfg.Network,
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to the redis: %w", err)
	}

	return client, nil
}
