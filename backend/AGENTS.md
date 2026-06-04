## Project
Chi Game backend: Go API for chigame.ir. Phase 1: minimal surface (healthcheck, DB bootstrap). Multiplayer, rooms, matchmaking, and chat come later — do not add them until the frontend needs them.

## Stack
- **Connect RPC** (`connectrpc.com/connect`) over HTTP/2 (h2c) on port **8383**
- **PostgreSQL** via **pgx/v5** pool
- **goose** migrations in `db/migrations/` (embedded in `db/migrate.go`)
- **sqlc** — queries in `db/query/`, schema from `db/migrations/`, generated code in `database/`
- **Protobuf** in repo `proto/`; generate with `buf generate` from monorepo root (`buf.gen.yaml`)

## Module
`github.com/arian-nj/chigame/backend` — import paths must match.

## Layout
- `cmd/api/main.go` — entry: config, migrate, DB pool, start API
- `api/` — HTTP router, Connect handlers; `APIApplication` implements generated handler interfaces
- `db/` — goose migrations and sqlc query SQL only (no Go business logic)
- `database/` — sqlc-generated types and `Queries` (do not hand-edit)
- `gen/` — buf/protoc Connect + protobuf Go code (regenerate, don't hand-edit)
- `internals/config/` — env parsing (`DATABASE_URL`, `RELEASE_MODE`, `JWT_SECRET`)

## Conventions
- New RPC: add/extend `.proto` → `buf generate` → implement method on `APIApplication` → register in `api/router.go`
- New tables: goose migration with `-- +goose Up/Down` → add sql queries → `sqlc generate`
- Errors: return `connect.NewError(connect.Code*, err)` from RPC handlers
- CORS: keep frontend dev origins (`localhost:5173`, production domain) in `api/router.go`
- Secrets: never commit; use environment variables

## Commands
From `backend/`:
- `go run ./cmd/api/`
- `sqlc generate`
From repo root:
- `buf generate`
- `buf lint`

Write Modular code where possible