package saga

import (
	"context"
	"strings"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

// CommandHandlerFunc function handlers for saga.Command
type CommandHandlerFunc func(context.Context, Command) ([]msg.Reply, error)

// CommandDispatcherOption options for CommandDispatcher
type CommandDispatcherOption interface {
	configureCommandDispatcher(*CommandDispatcher)
}

// CommandDispatcher is a MessageReceiver for Commands
type CommandDispatcher struct {
	publisher msg.ReplyMessagePublisher
	handlers  map[string]CommandHandlerFunc
	logger    log.Logger
}

var _ msg.MessageReceiver = (*CommandDispatcher)(nil)

// NewCommandDispatcher constructs a new CommandDispatcher
func NewCommandDispatcher(publisher msg.ReplyMessagePublisher, options ...CommandDispatcherOption) *CommandDispatcher {
	c := &CommandDispatcher{
		publisher: publisher,
		handlers:  map[string]CommandHandlerFunc{},
		logger:    log.DefaultLogger,
	}

	for _, option := range options {
		option.configureCommandDispatcher(c)
	}

	c.logger.Trace("saga.CommandDispatcher constructed")

	return c
}

// Handle adds a new Command that will be handled by handler
func (d *CommandDispatcher) Handle(cmd core.Command, handler CommandHandlerFunc) *CommandDispatcher {
	d.logger.Trace("saga command handler added", log.String("CommandName", cmd.CommandName()))
	d.handlers[cmd.CommandName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *CommandDispatcher) ReceiveMessage(ctx context.Context, message msg.Message) error {
	commandName, sagaID, sagaName, err := d.commandMessageInfo(message)
	if err != nil {
		return nil
	}

	logger := d.logger.Sub(
		log.String("CommandName", commandName),
		log.String("SagaName", sagaName),
		log.String("SagaID", sagaID),
		log.String("MessageID", message.ID()),
	)

	logger.Debug("received saga command message")

	// check first for a handler of the command; It is possible commands might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[commandName]
	if !exists {
		return nil
	}

	logger.Trace("saga command handler found")

	command, err := core.DeserializeCommand(commandName, message.Payload())
	if err != nil {
		logger.Error("error decoding saga command message payload", log.Error(err))
		return nil
	}

	replyChannel, err := message.Headers().GetRequired(msg.MessageCommandReplyChannel)
	if err != nil {
		logger.Error("error reading reply channel", log.Error(err))
		return nil
	}

	correlationHeaders := d.correlationHeaders(message.Headers())

	cmdMsg := commandMessage{sagaID, sagaName, command, correlationHeaders}

	replies, err := handler(ctx, cmdMsg)
	if err != nil {
		logger.Error("saga command handler returned an error", log.Error(err))
		rerr := d.sendReplies(ctx, replyChannel, []msg.Reply{msg.WithFailure()}, correlationHeaders)
		if rerr != nil {
			logger.Error("error sending replies", log.Error(rerr))
			return rerr
		}
		return nil
	}

	err = d.sendReplies(ctx, replyChannel, replies, correlationHeaders)
	if err != nil {
		logger.Error("error sending replies", log.Error(err))
		return err
	}

	return nil
}

func (d *CommandDispatcher) commandMessageInfo(message msg.Message) (string, string, string, error) {
	var err error
	var commandName, sagaID, sagaName string

	commandName, err = message.Headers().GetRequired(msg.MessageCommandName)
	if err != nil {
		d.logger.Error("error reading command name", log.Error(err))
		return "", "", "", err
	}

	sagaID, err = message.Headers().GetRequired(MessageCommandSagaID)
	if err != nil {
		d.logger.Error("error reading saga id", log.Error(err))
		return "", "", "", err
	}

	sagaName, err = message.Headers().GetRequired(MessageCommandSagaName)
	if err != nil {
		d.logger.Error("error reading saga name", log.Error(err))
		return "", "", "", err
	}

	return commandName, sagaID, sagaName, nil
}

func (d *CommandDispatcher) sendReplies(ctx context.Context, replyChannel string, replies []msg.Reply, correlationHeaders msg.Headers) error {
	for _, reply := range replies {
		if err := d.publisher.PublishReply(ctx, reply.Reply(),
			msg.WithHeaders(correlationHeaders),
			msg.WithHeaders(reply.Headers()),
			msg.WithDestinationChannel(replyChannel),
		); err != nil {
			return err
		}
	}

	return nil
}

func (d *CommandDispatcher) correlationHeaders(headers msg.Headers) msg.Headers {
	replyHeaders := make(map[string]string)
	for key, value := range headers {
		if key == msg.MessageCommandName {
			continue
		}

		if strings.HasPrefix(key, msg.MessageCommandPrefix) {
			replyHeader := msg.MessageReplyPrefix + key[len(msg.MessageCommandPrefix):]
			replyHeaders[replyHeader] = value
		}
	}

	return replyHeaders
}
