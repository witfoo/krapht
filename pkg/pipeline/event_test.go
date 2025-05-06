package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/mock"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (chan<- pipeline.Event, pipeline.Event)
		expected    bool
		description string
	}{
		{
			name: "send_to_nil_channel",
			setupFunc: func() (chan<- pipeline.Event, pipeline.Event) {
				return nil, mock.NewEvent(pipeline.EventLog, "test event")
			},
			expected:    false,
			description: "Sending to nil channel should return false",
		},
		{
			name: "send_to_buffered_channel_with_space",
			setupFunc: func() (chan<- pipeline.Event, pipeline.Event) {
				eventC := make(chan pipeline.Event, 1)
				return eventC, mock.NewEvent(pipeline.EventLog, "test event")
			},
			expected:    true,
			description: "Sending to a buffered channel with space should succeed",
		},
		{
			name: "send_to_full_buffered_channel",
			setupFunc: func() (chan<- pipeline.Event, pipeline.Event) {
				eventC := make(chan pipeline.Event, 1)
				// Fill the channel
				eventC <- mock.NewEvent(pipeline.EventLog, "filling event")
				return eventC, mock.NewEvent(pipeline.EventLog, "test event")
			},
			expected:    false,
			description: "Sending to a full buffered channel should return false",
		},
		{
			name: "send_to_unbuffered_channel_with_receiver",
			setupFunc: func() (chan<- pipeline.Event, pipeline.Event) {
				eventC := make(chan pipeline.Event)
				event := mock.NewEvent(pipeline.EventLog, "test event")

				// Use a channel to synchronize between the sender and receiver
				ready := make(chan struct{})
				sent := make(chan struct{})

				// Start a goroutine to receive from the channel
				go func() {
					// Signal we're ready to receive
					close(ready)

					// Block waiting for the event
					<-eventC

					// Signal that we've received the event
					close(sent)
				}()

				// Wait for the receiver goroutine to be ready
				<-ready

				return eventC, event
			},
			expected:    true,
			description: "Sending to an unbuffered channel with a receiver should succeed",
		},
		{
			name: "send_to_unbuffered_channel_without_receiver",
			setupFunc: func() (chan<- pipeline.Event, pipeline.Event) {
				eventC := make(chan pipeline.Event)
				return eventC, mock.NewEvent(pipeline.EventLog, "test event")
			},
			expected:    false,
			description: "Sending to an unbuffered channel without a receiver should return false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventC, event := tt.setupFunc()
			result := pipeline.SendEvent(eventC, event)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}
