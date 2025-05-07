/*
Package pipeline provides a simple pipeline implementation.

	The pipeline is composed of a Source, one or more Flows, and Sink.
	Runnable is a pipeline wrapper and manages the pipeline lifecyle.
	Source is the start of the pipeline.
	Flow is a transformation in the pipeline; it can be chained with other flows.
	Sink is the end of the pipeline.

	Methods use an event channel to communicate with the pipeline.
	Events are used to signal errors and other events in the pipeline.
*/
package pipeline

import (
	"context"
)

// Source interface represents the start of the pipeline.
// Extract method takes a context and returns a channel of type T.
// The context is used in Source to close/drain a pipeline.
type Source[T any] interface {
	Extract(ctx context.Context, eventC chan<- Event) (out <-chan T)
}

// Flow interface represents a transformation in the pipeline.
// Transform method takes an input channel and returns an output channel.
type Flow[In any, Out any] interface {
	Transform(in <-chan In, eventC chan<- Event) (out <-chan Out)
}

// FanOut interface represents a fan-out transformation in the pipeline.
// Split method takes an input channel and returns multiple output channels.
type FanOut[T any] interface {
	Split(in <-chan T, eventC chan<- Event, n uint8) []<-chan T
}

// FanIn interface represents a fan-in transformation in the pipeline.
// Merge method takes multiple input channels and returns a single output channel.
type FanIn[T any] interface {
	Merge(ins []<-chan T, eventC chan<- Event) (out <-chan T)
}

// Sink interface represents the end of the pipeline.
// Load method takes an input channel and returns nothing.
type Sink[In any] interface {
	Load(in <-chan In, eventC chan<- Event)
}

// Runnable interface represents a runnable pipeline.
// Run method takes a context and returns a channel of errors.
type Runnable interface {
	Run(ctx context.Context) <-chan Event
}
