# URL Shortener

A scalable URL shortening service built with **Go** and **PostgreSQL**.

The project is designed as a backend-focused learning project following the [roadmap.sh URL Shortening Service](https://roadmap.sh/projects/url-shortening-service) requirements, with an architecture that can later be extended with Docker, CI/CD, Kubernetes, caching, and distributed-system concepts.

## Features

- Create shortened URLs
- Generate unique short codes
- Retrieve original URLs using short codes
- Update shortened URLs
- Delete shortened URLs
- Track URL access count
- Store creation and update timestamps
- PostgreSQL persistence
- Connection pooling with `pgx`
- Layered backend architecture
- REST-style API
- Docker-based PostgreSQL development environment

### Planned

- URL statistics endpoint
- HTTP redirect endpoint
- URL validation
- Short-code collision handling
- Unit and integration tests
- Dockerize the Go application
- CI/CD with GitHub Actions
- Kubernetes deployment
- Caching with Redis

---

## Architecture

The application follows a layered architecture:

```text
                    HTTP Request
                         │
                         ▼
                  ┌─────────────┐
                  │   Handler   │
                  │ HTTP / JSON │
                  └──────┬──────┘
                         │
                         ▼
                  ┌─────────────┐
                  │   Service   │
                  │ Business    │
                  │   Logic     │
                  └──────┬──────┘
                         │
                         ▼
                  ┌─────────────┐
                  │ Repository  │
                  │    SQL      │
                  └──────┬──────┘
                         │
                         ▼
                  ┌─────────────┐
                  │ PostgreSQL  │
                  └─────────────┘
```

### Responsibilities

**Handler**

Responsible for:

- HTTP requests and responses
- JSON encoding/decoding
- HTTP status codes
- Request validation

**Service**

Responsible for:

- Business logic
- Short-code generation
- Coordinating repository operations
- Application-level validation

**Repository**

Responsible for:

- PostgreSQL queries
- CRUD operations
- Database interaction

**Database**

Responsible for:

- PostgreSQL connection pool
- Database connectivity

**Models**

Contains the application's data structures.

---

## Project Structure

```text
url-shortener/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── handler/
│   │   └── url_handler.go
│   │
│   ├── models/
│   │   └── url.go
│   │
│   ├── repository/
│   │   └── url_repository.go
│   │
│   ├── service/
│   │   └── url_service.go
│   │
│   └── utils/
│       └── random.go
│
├── migrations/
│   └── 001_create_urls.sql
│
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Tech Stack

| Technology | Purpose |
|---|---|
| Go | Backend/API |
| `net/http` | HTTP server and routing |
| PostgreSQL | Persistent database |
| `pgx/v5` | PostgreSQL driver |
| `pgxpool` | Database connection pooling |
| Docker | PostgreSQL development environment |
| SQL | Database operations |

---

## Database Schema

The application uses a `urls` table:

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    access_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Fields

| Field | Description |
|---|---|
| `id` | Unique database identifier |
| `short_code` | Generated shortened URL code |
| `original_url` | Original destination URL |
| `access_count` | Number of times the URL was accessed |
| `created_at` | URL creation timestamp |
| `updated_at` | Last modification timestamp |

---

## API

### Create URL

```http
POST /api/v1/urls
```

Request:

```json
{
  "url": "https://www.example.com"
}
```

Response:

```json
{
  "short_code": "aB72xK",
  "original_url": "https://www.example.com"
}
```

---

### Get URL

```http
GET /api/v1/urls/{code}
```

Example:

```http
GET /api/v1/urls/aB72xK
```

Returns information about the shortened URL and increments its access count.

---

### Update URL

```http
PUT /api/v1/urls/{code}
```

Request:

```json
{
  "url": "https://www.example.org"
}
```

Response:

```json
{
  "message": "URL updated successfully"
}
```

---

### Delete URL

```http
DELETE /api/v1/urls/{code}
```

Response:

```http
204 No Content
```

---

### Statistics

```http
GET /api/v1/urls/{code}/stats
```

Returns information such as:

```json
{
  "short_code": "aB72xK",
  "original_url": "https://www.example.com",
  "access_count": 42,
  "created_at": "2026-10-06T18:00:00Z",
  "updated_at": "2026-10-06T18:10:00Z"
}
```

> Statistics endpoint is currently planned/in progress.

---

## Getting Started

### Prerequisites

Make sure you have:

- Go
- Docker
- Docker Compose
- Git

installed.

---

### 1. Clone the repository

```bash
git clone <your-repository-url>
cd url-shortener
```

---

### 2. Start PostgreSQL

Start the PostgreSQL container:

```bash
docker compose up -d
```

Check that the container is running:

```bash
docker ps
```

---

### 3. Configure the database

The development database uses:

```text
Host:     localhost
Port:     5432
Database: url_shortener
User:     postgres
Password: postgres
```

The connection string is:

```text
postgres://postgres:postgres@localhost:5432/url_shortener
```

> For production deployments, credentials should be provided through environment variables or a secrets manager rather than being hard-coded.

---

### 4. Create the database table

Connect to PostgreSQL:

```bash
docker exec -it url-shortener-db psql -U postgres -d url_shortener
```

Run the migration:

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    access_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

### 5. Install Go dependencies

```bash
go mod download
```

---

### 6. Run the application

```bash
go run ./cmd
```

The server will start on:

```text
http://localhost:8080
```

---

## Testing the API

You can use **Postman**, `curl`, or any HTTP client.

### Create

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.google.com"}'
```

### Get

```bash
curl http://localhost:8080/api/v1/urls/aB72xK
```

### Update

```bash
curl -X PUT http://localhost:8080/api/v1/urls/aB72xK \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.example.com"}'
```

### Delete

```bash
curl -X DELETE http://localhost:8080/api/v1/urls/aB72xK
```

---

## Database Connection Pooling

The application uses `pgxpool` instead of creating a new PostgreSQL connection for every request.

```text
Application
     │
     ▼
 pgxpool
 ┌───┼───┐
 ▼   ▼   ▼
DB  DB   DB
```

The pool is created once when the application starts and reused by incoming requests.

This improves performance and avoids repeatedly establishing database connections.

---

## Error Handling

The API uses appropriate HTTP status codes:

| Status | Meaning |
|---|---|
| `201 Created` | URL successfully created |
| `200 OK` | Successful request |
| `204 No Content` | URL successfully deleted |
| `400 Bad Request` | Invalid request |
| `404 Not Found` | Short URL does not exist |
| `500 Internal Server Error` | Server/database error |

---

## Reference

Project requirements:

[roadmap.sh — URL Shortening Service](https://roadmap.sh/projects/url-shortening-service)

---

## Status

**🚧 In Progress**

The core URL shortening functionality is implemented. The project is currently being expanded with statistics, testing, containerization, CI/CD, and Kubernetes deployment.