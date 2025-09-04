import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../services/api';

function Task() {
  const { id } = useParams();
  const [task, setTask] = useState(null);
  const [code, setCode] = useState('');
  const [output, setOutput] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTask = async () => {
      try {
        const response = await api.get(`/tasks/${id}`);
        setTask(response.data);
      } catch (err) {
        setError('Ошибка загрузки задания');
      }
    };

    fetchTask();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      const response = await api.post('/submit', {
        task_id: parseInt(id),
        code
      });
      
      setOutput(response.data.output);
    } catch (err) {
      setError('Ошибка отправки решения');
    } finally {
      setLoading(false);
    }
  };

  if (!task) return <div>Загрузка...</div>;

  return (
    <div>
      <h2>{task.title}</h2>
      <p>{task.description}</p>
      
      <form onSubmit={handleSubmit}>
        <div>
          <textarea
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="Введите ваш код на C++"
            rows={10}
            style={{ width: '100%', fontFamily: 'monospace' }}
            required
          />
        </div>
        
        <button 
          type="submit" 
          disabled={loading}
          style={{ marginTop: '1rem', padding: '0.5rem 1rem' }}
        >
          {loading ? 'Отправка...' : 'Отправить решение'}
        </button>
      </form>
      
      {error && <div style={{ color: 'red', marginTop: '1rem' }}>{error}</div>}
      
      {output && (
        <div style={{ marginTop: '2rem' }}>
          <h3>Результат:</h3>
          <pre style={{ 
            backgroundColor: '#f8f9fa', 
            padding: '1rem', 
            borderRadius: '4px',
            whiteSpace: 'pre-wrap'
          }}>
            {output}
          </pre>
        </div>
      )}
    </div>
  );
}

export default Task;