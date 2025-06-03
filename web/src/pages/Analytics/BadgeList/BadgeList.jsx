import Badge from '../Badge/Badge';
import styles from './BadgeList.module.css';
import {
  ArrowLeft,
  Trophy,
  Target,
  Calendar,
  Clock,
  TrendingUp,
  TrendingDown,
  Flame,
  Award,
  Star,
  Zap,
  Shield,
  Crown,
  Activity,
  BarChart3,
} from 'lucide-react';

export default function BadgeList() {
  const badges = [
    {
      id: 1,
      icon: Flame,
      count: 3,
      title: 'Week Warrior',
      description: 'Maintain positive score for 7 days',
      type: 'positive',
    },
    {
      id: 2,
      icon: Trophy,
      count: 1,
      title: 'Century Club',
      description: 'Reach total score of 100',
      type: 'positive',
    },
    {
      id: 3,
      icon: Target,
      count: 2,
      title: 'Task Master',
      description: 'Complete 50 tasks',
      type: 'positive',
    },
    {
      id: 4,
      icon: Crown,
      count: 1,
      title: 'Streak King',
      description: 'Maintain 30-day streak',
      type: 'positive',
    },
    {
      id: 5,
      icon: Zap,
      count: 1,
      title: 'Lightning Round',
      description: 'Complete 10 tasks in one day',
      type: 'positive',
    },
    {
      id: 6,
      icon: Shield,
      count: 2,
      title: 'Consistency',
      description: 'Active for 30 days',
      type: 'positive',
    },
    {
      id: 7,
      icon: TrendingDown,
      count: 1,
      title: 'Rock Bottom',
      description: 'Reach total score of -50',
      type: 'negative',
    },
    {
      id: 8,
      icon: Activity,
      count: 1,
      title: 'Rough Week',
      description: 'Negative score for 7 days',
      type: 'negative',
    },
  ];

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>
        {' '}
        <Award size={20} />
        Badges
      </h3>

      <div className={styles.badgeContainer}>
        {badges.map((badge) => (
          <Badge key={badge.id} badge={badge} />
        ))}
      </div>
    </div>
  );
}
