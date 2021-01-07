package msg

type Success struct{}

func (Success) ReplyName() string { return "edat.msg.Success" }

type Failure struct{}

func (Failure) ReplyName() string { return "edat.msg.Failure" }
