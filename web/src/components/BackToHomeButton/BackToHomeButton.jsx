import styles from './BackToHomeButton.module.css';
import { ChevronLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function BackToHomeButton() {
  const navigate = useNavigate();

  return (
    <button onClick={() => navigate('/')} className={styles.button}>
      <ChevronLeft className={styles.arrow} />
    </button>
  );
}
