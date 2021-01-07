package es

import (
	"context"
)

// AggregateRootStore is the interface that infrastructures should implement to be used in AggregateRootRepositories
type AggregateRootStore interface {
	Load(ctx context.Context, aggregate *AggregateRoot) error
	Save(ctx context.Context, aggregate *AggregateRoot) error
}
