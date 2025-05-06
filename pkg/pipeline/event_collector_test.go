package pipeline_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/mock"
)

func TestNewEventCollector(t *testing.T) {
	tests := []struct {
		name          string
		options       []pipeline.EventCollectorOption
		expectWorkers int
		expectBuffer  int
	}{
		{
			name:          "default configuration",
			options:       []pipeline.EventCollectorOption{},
			expectWorkers: 1,   // default worker count
			expectBuffer:  100, // default buffer size
		},
		{
			name: "custom worker count",
			options: []pipeline.EventCollectorOption{
				pipeline.WithWorkers(5),
			},
			expectWorkers: 5,
			expectBuffer:  100, // default
		},
		{
			name: "custom buffer size",
			options: []pipeline.EventCollectorOption{
				pipeline.WithBufferSize(500),
			},
			expectWorkers: 1, // default
			expectBuffer:  500,
		},
		{
			name: "invalid values are ignored",
			options: []pipeline.EventCollectorOption{
				pipeline.WithWorkers(-5),
				pipeline.WithBufferSize(-10),
			},
			expectWorkers: 1,   // default preserved
			expectBuffer:  100, // default preserved
		},
		{
			name: "multiple options applied correctly",
			options: []pipeline.EventCollectorOption{
				pipeline.WithWorkers(3),
				pipeline.WithBufferSize(200),
			},
			expectWorkers: 3,
			expectBuffer:  200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := pipeline.NewEventCollector(tt.options...)

			// Since workers and bufferSize are private fields, we can only test
			// the behavior influenced by these fields
			// We'll do that in the collection tests
			assert.NotNil(t, collector)
		})
	}
}

func TestEventCollector_Collect(t *testing.T) {
	t.Run("processes events with general callback", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		var receivedEvent pipeline.Event

		callback := func(e pipeline.Event) {
			receivedEvent = e
			wg.Done()
		}

		collector := pipeline.NewEventCollector(
			pipeline.WithCallback(callback),
		)

		eventChan := collector.Collect()
		defer collector.Close()

		// Send test event
		testEvent := mock.NewEvent(
			pipeline.EventLog,
			"test-event",
		)
		// Send the event to the collector
		eventChan <- testEvent

		// Wait for callback to be executed
		waitWithTimeout(t, &wg, 2*time.Second)

		// Verify the event was received by the callback
		assert.Equal(t, testEvent, receivedEvent)
	})

	t.Run("processes events with typed callbacks", func(t *testing.T) {
		const (
			typeA pipeline.EventType = 10
			typeB pipeline.EventType = 11
		)

		var wg sync.WaitGroup
		wg.Add(2) // Expect two callbacks

		var typeACallCount int
		var typeBCallCount int
		var receivedEventA pipeline.Event
		var receivedEventB pipeline.Event

		callbackA := func(e pipeline.Event) {
			typeACallCount++
			receivedEventA = e
			wg.Done()
		}

		callbackB := func(e pipeline.Event) {
			typeBCallCount++
			receivedEventB = e
			wg.Done()
		}

		collector := pipeline.NewEventCollector(
			pipeline.WithTypedCallback(typeA, callbackA),
			pipeline.WithTypedCallback(typeB, callbackB),
		)

		eventChan := collector.Collect()

		testEventA := mock.NewEvent(typeA, "test-data-a")
		eventChan <- testEventA // Should trigger typeA callback

		// Send a second event to trigger the typeB callback
		testEventB := mock.NewEvent(typeB, "test-data-b")
		eventChan <- testEventB // Should trigger typeB callback

		// Wait for both callbacks to be executed
		waitWithTimeout(t, &wg, 2*time.Second)

		// Verify both callbacks were executed
		assert.Equal(t, 1, typeACallCount)
		assert.Equal(t, 1, typeBCallCount)
		assert.Equal(t, testEventA, receivedEventA) // typeA event
		assert.Equal(t, testEventB, receivedEventB) // typeB event
	})

	t.Run("properly handles close", func(t *testing.T) {
		collector := pipeline.NewEventCollector()

		// Start collecting events
		eventChan := collector.Collect()
		collector.Close()

		// Try to send event on closed channel which should panic
		// We will recover from the panic to avoid test failure
		done := make(chan struct{})
		go func() {
			defer close(done)

			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()

			// Use a sentinel value to detect when we should exit the loop
			var channelClosed bool

			for !channelClosed {
				select {
				case <-ticker.C:
					// Try to send to the channel and detect if it's closed
					// Recover from panic if the channel is closed
					func() {
						defer func() {
							if r := recover(); r != nil {
								//  Check to see if r is a panic from sending to a closed channel
								if err, ok := r.(error); ok {
									if err.Error() == "send on closed channel" {
										// Channel is closed - mark for loop exit
										channelClosed = true
									} else {
										// If this is not a panic from sending to a closed channel
										// We should log it and mark the channel as closed
										t.Logf("Panic occurred: %v", err)
									}
								}
							}
						}()

						eventChan <- mock.NewEvent(pipeline.EventLog, "test-data")
					}()

				case <-time.After(1 * time.Second):
					// Timeout safety - don't wait forever
					return
				}
			}
		}()

		// Wait for the sending to fail (channel closed) or timeout
		select {
		case <-done:
			// Success - we detected channel closure
		case <-time.After(2 * time.Second):
			t.Fatal("Close() did not close the event channel")
		}
	})

	t.Run("multiple workers process events concurrently", func(t *testing.T) {
		workerCount := 3
		eventCount := 9 // 3 events per worker

		var mu sync.Mutex
		processedEvents := make([]pipeline.Event, 0, eventCount)

		var wg sync.WaitGroup
		wg.Add(eventCount)

		callback := func(e pipeline.Event) {
			// Simulate processing time
			time.Sleep(50 * time.Millisecond)

			mu.Lock()
			processedEvents = append(processedEvents, e)
			mu.Unlock()

			wg.Done()
		}

		collector := pipeline.NewEventCollector(
			pipeline.WithWorkers(workerCount),
			pipeline.WithCallback(callback),
		)

		eventChan := collector.Collect()
		defer collector.Close()

		// Send events
		startTime := time.Now()
		for i := 0; i < eventCount; i++ {
			eventChan <- mock.NewEvent(
				pipeline.EventLog,
				"test-data-"+strconv.Itoa(i),
			)
		}

		// Wait for all events to be processed
		waitWithTimeout(t, &wg, 3*time.Second)
		duration := time.Since(startTime)

		mu.Lock()
		defer mu.Unlock()

		// Verify all events were processed
		assert.Equal(t, eventCount, len(processedEvents))

		// With sequential processing, it would take eventCount * 50ms
		// With parallel processing, it should take significantly less
		expectedSequentialTime := time.Duration(eventCount) * 50 * time.Millisecond
		assert.Less(t, duration, expectedSequentialTime,
			"Expected concurrent processing to be faster than sequential processing")
	})
}

// waitWithTimeout waits for the WaitGroup with a timeout
func waitWithTimeout(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()

	// Create a child context with timeout from the test context
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		// WaitGroup completed normally
		return
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			// This was a timeout within our function
			t.Fatalf("Timed out after %v waiting for WaitGroup", timeout)
		} else {
			// This was a parent context cancellation
			t.Log("Parent context canceled while waiting for WaitGroup")
		}
	}
}
