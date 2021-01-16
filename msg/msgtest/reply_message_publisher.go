package msgtest

import (
	"github.com/stackus/edat/msg/msgmocks"
)

func MockReplyMessagePublisher(setup func(m *msgmocks.ReplyMessagePublisher)) *msgmocks.ReplyMessagePublisher {
	m := &msgmocks.ReplyMessagePublisher{}
	setup(m)
	return m
}
