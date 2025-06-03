import styles from './Badge.module.css';

export default function Badge({ badge }) {
  const IconComponent = badge.icon;
  return (
    <div
      className={styles.badge}
      style={{
        border: `3px solid ${badge.type === 'positive' ? '#f472b6' : '#ef4444'}`,
      }}
      title={`${badge.title}: ${badge.description}`}
    >
      <IconComponent
        size={22}
        style={{
          color: badge.type === 'positive' ? '#f472b6' : '#ef4444',
        }}
      />

      <div
        className={styles.count}
        style={{
          backgroundColor: badge.type === 'positive' ? '#f472b6' : '#ef4444',
        }}
      >
        {badge.count}
      </div>
    </div>
  );
}
