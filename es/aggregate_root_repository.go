package es

import (
	"context"
	"errors"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
)

// AggregateRepository interface
type AggregateRepository interface {
	Load(ctx context.Context, aggregateID string) (*AggregateRoot, error)
	Save(ctx context.Context, command core.Command, options ...AggregateRootOption) (*AggregateRoot, error)
	Update(ctx context.Context, aggregateID string, command core.Command, options ...AggregateRootOption) (*AggregateRoot, error)
}

// AggregateRootRepository uses stores to load and save the changes to aggregates as events
type AggregateRootRepository struct {
	constructor func() Aggregate
	store       AggregateRootStore
	logger      log.Logger
}

// AggregateRootStoreMiddleware interface for embedating stores
type AggregateRootStoreMiddleware func(store AggregateRootStore) AggregateRootStore

// ErrAggregateNotFound is returned when no root was found for a given aggregate id
var ErrAggregateNotFound = errors.New("aggregate not found")

// NewAggregateRootRepository constructs a new AggregateRootRepository
func NewAggregateRootRepository(constructor func() Aggregate, store AggregateRootStore) *AggregateRootRepository {
	r := &AggregateRootRepository{
		constructor: constructor,
		store:       store,
		logger:      log.DefaultLogger,
	}

	r.logger.Trace("es.AggregateRootRepository constructed")

	return r
}

// Load finds aggregates in the provided store
func (r *AggregateRootRepository) Load(ctx context.Context, aggregateID string) (*AggregateRoot, error) {
	root := r.root(WithAggregateRootID(aggregateID))

	err := r.store.Load(ctx, root)
	if err != nil {
		return nil, err
	}

	if root.version == aggregateNeverCommitted {
		return nil, ErrAggregateNotFound
	}

	return root, r.store.Load(ctx, root)
}

// Save applies the given command to a new aggregate and persists it into the store
func (r *AggregateRootRepository) Save(ctx context.Context, command core.Command, options ...AggregateRootOption) (*AggregateRoot, error) {
	root := r.root(options...)

	return root, r.save(ctx, command, root)
}

// Update locates an existing aggregate, applies the commands and persists the result into the store
func (r *AggregateRootRepository) Update(ctx context.Context, aggregateID string, command core.Command, options ...AggregateRootOption) (*AggregateRoot, error) {
	root := r.root(append(options, WithAggregateRootID(aggregateID))...)

	err := r.store.Load(ctx, root)
	if err != nil {
		return nil, err
	}

	return root, r.save(ctx, command, root)
}

func (r *AggregateRootRepository) root(options ...AggregateRootOption) *AggregateRoot {
	return NewAggregateRoot(r.constructor(), options...)
}

func (r *AggregateRootRepository) save(ctx context.Context, command core.Command, root *AggregateRoot) error {
	err := root.ProcessCommand(command)
	if err != nil {
		return err
	}

	if root.PendingVersion() == root.Version() {
		return nil
	}

	err = r.store.Save(ctx, root)
	if err != nil {
		r.logger.Error("error saving aggregate root", log.Error(err))
		return err
	}

	return nil
}
