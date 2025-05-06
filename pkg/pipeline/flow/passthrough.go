package flow

import "github.com/witfoo/krapht/pkg/pipeline"

// Ensure that Passthrough implements the Flow interface.
var _ pipeline.Flow[any, any] = (*Passthrough[any])(nil)

// Passthrough is a struct that represents a passthrough operation on a data stream.
type Passthrough[I any] struct {
}

// NewPassthrough creates a new Passthrough.
func NewPassthrough[I any]() *Passthrough[I] {

	return &Passthrough[I]{}
}

// Transform applies the passthrough operation on the input channel and returns the output channel.
// The passthrough operation simply forwards the data from the input channel to the output channel.
func (pt Passthrough[I]) Transform(in <-chan I, _ chan<- pipeline.Event) <-chan I {
	out := make(chan I)

	go func() {
		defer close(out)
		for v := range in {
			out <- v
		}
	}()
	return out
}
