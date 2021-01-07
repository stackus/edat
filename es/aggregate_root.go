package es

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/stackus/edat/core"
)

const aggregateNeverCommitted = 0

var ErrPendingChanges = fmt.Errorf("cannot process command while pending changes exist")

// AggregateRoot is the base for Aggregates
type AggregateRoot struct {
	aggregate Aggregate
	version   int
}

// NewAggregateRoot constructor for *AggregateRoot
func NewAggregateRoot(aggregate Aggregate, options ...AggregateRootOption) *AggregateRoot {
	r := &AggregateRoot{
		aggregate: aggregate,
		version:   aggregateNeverCommitted,
	}

	for _, option := range options {
		option(r)
	}

	if r.aggregate.ID() == "" {
		r.aggregate.setID(uuid.New().String())
	}

	return r
}

// ID returns the ID for the root aggregate
func (r AggregateRoot) ID() string {
	return r.aggregate.ID()
}

// AggregateID returns the ID for the root aggregate
func (r AggregateRoot) AggregateID() string {
	return r.aggregate.ID()
}

// EntityName returns the Name of the root aggregate
func (r AggregateRoot) EntityName() string {
	return r.aggregate.EntityName()
}

// AggregateName returns the Name of the root aggregate
func (r AggregateRoot) AggregateName() string {
	return r.aggregate.EntityName()
}

// Aggregate returns the aggregate that resides at the root
func (r AggregateRoot) Aggregate() Aggregate {
	return r.aggregate
}

// PendingVersion is the version of the aggregate taking into account pending events
func (r AggregateRoot) PendingVersion() int {
	return r.version + len(r.aggregate.Events())
}

// Version is the version of the aggregate as it was created or loaded
func (r AggregateRoot) Version() int {
	return r.version
}

// ProcessCommand runs the command and records the changes as pending events or returns an error
func (r *AggregateRoot) ProcessCommand(command core.Command) error {
	if len(r.aggregate.Events()) != 0 {
		return ErrPendingChanges
	}

	err := r.aggregate.ProcessCommand(command)
	if err != nil {
		return err
	}

	for _, event := range r.aggregate.Events() {
		aErr := r.aggregate.ApplyEvent(event)
		if aErr != nil {
			return aErr
		}
	}

	return nil
}

// Events returns the list of pending events
func (r AggregateRoot) Events() []core.Event {
	return r.aggregate.Events()
}

// AddEvent stores entity events on the aggregate
func (r *AggregateRoot) AddEvent(events ...core.Event) {
	r.aggregate.AddEvent(events...)
}

// ClearEvents clears any pending events without committing them
func (r *AggregateRoot) ClearEvents() {
	r.aggregate.ClearEvents()
}

// CommitEvents clears any pending events and updates the last committed version value
func (r *AggregateRoot) CommitEvents() {
	r.version += len(r.aggregate.Events())
	r.aggregate.ClearEvents()
}

// LoadEvent is used to rerun events essentially left folding over the aggregate state
func (r *AggregateRoot) LoadEvent(events ...core.Event) error {
	for _, event := range events {
		err := r.aggregate.ApplyEvent(event)
		if err != nil {
			return err
		}
	}

	r.version += len(events)

	return nil
}

// LoadSnapshot is used to apply a snapshot to the aggregate to save having to rerun all events
func (r *AggregateRoot) LoadSnapshot(snapshot core.Snapshot, version int) error {
	err := r.aggregate.ApplySnapshot(snapshot)
	if err != nil {
		return err
	}

	r.version = version

	return nil
}
