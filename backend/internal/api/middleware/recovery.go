package middleware

import (
	response "dukkanim-api/internal/platform/http_response"
	"fmt"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				// type rr struct {
				// 	Message string `json:"message"`
				// 	Error   string `json:"error"`
				// }

				// res := rr{
				// 	Message: "Internal Server Error",
				// 	Error:   fmt.Sprint(err),
				// }

				response.RespondWithProblemDetails(
					w,
					r.Context(),
					http.StatusInternalServerError,
					fmt.Sprintf("%v", err),
					"SERVER_ERR",
					nil,
				)

				// data, err := json.Marshal(res)
				// if err != nil {
				// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// 	response.RespondWithProblemDetails(w, r.Context(), http.StatusInternalServerError, err.E)

				// 	return
				// }

				// w.WriteHeader(http.StatusInternalServerError)
				// _, _ = w.Write(data)

			}

		}()

		next.ServeHTTP(w, r)
	})
}
