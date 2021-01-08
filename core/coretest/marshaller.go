package coretest

import (
	"encoding/json"
	"reflect"
	"sync"
)

type TestMarshaller struct {
	types map[string]reflect.Type
	mu    sync.Mutex
}

func NewTestMarshaller() *TestMarshaller {
	return &TestMarshaller{
		types: map[string]reflect.Type{},
		mu:    sync.Mutex{},
	}
}

func (*TestMarshaller) Marshal(v interface{}) ([]byte, error)      { return json.Marshal(v) }
func (*TestMarshaller) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }
func (m *TestMarshaller) GetType(typeName string) reflect.Type     { return m.types[typeName] }
func (m *TestMarshaller) RegisterType(typeName string, v reflect.Type) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.types[typeName] = v
}

func (m *TestMarshaller) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.types = map[string]reflect.Type{}
}
