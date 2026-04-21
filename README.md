#  Internal Analytics Event Processor

A lightweight backend service built in Go for ingesting, processing, and storing structured analytics events. Designed with reliability, simplicity, and extensibility in mind, this project mirrors real-world internal data platform tooling.

---

##  Overview

Modern applications generate continuous streams of events such as user actions, system logs, and product interactions. This service provides a clean and reliable way to:

- Accept structured event data via HTTP
- Validate and process incoming requests
- Persist events in PostgreSQL
- Expose basic analytics endpoints
- Apply production-grade patterns like timeouts and graceful shutdown

---

##  Architecture

```
Client Application
        │
        ▼
   HTTP API (Go)
        │
        ▼
 Validation Layer
        │
        ▼
 PostgreSQL (JSONB storage)
        │
        ▼
 Analytics Endpoints
```

---

##  Tech Stack

- **Language:** Go  
- **Database:** PostgreSQL  
- **Driver:** pgx (connection pooling)  
- **Protocol:** REST (HTTP/JSON)  

---

##  Features

- Event ingestion via `/ingest`
- Flexible schema using JSONB
- PostgreSQL persistence with parameterized queries
- Connection pooling for efficient DB access
- Request validation and error handling
- Graceful shutdown with signal handling
- Configurable via environment variables
- Health check endpoint (`/health`)
- Analytics endpoint (`/analytics/count`)

---

## 📡 API Endpoints and examples

### 1. Ingest Event

**POST** `/ingest`

```json
{
  "event_type": "user_login",
  "source": "web_app",
  "payload": {
    "user_id": 42,
    "region": "us"
  }
}
```

**Response:**
```
202 Accepted
```

---

### 2. Health Check

**GET** `/health`

```json
{
  "status": "ok"
}
```

---

### 3. Event Count

**GET** `/analytics/count`

```json
{
  "total_events": 10
}
```

---

##  Database Schema

```sql
CREATE TABLE analytics_events (
    id SERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    source TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

##  Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Atharva1701/Internal-Analytics-Event-Processor.git
```

---

### 2. Set environment variables

#### Windows (PowerShell)

```powershell
setx DATABASE_URL "postgres://<user>:<your-password>@localhost:5432/analytics?sslmode=disable"
setx SERVER_ADDR ":8080"
```

Restart your terminal after setting variables.

---

### 3. Install dependencies

```bash
go mod tidy
```

---

### 4. Run the service

```bash
go run .
```

Server runs on:

```
http://localhost:8080
```

---

##  Testing the API

### Using curl (PowerShell)

```powershell
curl http://localhost:8080/ingest `
  -Method POST `
  -ContentType "application/json" `
  -Body '{
    "event_type": "user_login",
    "source": "web_app",
    "payload": {
      "user_id": 42
    }
  }'
```

---

##  Example Query

```sql
SELECT * FROM analytics_events;
```

---

##  Design Considerations

- Parameterized SQL queries to prevent injection  
- Context-based timeouts for DB operations  
- Connection pooling for scalability  
- Concurrent request handling using Go  
- JSONB for flexible and evolving schemas  

---

##  Future Improvements

- Authentication (JWT)
- Rate limiting
- Schema validation
- Pagination for analytics
- Indexing for performance
- Message queue integration (Kafka)
- Monitoring and metrics
- Docker-based deployment

---

##  What This Project Demonstrates

- Backend system design fundamentals  
- REST API development in Go  
- Database modeling and querying  
- Production-ready engineering practices  
- Handling concurrency, reliability, and scalability concerns  

---
