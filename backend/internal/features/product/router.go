package product

import (
	"database/sql"
	"dukkanim-api/internal/features/product/domain"
	"dukkanim-api/internal/features/product/internal"
	"dukkanim-api/internal/platform/observability/tracing"

	// "dukkanim-api/internal/platform/observability/tracing"
	"net/http"
)

type productRouteConfig struct {
	httpHandler domain.ProductHandlerInterface
	Mux         *http.ServeMux
	Db          *sql.DB
}

func RegisterProductRoutes(db *sql.DB, tracer tracing.Tracer) *http.ServeMux {

	mux := http.NewServeMux()

	cfg := productRouteConfig{
		httpHandler: &internal.ProductHandler{
			Tracer: tracer,
		},
	}

	mux.Handle("POST /products", http.HandlerFunc(cfg.httpHandler.CreateProduct))

	return mux

}
