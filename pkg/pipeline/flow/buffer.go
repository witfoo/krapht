package flow

import "github.com/witfoo/krapht/pkg/pipeline"

// Ensure that Buffer implements the Flow interface.
var _ pipeline.Flow[any, any] = (*Buffer[any])(nil)

// Buffer is a struct that represents a passthrough with buffered channel operation on a data stream.
type Buffer[I any] struct {
	size int
}

// NewBuffer creates a new Buffer flow.
// Size allows setting the buffer size for the output channel.
// This results in a very simple FIFO buffer implementation.
// No backoff, eviction, or drain policies are implemented.
// The buffer will grow to the size specified and then block until space is available.
// If size is less than or equal to 0, it will default to 1.
// This is a passthrough implementation with a buffered channel.
func NewBuffer[I any](size int) *Buffer[I] {
	if size <= 0 {
		size = 1
	}
	return &Buffer[I]{
		size: size,
	}
}

// Transform applies the passthrough operation on the input channel and returns the output channel.
// The passthrough operation simply forwards the data from the input channel to the buffered output channel.
func (b Buffer[I]) Transform(in <-chan I, _ chan<- pipeline.Event) <-chan I {
	out := make(chan I, b.size)

	go func() {
		defer close(out)
		for v := range in {
			out <- v
		}
	}()
	return out
}
