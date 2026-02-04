package caching

import (
	"context"
)

type Cacher interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl float64) error
	Exists(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
