package pipeline

// MetricType represents the type of a metric
type MetricType string

const (
	// MetricTypeCounter represents a counter metric
	MetricTypeCounter MetricType = "counter"
	// MetricTypeGauge represents a gauge metric
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeHistogram represents a histogram metric
	MetricTypeHistogram MetricType = "histogram"
	// MetricTypeSummary represents a summary metric
	MetricTypeSummary MetricType = "summary"
)

// Measurable is a specialized Event for metrics
type Measurable interface {
	Event
	// Name returns the name of the metric
	Name() string
	// Value returns the value of the metric
	Value() float64
	// Labels returns the tags associated with the metric
	Labels() map[string]string
	// MetricType returns the type of metric (counter, gauge, histogram, summary)
	MetricType() string
}

var _ Measurable = (*MetricEvent)(nil)

// MetricEvent represents a metric measurement in the pipeline
type MetricEvent struct {
	name       string
	value      float64
	labels     map[string]string
	metricType MetricType
}

// NewMetricEvent creates a new MetricEvent instance
func NewMetricEvent(name string, value float64, labels map[string]string, metricType MetricType) Event {
	return MetricEvent{
		name:       name,
		value:      value,
		labels:     labels,
		metricType: metricType,
	}
}

// Type returns the type of event
func (m MetricEvent) Type() EventType {
	return EventMetric
}

// String returns the string representation of the event
func (m MetricEvent) String() string {
	return m.name
}

// Name returns the name of the metric event
func (m MetricEvent) Name() string {
	return m.name
}

// Value returns the value of the metric event
func (m MetricEvent) Value() float64 {
	return m.value
}

// Labels returns the labels associated with the metric event
func (m MetricEvent) Labels() map[string]string {
	return m.labels
}

// MetricType returns the type of the metric
func (m MetricEvent) MetricType() string {
	return string(m.metricType)
}
