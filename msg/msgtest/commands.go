package msgtest

type (
	Command             struct{ Value string }
	UnregisteredCommand struct{ Value string }
)

func (Command) CommandName() string             { return "msgtest.Command" }
func (UnregisteredCommand) CommandName() string { return "msgtest.UnregisteredCommand" }

func (Command) DestinationChannel() string             { return "command-channel" }
func (UnregisteredCommand) DestinationChannel() string { return "command-channel" }
