package core

// Entity have identity and change tracking in the form of events
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

// Events returns all tracked changes made to the Entity as Events
func (e *EntityBase) Events() []Event {
	return e.events
}

// AddEvent adds a tracked change to the Entity
func (e *EntityBase) AddEvent(events ...Event) {
	e.events = append(e.events, events...)
}

// ClearEvents resets the tracked change list
func (e *EntityBase) ClearEvents() {
	e.events = []Event{}
}
