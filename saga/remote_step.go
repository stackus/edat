package saga

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
)

// RemoteStep is used to execute distributed saga business logic
type RemoteStep struct {
	actionHandlers map[bool]*remoteStepAction
	replyHandlers  map[bool]map[string]func(context.Context, core.SagaData, core.Reply) error
}

var _ Step = (*RemoteStep)(nil)

// NewRemoteStep constructor for RemoteStep
func NewRemoteStep() RemoteStep {
	return RemoteStep{
		actionHandlers: map[bool]*remoteStepAction{
			notCompensating: nil,
			isCompensating:  nil,
		},
		replyHandlers: map[bool]map[string]func(context.Context, core.SagaData, core.Reply) error{
			notCompensating: {},
			isCompensating:  {},
		},
	}
}

// Action adds a domain command constructor that will be called while the definition is advancing
func (s RemoteStep) Action(fn func(context.Context, core.SagaData) msg.DomainCommand, options ...RemoteStepActionOption) RemoteStep {
	handler := &remoteStepAction{
		handler: fn,
	}

	for _, option := range options {
		option(handler)
	}

	s.actionHandlers[notCompensating] = handler

	return s
}

// HandleActionReply adds additional handling for specific replies while advancing
//
// SuccessReply and FailureReply do not require any special handling unless desired
func (s RemoteStep) HandleActionReply(reply core.Reply, handler func(context.Context, core.SagaData, core.Reply) error) RemoteStep {
	s.replyHandlers[notCompensating][reply.ReplyName()] = handler

	return s
}

// Compensation adds a domain command constructor that will be called while the definition is compensating
func (s RemoteStep) Compensation(fn func(context.Context, core.SagaData) msg.DomainCommand, options ...RemoteStepActionOption) RemoteStep {
	handler := &remoteStepAction{
		handler: fn,
	}

	for _, option := range options {
		option(handler)
	}

	s.actionHandlers[isCompensating] = handler

	return s
}

// HandleCompensationReply adds additional handling for specific replies while compensating
//
// SuccessReply does not require any special handling unless desired
func (s RemoteStep) HandleCompensationReply(reply core.Reply, handler func(context.Context, core.SagaData, core.Reply) error) RemoteStep {
	s.replyHandlers[isCompensating][reply.ReplyName()] = handler

	return s
}

func (s RemoteStep) hasInvocableAction(ctx context.Context, sagaData core.SagaData, compensating bool) bool {
	return s.actionHandlers[compensating] != nil && s.actionHandlers[compensating].isInvocable(ctx, sagaData)
}

func (s RemoteStep) getReplyHandler(replyName string, compensating bool) func(context.Context, core.SagaData, core.Reply) error {
	return s.replyHandlers[compensating][replyName]
}

func (s RemoteStep) execute(ctx context.Context, sagaData core.SagaData, compensating bool) func(results *stepResults) {
	if commandToSend := s.actionHandlers[compensating].execute(ctx, sagaData); commandToSend != nil {
		return func(actions *stepResults) {
			actions.commands = []msg.DomainCommand{commandToSend}
		}
	}

	return func(actions *stepResults) {}
}
