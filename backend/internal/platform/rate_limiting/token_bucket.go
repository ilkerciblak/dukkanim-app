package ratelimiting

import (
	"context"
	"fmt"
	"math"
	"time"
)

type tokenBucketAlgorithm struct {
	storage       Storage
	refillRate    int64
	refillWindow  time.Duration
	tokenCapacity int64
}

func (r tokenBucketAlgorithm) Allow(ctx context.Context, identifier string) (bool, error) {
	tokenKey := fmt.Sprintf("rate_limiting:%s:token", identifier)
	refillTimeKey := fmt.Sprintf("rate_limiting:%s:refillTimeKey", identifier)

	current_token_count, err := r.storage.Get(ctx, tokenKey).Int64()
	if err != nil {
		current_token_count = r.tokenCapacity
	}

	last_refill_time, err := r.storage.Get(ctx, refillTimeKey).Time()
	if err != nil {
		// Time package does not support duration substraction, thus we adding minus duration to now LOL
		last_refill_time = time.Now().Add(-r.refillWindow)
	}
	// token_to_be_added := (time_elapsed_in_seconds / refill_tick_seconds) * refill_token_rate (token per seconds)
	token_to_be_added := int64(math.Floor(float64(time.Since(last_refill_time).Seconds()) / float64(r.refillWindow.Seconds()) * float64(r.refillRate)))
	if token_to_be_added > 0 {
		r.storage.Set(ctx, refillTimeKey, time.Now(), 0)
	}
	refilled_tokens := current_token_count + token_to_be_added
	refilled_tokens = min(r.tokenCapacity, refilled_tokens)

	if refilled_tokens >= 1 {
		_ = r.storage.Set(ctx, tokenKey, refilled_tokens-1, 0)
		return true, nil
	}
	_ = r.storage.Set(ctx, tokenKey, refilled_tokens, 0)

	return false, nil

}

func (r tokenBucketAlgorithm) Reset(ctx context.Context, identifier string) error {
	tokenKey := fmt.Sprintf("rate_limiting:%s:token", identifier)
	refillTimeKey := fmt.Sprintf("rate_limiting:%s:refillTimeKey", identifier)
	if err := r.storage.Set(ctx, tokenKey, r.tokenCapacity, 0); err != nil {
		return err
	}
	if err := r.storage.Set(ctx, refillTimeKey, time.Now(), 0); err != nil {
		return err
	}

	return nil
}

func (r tokenBucketAlgorithm) GetLimit(ctx context.Context, identifier string) (LimitInfo, error) {
	tokenKey := fmt.Sprintf("rate_limiting:%s:token", identifier)
	refillTimeKey := fmt.Sprintf("rate_limiting:%s:refillTimeKey", identifier)

	remaining_token, err := r.storage.Get(ctx, tokenKey).Int64()
	if err != nil {
		return LimitInfo{}, fmt.Errorf("[rate_limiting.TokenBucket.GetLimit] Token Query Failed:\n%v", err)
	}
	last_refill_time, err := r.storage.Get(ctx, refillTimeKey).Time()
	if err != nil {
		return LimitInfo{}, fmt.Errorf("[rate_limiting.TokenBucket.GetLimit] Refill Time Query Failed:\n%v", err)
	}

	reset_time := last_refill_time.Add(time.Duration(r.refillWindow.Seconds()))

	return LimitInfo{
		Remaining: remaining_token,
		ResetTime: reset_time,
	}, nil

}

type WithBucketConfig func(*algorithmConfig) *algorithmConfig

func WithRefillRate(rate int64) WithBucketConfig {
	return func(ac *algorithmConfig) *algorithmConfig {
		ac.refillRate = rate
		return ac
	}
}

func WithRefillWindow(duration time.Duration) WithBucketConfig {
	return func(ac *algorithmConfig) *algorithmConfig {
		ac.refillWindow = duration
		return ac
	}
}

func WithTokenCapacity(capacity int64) WithBucketConfig {
	return func(ac *algorithmConfig) *algorithmConfig {
		ac.tokenCapacity = capacity

		return ac
	}
}

// TODO: Should return RateLimitingStrategy interface
func TokenBucket(storage Storage, manipulators ...WithBucketConfig) RateLimitingStrategy {
	cfg := defaultCfg()

	for _, f := range manipulators {
		cfg = f(cfg)
	}

	return &tokenBucketAlgorithm{
		storage:       storage,
		refillRate:    cfg.refillRate,
		refillWindow:  cfg.refillWindow,
		tokenCapacity: cfg.tokenCapacity,
	}
}
