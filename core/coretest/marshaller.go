package coretest

import (
	"encoding/json"
	"reflect"
	"sync"
)

// TestMarshaller returns a Marshaller for testing purposes
//
// JSON encoding is used for the Marshal and Unmarshal methods.
//
// The registered types may be reset using Reset() at any time while testing.
type TestMarshaller struct {
	types map[string]reflect.Type
	mu    sync.Mutex
}

// NewTestMarshaller constructs a new TestMarshaller
func NewTestMarshaller() *TestMarshaller {
	return &TestMarshaller{
		types: map[string]reflect.Type{},
		mu:    sync.Mutex{},
	}
}

// Marshal returns v in byte form
func (*TestMarshaller) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

// Unmarshal returns the bytes marshalled into v
func (*TestMarshaller) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

// GetType returns the reflect.Type if it has been registered
func (m *TestMarshaller) GetType(typeName string) reflect.Type { return m.types[typeName] }

// RegisterType registers a new reflect.Type for the given name key
func (m *TestMarshaller) RegisterType(typeName string, v reflect.Type) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.types[typeName] = v
}

// Reset will remove all previously registered types
func (m *TestMarshaller) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.types = map[string]reflect.Type{}
}
