package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// ReadUserIP reads the user IP address from the request headers.
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// LogResult logs the result status from the request.
func LogResult(logger *zap.Logger) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := statusRecorder{w, http.StatusOK}

			next.ServeHTTP(&rec, r)

			ctx := r.Context()
			id := GetRequestID(ctx)

			fields := []zap.Field{
				zap.Int(LogKeyStatus, rec.status),
				zap.String(LogKeyLatency, time.Since(start).String()),
				zap.String(LogKeyID, id),
				zap.String(LogKeyMethod, r.Method),
				zap.String(LogKeyURI, r.RequestURI),
				zap.String(LogKeyHost, r.Host),
				zap.String(LogKeyRemoteIP, ReadUserIP(r)),
			}

			logger.Info(fmt.Sprintf("request: %s %s; response status: %d", r.Method, r.URL, rec.status), fields...)
		}
	}
}
