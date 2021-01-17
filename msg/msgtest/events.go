package msgtest

type (
	Event             struct{ Value string }
	UnregisteredEvent struct{ Value string }
)

func (Event) EventName() string             { return "msgtest.Event" }
func (UnregisteredEvent) EventName() string { return "msgtest.UnregisteredEvent" }

func (Event) DestinationChannel() string             { return "event-channel" }
func (UnregisteredEvent) DestinationChannel() string { return "event-channel" }
