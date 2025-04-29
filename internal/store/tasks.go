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

func (s *TaskStore) Create(ctx context.Context, task *Task) error {
	query := `
        INSERT INTO tasks (user_id, title, is_optional)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at, completions_count
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
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
