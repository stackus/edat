package saga

import (
	"context"

	"github.com/stackus/edat/core"
)

type RemoteStepActionOption func(action *remoteStepAction)

func WithRemoteStepPredicate(predicate func(context.Context, core.SagaData) bool) RemoteStepActionOption {
	return func(step *remoteStepAction) {
		step.predicate = predicate
	}
}
