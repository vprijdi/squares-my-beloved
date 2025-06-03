import styles from './StatisticsCard.module.css';
import {
  ArrowLeft,
  Trophy,
  Target,
  Calendar,
  Clock,
  TrendingUp,
  TrendingDown,
  Flame,
  Award,
  Star,
  Zap,
  Shield,
  Crown,
  Activity,
  BarChart3,
} from 'lucide-react';

export default function StatisticsCard() {
  const userStats = {
    highest_daily_score: 11,
    highest_daily_score_date: '2024-03-15',
    lowest_daily_score: -8,
    lowest_daily_score_date: '2024-02-20',
    total_tasks_completed: 247,
    total_tier0_completed: 156,
    total_tier1_completed: 67,
    total_tier2_completed: 24,
    total_active_days: 89,
    total_time_spent_minutes: 2340,
    avg_daily_score: 3.2,
    avg_weekly_score: 22.4,
    avg_monthly_score: 96.8,
    avg_tasks_per_active_day: 2.8,
    current_streak_days: 12,
    longest_streak_days: 28,
    longest_streak_start_date: '2024-01-15',
    longest_streak_end_date: '2024-02-12',
    total_score: 284,
  };

  const formatTime = (minutes) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
  };
  return (
    <div className={styles.container}>
      <h3 className={styles.title}>
        <BarChart3 size={20} />
        Your Statistics
      </h3>

      <div className={styles.statsContainer}>
        {/* Scores Section */}
        <div>
          <h4 className={styles.statHeading}>Scores</h4>

          <div className={styles.statList}>
            <div className={styles.statLine}>
              <span className={styles.statTitle}>Total Score</span>
              <span
                style={{
                  fontWeight: '600',
                  color: userStats.total_score >= 0 ? '#10b981' : '#ef4444',
                }}
              >
                {userStats.total_score >= 0 ? '+' : ''}
                {userStats.total_score}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Daily Average</span>
              <span
                style={{
                  fontWeight: '600',
                  color: userStats.avg_daily_score >= 0 ? '#10b981' : '#ef4444',
                }}
              >
                {userStats.avg_daily_score >= 0 ? '+' : ''}
                {userStats.avg_daily_score}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}> Weekly Average</span>
              <span
                style={{
                  fontWeight: '600',
                  color:
                    userStats.avg_weekly_score >= 0 ? '#10b981' : '#ef4444',
                }}
              >
                {userStats.avg_weekly_score >= 0 ? '+' : ''}
                {userStats.avg_weekly_score}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}> Highest Daily</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#10b981',
                }}
              >
                +{userStats.highest_daily_score}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Lowest Daily</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#ef4444',
                }}
              >
                {userStats.lowest_daily_score}
              </span>
            </div>
          </div>
        </div>

        {/* Tasks Section */}
        <div>
          <h4 className={styles.statHeading}>Tasks</h4>

          <div className={styles.statList}>
            <div className={styles.statLine}>
              <span className={styles.statTitle}>Total Completed</span>
              <span
                style={{
                  fontWeight: '600',
                }}
              >
                {userStats.total_tasks_completed}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}> Quick Tasks</span>
              <span
                style={{
                  fontSize: '0.875rem',
                  fontWeight: '600',
                  color: '#60a5fa',
                }}
              >
                {userStats.total_tier0_completed}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Medium Tasks</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#c084fc',
                }}
              >
                {userStats.total_tier1_completed}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Large Tasks</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#f472b6',
                }}
              >
                {userStats.total_tier2_completed}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Daily Average</span>
              <span
                style={{
                  color: '#f0f0f0',
                }}
              >
                {userStats.avg_tasks_per_active_day}
              </span>
            </div>
          </div>
        </div>

        {/* Activity Section */}
        <div>
          <h4 className={styles.statHeading}>Activity</h4>

          <div className={styles.statList}>
            <div className={styles.statLine}>
              <span className={styles.statTitle}>Active Days</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#f0f0f0',
                }}
              >
                {userStats.total_active_days}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Time Spent</span>
              <span>{formatTime(userStats.total_time_spent_minutes)}</span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Current Streak</span>
              <span>{userStats.current_streak_days} days</span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Longest Streak</span>
              <span>{userStats.longest_streak_days} days</span>
            </div>
          </div>
        </div>

        {/* Records Section */}
        <div>
          <h4 className={styles.statHeading}>Records</h4>

          <div className={styles.statList}>
            <div className={styles.statLine}>
              <span className={styles.statTitle}>Best Day</span>
              <span
                style={{
                  fontWeight: '600',
                  color: '#10b981',
                }}
              >
                {new Date(
                  userStats.highest_daily_score_date
                ).toLocaleDateString()}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Worst Day</span>
              <span
                style={{
                  fontSize: '0.875rem',
                  fontWeight: '600',
                  color: '#ef4444',
                }}
              >
                {new Date(
                  userStats.lowest_daily_score_date
                ).toLocaleDateString()}
              </span>
            </div>

            <div className={styles.statLine}>
              <span className={styles.statTitle}>Best Streak</span>
              <span>
                {new Date(
                  userStats.longest_streak_start_date
                ).toLocaleDateString()}{' '}
                -{' '}
                {new Date(
                  userStats.longest_streak_end_date
                ).toLocaleDateString()}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
