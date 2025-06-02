import styles from './Home.module.css';
import Header from '../../components/Header/Header';
import TaskList from './TaskList/TaskList';
import { BarChart3, Target } from 'lucide-react';

export default function Home() {
  return (
    <div>
      <Header />
      <div className={styles.container}>
        <div className={styles.leftPanel}>
          <div className={styles.leftHeader}>
            <h2 className={styles.title}>Today's Tasks</h2>
            <div className={styles.tabs}>
              <button className={styles.tabButton}>
                <BarChart3 size={16} />
                Analytics
              </button>
              <button className={styles.tabButton}>
                <Target size={16} />
                Goals
              </button>
            </div>
          </div>

          <TaskList />
        </div>

        <div className={styles.rightPanel}>
          <div className={styles.goalListContainer}>{/* <GoalList /> */}</div>
          <div className={styles.weeklyWidgetContainer}>
            {/* <WeeklyWidget /> */}
          </div>
        </div>
      </div>
    </div>
  );
}
