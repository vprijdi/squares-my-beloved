import styles from './Analytics.module.css';
import BackToHomeButton from '../../components/BackToHomeButton/BackToHomeButton.jsx';
import Header from '../../components/Header/Header.jsx';
import StatisticsCard from './StatisticsCard/StatisticsCard.jsx';
import BadgeList from './BadgeList/BadgeList.jsx';

export default function Analytics() {
  return (
    <div>
      <Header></Header>
      <div className={styles.container}>
        <div>
          <BackToHomeButton />
        </div>

        <div>
          <StatisticsCard />
          <BadgeList />
        </div>
      </div>
    </div>
  );
}
