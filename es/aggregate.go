package es

import (
	"github.com/stackus/edat/core"
)

// Aggregate is a domain object that will be used at the base of other domain objects
type Aggregate interface {
	core.Entity
	setID(id string)
	ProcessCommand(command core.Command) error
	ApplyEvent(event core.Event) error
	ApplySnapshot(snapshot core.Snapshot) error
	ToSnapshot() (core.Snapshot, error)

	AggregateID() string
	Version() int
	PendingVersion() int
	CommitEvents()
}

// AggregateBase provides aggregates a base to build on
type AggregateBase struct {
	core.EntityBase
	id      string
	version int
}

// ID returns to the immutable ID of the aggregate
func (a AggregateBase) ID() string {
	return a.id
}

// AggregateID returns the ID for the root aggregate
func (a AggregateBase) AggregateID() string {
	return a.id
}

// setID is used internally by the AggregateRoot to apply IDs when loading existing Aggregates
// nolint unused
func (a *AggregateBase) setID(id string) {
	a.id = id
}

// Version is the version of the aggregate as it was created or loaded
func (a *AggregateBase) Version() int {
	return a.version
}

// PendingVersion is the version of the aggregate taking into account pending events
func (a *AggregateBase) PendingVersion() int {
	return a.version + len(a.Events())
}

// CommitEvents clears any pending events and updates the last committed version value
func (a *AggregateBase) CommitEvents() {
	a.version += len(a.Events())
	a.ClearEvents()
}
