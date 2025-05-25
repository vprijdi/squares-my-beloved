package services

import (
	"context"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type ProgressionkService struct {
	store *store.Storage
}

func (s *ProgressionkService) Create(context.Context, *store.User) error {
	return nil
}
