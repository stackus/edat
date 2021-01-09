package msg

// Success reply type for generic successful replies to commands
type Success struct{}

// ReplyName implements core.Reply.ReplyName
func (Success) ReplyName() string { return "edat.msg.Success" }

// Failure reply type for generic failure replies to commands
type Failure struct{}

// ReplyName implements core.Reply.ReplyName
func (Failure) ReplyName() string { return "edat.msg.Failure" }
