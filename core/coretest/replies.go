package coretest

import (
	"github.com/stackus/edat/core/coremocks"
)

type (
	Reply             struct{ Value string }
	UnregisteredReply struct{ Value string }
)

func (Reply) ReplyName() string             { return "coretest.Reply" }
func (UnregisteredReply) ReplyName() string { return "coretest.UnregisteredReply" }

func MockReply(setup func(m *coremocks.Reply)) *coremocks.Reply {
	m := &coremocks.Reply{}
	setup(m)
	return m
}
