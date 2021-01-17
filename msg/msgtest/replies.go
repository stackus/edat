package msgtest

type (
	Reply             struct{ Value string }
	UnregisteredReply struct{ Value string }
)

func (Reply) ReplyName() string             { return "msgtest.Reply" }
func (UnregisteredReply) ReplyName() string { return "msgtest.UnregisteredReply" }

func (Reply) DestinationChannel() string             { return "reply-channel" }
func (UnregisteredReply) DestinationChannel() string { return "reply-channel" }
