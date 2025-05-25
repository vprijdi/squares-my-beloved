package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type UserService struct {
	store *store.Storage
	db    *sql.DB
}

var (
	ErrUserActivationFailed   = errors.New("failed to activate user")
	ErrProgressionSetupFailed = errors.New("failed to setup user progression")
	ErrStatisticsSetupFailed  = errors.New("failed to setup user statistics")
)

func (s *UserService) ActivateAndSetup(ctx context.Context, token string) error {
	return store.WithTx(s.db, ctx, func(tx *sql.Tx) error {
		userID, err := s.store.Users.Activate(ctx, tx, token)
		if err != nil {
			return err
		}

		if _, err := s.store.Progression.Create(ctx, tx, userID); err != nil {
			return err
		}

		if _, err := s.store.Statistics.Create(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}
