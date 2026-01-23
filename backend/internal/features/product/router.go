package product

import (
	"database/sql"
	"dukkanim-api/internal/features/product/domain"
	"dukkanim-api/internal/features/product/internal"
	"net/http"
)

type productRouteConfig struct {
	httpHandler domain.ProductHandlerInterface
	Mux         *http.ServeMux
	Db          *sql.DB
}

func RegisterProductRoutes(mux *http.ServeMux, db *sql.DB) {

	cfg := productRouteConfig{
		httpHandler: &internal.ProductHandler{},
	}

	mux.Handle("POST /products", http.HandlerFunc(cfg.httpHandler.CreateProduct))

}
