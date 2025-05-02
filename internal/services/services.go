package services

import (
	"context"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type Services struct {
	TaskServices interface {
		Create(context.Context, *store.User) error
	}
}

func NewServices(store *store.Storage) Services {
	return Services{
		TaskServices: &TaskService{store: store},
	}
}
