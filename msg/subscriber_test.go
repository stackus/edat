package msg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/log/logmocks"
	"github.com/stackus/edat/log/logtest"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/msg/msgmocks"
	"github.com/stackus/edat/msg/msgtest"
)

func TestSubscriber_Start(t *testing.T) {
	type fields struct {
		consumer    msg.Consumer
		logger      log.Logger
		middlewares []func(msg.MessageReceiver) msg.MessageReceiver
		receivers   map[string]msg.MessageReceiver
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
				consumer: msgtest.MockConsumer(func(m *msgmocks.Consumer) {
					m.On("Listen", mock.Anything, "channel", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
				}),
				middlewares: []func(msg.MessageReceiver) msg.MessageReceiver{
					func(next msg.MessageReceiver) msg.MessageReceiver {
						return next
					},
				},
				receivers: map[string]msg.MessageReceiver{
					"channel": msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := msg.NewSubscriber(tt.fields.consumer, msg.WithSubscriberLogger(tt.fields.logger))
			s.Use(tt.fields.middlewares...)
			for channel, receiver := range tt.fields.receivers {
				s.Subscribe(channel, receiver)
			}
			ctx, cancel := context.WithTimeout(tt.args.ctx, 1*time.Millisecond)
			defer cancel()

			err := s.Start(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.consumer, tt.fields.logger)
		})
	}
}

func TestSubscriber_Stop(t *testing.T) {
	type fields struct {
		consumer    msg.Consumer
		logger      log.Logger
		middlewares []func(msg.MessageReceiver) msg.MessageReceiver
		receivers   map[string]msg.MessageReceiver
	}
	type args struct {
		ctx context.Context
	}
	tests := map[string]struct {
		fields       fields
		args         args
		wantStartErr bool
		wantStopErr  bool
	}{
		"Success": {
			fields: fields{
				consumer: msgtest.MockConsumer(func(m *msgmocks.Consumer) {
					m.On("Listen", mock.Anything, "channel", mock.Anything).Return(nil)
					m.On("Close", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", "msg.Subscriber constructed", mock.Anything)
					m.On("Trace", "subscribed", mock.Anything)
					m.On("Trace", "all receivers are done", mock.Anything)
				}),
				middlewares: []func(msg.MessageReceiver) msg.MessageReceiver{},
				receivers: map[string]msg.MessageReceiver{
					"channel": msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantStopErr: false,
		},
		"ConsumerError": {
			fields: fields{
				consumer: msgtest.MockConsumer(func(m *msgmocks.Consumer) {
					m.On("Listen", mock.Anything, "channel", mock.Anything).Return(fmt.Errorf("consumer-error"))
					m.On("Close", mock.Anything).Return(nil)
				}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
					m.On("Error", "consumer stopped and returned an error", mock.Anything)
				}),
				middlewares: []func(msg.MessageReceiver) msg.MessageReceiver{},
				receivers: map[string]msg.MessageReceiver{
					"channel": msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantStartErr: true,
			wantStopErr:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := msg.NewSubscriber(tt.fields.consumer, msg.WithSubscriberLogger(tt.fields.logger))
			s.Use(tt.fields.middlewares...)
			for channel, receiver := range tt.fields.receivers {
				s.Subscribe(channel, receiver)
			}
			stopped := make(chan struct{})
			var startErr error
			go func() {
				startErr = s.Start(context.Background())
				close(stopped)
			}()
			time.Sleep(1 * time.Millisecond) // hack to give goroutine time to start and avoid a data race
			err := s.Stop(tt.args.ctx)
			<-stopped
			if (err != nil) != tt.wantStopErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantStopErr)
			}
			if (startErr != nil) != tt.wantStartErr {
				t.Errorf("Start() error = %v, wantErr %v", startErr, tt.wantStartErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.consumer, tt.fields.logger)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	type receivers struct {
		channel  string
		receiver msg.MessageReceiver
	}
	type fields struct {
		consumer msg.Consumer
		logger   log.Logger
	}
	type args struct {
		receivers []receivers
	}
	tests := map[string]struct {
		fields    fields
		args      args
		wantPanic bool
	}{
		"Success": {
			fields: fields{
				consumer: msgtest.MockConsumer(func(m *msgmocks.Consumer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				receivers: []receivers{
					{
						channel:  "channel",
						receiver: msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
					},
				},
			},
			wantPanic: false,
		},
		"Duplicate": {
			fields: fields{
				consumer: msgtest.MockConsumer(func(m *msgmocks.Consumer) {}),
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Trace", mock.Anything, mock.Anything)
				}),
			},
			args: args{
				receivers: []receivers{
					{
						channel:  "channel",
						receiver: msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
					},
					{
						channel:  "channel",
						receiver: msgtest.MockMessageReceiver(func(m *msgmocks.MessageReceiver) {}),
					},
				},
			},
			wantPanic: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := msg.NewSubscriber(tt.fields.consumer, msg.WithSubscriberLogger(tt.fields.logger))
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("Subscribe() = %v, wantPanic %v", r, tt.wantPanic)
				}
				mock.AssertExpectationsForObjects(t, tt.fields.consumer, tt.fields.logger)
			}()
			for _, r := range tt.args.receivers {
				s.Subscribe(r.channel, r.receiver)
			}
		})
	}
}
