package config

import (
	"os"
	"strconv"
)

type Config struct {
	APP_MODE APP_MODE
	APP_PORT string
	CONN_STR string
	// Observability
	LogLevel                       string
	RateLimiterTimeFrameSeconds    int
	RateLimiterRequestPerTimeFrame int
	// REDIS
	REDIS_ADDR string
	// MONGO
	MONGO_URL string
}

func Load() *Config {

	return &Config{
		APP_PORT:                       getEnv("APP_PORT", "8080"),
		LogLevel:                       getEnv("LOG_LEVEL", string(INFO)),
		RateLimiterTimeFrameSeconds:    getInt("RateLimiter_TimeFrame_Seconds", 60),
		RateLimiterRequestPerTimeFrame: getInt("RateLimiter_Request_Per_TimeFrame", 100),
		CONN_STR:                       getEnv("CONN_STR", ""),
		APP_MODE:                       APP_MODE(getEnv("APP_MODE", string(DEVELOPMENT))),
		REDIS_ADDR:                     getEnv("REDIS_ADDR", "localhost:6379"),
		MONGO_URL:                      getEnv("MONGO_URL", "mongodb://root:example@go-mongo:27017/"),
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

type LOG_LEVEL string

const (
	DEBUG   LOG_LEVEL = "DEBUG"
	INFO    LOG_LEVEL = "INFO"
	WARNING LOG_LEVEL = "WARN"
	ERROR   LOG_LEVEL = "ERROR"
)
