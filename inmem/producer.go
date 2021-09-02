package inmem

import (
	"context"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

// ProducerOption options for Producer
type ProducerOption interface {
	configureProducer(*Producer)
}

// Producer implements msg.Producer
type Producer struct {
	logger log.Logger
}

var _ msg.Producer = (*Producer)(nil)

// NewProducer constructs a new Producer
func NewProducer(options ...ProducerOption) *Producer {
	p := &Producer{
		logger: log.DefaultLogger,
	}

	for _, option := range options {
		option.configureProducer(p)
	}

	return p
}

// Send implements msg.Producer.Send
func (p *Producer) Send(_ context.Context, channel string, message msg.Message) error {
	if result, exists := channels.Load(channel); exists {
		destination := result.(chan msg.Message)

		destination <- message

		p.logger.Trace("message sent to inmem channel", log.String("Channel", channel))
	}

	return nil
}

// Close implements msg.Producer.Close
func (p *Producer) Close(context.Context) error {
	p.logger.Trace("closing message destination")
	return nil
}
