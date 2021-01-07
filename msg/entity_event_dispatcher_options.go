package msg

import (
	"github.com/stackus/edat/log"
)

// EntityEventDispatcherOption options for EntityEventDispatcher
type EntityEventDispatcherOption func(consumer *EntityEventDispatcher)

// WithEntityEventDispatcherLogger is an option to set the log.Logger of the EntityEventDispatcher
func WithEntityEventDispatcherLogger(logger log.Logger) EntityEventDispatcherOption {
	return func(dispatcher *EntityEventDispatcher) {
		dispatcher.logger = logger
	}
}
