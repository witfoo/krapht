package flow_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/flow"
)

func TestFilterMap_Transform(t *testing.T) {
	// Create a new instance of FilterMap
	filterMap, err := flow.NewFilterMap(
		func(i int) bool { return i%2 == 0 },                        // Predicate function: filter even numbers
		func(i int) (string, error) { return strconv.Itoa(i), nil }, // Transform function: convert int to string
	)
	assert.NoError(t, err)

	// Create an input channel with some test data
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()

	eventC := make(chan pipeline.Event)
	defer close(eventC)

	// Apply the filter and map operation
	output := filterMap.Transform(input, eventC)

	// Assert the output values
	expectedOutput := []string{"2", "4"}
	actualOutput := make([]string, 0)
	for o := range output {
		actualOutput = append(actualOutput, o)
	}

	assert.Equal(t, expectedOutput, actualOutput, "Output values should match the expected values")
}
