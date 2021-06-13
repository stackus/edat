package msg

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
)

// MessageSubscriber interface
type MessageSubscriber interface {
	Subscribe(channel string, receiver MessageReceiver)
}

// Subscriber receives domain events, commands, and replies from the consumer
type Subscriber struct {
	consumer     Consumer
	logger       log.Logger
	middlewares  []func(MessageReceiver) MessageReceiver
	receivers    map[string][]MessageReceiver
	stopping     chan struct{}
	subscriberWg sync.WaitGroup
	close        sync.Once
}

// NewSubscriber constructs a new Subscriber
func NewSubscriber(consumer Consumer, options ...SubscriberOption) *Subscriber {
	s := &Subscriber{
		consumer:  consumer,
		receivers: make(map[string][]MessageReceiver),
		stopping:  make(chan struct{}),
		logger:    log.DefaultLogger,
	}

	for _, option := range options {
		option(s)
	}

	s.logger.Trace("msg.Subscriber constructed")

	return s
}

// Use appends middleware receivers to the receiver stack
func (s *Subscriber) Use(mws ...func(MessageReceiver) MessageReceiver) {
	if len(s.receivers) > 0 {
		panic("middleware must be added before any subscriptions are made")
	}

	s.middlewares = append(s.middlewares, mws...)
}

// Subscribe connects the receiver with messages from the channel on the consumer
func (s *Subscriber) Subscribe(channel string, receiver MessageReceiver) {
	if _, exists := s.receivers[channel]; !exists {
		s.receivers[channel] = []MessageReceiver{}
	}
	s.logger.Trace("subscribed", log.String("Channel", channel))
	s.receivers[channel] = append(s.receivers[channel], s.chain(receiver))
}

// Start begins listening to all of the channels sending received messages into them
func (s *Subscriber) Start(ctx context.Context) error {
	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	group, gCtx := errgroup.WithContext(cCtx)

	group.Go(func() error {
		select {
		case <-s.stopping:
			cancel()
		case <-gCtx.Done():
		}

		return nil
	})

	for c, r := range s.receivers {
		// reassign to avoid issues with anonymous func
		channel := c
		receivers := r

		s.subscriberWg.Add(1)

		group.Go(func() error {
			defer s.subscriberWg.Done()
			receiveMessageFunc := func(mCtx context.Context, message Message) error {
				mCtx = core.SetRequestContext(
					mCtx,
					message.ID(),
					message.Headers().Get(MessageCorrelationID),
					message.Headers().Get(MessageCausationID),
				)

				s.logger.Trace("received message",
					log.String("MessageID", message.ID()),
					log.String("CorrelationID", message.Headers().Get(MessageCorrelationID)),
					log.String("CausationID", message.Headers().Get(MessageCausationID)),
					log.Int("PayloadSize", len(message.Payload())),
				)

				rGroup, rCtx := errgroup.WithContext(mCtx)
				for _, r2 := range receivers {
					receiver := r2
					rGroup.Go(func() error {
						return receiver.ReceiveMessage(rCtx, message)
					})
				}

				return rGroup.Wait()
			}
			err := s.consumer.Listen(gCtx, channel, receiveMessageFunc)
			if err != nil {
				s.logger.Error("consumer stopped and returned an error", log.Error(err))
				return err
			}

			return nil
		})
	}

	return group.Wait()
}

// Stop stops the consumer and underlying consumer
func (s *Subscriber) Stop(ctx context.Context) (err error) {
	s.close.Do(func() {
		close(s.stopping)

		done := make(chan struct{})
		go func() {
			err = s.consumer.Close(ctx)
			s.subscriberWg.Wait()
			close(done)
		}()

		select {
		case <-done:
			s.logger.Trace("all receivers are done")
		case <-ctx.Done():
			s.logger.Warn("timed out waiting for all receivers to close")
		}
	})

	return
}

func (s *Subscriber) chain(receiver MessageReceiver) MessageReceiver {
	if len(s.middlewares) == 0 {
		return receiver
	}

	r := s.middlewares[len(s.middlewares)-1](receiver)
	for i := len(s.middlewares) - 2; i >= 0; i-- {
		r = s.middlewares[i](r)
	}

	return r
}
