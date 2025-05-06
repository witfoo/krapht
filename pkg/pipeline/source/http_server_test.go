package source_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/witfoo/krapht/pkg/pipeline"
	"github.com/witfoo/krapht/pkg/pipeline/source"
)

func TestHTTP_Extract(t *testing.T) {
	conf := source.HTTPConfig{
		Addr:         ":8008",
		Endpoint:     "/test",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpInstance, err := source.NewHTTPServer(conf)
	assert.NoError(t, err)

	eventC := make(chan pipeline.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Process logs in a separate goroutine
	go func() {
		for event := range eventC {
			// In a test, we might want to assert on log events or just acknowledge them
			t.Logf("Received event: %v", event)
		}
	}()

	out := httpInstance.Extract(ctx, eventC)

	go func() {
		for log := range out {
			logP, err := log.Read()
			assert.NoError(t, err)
			assert.Equal(t, []byte("test log\n"), logP)
		}
	}()

	// wait for server to start
	time.Sleep(200 * time.Millisecond)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Send a test request
	logData := []byte("test log\n")
	req, err := http.NewRequest("POST", "http://127.0.0.1:8008/test", bytes.NewBuffer(logData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
