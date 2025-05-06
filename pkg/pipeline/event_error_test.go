package pipeline_test

import (
	"errors"
	"testing"

	"github.com/witfoo/krapht/pkg/pipeline"
)

func TestNewErrorEvent(t *testing.T) {
	msg := "operation failed"
	originalErr := errors.New("database connection failed")
	temporary := true

	event := pipeline.NewErrorEvent(msg, originalErr, temporary)

	// Check that the returned event is of the correct type
	_, ok := event.(pipeline.Errorable)
	if !ok {
		t.Fatal("Expected NewErrorEvent to return an event implementing Errorable")
	}

	// Check through error interface
	err, ok := event.(error)
	if !ok {
		t.Fatal("Expected ErrorEvent to implement error interface")
	}

	expectedErrMsg := msg + ": " + originalErr.Error()
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestErrorEventType(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")
	event := pipeline.NewErrorEvent(msg, originalErr, false)

	if event.Type() != pipeline.EventError {
		t.Errorf("Expected event type %v, got %v", pipeline.EventError, event.Type())
	}
}

func TestErrorEventString(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")
	event := pipeline.NewErrorEvent(msg, originalErr, false)

	expected := msg + ": " + originalErr.Error()
	if event.String() != expected {
		t.Errorf("Expected string %q, got %q", expected, event.String())
	}
}

func TestErrorEventError(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")
	event := pipeline.NewErrorEvent(msg, originalErr, false)

	errEvent, ok := event.(error)
	if !ok {
		t.Fatal("Event does not implement error interface")
	}

	expected := msg + ": " + originalErr.Error()
	if errEvent.Error() != expected {
		t.Errorf("Expected error string %q, got %q", expected, errEvent.Error())
	}
}

func TestErrorEventIsTemporary(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")

	// Test with temporary = true
	temporaryEvent := pipeline.NewErrorEvent(msg, originalErr, true)
	tErr, ok := temporaryEvent.(interface{ IsTemporary() bool })
	if !ok {
		t.Fatal("Event does not implement IsTemporary method")
	}
	if !tErr.IsTemporary() {
		t.Error("Expected IsTemporary() to return true")
	}

	// Test with temporary = false
	permanentEvent := pipeline.NewErrorEvent(msg, originalErr, false)
	pErr, ok := permanentEvent.(interface{ IsTemporary() bool })
	if !ok {
		t.Fatal("Event does not implement IsTemporary method")
	}
	if pErr.IsTemporary() {
		t.Error("Expected IsTemporary() to return false")
	}
}

func TestErrorEventUnwrap(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")
	event := pipeline.NewErrorEvent(msg, originalErr, false)

	unwrappable, ok := event.(interface{ Unwrap() error })
	if !ok {
		t.Fatal("Event does not implement Unwrap method")
	}

	unwrappedErr := unwrappable.Unwrap()
	if unwrappedErr != originalErr {
		t.Errorf("Expected unwrapped error to be %v, got %v", originalErr, unwrappedErr)
	}
}

func TestErrorableInterface(t *testing.T) {
	msg := "test error"
	originalErr := errors.New("underlying error")
	event := pipeline.NewErrorEvent(msg, originalErr, false)

	// Test it implements both Event and error interfaces (Errorable)
	_, ok1 := event.(pipeline.Event)
	if !ok1 {
		t.Error("Event does not implement Event interface")
	}

	_, ok2 := event.(error)
	if !ok2 {
		t.Error("Event does not implement error interface")
	}

	_, ok3 := event.(pipeline.Errorable)
	if !ok3 {
		t.Error("Event does not implement Errorable interface")
	}
}
