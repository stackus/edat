package saga

import (
	"github.com/stackus/edat/log"
)

type OrchestratorOption func(o *Orchestrator)

func WithOrchestratorLogger(logger log.Logger) OrchestratorOption {
	return func(o *Orchestrator) {
		o.logger = logger
	}
}
