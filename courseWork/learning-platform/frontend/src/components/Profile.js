import React, { useState, useEffect } from 'react';
import api from '../services/api';
import ProgressBar from './ProgressBar';

function Profile() {
  const [submissions, setSubmissions] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [userProgress, setUserProgress] = useState({ completed: 0, total: 0 });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [submissionsResponse, tasksResponse] = await Promise.all([
          api.get('/submissions'),
          api.get('/tasks')
        ]);

        console.log('Ответ от сервера:', submissionsResponse.data);
        setSubmissions(Array.isArray(submissionsResponse.data) ? submissionsResponse.data : []);
        setTasks(Array.isArray(tasksResponse.data) ? tasksResponse.data : []);

        // Вычисляем прогресс
        calculateProgress(submissionsResponse.data, tasksResponse.data);
      } catch (err) {
        console.error('Ошибка загрузки данных:', err);
        setError('Ошибка загрузки данных: ' + (err.response?.data?.error || err.message));
        setSubmissions([]);
        setTasks([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const calculateProgress = (userSubmissions, allTasks) => {
    if (!userSubmissions || !allTasks) return;

    // Находим уникальные успешно решенные задания
    const successfulSubmissions = userSubmissions.filter(sub => sub.status === 'success');
    const solvedTaskIds = new Set(successfulSubmissions.map(sub => sub.task_id));
    
    setUserProgress({
      completed: solvedTaskIds.size,
      total: allTasks.length
    });
  };

  if (loading) return (
    <div style={loadingStyle}>
      <div style={spinnerStyle}></div>
      <p>Загрузка данных...</p>
    </div>
  );

  if (error) return (
    <div style={errorStyle}>
      <div style={errorIconStyle}>⚠️</div>
      <p>{error}</p>
    </div>
  );

  return (
    <div style={containerStyle}>
      <div style={headerStyle}>
        <h2 style={titleStyle}>Ваш профиль</h2>
        <ProgressBar completed={userProgress.completed} total={userProgress.total} />
      </div>

      <div style={contentStyle}>
        <div style={sectionStyle}>
          <h3 style={sectionTitleStyle}>Статистика</h3>
          <div style={statsGridStyle}>
            <div style={statCardStyle}>
              <div style={statNumberStyle}>{userProgress.completed}</div>
              <div style={statLabelStyle}>Решено заданий</div>
            </div>
            <div style={statCardStyle}>
              <div style={statNumberStyle}>{submissions.length}</div>
              <div style={statLabelStyle}>Всего попыток</div>
            </div>
            <div style={statCardStyle}>
              <div style={statNumberStyle}>
                {submissions.length > 0 ? Math.round((userProgress.completed / submissions.length) * 100) : 0}%
              </div>
              <div style={statLabelStyle}>Эффективность</div>
            </div>
          </div>
        </div>

        <div style={sectionStyle}>
          <h3 style={sectionTitleStyle}>История решений</h3>
          
          {!submissions || submissions.length === 0 ? (
            <div style={emptyStateStyle}>
              <div style={emptyIconStyle}>📝</div>
              <p>Вы еще не отправляли решения.</p>
              <p>Перейдите на страницу заданий, чтобы начать решать!</p>
            </div>
          ) : (
            <div style={submissionsListStyle}>
              {submissions.map(submission => (
                <div key={submission.id} style={submissionItemStyle}>
                  <div style={submissionHeaderStyle}>
                    <h4 style={taskTitleStyle}>Задание #{submission.task_id}</h4>
                    <span style={statusStyle(submission.status)}>
                      {submission.status === 'success' ? '✅ Успешно' : 
                       submission.status === 'error' ? '❌ Ошибка' : '⏳ В обработке'}
                    </span>
                  </div>
                  
                  {submission.created_at && (
                    <div style={timeStyle}>
                      {new Date(submission.created_at).toLocaleString('ru-RU')}
                    </div>
                  )}

                  {submission.output && (
                    <div style={outputSectionStyle}>
                      <h5 style={outputTitleStyle}>Результат:</h5>
                      <pre style={outputStyle}>
                        {submission.output}
                      </pre>
                    </div>
                  )}

                  {submission.test_results && submission.test_results.length > 0 && (
                    <div style={testsSectionStyle}>
                      <h5 style={testsTitleStyle}>Результаты тестов:</h5>
                      {submission.test_results.map((test, index) => (
                        <div key={index} style={testItemStyle(test.passed)}>
                          <div style={testHeaderStyle}>
                            <span style={testStatusStyle(test.passed)}>
                              {test.passed ? '✅' : '❌'} Тест {index + 1}
                            </span>
                            <span style={testDescriptionStyle}>{test.description}</span>
                          </div>
                          {!test.passed && (
                            <div style={testDetailsStyle}>
                              {test.input && <div>Вход: <strong>{test.input}</strong></div>}
                              <div>Ожидалось: <strong>{test.expected}</strong></div>
                              <div>Получено: <strong>{test.actual}</strong></div>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}

                  <details style={codeSectionStyle}>
                    <summary style={codeSummaryStyle}>Показать код решения</summary>
                    <pre style={codeStyle}>
                      {submission.code}
                    </pre>
                  </details>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// Стили
const containerStyle = {
  padding: '2rem',
  maxWidth: '1200px',
  margin: '0 auto'
};

const headerStyle = {
  marginBottom: '2rem'
};

const titleStyle = {
  color: '#2c3e50',
  marginBottom: '1.5rem',
  fontSize: '2rem'
};

const contentStyle = {
  display: 'flex',
  flexDirection: 'column',
  gap: '2rem'
};

const sectionStyle = {
  backgroundColor: '#fff',
  padding: '1.5rem',
  borderRadius: '12px',
  boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
};

const sectionTitleStyle = {
  color: '#2c3e50',
  marginBottom: '1rem',
  fontSize: '1.3rem'
};

const statsGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
  gap: '1rem',
  marginBottom: '1rem'
};

const statCardStyle = {
  backgroundColor: '#f8f9fa',
  padding: '1.5rem',
  borderRadius: '8px',
  textAlign: 'center',
  border: '2px solid #e9ecef'
};

const statNumberStyle = {
  fontSize: '2rem',
  fontWeight: 'bold',
  color: '#3498db',
  marginBottom: '0.5rem'
};

const statLabelStyle = {
  color: '#6c757d',
  fontSize: '0.9rem'
};

const submissionsListStyle = {
  display: 'flex',
  flexDirection: 'column',
  gap: '1.5rem'
};

const submissionItemStyle = {
  backgroundColor: '#f8f9fa',
  padding: '1.5rem',
  borderRadius: '8px',
  border: '1px solid #e9ecef'
};

const submissionHeaderStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginBottom: '0.5rem'
};

const taskTitleStyle = {
  margin: 0,
  color: '#2c3e50'
};

const timeStyle = {
  color: '#6c757d',
  fontSize: '0.9rem',
  marginBottom: '1rem'
};

const outputSectionStyle = {
  marginBottom: '1rem'
};

const outputTitleStyle = {
  margin: '0 0 0.5rem 0',
  color: '#2c3e50',
  fontSize: '1rem'
};

const outputStyle = {
  backgroundColor: '#fff',
  padding: '1rem',
  borderRadius: '6px',
  border: '1px solid #dee2e6',
  fontSize: '0.9rem',
  overflow: 'auto',
  maxHeight: '200px',
  margin: 0
};

const testsSectionStyle = {
  marginBottom: '1rem'
};

const testsTitleStyle = {
  margin: '0 0 0.5rem 0',
  color: '#2c3e50',
  fontSize: '1rem'
};

const testItemStyle = (passed) => ({
  backgroundColor: passed ? '#d4edda' : '#f8d7da',
  padding: '1rem',
  borderRadius: '6px',
  marginBottom: '0.5rem',
  border: `1px solid ${passed ? '#c3e6cb' : '#f5c6cb'}`
});

const testHeaderStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginBottom: '0.5rem'
};

const testStatusStyle = (passed) => ({
  fontWeight: 'bold',
  color: passed ? '#155724' : '#721c24'
});

const testDescriptionStyle = {
  color: '#6c757d',
  fontSize: '0.9rem'
};

const testDetailsStyle = {
  fontSize: '0.9rem',
  color: '#495057'
};

const codeSectionStyle = {
  marginTop: '1rem'
};

const codeSummaryStyle = {
  cursor: 'pointer',
  color: '#3498db',
  fontWeight: '500',
  fontSize: '0.9rem'
};

const codeStyle = {
  backgroundColor: '#2d3748',
  color: '#e2e8f0',
  padding: '1rem',
  borderRadius: '6px',
  fontSize: '0.8rem',
  overflow: 'auto',
  maxHeight: '300px',
  marginTop: '0.5rem'
};

const emptyStateStyle = {
  textAlign: 'center',
  padding: '3rem',
  color: '#6c757d'
};

const emptyIconStyle = {
  fontSize: '3rem',
  marginBottom: '1rem'
};

const loadingStyle = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '3rem',
  color: '#6c757d'
};

const spinnerStyle = {
  border: '4px solid #f3f3f3',
  borderTop: '4px solid #3498db',
  borderRadius: '50%',
  width: '40px',
  height: '40px',
  animation: 'spin 1s linear infinite',
  marginBottom: '1rem'
};

const errorStyle = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '3rem',
  color: '#dc3545',
  textAlign: 'center'
};

const errorIconStyle = {
  fontSize: '3rem',
  marginBottom: '1rem'
};

// Добавляем анимацию для спиннера
const style = document.createElement('style');
style.textContent = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;
document.head.appendChild(style);

const statusStyle = (status) => ({
  padding: '0.4rem 0.8rem',
  borderRadius: '20px',
  fontSize: '0.8rem',
  fontWeight: '600',
  backgroundColor: status === 'success' ? '#d4edda' : 
                  status === 'error' ? '#f8d7da' : '#fff3cd',
  color: status === 'success' ? '#155724' : 
         status === 'error' ? '#721c24' : '#856404',
  border: `1px solid ${status === 'success' ? '#c3e6cb' : 
                       status === 'error' ? '#f5c6cb' : '#ffeaa7'}`
});

export default Profile;