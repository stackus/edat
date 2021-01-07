package msg

import (
	"github.com/stackus/edat/core"
)

// EntityEvent is an event with message header information
type EntityEvent interface {
	EntityID() string
	EntityName() string
	Event() core.Event
	Headers() Headers
}

type entityEventMessage struct {
	entityID   string
	entityName string
	event      core.Event
	headers    Headers
}

func (m entityEventMessage) EntityID() string {
	return m.entityID
}

func (m entityEventMessage) EntityName() string {
	return m.entityName
}

func (m entityEventMessage) Event() core.Event {
	return m.event
}

func (m entityEventMessage) Headers() Headers {
	return m.headers
}
