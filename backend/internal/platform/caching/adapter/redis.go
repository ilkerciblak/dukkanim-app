package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheAdapter struct {
	client *redis.Client
}

func RedisAdapter(ctx context.Context, client *redis.Client) *redisCacheAdapter {
	return &redisCacheAdapter{
		client: client,
	}
}

func (c *redisCacheAdapter) Get(ctx context.Context, key string) (any, error) {
	res, err := c.client.JSONGet(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (c *redisCacheAdapter) Set(ctx context.Context, key string, value any, ttl float64) error {
	if _, err := c.client.Set(ctx, key, value, time.Minute*time.Duration(ttl)).Result(); err != nil {
		return err
	}
	return nil
}
func (c *redisCacheAdapter) Exists(ctx context.Context, key string) bool {
	if _, err := c.client.Exists(ctx, key).Result(); err != nil {
		return false
	}

	return true
}
func (c *redisCacheAdapter) Delete(ctx context.Context, key string) error {

	return c.client.Del(ctx, key).Err()

}
func (c *redisCacheAdapter) Clear(ctx context.Context) error {
	// c.client.
	return fmt.Errorf("Not Impelemented For Redis Adapter")
}
