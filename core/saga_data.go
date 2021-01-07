package core

import (
	"fmt"
	"reflect"
)

// SagaData interface
type SagaData interface {
	SagaDataName() string
}

// SerializeSagaData serializes saga data with a registered marshaller
func SerializeSagaData(v SagaData) ([]byte, error) {
	return marshal(v.SagaDataName(), v)
}

// DeserializeSagaData deserializes the saga data data using a registered marshaller returning a *SagaData
func DeserializeSagaData(sagaDataName string, data []byte) (SagaData, error) {
	sagaData, err := unmarshal(sagaDataName, data)
	if err != nil {
		return nil, err
	}

	if sagaData != nil {
		if _, ok := sagaData.(SagaData); !ok {
			return nil, fmt.Errorf("`%s` was registered but not registered as a saga data", sagaDataName)
		}
	}

	return sagaData.(SagaData), nil
}

// RegisterSagaData registers one or more saga data with a registered marshaller
//
// Register saga data using any form desired "&MySagaData{}", "MySagaData{}", "(*MySagaData)(nil)"
//
// SagaData must be registered after first registering a marshaller you wish to use
func RegisterSagaData(sagaDatas ...SagaData) {
	for _, sagaData := range sagaDatas {
		if v := reflect.ValueOf(sagaData); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
			sagaDataName := reflect.Zero(reflect.TypeOf(sagaData).Elem()).Interface().(SagaData).SagaDataName()
			registerType(sagaDataName, sagaData)
		} else {
			registerType(sagaData.SagaDataName(), sagaData)
		}
	}
}
