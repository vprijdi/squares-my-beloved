import ScoreTrendGraph from './ScoreTrendGraph/ScoreTrendGraph';
import styles from './ScoreView.module.css';

export default function ScoreView() {
  return (
    <div className={styles.container}>
      <div className={styles.buttonContainer}>
        <div></div>
      </div>
      <ScoreTrendGraph />
    </div>
  );
}
