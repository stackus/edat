package core

import (
	"context"
)

type contextKey int

// Known contextKey key values
const (
	requestIDKey contextKey = iota + 1
	correlationIDKey
	causationIDKey
)

// GetRequestID returns the RequestID from the context or a blank if not set
func GetRequestID(ctx context.Context) string {
	requestID := ctx.Value(requestIDKey)
	if requestID == nil {
		return ""
	}

	return requestID.(string)
}

// GetCorrelationID returns the CorrelationID from the context or a blank if not set
//
// In a long line of events, commands and messages this ID will match the original RequestID
func GetCorrelationID(ctx context.Context) string {
	correlationID := ctx.Value(correlationIDKey)
	if correlationID == nil {
		return GetRequestID(ctx)
	}

	return correlationID.(string)
}

// GetCausationID returns the CausationID from the context or a blank if not set
//
// In a long line of events, commands and messages this ID will match the previous RequestID
func GetCausationID(ctx context.Context) string {
	causationID := ctx.Value(causationIDKey)
	if causationID == nil {
		return GetRequestID(ctx)
	}

	return causationID.(string)
}

// SetRequestContext sets the Request, Correlation, and Causation IDs on the context
//
// Correlation and Causation IDs will use the RequestID if blank ID values are provided
func SetRequestContext(ctx context.Context, requestID, correlationID, causationID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	// CorrelationIDs point back to the first request
	if correlationID == "" {
		ctx = context.WithValue(ctx, correlationIDKey, requestID)
	} else {
		ctx = context.WithValue(ctx, correlationIDKey, correlationID)
	}

	// CausationIDs point back to the previous request
	if causationID == "" {
		ctx = context.WithValue(ctx, causationIDKey, requestID)
	} else {
		ctx = context.WithValue(ctx, causationIDKey, causationID)
	}
	return ctx
}
