package msg

import (
	"github.com/stackus/edat/core"
)

// DomainEvent interface for events that are shared across the domain
type DomainEvent interface {
	core.Event
	DestinationChannel() string
}

// Event is an event with message header information
type Event interface {
	Event() core.Event
	Headers() Headers
}

type eventMessage struct {
	event   core.Event
	headers Headers
}

func (m eventMessage) Event() core.Event {
	return m.event
}

func (m eventMessage) Headers() Headers {
	return m.headers
}
