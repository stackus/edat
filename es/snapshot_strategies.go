package es

var DefaultSnapshotStrategy = NewMaxChangesSnapshotStrategy(10)

type SnapshotStrategy interface {
	ShouldSnapshot(aggregate *AggregateRoot) bool
}

type maxChangesSnapshotStrategy struct {
	maxChanges int
}

func NewMaxChangesSnapshotStrategy(maxChanges int) SnapshotStrategy {
	return &maxChangesSnapshotStrategy{maxChanges: maxChanges}
}

func (s *maxChangesSnapshotStrategy) ShouldSnapshot(aggregate *AggregateRoot) bool {
	return aggregate.PendingVersion() >= s.maxChanges && ((len(aggregate.Events()) >= s.maxChanges) ||
		(aggregate.PendingVersion()%s.maxChanges < len(aggregate.Events())) ||
		(aggregate.PendingVersion()%s.maxChanges == 0))
}
