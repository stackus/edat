package saga

import (
	"github.com/stackus/edat/core"
)

type Instance struct {
	sagaID       string
	sagaName     string
	sagaData     core.SagaData
	currentStep  int
	endState     bool
	compensating bool
}

// NewSagaInstance constructor for *SagaInstances
func NewSagaInstance(sagaName, sagaID string, sagaData core.SagaData, currentStep int, endState, compensating bool) *Instance {
	return &Instance{
		sagaID:       sagaID,
		sagaName:     sagaName,
		sagaData:     sagaData,
		currentStep:  currentStep,
		endState:     endState,
		compensating: compensating,
	}
}

func (i *Instance) SagaID() string {
	return i.sagaID
}

func (i *Instance) SagaName() string {
	return i.sagaName
}

func (i *Instance) SagaData() core.SagaData {
	return i.sagaData
}

func (i *Instance) CurrentStep() int {
	return i.currentStep
}

func (i *Instance) EndState() bool {
	return i.endState
}

func (i *Instance) Compensating() bool {
	return i.compensating
}

func (i *Instance) getStepContext() stepContext {
	return stepContext{
		step:         i.currentStep,
		compensating: i.compensating,
		ended:        i.endState,
	}
}

func (i *Instance) updateStepContext(stepCtx stepContext) {
	i.currentStep = stepCtx.step
	i.endState = stepCtx.ended
	i.compensating = stepCtx.compensating
}
