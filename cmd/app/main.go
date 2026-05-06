package main

import (
	"log"

	"task-tracker-clean/config"
	"task-tracker-clean/internal/app"
)

func main() {
	cfg := config.New()

	if err := app.New(cfg).Run(); err != nil {
		log.Fatalf("app failed: %v", err)
	}
}
