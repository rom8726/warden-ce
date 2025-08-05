# Warden Community Edition: Self-Hosted Sentry-Compatible Error Monitoring Platform

<div align="center">
  <img src="docs/logo_full.png" alt="Warden Logo" height="150"/>
</div>

---

## Brief Overview

**Warden** is a self-hosted error monitoring platform fully compatible with the official Sentry SDKs. It enables organizations to collect, store, and analyze application errors on their own infrastructure, with a modern web UI and advanced incident management features. Warden is designed as a drop-in Sentry alternative: you can use existing Sentry SDKs without any changes to client code.

---

## Description

Warden provides a robust, extensible, and privacy-friendly solution for error tracking and alerting. It supports the Sentry ingestion protocol, allowing seamless integration with popular SDKs (`sentry-go`, `@sentry/node`, `sentry-python`, etc.). The platform features a clean architecture, API-first development, and a powerful React-based UI for error analysis, team collaboration, and project management.

---

## Table of Contents

1. [About the Project](#about-the-project)
2. [Key Features](#key-features)
3. [Architecture & Components](#architecture--components)
4. [Development Practices & Highlights](#development-practices--highlights)
5. [Getting Started](#getting-started)
6. [Installer](#installer)
7. [Testing](#testing)

---

## About the Project

Warden is designed for teams and companies that need full control over their error monitoring stack, data privacy, and the ability to customize or extend the platform. Community Edition version is suitable for small teams and individual pet projects. There is also an Enterprise version of the project exist. To learn more, visit the [warden-project.tech](https://warden-project.tech) website.

---

## Key Features

- **Sentry SDK Compatibility:** Accepts events via `/api/:project_id/store/` and `/api/:project_id/envelope/` endpoints, using standard Sentry DSN and authentication headers.
- **Modern Web UI:** Powerful React-based interface for error analysis, filtering, search, and team workflows.
- **Project & Team Management:** RBAC, 2FA, user and team management, project settings.
- **Event Grouping & Fingerprinting:** Advanced grouping of errors and exceptions for efficient triage.
- **Notifications:** Integrations with Email, Slack, Telegram, Mattermost, and Webhooks.
- **Metrics & Monitoring:** Prometheus metrics, health checks, and rate limiting.
- **API-First:** OpenAPI specification (`specs/server.yml`) is the single source of truth for the API. Code and DTOs are generated from the spec.
- **Scalable Storage:**
  - **PostgreSQL:** Metadata, users, projects, teams, notifications.
  - **ClickHouse:** High-performance storage for events and stack traces (via Kafka ingestion).
  - **Kafka:** Buffering and streaming of incoming events.
  - **Redis:** Caching and rate limiting.
- **Docker-Ready:** Full dev and production setup via Docker Compose and multi-stage Dockerfile.

## Deployment Instructions

After reviewing the key features, you can deploy Warden using Docker Compose for production.  
**Follow these steps:**

1. **Copy the production Docker Compose file:**

   Use the `docker-compose.yml` from the `prod` directory.

2. **Prepare configuration files:**

   Create the following files in `/opt/warden` (must be owned by `root`):

   - `/opt/warden/platform.env` — platform settings
   - `/opt/warden/config.env` — application configuration

3. **Fill in `platform.env`:**

   ```
   DOMAIN=warden.your-company.tech         # Domain for the platform
   DB_PASSWORD=your_postgres_password      # PostgreSQL password
   CLICKHOUSE_DB_PASSWORD=your_ch_password # ClickHouse password
   SSL_CERT=warden.crt                     # SSL certificate filename (should be in /opt/warden/nginx/ssl)
   SSL_CERT_KEY=warden.key                 # SSL key filename (should be in /opt/warden/nginx/ssl)
   PLATFORM_VERSION=latest                 # Platform version (e.g., v1.2.3 or latest)
   ```

4. **Fill in `config.env`:**

   ```
   WARDEN_FRONTEND_URL=https://warden.your-company.tech   # URL for emails and redirects
   WARDEN_SECRET_KEY=your_very_secret_key                # 16 or 32 chars
   WARDEN_JWT_SECRET_KEY=your_jwt_secret_key             # 16 or 32 chars
   WARDEN_POSTGRES_PASSWORD=your_postgres_password       # Same as DB_PASSWORD in platform.env
   WARDEN_CLICKHOUSE_PASSWORD=your_ch_password           # Same as CLICKHOUSE_DB_PASSWORD in platform.env
   WARDEN_MAILER_ADDR=smtp.yandex.ru:465                 # SMTP server address
   WARDEN_MAILER_USER=mailer_user                        # SMTP user
   WARDEN_MAILER_PASSWORD=mailer_password                # SMTP password
   WARDEN_MAILER_FROM=warden@your-company.tech           # Sender email
   WARDEN_MAILER_USE_TLS=true                            # Use TLS (true/false)
   WARDEN_MAILER_CERT_FILE=                              # (optional) Path to TLS cert
   WARDEN_MAILER_KEY_FILE=                               # (optional) Path to TLS key
   WARDEN_MAILER_ALLOW_INSECURE=false                    # Allow insecure certs (for self-signed)
   WARDEN_ADMIN_EMAIL=admin@your-company.tech            # Platform admin email
   WARDEN_ADMIN_TMP_PASSWORD=your_admin_temp_password    # Initial admin password
   ```

5. **Place SSL certificates:**

   Place your SSL certificate and key files in `/opt/warden/nginx/ssl` with the names specified in `platform.env`.

6. **Start the platform:**

   ```bash
   cd /opt/warden
   docker compose --env-file platform.env up -d
   ```

---

## Architecture & Components

Warden follows a strict **Layered Clean Architecture**:

- **Domain Layer (`internal/domain`):** Core business entities (e.g., User, Project, Issue) with no dependencies on frameworks or storage.
- **Contract Layer (`internal/contract`):** Interfaces (contracts) for repositories and services, enabling dependency inversion.
- **Use Case Layer (`internal/usecases`):** Business logic orchestrating domain entities and repositories via interfaces.
- **Repository Layer (`internal/repository`):** Data access implementations, mapping between storage models and domain entities.
- **API Layer (`internal/api/rest`):** HTTP handlers, request validation, DTO conversion, and OpenAPI-based codegen.
- **Services Layer (`internal/services`):** Integrations with external systems (email, messengers, etc.).

**Dependency Rule:** Dependencies always point inward; outer layers depend on interfaces of inner layers, never the other way around.

**Transaction Management:** Multi-repository operations are wrapped in transactions using a transaction manager (`pkg/db/TxManager`).

**Dependency Injection:** All components are initialized via a DI container ([rom8726/di](https://github.com/rom8726/di)).

---

## Development Practices & Highlights

- **Strict Layer Boundaries:**
  - API layer only calls use cases
  - Use cases depend only on interfaces from `internal/contract`
  - Repositories implement only interfaces, never expose storage details to business logic
- **Dependency Inversion:** All dependencies are injected via interfaces, not concrete types.
- **Domain-Driven:** Domain entities are pure Go structs, free from storage or framework tags.
- **API-First:** All API changes start with the OpenAPI spec (`specs/server.yml`). DTOs and handlers are generated and mapped to domain models in `internal/dto`.
- **Testing:**
  - **Unit tests:** Use mocks (generated via `mockery`, stored in `test_mocks/`)
  - **Functional tests:** Use [testy](https://github.com/rom8726/testy) — declarative YAML scenarios in `tests/cases/`, with fixtures in `tests/fixtures/`. Testy spins up the app, loads fixtures, runs HTTP requests, and checks DB state.
  - **Integration tests:** Run with `go test -tags=integration ./tests/...`
- **Linting & Formatting:**
  - `golangci-lint` with strict rules (see `.golangci.yml`)
  - Max line length: 120 characters
  - Auto-formatting via `gofmt` and `gofumpt`
- **Dev Environment:**
  - `make dev-up` — start full dev stack (Postgres, ClickHouse, Kafka, Redis, etc.)
  - `make dev-down`, `make dev-clean`, `make dev-logs` for management
  - Environment variables in `dev/config.env.example`
- **Migration Management:**
  - PostgreSQL and ClickHouse migrations in `migrations/`
  - Run automatically on startup

---

## Getting Started

### Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rom8726/warden-ce.git
   cd warden
   ```
2. **Start the development environment:**
   ```bash
   make dev-up
   ```
3. **Generate server code from OpenAPI spec:**
   ```bash
   make generate-backend
   ```
4. **Access the UI:**
   - Open [http://localhost:3000](http://localhost:3000) in your browser
5. **View logs:**
   ```bash
   make dev-logs
   ```

---

## Testing

Warden uses its own [testy](https://github.com/rom8726/testy) framework for functional API testing.

### Functional Testing

**testy** is a declarative framework for testing HTTP APIs that allows:

- Writing tests in YAML format with scenario descriptions
- Running tests through `http.Handler` with debugging capabilities (breakpoints)
- Working with PostgreSQL fixtures via [pgfixtures](https://github.com/rom8726/pgfixtures)
- Checking JSON responses, SQL queries, and external service mocks
- Using templates, variables, and step chains

### Running Tests

```bash
# Unit tests
make test

# Functional tests
make test.integration
```

---

## DSN Structure and Protocol

### DSN Format

DSN (Data Source Name) is used by the SDK to configure the connection to the server. It includes:

```
https://<public_key>@<host>/<project_id>
```

Example:

```
https://abc123@warden.local/42
```

### Authorization

Sentry SDK uses one of the following headers:

#### `X-Sentry-Auth` (preferred)

```
X-Sentry-Auth: Sentry sentry_key=abc123, sentry_version=7, sentry_client=sentry.go/0.13.0
```

#### Or `Authorization: Sentry ...`

```
Authorization: Sentry sentry_key=abc123, sentry_version=7
```

---

## API: Event Reception

### `POST /api/:project_id/store/`

#### Purpose

Reception of error events and exceptions sent using official Sentry SDKs.

#### Parameters

| Parameter   | Type       | Required | Description                                        |
| ----------- | ---------- | -------- | -------------------------------------------------- |
| `project_id`| path param | ✅        | Project identifier, matches `project_id` from DSN  |

#### Headers

* `X-Sentry-Auth` or `Authorization: Sentry ...`

#### Request Body

JSON format corresponding to [Sentry Event Payload](https://develop.sentry.dev/sdk/event-payloads/).

Example:

```json
{
  "event_id": "8e4f5d83f65b4173b0e4036a64042fda",
  "timestamp": "2025-06-04T17:00:00Z",
  "level": "error",
  "platform": "go",
  "message": "Unhandled panic: index out of range",
  "exception": {
    "values": [{
      "type": "panic",
      "value": "index out of range",
      "stacktrace": {
        "frames": [...]
      }
    }]
  },
  "user": {
    "id": "123"
  },
  "tags": {
    "env": "prod"
  }
}
```

#### Responses

| Code               | Meaning                               |
| ------------------ | ------------------------------------- |
| `200 OK`           | Event accepted                        |
| `400 Bad Request`  | Invalid event format                  |
| `401 Unauthorized` | Invalid or missing key                |
| `404 Not Found`    | Project not found                     |

---

## Project Architecture

The project is built on the principles of **Clean Architecture** with a clear separation into layers:

### Architecture Layers

1. **Domain Layer** (`internal/domain`)
   - Contains the core business models: `Event`, `Exception`, `Project`
   - Does not depend on other layers or external frameworks

2. **Business Logic Layer** (`internal/usecases`)
   - Implements the application's business logic
   - Interacts with repositories through interfaces
   - Divided by functional areas: `events`, `exceptions`, `projects`

3. **Repository Layer** (`internal/repository`)
   - Responsible for data access and storage
   - Implements interfaces defined in the business logic layer
   - Divided by entity types: `events`, `exceptions`, `projects`

4. **API Layer** (`internal/api/rest`)
   - Implements HTTP API for client interaction
   - Uses generated interfaces from OpenAPI specification
   - Transforms HTTP requests into business logic calls

### System Components

| Component           | Purpose |
|---------------------|---------|
| **PostgreSQL**      | Storage of project information and known fingerprints |
| **ClickHouse**      | Storage of events and exceptions with TTL and partitioning |
| **Kafka**           | Message broker between API and ClickHouse |
| **Redis**           | Cache for duplicate checking and rate limiting |
| **API Server**      | Processing incoming events, generating fingerprints, sending to Kafka |

### Data Flow

```
Sentry SDK → [HTTP POST] → API Server → [Converts + generates fingerprint] → [Writes to Kafka] → ClickHouse Kafka Engine → Materialized View → ClickHouse Tables
```

---

## Launch and Deployment

### Local Deployment

For local development and testing, use Docker Compose:

1. Clone the repository:
   ```bash
   git clone https://github.com/rom8726/warden-ce.git
   cd warden
   ```

2. Create environment configuration file:
   ```bash
   cp dev/compose.env.example dev/compose.env
   cp dev/config.env.example dev/config.env
   ```

3. Start services using Docker Compose:
   ```bash
   make dev-up
   ```

4. Check that all services are running:
   ```bash
   docker ps
   ```

5. To stop services:
   ```bash
   make dev-down
   ```

### Building and Running

To build the application:

```bash
make build
```

To run the server:

```bash
./bin/app server --env-file=./config.env
```

### Migrations

Migrations run automatically when the server starts. For manual migration execution:

```bash
./bin/app migrate --env-file=./config.env
```

## Monitoring and Metrics

The application provides Prometheus metrics at the `/metrics` endpoint of the technical server (default port 8081).

Main metrics:

- `warden_events_received_total` - number of events received
- `warden_events_processed_total` - number of events processed
- `warden_exceptions_received_total` - number of exceptions received
- `warden_exceptions_processed_total` - number of exceptions processed
- `warden_validation_errors_total` - number of validation errors
- `warden_processing_time_seconds` - event and exception processing time
- `warden_kafka_messages_produced_total` - number of messages sent to Kafka
- `warden_kafka_messages_consumed_total` - number of messages received from Kafka

## Development

### Project Structure

```
warden/
├── cmd/                  # Application entry points
├── dev/                  # Files for local development
├── internal/             # Internal application code
│   ├── backend/          # Backend code
│   │   ├── api/          # API layer
│   │   ├── contract/     # Interfaces (contracts)
│   │   ├── dto/          # Data Transfer Objects
│   │   └── usecases/     # Business logic
│   ├── context/          # Context utilities
│   ├── domain/           # Domain models
│   └── repository/       # Repositories for data access
├── migrations/           # Database migrations
│   ├── clickhouse/       # Migrations for ClickHouse
│   └── postgresql/       # Migrations for PostgreSQL
├── pkg/                  # Reusable packages
│   ├── db/               # Database utilities
│   ├── httpserver/       # HTTP server
│   ├── kafka/            # Kafka utilities
│   └── metrics/          # Prometheus metrics
├── specs/                # OpenAPI specifications
├── test_mocks/           # Generated mocks for testing
└── tests/                # Tests
```

### Code Generation

Server code is generated from OpenAPI specification using ogen:

```bash
make generate
```

### Linting

For code checking, golangci-lint is used:

```bash
make lint
```

---

## References and Resources

* [Sentry Event Payloads](https://develop.sentry.dev/sdk/event-payloads/)
* [Sentry Envelope Format](https://develop.sentry.dev/sdk/envelopes/)

---

## API Specification First

The project uses the **API Specification First** approach. This means that:

* Description of all endpoints is maintained in an **OpenAPI YAML document** (`specs/server.yml`)
* This document is the **single source of truth** for model generation, validation, auto-documentation, and client
* Server code is generated using [ogen](https://github.com/ogen-go/ogen)

Main endpoints:

1. `POST /api/{project_id}/store/` - for receiving events in JSON format
2. `POST /api/{project_id}/envelope/` - for receiving events in Envelope format

---

## Fingerprint Calculation

Fingerprint is used for grouping events and exceptions:

### Fingerprint For Event
- For events: hash of message + level + platform

### Fingerprint For Exception
- For exceptions: hash of exception_type, exception_value, stacktrace

SHA1 is used for fingerprint calculation.

---
