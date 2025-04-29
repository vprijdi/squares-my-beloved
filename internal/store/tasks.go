package store

import (
	"context"
	"database/sql"
)

type TaskStore struct {
	db *sql.DB
}

func (s *TaskStore) Create(ctx context.Context) error {
	return nil
}
