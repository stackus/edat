package msg

import (
	"github.com/google/uuid"
)

// Message interface for messages containing payloads and headers
type Message interface {
	ID() string
	Headers() Headers
	Payload() []byte
}

// Message is used to pass events, commands, and replies to and from servers
type message struct {
	id      string
	headers Headers
	payload []byte
}

// NewMessage message constructor
func NewMessage(payload []byte, options ...MessageOption) Message {
	id := uuid.New().String()

	m := message{
		id:      id,
		payload: payload,
		headers: map[string]string{MessageID: id},
	}

	for _, option := range options {
		option(&m)
	}

	return m
}

// ID returns the message ID
func (m message) ID() string {
	return m.id
}

// Headers returns the message Headers
func (m message) Headers() Headers {
	return m.headers
}

// Payload returns the message payload
func (m message) Payload() []byte {
	return m.payload
}
