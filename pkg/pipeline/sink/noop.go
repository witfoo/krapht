package sink

import "github.com/witfoo/krapht/pkg/pipeline"

var _ pipeline.Sink[any] = (*Noop)(nil)

// Noop is a struct that represents a no-op on a pipline sink.
type Noop struct {
}

// Load just reads input channel.
func (n Noop) Load(in <-chan any, _ chan<- pipeline.Event) {
	// revive:disable-next-line:empty-block
	for range in {
	} // just read the data
}
