package inmem

import (
	"context"
	"sync"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

var channels = sync.Map{}

// Consumer implements msg.Consumer
type Consumer struct {
	logger log.Logger
}

var _ msg.Consumer = (*Consumer)(nil)

// NewConsumer constructs a new Consumer
func NewConsumer(options ...ConsumerOption) *Consumer {
	c := &Consumer{
		logger: log.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// Listen implements msg.Consumer.Listen
func (c *Consumer) Listen(ctx context.Context, channel string, consumer msg.ReceiveMessageFunc) error {
	result, _ := channels.LoadOrStore(channel, make(chan msg.Message))

	messages := result.(chan msg.Message)

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return nil
			}
			err := consumer(ctx, message)
			if err != nil {
				c.logger.Error("error consuming message", log.Error(err))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// Close implements msg.Consumer.Close
func (c *Consumer) Close(context.Context) error {
	channels.Range(func(key, value interface{}) bool {
		messages := value.(chan msg.Message)
		close(messages)

		c.logger.Trace("closed channel", log.String("Channel", key.(string)))

		return true
	})

	c.logger.Trace("closing message source")
	return nil
}
