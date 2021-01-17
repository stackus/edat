package msgtest

import (
	"github.com/stackus/edat/msg/msgmocks"
)

func MockMessageReceiver(setup func(m *msgmocks.MessageReceiver)) *msgmocks.MessageReceiver {
	m := &msgmocks.MessageReceiver{}
	setup(m)
	return m
}
