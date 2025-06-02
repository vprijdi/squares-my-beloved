import styles from './WeeklyWidget.module.css';

export default function WeeklyWidget() {
  const weeklyPoints = [
    { day: 'Mon', points: 6 },
    { day: 'Tue', points: -2 },
    { day: 'Wed', points: 8 },
    { day: 'Thu', points: 4 },
    { day: 'Fri', points: -1 },
    { day: 'Sat', points: 5 },
    { day: 'Sun', points: 0 },
  ];

  return (
    <div>
      <h3 className={styles.title}>Weekly Progress</h3>

      <div className={styles.container}>
        <div className={styles.chart}>
          {weeklyPoints.map(({ day, points }) => (
            <div key={day} className={styles.barContainer}>
              <div
                style={{
                  height:
                    points === 0 ? '2px' : `${Math.abs(points) * 8 + 20}px`,
                  width: '24px',
                  backgroundColor:
                    points === 0
                      ? '#f472b6'
                      : points > 0
                        ? '#f472b6'
                        : '#d4d4d8',
                  borderRadius: '0',
                  display: 'flex',
                  alignItems: points > 0 ? 'flex-start' : 'flex-end',
                  justifyContent: 'center',
                  marginBottom: points < 0 ? `${Math.abs(points) * 8}px` : '0',
                  marginTop: points > 0 ? `${(8 - points) * 8}px` : '0',
                }}
              >
                <span
                  style={{
                    fontSize: '0.75rem',
                    fontWeight: 'bold',
                    color: '#0f0f0f',
                    padding: '2px 0',
                  }}
                >
                  {points > 0 ? `+${points}` : points}
                </span>
              </div>
              <span
                style={{
                  fontSize: '0.75rem',
                  color: '#9ca3af',
                  fontWeight: '500',
                }}
              >
                {day}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
