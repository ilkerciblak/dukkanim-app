package middleware

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler

func CreateMiddlewareChain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(routeHandler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			routeHandler = middlewares[i](routeHandler)
		}

		return routeHandler
	}
}
