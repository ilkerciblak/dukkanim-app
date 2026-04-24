package usermanagement

import (
	"database/sql"
	"dukkanim-api/internal/platform/observability/tracing"
	"net/http"
)

func RegisterUserManagementHttpRoutes(
	db *sql.DB,
	tracer tracing.Tracer,
) *http.ServeMux {
	mux := http.NewServeMux()

	// define domain config

	// define restful endpoints


	return mux
}
