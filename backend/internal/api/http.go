package api

import (
	"dukkanim-api/internal/api/middleware"
	"dukkanim-api/internal/features/product"
	"dukkanim-api/internal/platform/caching"
	"dukkanim-api/internal/platform/config"
	"dukkanim-api/internal/platform/observability/logging"
	"dukkanim-api/internal/platform/observability/metrics"
	"dukkanim-api/internal/platform/observability/tracing"
	ratelimiting "dukkanim-api/internal/platform/rate_limiting"
	"fmt"
	"time"

	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func ServeHttp(cfg *config.Config, db *sql.DB, logger logging.Logger, tracer tracing.Tracer, metric metrics.Metrics, cache caching.Cacher, redis redis.Client, errChan chan<- error) {
	mux := http.NewServeMux()

	//TODO: Yukaridakileri disari tasi cunku ayni seyleri ws' icin falan da kullanabilirsin sonucta

	http_metrics := metrics.HTTPMetrics(metric)

	rl := ratelimiting.FixedWindowAlgorithm(
		ratelimiting.RedisStorageAdapter(&redis),
		ratelimiting.WithWindowDuration(time.Minute),
		ratelimiting.WithRequestCountLimit(10),
	)
	mux.Handle("/health", http.HandlerFunc(healthHandler(db)))
	mux.Handle("/metrics", metric.Handler())

	muxChain := middleware.CreateMiddlewareChain(
		middleware.Recovery,
		middleware.CORS,
		middleware.JSONContentType,
		middleware.RequestLogging(logger),
		// rl.RateLimiting,
		ratelimiting.Middleware(rl),
		middleware.HttpMetrics(*http_metrics),
		middleware.DistributedTracing(tracer),
		middleware.SecurityHeader,
	)

	registerRoutes(
		mux,
		db,
		tracer,
		muxChain,
		product.RegisterProductRoutes,
	)

	server := &http.Server{
		Addr:    fmt.Sprint(":", cfg.APP_PORT),
		Handler: mux,
	}

	go func() {

		if err := server.ListenAndServe(); err != nil {

			errChan <- err
		}
	}()

	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("🛜  Server Running on the Port:\t ", cfg.APP_PORT)
	fmt.Println("📌 APP MODE:\t ", cfg.APP_MODE)
	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("")

}

func registerRoutes(mux *http.ServeMux, db *sql.DB, tracer tracing.Tracer, middlewareChain middleware.MiddlewareFunc, domainRouteRegisteringFunctions ...func(db *sql.DB, tracer tracing.Tracer) *http.ServeMux) {
	for _, f := range domainRouteRegisteringFunctions {
		mux.Handle("/api/", http.StripPrefix("/api", middlewareChain(
			f(db, tracer),
		)))

	}
}

func healthHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"message": "Only HTTP GET is Allowed" }`))
			return
		}
		if _, err := db.ExecContext(r.Context(), "SELECT 1"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e := fmt.Sprintf(`{"status":"Error!", "message":%v}`, err)
			w.Write([]byte(e))
			return
		}

		// if err := mongo.Ping(r.Context(), readpref.Primary()); err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	e := fmt.Sprintf(`{"status":"Error!", "message":%v}`, err)
		// 	w.Write([]byte(e))
		// }

		// if err := redis.Ping(r.Context()).Err(); err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	e := fmt.Sprintf(`{"status":"Error!", "message":%v}`, err)
		// 	w.Write([]byte(e))
		// }

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "OK!"}`))
	}
}
