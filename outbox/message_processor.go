package outbox

import (
	"context"
)

type MessageProcessor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
