# AGENTS.md

## Stack-Specific Notes
- HTTP: Chi framework, handlers in `internal/controller/restapi/`, server wrapper in `pkg/httpserver/`
- Storage: Postgres-only, implementations in `internal/repo/persistent/`, client wrapper in `pkg/postgres/`
- Logger: Standard library `log` (no third-party loggers), wrapper in `pkg/logger/`
- Swagger: `swag` CLI, entrypoint is the Chi router file in `internal/controller/restapi/`
- Clean arch layers: `controller → usecase → repo → entity`; interfaces in `internal/{repo,usecase}/contracts.go`
- No auth yet: No middleware or auth code in `restapi/`
- Migrations: `migrations/` uses `golang-migrate`

## Common Commands
- Swagger gen: `swag init --parseDependency -g internal/controller/restapi/router.go`
- Tests: `go test -v -race ./internal/... ./pkg/...`
- Run app: `go run ./cmd/app` (requires running Postgres, applied migrations)
- Migrations: `migrate -path migrations -database "$PG_URL?sslmode=disable" up`
