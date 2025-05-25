package store

import (
	"context"
	"database/sql"
	"time"
)

type UserStatistics struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

	// All-time statistics
	HighestDailyScore     int        `json:"highest_daily_score"`
	HighestDailyScoreDate *time.Time `json:"highest_daily_score_date"`
	LowestDailyScore      *int       `json:"lowest_daily_score"`
	LowestDailyScoreDate  *time.Time `json:"lowest_daily_score_date"`

	// Tasks statistics
	TotalTasksCompleted int `json:"total_tasks_completed"`
	TotalTier0Completed int `json:"total_tier0_completed"`
	TotalTier1Completed int `json:"total_tier1_completed"`
	TotalTier2Completed int `json:"total_tier2_completed"`

	// Time statistics
	TotalActiveDays       int `json:"total_active_days"`
	TotalTimeSpentMinutes int `json:"total_time_spent_minutes"`

	// Averages
	AvgDailyScore        float64 `json:"avg_daily_score"`
	AvgWeeklyScore       float64 `json:"avg_weekly_score"`
	AvgMonthlyScore      float64 `json:"avg_monthly_score"`
	AvgTasksPerActiveDay float64 `json:"avg_tasks_per_active_day"`

	// Streaks
	CurrentStreakDays      int        `json:"current_streak_days"`
	LongestStreakDays      int        `json:"longest_streak_days"`
	LongestStreakStartDate *time.Time `json:"longest_streak_start_date"`
	LongestStreakEndDate   *time.Time `json:"longest_streak_end_date"`

	UpdatedAt time.Time `json:"updated_at"`
}

type StatisticsStore struct {
	db *sql.DB
}

