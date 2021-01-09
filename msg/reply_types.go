package msg

// Success reply type for generic successful replies to commands
type Success struct{}

func (Success) ReplyName() string { return "edat.msg.Success" }

// Failure reply type for generic failure replies to commands
type Failure struct{}

func (Failure) ReplyName() string { return "edat.msg.Failure" }
