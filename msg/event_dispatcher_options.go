package msg

import (
	"github.com/stackus/edat/log"
)

// EventDispatcherOption options for EventDispatcher
type EventDispatcherOption func(consumer *EventDispatcher)

// WithEventDispatcherLogger is an option to set the log.Logger of the EventDispatcher
func WithEventDispatcherLogger(logger log.Logger) EventDispatcherOption {
	return func(dispatcher *EventDispatcher) {
		dispatcher.logger = logger
	}
}
