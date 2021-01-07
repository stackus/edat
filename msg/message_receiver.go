package msg

import (
	"context"
)

// MessageReceiver interface for channel subscription receivers
type MessageReceiver interface {
	ReceiveMessage(context.Context, Message) error
}

// ReceiveMessageFunc makes it easy to drop in functions as receivers
type ReceiveMessageFunc func(context.Context, Message) error

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (f ReceiveMessageFunc) ReceiveMessage(ctx context.Context, message Message) error {
	return f(ctx, message)
}
