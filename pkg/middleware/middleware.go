package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// Wraps wraps the router with the middleware functions.
func Wraps(router http.Handler, logger *zap.Logger) http.Handler {
	funcs := []func(http.Handler) http.HandlerFunc{
		LogResult(logger),
		InjectRequestID,
	}

	r := router
	for _, f := range funcs {
		r = f(r)
	}

	return r

}
