# Task Management Application

A full-stack task management application with a Go REST API backend and Next.js frontend.

## Tech Stack

- **Frontend**: Next.js 15, Tailwind CSS v4, TypeScript
- **Backend**: Go 1.21, gorilla/mux
- **Database**: PostgreSQL 16
- **Authentication**: JWT (HS256) with bcrypt password hashing
- **Infrastructure**: Docker Compose, Makefile

## Quick Start

### Prerequisites

- Docker & Docker Compose (everything else runs inside containers)

### One-Command Setup

```bash
make up
```

This starts all three services (Postgres, backend, frontend), runs migrations automatically, and makes them available at:

| Service  | URL                       |
|----------|---------------------------|
| Frontend | http://localhost:3000      |
| Backend  | http://localhost:8080      |
| Database | localhost:5432             |

### Seed Sample Data

```bash
make seed
```

This creates a test account and 3 sample tasks:

| Field    | Value              |
|----------|--------------------|
| Email    | test@example.com   |
| Password | password123        |

## Available Commands

```bash
make up             # Start all services
make down           # Stop all services
make logs           # Stream logs from all services
make test           # Run all tests (backend + frontend, via Docker)
make test-backend   # Run Go tests only
make test-frontend  # Run Jest tests only
make lint           # Lint all code
make migrate        # Run database migrations
make seed           # Seed sample data
make reset          # Stop everything and wipe volumes
```

## API Reference

### Authentication

| Method | Endpoint          | Auth | Description          |
|--------|-------------------|------|----------------------|
| POST   | /api/auth/signup  | —    | Register new user    |
| POST   | /api/auth/login   | —    | Login, returns JWT   |
| GET    | /api/auth/me      | JWT  | Get current user     |

### Tasks

All task endpoints require `Authorization: Bearer <token>`.

| Method | Endpoint        | Description                              |
|--------|-----------------|------------------------------------------|
| POST   | /api/tasks      | Create a task                            |
| GET    | /api/tasks      | List tasks (filters, search, pagination) |
| GET    | /api/tasks/:id  | Get a single task                        |
| PATCH  | /api/tasks/:id  | Partially update a task                  |
| DELETE | /api/tasks/:id  | Delete a task                            |

#### GET /api/tasks query parameters

| Parameter  | Type   | Description                                         |
|------------|--------|-----------------------------------------------------|
| status     | string | Filter: `pending`, `in_progress`, `completed`       |
| priority   | string | Filter: `low`, `medium`, `high`                     |
| search     | string | Full-text search on title and description (ILIKE)   |
| sort_by    | string | `created_at`, `updated_at`, `due_date`, `priority`, `status`, `title` |
| sort_order | string | `ASC` or `DESC` (default: `DESC`)                   |
| page       | int    | Page number (default: 1)                            |
| page_size  | int    | Results per page, max 100 (default: 10)             |

Response includes `total`, `page`, `page_size`, and `total_pages` for pagination.

#### Health check

```
GET /api/health → { "status": "healthy" }
```

## Environment Variables

All config is read from environment — nothing is hardcoded. See `.env.example` for the full list.

| Variable            | Required | Default               | Description                          |
|---------------------|----------|-----------------------|--------------------------------------|
| DATABASE_URL        | Yes      | —                     | PostgreSQL connection string         |
| BACKEND_PORT        | No       | 8080                  | Backend listen port                  |
| JWT_SECRET          | No       | dev-secret-key        | JWT signing key (change in prod)     |
| JWT_EXPIRY          | No       | 24h                   | Token lifetime (e.g. `12h`, `7d`)    |
| CORS_ORIGINS        | No       | http://localhost:3000 | Allowed CORS origin                  |
| DB_MAX_OPEN_CONNS   | No       | 25                    | DB connection pool size              |
| DB_MAX_IDLE_CONNS   | No       | 5                     | DB idle connection pool size         |
| LOG_LEVEL           | No       | info                  | `debug`, `info`, `warn`, `error`     |
| NEXT_PUBLIC_API_URL | No       | http://localhost:8080 | API base URL (browser-visible)       |

