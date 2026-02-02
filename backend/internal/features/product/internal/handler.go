package internal

import (
	"dukkanim-api/internal/features/product/domain"
	response "dukkanim-api/internal/platform/http_response"
	"dukkanim-api/internal/platform/observability/logging"
	"dukkanim-api/internal/platform/observability/tracing"

	// "dukkanim-api/internal/platform/observability/tracing"
	"dukkanim-api/internal/platform/problem"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	Tracer tracing.Tracer
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateProductRequest

	// spanCtx, span := otel.Tracer("dukkanim-api").Start(r.Context(), "ProductHandler.CreateProduct")
	spanCtx, span := h.Tracer.Start(r.Context(), "ProductHandler.CreateProduct")

	logger := logging.FromContext(spanCtx)

	logger.With(
		spanCtx,
		logging.SetService("ProductHandler.CreateProduct"),
	)

	defer span.End()

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(tracing.Error, "[FAILED] to Decode Request Body")
		logger.Warning(spanCtx, "Request Decoding Error", "err", err.Error(), "req", r.Body)
		response.RespondWithProblemDetails(w, spanCtx, http.StatusBadRequest, err.Error(), "", nil)
		return
	}

	// logger.With(spanCtx, logging.AppendField("command", request))

	if _, err := request.ToEntity(); err != nil {
		if appErr, k := err.(*problem.Problem); k {
			span.RecordError(err)
			response.RespondWithProblemDetails(w, spanCtx, appErr.StatusCode, appErr.Detail, "", appErr.Errors)
			return
		}
		response.RespondWithProblemDetails(w, spanCtx, http.StatusUnprocessableEntity, err.Error(), "", nil)
		return
	}
	span.SetStatus(tracing.Success, "")
	response.RespondWithJSON(w, http.StatusCreated, nil, spanCtx)
}
