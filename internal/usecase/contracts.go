package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"task-tracker-clean/internal/entity"
)

type TaskFilter struct {
	Status    *entity.TaskStatus
	Name      string
	CreatedAt TimeRange
}

type TimeRange struct {
	From *time.Time
	To   *time.Time
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, title string) (*entity.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	ListTasks(ctx context.Context, filter TaskFilter) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, title *string, status *entity.TaskStatus) (*entity.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
