package entity

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the current state of a task.
type TaskStatus string // @name entity.TaskStatus

const (
	TaskStatusToDo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusTrashed    TaskStatus = "trashed"
)

// Valid reports whether s is a known task status.
func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusToDo, TaskStatusInProgress, TaskStatusDone, TaskStatusTrashed:
		return true
	default:
		return false
	}
}

// Task represents a task in the system.
type Task struct {
	ID        uuid.UUID  `json:"id"          example:"550e8400-e29b-41d4-a716-446655440000"`
	Title     string     `json:"title"       example:"My task"`
	Status    TaskStatus `json:"status"      example:"todo"`
	CreatedAt time.Time  `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
} // @name entity.Task

// validTransitions defines allowed status transitions.
// This map is read-only and should not be modified at runtime.
var validTransitions = map[TaskStatus][]TaskStatus{
	TaskStatusToDo:       {TaskStatusInProgress, TaskStatusTrashed},
	TaskStatusInProgress: {TaskStatusDone, TaskStatusToDo, TaskStatusTrashed},
	TaskStatusDone:       {TaskStatusTrashed},
	TaskStatusTrashed:    {},
}

// Transition validates and applies a status transition.
// Note: UpdatedAt is not modified here; the caller (usecase/repository) is responsible for updating it.
func (t *Task) Transition(newStatus TaskStatus) error {
	if !newStatus.Valid() {
		return ErrInvalidTaskStatus
	}

	if t.Status == newStatus {
		return nil
	}

	allowed := validTransitions[t.Status]
	for _, s := range allowed {
		if s == newStatus {
			t.Status = newStatus
			return nil
		}
	}

	return ErrInvalidTransition
}

// Validate checks the task fields are valid.
func (t *Task) Validate() error {
	if t.ID == uuid.Nil {
		return ErrMissingTaskID
	}
	if t.Title == "" {
		return ErrMissingTaskTitle
	}
	if !t.Status.Valid() {
		return ErrInvalidTaskStatus
	}
	return nil
}
