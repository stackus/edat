package saga

import (
	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
)

type stepResults struct {
	commands           []msg.DomainCommand
	updatedSagaData    core.SagaData
	updatedStepContext stepContext
	local              bool
	failure            error
}
