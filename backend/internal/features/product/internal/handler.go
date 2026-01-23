package internal

import (
	"dukkanim-api/internal/features/product/domain"
	response "dukkanim-api/internal/platform/http_response"
	"dukkanim-api/internal/platform/logging"
	"dukkanim-api/internal/platform/problem"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateProductRequest
	logger := logging.FromContext(r.Context())

	logger.DEBUG(r.Context(), "Create Product Service", "service", "product-create")

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.ERROR(r.Context(), "Failed to Decode Request Body", "body", r.Body, "err", err.Error())
		response.RespondWithProblemDetails(w, r.Context(), http.StatusBadRequest, err.Error(), "", nil)
		return
	}

	if _, err := request.ToEntity(); err != nil {
		if appErr, k := err.(*problem.Problem); k {
			logger.ERROR(r.Context(), "Validation Error", "err", appErr.Errors)
			response.RespondWithProblemDetails(w, r.Context(), appErr.StatusCode, appErr.Detail, "", appErr.Errors)
			return
		}
		logger.ERROR(r.Context(), "Validation Error", "err", err)
		response.RespondWithProblemDetails(w, r.Context(), http.StatusUnprocessableEntity, err.Error(), "", nil)
		return
	}

	logger.DEBUG(r.Context(), "Create Product Handler Succeed")
	response.RespondWithJSON(w, http.StatusCreated, nil, r.Context())
}
