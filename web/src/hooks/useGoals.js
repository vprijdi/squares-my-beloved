import { useContext } from 'react';
import { GoalContext } from '../context/GoalContext';

export function useGoals() {
  const context = useContext(GoalContext);
  if (!context) throw new Error('Missing GoalProvider');
  return context;
}
