package core

import (
	"fmt"
	"reflect"
	"sync"
)

// Marshaller provides marshaling functions and type tracking capabilities
//
// This is how the library avoids requiring boilerplate written to convert each data type to
// and from a marshalled form.
//
// An example marshaller that uses gogoproto with "gogoslick_out" code generation. Adding
// a protobuf based Marshaller will result in significant speed improvements at the
// expense of having to maintain generated code.
//  // Define a marshaller
//  type MyProtoMarshaller struct{}
//  func (MyProtoMarshaller) Marshal(v interface{}) ([]byte, error) { return proto.Marshal(v.(proto.Message))}
//  func (MyProtoMarshaller) Unmarshal(data []byte, v interface{}) error { return proto.Unmarshal(data, v.(proto.Message))}
//  func (MyProtoMarshaller) GetType(typeName string) reflect.Type {
//  	t := proto.MessageType(typeName)
//  	if t != nil {
//  		return t.Elem()
//  	}
//  	return nil
//  }
//  func (ProtoMarshaller) RegisterType(string, reflect.Type) {}
//
//  // Register your marshaller and a function to test for the types it should be given to handle
//  core.RegisterMarshaller(MyProtoMarshaller{}, func(i interface{}) bool {
//  	_, ok := i.(proto.Message)
//  	return ok
//  })
type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	GetType(typeName string) reflect.Type
	RegisterType(typeName string, v reflect.Type)
}

type registeredMarshaller struct {
	marshaller Marshaller
	affinity   func(interface{}) bool
}

var registry = struct {
	defaultMarshaller Marshaller
	marshallers       []registeredMarshaller
	mu                sync.Mutex
}{
	marshallers: []registeredMarshaller{},
	mu:          sync.Mutex{},
}

func registerType(typeName string, v interface{}) {
	marshaller := registry.defaultMarshaller

	for _, s := range registry.marshallers {
		if s.affinity(v) {
			marshaller = s.marshaller
			break
		}
	}

	if marshaller == nil {
		panic("no marshallers have been set")
	}

	var t reflect.Type

	if value := reflect.ValueOf(v); value.Kind() == reflect.Ptr && value.Pointer() == 0 {
		t = reflect.TypeOf(v).Elem()
	} else {
		t = reflect.TypeOf(v)

		if value.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	marshaller.RegisterType(typeName, t)
}

func marshal(typeName string, v interface{}) ([]byte, error) {
	var t reflect.Type

	marshaller := registry.defaultMarshaller

	if marshaller != nil {
		t = marshaller.GetType(typeName)
	}

	if marshaller == nil || marshaller.GetType(typeName) == nil {
		for _, s := range registry.marshallers {
			if t = s.marshaller.GetType(typeName); t != nil {
				marshaller = s.marshaller
				break
			}
		}
	}

	if marshaller == nil || t == nil {
		return nil, fmt.Errorf("`%s` was not registered with any marshaller", typeName)
	}

	return marshaller.Marshal(v)
}

func unmarshal(typeName string, data []byte) (interface{}, error) {
	var t reflect.Type

	marshaller := registry.defaultMarshaller

	if marshaller != nil {
		t = marshaller.GetType(typeName)
	}

	if t == nil {
		for _, s := range registry.marshallers {
			if t = s.marshaller.GetType(typeName); t != nil {
				marshaller = s.marshaller
				break
			}
		}
	}

	if marshaller == nil || t == nil {
		return nil, fmt.Errorf("`%s` was not registered with any marshaller", typeName)
	}

	dst := reflect.New(t).Interface()

	err := marshaller.Unmarshal(data, dst)
	return dst, err
}

// RegisterMarshaller allows applications to register a new optimized marshaller for specific types or situations
func RegisterMarshaller(marshaller Marshaller, affinityFn func(interface{}) bool) {
	registerMarshaller(marshaller, affinityFn, false)
}

// RegisterDefaultMarshaller registers a marshaller to be used when no other marshaller should be used
func RegisterDefaultMarshaller(marshaller Marshaller) {
	registerMarshaller(marshaller, nil, true)
}

func registerMarshaller(marshaller Marshaller, affinityFn func(interface{}) bool, asDefault bool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if asDefault {
		registry.defaultMarshaller = marshaller
		return
	}

	rm := registeredMarshaller{
		marshaller: marshaller,
		affinity:   affinityFn,
	}

	registry.marshallers = append(registry.marshallers, rm)
}
