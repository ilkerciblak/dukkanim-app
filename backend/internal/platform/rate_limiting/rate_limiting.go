package ratelimiting

import (
	"context"
	"strconv"
	"time"
)

type RateLimitingStrategy interface {
	Allow(ctx context.Context, identifier string) (bool, error)
	Reset(ctx context.Context, identifier string) error
	GetLimit(ctx context.Context, identifier string) (LimitInfo, error)
}

type LimitInfo struct {
	Remaining int64
	ResetTime time.Time
}

type Storage interface {
	Get(ctx context.Context, key string) *storageResult
	Set(ctx context.Context, key string, value any, expirety time.Duration) error
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) error
	TTL(ctx context.Context, key string) (*time.Duration, error)
}

type storageResult struct {
	val string
	err error
}

func (cmd *storageResult) Result() (string, error) {
	if cmd.err != nil {
		return "", cmd.err
	}
	return cmd.val, cmd.err
}

func (cmd *storageResult) Val() string {
	return cmd.val
}

func (cmd *storageResult) Bool() (bool, error) {
	if cmd.err != nil {
		return false, cmd.err
	}
	return strconv.ParseBool(cmd.val)
}

func (cmd *storageResult) Int() (int, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.Atoi(cmd.Val())
}

func (cmd *storageResult) Int64() (int64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseInt(cmd.Val(), 10, 64)
}

func (cmd *storageResult) Uint64() (uint64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseUint(cmd.Val(), 10, 64)
}

func (cmd *storageResult) Float32() (float32, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	f, err := strconv.ParseFloat(cmd.Val(), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func (cmd *storageResult) Float64() (float64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseFloat(cmd.Val(), 64)
}

func (cmd *storageResult) Time() (time.Time, error) {
	if cmd.err != nil {
		return time.Time{}, cmd.err
	}
	return time.Parse(time.RFC3339Nano, cmd.Val())
}

type algorithmConfig struct {
	refillRate    int64
	tokenCapacity int64
	refillWindow  time.Duration
}

func defaultCfg() *algorithmConfig {

	return &algorithmConfig{
		refillRate:    5,
		refillWindow:  time.Second * 60,
		tokenCapacity: 5,
	}
}
