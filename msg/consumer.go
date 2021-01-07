package msg

import (
	"context"
)

// Consumer is the interface that infrastructures should implement to be used in MessageDispatchers
type Consumer interface {
	Listen(ctx context.Context, channel string, consumer ReceiveMessageFunc) error
	Close(ctx context.Context) error
}
