# Safari Payment Service - Merchant-Initiated Payment API

This project implements a **merchant-initiated M-Pesa (Safari) payment service** in Go. It integrates with the [Safaricom Ethiopia Sandbox](https://developer.safaricom.et) to initiate and process payments, storing confirmed transactions in a PostgreSQL database.

---

# Features

- M-Pesa API Integration – Supports merchant-initiated payment requests via Safaricom API.
- PostgreSQL Database – Persists successful payments via Dockerized PostgreSQL.
- Robust Tooling:
  - `zap` for structured logging
  - `ozzo-validation` for input validation
  - `golang-migrate` for schema migrations
  - `sqlc` for type-safe PostgreSQL query generation
  - `errorx` for error tracing and propagation
  - `gin` for ease of middleware setup and routing

---

Project Setup

Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Ngrok](https://ngrok.com/) (for public webhook URL)

---

# Configuration

Create a `.env` file under `config/` directory with the following environment variables:

```env
# Environment
DEBUG=true
DEV=true

# Database
DATABASE_URL=postgres://user:password@localhost:5432/db?sslmode=disable
IDLE_CONN_TIMEOUT=5m

# Migration
MIGRATION_ACTIVE=true
MIGRATION_PATH=internal/query/schemas

# Server
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_TIMEOUT=30s
SERVER_READ_HEADER_TIMEOUT=30s
MAX_RETRIES=3
API_KEY=somevalue

# Safari API
SAFARI_CONSUMER_KEY=your_consumer_key
SAFARI_CONSUMER_SECRET=your_consumer_secret
SAFARI_BUSINESS_SHORT_CODE=short_code form safari sandbox
SAFARI_PASSWORD=your_encoded_password
SAFARI_BASE_URL=https://apisandbox.safaricom.et
SAFARI_CALLBACK_URL=https://your-ngrok-url.ngrok-free.app/v1/payment/webhook
```

# How to run

# 1. Clone the repository

git clone https://github.com/Kirubel-Enyew27/safari-payment.git
cd safari-payment

# 2. Configure your environment

config/.env

# 3. Start PostgreSQL container

docker compose up -d

# 4. Run migrations (optional if MIGRATION_ACTIVE=true handles it)

go run cmd/main.go migrate

# 5. Start the service

go run cmd/main.go

# Testing the webhook

Use Ngrok to expose your local server to the internet:

ngrok http 8080

Replace the SAFARI_CALLBACK_URL in your .env file with the generated Ngrok URL.
