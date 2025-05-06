// Package sink provides data sink functions used in the pipeline package.
package sink

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	pl "github.com/witfoo/krapht/pkg/pipeline"
)

// ackable is a type that has an Ack method
type ackable interface {
	Ack() error
}

// Static check that Nats implements the sink interface
var _ pl.Sink[pl.DataRawReadable] = (*NatsStream)(nil)

// ErrNatsSink is the error returned by the Nats sink
var ErrNatsSink = errors.New("sink error")

// NatsSinkPrefix is the prefix for the Nats sink error
const NatsSinkPrefix = "nats stream sink"

// NatsStream implements the source interface
type NatsStream struct {
	js      jetstream.JetStream // nats jetstream
	subject string              // nats jetstream subject
}

// NewNatsStream creates a new Nats sink
func NewNatsStream(js jetstream.JetStream, subject string) (NatsStream, error) {
	var empty NatsStream

	if js == nil {
		return empty, fmt.Errorf("%s: %w: nats connection is nil", NatsSinkPrefix, ErrNatsSink)
	}

	if subject == "" {
		return empty, fmt.Errorf("%s: %w: subject is empty", NatsSinkPrefix, ErrNatsSink)
	}

	return NatsStream{
		js:      js,
		subject: subject,
	}, nil
}

// Load connects to the nats stream and sends messages to the stream
func (b NatsStream) Load(in <-chan pl.DataRawReadable, eventC chan<- pl.Event) {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// connect in channel to nats producer
		for drr := range in {
			p, err := drr.Data().Read()
			if err != nil {
				eventC <- pl.NewErrorEvent(
					"failed to read data",
					err,
					true)
				continue
			}
			_, err = b.js.Publish(ctx, b.subject, p)
			if err != nil {
				eventC <- pl.NewErrorEvent(
					"failed to publish message to nats",
					err,
					true)
				continue
			}

			// ack the message if it is ackable
			if a, ok := drr.Raw().(ackable); ok {
				if err := a.Ack(); err != nil {
					eventC <- pl.NewErrorEvent(
						"failed to ack message",
						err,
						true)
					continue
				}
			}
		}
	}()
}
