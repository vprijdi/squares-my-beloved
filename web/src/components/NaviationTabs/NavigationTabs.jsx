import styles from './NavigationTabs.module.css';

export default function NavigationTabs() {
  return (
    <div className={styles.leftHeader}>
      <h2 className={styles.title}>Today's Tasks</h2>
      <div className={styles.tabs}>
        <label className={styles.label}>Score: 137</label>
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
  );
}
