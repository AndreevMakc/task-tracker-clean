package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"task-tracker-clean/internal/entity"
)

var ErrNotFound = errors.New("entity not found")

type TimeRange struct {
	From *time.Time
	To   *time.Time
}

type TaskFilter struct {
	Status    *entity.TaskStatus
	Name      string
	CreatedAt TimeRange
}

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	List(ctx context.Context, filter TaskFilter) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}