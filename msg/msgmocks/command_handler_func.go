// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package msgmocks

import (
	context "context"

	msg "github.com/stackus/edat/msg"
	mock "github.com/stretchr/testify/mock"
)

// CommandHandlerFunc is an autogenerated mock type for the CommandHandlerFunc type
type CommandHandlerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *CommandHandlerFunc) Execute(_a0 context.Context, _a1 msg.Command) ([]msg.Reply, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []msg.Reply
	if rf, ok := ret.Get(0).(func(context.Context, msg.Command) []msg.Reply); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msg.Reply)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, msg.Command) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
