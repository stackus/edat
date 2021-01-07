package inmem

import (
	"github.com/stackus/edat/es"
)

type SnapshotStoreOption func(store *SnapshotStore)

func SnapshotStoreStrategy(strategy es.SnapshotStrategy) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.strategy = strategy
	}
}
