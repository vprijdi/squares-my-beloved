import React, { useState } from 'react';
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
import ScoreTrendGraph from '../Analytics/ScoreView/ScoreTrendGraph/ScoreTrendGraph';

const OtherPage = () => {
  const [chartView, setChartView] = useState('graph'); // 'graph' or 'heatmap'

  // Mock user statistics based on the table structure
  const userStats = {
    highest_daily_score: 15,
    highest_daily_score_date: '2024-03-15',
    lowest_daily_score: -8,
    lowest_daily_score_date: '2024-02-20',
    total_tasks_completed: 247,
    total_tier0_completed: 156, // Quick tasks
    total_tier1_completed: 67, // Medium tasks
    total_tier2_completed: 24, // Large tasks
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

  // Badge definitions with earned counts
  const badges = [
    {
      id: 1,
      icon: Flame,
      count: 3,
      title: 'Week Warrior',
      description: 'Maintain positive score for 7 days',
      type: 'positive',
    },
    {
      id: 2,
      icon: Trophy,
      count: 1,
      title: 'Century Club',
      description: 'Reach total score of 100',
      type: 'positive',
    },
    {
      id: 3,
      icon: Target,
      count: 2,
      title: 'Task Master',
      description: 'Complete 50 tasks',
      type: 'positive',
    },
    {
      id: 4,
      icon: Crown,
      count: 1,
      title: 'Streak King',
      description: 'Maintain 30-day streak',
      type: 'positive',
    },
    {
      id: 5,
      icon: Zap,
      count: 1,
      title: 'Lightning Round',
      description: 'Complete 10 tasks in one day',
      type: 'positive',
    },
    {
      id: 6,
      icon: Shield,
      count: 2,
      title: 'Consistency',
      description: 'Active for 30 days',
      type: 'positive',
    },
    {
      id: 7,
      icon: TrendingDown,
      count: 1,
      title: 'Rock Bottom',
      description: 'Reach total score of -50',
      type: 'negative',
    },
    {
      id: 8,
      icon: Activity,
      count: 1,
      title: 'Rough Week',
      description: 'Negative score for 7 days',
      type: 'negative',
    },
  ];

  // Generate heatmap data (last 365 days)
  const generateHeatmapData = () => {
    const data = [];
    const today = new Date();

    for (let i = 364; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);
      const score = Math.max(0, Math.floor(Math.random() * 10)); // 0 to 14 points per day

      data.push({
        date: date.toISOString().split('T')[0],
        score: score,
      });
    }
    return data;
  };

  const heatmapData = generateHeatmapData();

  const getHeatmapColor = (score) => {
    if (score === 0) return '#161616';
    if (score <= 3) return '#f472b6';
    if (score <= 6) return '#ec4899';
    if (score <= 9) return '#db2777';
    if (score <= 12) return '#be185d';
    return '#9d174d';
  };

  const formatTime = (minutes) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
  };

  return (
    <div
      style={{
        backgroundColor: '#0f0f0f',
        color: '#f0f0f0',
        minHeight: '100vh',
        fontFamily: 'system-ui, -apple-system, sans-serif',
      }}
    >
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
        {/* Header */}
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '1rem',
            marginBottom: '2rem',
          }}
        >
          <button
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              padding: '0.5rem 1rem',
              backgroundColor: 'transparent',
              border: '1px solid #2e2e2e',
              borderRadius: '0.5rem',
              color: '#f0f0f0',
              cursor: 'pointer',
              fontSize: '0.875rem',
            }}
          >
            <ArrowLeft size={16} />
            Back to Dashboard
          </button>
          <h1
            style={{
              fontSize: '2rem',
              fontWeight: 'bold',
              margin: 0,
              color: '#f0f0f0',
            }}
          >
            Analytics
          </h1>
        </div>

        <div
          style={{
            display: 'grid',
            gridTemplateColumns: '2fr 1fr',
            gap: '2rem',
            marginBottom: '2rem',
          }}
        >
          {/* Statistics Card */}
          <div
            style={{
              backgroundColor: '#1a1a1a',
              border: '1px solid #2e2e2e',
              borderRadius: '0.75rem',
              padding: '1.5rem',
            }}
          >
            <h2
              style={{
                fontSize: '1.25rem',
                fontWeight: 'bold',
                marginBottom: '1.5rem',
                color: '#f0f0f0',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
              }}
            >
              <BarChart3 size={20} />
              Your Statistics
            </h2>

            <div
              style={{
                display: 'grid',
                gridTemplateColumns: '1fr 1fr',
                gap: '1.5rem',
              }}
            >
              {/* Scores Section */}
              <div>
                <h3
                  style={{
                    fontSize: '0.875rem',
                    fontWeight: '600',
                    color: '#9ca3af',
                    marginBottom: '0.75rem',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}
                >
                  Scores
                </h3>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '0.5rem',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Total Score
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color:
                          userStats.total_score >= 0 ? '#10b981' : '#ef4444',
                      }}
                    >
                      {userStats.total_score >= 0 ? '+' : ''}
                      {userStats.total_score}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Daily Average
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color:
                          userStats.avg_daily_score >= 0
                            ? '#10b981'
                            : '#ef4444',
                      }}
                    >
                      {userStats.avg_daily_score >= 0 ? '+' : ''}
                      {userStats.avg_daily_score}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Weekly Average
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color:
                          userStats.avg_weekly_score >= 0
                            ? '#10b981'
                            : '#ef4444',
                      }}
                    >
                      {userStats.avg_weekly_score >= 0 ? '+' : ''}
                      {userStats.avg_weekly_score}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Highest Daily
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#10b981',
                      }}
                    >
                      +{userStats.highest_daily_score}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Lowest Daily
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
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
                <h3
                  style={{
                    fontSize: '0.875rem',
                    fontWeight: '600',
                    color: '#9ca3af',
                    marginBottom: '0.75rem',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}
                >
                  Tasks
                </h3>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '0.5rem',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Total Completed
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f472b6',
                      }}
                    >
                      {userStats.total_tasks_completed}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#60a5fa' }}>
                      Quick Tasks
                    </span>
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
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#c084fc' }}>
                      Medium Tasks
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#c084fc',
                      }}
                    >
                      {userStats.total_tier1_completed}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#f472b6' }}>
                      Large Tasks
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f472b6',
                      }}
                    >
                      {userStats.total_tier2_completed}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Daily Average
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
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
                <h3
                  style={{
                    fontSize: '0.875rem',
                    fontWeight: '600',
                    color: '#9ca3af',
                    marginBottom: '0.75rem',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}
                >
                  Activity
                </h3>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '0.5rem',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Active Days
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f0f0f0',
                      }}
                    >
                      {userStats.total_active_days}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Time Spent
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f0f0f0',
                      }}
                    >
                      {formatTime(userStats.total_time_spent_minutes)}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Current Streak
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f472b6',
                      }}
                    >
                      {userStats.current_streak_days} days
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Longest Streak
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f472b6',
                      }}
                    >
                      {userStats.longest_streak_days} days
                    </span>
                  </div>
                </div>
              </div>

              {/* Records Section */}
              <div>
                <h3
                  style={{
                    fontSize: '0.875rem',
                    fontWeight: '600',
                    color: '#9ca3af',
                    marginBottom: '0.75rem',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}
                >
                  Records
                </h3>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '0.5rem',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Best Day
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#10b981',
                      }}
                    >
                      {new Date(
                        userStats.highest_daily_score_date
                      ).toLocaleDateString()}
                    </span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                      borderBottom: '1px solid #2e2e2e',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Worst Day
                    </span>
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
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '0.5rem 0',
                    }}
                  >
                    <span style={{ fontSize: '0.875rem', color: '#d1d5db' }}>
                      Best Streak
                    </span>
                    <span
                      style={{
                        fontSize: '0.875rem',
                        fontWeight: '600',
                        color: '#f472b6',
                      }}
                    >
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

          {/* Badges Card */}
          <div
            style={{
              backgroundColor: '#1a1a1a',
              border: '1px solid #2e2e2e',
              borderRadius: '0.75rem',
              padding: '1.5rem',
            }}
          >
            <h2
              style={{
                fontSize: '1.25rem',
                fontWeight: 'bold',
                marginBottom: '1.5rem',
                color: '#f0f0f0',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
              }}
            >
              <Award size={20} />
              Badges
            </h2>

            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(4, 1fr)',
                gap: '1rem',
              }}
            >
              {badges.map((badge) => {
                const IconComponent = badge.icon;
                return (
                  <div
                    key={badge.id}
                    style={{
                      position: 'relative',
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'center',
                      justifyContent: 'center',
                      width: '64px',
                      height: '64px',
                      backgroundColor: '#161616',
                      border: `3px solid ${badge.type === 'positive' ? '#f472b6' : '#ef4444'}`,
                      borderRadius: '50%',
                      cursor: 'pointer',
                      transition: 'all 0.2s ease',
                    }}
                    title={`${badge.title}: ${badge.description}`}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.transform = 'scale(1.05)';
                      e.currentTarget.style.backgroundColor = '#202020';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.transform = 'scale(1)';
                      e.currentTarget.style.backgroundColor = '#161616';
                    }}
                  >
                    <IconComponent
                      size={28}
                      style={{
                        color:
                          badge.type === 'positive' ? '#f472b6' : '#ef4444',
                      }}
                    />
                    <div
                      style={{
                        position: 'absolute',
                        top: '-6px',
                        right: '-6px',
                        backgroundColor:
                          badge.type === 'positive' ? '#f472b6' : '#ef4444',
                        color: '#0f0f0f',
                        borderRadius: '50%',
                        width: '24px',
                        height: '24px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: '0.75rem',
                        fontWeight: '700',
                        border: '2px solid #0f0f0f',
                      }}
                    >
                      {badge.count}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Chart Section */}
        <div
          style={{
            backgroundColor: '#1a1a1a',
            border: '1px solid #2e2e2e',
            borderRadius: '0.75rem',
            padding: '1.5rem',
          }}
        >
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: '1.5rem',
            }}
          >
            <h2
              style={{
                fontSize: '1.25rem',
                fontWeight: 'bold',
                color: '#f0f0f0',
                margin: 0,
              }}
            >
              Progress Visualization
            </h2>
            <div
              style={{
                display: 'flex',
                backgroundColor: '#161616',
                border: '1px solid #2e2e2e',
                borderRadius: '0.5rem',
                padding: '0.25rem',
              }}
            >
              <button
                onClick={() => setChartView('graph')}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor:
                    chartView === 'graph' ? '#f472b6' : 'transparent',
                  color: chartView === 'graph' ? '#0f0f0f' : '#f0f0f0',
                  border: 'none',
                  borderRadius: '0.375rem',
                  cursor: 'pointer',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                }}
              >
                Graph
              </button>
              <button
                onClick={() => setChartView('heatmap')}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor:
                    chartView === 'heatmap' ? '#f472b6' : 'transparent',
                  color: chartView === 'heatmap' ? '#0f0f0f' : '#f0f0f0',
                  border: 'none',
                  borderRadius: '0.375rem',
                  cursor: 'pointer',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                }}
              >
                Heatmap
              </button>
            </div>
          </div>

          {chartView === 'graph' ? (
            <ScoreTrendGraph />
          ) : (
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                gap: '0.5rem',
              }}
            >
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  marginBottom: '1rem',
                }}
              >
                <span style={{ fontSize: '0.875rem', color: '#9ca3af' }}>
                  Last 365 days
                </span>
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '0.5rem',
                    fontSize: '0.75rem',
                    color: '#9ca3af',
                  }}
                >
                  <span>Less</span>
                  <div style={{ display: 'flex', gap: '2px' }}>
                    {[0, 3, 6, 9, 12, 15].map((score) => (
                      <div
                        key={score}
                        style={{
                          width: '10px',
                          height: '10px',
                          backgroundColor: getHeatmapColor(score),
                          borderRadius: '2px',
                        }}
                      />
                    ))}
                  </div>
                  <span>More</span>
                </div>
              </div>

              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(53, 1fr)',
                  gap: '2px',
                  maxWidth: '100%',
                }}
              >
                {heatmapData.map((day, index) => (
                  <div
                    key={index}
                    style={{
                      width: '12px',
                      height: '12px',
                      backgroundColor: getHeatmapColor(day.score),
                      borderRadius: '2px',
                      cursor: 'pointer',
                    }}
                    title={`${day.date}: ${day.score} points`}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default OtherPage;
