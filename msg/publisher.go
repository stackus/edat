package msg

import (
	"context"
	"sync"
	"time"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
)

// CommandMessagePublisher interface
type CommandMessagePublisher interface {
	PublishCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error
}

// EntityEventMessagePublisher interface
type EntityEventMessagePublisher interface {
	PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error
}

// EventMessagePublisher interface
type EventMessagePublisher interface {
	PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error
}

// ReplyMessagePublisher interface
type ReplyMessagePublisher interface {
	PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error
}

// MessagePublisher interface
type MessagePublisher interface {
	Publish(ctx context.Context, message Message) error
}

var _ CommandMessagePublisher = (*Publisher)(nil)
var _ EntityEventMessagePublisher = (*Publisher)(nil)
var _ EventMessagePublisher = (*Publisher)(nil)
var _ ReplyMessagePublisher = (*Publisher)(nil)
var _ MessagePublisher = (*Publisher)(nil)

// Publisher send domain events, commands, and replies to the publisher
type Publisher struct {
	producer Producer
	logger   log.Logger
	close    sync.Once
}

// NewPublisher constructs a new Publisher
func NewPublisher(producer Producer, options ...PublisherOption) *Publisher {
	p := &Publisher{
		producer: producer,
		logger:   log.DefaultLogger,
	}

	for _, option := range options {
		option(p)
	}

	p.logger.Trace("msg.Publisher constructed")

	return p
}

// PublishCommand serializes a command into a message with command specific headers and publishes it to a producer
func (p *Publisher) PublishCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error {
	msgOptions := []MessageOption{WithHeaders(map[string]string{
		MessageCommandName:         command.CommandName(),
		MessageCommandReplyChannel: replyChannel,
	})}

	if v, ok := command.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		log.String("CommandName", command.CommandName()),
	)

	logger.Trace("publishing command")

	payload, err := core.SerializeCommand(command)
	if err != nil {
		logger.Error("error serializing command payload", log.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing command", log.Error(err))
	}

	return err
}

// PublishReply serializes a reply into a message with reply specific headers and publishes it to a producer
func (p *Publisher) PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error {
	msgOptions := []MessageOption{WithHeaders(map[string]string{
		MessageReplyName: reply.ReplyName(),
	})}

	if v, ok := reply.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		log.String("ReplyName", reply.ReplyName()),
	)

	logger.Trace("publishing reply")

	payload, err := core.SerializeReply(reply)
	if err != nil {
		logger.Error("error serializing reply payload", log.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing reply", log.Error(err))
	}

	return err
}

// PublishEntityEvents serializes entity events into messages with entity specific headers and publishes it to a producer
func (p *Publisher) PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error {
	msgOptions := []MessageOption{WithHeaders(map[string]string{
		MessageEventEntityID:   entity.ID(),
		MessageEventEntityName: entity.EntityName(),
		MessageChannel:         entity.EntityName(), // allow entity name and channel to overlap
	})}

	if v, ok := entity.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	for _, event := range entity.Events() {
		logger := p.logger.Sub(
			log.String("EntityID", entity.ID()),
			log.String("EntityName", entity.EntityName()),
		)

		err := p.PublishEvent(ctx, event, msgOptions...)
		if err != nil {
			logger.Error("error publishing entity event", log.Error(err))
			return err
		}
	}

	return nil
}

// PublishEvent serializes an event into a message with event specific headers and publishes it to a producer
func (p *Publisher) PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error {
	msgOptions := []MessageOption{WithHeaders(map[string]string{
		MessageEventName: event.EventName(),
	})}

	if v, ok := event.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		log.String("EventName", event.EventName()),
	)

	logger.Trace("publishing event")

	payload, err := core.SerializeEvent(event)
	if err != nil {
		logger.Error("error serializing event payload", log.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing event", log.Error(err))
	}

	return err
}

// Publish sends a message off to a producer
func (p *Publisher) Publish(ctx context.Context, message Message) error {
	var err error
	var channel string

	channel, err = message.Headers().GetRequired(MessageChannel)
	if err != nil {
		return err
	}

	message.Headers()[MessageDate] = time.Now().Format(time.RFC3339)

	// Published messages are request boundaries
	if id, exists := message.Headers()[MessageCorrelationID]; !exists || id == "" {
		message.Headers()[MessageCorrelationID] = core.GetCorrelationID(ctx)
	}

	if id, exists := message.Headers()[MessageCausationID]; !exists || id == "" {
		message.Headers()[MessageCausationID] = core.GetRequestID(ctx)
	}

	logger := p.logger.Sub(
		log.String("MessageID", message.ID()),
		log.String("CorrelationID", message.Headers()[MessageCorrelationID]),
		log.String("CausationID", message.Headers()[MessageCausationID]),
		log.String("Destination", channel),
		log.Int("PayloadSize", len(message.Payload())),
	)

	logger.Trace("publishing message")

	err = p.producer.Send(ctx, channel, message)
	if err != nil {
		logger.Error("error publishing message", log.Error(err))
		return err
	}

	return nil
}

// Stop stops the publisher and underlying producer
func (p *Publisher) Stop(ctx context.Context) (err error) {
	defer p.logger.Trace("publisher stopped")
	p.close.Do(func() {
		err = p.producer.Close(ctx)
	})

	return
}
