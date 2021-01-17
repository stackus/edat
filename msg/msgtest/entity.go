package msgtest

import (
	"github.com/stackus/edat/core/coretest"
)

type Entity struct {
	coretest.Entity
}

func (Entity) DestinationChannel() string { return "entity-channel" }
