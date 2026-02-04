package ratelimiting

import (
	"context"
	"fmt"
	"time"
)

type fixedWindowAlgorithm struct {
	storage           Storage
	refillWindow      time.Duration
	requestCountLimit int64
}

type WithFixedWindowConfig func(*algorithmConfig) *algorithmConfig

func WithWindowDuration(dur time.Duration) WithFixedWindowConfig {
	return func(ac *algorithmConfig) *algorithmConfig {
		ac.refillWindow = dur
		return ac
	}
}

func WithRequestCountLimit(limit int64) WithFixedWindowConfig {
	return func(ac *algorithmConfig) *algorithmConfig {
		ac.tokenCapacity = limit
		return ac
	}
}

func FixedWindowAlgorithm(storage Storage, manipulators ...WithFixedWindowConfig) RateLimitingStrategy {
	cfg := defaultCfg()
	for _, f := range manipulators {
		cfg = f(cfg)
	}

	return &fixedWindowAlgorithm{
		storage:           storage,
		refillWindow:      cfg.refillWindow,
		requestCountLimit: cfg.tokenCapacity,
	}
}

func (r *fixedWindowAlgorithm) Allow(ctx context.Context, identifier string) (bool, error) {

	counterKey := fmt.Sprintf("rate_limiting:%s:request_counter", identifier)

	requets_count, err := r.storage.Get(ctx, counterKey).Int64()
	if err != nil {
		requets_count = 0
		r.storage.Set(ctx, counterKey, requets_count, r.refillWindow)
	}

	if requets_count < r.requestCountLimit {
		if err := r.storage.Increment(ctx, counterKey); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil

}

func (r *fixedWindowAlgorithm) Reset(ctx context.Context, identifier string) error {
	counterKey := fmt.Sprintf("rate_limiting:%s:request_counter", identifier)
	r.storage.Set(ctx, counterKey, 0, r.refillWindow)
	return nil
}

func (r *fixedWindowAlgorithm) GetLimit(ctx context.Context, identifier string) (LimitInfo, error) {
	counterKey := fmt.Sprintf("rate_limiting:%s:request_counter", identifier)

	current_count, err := r.storage.Get(ctx, counterKey).Int64()
	if err != nil {
		return LimitInfo{}, fmt.Errorf("[rate_limiting.FixedWindow.GetLimit] Token Query Failed:\n%v", err)
	}

	remaning_ttl, err := r.storage.TTL(ctx, counterKey)
	if err != nil {
		return LimitInfo{}, fmt.Errorf("[rate_limiting.FixedWindow.GetLimit] TTL Query Failed:\n%v", err)
	}
	reset_time := time.Now().Add(*remaning_ttl)

	return LimitInfo{
		Remaining: r.requestCountLimit - current_count,
		ResetTime: reset_time,
	}, nil
}
