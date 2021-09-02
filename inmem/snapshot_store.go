package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
)

// SnapshotStore implements es.AggregateRootStore
type SnapshotStore struct {
	strategy  es.SnapshotStrategy
	snapshots sync.Map
	next      es.AggregateRootStore
}

type snapshotMsg struct {
	name     string
	version  int
	snapshot json.RawMessage
}

var _ es.AggregateRootStore = (*SnapshotStore)(nil)

// NewSnapshotStore constructs a new SnapshotStore and returns es.AggregateRootStoreMiddleware
func NewSnapshotStore(options ...SnapshotStoreOption) es.AggregateRootStoreMiddleware {
	s := &SnapshotStore{
		strategy:  es.DefaultSnapshotStrategy,
		snapshots: sync.Map{},
	}

	for _, option := range options {
		option.configureSnapshotStore(s)
	}

	return func(next es.AggregateRootStore) es.AggregateRootStore {
		s.next = next
		return s
	}
}

// Load implements es.AggregateRootStore.Load
func (s *SnapshotStore) Load(ctx context.Context, root *es.AggregateRoot) error {
	name := root.AggregateName()
	id := root.AggregateID()

	if result, exists := s.snapshots.Load(s.streamID(name, id)); exists {
		message := result.(snapshotMsg)
		snapshot, err := core.DeserializeSnapshot(message.name, message.snapshot)
		if err != nil {
			return err
		}

		err = root.LoadSnapshot(snapshot, message.version)
		if err != nil {
			return err
		}
	}

	return s.next.Load(ctx, root)
}

// Save implements es.AggregateRootStore.Save
func (s *SnapshotStore) Save(ctx context.Context, root *es.AggregateRoot) error {
	err := s.next.Save(ctx, root)
	if err != nil {
		return err
	}

	if !s.strategy.ShouldSnapshot(root) {
		return nil
	}

	snapshot, err := root.ToSnapshot()
	if err != nil {
		return err
	}

	data, err := core.SerializeSnapshot(snapshot)
	if err != nil {
		return err
	}

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()

	s.snapshots.Store(s.streamID(name, id), snapshotMsg{
		name:     snapshot.SnapshotName(),
		version:  version,
		snapshot: data,
	})

	return nil
}

func (s *SnapshotStore) streamID(name, id string) string {
	return fmt.Sprintf("%s:%s", name, id)
}
