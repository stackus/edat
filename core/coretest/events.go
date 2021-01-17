package coretest

import (
	"github.com/stackus/edat/core/coremocks"
)

type (
	Event             struct{ Value string }
	UnregisteredEvent struct{ Value string }
)

func (Event) EventName() string             { return "coretest.Event" }
func (UnregisteredEvent) EventName() string { return "coretest.UnregisteredEvent" }

func MockEvent(setup func(m *coremocks.Event)) *coremocks.Event {
	m := &coremocks.Event{}
	setup(m)
	return m
}
