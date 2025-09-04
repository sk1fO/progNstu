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
        console.log('Ответ от сервера:', response.data);
        setSubmissions(Array.isArray(response.data) ? response.data : []);
      } catch (err) {
        console.error('Ошибка загрузки отправленных решений:', err);
        setError('Ошибка загрузки отправленных решений: ' + (err.response?.data?.error || err.message));
        setSubmissions([]);
      } finally {
        setLoading(false);
      }
    };

    fetchSubmissions();
  }, []);

  if (loading) return <div style={{ padding: '2rem', textAlign: 'center' }}>Загрузка решений...</div>;
  if (error) return <div style={{ color: 'red', padding: '1rem', backgroundColor: '#ffe6e6', borderRadius: '4px', margin: '1rem' }}>{error}</div>;

  return (
    <div style={{ padding: '1rem' }}>
      <h2>Ваши отправленные решения</h2>
      
      {!submissions || submissions.length === 0 ? (
        <div style={{ padding: '2rem', textAlign: 'center', color: '#666' }}>
          <p>Вы еще не отправляли решения.</p>
          <p>Перейдите на страницу заданий, чтобы начать решать!</p>
        </div>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {submissions.map(submission => (
            <li key={submission.id} style={submissionItemStyle}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                <h3 style={{ margin: 0 }}>Задание #{submission.task_id}</h3>
                <span style={statusStyle(submission.status)}>
                  {submission.status === 'success' ? 'Успешно' : 
                   submission.status === 'error' ? 'Ошибка' : 'В обработке'}
                </span>
              </div>
              
              {submission.created_at && (
                <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '1rem' }}>
                  Отправлено: {new Date(submission.created_at).toLocaleString('ru-RU')}
                </div>
              )}

              {submission.output && (
                <div>
                  <h4 style={{ margin: '0 0 0.5rem 0' }}>Результат:</h4>
                  <pre style={{ 
                    backgroundColor: '#f8f9fa', 
                    padding: '1rem', 
                    borderRadius: '4px',
                    whiteSpace: 'pre-wrap',
                    fontSize: '0.875rem',
                    maxHeight: '200px',
                    overflow: 'auto',
                    border: '1px solid #e9ecef',
                    margin: 0
                  }}>
                    {submission.output}
                  </pre>
                </div>
              )}

              <details style={{ marginTop: '1rem' }}>
                <summary style={{ cursor: 'pointer', color: '#007bff', fontSize: '0.875rem' }}>
                  Показать код
                </summary>
                <pre style={{
                  backgroundColor: '#f8f9fa',
                  padding: '1rem',
                  borderRadius: '4px',
                  whiteSpace: 'pre-wrap',
                  fontSize: '0.75rem',
                  maxHeight: '300px',
                  overflow: 'auto',
                  border: '1px solid #e9ecef',
                  marginTop: '0.5rem'
                }}>
                  {submission.code}
                </pre>
              </details>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

const submissionItemStyle = {
  marginBottom: '1.5rem',
  padding: '1.5rem',
  border: '1px solid #ddd',
  borderRadius: '8px',
  backgroundColor: '#fff',
  boxShadow: '0 1px 3px rgba(0,0,0,0.1)'
};

const statusStyle = (status) => {
  const colors = {
    success: '#28a745',
    error: '#dc3545',
    pending: '#ffc107'
  };
  
  return {
    padding: '0.25rem 0.75rem',
    borderRadius: '4px',
    backgroundColor: colors[status] || '#6c757d',
    color: 'white',
    fontSize: '0.875rem',
    fontWeight: 'bold'
  };
};

export default Profile;