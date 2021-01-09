package msg

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
)

// EntityEventHandlerFunc function handlers for msg.EntityEvent
type EntityEventHandlerFunc func(context.Context, EntityEvent) error

// EntityEventDispatcher is a MessageReceiver for DomainEvents
type EntityEventDispatcher struct {
	handlers map[string]EntityEventHandlerFunc
	logger   log.Logger
}

var _ MessageReceiver = (*EntityEventDispatcher)(nil)

// NewEntityEventDispatcher constructs a new EntityEventDispatcher
func NewEntityEventDispatcher(options ...EntityEventDispatcherOption) *EntityEventDispatcher {
	c := &EntityEventDispatcher{
		handlers: map[string]EntityEventHandlerFunc{},
		logger:   log.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Trace("msg.EntityEventDispatcher constructed")

	return c
}

// Handle adds a new Event that will be handled by EventMessageFunc handler
func (d *EntityEventDispatcher) Handle(evt core.Event, handler EntityEventHandlerFunc) *EntityEventDispatcher {
	d.logger.Trace("entity event handler added", log.String("EventName", evt.EventName()))
	d.handlers[evt.EventName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *EntityEventDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	eventName, err := message.Headers().GetRequired(MessageEventName)
	if err != nil {
		d.logger.Error("error reading event name", log.Error(err))
		return nil
	}

	entityName, err := message.Headers().GetRequired(MessageEventEntityName)
	if err != nil {
		d.logger.Error("error reading entity name", log.Error(err))
		return nil
	}

	entityID, err := message.Headers().GetRequired(MessageEventEntityID)
	if err != nil {
		d.logger.Error("error reading entity id", log.Error(err))
		return nil
	}

	logger := d.logger.Sub(
		log.String("EntityName", entityName),
		log.String("EntityID", entityID),
		log.String("EventName", eventName),
		log.String("MessageID", message.ID()),
	)

	logger.Debug("received entity event message")

	// check first for a handler of the event; It is possible events might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[eventName]
	if !exists {
		return nil
	}

	logger.Trace("entity event handler found")

	event, err := core.DeserializeEvent(eventName, message.Payload())
	if err != nil {
		logger.Error("error decoding entity event message payload", log.Error(err))
		return nil
	}

	evtMsg := entityEventMessage{entityID, entityName, event, message.Headers()}

	err = handler(ctx, evtMsg)
	if err != nil {
		logger.Error("entity event handler returned an error", log.Error(err))
	}

	return err
}
