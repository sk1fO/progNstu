import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';

function TaskList() {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await api.get('/tasks');
        setTasks(response.data);
      } catch (err) {
        setError('Ошибка загрузки заданий');
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, []);

  if (loading) return <div>Загрузка...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

  return (
    <div>
      <h2>Список заданий</h2>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {tasks.map(task => (
          <li key={task.id} style={taskItemStyle}>
            <Link to={`/task/${task.id}`} style={taskLinkStyle}>
              <h3>{task.title}</h3>
              <p>{task.description}</p>
              <span style={difficultyStyle(task.difficulty)}>
                Сложность: {task.difficulty}
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

const taskItemStyle = {
  marginBottom: '1rem',
  padding: '1rem',
  border: '1px solid #ddd',
  borderRadius: '4px'
};

const taskLinkStyle = {
  textDecoration: 'none',
  color: 'inherit'
};

const difficultyStyle = (difficulty) => {
  const colors = {
    easy: '#28a745',
    medium: '#ffc107',
    hard: '#dc3545'
  };
  
  return {
    display: 'inline-block',
    padding: '0.25rem 0.5rem',
    borderRadius: '4px',
    backgroundColor: colors[difficulty] || '#6c757d',
    color: 'white',
    fontSize: '0.875rem'
  };
};

export default TaskList;