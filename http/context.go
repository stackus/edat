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

// RequestContext is an http.Handler middleware that sets the id, correlation, and causation ids into context
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

		w.Header().Set(RequestIDHeader, core.GetRequestID(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetResponseHeaders puts the id, correlation, and causation ids into the outgoing http response
func SetResponseHeaders(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set(RequestIDHeader, core.GetRequestID(ctx))
	w.Header().Set(CorrelationIDHeader, core.GetCorrelationID(ctx))
	w.Header().Set(CausationIDHeader, core.GetCausationID(ctx))
}
