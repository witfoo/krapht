package sink

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/witfoo/krapht/pkg/pipeline"
)

// Static check to ensure that the Sink interface is implemented by the Logger struct
var _ pipeline.Sink[any] = (*Logger[any])(nil)

// Logger is a sink that logs the data
type Logger[I any] struct {
	logger *log.Logger
}

// NewLogger creates a new Logger instance
func NewLogger[I any](logger *log.Logger) *Logger[I] {
	return &Logger[I]{
		logger: logger,
	}
}

// Load logs the data
func (l Logger[I]) Load(in <-chan I, eventC chan<- pipeline.Event) {
	var formattedStr string

	for data := range in {
		// try to format the data as JSON
		formattedB, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			// Fall back to standard Go formatting
			formattedStr = fmt.Sprintf("%+v", data)
		}
		formattedStr = string(formattedB)

		if len(formattedStr) == 0 {
			// If the formatted string is empty, send a warning event
			eventC <- pipeline.NewErrorEvent("formatted data is empty", nil, true)
		} else {
			// Log the formatted string
			l.logger.Println(formattedStr)
		}
	}
}
