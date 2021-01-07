package register_types

import (
	"github.com/stackus/edat/msg"
)

// RegisterTypes is called automatically after registering a new default marshaller
//
// There shouldn't be any reason to call this directly.
func RegisterTypes() {
	msg.RegisterTypes()
}
