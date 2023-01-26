package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Provider struct {
	client     *redis.Client
	timeExpire time.Duration
}

func NewRedisProvider(client *redis.Client) *Provider {
	return &Provider{
		client:     client,
		timeExpire: time.Hour,
	}
}

func (r *Provider) Write(ctx context.Context, key string, value int64) error {

	res := r.client.Set(ctx, key, value, r.timeExpire)
	return res.Err()
}

func (r *Provider) Exist(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return res != 0, nil
}

func (r *Provider) Read(ctx context.Context, key string) (int64, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("couldn't find value for key %s: %v", key, err)
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed cast: %v", err)
	}

	return int64(res), nil
}