func (s *StatisticsStore) Create(ctx context.Context, tx *sql.Tx, userID int64) (*UserStatistics, error) {
	query := `
        INSERT INTO user_statistics (user_id)
        VALUES ($1)
        RETURNING 
            id,
            highest_daily_score,
            highest_daily_score_date,
            lowest_daily_score,
            lowest_daily_score_date,
            total_tasks_completed,
            total_tier0_completed,
            total_tier1_completed,
            total_tier2_completed,
            total_active_days,
            total_time_spent_minutes,
            avg_daily_score,
            avg_weekly_score,
            avg_monthly_score,
            avg_tasks_per_active_day,
            current_streak_days,
            longest_streak_days,
            longest_streak_start_date,
            longest_streak_end_date,
            updated_at
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stats := &UserStatistics{UserID: userID}

	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&stats.ID,
		&stats.HighestDailyScore,
		&stats.HighestDailyScoreDate,
		&stats.LowestDailyScore,
		&stats.LowestDailyScoreDate,
		&stats.TotalTasksCompleted,
		&stats.TotalTier0Completed,
		&stats.TotalTier1Completed,
		&stats.TotalTier2Completed,
		&stats.TotalActiveDays,
		&stats.TotalTimeSpentMinutes,
		&stats.AvgDailyScore,
		&stats.AvgWeeklyScore,
		&stats.AvgMonthlyScore,
		&stats.AvgTasksPerActiveDay,
		&stats.CurrentStreakDays,
		&stats.LongestStreakDays,
		&stats.LongestStreakStartDate,
		&stats.LongestStreakEndDate,
		&stats.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *StatisticsStore) GetByUserID(ctx context.Context, userID int64) (*UserStatistics, error) {
	query := `
		SELECT 
			id, user_id,
			highest_daily_score, highest_daily_score_date,
			lowest_daily_score, lowest_daily_score_date,
			total_tasks_completed,
			total_tier0_completed, total_tier1_completed, total_tier2_completed,
			total_active_days, total_time_spent_minutes,
			avg_daily_score, avg_weekly_score, avg_monthly_score, avg_tasks_per_active_day,
			current_streak_days, longest_streak_days,
			longest_streak_start_date, longest_streak_end_date,
			updated_at
		FROM user_statistics
		WHERE user_id = $1
	`

	var stats UserStatistics
	var (
		highestDate, lowestDate sql.NullTime
		lowestScore             sql.NullInt64
		streakStart, streakEnd  sql.NullTime
	)

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.ID,
		&stats.UserID,
		&stats.HighestDailyScore,
		&highestDate,
		&lowestScore,
		&lowestDate,
		&stats.TotalTasksCompleted,
		&stats.TotalTier0Completed,
		&stats.TotalTier1Completed,
		&stats.TotalTier2Completed,
		&stats.TotalActiveDays,
		&stats.TotalTimeSpentMinutes,
		&stats.AvgDailyScore,
		&stats.AvgWeeklyScore,
		&stats.AvgMonthlyScore,
		&stats.AvgTasksPerActiveDay,
		&stats.CurrentStreakDays,
		&stats.LongestStreakDays,
		&streakStart,
		&streakEnd,
		&stats.UpdatedAt,
	)

	if highestDate.Valid {
		stats.HighestDailyScoreDate = &highestDate.Time
	}
	if lowestScore.Valid {
		val := int(lowestScore.Int64)
		stats.LowestDailyScore = &val
	}
	if lowestDate.Valid {
		stats.LowestDailyScoreDate = &lowestDate.Time
	}
	if streakStart.Valid {
		stats.LongestStreakStartDate = &streakStart.Time
	}
	if streakEnd.Valid {
		stats.LongestStreakEndDate = &streakEnd.Time
	}

	switch err {
	case nil:
		return &stats, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (s *StatisticsStore) Update(ctx context.Context, stats *UserStatistics) error {
	query := `
		INSERT INTO user_statistics (
			user_id,
			highest_daily_score, highest_daily_score_date,
			lowest_daily_score, lowest_daily_score_date,
			total_tasks_completed,
			total_tier0_completed, total_tier1_completed, total_tier2_completed,
			total_active_days, total_time_spent_minutes,
			avg_daily_score, avg_weekly_score, avg_monthly_score, avg_tasks_per_active_day,
			current_streak_days, longest_streak_days,
			longest_streak_start_date, longest_streak_end_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (user_id) DO UPDATE SET
			highest_daily_score = EXCLUDED.highest_daily_score,
			highest_daily_score_date = EXCLUDED.highest_daily_score_date,
			lowest_daily_score = EXCLUDED.lowest_daily_score,
			lowest_daily_score_date = EXCLUDED.lowest_daily_score_date,
			total_tasks_completed = EXCLUDED.total_tasks_completed,
			total_tier0_completed = EXCLUDED.total_tier0_completed,
			total_tier1_completed = EXCLUDED.total_tier1_completed,
			total_tier2_completed = EXCLUDED.total_tier2_completed,
			total_active_days = EXCLUDED.total_active_days,
			total_time_spent_minutes = EXCLUDED.total_time_spent_minutes,
			avg_daily_score = EXCLUDED.avg_daily_score,
			avg_weekly_score = EXCLUDED.avg_weekly_score,
			avg_monthly_score = EXCLUDED.avg_monthly_score,
			avg_tasks_per_active_day = EXCLUDED.avg_tasks_per_active_day,
			current_streak_days = EXCLUDED.current_streak_days,
			longest_streak_days = EXCLUDED.longest_streak_days,
			longest_streak_start_date = EXCLUDED.longest_streak_start_date,
			longest_streak_end_date = EXCLUDED.longest_streak_end_date,
			updated_at = NOW()
		RETURNING id, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var (
		lowestScore *int
		lowestDate  *time.Time
	)

	if stats.LowestDailyScore != nil {
		lowestScore = stats.LowestDailyScore
		lowestDate = stats.LowestDailyScoreDate
	}

	return s.db.QueryRowContext(
		ctx,
		query,
		stats.UserID,
		stats.HighestDailyScore,
		stats.HighestDailyScoreDate,
		lowestScore,
		lowestDate,
		stats.TotalTasksCompleted,
		stats.TotalTier0Completed,
		stats.TotalTier1Completed,
		stats.TotalTier2Completed,
		stats.TotalActiveDays,
		stats.TotalTimeSpentMinutes,
		stats.AvgDailyScore,
		stats.AvgWeeklyScore,
		stats.AvgMonthlyScore,
		stats.AvgTasksPerActiveDay,
		stats.CurrentStreakDays,
		stats.LongestStreakDays,
		stats.LongestStreakStartDate,
		stats.LongestStreakEndDate,
	).Scan(&stats.ID, &stats.UpdatedAt)
}

func (s *StatisticsStore) UpdateAllTimeStats(ctx context.Context, userID int64, highestScore int, highestDate time.Time, lowestScore *int, lowestDate *time.Time) error {
	query := `
        INSERT INTO user_statistics (user_id, highest_daily_score, highest_daily_score_date, lowest_daily_score, lowest_daily_score_date)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE SET
            highest_daily_score = GREATEST(user_statistics.highest_daily_score, EXCLUDED.highest_daily_score),
            highest_daily_score_date = CASE 
                WHEN EXCLUDED.highest_daily_score > user_statistics.highest_daily_score THEN EXCLUDED.highest_daily_score_date
                ELSE user_statistics.highest_daily_score_date
            END,
            lowest_daily_score = LEAST(COALESCE(user_statistics.lowest_daily_score, EXCLUDED.lowest_daily_score), EXCLUDED.lowest_daily_score),
            lowest_daily_score_date = CASE 
                WHEN EXCLUDED.lowest_daily_score < COALESCE(user_statistics.lowest_daily_score, EXCLUDED.lowest_daily_score) THEN EXCLUDED.lowest_daily_score_date
                ELSE user_statistics.lowest_daily_score_date
            END,
            updated_at = NOW()
    `

	_, err := s.db.ExecContext(ctx, query,
		userID,
		highestScore,
		highestDate,
		lowestScore,
		lowestDate,
	)
	return err
}

func (s *StatisticsStore) UpdateTaskStats(ctx context.Context, userID int64, tierCompletions map[int]int) error {
	query := `
        INSERT INTO user_statistics (user_id, total_tasks_completed, total_tier0_completed, total_tier1_completed, total_tier2_completed)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE SET
            total_tasks_completed = user_statistics.total_tasks_completed + EXCLUDED.total_tasks_completed,
            total_tier0_completed = user_statistics.total_tier0_completed + EXCLUDED.total_tier0_completed,
            total_tier1_completed = user_statistics.total_tier1_completed + EXCLUDED.total_tier1_completed,
            total_tier2_completed = user_statistics.total_tier2_completed + EXCLUDED.total_tier2_completed,
            updated_at = NOW()
    `

	total := tierCompletions[0] + tierCompletions[1] + tierCompletions[2]
	_, err := s.db.ExecContext(ctx, query,
		userID,
		total,
		tierCompletions[0],
		tierCompletions[1],
		tierCompletions[2],
	)
	return err
}

func (s *StatisticsStore) UpdateTimeStats(ctx context.Context, userID int64, activeDaysIncrement int, minutesSpent int) error {
	query := `
        INSERT INTO user_statistics (user_id, total_active_days, total_time_spent_minutes)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE SET
            total_active_days = user_statistics.total_active_days + EXCLUDED.total_active_days,
            total_time_spent_minutes = user_statistics.total_time_spent_minutes + EXCLUDED.total_time_spent_minutes,
            updated_at = NOW()
    `

	_, err := s.db.ExecContext(ctx, query,
		userID,
		activeDaysIncrement,
		minutesSpent,
	)
	return err
}

func (s *StatisticsStore) UpdateAverages(ctx context.Context, userID int64, dailyScore, weeklyScore, monthlyScore, tasksPerDay float64) error {
	query := `
        INSERT INTO user_statistics (user_id, avg_daily_score, avg_weekly_score, avg_monthly_score, avg_tasks_per_active_day)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE SET
            avg_daily_score = EXCLUDED.avg_daily_score,
            avg_weekly_score = EXCLUDED.avg_weekly_score,
            avg_monthly_score = EXCLUDED.avg_monthly_score,
            avg_tasks_per_active_day = EXCLUDED.avg_tasks_per_active_day,
            updated_at = NOW()
    `

	_, err := s.db.ExecContext(ctx, query,
		userID,
		dailyScore,
		weeklyScore,
		monthlyScore,
		tasksPerDay,
	)
	return err
}

func (s *StatisticsStore) UpdateStreaks(ctx context.Context, userID int64, currentStreak int, longestStreak int, streakStart *time.Time, streakEnd *time.Time) error {
	query := `
        INSERT INTO user_statistics (user_id, current_streak_days, longest_streak_days, longest_streak_start_date, longest_streak_end_date)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE SET
            current_streak_days = EXCLUDED.current_streak_days,
            longest_streak_days = GREATEST(user_statistics.longest_streak_days, EXCLUDED.longest_streak_days),
            longest_streak_start_date = CASE
                WHEN EXCLUDED.longest_streak_days > user_statistics.longest_streak_days THEN EXCLUDED.longest_streak_start_date
                ELSE user_statistics.longest_streak_start_date
            END,
            longest_streak_end_date = CASE
                WHEN EXCLUDED.longest_streak_days > user_statistics.longest_streak_days THEN EXCLUDED.longest_streak_end_date
                ELSE user_statistics.longest_streak_end_date
            END,
            updated_at = NOW()
    `

	_, err := s.db.ExecContext(ctx, query,
		userID,
		currentStreak,
		longestStreak,
		streakStart,
		streakEnd,
	)
	return err
}
