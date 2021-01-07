package msg

import (
	"github.com/stackus/edat/es"
)

// MessageOption options for Message
type MessageOption func(m *message)

// WithMessageID is an option to set the ID of the Message
func WithMessageID(id string) MessageOption {
	return func(m *message) {
		m.id = id
		m.headers[MessageID] = id
	}
}

// WithDestinationChannel is and option to set the destination of the outgoing Message
//
// This will override the previous value set by interface { DestinationChannel() string }
func WithDestinationChannel(destinationChannel string) MessageOption {
	return func(m *message) {
		m.headers[MessageChannel] = destinationChannel
	}
}

// WithHeaders is an option to set additional headers onto the Message
func WithHeaders(headers Headers) MessageOption {
	return func(m *message) {
		for key, value := range headers {
			m.headers[key] = value
		}
	}
}

// WithAggregateInfo is an option to set additional Aggregate specific headers
func WithAggregateInfo(a *es.AggregateRoot) MessageOption {
	return func(m *message) {
		m.headers[MessageEventEntityName] = a.AggregateName()
		m.headers[MessageEventEntityID] = a.AggregateID()
	}
}
