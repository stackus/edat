package es

// AggregateRootOption options for AggregateRoots
type AggregateRootOption func(r *AggregateRoot)

// WithAggregateRootID is an option to set the ID of the AggregateRoot
func WithAggregateRootID(aggregateID string) AggregateRootOption {
	return func(r *AggregateRoot) {
		r.Aggregate.setID(aggregateID)
	}
}
