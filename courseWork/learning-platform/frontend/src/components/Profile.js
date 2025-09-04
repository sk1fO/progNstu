import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Profile() {
  const [submissions, setSubmissions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchSubmissions = async () => {
      try {
        const response = await api.get('/submissions');
        setSubmissions(response.data);
      } catch (err) {
        setError('Ошибка загрузки отправленных решений');
      } finally {
        setLoading(false);
      }
    };

    fetchSubmissions();
  }, []);

  if (loading) return <div>Загрузка...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

  return (
    <div>
      <h2>Ваши отправленные решения</h2>
      
      {submissions.length === 0 ? (
        <p>Вы еще не отправляли решения.</p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {submissions.map(submission => (
            <li key={submission.id} style={submissionItemStyle}>
              <h3>Задание #{submission.task_id}</h3>
              <p>Статус: 
                <span style={statusStyle(submission.status)}>
                  {submission.status === 'success' ? 'Успешно' : 
                   submission.status === 'error' ? 'Ошибка' : 'В обработке'}
                </span>
              </p>
              {submission.output && (
                <div>
                  <h4>Вывод:</h4>
                  <pre style={{ 
                    backgroundColor: '#f8f9fa', 
                    padding: '0.5rem', 
                    borderRadius: '4px',
                    whiteSpace: 'pre-wrap',
                    fontSize: '0.875rem'
                  }}>
                    {submission.output}
                  </pre>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

const submissionItemStyle = {
  marginBottom: '1rem',
  padding: '1rem',
  border: '1px solid #ddd',
  borderRadius: '4px'
};

const statusStyle = (status) => {
  const colors = {
    success: '#28a745',
    error: '#dc3545',
    pending: '#ffc107'
  };
  
  return {
    display: 'inline-block',
    marginLeft: '0.5rem',
    padding: '0.25rem 0.5rem',
    borderRadius: '4px',
    backgroundColor: colors[status] || '#6c757d',
    color: 'white',
    fontSize: '0.875rem'
  };
};

export default Profile;