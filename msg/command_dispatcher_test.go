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
	"github.com/stackus/edat/msg/msgmocks"
	"github.com/stackus/edat/msg/msgtest"
)

func TestCommandDispatcher_ReceiveMessage(t *testing.T) {
	type handler struct {
		cmd core.Command
		fn  msg.CommandHandlerFunc
	}
	type fields struct {
		publisher msg.ReplyMessagePublisher
		handlers  []handler
		logger    log.Logger
	}
	type args struct {
		ctx     context.Context
		message msg.Message
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterCommands(coretest.Command{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {
					m.On("PublishReply", mock.Anything, mock.AnythingOfType("msg.Success"), mock.Anything, mock.Anything, mock.Anything).Return(nil)
				}),
				handlers: []handler{
					{
						cmd: coretest.Command{},
						fn: func(ctx context.Context, command msg.Command) ([]msg.Reply, error) {
							return []msg.Reply{msg.WithSuccess()}, nil
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
					msg.MessageCommandName:         coretest.Command{}.CommandName(),
					msg.MessageCommandReplyChannel: "reply-channel",
				})),
			},
			wantErr: false,
		},
		"HandlerError": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {
					m.On("PublishReply", mock.Anything, mock.AnythingOfType("msg.Failure"), mock.Anything, mock.Anything, mock.Anything).Return(nil)
				}),
				handlers: []handler{
					{
						cmd: coretest.Command{},
						fn: func(ctx context.Context, command msg.Command) ([]msg.Reply, error) {
							return nil, fmt.Errorf("handler error")
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "command handler returned an error", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageCommandName:         coretest.Command{}.CommandName(),
					msg.MessageCommandReplyChannel: "reply-channel",
				})),
			},
			wantErr: false,
		},
		"UnregisteredCommand": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {}),
				handlers: []handler{
					{
						cmd: coretest.UnregisteredCommand{},
						fn: func(ctx context.Context, command msg.Command) ([]msg.Reply, error) {
							return []msg.Reply{msg.WithSuccess()}, nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error decoding command message payload", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageCommandName:         coretest.UnregisteredCommand{}.CommandName(),
					msg.MessageCommandReplyChannel: "reply-channel",
				})),
			},
			wantErr: false,
		},
		"MissingCommandName": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {}),
				handlers: []handler{
					{
						cmd: coretest.Command{},
						fn: func(ctx context.Context, command msg.Command) ([]msg.Reply, error) {
							return []msg.Reply{msg.WithSuccess()}, nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error reading command name", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageCommandReplyChannel: "reply-channel",
				})),
			},
			wantErr: false,
		},
		"MissingReplyChannel": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {}),
				handlers: []handler{
					{
						cmd: coretest.Command{},
						fn: func(ctx context.Context, command msg.Command) ([]msg.Reply, error) {
							return []msg.Reply{msg.WithSuccess()}, nil
						},
					},
				},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
					m.On("Error", "error reading reply channel", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(map[string]string{
					msg.MessageCommandName: coretest.Command{}.CommandName(),
				})),
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := msg.NewCommandDispatcher(tt.fields.publisher, msg.WithCommandDispatcherLogger(tt.fields.logger))
			for _, handler := range tt.fields.handlers {
				d.Handle(handler.cmd, handler.fn)
			}
			if err := d.ReceiveMessage(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.publisher, tt.fields.logger)
		})
	}
}
