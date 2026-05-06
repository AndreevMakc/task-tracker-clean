package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"task-tracker-clean/internal/controller/restapi/v1/request"
	"task-tracker-clean/internal/controller/restapi/v1/response"
	"task-tracker-clean/internal/entity"
	"task-tracker-clean/internal/usecase"
)

type TaskHandler struct {
	uc usecase.TaskUsecase
}

func NewTaskHandler(r chi.Router, uc usecase.TaskUsecase) {
	h := &TaskHandler{uc: uc}

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

func parseUUIDFromPath(r *http.Request, key string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, key)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, entity.ErrInvalidTaskID
	}
	return id, nil
}

// Create godoc
// @Summary     Create a new task
// @Description Create a new task with the given title
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       input body     request.CreateTask true "Task title"
// @Success     201   {object} entity.Task
// @Failure     400   {object} response.Error
// @Failure     500   {object} response.Error
// @Router      /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	task, err := h.uc.CreateTask(r.Context(), req.Title)
	if err != nil {
		handleUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

// Get godoc
// @Summary     Get a task by ID
// @Description Get a task by its UUID
// @Tags        tasks
// @Produce     json
// @Param       id   path     string true "Task ID"
// @Success     200  {object} entity.Task
// @Failure     400  {object} response.Error
// @Failure     404  {object} response.Error
// @Failure     500  {object} response.Error
// @Router      /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromPath(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.uc.GetTask(r.Context(), id)
	if err != nil {
		handleUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// List godoc
// @Summary     List tasks
// @Description List tasks with optional filters
// @Tags        tasks
// @Produce     json
// @Param       status query    string false "Filter by status (todo, in_progress, done, trashed)"
// @Param       name   query    string false "Filter by name (partial match)"
// @Success     200    {array}  entity.Task
// @Failure     400    {object} response.Error
// @Failure     500    {object} response.Error
// @Router      /tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := usecase.TaskFilter{}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		s := entity.TaskStatus(statusStr)
		if s.Valid() {
			filter.Status = &s
		} else {
			writeError(w, http.StatusBadRequest, "invalid status value")
			return
		}
	}

	filter.Name = r.URL.Query().Get("name")

	tasks, err := h.uc.ListTasks(r.Context(), filter)
	if err != nil {
		handleUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

// Update godoc
// @Summary     Update a task
// @Description Update task title and/or status
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id   path     string         true "Task ID"
// @Param       input body     request.UpdateTask true "Update data"
// @Success     200  {object} entity.Task
// @Failure     400  {object} response.Error
// @Failure     404  {object} response.Error
// @Failure     500  {object} response.Error
// @Router      /tasks/{id} [put]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromPath(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req request.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != nil && *req.Title == "" {
		writeError(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	var status *entity.TaskStatus
	if req.Status != nil {
		s := entity.TaskStatus(*req.Status)
		if !s.Valid() {
			writeError(w, http.StatusBadRequest, "invalid task status")
			return
		}
		status = &s
	}

	task, err := h.uc.UpdateTask(r.Context(), id, req.Title, status)
	if err != nil {
		handleUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// Delete godoc
// @Summary     Delete a task
// @Description Mark a task as trashed (soft delete)
// @Tags        tasks
// @Param       id   path     string true "Task ID"
// @Success     204  "No Content"
// @Failure     400  {object} response.Error
// @Failure     404  {object} response.Error
// @Failure     500  {object} response.Error
// @Router      /tasks/{id} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromPath(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.uc.DeleteTask(r.Context(), id); err != nil {
		handleUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, response.Error{Error: message})
}

func handleUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "task not found")
	case errors.Is(err, entity.ErrInvalidTaskStatus):
		writeError(w, http.StatusBadRequest, "invalid task status")
	case errors.Is(err, entity.ErrMissingTaskTitle):
		writeError(w, http.StatusBadRequest, "missing task title")
	case errors.Is(err, entity.ErrInvalidTransition):
		writeError(w, http.StatusBadRequest, "invalid status transition")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
