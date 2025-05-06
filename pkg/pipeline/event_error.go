package pipeline

// Errorable is a specialized Event for errors that also implements Go's error interface
type Errorable interface {
	Event
	error // implements the error interface
}

// ErrorEvent represents an error in the pipeline
type ErrorEvent struct {
	msg       string
	err       error
	temporary bool
}

// NewErrorEvent creates a new pipeline ErrorEvent
func NewErrorEvent(msg string, err error, temporary bool) Event {
	return ErrorEvent{
		msg:       msg,
		err:       err,
		temporary: temporary,
	}
}

// Type returns the event type
func (e ErrorEvent) Type() EventType {
	return EventError
}

// String returns a string representation of the error
func (e ErrorEvent) String() string {
	return e.msg + ": " + e.err.Error()
}

// Error implements the error interface
func (e ErrorEvent) Error() string {
	return e.String()
}

// IsTemporary indicates if this is a temporary error
func (e ErrorEvent) IsTemporary() bool {
	return e.temporary
}

// Unwrap returns the underlying error
func (e ErrorEvent) Unwrap() error {
	return e.err
}
