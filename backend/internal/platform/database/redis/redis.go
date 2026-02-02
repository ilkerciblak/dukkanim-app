package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func newClient(ctx context.Context, addr, password string, defaultDb int) (*redis.Client, error) {
	opt := &redis.Options{

		Addr:     addr,
		DB:       defaultDb,
		Password: password,
	}

	client := redis.NewClient(opt)

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func Redis(ctx context.Context, addr, password string, defaultdb int) (*redis.Client, error) {
	for i := range 4 {
		fmt.Printf("Attempt to instrument Redis connection - attempt#[%d]\n", i+1)
		client, err := newClient(ctx, addr, password, defaultdb)
		if err != nil {
			fmt.Printf("Attempt Failed due to :%v", err)
			continue
		}

		return client, nil
	}

	return nil, fmt.Errorf("Redis client instrumentation failed")
}
