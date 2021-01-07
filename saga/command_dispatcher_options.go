package saga

import (
	"github.com/stackus/edat/log"
)

// CommandDispatcherOption options for CommandConsumers
type CommandDispatcherOption func(consumer *CommandDispatcher)

// WithCommandDispatcherLogger is an option to set the log.Logger of the CommandDispatcher
func WithCommandDispatcherLogger(logger log.Logger) CommandDispatcherOption {
	return func(dispatcher *CommandDispatcher) {
		dispatcher.logger = logger
	}
}
