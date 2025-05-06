package pipeline_test

import (
	"testing"

	"github.com/witfoo/krapht/pkg/pipeline"
)

func TestNewMetricEvent(t *testing.T) {
	name := "test_metric"
	value := 42.0
	labels := map[string]string{"env": "test", "service": "api"}

	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)

	// Check that we get the correct type
	_, ok := event.(pipeline.Measurable)
	if !ok {
		t.Fatal("Expected NewMetricEvent to return an event implementing Measurable")
	}

}

func TestMetricEventType(t *testing.T) {
	name := "test_metric"
	value := 10.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	if event.Type() != pipeline.EventMetric {
		t.Errorf("Expected event type %v, got %v", pipeline.EventMetric, event.Type())
	}
}

func TestMetricEventString(t *testing.T) {
	name := "test_metric"
	value := 10.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	if event.String() != name {
		t.Errorf("Expected string representation %s, got %s", name, event.String())
	}
}

func TestMetricEventName(t *testing.T) {
	name := "test_metric"
	value := 10.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	measurable := event.(pipeline.Measurable)
	if measurable.Name() != name {
		t.Errorf("Expected name %s, got %s", name, measurable.Name())
	}
}

func TestMetricEventValue(t *testing.T) {
	name := "test_metric"
	value := 42.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	measurable := event.(pipeline.Measurable)
	if measurable.Value() != value {
		t.Errorf("Expected value %f, got %f", value, measurable.Value())
	}
}

func TestMetricEventLabels(t *testing.T) {
	name := "test_metric"
	value := 10.0
	labels := map[string]string{"env": "test", "service": "api"}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	measurable := event.(pipeline.Measurable)

	returnedLabels := measurable.Labels()
	if len(returnedLabels) != len(labels) {
		t.Errorf("Expected %d labels, got %d", len(labels), len(returnedLabels))
	}

	for k, v := range labels {
		if returnedLabels[k] != v {
			t.Errorf("Expected label %s to be %s, got %s", k, v, returnedLabels[k])
		}
	}
}

func TestMetricEventMetricType(t *testing.T) {
	name := "test_metric"
	value := 10.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)
	measurable := event.(pipeline.Measurable)

	if measurable.MetricType() != string(pipeline.MetricTypeCounter) {
		t.Errorf("Expected metric type %v, got %v", pipeline.MetricTypeCounter, measurable.MetricType())
	}
}

func TestMeasurableInterface(t *testing.T) {
	// Test that the event returned by NewMetricEvent implements the Measurable interface
	name := "test_metric"
	value := 10.0
	labels := map[string]string{}
	event := pipeline.NewMetricEvent(name, value, labels, pipeline.MetricTypeCounter)

	_, ok := event.(pipeline.Measurable)
	if !ok {
		t.Error("Event created with NewMetricEvent does not implement Measurable")
	}
}
