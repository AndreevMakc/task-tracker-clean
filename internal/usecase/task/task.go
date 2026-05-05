package task

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"task-tracker-clean/internal/entity"
	"task-tracker-clean/internal/repo"
	"task-tracker-clean/internal/usecase"
)

type taskUsecase struct {
	taskRepo repo.TaskRepository
}

func NewTaskUsecase(tr repo.TaskRepository) usecase.TaskUsecase {
	return &taskUsecase{taskRepo: tr}
}

func (uc *taskUsecase) CreateTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		ID:        uuid.New(),
		Title:     title,
		Status:    entity.TaskStatusToDo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := uc.taskRepo.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (uc *taskUsecase) GetTask(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	t, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, entity.ErrTaskNotFound
		}
		return nil, err
	}
	return t, nil
}

func toRepoFilter(f usecase.TaskFilter) repo.TaskFilter {
	return repo.TaskFilter{
		Status:    f.Status,
		Name:      f.Name,
		CreatedAt: repo.TimeRange{From: f.CreatedAt.From, To: f.CreatedAt.To},
	}
}

func (uc *taskUsecase) ListTasks(ctx context.Context, filter usecase.TaskFilter) ([]entity.Task, error) {
	return uc.taskRepo.List(ctx, toRepoFilter(filter))
}

func (uc *taskUsecase) UpdateTask(ctx context.Context, id uuid.UUID, title *string, status *entity.TaskStatus) (*entity.Task, error) {
	t, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, entity.ErrTaskNotFound
		}
		return nil, err
	}

	if title != nil {
		if *title == "" {
			return nil, entity.ErrMissingTaskTitle
		}
		t.Title = *title
	}

	if status != nil {
		if err := t.Transition(*status); err != nil {
			return nil, err
		}
	}

	t.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (uc *taskUsecase) DeleteTask(ctx context.Context, id uuid.UUID) error {
	t, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return entity.ErrTaskNotFound
		}
		return err
	}

	if err := t.Transition(entity.TaskStatusTrashed); err != nil {
		return err
	}

	t.UpdatedAt = time.Now()

	return uc.taskRepo.Update(ctx, t)
}
