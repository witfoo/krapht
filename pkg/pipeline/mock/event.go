package mock

import "github.com/witfoo/krapht/pkg/pipeline"

// Event is a simple implementation of the Event interface for testing.
type Event struct {
	eventType pipeline.EventType
	message   string
}

// Type returns the type of the mock Event.
func (e Event) Type() pipeline.EventType {
	return e.eventType
}

// String returns the string representation of the mock Event.
func (e Event) String() string {
	return e.message
}

// NewEvent creates a new mock Event with the specified type and message.
func NewEvent(eventType pipeline.EventType, message string) Event {
	return Event{
		eventType: eventType,
		message:   message,
	}
}
