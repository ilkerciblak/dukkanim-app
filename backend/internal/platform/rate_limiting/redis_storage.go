package ratelimiting

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	client redis.Client
}

func RedisStorageAdapter(client *redis.Client) *redisStorage {
	return &redisStorage{
		client: *client,
	}

}

func (r redisStorage) Get(ctx context.Context, key string) *storageResult {
	res, err := r.client.Get(ctx, key).Result()

	return &storageResult{
		val: res,
		err: err,
	}

}

func (r redisStorage) Set(ctx context.Context, key string, value any, expirety time.Duration) error {
	if _, err := r.client.Set(ctx, key, value, expirety).Result(); err != nil {
		return fmt.Errorf("[RedisStorageAdapter (Set) Error]: %v", err)
	}

	return nil
}

func (r redisStorage) Delete(ctx context.Context, key string) error {
	if _, err := r.client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("[RedisStorageAdapter (Delete) Error]: %v", err)
	}

	return nil
}

func (r redisStorage) Increment(ctx context.Context, key string) error {
	if _, err := r.client.Incr(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

// TODO: Refactor TTL Method to return error or time.duration
func (r redisStorage) TTL(ctx context.Context, key string) (*time.Duration, error) {
	dur, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &dur, nil

}
