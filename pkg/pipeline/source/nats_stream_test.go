package source_test

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
	"github.com/witfoo/krapht/pkg/pipeline/source"
)

func TestBroker_Extract(t *testing.T) {

	if test := os.Getenv("INTEGRATION_TESTS"); test != "true" {
		t.Skip("Skipping integration test")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start a NATS server in a Docker container
	natsContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "nats:latest",
				Cmd:          []string{"-js", "-m", "8222"},
				ExposedPorts: []string{"4222/tcp", "8222/tcp"},
				WaitingFor: wait.ForAll(
					wait.ForListeningPort("4222/tcp"),
					wait.ForHTTP("/").WithPort("8222/tcp").WithStartupTimeout(20*time.Second),
				),
			},
			Started: true,
		},
	)
	assert.NoError(t, err)
	//defer func() { _ = natsContainer.Terminate(ctx) }()

	// Get the NATS server host and port
	natsHost, err := natsContainer.Host(ctx)
	assert.NoError(t, err)
	natsPort, err := natsContainer.MappedPort(ctx, "4222")
	assert.NoError(t, err)

	// Connect to the NATS server
	nc, err := nats.Connect("nats://" + natsHost + ":" + natsPort.Port())
	assert.NoError(t, err)
	defer nc.Close()

	// Create jetstream
	js, err := jetstream.New(nc)
	assert.NoError(t, err)

	// Create stream
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "test-stream",
		Subjects: []string{"test"},
	})
	assert.NoError(t, err)

	// Create a jetstream consumer config
	conf := jetstream.ConsumerConfig{
		Durable:       "source",
		FilterSubject: "test",
	}

	eventC := make(chan pipeline.Event, 1)

	// Create a new broker source
	b, err := source.NewBrokerStream(js, "test-stream", conf)
	assert.NoError(t, err)

	// Start extracting messages
	msgCh := b.Extract(ctx, eventC)

	// Create headers for the test message
	headers := nats.Header{}
	headers.Add("orgID", "test-org")

	// Create a message with headers
	msg := &nats.Msg{
		Subject: "test",
		Data:    []byte("test message"),
		Header:  headers,
	}

	// Publish a test message with headers to the NATS server
	_, err = js.PublishMsg(ctx, msg)
	assert.NoError(t, err)

	type ackable interface {
		Ack() error
	}

	// Wait for the message to be received
	select {
	case msg := <-msgCh:
		data, err := msg.Read()
		assert.NoError(t, err)
		assert.Equal(t, []byte("test message"), data)
		err = msg.Ack()
		assert.NoError(t, err)
	case event := <-eventC:
		t.Errorf("Received event: %v", event)
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for message")
	}

	_ = natsContainer.Terminate(ctx)

}
