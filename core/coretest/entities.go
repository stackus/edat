package coretest

import (
	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coremocks"
)

type Entity struct {
	core.EntityBase
	Value string
}

func (Entity) EntityName() string   { return "msgtest.Entity" }
func (Entity) ID() string           { return "entity-id" }
func (Entity) Events() []core.Event { return []core.Event{&Event{}} }

func MockEntity(setup func(m *coremocks.Entity)) *coremocks.Entity {
	m := &coremocks.Entity{}
	setup(m)
	return m
}
