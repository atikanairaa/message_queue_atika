# Message Queue Explore (Python)

A clean-code reimplementation of a RabbitMQ pub/sub sample in Python. The original Go project has separate publisher and consumer projects; this version uses the same architecture and behavior in `publisher/` and `consumer/` packages.

## Architecture

- Producer: HTTP API accepts JSON messages and publishes to RabbitMQ fanout exchange.
- Consumer: separate service for EMAIL, SMS, and FCM that binds to fanout exchange and processes messages.
- Domain layer: message entity and use cases.
- Infrastructure layer: RabbitMQ connection and publisher/consumer repository code.

```
[HTTP Client] -> /publish -> [Publisher app] -> RabbitMQ Exchange (fanout)
                               |                         \
                               +-> queue EMAIL -> [Email consumer]
                               +-> queue SMS   -> [SMS consumer]
                               +-> queue FCM   -> [FCM consumer]
```

## Target language choice

Python was selected because:
- RabbitMQ has mature `pika` client
- Fast iteration with clean structure via `FastAPI`
- Concurrency model for I/O-bound messaging fits `async`/threading-friendly design
- Equivalent feature set (HTTP + AMQP) with minimal boilerplate in a single environment

## Setup

1. Create virtual environment

```bash
python -m venv .venv
source .venv/bin/activate   # or .\.venv\Scripts\activate on Windows
pip install -U pip
pip install -e .
```

2. Ensure RabbitMQ is running locally (default `amqp://guest:guest@localhost:5672/%2F`).
3. Configure env vars as needed:

- `RABBITMQ_URL` (default `amqp://guest:guest@localhost:5672/%2F`)
- `EXCHANGE_NAME` (default `notifications`)
- `PORT` (publisher, default `8080`)
- `SERVICE_NAME` (consumer optional default `CONSUMER`)

## Run Publisher

```bash
python publisher/app.py
```

POST to `http://localhost:8080/publish` with JSON:

```json
{
  "order_id": "12345",
  "user_id": "67890",
  "content": "New order received",
  "timestamp": "2025-03-11T10:00:00Z"
}
```

## Run Consumer

Run each service in its own process:

```bash
python consumer/app.py --service EMAIL
python consumer/app.py --service SMS
python consumer/app.py --service FCM
```

## Consumer log output

Each consumer prints structured logs to stdout that include timestamp and level. Example output:

```text
2026-03-29 14:10:00 | INFO | [EMAIL] Bound queue 'EMAIL' to exchange 'notifications'
2026-03-29 14:10:00 | INFO | [EMAIL] Waiting for messages...
2026-03-29 14:10:05 | INFO | [EMAIL] Received message: Message(order_id='12345', user_id='67890', content='New order received', timestamp='2025-03-11T10:00:00Z')
2026-03-29 14:10:05 | INFO | sending email New order received
```

This confirms the consumer is receiving queue messages and logging each processing step.

## Testing

```bash
pytest
```

## Clean code decisions

- Single responsibility: each module has one role (config, repository, usecase, service)
- Fail-fast config validation
- No magic values; env config with defaults and constants
- Structured logging via `loguru`
- Repo methods raise and return explicit errors
- docstrings for public functions and classes
