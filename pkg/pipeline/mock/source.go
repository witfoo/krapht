package mock

import (
	"context"
	"fmt"

	"github.com/witfoo/krapht/pkg/pipeline"
)

// Ensure mockSource implements the pipeline.Source interface
var _ pipeline.Source[ReadableImpl] = (*SourceImpl)(nil)

// SourceImpl is a simple source implementation for testing
type SourceImpl struct {
	data []ReadableImpl // mock data that implements the Readable interface
}

// NewSourceImpl creates a new SourceImpl with the given data
// It implements the pipeline.Source interface
// and is used for testing purposes
// It simulates a source that extracts data from a slice of ReadableImpl
func NewSourceImpl(data []ReadableImpl) *SourceImpl {
	return &SourceImpl{
		data: data,
	}
}

// Extract simulates data extraction from a source
func (s *SourceImpl) Extract(ctx context.Context, eventC chan<- pipeline.Event) <-chan ReadableImpl {
	out := make(chan ReadableImpl)
	go func() {
		defer close(out)
		for _, d := range s.data {
			// Read the data from the Readable interface
			source, err := d.Read()
			if err != nil {
				// Send an error event to the event channel
				pipeline.SendEvent(eventC, pipeline.NewErrorEvent("mockSource error reading data", err, true))
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- d:
				// Send an event to the event channel
				pipeline.SendEvent(eventC, pipeline.NewLogEvent("mockSource", pipeline.LevelInfo, fmt.Sprintf("Extracted item: %s", source)))
			}
		}
	}()
	return out
}
