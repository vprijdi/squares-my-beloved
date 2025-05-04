CREATE TABLE user_statistics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,

    -- All-time statistics
    highest_daily_score INTEGER DEFAULT 0,
    highest_daily_score_date DATE,
    lowest_daily_score INTEGER,
    lowest_daily_score_date DATE,

    -- Tasks statistics
    total_tasks_completed INTEGER NOT NULL DEFAULT 0,
    total_tier0_completed INTEGER NOT NULL DEFAULT 0,
    total_tier1_completed INTEGER NOT NULL DEFAULT 0,
    total_tier2_completed INTEGER NOT NULL DEFAULT 0,
    
    -- Time statistics
    total_active_days INTEGER NOT NULL DEFAULT 0,
    total_time_spent_minutes INTEGER NOT NULL DEFAULT 0,
    
    -- Averages
    avg_daily_score FLOAT NOT NULL DEFAULT 0,
    avg_weekly_score FLOAT NOT NULL DEFAULT 0,
    avg_monthly_score FLOAT NOT NULL DEFAULT 0,
    avg_tasks_per_active_day FLOAT NOT NULL DEFAULT 0,
    
    -- Streaks
    current_streak_days INTEGER NOT NULL DEFAULT 0,
    longest_streak_days INTEGER NOT NULL DEFAULT 0,
    longest_streak_start_date DATE,
    longest_streak_end_date DATE,
    
    -- Last updated
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_user_statistics_user_id ON user_statistics (user_id);