package msg

import (
	"github.com/stackus/edat/core"
)

// Reply outcomes
const (
	ReplyOutcomeSuccess = "SUCCESS"
	ReplyOutcomeFailure = "FAILURE"
)

// Reply interface
type Reply interface {
	Reply() core.Reply
	Headers() Headers
}

type replyMessage struct {
	reply   core.Reply
	headers Headers
}

// NewReply constructs a new reply with headers
func NewReply(reply core.Reply, headers Headers) Reply {
	return &replyMessage{reply, headers}
}

// Reply returns the core.Reply
func (m replyMessage) Reply() core.Reply {
	return m.reply
}

// Headers returns the msg.Headers
func (m replyMessage) Headers() Headers {
	return m.headers
}

// SuccessReply wraps a reply and returns it as a Success reply
// Deprecated: Use the WithReply() reply builder
func SuccessReply(reply core.Reply) Reply {
	if reply == nil {
		return &replyMessage{
			reply: Success{},
			headers: map[string]string{
				MessageReplyOutcome: ReplyOutcomeSuccess,
				MessageReplyName:    Success{}.ReplyName(),
			},
		}
	}

	return &replyMessage{
		reply: reply,
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeSuccess,
			MessageReplyName:    reply.ReplyName(),
		},
	}
}

// FailureReply wraps a reply and returns it as a Failure reply
// Deprecated: Use the WithReply() reply builder
func FailureReply(reply core.Reply) Reply {
	if reply == nil {
		return &replyMessage{
			reply: Failure{},
			headers: map[string]string{
				MessageReplyOutcome: ReplyOutcomeFailure,
				MessageReplyName:    Failure{}.ReplyName(),
			},
		}
	}

	return &replyMessage{
		reply: reply,
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeFailure,
			MessageReplyName:    reply.ReplyName(),
		},
	}
}

// WithSuccess returns a generic Success reply
func WithSuccess() Reply {
	return &replyMessage{
		reply: Success{},
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeSuccess,
			MessageReplyName:    Success{}.ReplyName(),
		},
	}
}

// WithFailure returns a generic Failure reply
func WithFailure() Reply {
	return &replyMessage{
		reply: Failure{},
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeFailure,
			MessageReplyName:    Failure{}.ReplyName(),
		},
	}
}
