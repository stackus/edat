package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/saga"
)

type SagaInstanceStore struct {
	instances sync.Map
}

type instanceData struct {
	sagaID       string
	sagaName     string
	sagaData     core.SagaData
	currentStep  int
	endState     bool
	compensating bool
}

var _ saga.InstanceStore = (*SagaInstanceStore)(nil)

func NewSagaInstanceStore() *SagaInstanceStore {
	return &SagaInstanceStore{
		instances: sync.Map{},
	}
}

func (s *SagaInstanceStore) Find(_ context.Context, sagaName, sagaID string) (*saga.Instance, error) {
	if dataT, exists := s.instances.Load(s.instanceID(sagaName, sagaID)); exists {
		data := dataT.(instanceData)

		instance := saga.NewSagaInstance(data.sagaName, sagaID, data.sagaData, data.currentStep, data.endState, data.compensating)

		return instance, nil
	}

	return nil, nil
}

func (s *SagaInstanceStore) Save(_ context.Context, instance *saga.Instance) error {
	return s.save(instance)
}

func (s *SagaInstanceStore) Update(_ context.Context, instance *saga.Instance) error {
	return s.save(instance)
}

func (s *SagaInstanceStore) save(instance *saga.Instance) error {
	instanceID := s.instanceID(instance.SagaName(), instance.SagaID())

	s.instances.Store(instanceID, instanceData{
		sagaID:       instance.SagaID(),
		sagaName:     instance.SagaName(),
		sagaData:     instance.SagaData(),
		currentStep:  instance.CurrentStep(),
		endState:     instance.EndState(),
		compensating: instance.Compensating(),
	})

	return nil
}

func (s *SagaInstanceStore) instanceID(sagaName, sagaID string) string {
	return fmt.Sprintf("%s:%s", sagaName, sagaID)
}
