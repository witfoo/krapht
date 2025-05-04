# Why NATS?

KraphT is built on NATS as its core messaging layer. This choice reflects a need for high-performance, low-latency, event-driven communication between independent microservices across potentially distributed deployments.

## 1. Event-Driven by Design

KraphT operates as a streaming system where data is transformed into triples, enriched by rules, and propagated to various subsystems. NATS supports this model with fast, loosely coupled pub-sub and queue-based messaging.

Each service in the pipeline can independently emit and consume data, supporting:

- Stateless triple extractors
- Stateful rule engines
- Graph output services
- Optional analytics or visualization layers

## 2. Lightweight and High Throughput

NATS provides sub-millisecond latencies and high message throughput with minimal resource usage. It avoids the overhead and complexity of heavier systems like Kafka or MQTT brokers, making it ideal for real-time processing and flexible service scaling.

## 3. JetStream for Durable Streams

NATS JetStream provides:

- Message persistence
- Replay and temporal logic
- Consumer groups for horizontal scaling
- Backpressure and flow control

This enables KraphT to handle both real-time and batch workloads within the same architecture.

## 4. Federation with Hub-Spoke and Mesh Topologies

NATS supports **leaf nodes**, **clusters**, and **gateways**—allowing KraphT to operate in diverse topologies suited for:

### Hub-Spoke Architectures

- Leaf nodes in remote locations (e.g., offshore data collectors or research vessels) can connect to a central hub.
- Local services can process, normalize, and cache data without requiring full connectivity to upstream components.
- When connectivity is restored, leaf nodes stream data to the central system.

### Federated/Regional Clusters

- Each region or site can run an independent KraphT instance with local NATS clusters.
- Gateways can bridge those regions, enabling cross-site rule propagation or triple synchronization.
- This supports data sovereignty, local processing, and coordinated global reasoning.

### Mesh/Multi-Tenant Systems

- Mesh topologies support peer-to-peer communication between independent processing nodes.
- Ideal for collaborative research networks or multi-tenant deployments with isolation guarantees.

## 5. Microservice-Friendly

NATS supports independent, stateless services that can be updated, scaled, or versioned independently. It encourages composability, observability, and isolation between components, all of which align with KraphT’s modular design.

## 6. Deployment Simplicity

NATS is operationally simple:

- Single static binary
- Native support for clustering and TLS
- First-class Kubernetes integration
- Small resource footprint

## Summary

NATS provides more than messaging—it provides the backbone for a resilient, federated, event-driven system. In KraphT, it enables:

- Real-time triple extraction and rule processing
- Durable state synchronization via JetStream
- Scalable deployment via microservice orchestration
- Distributed reasoning across local and global contexts

With support for hub-spoke, cluster, and mesh topologies, NATS empowers KraphT to grow from single-node pipelines to global graph reasoning platforms without changing its architectural foundation.
