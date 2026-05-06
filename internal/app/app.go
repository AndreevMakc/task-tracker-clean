package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-tracker-clean/config"
	"task-tracker-clean/internal/controller/restapi"
	"task-tracker-clean/internal/repo/persistent"
	taskusecase "task-tracker-clean/internal/usecase/task"
	"task-tracker-clean/pkg/httpserver"
	"task-tracker-clean/pkg/postgres"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	log.Printf("starting %s", a.cfg.App.Name)

	ctx := context.Background()

	pg, err := postgres.New(ctx, a.cfg.PG.URL,
		postgres.MaxPoolSize(a.cfg.PG.MaxPoolSize),
		postgres.ConnAttempts(a.cfg.PG.ConnAttempts),
		postgres.ConnTimeout(a.cfg.PG.ConnTimeout),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pg.Close()

	taskRepo := persistent.NewTaskRepo(pg.Pool())
	taskUC := taskusecase.NewTaskUsecase(taskRepo)

	router := restapi.NewRouter(taskUC)

	server := httpserver.New(
		router,
		httpserver.Port(a.cfg.HTTP.Port),
		httpserver.ReadTimeout(a.cfg.HTTP.Timeout),
		httpserver.WriteTimeout(a.cfg.HTTP.Timeout),
		httpserver.IdleTimeout(a.cfg.HTTP.IdleTimeout),
	)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("starting http server on port %s", a.cfg.HTTP.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		// Server crashed on its own, graceful shutdown is meaningless
		return fmt.Errorf("http server error: %w", err)
	case <-quit:
	}

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
	return nil
}
