# AGENTS.md

## Project Structure

```
cmd/
  app/            — entrypoint; wires all dependencies
  migrate/        — standalone migration runner
config/           — typed config loaded from environment
internal/
  app/            — DI root: connects repo → usecase → controller
  entity/         — domain types and business rules (no internal imports)
  usecase/
    contracts.go  — TaskUsecase interface + filter DTOs
    task/         — business logic implementation
  repo/
    contracts.go  — TaskRepository interface + repo-level sentinel errors
    persistent/   — pgx-backed implementation
  controller/
    restapi/      — Chi router setup + Swagger annotations
    restapi/v1/   — HTTP handlers, request/response DTOs
migrations/       — SQL files for golang-migrate (up/down pairs)
pkg/
  postgres/       — pgxpool wrapper with retry logic
  httpserver/     — http.Server wrapper with graceful shutdown
  logger/         — stdlib log wrapper
```

## Clean Architecture Rules

Dependency direction flows **inward only**:

```
controller → usecase → repo → entity
```

| Layer | Responsibility | Hard constraints |
|---|---|---|
| `entity` | Domain types, transitions, validation | Zero imports from other `internal/` packages |
| `repo` | Data access interface + sentinel errors | `contracts.go` owns `TaskRepository`; concrete impl in `persistent/` |
| `usecase` | Orchestrates repo calls, enforces business rules | Depends on `repo.TaskRepository` interface only, never `persistent.*` |
| `controller` | HTTP decode/encode, routing | Depends on `usecase.TaskUsecase` interface only, never imports `repo` |

**Never** import a concrete implementation across layers.  
**Never** let `entity` depend on persistence or transport types.  
Error translation flows inward → outward: `pgx.ErrNoRows` → `repo.ErrNotFound` → `entity.Err*` → HTTP status.

## Stack & context7 IDs

| Concern | Module | context7 ID |
|---|---|---|
| HTTP router | `github.com/go-chi/chi/v5` | `/go-chi/docs` |
| PostgreSQL driver | `github.com/jackc/pgx/v5` | `/jackc/pgx` |
| Migrations | `github.com/golang-migrate/migrate/v4` | resolve with context7 if needed |
| UUID | `github.com/google/uuid` | resolve with context7 if needed |

### Fetching docs with context7

Always fetch current docs before using a library API:

```
# Chi: middleware, routing, subrouters
context7: /go-chi/docs  →  query: "middleware subrouter URL params"

# pgx: pool, transactions, error handling
context7: /jackc/pgx    →  query: "pgxpool acquire transaction pgx.ErrNoRows"
```

## Common Commands

```bash
# Run (requires running Postgres + applied migrations)
go run ./cmd/app

# Tests with race detector
go test -v -race ./internal/... ./pkg/...

# Regenerate Swagger docs
swag init --parseDependency -g internal/controller/restapi/router.go

# Apply all migrations
migrate -path migrations -database "$PG_URL?sslmode=disable" up

# Roll back one migration
migrate -path migrations -database "$PG_URL?sslmode=disable" down 1
```

## Go Conventions

### Errors
- Wrap with context at each layer boundary: `fmt.Errorf("usecase.GetTask: %w", err)`.
- Use `errors.Is` / `errors.As` — never string-match error messages.
- Sentinel errors belong to the layer that produces them (`repo.ErrNotFound`, `entity.ErrInvalidTransition`).
- Usecases translate repo errors; controllers translate usecase errors to HTTP status codes.

### Interfaces
- Defined in the package that **consumes** them (or in `contracts.go` of the package being consumed).
- Keep interfaces small — include only what the consumer actually calls.
- Accept interfaces, return concrete types.

### Context
- `context.Context` is always the first argument of any function that does I/O.
- Never store `context.Context` in a struct field.

### Naming
- Acronyms uppercase: `ID`, `URL`, `HTTP`, `API` — not `Id`, `Url`.
- Unexported test helpers: `testHelperName(t *testing.T, ...)`.
- Avoid redundant package prefixes: `entity.Task`, not `entity.EntityTask`.

### Testing
- Table-driven tests with `t.Run`.
- Always run with `-race` flag.
- Test file lives next to the file under test (`task_test.go` beside `task.go`).

## Adding a Feature (Checklist)

1. **entity** — add/extend domain type; add `Validate` and transition logic if needed.
2. **repo interface** — extend `TaskRepository` in `internal/repo/contracts.go`; write migration if schema changes.
3. **repo implementation** — implement in `internal/repo/persistent/task.go`; translate `pgx.ErrNoRows` → `repo.ErrNotFound`.
4. **usecase interface** — extend `TaskUsecase` in `internal/usecase/contracts.go`.
5. **usecase implementation** — implement in `internal/usecase/task/task.go`; translate repo errors to entity/domain errors.
6. **handler** — add/update in `internal/controller/restapi/v1/task_handler.go`; use DTOs from `request/` and `response/`.
7. **swagger** — regenerate docs (see command above).
8. **tests** — write table-driven tests; verify with `go test -v -race ./...`.

## Commit Conventions

Format: `<type>(<scope>): <subject>`

| Type | When |
|---|---|
| `feat` | new capability |
| `fix` | bug fix |
| `refactor` | restructuring without behavior change |
| `test` | tests only |
| `chore` | deps, tooling, config |
| `docs` | documentation only |
| `perf` | performance improvement |

Rules:
- Subject lowercase, no trailing period.
- Scope = package name: `entity`, `repo`, `usecase`, `restapi`.
- Breaking change: append `!` — `feat!: change task ID to uuid`.

```
feat(entity): add status transition validation
fix(repo): handle not found error from postgres
chore: update dependencies
```

## Behavioral Guidelines

### Think before coding
- State assumptions explicitly. If uncertain, ask before writing code.
- If multiple interpretations exist, present the options — don't pick silently.

### Simplicity first
- Minimum code that solves the problem. No speculative features.
- No abstractions for single-use code. No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

### Surgical changes
- Don't "improve" adjacent code or formatting — match existing style.
- Remove only imports/vars/functions that **your** changes made unused.
- If you notice unrelated dead code, mention it; don't delete it.

### Goal-driven execution
- Turn every task into a verifiable goal: write the test first, then make it pass.
- For multi-step tasks, state a brief numbered plan before starting.
