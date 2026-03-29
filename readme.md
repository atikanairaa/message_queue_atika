# Message Queue Atika

## Overview

This repository is a multi-language example for RabbitMQ pub/sub notification processing, including:
- Go implementation:
  - `message_queue_cloning/notification_publisher`
  - `message_queue_cloning/notification_consumer`
- Python implementation:
  - `message_queue_explore/publisher`
  - `message_queue_explore/consumer`

Functionality:
- Publisher accepts HTTP request, publishes to RabbitMQ exchange.
- Consumer(s) subscribe to queue(s) and process notifications:
  - Email
  - SMS
  - FCM

Purpose:
- Learning message queue architecture
- Decoupled distributed service pattern
- RabbitMQ fanout exchange pub/sub

---

## Repo layout

- `readme.md` (root, this generated README)
- `message_queue_cloning/notification_publisher/`
  - Go app with Echo HTTP API
- `message_queue_cloning/notification_consumer/`
  - Go console consumers in `cmd/consumer/{email,sms,fcm}`
- `message_queue_explore/`
  - Python app with equivalent publisher+consumer logic

---

## Prerequisites

1. Go 1.20+ (for Go modules)
2. Python 3.11+ (for Python sample)
3. RabbitMQ server installed and running:
   - default `amqp://guest:guest@localhost:5672/%2F`
4. `git` (for clone)
5. Windows/macOS/Linux + language-specific tools

---

## RabbitMQ setup (quick)

### Windows
- Install Erlang + RabbitMQ
- Run:
  - `rabbitmq-server.bat`
  - optionally service:
    - `rabbitmq-service install`
    - `rabbitmq-service start`

### macOS
- `brew install rabbitmq`
- `brew services start rabbitmq`

### Ubuntu
- `sudo apt install rabbitmq-server -y`
- `sudo systemctl start rabbitmq-server`

---

## Go service: notification_publisher

Path:
- `message_queue_cloning/notification_publisher`

### Build/Install
```powershell
cd message_queue_cloning/notification_publisher
go mod tidy
go mod vendor
```

### Run
```powershell
go run main.go
```

### Test publish:
```bash
curl -X POST http://localhost:8080/publish \
 -H "Content-Type: application/json" \
 -d '{
   "order_id":"12345",
   "user_id":"67890",
   "content":"New order received",
   "timestamp":"2025-03-11T10:00:00Z"
 }'
```

Expected response:
```json
{"code":200,"message":"Message published successfully"}
```

---

## Go consumers: notification_consumer

Path:
- `message_queue_cloning/notification_consumer`

### Build/Install
```powershell
cd message_queue_cloning/notification_consumer
go mod tidy
go mod vendor
```

### Run consumers
- Email:
  `go run cmd/consumer/email/main.go`
- SMS:
  `go run cmd/consumer/sms/main.go`
- FCM:
  `go run cmd/consumer/fcm/main.go`

### Expected logs
- `Received message: {...}`
- `sending email ...` or `sending SMS ...` or `sending FCM ...`

---

## Python version: message_queue_explore

Path:
- `message_queue_explore`

### Setup
```bash
cd message_queue_explore
python -m venv .venv
source .venv/bin/activate     # Windows: .\.venv\Scripts\activate
pip install -U pip
pip install -e .
```

### Config env vars
- `RABBITMQ_URL` (default `amqp://guest:guest@localhost:5672/%2F`)
- `EXCHANGE_NAME` (default `notifications`)
- `PORT` (publisher default `8080`)
- `SERVICE_NAME` (consumer)

### Run publisher
```bash
python publisher/app.py
```

### Run consumer(s)
```bash
python consumer/app.py --service EMAIL
python consumer.app.py --service SMS
python consumer/app.py --service FCM
```

### Publish message
```bash
curl -X POST http://localhost:8080/publish \
  -H "Content-Type: application/json" \
  -d '{"order_id":"12345","user_id":"67890","content":"New order received","timestamp":"2025-03-11T10:00:00Z"}'
```

### Test
```bash
pytest
```

---

## Architecture

- Publisher -> RabbitMQ fanout exchange
- Exchange copies to each queue (EMAIL, SMS, FCM)
- Consumer per channel processes message independently
- Decoupling, scaling, fault isolation

---

## Notes

- For both versions, ensure RabbitMQ is accessible before starting publisher/consumer.
- If using Docker, containerize RabbitMQ + services as needed.
- The sample is designed for learning stage; extension points:
  - durable queues
  - ack/nack handling
  - retry + dead-letter policy
  - persisted storage or monitoring integration

---

## Contributor info

- Author: likely `kusnadi8605` (based on linked samples)
- University semester project
- Focus on clean architecture and message-driven pattern

---

## Quick start in one line (Go)

1. Start RabbitMQ
2. Run publisher + consumers
3. Use curl to send message
4. Observe consumer logs
