package registertypes

import (
	"github.com/stackus/edat/msg"
)

// RegisterTypes registers internal library types
//
// Marshaller implementors: This should be called automatically after registering a new default marshaller.
//
// Users: There shouldn't be any reason to call this directly.
func RegisterTypes() {
	msg.RegisterTypes()
}
