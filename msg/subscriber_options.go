package msg

import (
	"github.com/stackus/edat/log"
)

// SubscriberOption options for MessageConsumers
type SubscriberOption func(*Subscriber)

// WithSubscriberLogger is an option to set the log.Logger of the Subscriber
func WithSubscriberLogger(logger log.Logger) SubscriberOption {
	return func(subscriber *Subscriber) {
		subscriber.logger = logger
	}
}
