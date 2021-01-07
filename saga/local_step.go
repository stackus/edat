package saga

import (
	"context"

	"github.com/stackus/edat/core"
)

// LocalStep is used to execute local saga business logic
type LocalStep struct {
	actions map[bool]func(context.Context, core.SagaData) error
}

var _ Step = (*LocalStep)(nil)

// NewLocalStep constructor for LocalStep
func NewLocalStep(action func(context.Context, core.SagaData) error) LocalStep {
	return LocalStep{
		actions: map[bool]func(context.Context, core.SagaData) error{
			notCompensating: action,
		},
	}
}

func (s LocalStep) Compensation(compensation func(context.Context, core.SagaData) error) LocalStep {
	s.actions[isCompensating] = compensation
	return s
}

func (s LocalStep) hasInvocableAction(_ context.Context, _ core.SagaData, compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s LocalStep) getReplyHandler(string, bool) func(context.Context, core.SagaData, core.Reply) error {
	return nil
}

func (s LocalStep) execute(ctx context.Context, sagaData core.SagaData, compensating bool) func(results *stepResults) {
	err := s.actions[compensating](ctx, sagaData)
	return func(results *stepResults) {
		results.local = true
		results.failure = err
	}
}
