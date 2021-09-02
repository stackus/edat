package msg

import (
	"github.com/stackus/edat/log"
)

type LoggerOption struct {
	log.Logger
}

func WithLogger(logger log.Logger) LoggerOption {
	return LoggerOption{logger}
}

func (o LoggerOption) configureCommandDispatcher(d *CommandDispatcher) {
	d.logger = o
}

func (o LoggerOption) configureEntityEventDispatcher(d *EntityEventDispatcher) {
	d.logger = o
}

func (o LoggerOption) configureEventDispatcher(d *EventDispatcher) {
	d.logger = o
}

func (o LoggerOption) configurePublisher(p *Publisher) {
	p.logger = o
}

func (o LoggerOption) configureSubscriber(s *Subscriber) {
	s.logger = o
}
