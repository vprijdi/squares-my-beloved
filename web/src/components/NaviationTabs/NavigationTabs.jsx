import styles from './NagivationTabs.module.css';
import { BarChart3, Target } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function NavigationTabs() {
  const navigate = useNavigate();

  return (
    <div className={styles.leftHeader}>
      <h2 className={styles.title}>Today's Tasks</h2>
      <div className={styles.tabs}>
        <label className={styles.label}>Score: 137</label>
        <button
          className={styles.tabButton}
          onClick={() => navigate('/analytics')}
        >
          <BarChart3 size={16} />
          Analytics
        </button>
        <button className={styles.tabButton} onClick={() => navigate('/goals')}>
          <Target size={16} />
          Goals
        </button>
      </div>
    </div>
  );
}
