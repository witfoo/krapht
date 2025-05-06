package pipeline_test

import (
	"testing"

	"github.com/witfoo/krapht/pkg/pipeline"
)

func TestNewLogEvent(t *testing.T) {
	source := "test_source"
	level := pipeline.LevelInfo
	msg := "test message"

	event := pipeline.NewLogEvent(source, level, msg)

	// Check that we get the correct type
	loggable, ok := event.(pipeline.Loggable)
	if !ok {
		t.Fatal("Expected NewLogEvent to return an event implementing Loggable")
	}

	// Test through the public interface
	if loggable.Level() != level {
		t.Errorf("Expected level %v, got %v", level, loggable.Level())
	}

	// Test string representation
	expectedStr := source + ": " + msg
	if event.String() != expectedStr {
		t.Errorf("Expected string representation %s, got %s", expectedStr, event.String())
	}
}

func TestLogEventType(t *testing.T) {
	source := "test_source"
	level := pipeline.LevelInfo
	msg := "test message"

	event := pipeline.NewLogEvent(source, level, msg)

	if event.Type() != pipeline.EventLog {
		t.Errorf("Expected event type %v, got %v", pipeline.EventLog, event.Type())
	}
}

func TestLogEventString(t *testing.T) {
	source := "test_source"
	level := pipeline.LevelInfo
	msg := "test message"

	event := pipeline.NewLogEvent(source, level, msg)

	expectedStr := source + ": " + msg
	if event.String() != expectedStr {
		t.Errorf("Expected string representation %s, got %s", expectedStr, event.String())
	}
}

func TestLogEventLevel(t *testing.T) {
	source := "test_source"
	level := pipeline.LevelDebug
	msg := "test message"

	event := pipeline.NewLogEvent(source, level, msg)
	loggable, ok := event.(pipeline.Loggable)
	if !ok {
		t.Fatal("Event doesn't implement Loggable interface")
	}

	if loggable.Level() != level {
		t.Errorf("Expected level %v, got %v", level, loggable.Level())
	}
}

func TestLogLevels(t *testing.T) {
	// Test that the log levels are defined correctly
	tests := []struct {
		level    pipeline.LogLevel
		expected uint8
	}{
		{pipeline.LevelError, 3},
		{pipeline.LevelWarn, 4},
		// Notice level is skipped (5)
		{pipeline.LevelInfo, 6},
		{pipeline.LevelDebug, 7},
	}

	for _, test := range tests {
		if uint8(test.level) != test.expected {
			t.Errorf("Expected log level %v to have value %d, got %d", test.level, test.expected, uint8(test.level))
		}
	}
}

func TestLoggableInterface(t *testing.T) {
	// Test that the event returned by NewLogEvent implements the Loggable interface
	source := "test_source"
	level := pipeline.LevelInfo
	msg := "test message"

	event := pipeline.NewLogEvent(source, level, msg)

	_, ok := event.(pipeline.Loggable)
	if !ok {
		t.Error("Event created with NewLogEvent does not implement Loggable")
	}
}
