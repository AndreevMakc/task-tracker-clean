package restapi

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"task-tracker-clean/internal/controller/restapi/v1"
	"task-tracker-clean/internal/usecase"
)

// @title       Task Tracker API
// @version     1.0
// @description API for managing tasks in a task tracker application.
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(uc usecase.TaskUsecase) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		v1.NewTaskHandler(r, uc)
	})

	return r
}
