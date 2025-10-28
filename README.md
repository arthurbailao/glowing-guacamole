# Anomaly Detection with Z-Score

This is a simple anomaly detection system built with Go and Redis.

## Local development

### Prerequisites

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Setup

Install the project's dependencies:

```bash
make setup
```

### Linting

Use [golangci-lint](https://golangci-lint.run/) to lint the project:

```bash
make lint
```

### Testing

Run the tests:

```bash
make test
```

## Running the project

Just run via Docker Compose:

```bash
docker compose up
```

The project will start a Redis instance, the consumer and producer.

### Environment variables

The consumer accepts the following environment variables:

| Name | Description | Default |
| --- | --- | --- |
| REDIS_URL | Redis URL | |
| REDIS_PUBSUB_CHANNEL | Redis pubsub channel | data |
| Z_SCORE_THRESHOLD | Z-score threshold | 3.0 |
| DETECTOR_WINDOW_SIZE | Detector window size | 50 |


And the producer:

| Name | Description | Default |
| --- | --- | --- |
| REDIS_URL | Redis URL | |
| PRODUCER_THROUGHPUT_PER_SECOND | Producer throughput per second | 5 |

## Design decisions

- Warm-up period before detecting anomalies: The consumer waits for a warm-up period before it starts reporting anomalies — basically, it only kicks in after the rolling window is full. This avoids false positives when there isn’t enough data yet to detect anomalies.
- Redis as the message broker: I used Redis because it’s lightweight and easy to set up for a prototype like this. In a production environment, though, I’d go with something that supports message acknowledgments and checkpoints — like Kafka, RabbitMQ, or even SQS — to handle delivery guarantees and fault tolerance.

## Up Next

### Additional tooling

- CI/CD pipeline: Set up a proper pipeline to run tests, lint, build, and push images automatically. That keeps quality consistent and reduces the chance of manual mistakes.
- GitOps deployment: Use something like FluxCD or ArgoCD to watch for new images and deploy them automatically. It makes releases safer and easier to audit.
- Monitoring: Add metrics and alerts around message throughput, CPU, memory usage, and processing latency. It’s important to know if the consumer is keeping up or getting overwhelmed.
- Health check/Readiness probe: Make the consumer report “healthy” only after the warm-up period finishes. This way, during a rolling update, Kubernetes will wait for the new pod to actually be ready before marking the deployment complete. It helps avoid situations where the rollout finishes while the new pod hasn’t processed enough data yet to start detecting anomalies properly.

### Deployment

Kubernetes would fit this project really well. Rolling updates are especially useful here, since we don’t want downtime while deploying new versions. Setting a proper rolling update strategy would make sure there’s always at least one consumer up and running, processing messages and detecting anomalies. That’s also why the health check mentioned above matters.

```yaml
spec:
  strategy:
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
```

### Missing requirements

- Expected load: It’s important to know the expected traffic rate (messages per second, for example). Without that, we can’t really size or tune the system properly. We might either over-engineer it or miss scaling limits.
- Horizontal scalability: For higher throughput, the system should allow sharding messages so multiple consumers can run in parallel. In a multi-tenant setup, sharding by tenant would make a lot of sense — it keeps workloads isolated and makes scaling much easier.
- State persistence: Right now, the consumer keeps the rolling window in memory. If the service restarts or a new pod spins up during a rolling update, that state is lost, and the window has to rebuild from scratch. In a production setup, it would be better to persist the current window (for example, in Redis) so the service can warmup quickly and maintain detection continuity after restarts.
