package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// InjectRequestID injects a request ID into the context of the request.
func InjectRequestID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyRequestID, id.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(ContextKeyRequestID).(string)
	if !ok {
		return RequestIDUnknown
	}

	return reqID
}
