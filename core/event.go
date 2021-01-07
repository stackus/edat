package core

import (
	"fmt"
	"reflect"
)

// Event interface
type Event interface {
	EventName() string
}

// SerializeEvent serializes events with a registered marshaller
func SerializeEvent(v Event) ([]byte, error) {
	return marshal(v.EventName(), v)
}

// DeserializeEvent deserializes the event data using a registered marshaller returning an *Event
func DeserializeEvent(eventName string, data []byte) (Event, error) {
	evt, err := unmarshal(eventName, data)
	if err != nil {
		return nil, err
	}

	if evt != nil {
		if _, ok := evt.(Event); !ok {
			return nil, fmt.Errorf("`%s` was registered but not registered as an event", eventName)
		}
	}

	return evt.(Event), nil
}

// RegisterEvents registers one or more events with a registered marshaller
//
// Register events using any form desired "&MyEvent{}", "MyEvent{}", "(*MyEvent)(nil)"
//
// Events must be registered after first registering a marshaller you wish to use
func RegisterEvents(events ...Event) {
	for _, event := range events {
		if v := reflect.ValueOf(event); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
			eventName := reflect.Zero(reflect.TypeOf(event).Elem()).Interface().(Event).EventName()
			registerType(eventName, event)
		} else {
			registerType(event.EventName(), event)
		}
	}
}
