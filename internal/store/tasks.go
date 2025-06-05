package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
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
	Completed      []bool   `json:"completed"`
	Tier             int32  `json:"tier"`
}

func (s *TaskStore) create(ctx context.Context, tx *sql.Tx, task *Task) error {
	query := `
        INSERT INTO tasks (user_id, title, is_optional, tier)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at, completed, completions_count
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.IsOptional,
		task.Tier,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
		pq.Array(&task.Completed),
		&task.CompletionsCount,
	)

	return err
}

func (s *TaskStore) CreateTasks(ctx context.Context, tasks []*Task) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
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
	SELECT id, user_id, title, is_optional, created_at, updated_at, completions_count, completed, tier
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
		pq.Array(&task.Completed),
		&task.Tier,
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

func (s *TaskStore) Complete(ctx context.Context, taskID int64, newCompleted[]bool) error {
	query := `
        UPDATE tasks 
        SET 
            completions_count = completions_count + 1,
            updated_at = NOW(),
			completed = $2 
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, taskID, pq.Array(newCompleted))
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

func (s *TaskStore) GetAllUserTasks(ctx context.Context, userID int64, tf *TaskFilters) ([]Task, error) {
	query := `
	SELECT t.id, t.user_id, t.title, t.is_optional, t.created_at, t.updated_at, 
		   t.completions_count, t.completed, t.tier
	FROM tasks t
	WHERE t.user_id = $1
	  AND ($4::text IS NULL OR t.created_at::date = $4::date)
	ORDER BY t.created_at DESC
	LIMIT $2
	OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, tf.Limit, tf.Offset, tf.Date )
	if err != nil {
		return nil, err
	}

	var filteredResult []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.IsOptional,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.CompletionsCount,
			pq.Array(&t.Completed),
			&t.Tier,
		)
		if err != nil {
			return nil, err
		}
		filteredResult = append(filteredResult, t)
	}

	return filteredResult, err
}
