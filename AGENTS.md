# AGENTS.md

## Stack-Specific Notes
- HTTP: Chi framework, handlers in `internal/controller/restapi/`, server wrapper in `pkg/httpserver/`
- Storage: Postgres-only, implementations in `internal/repo/persistent/`, client wrapper in `pkg/postgres/`
- Logger: Standard library `log` (no third-party loggers), wrapper in `pkg/logger/`
- Swagger: `swag` CLI, entrypoint is the Chi router file in `internal/controller/restapi/`
- Clean arch layers: `controller → usecase → repo → entity`; interfaces in `internal/{repo,usecase}/contracts.go`
- No auth yet: No middleware or auth code in `restapi/`
- Migrations: `migrations/` uses `golang-migrate`

## Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) format: `<type>(<scope>): <subject>`

Types:
- `feat` — new feature
- `fix` — bug fix
- `refactor` — code change that is not a fix or feature
- `test` — adding or updating tests
- `chore` — dependency updates, tooling, config
- `docs` — documentation only
- `perf` — performance improvement

Rules:
- Subject in lowercase, no period at end
- Scope is optional but preferred; use the package name (e.g. `entity`, `repo`, `usecase`)
- Breaking changes: append `!` after type — `feat!: change task ID to uuid`

Examples:
```
feat(entity): add status transition validation
fix(repo): handle not found error from postgres
chore: update dependencies
```

## Common Commands
- Swagger gen: `swag init --parseDependency -g internal/controller/restapi/router.go`
- Tests: `go test -v -race ./internal/... ./pkg/...`
- Run app: `go run ./cmd/app` (requires running Postgres, applied migrations)
- Migrations: `migrate -path migrations -database "$PG_URL?sslmode=disable" up`

