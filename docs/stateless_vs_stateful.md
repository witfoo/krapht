# Understanding Stateless Triple Extractors and Stateful Rule Engines

These two architectural components represent key design approaches within KraphT's event-driven pipeline:

## Stateless Triple Extractors

Triple extractors convert raw data into semantic triples (subject-predicate-object relationships). They are designed to be **stateless** because:

1. **Independent Processing:** Each extractor processes a single input document or data point without needing context from previous extractions
2. **Horizontal Scalability:** Multiple instances can run in parallel without coordination
3. **Fault Tolerance:** If an extractor fails, it can be restarted without losing state
4. **Deployment Flexibility:** Can be deployed as ephemeral containers or serverless functions
5. **No Memory Constraints:** Suitable for resource-constrained environments like edge devices

**Example scenario:** A triple extractor might process a sensor reading from an ocean buoy, converting "Buoy-7 recorded temperature 21.5°C at 2025-05-03T14:30Z" into triples like:

- (Buoy-7, hasReading, Reading-1234)
- (Reading-1234, measuresProperty, Temperature)
- (Reading-1234, hasValue, "21.5°C")
- (Reading-1234, recordedAt, "2025-05-03T14:30Z")

This processing doesn't require knowledge of previous readings.

## Stateful Rule Engines

Rule engines analyze streams of triples to identify patterns and infer new relationships. They are **stateful** because:

1. **Temporal Context:** Need to retain historical triple information to identify patterns over time
2. **Aggregation:** Often combine information from multiple triple streams
3. **Pattern Recognition:** Must track partial pattern matches until they complete or expire
4. **Inference State:** Maintain derived knowledge and relationships
5. **Windowing Operations:** Support time or count-based windows for analysis

**Example scenario:** A rule engine might track ocean temperature readings over time to identify a warming trend:

- It maintains state about previous temperature readings from multiple buoys
- It applies a rule like "if temperature increases by more than 2°C over 24 hours across 3+ buoys in the same region, flag as anomalous warming event"
- This requires maintaining state about previous readings, partially matched patterns, and temporal context

## Architectural Advantages

This separation provides several benefits in a NATS-based architecture:

1. **Independent Scaling:** Scale extractors based on input volume and rule engines based on complexity
2. **Resource Optimization:** Deploy lightweight extractors on edge devices and more resource-intensive rule engines on more powerful nodes
3. **Resilience:** Stateless extractors can restart without data loss, while stateful rule engines can recover from durable streams
4. **Progressive Enhancement:** Basic extraction can happen locally, while advanced rule processing can occur where more resources are available
5. **Data Flow Optimization:** Extractors can run close to data sources, reducing network requirements for raw data

This design pattern aligns perfectly with NATS' streaming model, where stateless extractors publish triples that stateful rule engines subscribe to and process based on configured patterns.
