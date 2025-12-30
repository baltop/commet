# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Start PostgreSQL database
docker compose up -d

# Install dependencies
go mod tidy

# Build
go build -o bin/server ./cmd/server

# Run server (requires PostgreSQL to be running)
./bin/server

# Development with hot reload (requires air)
go install github.com/air-verse/air@latest
air
```

## Architecture Overview

This is a Go Gin web application using a layered architecture with HTMX + Alpine.js frontend.

### Backend Layers (internal/)

**Request Flow:** Handler → Service → Repository → Database

- **handlers/**: HTTP request handlers, render templates, call services
- **services/**: Business logic, authentication (JWT generation/validation)
- **repository/**: Data access layer using GORM
- **models/**: GORM models and DTOs
- **middleware/**: JWT auth middleware (`AuthMiddleware`), guest-only middleware (`GuestMiddleware`)
- **config/**: Viper-based configuration from `.env`
- **database/**: GORM connection, migrations, seed data

### Authentication

- JWT tokens stored in HTTP-Only cookies (cookie name: `auth_token`)
- `AuthMiddleware` validates tokens and sets claims in Gin context
- `GuestMiddleware` redirects authenticated users away from login/register

### Frontend (web/templates/)

Templates use Go html/template with custom functions:
- `safeJS`: Output JSON data unescaped for Chart.js
- `dict`: Create maps in templates
- `slice`: String slicing for user initials

**Template naming**: Templates are referenced by path relative to `web/templates/` (e.g., `auth/login.html`)

### HTMX Pattern

Dashboard charts load asynchronously via HTMX partials:
- Main page has containers with `hx-get="/dashboard/charts/line"` etc.
- Handlers return partial HTML snippets (not full pages)
- Check `HX-Request` header for HTMX-specific responses

## Environment Variables

Key variables in `.env`:
- `DB_PORT=5435` (PostgreSQL port, differs from default 5432)
- `JWT_SECRET` - Required for production
- `GIN_MODE` - Set to `release` for production
