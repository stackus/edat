package saga

import (
	"github.com/stackus/edat/log"
)

// OrchestratorOption options for Orchestrator
type OrchestratorOption func(o *Orchestrator)

// WithOrchestratorLogger is an option to set the log.Logger of the Orchestrator
func WithOrchestratorLogger(logger log.Logger) OrchestratorOption {
	return func(o *Orchestrator) {
		o.logger = logger
	}
}
