import styles from './TaskList.module.css';
import { useState } from 'react';
import { useTasks } from '../../../hooks/useTasks';
import TaskCard from '../TaskCard/TaskCard';
import { Plus } from 'lucide-react';

export default function TaskList() {
  const { tasks, loading, newTask } = useTasks();
  const [showNewTask, setShowNewTask] = useState(false);
  const [newTaskLength, setNewTaskLength] = useState('Quick');
  const [newTaskTitle, setNewTaskTitle] = useState('');

  if (loading) return <div>Loading tasks...</div>;

  const getTierFromLength = (length) => {
    switch (length) {
      case 'Quick':
        return 0;
      case 'Medium':
        return 1;
      case 'Large':
        return 2;
      default:
        return 0;
    }
  };

  const addTask = () => {
    if (newTaskTitle.trim()) {
      const task = {
        id: Date.now(),
        title: newTaskTitle,
        tier: getTierFromLength(newTaskLength),
        completed: [],
        repetitions: 1,
      };
      newTask(task);
      setNewTaskTitle('');
      setShowNewTask(false);
    }
  };

  return (
    <div className={styles.container}>
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} />
      ))}

      {showNewTask ? (
        <div className={styles.newTaskContainer}>
          <div style={{ marginBottom: '1rem' }}>
            <label className={styles.label}>Task Length</label>
            <div className={styles.lengthChoice}>
              {['Quick', 'Medium', 'Large'].map((length) => (
                <button
                  className={styles.lengthButton}
                  key={length}
                  onClick={() => setNewTaskLength(length)}
                  style={{
                    backgroundColor:
                      newTaskLength === length
                        ? length === 'Quick'
                          ? '#60a5fa'
                          : length === 'Medium'
                            ? '#c084fc'
                            : '#f472b6'
                        : '#2e2e2e',
                  }}
                >
                  {length}
                </button>
              ))}
            </div>
          </div>

          <input
            className={styles.newTaskInput}
            type="text"
            placeholder="Enter task title..."
            value={newTaskTitle}
            onChange={(e) => setNewTaskTitle(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && newTaskTitle.trim()) {
                addTask();
              }
            }}
            autoFocus
          />
          <div
            style={{
              display: 'flex',
              gap: '0.5rem',
              justifyContent: 'flex-end',
            }}
          >
            <button
              onClick={() => setShowNewTask(false)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#2e2e2e',
                border: 'none',
                borderRadius: '0.375rem',
                color: '#f0f0f0',
                cursor: 'pointer',
                fontSize: '0.875rem',
              }}
            >
              Cancel
            </button>
            <button
              onClick={addTask}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#f472b6',
                border: 'none',
                borderRadius: '0.375rem',
                color: '#0f0f0f',
                cursor: 'pointer',
                fontSize: '0.875rem',
                fontWeight: '600',
              }}
            >
              Add Task
            </button>
          </div>
        </div>
      ) : (
        <button className={styles.newTask} onClick={() => setShowNewTask(true)}>
          <Plus size={16} />
          Add New Task
        </button>
      )}
    </div>
  );
}
