package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/stackus/edat/core"
)

// Request tracking
const (
	RequestIDHeader     = "X-Request-Id"
	CorrelationIDHeader = "X-Correlation-Id"
	CausationIDHeader   = "X-Causation-Id"
)

func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = requestID
		}
		causationID := r.Header.Get(CausationIDHeader)
		if causationID == "" {
			causationID = requestID
		}

		ctx := core.SetRequestContext(r.Context(), requestID, correlationID, causationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetRequestHeaders(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set(RequestIDHeader, core.GetRequestID(ctx))
	w.Header().Set(CorrelationIDHeader, core.GetCorrelationID(ctx))
	w.Header().Set(CausationIDHeader, core.GetCausationID(ctx))
}
