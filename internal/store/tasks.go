package store

import (
	"context"
	"database/sql"
)

type TaskStore struct {
	db *sql.DB
}

type Task struct {
	ID               int64  `json:"id"`
	UserID           int64  `json:"user_id"`
	Title            string `json:"title"`
	IsOptional       bool   `json:"is_optional"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	CompletionsCount int32  `json:"completions_count"`
}

func (s *TaskStore) create(ctx context.Context, tx *sql.Tx, task *Task) error {
	query := `
        INSERT INTO tasks (user_id, title, is_optional)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at, completions_count
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.IsOptional,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CompletionsCount,
	)

	return err
}

func (s *TaskStore) CreateTasks(ctx context.Context, tasks []*Task) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		for _, task := range tasks {
			if err := s.create(ctx, tx, task); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TaskStore) GetByID(ctx context.Context, taskID int64) (*Task, error) {
	query := `
	SELECT id, user_id, title, is_optional, created_at, updated_at, completions_count
	FROM tasks
	WHERE id = $1
	`

	var task Task
	err := s.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.IsOptional,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CompletionsCount,
	)

	switch err {
	case nil:
		return &task, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (s *TaskStore) Complete(ctx context.Context, taskID int64) error {
	query := `
        UPDATE tasks 
        SET 
            completions_count = completions_count + 1,
            updated_at = NOW()
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
