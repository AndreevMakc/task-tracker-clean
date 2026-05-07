# Task Tracker

A task management and tracking application operated via Telegram bot. Built with Go using Clean Architecture.

## Tech Stack

- **Language**: Go
- **HTTP**: Chi framework (webhooks, health checks)
- **Database**: PostgreSQL
- **Logger**: Standard library `log`
- **API Docs**: Swagger
- **Bot**: Telegram Bot API

## Architecture

Clean Architecture with four layers:

```
controller → usecase → repo → entity
```

### Directory Structure

```
cmd/app/              # Application entrypoint
internal/
  app/                # Bootstrap and shutdown logic
  controller/         # Adapters: Telegram bot handlers, REST API
    restapi/          # Chi HTTP routes (webhooks, Swagger)
  entity/             # Domain models and business errors
  usecase/            # Business logic
  repo/
    persistent/       # PostgreSQL implementations
pkg/
  httpserver/         # HTTP server wrapper
  logger/             # Logger wrapper
  postgres/           # Database client wrapper
migrations/           # SQL migrations (golang-migrate)
config/               # Configuration
docs/                 # Swagger documentation
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- [migrate CLI](https://github.com/golang-migrate/migrate) (for database migrations)
- Telegram Bot Token (from @BotFather)

### Quick Start (using Make)

1. Create `.env` file (copy from example):
```bash
cp .env.example .env
```

2. Edit `.env` with your settings (update `PG_URL` with your PostgreSQL credentials).

3. Run the application (auto-applies migrations):
```bash
make run
```

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make run` | Apply migrations and start the application |
| `make migrate` | Apply all pending database migrations |
| `make check-migrations` | Show current migration version in database |
| `make help` | Show all available commands |

### Manual Setup (without Make)

1. Create `.env` file:
```env
PG_URL=postgres://user:pass@localhost:5432/tasktracker?sslmode=disable
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

2. Apply migrations:
```bash
migrate -path migrations -database "$PG_URL" up
```

3. Run the app:
```bash
go run ./cmd/app
```

## API Documentation

After running the app, Swagger UI is available at `/swagger/index.html` (if HTTP endpoints are enabled).

## License

MIT
