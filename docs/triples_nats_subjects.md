# Efficient Subject Routing in Triple-Based Architecture

The triple-based architecture enables exceptionally efficient subject routing within NATS. Let's explore this benefit in greater detail:

## Subject Hierarchy Design for Triples

```
triple.<predicate>.<source>
```

This hierarchical structure provides several significant advantages:

### 1. Predicate-Based Subscription

```go
// Subscribe to all authentication events regardless of source
_, err := js.PullSubscribe("triple.authenticated_to.>", "AUTH_MONITOR")

// Subscribe to all network connections
_, err := js.PullSubscribe("triple.connected_to.>", "NETWORK_MONITOR") 

// Subscribe to all process executions from Windows sources
_, err := js.PullSubscribe("triple.executes.windows_security", "WINDOWS_PROCESS_MONITOR")
```

This approach allows consumers to:
- Focus on specific relationship types
- Process only the events relevant to their function
- Scale independently based on subject volume

### 2. Fanout Optimizations

NATS efficiently handles subject-based fanout:

```go
// Different consumers can process the same triples differently
type ProcessExecutionMonitor struct {
    js nats.JetStreamContext
    kv nats.KeyValue
    mu sync.RWMutex
}

func (m *ProcessExecutionMonitor) Start(ctx context.Context) error {
    // Subscribe only to process execution events
    sub, err := m.js.PullSubscribe("triple.executes.>", "PROCESS_MONITOR")
    if err != nil {
        return fmt.Errorf("subscribing to process events: %w", err)
    }
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                msgs, err := sub.Fetch(100, nats.MaxWait(1*time.Second))
                if err == nats.ErrTimeout {
                    continue
                }
                if err != nil {
                    continue
                }
                
                for _, msg := range msgs {
                    var triple Triple
                    if err := json.Unmarshal(msg.Data, &triple); err != nil {
                        msg.Ack()
                        continue
                    }
                    
                    // Process execution-specific logic
                    if err := m.processExecutionEvent(ctx, triple); err != nil {
                        // Handle error but ensure message is acknowledged
                    }
                    
                    msg.Ack()
                }
            }
        }
    }()
    
    return nil
}
```

### 3. Load Distribution

This subject hierarchy enables natural load distribution:

```go
// For high-volume predicates, distribute processing across multiple consumers
func NewProcessMonitorPool(ctx context.Context, nc *nats.Conn, poolSize int) ([]*ProcessExecutionMonitor, error) {
    js, err := nc.JetStream()
    if err != nil {
        return nil, fmt.Errorf("getting jetstream context: %w", err)
    }
    
    // Create the consumer with explicit ack policy
    _, err = js.AddConsumer("TRIPLES", &nats.ConsumerConfig{
        Durable:       "PROCESS_MONITOR",
        AckPolicy:     nats.AckExplicitPolicy,
        FilterSubject: "triple.executes.>",
        MaxAckPending: 1000,
        AckWait:       30 * time.Second,
    })
    if err != nil && !errors.Is(err, nats.ErrConsumerNameExists) {
        return nil, fmt.Errorf("creating consumer: %w", err)
    }
    
    // Create pool of monitors sharing the same consumer
    monitors := make([]*ProcessExecutionMonitor, poolSize)
    for i := 0; i < poolSize; i++ {
        monitor, err := NewProcessExecutionMonitor(ctx, js)
        if err != nil {
            return nil, fmt.Errorf("creating monitor %d: %w", i, err)
        }
        monitors[i] = monitor
    }
    
    return monitors, nil
}
```

### 4. Selective Filtering

The subject hierarchy enables precise filtering without wasting resources:

```go
// Filter by specific predicates and sources
// filepath: /Users/coby/Code/witfoo-dev/common/triple/processor/routing.go
package processor

import (
    "context"
    "fmt"
    "sync"

    "github.com/nats-io/nats.go"
)

// RouteTriples routes triples to appropriate consumers based on configuration
func RouteTriples(ctx context.Context, js nats.JetStreamContext, routing map[string]string) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(routing))
    
    // Start a goroutine for each routing rule
    for subject, target := range routing {
        wg.Add(1)
        go func(subject, target string) {
            defer wg.Done()
            
            sub, err := js.PullSubscribe(subject, fmt.Sprintf("ROUTER_%s", subject))
            if err != nil {
                errCh <- fmt.Errorf("subscribing to %s: %w", subject, err)
                return
            }
            
            for {
                select {
                case <-ctx.Done():
                    return
                default:
                    msgs, err := sub.Fetch(100, nats.MaxWait(1*time.Second))
                    if err == nats.ErrTimeout {
                        continue
                    }
                    if err != nil {
                        // Log error but continue
                        continue
                    }
                    
                    for _, msg := range msgs {
                        // Route message to target subject
                        _, err := js.Publish(target, msg.Data)
                        if err != nil {
                            // Log error but continue
                        }
                        msg.Ack()
                    }
                }
            }
        }(subject, target)
    }
    
    // Handle errors
    go func() {
        wg.Wait()
        close(errCh)
    }()
    
    // Collect any errors
    for err := range errCh {
        return err // Return first error
    }
    
    return nil
}
```

### 5. Dynamic Subject Construction

Triple-based architecture allows for dynamic subject routing based on message content:

