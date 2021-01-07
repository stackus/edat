package msg

import (
	"github.com/stackus/edat/core"
)

// DomainCommand interface for commands that are shared across the domain
type DomainCommand interface {
	core.Command
	DestinationChannel() string
}

// Command is a core.Command with message header information
type Command interface {
	Command() core.Command
	Headers() Headers
}

type commandMessage struct {
	command core.Command
	headers Headers
}

func (m commandMessage) Command() core.Command {
	return m.command
}

func (m commandMessage) Headers() Headers {
	return m.headers
}
