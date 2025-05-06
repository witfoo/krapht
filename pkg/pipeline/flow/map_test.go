package flow_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/flow"
)

func TestMap_Transform(t *testing.T) {
	// Create a test input channel
	input := make(chan int)

	// Create an event channel instead of a log channel
	eventC := make(chan pipeline.Event)
	defer close(eventC)

	// Create a Map instance with a transform function
	mapper, err := flow.NewMap(func(in int) (string, error) {
		return strconv.Itoa(in), nil
	})
	assert.NoError(t, err)

	// Start the Transform operation
	output := mapper.Transform(input, eventC)

	// Listen for events
	go func() {
		for event := range eventC {
			t.Logf("Received event: %v", event)
		}
	}()

	// Send test values to the input channel
	go func() {
		input <- 1
		input <- 2
		input <- 3
		close(input)
	}()

	var result []string

	for val := range output {
		result = append(result, val)
	}

	// Assert the expected output values
	expectedOutput := []string{"1", "2", "3"}
	assert.Equal(t, expectedOutput, result)
}
