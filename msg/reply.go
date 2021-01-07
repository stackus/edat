package msg

import (
	"github.com/stackus/edat/core"
)

// Reply outcomes
const (
	ReplyOutcomeSuccess = "SUCCESS"
	ReplyOutcomeFailure = "FAILURE"
)

type Reply interface {
	Reply() core.Reply
	Headers() Headers
}

type replyMessage struct {
	reply   core.Reply
	headers Headers
}

func NewReply(reply core.Reply, headers Headers) Reply {
	return replyMessage{reply, headers}
}

func (m replyMessage) Reply() core.Reply {
	return m.reply
}

func (m replyMessage) Headers() Headers {
	return m.headers
}

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

func WithSuccess() Reply {
	return &replyMessage{
		reply: Success{},
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeSuccess,
			MessageReplyName:    Success{}.ReplyName(),
		},
	}
}

func WithFailure() Reply {
	return &replyMessage{
		reply: Failure{},
		headers: map[string]string{
			MessageReplyOutcome: ReplyOutcomeFailure,
			MessageReplyName:    Failure{}.ReplyName(),
		},
	}
}
