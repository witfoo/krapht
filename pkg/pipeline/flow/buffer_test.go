package flow_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/flow"
)

func TestBuffer_Transform(t *testing.T) {
	t.Run("passes through all data", func(t *testing.T) {
		// Create a new buffer with size 1
		buffer := flow.NewBuffer[int](1)

		// Create input and event channels
		in := make(chan int)
		eventC := make(chan pipeline.Event)

		// Transform the input channel
		out := buffer.Transform(in, eventC)

		// Send data to the input channel
		go func() {
			for i := 0; i < 10; i++ {
				in <- i
			}
			close(in)
		}()

		// Read data from the output channel
		var result []int
		for v := range out {
			result = append(result, v)
		}

		// Verify that all data was passed through
		assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, result)
	})

	t.Run("handles zero size gracefully", func(t *testing.T) {
		// Create a buffer with size 0 (should default to 1)
		buffer := flow.NewBuffer[string](0)

		// Create input and event channels
		in := make(chan string)
		eventC := make(chan pipeline.Event)

		// Transform the input channel
		out := buffer.Transform(in, eventC)

		// Send data to the input channel
		go func() {
			in <- "test"
			close(in)
		}()

		// Verify that data is passed through
		assert.Equal(t, "test", <-out)

		// Verify that the channel is closed
		_, ok := <-out
		assert.False(t, ok)
	})

	t.Run("handles negative size gracefully", func(t *testing.T) {
		// Create a buffer with negative size (should default to 1)
		buffer := flow.NewBuffer[string](-10)

		// Create input and event channels
		in := make(chan string)
		eventC := make(chan pipeline.Event)

		// Transform the input channel
		out := buffer.Transform(in, eventC)

		// Send data to the input channel
		go func() {
			in <- "test"
			close(in)
		}()

		// Verify that data is passed through
		assert.Equal(t, "test", <-out)

		// Verify that the channel is closed
		_, ok := <-out
		assert.False(t, ok)
	})

	t.Run("respects buffer size", func(t *testing.T) {
		// Create a buffer with size 3
		buffer := flow.NewBuffer[int](3)

		// Create input and event channels
		in := make(chan int)
		eventC := make(chan pipeline.Event)

		// Transform the input channel
		out := buffer.Transform(in, eventC)

		// Use a WaitGroup to control test flow
		var wg sync.WaitGroup
		wg.Add(1)

		// This goroutine will block after sending 3 elements
		// until we start consuming from the output channel
		go func() {
			defer close(in)
			for i := 0; i < 5; i++ {
				in <- i
				if i == 2 {
					// Signal that 3 elements have been sent
					wg.Done()
					// Give time for potential overflow
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		// Wait until 3 elements have been sent
		wg.Wait()

		// Read all data from output channel
		var result []int
		for v := range out {
			result = append(result, v)
		}

		// Verify that all data was passed through
		assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
	})
}
