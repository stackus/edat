package msg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/log/logmocks"
	"github.com/stackus/edat/log/logtest"
	"github.com/stackus/edat/msg"
)

func TestEntityEventDispatcher_ReceiveMessage(t *testing.T) {
	type handler struct {
		evt core.Event
		fn  msg.EntityEventHandlerFunc
	}
	type fields struct {
		handlers []handler
		logger   log.Logger
	}
	type args struct {
		ctx     context.Context
		message msg.Message
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterEvents(coretest.Event{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.Event{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventName:       coretest.Event{}.EventName(),
					msg.MessageEventEntityName: "entity-name",
					msg.MessageEventEntityID:   "entity-id",
				})),
			},
			wantErr: false,
		},
		"HandlerError": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.Event{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return fmt.Errorf("handler error")
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "entity event handler returned an error", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventName:       coretest.Event{}.EventName(),
					msg.MessageEventEntityName: "entity-name",
					msg.MessageEventEntityID:   "entity-id",
				})),
			},
			wantErr: true,
		},
		"coretest.UnregisteredEvent": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.UnregisteredEvent{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error decoding entity event message payload", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventName:       coretest.UnregisteredEvent{}.EventName(),
					msg.MessageEventEntityName: "entity-name",
					msg.MessageEventEntityID:   "entity-id",
				})),
			},
			wantErr: false,
		},
		"MissingEventName": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.Event{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error reading event name", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventEntityName: "entity-name",
					msg.MessageEventEntityID:   "entity-id",
				})),
			},
			wantErr: false,
		},
		"MissingEntityName": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.Event{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error reading entity name", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventName:     coretest.Event{}.EventName(),
					msg.MessageEventEntityID: "entity-id",
				})),
			},
			wantErr: false,
		},
		"MissingEntityID": {
			fields: fields{
				handlers: []handler{
					{
						evt: coretest.Event{},
						fn: func(ctx context.Context, evtMsg msg.EntityEvent) error {
							return nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error reading entity id", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageEventName:       coretest.Event{}.EventName(),
					msg.MessageEventEntityName: "entity-name",
				})),
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := msg.NewEntityEventDispatcher(msg.WithEntityEventDispatcherLogger(tt.fields.logger))
			for _, handler := range tt.fields.handlers {
				d.Handle(handler.evt, handler.fn)
			}
			if err := d.ReceiveMessage(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.logger)
		})
	}
}
