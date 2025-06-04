import styles from './Analytics.module.css';
import BackToHomeButton from '../../components/BackToHomeButton/BackToHomeButton.jsx';
import Header from '../../components/Header/Header.jsx';
import StatisticsCard from './StatisticsCard/StatisticsCard.jsx';
import BadgeList from './BadgeList/BadgeList.jsx';
import ScoreView from './ScoreView/ScoreView.jsx';

export default function Analytics() {
  return (
    <div>
      <Header />
      <div className={styles.container}>
        <div>
          <BackToHomeButton />
        </div>

        <div>
          <StatisticsCard />
          <BadgeList />
          <ScoreView />
        </div>
      </div>
    </div>
  );
}
