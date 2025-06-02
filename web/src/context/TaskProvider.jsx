import { useState, useEffect } from 'react';
import { TaskContext } from './TaskContext';

export function TaskProvider({ children }) {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/tasks.json')
      .then((res) => res.json())
      .then((data) => {
        setTasks(data.tasks);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to load tasks:', err);
        setLoading(false);
      });
  }, []);

  const newTask = (newTask) => {
    setTasks((prevTasks) => [
      ...prevTasks,
      {
        id: Date.now().toString(),
        title: newTask.title,
        tier: newTask.tier || 2,
        repetitions: newTask.repetitions || 1,
        completed: Array(newTask.repetitions || 1).fill(false),
        completions_today: 0,
        completion_count: 0,
        target_repetitions: newTask.repetitions || 1,
      },
    ]);
  };

  const incrementCompletion = (taskId, checkboxIndex) => {
    setTasks(
      tasks.map((task) => {
        if (task.id !== taskId) return task;

        const newCompletionsToday = Math.max(
          checkboxIndex + 1,
          task.completions_today
        );

        return {
          ...task,
          completions_today: newCompletionsToday,
          completion_count: task.completion_count + 1,
          last_updated: new Date().toISOString(),
          is_complete: newCompletionsToday >= task.target_repetitions,
        };
      })
    );

    setTasks(
      tasks.map((task) =>
        task.id === taskId
          ? {
              ...task,
              completions_today: task.completions_today + 1,
              completion_count: task.completion_count + 1,
              last_updated: new Date().toISOString().split('T')[0], // YYYY-MM-DD
            }
          : task
      )
    );
  };

  const toggleTaskCheckmark = (taskId, checkmarkIndex) => {
    setTasks(
      tasks.map((task) =>
        task.id === taskId
          ? {
              ...task,
              completed: task.completed.map((checked, index) =>
                index === checkmarkIndex ? !checked : checked
              ),
            }
          : task
      )
    );
  };

  const value = {
    tasks,
    loading,
    incrementCompletion,
    toggleTaskCheckmark,
    newTask,
  };

  useEffect(() => {
    const today = new Date().toDateString();
    const lastUpdated = localStorage.getItem('tasksLastUpdated');

    if (lastUpdated !== today) {
      setTasks(
        tasks.map((task) => ({
          ...task,
          completions_today: 0,
        }))
      );
      localStorage.setItem('tasksLastUpdated', today);
    }
  }, [tasks]);

  return <TaskContext.Provider value={value}>{children}</TaskContext.Provider>;
}