> `DATABASE_URL` has no default — the backend exits immediately if it is missing.

## Database Schema

### users
```sql
id         uuid        PRIMARY KEY
email      varchar     UNIQUE NOT NULL
password   varchar     NOT NULL  -- bcrypt hash
created_at timestamptz NOT NULL
updated_at timestamptz NOT NULL
```

### tasks
```sql
id          uuid        PRIMARY KEY
user_id     uuid        NOT NULL REFERENCES users(id)
title       varchar     NOT NULL
description text
status      varchar     NOT NULL  -- pending | in_progress | completed
priority    varchar     NOT NULL  -- low | medium | high
due_date    timestamptz
created_at  timestamptz NOT NULL
updated_at  timestamptz NOT NULL
```

Indexes on `user_id`, `status`, `created_at`.

## Testing

Tests run entirely inside Docker — no local Go or Node.js installation required.

```bash
# All tests
make test

# Backend only (Go)
make test-backend

# Frontend only (Jest)
make test-frontend
```

Current coverage:

| Suite    | Tests | What's covered                                  |
|----------|-------|-------------------------------------------------|
| Backend  | 3     | Password hashing, JWT generation & verification |
| Frontend | 4     | AuthCard rendering, form fields, mode toggle    |

## Project Structure

```
.
├── backend/
│   ├── cmd/main.go              # Entry point (migrate / seed / serve)
│   ├── internal/
│   │   ├── api/                 # HTTP handlers, middleware, helpers
│   │   ├── auth/                # JWT + bcrypt
│   │   ├── db/                  # SQL repositories
│   │   ├── log/                 # Structured logger
│   │   └── models/              # Shared data types
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── app/                     # Next.js App Router pages
│   ├── components/              # React components
│   ├── lib/                     # Axios client, Zustand auth store
│   ├── __tests__/               # Jest tests
│   └── package.json
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Features

- Task CRUD with full input validation on all write endpoints
- JWT authentication — signup, login, protected routes
- Per-user data isolation enforced at the database query level
- Auth state persisted in `localStorage` — no flash or logout on page refresh
- Task filtering by status and priority
- Full-text search across title and description
- Sorting by any column, ascending or descending
- Pagination with total page count returned on every list response
- Optimistic UI updates — no full reload after create, update, or delete
- Dark mode toggle (persisted across sessions)
- Responsive layout for mobile and desktop

## Security Notes

- CORS is handled at the HTTP handler level (not gorilla/mux middleware), ensuring preflight `OPTIONS` requests are always answered correctly
- All dynamic SQL parameters (status, priority, sort column, sort order) are validated against whitelists before use — no string interpolation of user input
- Request bodies are capped at 1 MB to prevent memory exhaustion
- JWT algorithm is explicitly verified as HMAC — algorithm confusion attacks rejected
- Context keys use an unexported struct type — cannot be forged or collided with by external packages
- `DATABASE_URL` required at startup — no silent fallback to a default database

## Troubleshooting

**Backend fails to start (database connection error)**
```bash
docker-compose ps postgres      # Check Postgres is healthy
docker-compose logs postgres    # Inspect Postgres logs
```

**Frontend can't reach the backend**
```bash
curl http://localhost:8080/api/health    # Verify backend is up
# In docker-compose.yml, NEXT_PUBLIC_API_URL must be http://localhost:8080
# (not http://backend:8080 — browsers cannot resolve Docker service names)
```

**401 on all task requests after login**
- Token may have expired (default 24h) — log out and log back in
- Verify `JWT_SECRET` is consistent between runs (not regenerated on restart)

**Full reset**
```bash
make reset    # Stops containers and wipes all volumes including the database
make up       # Fresh start
make seed     # Recreate sample data
```
