package pipeline

import (
	"sync"
	"sync/atomic"
)

// EventCallback is a function called when an event is received.
type EventCallback func(Event)

// TypedEventCallback pairs an event type with a callback function.
type TypedEventCallback struct {
	eventType EventType
	callback  EventCallback
}

// EventCollectorOption represents a functional option for configuring EventCollector.
type EventCollectorOption func(*EventCollector)

// WithWorkers configures the number of worker goroutines processing events.
func WithWorkers(workers int) EventCollectorOption {
	return func(c *EventCollector) {
		if workers > 0 {
			c.workers = workers
		}
	}
}

// WithBufferSize configures the buffer size for the event channel.
func WithBufferSize(size int) EventCollectorOption {
	return func(c *EventCollector) {
		if size > 0 {
			c.bufferSize = size
		}
	}
}

// WithCallback adds a callback that will be called for all events.
func WithCallback(callback EventCallback) EventCollectorOption {
	return func(c *EventCollector) {
		if callback != nil {
			c.callbacks = append(c.callbacks, callback)
		}
	}
}

// WithTypedCallback adds a callback for a specific event type.
func WithTypedCallback(eventType EventType, callback EventCallback) EventCollectorOption {
	return func(c *EventCollector) {
		if callback != nil {
			c.typedCallbacks = append(c.typedCallbacks, TypedEventCallback{
				eventType: eventType,
				callback:  callback,
			})
		}
	}
}

// EventCollector captures events and provides callback processing.
// It implements thread-safe event collection with proper lifecycle management.
type EventCollector struct {
	bufferSize     int
	workers        int
	callbacks      []EventCallback
	typedCallbacks []TypedEventCallback
	wg             sync.WaitGroup
	eventChan      chan Event
	isOpen         atomic.Bool
}

// NewEventCollector creates a new event collector with default settings.
// It accepts optional EventCollectorOption functions to configure the collector.
func NewEventCollector(opts ...EventCollectorOption) *EventCollector {
	collector := &EventCollector{
		bufferSize:     100, // Default buffer size
		workers:        1,   // Default single worker
		callbacks:      make([]EventCallback, 0),
		typedCallbacks: make([]TypedEventCallback, 0),
	}

	// Apply all options
	for _, opt := range opts {
		opt(collector)
	}

	return collector
}

// Collect returns a channel that collects events.
// The context controls the lifecycle of the collection process.
// It starts worker goroutines to process events concurrently.
// It is up to caller to close the Collector when done.
func (c *EventCollector) Collect() chan<- Event {
	// Check if the collector is already open
	if c.isOpen.Load() {
		return c.eventChan
	}

	// Create a buffered channel for event collection
	eventChan := make(chan Event, c.bufferSize)
	// Mark the collector as open
	c.isOpen.Store(true)
	c.eventChan = eventChan

	// Start worker goroutines to process events
	for range c.workers {
		// Increment the wait group counter for each worker
		c.wg.Add(1)

		// Start a goroutine for each worker
		go func() {
			// Ensure the goroutine signals completion
			defer c.wg.Done()

			// Process events until context is done
			// or the event channel is closed
			for event := range eventChan {
				// Process the event
				c.processEvent(event)
			}

		}()
	}
	// Return the event channel for sending events
	return eventChan
}

// processEvent handles a single event by applying callbacks.
// It's called for each event received by a worker goroutine.
func (c *EventCollector) processEvent(event Event) {
	// Apply general callbacks
	for _, callback := range c.callbacks {
		callback(event)
	}

	// Apply type-specific callbacks
	eventType := event.Type()
	for _, typed := range c.typedCallbacks {
		if typed.eventType == eventType {
			typed.callback(event)
		}
	}
}

// Close the event channel to signal all workers to stop
func (c *EventCollector) Close() {
	// Check if the collector is already closed
	if !c.isOpen.Load() {
		return
	}
	// Mark the collector as closed
	c.isOpen.Store(false)
	// Close the event channel to signal all workers to stop
	close(c.eventChan)
	// Wait for all workers to finish processing
	c.wg.Wait()
}