```go
// GenerateTripleSubject generates a NATS subject for a triple
// filepath: /Users/coby/Code/witfoo-dev/common/triple/subject/generator.go
package subject

import (
    "strings"

    "github.com/witfoo/common/triple"
)

// GenerateTripleSubject creates a subject for publishing a triple
func GenerateTripleSubject(t triple.Triple, options ...Option) string {
    // Start with base subject
    subject := "triple"
    
    // Apply default configuration
    config := defaultConfig()
    
    // Apply options
    for _, opt := range options {
        opt(config)
    }
    
    // Add predicate level
    if config.includePredicate {
        // Sanitize predicate for NATS subject compatibility
        predicate := sanitizeForSubject(t.Predicate)
        subject = subject + "." + predicate
    }
    
    // Add source level
    if config.includeSource && t.Source != "" {
        // Sanitize source for NATS subject compatibility
        source := sanitizeForSubject(t.Source)
        subject = subject + "." + source
    }
    
    // Add custom levels from config
    for _, level := range config.additionalLevels {
        if level.value != "" {
            subject = subject + "." + sanitizeForSubject(level.value)
        }
    }
    
    return subject
}

// sanitizeForSubject ensures a string is valid for use in a NATS subject
func sanitizeForSubject(s string) string {
    // Replace invalid characters with underscores
    replacer := strings.NewReplacer(
        " ", "_",
        ".", "_",
        ">", "_",
        "*", "_",
    )
    return replacer.Replace(s)
}

// Configuration options
type config struct {
    includePredicate bool
    includeSource    bool
    additionalLevels []subjectLevel
}

type subjectLevel struct {
    value string
}

func defaultConfig() *config {
    return &config{
        includePredicate: true,
        includeSource:    true,
    }
}

// Option configures subject generation
type Option func(*config)

// WithoutPredicate removes predicate from subject
func WithoutPredicate() Option {
    return func(c *config) {
        c.includePredicate = false
    }
}

// WithoutSource removes source from subject
func WithoutSource() Option {
    return func(c *config) {
        c.includeSource = false
    }
}

// WithAdditionalLevel adds a custom level to the subject
func WithAdditionalLevel(value string) Option {
    return func(c *config) {
        c.additionalLevels = append(c.additionalLevels, subjectLevel{value: value})
    }
}
```

### 6. Stream Partitioning

For high-volume deployments, predicate-based routing enables efficient stream partitioning:

```go
// ConfigureTripleStreams sets up optimized streams for different triple types
// filepath: /Users/coby/Code/witfoo-dev/common/triple/stream/configurator.go
package stream

import (
    "context"
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
)

// StreamConfig defines configuration for a triple stream
type StreamConfig struct {
    Name       string
    Subjects   []string
    MaxAge     time.Duration
    Storage    nats.StorageType
    Replicas   int
    Retention  nats.RetentionPolicy
    Discard    nats.DiscardPolicy
    MaxBytes   int64
}

// ConfigureTripleStreams creates optimized streams for different triple predicates
func ConfigureTripleStreams(ctx context.Context, js nats.JetStreamContext, configs []StreamConfig) error {
    for _, cfg := range configs {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Create or update stream
            _, err := js.AddStream(&nats.StreamConfig{
                Name:      cfg.Name,
                Subjects:  cfg.Subjects,
                MaxAge:    cfg.MaxAge,
                Storage:   cfg.Storage,
                Replicas:  cfg.Replicas,
                Retention: cfg.Retention,
                Discard:   cfg.Discard,
                MaxBytes:  cfg.MaxBytes,
            })
            if err != nil {
                return fmt.Errorf("creating stream %s: %w", cfg.Name, err)
            }
        }
    }
    
    return nil
}

// Default configurations for common triple types
func DefaultTripleStreamConfigs() []StreamConfig {
    return []StreamConfig{
        {
            Name:      "AUTHENTICATION",
            Subjects:  []string{"triple.authenticated_to.>", "triple.uses_auth_type.>"},
            MaxAge:    30 * 24 * time.Hour, // 30 days
            Storage:   nats.FileStorage,
            Replicas:  1,
            Retention: nats.LimitsPolicy,
            Discard:   nats.DiscardOld,
            MaxBytes:  -1, // Unlimited
        },
        {
            Name:      "PROCESS_EXECUTION",
            Subjects:  []string{"triple.executes.>"},
            MaxAge:    14 * 24 * time.Hour, // 14 days
            Storage:   nats.FileStorage,
            Replicas:  1,
            Retention: nats.LimitsPolicy,
            Discard:   nats.DiscardOld,
            MaxBytes:  -1, // Unlimited
        },
        {
            Name:      "NETWORK_CONNECTIONS",
            Subjects:  []string{"triple.connected_to.>"},
            MaxAge:    7 * 24 * time.Hour, // 7 days
            Storage:   nats.FileStorage,
            Replicas:  1,
            Retention: nats.LimitsPolicy,
            Discard:   nats.DiscardOld,
            MaxBytes:  -1, // Unlimited
        },
        // Default catch-all stream for other triples
        {
            Name:      "TRIPLES_OTHER",
            Subjects:  []string{"triple.>"},
            MaxAge:    7 * 24 * time.Hour, // 7 days
            Storage:   nats.FileStorage,
            Replicas:  1,
            Retention: nats.LimitsPolicy,
            Discard:   nats.DiscardOld,
            MaxBytes:  -1, // Unlimited
        },
    }
}
```

This approach to subject routing provides exceptional flexibility, scalability, and performance for your triple-based security incident graph architecture on NATS JetStream.