package saga

import (
	"github.com/stackus/edat/msg"
)

const (
	notCompensating = false
	isCompensating  = true
)

type LifecycleHook int

// Definition lifecycle hooks
const (
	SagaStarting LifecycleHook = iota
	SagaCompleted
	SagaCompensated
)

// Saga message headers
const (
	MessageCommandSagaID   = msg.MessageCommandPrefix + "SAGA_ID"
	MessageCommandSagaName = msg.MessageCommandPrefix + "SAGA_NAME"
	MessageCommandResource = msg.MessageCommandPrefix + "RESOURCE"

	MessageReplySagaID   = msg.MessageReplyPrefix + "SAGA_ID"
	MessageReplySagaName = msg.MessageReplyPrefix + "SAGA_NAME"
)
