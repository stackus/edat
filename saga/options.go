package saga

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

func (o LoggerOption) configureOrchestrator(orchestrator *Orchestrator) {
	orchestrator.logger = o
}
