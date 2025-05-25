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
		Activate(ctx context.Context, tx *sql.Tx, token string) (int64, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
	}
	Tasks interface {
		CreateTasks(context.Context, []*Task) error
		GetByID(ctx context.Context, taskID int64) (*Task, error)
		Complete(ctx context.Context, taskID int64) error
		GetAllUserTasks(ctx context.Context, userID int64, tf *TaskFilters) ([]Task, error)
	}
	Goals interface {
		CreateGoals(context.Context, []*Goal) error
		GetByID(ctx context.Context, goalID int64) (*Goal, error)
		Achieve(ctx context.Context, goalID int64) error
		GetAllUserGoals(ctx context.Context, userID int64, tf *GoalFilters) ([]Goal, error)
	}
	Progression interface {
		Create(ctx context.Context, tx *sql.Tx, userID int64) (*UserProgression, error)
		Update(ctx context.Context, progression *UserProgression) error
		GetByUserID(ctx context.Context, userID int64) (*UserProgression, error)
	}

	Statistics interface {
		// Core operations
		Create(ctx context.Context, tx *sql.Tx, userID int64) (*UserStatistics, error)
		GetByUserID(ctx context.Context, userID int64) (*UserStatistics, error)
		Update(ctx context.Context, stats *UserStatistics) error

		// Section-specific updates
		UpdateAllTimeStats(
			ctx context.Context,
			userID int64,
			highestScore int,
			highestDate time.Time,
			lowestScore *int,
			lowestDate *time.Time,
		) error

		UpdateTaskStats(
			ctx context.Context,
			userID int64,
			tierCompletions map[int]int,
		) error

		UpdateTimeStats(
			ctx context.Context,
			userID int64,
			activeDaysIncrement int,
			minutesSpent int,
		) error

		UpdateAverages(
			ctx context.Context,
			userID int64,
			dailyScore, weeklyScore, monthlyScore, tasksPerDay float64,
		) error

		UpdateStreaks(
			ctx context.Context,
			userID int64,
			currentStreak int,
			longestStreak int,
			streakStart *time.Time,
			streakEnd *time.Time,
		) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:       &UserStore{db},
		Tasks:       &TaskStore{db},
		Goals:       &GoalStore{db},
		Progression: &ProgressionStore{db},
		Statistics:  &StatisticsStore{db},
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
