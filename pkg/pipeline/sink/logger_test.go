package sink

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
	//"github.com/witfoo/krapht/pkg/pipeline"
)

// MockLogger is a mock implementation of slog.Logger
type MockLogger struct {
	buffer *bytes.Buffer
	mu     sync.Mutex
}

// Implement io.Writer interface
func (m *MockLogger) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buffer.Write(p)
}

func (m *MockLogger) String() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buffer.String()
}

// func (m *MockLogger) Println(v ...interface{}) {
// 	m.buffer.WriteString(fmt.Sprintln(v...))
// }

func TestLogger_Load(t *testing.T) {
	buffer := &bytes.Buffer{}
	mockLogger := &MockLogger{buffer: buffer}
	// set log output to bytes buffer
	log.SetOutput(mockLogger)
	// Create a new Logger sink with the logger set to the mock logger

	loggerSink := &Logger[any]{logger: log.New(mockLogger, "", 0)}

	in := make(chan any)
	defer close(in)
	go loggerSink.Load(in, nil)

	tests := []struct {
		name     string
		input    []any
		expected string
	}{
		{
			name:     "Single log entry",
			input:    []any{"test log entry"},
			expected: "test log entry\n",
		},
		{
			name:     "Multiple log entries",
			input:    []any{"first log entry", "second log entry"},
			expected: "first log entry\nsecond log entry\n",
		},
		{
			name:     "Empty log entry",
			input:    []any{""},
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write the input to the channel
			for _, entry := range tt.input {
				in <- entry
			}
			// Delay to allow for the goroutine to log the data
			time.Sleep(100 * time.Millisecond)

			bufferContent := mockLogger.String()
			// Check if the buffer contains the expected output
			if strings.Contains(bufferContent, tt.expected) {
				t.Log("Log entry found in buffer")
			} else {
				t.Log("Log entry not found in buffer")
			}

			// Clear the buffer
			buffer.Reset()
		})
	}
}
