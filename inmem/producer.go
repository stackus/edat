package inmem

import (
	"context"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

type Producer struct {
	logger log.Logger
}

var _ msg.Producer = (*Producer)(nil)

func NewProducer(options ...ProducerOption) *Producer {
	p := &Producer{
		logger: log.DefaultLogger,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *Producer) Send(_ context.Context, channel string, message msg.Message) error {
	if result, exists := channels.Load(channel); exists {
		destination := result.(chan msg.Message)

		destination <- message

		p.logger.Trace("message sent to inmem channel", log.String("Channel", channel))
	}

	return nil
}

func (p *Producer) Close(context.Context) error {
	p.logger.Trace("closing message destination")
	return nil
}
