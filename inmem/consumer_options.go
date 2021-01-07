package inmem

import (
	"github.com/stackus/edat/log"
)

type ConsumerOption func(*Consumer)

func WithConsumerLogger(logger log.Logger) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.logger = logger
	}
}
