package outbox

import (
	"context"
)

// MessageProcessor interface
type MessageProcessor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
