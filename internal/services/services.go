package services

import (
	"context"
	"database/sql"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type Services struct {
	UserServices interface {
		ActivateAndSetup(ctx context.Context, token string) error
	}
}

func NewServices(store *store.Storage, db *sql.DB) Services {
	return Services{
		UserServices: &UserService{
			store: store,
			db:    db,
		},
	}
}
