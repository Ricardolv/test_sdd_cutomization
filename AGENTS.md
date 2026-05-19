# AGENTS.md

## Stack

- **Backend**: Go + chi router, PostgreSQL via `database/sql` (no ORM)
- **Frontend**: Next.js (React, TypeScript) + shadcn/ui
- **Dev server**: `air` for Go hot reload, `npm run dev` for Next.js
- **Default port**: 9090 (backend, from `config/config.go`)

## Project Structure

```
apps/
  backend/
    cmd/server/          # entry point
    config/              # env config
    scripts/             # migrations (migrate tool), HTTP test files
    internal/
      database/          # DB connection helper
      modules/           # business logic (models, services, repositories, handlers)
  frontend/
    src/app/             # Next.js pages
    src/shared/          # reusable code (components/ui = shadcn)
  docker-compose.yml     # PostgreSQL service
```

## Commands

| Task | Command |
|------|---------|
| Run backend | `go run ./cmd/server` (from `apps/backend/`) |
| Hot reload | `air` (from `apps/backend/`) |
| Run tests | `go test ./...` (from `apps/backend/`) |
| Run single package tests | `go test ./internal/modules/customer/services/...` |
| Lint | `golangci-lint run` (from `apps/backend/`) |
| Format | `gofmt -w .` or `go fmt ./...` |
| Run frontend | `npm run dev` (from `apps/frontend/`) |
| Typecheck | `npx tsc --noEmit` (from `apps/frontend/`) |
| Start DB | `docker compose up -d` (from `apps/`) |
| Migrations | `migrate -path scripts/migrations -database "$DATABASE_URL" up` |

## Setup Prerequisites

- PostgreSQL must be running before backend starts (backend panics on DB connection failure)
- Run `docker compose up -d` from `apps/` to start the database
- Apply migrations before first run: `migrate -path scripts/migrations -database "postgres://postgres:postgres@localhost:5432/testapp?sslmode=disable" up`
- Repository integration tests require a running PostgreSQL database; they skip gracefully if unavailable

## Spec-Driven Workflow

All changes start as a spec in `.spec/changes/`. Read before implementing:

- `.spec/memory/contexto-tecnico.md` — global technical context
- `.spec/memory/estrutura.md` — project structure and conventions
- `.spec/shared/como-executar.md` — execution rules and evidence format
- `.spec/shared/regras-de-nomenclatura.md` — naming conventions (Go: `*_handler.go`, Next.js: `*.page.tsx`)
- `.spec/templates/modelo-base.md` — base spec template
- `.spec/templates/modelo-crud.md` — CRUD spec template

**Rule**: spec describes the change **before** implementation. One spec per delivery.

## Architecture Rules

- Frontend never knows about the database
- Backend exposes use cases via REST API (JSON)
- Business rules don't depend on the web interface
- Each module: `models/` → `services/` → `repositories/` → `handlers/`
- Services contain business logic and orchestrate repositories
- Handlers use chi for routing, call services
- Centralized error handling via middleware
- No premature abstractions — keep it simple
- UUIDs generated in Go via `github.com/google/uuid` (not DB-side)

## Conventions

- `internal/` = private app code (not importable externally)
- `shared/` = reusable, project-agnostic code only
- Auth, session, routes, tenant logic live in `app/`, not `shared/`
- API responses: map domain entities to simple JSON objects before returning
- Validations happen on both client and server
- Inactivation = soft-delete (`active` field set to `false`, not row removal)
- Queries filter `active = true` by default
