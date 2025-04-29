package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create(context.Context) error
	}
	Tasks interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: UserStore{db},
		Tasks: UserStore{db},
	}
}
