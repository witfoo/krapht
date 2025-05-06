package flow

import (
	"errors"

	"github.com/witfoo/krapht/pkg/pipeline"
)

var _ pipeline.Flow[any, any] = (*Map[any, any])(nil)

// MapFunc is a function that transforms an input value to an output value.
type MapFunc[I, O any] func(in I) (O, error)

// Map is a struct that represents a map operation on a data stream.
type Map[I, O any] struct {
	transform MapFunc[I, O]
}

// NewMap creates a new Map with the given transform function.
func NewMap[I, O any](transform MapFunc[I, O]) (*Map[I, O], error) {
	if transform == nil {
		return nil, errors.New("transform func is nil")
	}

	return &Map[I, O]{
		transform: transform,
	}, nil
}

// Transform applies the map operation on the input channel in a goroutine and returns the output channel.
func (m Map[I, O]) Transform(in <-chan I, eventC chan<- pipeline.Event) <-chan O {
	out := make(chan O)
	go func() {
		defer close(out)
		for v := range in {
			val, err := m.transform(v)
			if err != nil {
				// Send an error event when transform fails
				eventC <- pipeline.NewErrorEvent(
					"map transform error",
					err,
					true)
				continue
			}
			out <- val
		}
	}()
	return out
}
