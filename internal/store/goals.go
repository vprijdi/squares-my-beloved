package store

import (
	"context"
	"database/sql"
)

type GoalStore struct {
	db *sql.DB
}

type Goal struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	Title      string  `json:"title"`
	Achieved   bool    `json:"achieved"`
	CreatedAt  string  `json:"created_at"`
	AchievedAt *string `json:"achieved_at"`
}

func (s *GoalStore) create(ctx context.Context, tx *sql.Tx, goal *Goal) error {
	query := `
        INSERT INTO goals (user_id, title, achieved)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, achieved_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var achievedAt sql.NullString
	err := tx.QueryRowContext(
		ctx,
		query,
		goal.UserID,
		goal.Title,
		goal.Achieved,
	).Scan(
		&goal.ID,
		&goal.CreatedAt,
		&achievedAt,
	)

	if achievedAt.Valid {
		goal.AchievedAt = &achievedAt.String
	}

	return err
}

func (s *GoalStore) CreateGoals(ctx context.Context, goals []*Goal) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		for _, goal := range goals {
			if err := s.create(ctx, tx, goal); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *GoalStore) GetByID(ctx context.Context, goalID int64) (*Goal, error) {
	query := `
	SELECT id, user_id, title, achieved, created_at, achieved_at
	FROM goals
	WHERE id = $1
	`

	var goal Goal
	var achievedAt sql.NullString
	err := s.db.QueryRowContext(ctx, query, goalID).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Title,
		&goal.Achieved,
		&goal.CreatedAt,
		&achievedAt,
	)

	if achievedAt.Valid {
		goal.AchievedAt = &achievedAt.String
	}

	switch err {
	case nil:
		return &goal, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (s *GoalStore) Achieve(ctx context.Context, goalID int64) error {
	query := `
        UPDATE goals 
        SET 
            achieved = TRUE,
            achieved_at = NOW()
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, goalID)
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

// FIXME
func (s *GoalStore) GetAllUserGoals(ctx context.Context, userID int64, gf *GoalFilters) ([]Goal, error) {
	query := `
	SELECT id, user_id, title, achieved, created_at, achieved_at
	FROM goals
	WHERE user_id = $1
	  AND ($4::text IS NULL OR created_at::date = $4::date)
	  AND ($5::boolean IS NULL OR achieved = $5)
	ORDER BY created_at DESC
	LIMIT $2
	OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, gf.Limit, gf.Offset, gf.Achieved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var g Goal
		var achievedAt sql.NullString
		err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Title,
			&g.Achieved,
			&g.CreatedAt,
			&achievedAt,
		)
		if err != nil {
			return nil, err
		}
		if achievedAt.Valid {
			g.AchievedAt = &achievedAt.String
		}
		goals = append(goals, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return goals, nil
}
