package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
)

type EventStore struct {
	events map[string][]eventMsg
	mu     sync.Mutex
}

type eventMsg struct {
	eventName string
	event     json.RawMessage
}

var _ es.AggregateRootStore = (*EventStore)(nil)

func NewEventStore(options ...EventStoreOption) *EventStore {
	s := &EventStore{
		events: make(map[string][]eventMsg),
		mu:     sync.Mutex{},
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *EventStore) Load(_ context.Context, root *es.AggregateRoot) error {
	// just lock it all
	s.mu.Lock()
	defer s.mu.Unlock()

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()

	if messages, exists := s.events[s.streamID(name, id)]; exists {
		if len(messages) < version {
			return nil
		}

		for _, message := range messages[version+1:] {
			event, err := core.DeserializeEvent(message.eventName, message.event)
			if err != nil {
				return err
			}
			err = root.LoadEvent(event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *EventStore) Save(_ context.Context, root *es.AggregateRoot) error {
	// just lock it all
	s.mu.Lock()
	defer s.mu.Unlock()

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.Version()
	streamID := s.streamID(name, id)

	if _, exists := s.events[streamID]; !exists {
		s.events[streamID] = []eventMsg{}
	}

	streamLength := len(s.events[streamID])

	if streamLength != version {
		return fmt.Errorf("stream has been modified and cannot write event: expected: %d, got: %d", version, streamLength)
	}

	for _, event := range root.Events() {
		data, err := core.SerializeEvent(event)
		if err != nil {
			return err
		}

		s.events[streamID] = append(s.events[streamID], eventMsg{
			eventName: event.EventName(),
			event:     data,
		})
	}

	return nil
}

func (s *EventStore) streamID(name, id string) string {
	return fmt.Sprintf("%s:%s", name, id)
}
