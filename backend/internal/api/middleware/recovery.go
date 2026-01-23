package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				type response struct {
					Message string `json:"message"`
					Error   string `json:"error"`
				}

				res := response{
					Message: "Internal Server Error",
					Error:   fmt.Sprint(err),
				}

				data, err := json.Marshal(res)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(data)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
