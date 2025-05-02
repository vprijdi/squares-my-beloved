package services

import (
	"context"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type TaskService struct {
	store *store.Storage
}

func (s *TaskService) Create(context.Context, *store.User) error {
	return nil
}
