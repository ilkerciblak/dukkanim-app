package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ProblemDetails struct {
	Type      string            `json:"type"`  // comes from status
	Title     string            `json:"title"` // comes from status
	Status    int               `json:"status"`
	Detail    string            `json:"detail,omitempty"`
	Code      string            `json:"code"`
	Instance  string            `json:"instance"`   // ctx
	TraceId   string            `json:"trace_id"`   // ctx
	RequestId string            `json:"request_id"` // ctx
	Errors    map[string]string `json:"errors,omitempty"`
}

var detailsMap map[int]map[string]string = map[int]map[string]string{

	400: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/400",
		"title": "Bad Request",
	},
	401: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/401",
		"title": "Unauthorized",
	},
	403: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/403",
		"title": "Forbidden",
	},
	404: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/404",
		"title": "Not Found",
	},
	405: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/405",
		"title": "Method Not Allowed",
	},
	406: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/406",
		"title": "Not Acceptable",
	},

	407: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/407",
		"title": "Proxy Auth Required",
	},

	408: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/408",
		"title": "Request Timeout",
	},

	409: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/409",
		"title": "Conflict",
	},

	410: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/410",
		"title": "Gone",
	},

	411: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/411",
		"title": "Length Required",
	},

	412: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/412",
		"title": "Precondition Failed",
	},

	413: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/413",
		"title": "Request Entity Too Large",
	},

	414: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/414",
		"title": "Request URI Too Long",
	},

	415: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/415",
		"title": "Unsupported Media Type",
	},

	416: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/416",
		"title": "Bad Request",
	},

	417: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/417",
		"title": "Bad Request",
	},

	418: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/418",
		"title": "I'm a teapot",
	},

	421: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/421",
		"title": "Bad Request",
	},

	422: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/422",
		"title": "Unproccessable Entity",
	},

	423: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/423",
		"title": "Locked",
	},

	424: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/424",
		"title": "Bad Request",
	},

	425: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/425",
		"title": "Locked",
	},

	426: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/426",
		"title": "Bad Request",
	},

	428: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/428",
		"title": "Bad Request",
	},

	429: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/429",
		"title": "Too Many Request",
	},

	431: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/431",
		"title": "Bad Request",
	},

	451: {
		"type":  "https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/451",
		"title": "Bad Request",
	},
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any, ctx context.Context) {
	w.WriteHeader(statusCode)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			RespondWithProblemDetails(w, ctx, http.StatusInternalServerError, fmt.Sprintf("Parsing Error With: %v", err), "", nil)
		}
	}

}

func RespondWithProblemDetails(w http.ResponseWriter, ctx context.Context, statusCode int, detail, code string, errors map[string]string) {
	pd := constructProblemDetails(statusCode, detail, code, ctx, errors)

	data, err := json.Marshal(pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(pd.Status)
	w.Write(data)

}

func constructProblemDetails(statusCode int, detail, code string, ctx context.Context, errors map[string]string) ProblemDetails {
	var pd ProblemDetails = ProblemDetails{
		Status: statusCode,
		Type:   detailsMap[statusCode]["type"],
		Title:  detailsMap[statusCode]["title"],
		Code:   code,
		Detail: detail}

	if id, k := ctx.Value(RequestIdKey).(uuid.UUID); k {
		pd.RequestId = id.String()
	}

	if traceId, k := ctx.Value("trace-id").(string); k {
		pd.TraceId = traceId
	}

	if instace, k := ctx.Value("request-instance").(string); k {
		pd.Instance = instace
	}

	if len(errors) != 0 {
		pd.Errors = errors
	}

	return pd
}

type key int

var RequestIdKey key
