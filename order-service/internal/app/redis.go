package app

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"order-service/internal/config"
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
		return nil, fmt.Errorf("failed to connect to the redis %s:%s, %s, %s : %v", cfg.Host, cfg.Port, cfg.Username, cfg.Password, err)
	}

	return client, nil
}
