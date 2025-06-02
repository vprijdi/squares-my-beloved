import styles from './GoalList.module.css';
import { useGoals } from '../../../hooks/useGoals';
import GoalCard from '../../../components/GoalCard/GoalCard';

export default function GoalList() {
  const { goals } = useGoals();

  return (
    <div>
      <h3 className={styles.header}>Learning Goals</h3>

      <div className={styles.container}>
        {goals.map((goal) => (
          <GoalCard key={goal.id} goal={goal} />
        ))}
      </div>
    </div>
  );
}
