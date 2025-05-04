# What Types of Data Are Best Suited for KraphT?

KraphT is optimized for structured, event-like data that can be expressed as semantic triples. While it works well with time-stamped information, it is not limited to traditional time series data.

## 1. Event-Centric, Structured Data

**Definition**: Data that records discrete events or observations with contextual structure.

**Examples**:

- Security logs (e.g., login attempt, network connection)
- Sensor observations (e.g., temperature reading, salinity sample)
- User or system actions (e.g., alert fired, task completed)

**Why it fits**:

- Naturally maps to subject-predicate-object triples
- Enables rule-based inference and pattern matching
- Time is often included, but the key value lies in the facts or relationships

## 2. Stateful Entity Relationships

**Definition**: Descriptive links between actors, objects, and systems.

**Examples**:

- User to role mapping
- Device to location binding
- Host to software inventory
- Sensor to deployment region

**Why it fits**:

- Central to building a useful property graph
- Enables joining streaming events with reference data
- Facilitates inference like "who is affected" or "what else is related"

## 3. Streaming Time Series Data

**Definition**: High-frequency measurements over time, typically numeric and dense.

**Examples**:

- CPU usage per second
- Environmental sampling at fixed intervals
- Financial tick data

**Fit**: Partial

**How to use with KraphT**:

- Convert raw time series into event-driven triggers (e.g., threshold violations)
- Use external systems (e.g., Prometheus, InfluxDB) for storage and aggregation
- Feed summarized or anomalous events into KraphT for graph enrichment

## 4. Unstructured Content

**Definition**: Data that lacks structure without prior processing.

**Examples**:

- Free-text documents
- PDFs, chat logs, images
- Audio transcripts

**Fit**: Not directly compatible

**How to integrate**:

- Pair with external systems that extract structure (e.g., log parsers, LLMs)
- Emit structured facts or tags into KraphT from pre-processing pipelines

---

### Summary Table

| Data Type                  | Fit for KraphT | Notes                                               |
|----------------------------|----------------|-----------------------------------------------------|
| Structured events          | Yes             | Logs, telemetry, sensor outputs                     |
| Entity relationships       | Yes             | Descriptive and contextual graph edges              |
| Dense time series          | Partial         | Best when reduced to event form                     |
| Unstructured documents     | No              | Requires external preprocessing                     |
| Batch ETL outputs          | Yes, if structured | Suitable if they can be expressed as triples     |
