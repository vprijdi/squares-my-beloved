import styles from './Home.module.css';
import Header from '../../components/Header/Header';
import TaskList from './TaskList/TaskList';
import { BarChart3, Target } from 'lucide-react';
import WeeklyWidget from './WeeklyWidget/WeeklyWidget';
import GoalList from './GoalList/GoalList';
import NavigationTabs from '../../components/NaviationTabs/NavigationTabs';

export default function Home() {
  return (
    <div>
      <Header />
      <div className={styles.container}>
        <div className={styles.leftPanel}>
          <NavigationTabs />
          <TaskList />
        </div>

        <div className={styles.rightPanel}>
          <div className={styles.weeklyWidgetContainer}>
            <WeeklyWidget />
          </div>
          <div className={styles.goalListContainer}>
            <GoalList />
          </div>
        </div>
      </div>
    </div>
  );
}
