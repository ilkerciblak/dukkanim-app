package main

import (
	"context"
	"dukkanim-api/internal/api"
	"dukkanim-api/internal/platform/config"
	"dukkanim-api/internal/platform/database/mongodb"
	"dukkanim-api/internal/platform/database/postgres"
	"dukkanim-api/internal/platform/database/redis"
	logging "dukkanim-api/internal/platform/observability/logging/adapter"
	metrics "dukkanim-api/internal/platform/observability/metrics/adapter"
	tracing "dukkanim-api/internal/platform/observability/tracing/adapter"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// 1. Configurations

	errChan := make(chan error, 1)

	// Environment Variables Initialization
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. PLATFORM INITIALIZATIONS

	// 2.1 DB CONNECTION
	// db_config := db.NewSqlConnectionConfig("postgres", cfg.CONN_STR)

	// conn := db_config.InitializeConnection(errChan)

	sql_db, err := postgres.PostgresDB(ctx, cfg.CONN_STR)
	if err != nil {
		errChan <- err
	}
	defer sql_db.Close()

	mongodb, err := mongodb.MongoDB(ctx, cfg.MONGO_URL)
	if err != nil {
		errChan <- err
	}

	go func() {
		if err = mongodb.Disconnect(ctx); err != nil {
			errChan <- err
		}
	}()

	redisDb, err := redis.Redis(ctx, cfg.REDIS_ADDR, "", 0)
	if err != nil {
		errChan <- err
	}

	defer redisDb.Conn().Close()

	if err := goose.Up(sql_db.Connection, "migrations/"); err != nil {
		errChan <- err
	}

	// 3. Initialize Platform Integrations
	logger := logging.NewslogLogger(cfg.LogLevel)
	tracer := tracing.NewOtelTracer()
	metric := metrics.PrometheusMetrics()

	api.ServeHttp(cfg, sql_db.Connection, logger, tracer, metric, errChan)

	if err := <-errChan; err != nil {
		fmt.Printf("\nSERVER CLOSING DUE TO: \n%v", err)
		os.Exit(1)
	}

}
