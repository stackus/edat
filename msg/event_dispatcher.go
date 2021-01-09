package msg

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
)

// EventHandlerFunc function handlers for msg.Event
type EventHandlerFunc func(context.Context, Event) error

// EventDispatcher is a MessageReceiver for Events
type EventDispatcher struct {
	handlers map[string]EventHandlerFunc
	logger   log.Logger
}

var _ MessageReceiver = (*EventDispatcher)(nil)

// NewEventDispatcher constructs a new EventDispatcher
func NewEventDispatcher(options ...EventDispatcherOption) *EventDispatcher {
	c := &EventDispatcher{
		handlers: map[string]EventHandlerFunc{},
		logger:   log.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Trace("msg.EventDispatcher constructed")

	return c
}

// Handle adds a new Event that will be handled by EventMessageFunc handler
func (d *EventDispatcher) Handle(evt core.Event, handler EventHandlerFunc) *EventDispatcher {
	d.logger.Trace("event handler added", log.String("EventName", evt.EventName()))
	d.handlers[evt.EventName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *EventDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	eventName, err := message.Headers().GetRequired(MessageEventName)
	if err != nil {
		d.logger.Error("error reading event name", log.Error(err))
		return nil
	}

	logger := d.logger.Sub(
		log.String("EventName", eventName),
		log.String("MessageID", message.ID()),
	)

	logger.Debug("received event message")

	// check first for a handler of the event; It is possible events might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[eventName]
	if !exists {
		return nil
	}

	logger.Trace("event handler found")

	event, err := core.DeserializeEvent(eventName, message.Payload())
	if err != nil {
		logger.Error("error decoding event message payload", log.Error(err))
		return nil
	}

	evtMsg := eventMessage{event, message.Headers()}

	err = handler(ctx, evtMsg)
	if err != nil {
		logger.Error("event handler returned an error", log.Error(err))
	}

	return err
}
