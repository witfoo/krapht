package flow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/flow"
)

func TestFilter_Transform(t *testing.T) {

	eventC := make(chan pipeline.Event)
	defer close(eventC)

	// Create a channel to receive input data
	input := make(chan int)

	// Create a channel to receive filtered output data
	filterFlow, err := flow.NewFilter(func(in int) bool {
		return in%2 == 0 // Filter even numbers
	})
	assert.NoError(t, err)

	outputC := filterFlow.Transform(input, eventC)

	// Send some test data to the input channel
	go func() {
		input <- 1
		input <- 2
		input <- 3
		input <- 4
		close(input)
	}()

	// Create a slice to store the filtered output
	var result []int

	// Receive the filtered output data from the output channel
	for val := range outputC {
		result = append(result, val)
	}

	// Assert that the filtered output matches the expected result
	expected := []int{2, 4}
	assert.Equal(t, expected, result)
}
