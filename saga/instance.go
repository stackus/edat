package saga

import (
	"github.com/stackus/edat/core"
)

// Instance is the container for saga data
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

// SagaID returns the instance saga id
func (i *Instance) SagaID() string {
	return i.sagaID
}

// SagaName returns the instance saga name
func (i *Instance) SagaName() string {
	return i.sagaName
}

// SagaData returns the instance saga data
func (i *Instance) SagaData() core.SagaData {
	return i.sagaData
}

// CurrentStep returns the step currently being processed
func (i *Instance) CurrentStep() int {
	return i.currentStep
}

// EndState returns whether or not all steps have completed
func (i *Instance) EndState() bool {
	return i.endState
}

// Compensating returns whether or not the instance is compensating (rolling back)
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
