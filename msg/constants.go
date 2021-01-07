package msg

// Message header keys
const (
	MessageID            = "ID"
	MessageDate          = "DATE"
	MessageChannel       = "CHANNEL"
	MessageCorrelationID = "CORRELATION_ID"
	MessageCausationID   = "CAUSATION_ID"

	MessageEventPrefix     = "EVENT_"
	MessageEventName       = MessageEventPrefix + "NAME"
	MessageEventEntityName = MessageEventPrefix + "ENTITY_NAME"
	MessageEventEntityID   = MessageEventPrefix + "ENTITY_ID"

	MessageCommandPrefix       = "COMMAND_"
	MessageCommandName         = MessageCommandPrefix + "NAME"
	MessageCommandChannel      = MessageCommandPrefix + "CHANNEL"
	MessageCommandReplyChannel = MessageCommandPrefix + "REPLY_CHANNEL"

	MessageReplyPrefix  = "REPLY_"
	MessageReplyName    = MessageReplyPrefix + "NAME"
	MessageReplyOutcome = MessageReplyPrefix + "OUTCOME"
)
