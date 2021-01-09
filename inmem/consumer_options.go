package inmem

import (
	"github.com/stackus/edat/log"
)

// ConsumerOption options for Consumer
type ConsumerOption func(*Consumer)

// WithConsumerLogger sets the log.Logger for Consumer
func WithConsumerLogger(logger log.Logger) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.logger = logger
	}
}
