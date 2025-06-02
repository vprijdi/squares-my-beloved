import { useTasks } from '../../../context/useTasks';
import styles from './TaskCard.module.css';
import { CheckSquare, Square } from 'lucide-react';

export default function TaskCard({ task }) {
  const { toggleTaskCheckmark } = useTasks();
  // const { incrementCompletion } = useTasks();
  // const isDone = task.completions_today >= task.repetitions;

  const getTierClass = () => {
    switch (task.tier) {
      case 2:
        return styles.long;
      case 1:
        return styles.medium;
      case 0:
      default:
        return styles.quick;
    }
  };
  const getTierLabel = () => {
    switch (task.tier) {
      case 2:
        return 'Large';
      case 1:
        return 'Medium';
      case 0:
      default:
        return 'Quick';
    }
  };

  const getDuration = () => {
    const remainingReps =
      task.repetitions - task.completed.filter((c) => c).length;

    let baseDuration;
    switch (task.tier) {
      case 2:
        baseDuration = 60;
        break;
      case 1:
        baseDuration = 30;
        break;
      case 0:
        baseDuration = 10;
        break;
      default:
        return '0m';
    }

    const totalMinutes = baseDuration * remainingReps;

    if (totalMinutes >= 60) {
      console.log(totalMinutes);
      const hours = Math.floor(totalMinutes / 60);
      const mins = totalMinutes % 60;
      return `${hours}h${mins > 0 ? ` ${mins}m` : ''}`;
    }
    return totalMinutes === 10 ? `<10m` : `${totalMinutes}m`;
  };

  function getStatusColor() {
    const isComplete =
      task.completed.filter(Boolean).length === task.repetitions;

    if (isComplete) return '#10b981';

    switch (task.tier) {
      case 2:
        return '#f472b6';
      case 1:
        return '#c084fc';
      case 0:
      default:
        return '#60a5fa';
    }
  }
  const cardClasses = `${styles.card} ${getTierClass()}`;

  return (
    <div className={cardClasses}>
      <div style={{ flex: 1 }}>
        <div
          style={{
            textDecoration: task.completed.every((c) => c)
              ? 'line-through'
              : 'none',
            opacity: task.completed.every((c) => c) ? 0.7 : 1,
            fontWeight: '500',
          }}
        >
          {task.title}
        </div>

        <div className={styles.taskMeta}>
          <span className={styles.duration}>{getTierLabel()}</span>
          <span className={styles.grey}>•</span>
          <span className={styles.grey}>{getDuration()}</span>
          <span className={styles.grey}>•</span>
          <span
            style={{
              fontSize: '0.75rem',
              color: getStatusColor(),
            }}
          >
            {task.completed.filter((c) => c).length}/{task.repetitions} complete
          </span>
        </div>
      </div>

      <div className={styles.something}>
        {task.completed.map((isChecked, index) => (
          <button
            key={index}
            onClick={(e) => {
              e.stopPropagation();
              toggleTaskCheckmark(task.id, index);
            }}
            className={styles.checkbox}
          >
            {isChecked ? (
              <CheckSquare className={styles.checkboxIcon} />
            ) : (
              <Square className={styles.checkboxIcon} />
            )}
          </button>
        ))}
      </div>
    </div>
  );
}
