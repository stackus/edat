package inmem

import (
	"github.com/stackus/edat/log"
)

type LoggerOption struct {
	log.Logger
}

func WithLogger(s log.Logger) LoggerOption {
	return LoggerOption{s}
}

func (o LoggerOption) configureConsumer(c *Consumer) {
	c.logger = o
}

func (o LoggerOption) configureProducer(p *Producer) {
	p.logger = o
}
