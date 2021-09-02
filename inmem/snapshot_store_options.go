package inmem

import (
	"github.com/stackus/edat/es"
)

// SnapshotStoreOption options for SnapshotStore
type SnapshotStoreOption interface {
	configureSnapshotStore(*SnapshotStore)
}

type StrategyOption struct {
	es.SnapshotStrategy
}

// WithStrategy sets the snapshotting strategy for SnapshotStore
func WithStrategy(strategy es.SnapshotStrategy) SnapshotStoreOption {
	return StrategyOption{strategy}
}

func (o StrategyOption) configureSnapshotStore(store *SnapshotStore) {
	store.strategy = o
}
