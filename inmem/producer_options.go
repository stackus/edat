package inmem

import (
	"github.com/stackus/edat/log"
)

type ProducerOption func(destination *Producer)

func WithProducerLogger(logger log.Logger) ProducerOption {
	return func(producer *Producer) {
		producer.logger = logger
	}
}
