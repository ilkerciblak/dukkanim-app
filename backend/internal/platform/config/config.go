package config

import (
	"os"
	"strconv"
)

type Config struct {
	APP_MODE                       APP_MODE
	APP_PORT                       string
	CONN_STR                       string
	RateLimiterTimeFrameSeconds    int
	RateLimiterRequestPerTimeFrame int
}

func Load() *Config {
	return &Config{
		APP_MODE:                       APP_MODE(getEnv("APP_MODE", string(DEVELOPMENT))),
		APP_PORT:                       getEnv("APP_PORT", "8080"),
		CONN_STR:                       getEnv("CONN_STR", ""),
		RateLimiterTimeFrameSeconds:    getInt("RateLimiter_TimeFrame_Seconds", 60),
		RateLimiterRequestPerTimeFrame: getInt("RateLimiter_Request_Per_TimeFrame", 100),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.ParseInt(value, 2, 8)
		if err != nil {
			return defaultValue
		}

		return int(val)
	}
	return defaultValue
}

type APP_MODE string

const (
	DEVELOPMENT  APP_MODE = "DEV"
	PRODUCTION   APP_MODE = "PROD"
	MAINTAINENCE APP_MODE = "maintainence"
)
