package persistent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-tracker-clean/internal/entity"
	"task-tracker-clean/internal/repo"
)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *entity.Task) error {
	query := `
		INSERT INTO tasks (id, title, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("task create: %w", err)
	}
	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	query := `
		SELECT id, title, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	var task entity.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("task getbyid: %w", err)
	}
	return &task, nil
}

func (r *TaskRepo) List(ctx context.Context, filter repo.TaskFilter) ([]entity.Task, error) {
	var conditions []string
	var args []any

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *filter.Status)
	}
	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.CreatedAt.From != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, *filter.CreatedAt.From)
	}
	if filter.CreatedAt.To != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", len(args)+1))
		args = append(args, *filter.CreatedAt.To)
	}

	query := `
		SELECT id, title, status, created_at, updated_at
		FROM tasks
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("task list: %w", err)
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("task list scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("task list rows: %w", err)
	}
	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *entity.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, status = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(ctx, query,
		task.Title,
		task.Status,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("task update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repo.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("task delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repo.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) UpdateWithLock(ctx context.Context, id uuid.UUID, fn func(*entity.Task) error) (*entity.Task, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("task begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, title, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
		FOR UPDATE
	`
	var task entity.Task
	err = tx.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("task getbyid for update: %w", err)
	}

	if err := fn(&task); err != nil {
		return nil, err
	}

	updateQuery := `
		UPDATE tasks
		SET title = $1, status = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := tx.Exec(ctx, updateQuery,
		task.Title,
		task.Status,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("task update tx: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, repo.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("task commit tx: %w", err)
	}

	return &task, nil
}