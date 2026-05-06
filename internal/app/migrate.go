//go:build migrate

package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"task-tracker-clean/config"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
	_migrationsPath  = "file://migrations"
)

// Init applies all up migrations from the migrations folder.
func Init(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(_defaultAttempts)*_defaultTimeout)
	defer cancel()

	var (
		m   *migrate.Migrate
		err error
	)

	for attempts := _defaultAttempts; attempts > 0; attempts-- {
		m, err = migrate.New(
			_migrationsPath,
			cfg.PG.URL,
		)
		if err == nil {
			break
		}

		log.Printf("failed to create migrate instance (attempts left: %d): %v", attempts, err)

		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while creating migrate instance: %w", ctx.Err())
		case <-time.After(_defaultTimeout):
		}
	}

	if err != nil {
		return fmt.Errorf("failed to create migrate instance after retries: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("migrations applied successfully")
	return nil
}
