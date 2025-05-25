package store

import (
	"context"
	"database/sql"
)

type UserProgression struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Level      int    `json:"level"`
	CurrentExp int    `json:"current_exp"`
	TotalScore int    `json:"total_score"`
	UpdatedAt  string `json:"updated_at"`
}

type ProgressionStore struct {
	db *sql.DB
}

const (
	defaultLevel = 1
	defaultExp   = 0
	defaultScore = 0
)

func (s *ProgressionStore) Create(ctx context.Context, tx *sql.Tx, userID int64) (*UserProgression, error) {
	query := `
        INSERT INTO user_progression (user_id, level, current_exp, total_score)
        VALUES ($1, $2, $3, $4)
        RETURNING id, updated_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	progression := &UserProgression{
		UserID:     userID,
		Level:      defaultLevel,
		CurrentExp: defaultExp,
		TotalScore: defaultScore,
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		userID,
		progression.Level,
		progression.CurrentExp,
		progression.TotalScore,
	).Scan(
		&progression.ID,
		&progression.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return progression, nil
}

func (s *ProgressionStore) Update(ctx context.Context, progression *UserProgression) error {
	query := `
        UPDATE user_progression
        SET 
            level = $2,
            current_exp = $3,
            total_score = $4,
            updated_at = NOW()
        WHERE user_id = $1
        RETURNING updated_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return s.db.QueryRowContext(
		ctx,
		query,
		progression.UserID,
		progression.Level,
		progression.CurrentExp,
		progression.TotalScore,
	).Scan(&progression.UpdatedAt)
}

func (s *ProgressionStore) GetByUserID(ctx context.Context, userID int64) (*UserProgression, error) {
	query := `
        SELECT id, user_id, level, current_exp, total_score, updated_at
        FROM user_progression
        WHERE user_id = $1
    `

	var progression UserProgression
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&progression.ID,
		&progression.UserID,
		&progression.Level,
		&progression.CurrentExp,
		&progression.TotalScore,
		&progression.UpdatedAt,
	)

	switch err {
	case nil:
		return &progression, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
