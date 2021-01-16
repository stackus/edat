// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package msgmocks

import (
	context "context"

	core "github.com/stackus/edat/core"
	mock "github.com/stretchr/testify/mock"

	msg "github.com/stackus/edat/msg"
)

// EventMessagePublisher is an autogenerated mock type for the EventMessagePublisher type
type EventMessagePublisher struct {
	mock.Mock
}

// PublishEvent provides a mock function with given fields: ctx, event, options
func (_m *EventMessagePublisher) PublishEvent(ctx context.Context, event core.Event, options ...msg.MessageOption) error {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, event)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Event, ...msg.MessageOption) error); ok {
		r0 = rf(ctx, event, options...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
