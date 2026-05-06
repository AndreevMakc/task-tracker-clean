//go:build migrate

package main

import (
	"log"

	"task-tracker-clean/config"
	"task-tracker-clean/internal/app"
)

func main() {
	cfg := config.New()

	if err := app.Init(cfg); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
