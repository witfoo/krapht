// Package source provides data source functions used in the pipeline package.
package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/witfoo/krapht/pkg/pipeline"
)

// Static check that Broker implements the source interface
var _ pipeline.Source[JSMsg] = (*BrokerStream)(nil)

// ErrBrokerSource is the error returned by the broker source
var ErrBrokerSource = errors.New("source error")

// BrokerSourcePrefix is the prefix for the broker source error
var BrokerSourcePrefix = "broker stream source"

// BrokerStream implements the source interface
type BrokerStream struct {
	js     jetstream.JetStream      // nats jetstream
	conf   jetstream.ConsumerConfig // nats jetstream consumer config
	stream string                   // nats jetstream stream name
}

// NewBrokerStream creates a new broker source
func NewBrokerStream(js jetstream.JetStream, stream string, conf jetstream.ConsumerConfig) (BrokerStream, error) {

	if js == nil {
		return BrokerStream{}, fmt.Errorf("%s: nats connection is nil %s", BrokerSourcePrefix, ErrBrokerSource)
	}

	if stream == "" {
		return BrokerStream{}, fmt.Errorf("%s: stream name is empty %s", BrokerSourcePrefix, ErrBrokerSource)
	}

	return BrokerStream{
		js:     js,
		stream: stream,
		conf:   conf,
	}, nil
}

// Extract connects to the nats stream and returns a channel of messages
func (b BrokerStream) Extract(ctx context.Context, eventC chan<- pipeline.Event) <-chan JSMsg {
	out := make(chan JSMsg)

	go func() {
		defer close(out)

		// Add a consumer to the stream
		consumer, err := b.js.CreateOrUpdateConsumer(ctx, b.stream, b.conf)
		if err != nil {
			eventC <- pipeline.NewErrorEvent(
				"failed to create or update consumer",
				err,
				false)
			return
		}

		// connect nats consumer msgs to source out channel
		// TODO:  Eval errors for temp vs perm
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgs, err := consumer.Messages(jetstream.PullMaxMessages(10))
				if err != nil {
					eventC <- pipeline.NewErrorEvent(
						"failed to get messages from nats",
						err,
						true)
					continue
				}
				// iterate over messages
				for {
					msg, err := msgs.Next()
					if err != nil {
						eventC <- pipeline.NewErrorEvent(
							"failed to get next message from nats",
							err,
							true)
						continue
					}

					jsMsg, err := NewJSMsg(msg)
					if err != nil {
						eventC <- pipeline.NewErrorEvent(
							"failed to create JSMsg from message",
							err,
							true)
						// TODO: If we do not ack the message, it may be redelivered
						continue
					}

					out <- jsMsg
				}
			}
		}
	}()

	return out
}

// Ensure that JSMsg implements the Readable interface.
var _ pipeline.Readable = (*JSMsg)(nil)

// JSMsg is a struct that wraps a jetstream message interface with a Readable interface.
type JSMsg struct {
	jetstream.Msg // embedding jetstream message interface
}

// NewJSMsg creates a new jsMsg struct.
func NewJSMsg(jsMsg jetstream.Msg) (JSMsg, error) {
	if jsMsg == nil {
		return JSMsg{}, nil
	}

	return JSMsg{
		Msg: jsMsg,
	}, nil
}

func (j JSMsg) Read() ([]byte, error) {
	return j.Data(), nil
}
