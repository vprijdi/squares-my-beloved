import styles from './GoalCard.module.css';
import { useState } from 'react';
import { ChevronDown, ChevronRight, CheckCircle, Circle } from 'lucide-react';

export default function GoalCard({ goal }) {
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div className={styles.card}>
      <div
        onClick={() => setIsExpanded(!isExpanded)}
        className={styles.something}
      >
        {isExpanded ? (
          <ChevronDown className={styles.chevron} />
        ) : (
          <ChevronRight sclassName={styles.chevron} />
        )}

        <div style={{ flex: 1 }}>
          <div className={styles.title}> {goal.title}</div>
          <div className={styles.meta}>
            <div className={styles.bar}></div>
            <span className={styles.percentage}>{goal.completion}%</span>
          </div>

          <div className={styles.label}>
            {goal.subgoals.filter((sg) => sg.completed).length}/
            {goal.subgoals.length} subgoals
          </div>
        </div>
      </div>

      {isExpanded && (
        <div
          style={{
            borderTop: '1px solid #2e2e2e',
            backgroundColor: '#161616',
          }}
        >
          {goal.subgoals.map((subgoal) => (
            <div
              key={subgoal.id}
              style={{
                padding: '0.75rem 1.25rem',
                display: 'flex',
                alignItems: 'center',
                gap: '0.75rem',
                borderBottom: '1px solid #1f1f1f',
              }}
            >
              {subgoal.completed ? (
                <CheckCircle
                  size={16}
                  style={{ color: '#f472b6', flexShrink: 0 }}
                />
              ) : (
                <Circle size={16} style={{ color: '#6b7280', flexShrink: 0 }} />
              )}
              <span
                style={{
                  textDecoration: subgoal.completed ? 'line-through' : 'none',
                  opacity: subgoal.completed ? 0.7 : 1,
                  fontSize: '0.875rem',
                }}
              >
                {subgoal.title}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
