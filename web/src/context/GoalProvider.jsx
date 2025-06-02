import { useState, useEffect } from 'react';
import { GoalContext } from './GoalContext';

export function GoalProvider({ children }) {
  const [goals, setGoals] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/goals.json')
      .then((res) => res.json())
      .then((data) => {
        setGoals(data.goals);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to load goals:', err);
        setLoading(false);
      });
  }, []);

  const newGoal = (newGoal) => {
    setGoals((prevGoals) => [
      ...prevGoals,
      {
        id: Date.now().toString(),
        title: newGoal.title,
        completion: 0,
        subgoals: newGoal.subgoals || [],
      },
    ]);
  };

  const value = {
    goals,
    loading,
    newGoal,
  };

  return <GoalContext.Provider value={value}>{children}</GoalContext.Provider>;
}
