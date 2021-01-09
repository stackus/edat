package es

// DefaultSnapshotStrategy is a strategy that triggers snapshots every 10 changes
var DefaultSnapshotStrategy = NewMaxChangesSnapshotStrategy(10)

// SnapshotStrategy interface
type SnapshotStrategy interface {
	ShouldSnapshot(aggregate *AggregateRoot) bool
}

type maxChangesSnapshotStrategy struct {
	maxChanges int
}

// NewMaxChangesSnapshotStrategy constructs a new SnapshotStrategy with "max changes" rules
func NewMaxChangesSnapshotStrategy(maxChanges int) SnapshotStrategy {
	return &maxChangesSnapshotStrategy{maxChanges: maxChanges}
}

// ShouldSnapshot implements es.SnapshotStrategy.ShouldSnapshot
func (s *maxChangesSnapshotStrategy) ShouldSnapshot(aggregate *AggregateRoot) bool {
	return aggregate.PendingVersion() >= s.maxChanges && ((len(aggregate.Events()) >= s.maxChanges) ||
		(aggregate.PendingVersion()%s.maxChanges < len(aggregate.Events())) ||
		(aggregate.PendingVersion()%s.maxChanges == 0))
}
