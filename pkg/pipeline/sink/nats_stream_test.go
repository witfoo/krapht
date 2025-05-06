package sink_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/mock"
	"github.com/witfoo/krapht/pkg/pipeline/sink"
)

func TestBrokerSink(t *testing.T) {

	if test := os.Getenv("INTEGRATION_TESTS"); test != "true" {
		t.Skip("Skipping integration test")
	}

	// Create a context
	ctx, canc := context.WithCancel(context.Background())
	defer canc()

	// Start a NATS container
	natsContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "nats:latest",
				ExposedPorts: []string{"4222/tcp", "8222/tcp"},
				Cmd:          []string{"-js", "-m", "8222"},
				WaitingFor: wait.ForAll(
					wait.ForListeningPort("4222/tcp"),
					wait.ForHTTP("/").WithPort("8222/tcp").WithStartupTimeout(20*time.Second),
				),
			},
			Started: true,
		},
	)
	assert.NoError(t, err)

	// Get the NATS container host and port
	natsHost, err := natsContainer.Host(ctx)
	assert.NoError(t, err)
	natsPort, err := natsContainer.MappedPort(ctx, "4222")
	assert.NoError(t, err)

	// Create a connection to the NATS container
	natsConn, err := nats.Connect("nats://" + natsHost + ":" + natsPort.Port())
	assert.NoError(t, err)
	defer natsConn.Close()

	// Create a jetstream stream
	js, err := jetstream.New(natsConn)
	assert.NoError(t, err)

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "test-stream",
		Subjects: []string{"test"},
	})
	assert.NoError(t, err)

	// Create a jetstream consumer
	consumer, err := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "test-durable",
	})
	assert.NoError(t, err)

	// Create log channel
	eventC := make(chan pipeline.Event, 1)
	defer close(eventC)

	// Listen for logs from sink
	go func() {
		for event := range eventC {
			t.Logf("Received event: %v", event)
		}
	}()

	// Create a new Broker sink with test configuration
	sink, err := sink.NewNatsStream(js, "test")
	assert.NoError(t, err)

	// Create a channel to receive the message
	in := make(chan pipeline.DataRawReadable)
	defer close(in)

	// Start the sink
	sink.Load(in, eventC)

	// Create a mock DataRawReadable
	drr := mock.NewDataRawReadableImpl(
		mock.NewReadableImpl([]byte("test raw data")),
		mock.NewReadableImpl([]byte("test data")),
	)
	assert.NoError(t, err)

	// Send a message to the sink in channel
	in <- drr

	// Wait for the message to be received by nats
	msgs, err := consumer.Fetch(1)
	assert.NoError(t, err)
	assert.NoError(t, msgs.Error())
	for msg := range msgs.Messages() {

		assert.Equal(t, "test data", string(msg.Data()))
	}

	_ = natsContainer.Terminate(ctx)

}
