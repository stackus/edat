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
}

// AggregateBase provides aggregates a base to build on
type AggregateBase struct {
	core.EntityBase
	id string
}

// ID returns to the immutable ID of the aggregate
func (a AggregateBase) ID() string {
	return a.id
}

// setID is used internally by the AggregateRoot to apply IDs when loading existing Aggregates
// nolint unused
func (a *AggregateBase) setID(id string) {
	a.id = id
}
