package msgtest

import (
	"github.com/stackus/edat/msg/msgmocks"
)

func MockConsumer(setup func(m *msgmocks.Consumer)) *msgmocks.Consumer {
	m := &msgmocks.Consumer{}
	setup(m)
	return m
}
