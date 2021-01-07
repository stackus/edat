package core

type Entity interface {
	ID() string
	EntityName() string
	Events() []Event
	AddEvent(events ...Event)
	ClearEvents()
}

// EntityBase provides entities a base to build on
type EntityBase struct {
	events []Event
}

func (e *EntityBase) Events() []Event {
	return e.events
}

func (e *EntityBase) AddEvent(events ...Event) {
	e.events = append(e.events, events...)
}

func (e *EntityBase) ClearEvents() {
	e.events = []Event{}
}
