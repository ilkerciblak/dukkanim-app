package api

import (
	"dukkanim-api/internal/api/middleware"
	"dukkanim-api/internal/features/product"
	"dukkanim-api/internal/platform/config"
	"dukkanim-api/internal/platform/logging"
	loggers "dukkanim-api/internal/platform/logging/adapters"
	"fmt"
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
)

func Serve(cfg *config.Config, db *sql.DB) *http.Server {
	mux := http.NewServeMux()

	logger := registerPlatformFeatures(cfg)

	rl := middleware.NewRateLimitingMiddleware(middleware.RatelimiterConfig{
		RequestPerTimeFrame: cfg.RateLimiterRequestPerTimeFrame,
		TimeFrameSeconds:    cfg.RateLimiterTimeFrameSeconds,
	})

	muxChain := middleware.CreateMiddlewareChain(
		middleware.Recovery,
		middleware.CORS,
		middleware.JSONContentType,
		middleware.RequestLogging(logger),
		rl.RateLimiting,
		middleware.SecurityHeader,
	)

	server := &http.Server{
		Addr: fmt.Sprint(":", cfg.APP_PORT),

		Handler: muxChain(mux),
	}
	healthService := &HealthServiceHandler{}

	mux.Handle("/health", http.HandlerFunc(healthService.healthHandler))

	registerRoutes(mux, db)

	return server
}

func registerRoutes(mux *http.ServeMux, db *sql.DB) {
	product.RegisterProductRoutes(mux, db)
}

type HealthServiceHandler struct {
}

func (h *HealthServiceHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"message": "Only HTTP GET is Allowed" }`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "OK!"}`))

}

func registerPlatformFeatures(cfg *config.Config) logging.Logger {

	var logger logging.Logger

	if cfg.APP_MODE == config.DEVELOPMENT {
		logger = loggers.NewSlogLogger()
	}

	return logger
}
