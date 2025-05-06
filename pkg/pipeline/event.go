package pipeline

// EventType represents the type of pipeline event
type EventType uint8

const (
	// EventLog represents a log entry
	EventLog EventType = iota
	// EventError represents a pipeline error
	EventError
	// EventMetric represents a metric measurement
	EventMetric
)

// Event interface represents an event that is sent vie Event Bus
type Event interface {
	// Type returns the type of event
	Type() EventType
	// String returns the string representation of the event
	String() string
}

// SendEvent attempts to send an event to the provided channel without blocking.
// Returns true if the event was sent, false if the channel is full or nil.
func SendEvent(eventC chan<- Event, event Event) bool {
	if eventC == nil {
		return false
	}

	select {
	case eventC <- event:
		return true
	default:
		return false
	}
}
