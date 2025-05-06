package source

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/witfoo/krapht/pkg/pipeline"
)

// Ensure that HTTPLog implements the Readable interface.
var _ pipeline.Readable = (*HTTPLog)(nil)

// HTTPLog is a simple struct that implements the Readable interface.
type HTTPLog struct {
	addr string
	log  []byte
	id   uuid.UUID
}

// NewHTTPLog creates a new HTTPLog with the given log and address.
// It creates a new UUID for the log.
func NewHTTPLog(log []byte, addr string) (HTTPLog, error) {
	if log == nil {
		return HTTPLog{}, errors.New("http log error: log is nil")
	}
	if addr == "" {
		addr = "unknown"
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return HTTPLog{}, err
	}

	return HTTPLog{
		addr: addr,
		log:  log,
		id:   id,
	}, err
}

// Read returns the log and nil error.
func (h HTTPLog) Read() ([]byte, error) {
	return h.log, nil
}

// Addr returns the address of the HTTPLog.
func (h HTTPLog) Addr() string {
	return h.addr
}

// ID returns the creation time of the HTTPLog.
func (h HTTPLog) ID() uuid.UUID {
	return h.id
}

// HTTPConfig is the configuration for the HTTP source.
type HTTPConfig struct {
	Addr         string
	Endpoint     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Ensure that HTTP implements Source interface.
var _ pipeline.Source[HTTPLog] = (*HTTPServer)(nil)

// HTTPServer is a struct that represents an HTTP source.
type HTTPServer struct {
	addr         string
	endpoint     string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewHTTPServer creates a new HTTP source with the given configuration.
func NewHTTPServer(conf HTTPConfig) (HTTPServer, error) {
	if conf.Addr == "" {
		conf.Addr = ":8008"
	}

	if conf.Endpoint == "" {
		conf.Endpoint = "/log"
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 5 * time.Second
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 5 * time.Second
	}

	return HTTPServer{
		addr:         conf.Addr,
		endpoint:     conf.Endpoint,
		readTimeout:  conf.ReadTimeout,
		writeTimeout: conf.WriteTimeout,
	}, nil
}

// Extract creates a new HTTP server and returns a channel of RawReadable.
func (h HTTPServer) Extract(ctx context.Context, eventC chan<- pipeline.Event) <-chan HTTPLog {
	out := make(chan HTTPLog)

	handler := &logHandler{
		outC:   out,
		eventC: eventC,
	}

	mux := http.NewServeMux()
	mux.Handle("POST "+h.endpoint, handler)

	server := &http.Server{
		Addr:         h.addr,
		Handler:      mux,
		ReadTimeout:  h.readTimeout,
		WriteTimeout: h.writeTimeout,
	}

	go func() {
		defer close(out)
		// Listen for the context's done signal in a separate goroutine
		go func() {
			<-ctx.Done()
			if err := server.Shutdown(context.Background()); err != nil {
				eventC <- pipeline.NewErrorEvent(
					"HTTP source server shutdown error",
					err,
					true)
			}
		}()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				eventC <- pipeline.NewErrorEvent(
					"HTTP source server error",
					err,
					true)
			}
		}
	}()

	return out
}

type logHandler struct {
	outC   chan<- HTTPLog
	eventC chan<- pipeline.Event
}

// ServeHTTP handles the incoming HTTP request and sends the log to the output channel.
func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// make sure body is not empty
	if len(body) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newlineByte := byte('\n')

	// make sure body ends with newline
	if body[len(body)-1] != newlineByte {
		body = append(body, newlineByte)
	}

	// wrap body and send HTTPLog to output channel
	h.outC <- HTTPLog{log: body, addr: r.RemoteAddr}

	// send OK status
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		h.eventC <- pipeline.NewErrorEvent(
			"failed to write response",
			err,
			true)
	}

}
