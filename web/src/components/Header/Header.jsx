import styles from './Header.module.css';
import { format } from 'date-fns';

import { Box, Calendar, Blocks } from 'lucide-react';

export default function Header() {
  const today = new Date();
  const formattedDate = format(today, 'MMMM d, yyyy');

  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <div className={styles.headerLeft}>
          <Blocks className={styles.headerIcon} />
          <h1 className={styles.headerTitle}>Learning Tracker</h1>
          <div className={styles.headerDate}>
            <Calendar className={styles.dateIcon} />
            <span>Today - {formattedDate}</span>
          </div>
        </div>

        <div className={styles.headerStats}>
          <div className={styles.statItem}>
            <div className={styles.statContent}>
              <div className={styles.statLabel}>Level 10</div>
              <div className={styles.statSublabel}>800 / 1000 XP</div>
            </div>
          </div>

          <div className={styles.xpBar}>
            <div
              className={styles.xpProgress}
              style={{
                width: `80%`,
              }}
            ></div>
          </div>
        </div>
      </div>
    </header>
  );
}
