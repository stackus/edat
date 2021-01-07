package msg

import (
	"context"
)

// Producer is the interface that infrastructures should implement to be used in a Publisher
type Producer interface {
	Send(ctx context.Context, channel string, message Message) error
	Close(ctx context.Context) error
}
