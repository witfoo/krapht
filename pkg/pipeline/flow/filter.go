// Package flow provides a set of functions to transform data streams
package flow

import (
	"errors"

	"github.com/witfoo/krapht/pkg/pipeline"
)

// Ensure Filter implements the Flow interface
var _ pipeline.Flow[any, any] = (*Filter[any])(nil)

// PredicateFunc is a function type that takes an input of type I and returns a boolean value.
type PredicateFunc[I any] func(in I) bool

// Filter is a struct that represents a filter on a data stream.
type Filter[I any] struct {
	predicate PredicateFunc[I] // The predicate function used for filtering the input data.
}

// NewFilter creates a new instance of Filter with the given predicate function and error channel.
func NewFilter[I any](predicate PredicateFunc[I]) (Filter[I], error) {
	if predicate == nil {
		return Filter[I]{}, errors.New("predicate func is nil")
	}

	return Filter[I]{
		predicate: predicate,
	}, nil
}

// Transform applies the filter operation on the input channel and returns the output channel.
// The filter only passes the data to out that satisfies the predicate function.
func (f Filter[I]) Transform(in <-chan I, _ chan<- pipeline.Event) <-chan I {
	out := make(chan I)
	go func() {
		defer close(out)
		for v := range in {
			if f.predicate(v) {
				out <- v
			}
		}
	}()
	return out
}
