package inmem

import (
	"github.com/stackus/edat/log"
)

// ProducerOption options for Producer
type ProducerOption func(destination *Producer)

// WithProducerLogger sets the log.Logger for Producer
func WithProducerLogger(logger log.Logger) ProducerOption {
	return func(producer *Producer) {
		producer.logger = logger
	}
}
