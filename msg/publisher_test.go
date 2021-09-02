package msg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coremocks"
	"github.com/stackus/edat/core/coretest"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/log/logmocks"
	"github.com/stackus/edat/log/logtest"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/msg/msgmocks"
	"github.com/stackus/edat/msg/msgtest"
)

func TestPublisher_Publish(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx     context.Context
		message msg.Message
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "message-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{}`), msg.WithHeaders(map[string]string{
					msg.MessageChannel: "message-channel",
				})),
			},
			wantErr: false,
		},
		"MissingChannel": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				message: msg.NewMessage([]byte(`{}`)),
			},
			wantErr: true,
		},
		"ProducerError": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "message-channel", mock.Anything).Return(fmt.Errorf("producer-error"))
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error publishing message", mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{}`), msg.WithHeaders(map[string]string{
					msg.MessageChannel: "message-channel",
				})),
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.Publish(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}

func TestPublisher_PublishCommand(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx          context.Context
		replyChannel string
		command      core.Command
		options      []msg.MessageOption
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterCommands(coretest.Command{}, msgtest.Command{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "command-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:          context.Background(),
				replyChannel: "reply-channel",
				command:      coretest.Command{},
				options:      []msg.MessageOption{msg.WithDestinationChannel("command-channel")},
			},
			wantErr: false,
		},
		"UnregisteredCommand": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error serializing command payload", mock.Anything)
				}),
			},
			args: args{
				ctx:          context.Background(),
				replyChannel: "reply-channel",
				command:      coretest.UnregisteredCommand{},
				options:      []msg.MessageOption{msg.WithDestinationChannel("command-channel")},
			},
			wantErr: true,
		},
		"DomainCommand": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "command-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:          context.Background(),
				replyChannel: "reply-channel",
				command:      msgtest.Command{},
				options:      []msg.MessageOption{},
			},
			wantErr: false,
		},
		"MissingDestination": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error publishing command", mock.Anything)
				}),
			},
			args: args{
				ctx:          context.Background(),
				replyChannel: "reply-channel",
				command:      coretest.Command{},
				options:      []msg.MessageOption{},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.PublishCommand(tt.args.ctx, tt.args.replyChannel, tt.args.command, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("PublishCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}

func TestPublisher_PublishEvent(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx     context.Context
		event   core.Event
		options []msg.MessageOption
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterEvents(coretest.Event{}, msgtest.Event{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "event-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				event:   coretest.Event{},
				options: []msg.MessageOption{msg.WithDestinationChannel("event-channel")},
			},
			wantErr: false,
		},
		"UnregisteredEvent": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error serializing event payload", mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				event:   coretest.UnregisteredEvent{},
				options: []msg.MessageOption{msg.WithDestinationChannel("event-channel")},
			},
			wantErr: true,
		},
		"DomainEvent": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "event-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				event:   msgtest.Event{},
				options: []msg.MessageOption{},
			},
			wantErr: false,
		},
		"MissingDestination": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error publishing event", mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				event:   coretest.Event{},
				options: []msg.MessageOption{},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.PublishEvent(tt.args.ctx, tt.args.event, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("PublishEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}

func TestPublisher_PublishEntityEvents(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx     context.Context
		entity  core.Entity
		options []msg.MessageOption
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterEvents(coretest.Event{}, msgtest.Event{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, coretest.Entity{}.EntityName(), mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				entity:  &coretest.Entity{},
				options: []msg.MessageOption{},
			},
			wantErr: false,
		},
		"UnregisteredEvent": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error publishing entity event", mock.Anything).Once()
					m.On("Error", "error serializing event payload", mock.Anything).Once()
				}),
			},
			args: args{
				ctx: context.Background(),
				entity: coretest.MockEntity(func(m *coremocks.Entity) {
					m.On("ID").Return("entity-id")
					m.On("EntityName").Return("entity-name")
					m.On("Events").Return([]core.Event{&coretest.UnregisteredEvent{}})
				}),
				options: []msg.MessageOption{msg.WithDestinationChannel("event-channel")},
			},
			wantErr: true,
		},
		"DomainEventIgnored": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "entity-name", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				entity: coretest.MockEntity(func(m *coremocks.Entity) {
					m.On("ID").Return("entity-id")
					m.On("EntityName").Return("entity-name")
					m.On("Events").Return([]core.Event{&msgtest.Event{}})
				}),
				options: []msg.MessageOption{},
			},
			wantErr: false,
		},
		"DomainEntity": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "entity-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				entity:  &msgtest.Entity{},
				options: []msg.MessageOption{},
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.PublishEntityEvents(tt.args.ctx, tt.args.entity, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("PublishEntityEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}

func TestPublisher_PublishReply(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx     context.Context
		reply   core.Reply
		options []msg.MessageOption
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterReplies(coretest.Reply{}, msgtest.Reply{})

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "reply-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				reply:   coretest.Reply{},
				options: []msg.MessageOption{msg.WithDestinationChannel("reply-channel")},
			},
			wantErr: false,
		},
		"UnregisteredReply": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error serializing reply payload", mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				reply:   coretest.UnregisteredReply{},
				options: []msg.MessageOption{msg.WithDestinationChannel("reply-channel")},
			},
			wantErr: true,
		},
		"DomainReply": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Send", mock.Anything, "reply-channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				reply:   msgtest.Reply{},
				options: []msg.MessageOption{},
			},
			wantErr: false,
		},
		"MissingDestination": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "error publishing reply", mock.Anything)
				}),
			},
			args: args{
				ctx:     context.Background(),
				reply:   coretest.Reply{},
				options: []msg.MessageOption{},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.PublishReply(tt.args.ctx, tt.args.reply, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("PublishReply() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}

func TestPublisher_Stop(t *testing.T) {
	type fields struct {
		producer msg.Producer
		logger   log.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				producer: msgtest.MockProducer(func(m *msgmocks.Producer) {
					m.On("Close", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := msg.NewPublisher(tt.fields.producer, msg.WithLogger(tt.fields.logger))
			if err := p.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.producer, tt.fields.logger)
		})
	}
}
