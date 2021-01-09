package outbox

import (
	"encoding/json"

	"github.com/stackus/edat/msg"
)

// Message struct for the temporary form of a Producers msg.Message
type Message struct {
	MessageID   string
	Destination string
	Payload     []byte
	Headers     []byte
}

type message struct {
	id      string
	headers msg.Headers
	payload []byte
}

// ToMessage converts this form back to msg.Message or returns an error when headers cannot be unmarshalled
func (m Message) ToMessage() (msg.Message, error) {
	var headers map[string]string

	err := json.Unmarshal(m.Headers, &headers)
	if err != nil {
		return nil, err
	}

	return message{
		id:      m.MessageID,
		headers: headers,
		payload: m.Payload,
	}, nil
}

func (m message) ID() string {
	return m.id
}

func (m message) Headers() msg.Headers {
	return m.headers
}

func (m message) Payload() []byte {
	return m.payload
}
