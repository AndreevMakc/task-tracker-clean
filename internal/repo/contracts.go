package repo

import (
	"context"
	"errors"
	"time"

	"task-tracker-clean/internal/entity"

	"github.com/google/uuid"
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
	//TODO: to be implement later for admins
	Delete(ctx context.Context, id uuid.UUID) error

	UpdateWithLock(ctx context.Context, id uuid.UUID, fn func(*entity.Task) error) (*entity.Task, error)
}
