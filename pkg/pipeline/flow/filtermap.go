package flow

import (
	"errors"

	"github.com/witfoo/krapht/pkg/pipeline"
)

var _ pipeline.Flow[any, any] = (*FilterMap[any, any])(nil)

// FilterMap is a struct that represents a filter and map operation on a data stream.
type FilterMap[I, O any] struct {
	predicate PredicateFunc[I] // The predicate function used for filtering the input data.
	transform MapFunc[I, O]    // The transform function used for mapping the filtered data.
}

// NewFilterMap creates a new instance of FilterMap with the given predicate function, transform function, and error channel.
func NewFilterMap[I, O any](predicate PredicateFunc[I], transform MapFunc[I, O]) (FilterMap[I, O], error) {

	if predicate == nil {
		return FilterMap[I, O]{}, errors.New("predicate func is nil")
	}

	if transform == nil {
		return FilterMap[I, O]{}, errors.New("transform func is nil")
	}

	return FilterMap[I, O]{
		predicate: predicate,
		transform: transform,
	}, nil
}

// Transform applies the filter and map operation on the input channel and returns the output channel.
func (f FilterMap[I, O]) Transform(in <-chan I, eventC chan<- pipeline.Event) <-chan O {
	out := make(chan O)
	go func() {
		defer close(out)
		for v := range in {
			if f.predicate(v) {
				val, err := f.transform(v)
				if err != nil {
					eventC <- pipeline.NewErrorEvent("filtermap error transforming data", err, true)
					continue
				}
				out <- val
			}
		}
	}()
	return out
}
