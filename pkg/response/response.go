package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	apierror "github.com/nkitlabs/go-http-gorm-example/pkg/errors"
	"github.com/nkitlabs/go-http-gorm-example/pkg/middleware"
)

func WriteError(ctx context.Context, w http.ResponseWriter, err error, log *zap.Logger) {
	reqID, ok := ctx.Value(middleware.ContextKeyRequestID).(string)
	if !ok {
		reqID = middleware.RequestIDUnknown
	}

	log.Error(err.Error(), zap.String(middleware.LogKeyID, reqID))

	e := apierror.ToError(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)

	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Error(fmt.Sprintf("failed to write error response: %v", err), zap.String(middleware.LogKeyID, reqID))
		http.Error(w, e.Error(), e.Code)
	}
}

func Write(ctx context.Context, w http.ResponseWriter, code int, data any, log *zap.Logger) {
	reqID, ok := ctx.Value(middleware.ContextKeyRequestID).(string)
	if !ok {
		reqID = middleware.RequestIDUnknown
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error(fmt.Sprintf("failed to write error response: %v", err), zap.String(middleware.LogKeyID, reqID))
		http.Error(w, apierror.ErrInternal.Error(), apierror.ErrInternal.Code)
	}
}
