package msg

import (
	"github.com/stackus/edat/core"
)

func WithReply(reply core.Reply) *ReplyBuilder {
	return &ReplyBuilder{
		reply:   reply,
		headers: map[string]string{},
	}
}

type ReplyBuilder struct {
	reply   core.Reply
	headers map[string]string
}

func (b *ReplyBuilder) Reply(reply core.Reply) *ReplyBuilder {
	b.reply = reply
	return b
}

func (b *ReplyBuilder) Headers(headers map[string]string) *ReplyBuilder {
	for key, value := range headers {
		b.headers[key] = value
	}
	return b
}

func (b *ReplyBuilder) Success() Reply {
	if b.reply == nil {
		b.reply = Success{}
	}

	b.headers[MessageReplyOutcome] = ReplyOutcomeSuccess
	b.headers[MessageReplyName] = b.reply.ReplyName()

	return &replyMessage{
		reply:   b.reply,
		headers: b.headers,
	}
}

func (b *ReplyBuilder) Failure() Reply {
	if b.reply == nil {
		b.reply = Failure{}
	}

	b.headers[MessageReplyOutcome] = ReplyOutcomeFailure
	b.headers[MessageReplyName] = b.reply.ReplyName()

	return &replyMessage{
		reply:   b.reply,
		headers: b.headers,
	}
}
