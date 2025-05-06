package pipeline

// LogLevel represents the log level of the event.
type LogLevel uint8

const (
	// LevelError represents an error event.
	LevelError LogLevel = iota + 3
	// LevelWarn represents a warning event.
	LevelWarn
	// Skip Notice level
	_
	// LevelInfo represents an informational event.
	LevelInfo
	// LevelDebug represents a debug event.
	LevelDebug
)

// Loggable is a specialized Event for logging
type Loggable interface {
	Event
	// Level returns the log level of the event
	Level() LogLevel
}

// LogEvent represents a log entry in the pipeline
type LogEvent struct {
	source string
	level  LogLevel
	msg    string
}

// NewLogEvent creates a new LogEvent instance
func NewLogEvent(source string, level LogLevel, msg string) Event {
	return LogEvent{
		source: source,
		level:  level,
		msg:    msg,
	}
}

// Type returns the type of event
func (l LogEvent) Type() EventType {
	return EventLog
}

// String returns the string representation of the event
func (l LogEvent) String() string {
	return l.source + ": " + l.msg
}

// Level returns the log level of the event
func (l LogEvent) Level() LogLevel {
	return l.level
}

// Send sends the log event to the event channel
func (l LogEvent) Send(eventC chan<- Event) {
	eventC <- l
}
