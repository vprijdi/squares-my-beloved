package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(ctx context.Context, userID int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(context.Context, string) error
		GetByEmail(ctx context.Context, email string) (*User, error)
	}
	Tasks interface {
		CreateTasks(context.Context, []*Task) error
		GetByID(ctx context.Context, taskID int64) (*Task, error)
		Complete(ctx context.Context, taskID int64) error
		GetAllUserTasks(ctx context.Context, userID int64, tf *TaskFilters) ([]Task, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
		Tasks: &TaskStore{db},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
