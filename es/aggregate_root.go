package es

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/stackus/edat/core"
)

const aggregateNeverCommitted = 0

// ErrPendingChanges is the error returned when a second command is applied to an aggregate without clearing changes
var ErrPendingChanges = fmt.Errorf("cannot process command while pending changes exist")

// AggregateRoot is the base for Aggregates
type AggregateRoot struct {
	Aggregate
	version int
}

// NewAggregateRoot constructor for *AggregateRoot
func NewAggregateRoot(aggregate Aggregate, options ...AggregateRootOption) *AggregateRoot {
	r := &AggregateRoot{
		Aggregate: aggregate,
		version:   aggregateNeverCommitted,
	}

	for _, option := range options {
		option(r)
	}

	if r.Aggregate.ID() == "" {
		r.Aggregate.setID(uuid.New().String())
	}

	return r
}

// AggregateName returns the Name of the root aggregate
func (r AggregateRoot) AggregateName() string {
	return r.Aggregate.EntityName()
}

// ProcessCommand runs the command and records the changes as pending events or returns an error
func (r *AggregateRoot) ProcessCommand(command core.Command) error {
	if len(r.Aggregate.Events()) != 0 {
		return ErrPendingChanges
	}

	err := r.Aggregate.ProcessCommand(command)
	if err != nil {
		return err
	}

	for _, event := range r.Aggregate.Events() {
		aErr := r.Aggregate.ApplyEvent(event)
		if aErr != nil {
			return aErr
		}
	}

	return nil
}

// LoadEvent is used to rerun events essentially left folding over the aggregate state
func (r *AggregateRoot) LoadEvent(events ...core.Event) error {
	for _, event := range events {
		err := r.Aggregate.ApplyEvent(event)
		if err != nil {
			return err
		}
	}

	r.version += len(events)

	return nil
}

// LoadSnapshot is used to apply a snapshot to the aggregate to save having to rerun all events
func (r *AggregateRoot) LoadSnapshot(snapshot core.Snapshot, version int) error {
	err := r.Aggregate.ApplySnapshot(snapshot)
	if err != nil {
		return err
	}

	r.version = version

	return nil
}
