import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';
import api from '../services/api';
import MonacoEditor from './MonacoEditor';

function Task() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [task, setTask] = useState(null);
  const [code, setCode] = useState('');
  const [output, setOutput] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchTask = async () => {
      try {
        const response = await api.get(`/tasks/${id}`);
        setTask(response.data);
        setCode(response.data.template || `#include <iostream>\n\nint main() {\n    // Ваш код здесь\n    return 0;\n}`);
      } catch (err) {
        setError('Ошибка загрузки задания');
        console.error('Ошибка загрузки задания:', err);
      }
    };

    fetchTask();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!code.trim()) {
      setError('Введите код решения');
      return;
    }

    setIsSubmitting(true);
    setError('');
    setOutput('');
    
    try {
      const response = await api.post('/submit', {
        task_id: parseInt(id),
        code
      });
      
      setOutput(response.data.output);
      
      
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка отправки решения');
      console.error('Ошибка отправки:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCodeChange = (newCode) => {
    setCode(newCode);
    setError('');
  };

  if (!task) {
    return (
      <div style={loadingStyle}>
        <div style={spinnerStyle}></div>
        <p>Загрузка задания...</p>
      </div>
    );
  }

  return (
    <div style={containerStyle}>
      <div style={headerStyle}>
        <button 
          onClick={() => navigate('/')}
          style={backButtonStyle}
        >
          ← Назад к заданиям
        </button>
        <h1 style={titleStyle}>{task.title}</h1>
        <span style={difficultyStyle(task.difficulty)}>
          Сложность: {task.difficulty}
        </span>
      </div>

      <div style={contentStyle}>
        <div style={descriptionSectionStyle}>
          <div style={cardStyle}>
            <h2 style={sectionTitleStyle}>Описание задания</h2>
            <p style={descriptionStyle}>{task.description}</p>
            
            {task.template && (
              <div style={templateSectionStyle}>
                <h3 style={subTitleStyle}>Шаблон кода</h3>
                <div style={templateStyle}>
                  <SyntaxHighlighter
                    language="cpp"
                    style={atomDark}
                    customStyle={{
                      borderRadius: '6px',
                      padding: '1rem',
                      fontSize: '12px'
                    }}
                  >
                    {task.template}
                  </SyntaxHighlighter>
                </div>
              </div>
            )}
          </div>
        </div>

        <div style={solutionSectionStyle}>
          <div style={cardStyle}>
            <h2 style={sectionTitleStyle}>Ваше решение</h2>
            
            <form onSubmit={handleSubmit}>
              <MonacoEditor
                code={code}
                onChange={handleCodeChange}
                readOnly={false}
              />
              
              <div style={buttonGroupStyle}>
                <button 
                  type="submit" 
                  disabled={isSubmitting}
                  style={submitButtonStyle(isSubmitting)}
                >
                  {isSubmitting ? (
                    <>
                      <div style={buttonSpinnerStyle}></div>
                      Проверка...
                    </>
                  ) : (
                    'Отправить решение'
                  )}
                </button>
                
                <button 
                  type="button"
                  onClick={() => setCode(task.template || '')}
                  style={resetButtonStyle}
                  disabled={isSubmitting}
                >
                  Сбросить
                </button>
              </div>
            </form>

            {error && (
              <div style={errorStyle}>
                <div style={errorIconStyle}>⚠️</div>
                <p>{error}</p>
              </div>
            )}

            {output && (
              <div style={outputSectionStyle}>
                <h3 style={outputTitleStyle}>
                  Результат проверки
                  {output.includes('✅ Все тесты пройдены') && (
                    <span style={successIconStyle}> </span>
                  )}
                </h3>
                <div style={outputStyle}>
                  <pre style={outputTextStyle}>{output}</pre>
                </div>
                
                {output.includes('✅ Все тесты пройдены') && (
                  <div style={successMessageStyle}>
                    <div style={successCheckStyle}>✓</div>
                    <p>Поздравляем! Задание выполнено успешно!</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

// Стили остаются такими же как в предыдущей версии...
// Стили
const containerStyle = {
  padding: '2rem',
  maxWidth: '1400px',
  margin: '0 auto',
  minHeight: '100vh'
};

const headerStyle = {
  marginBottom: '2rem',
  display: 'flex',
  alignItems: 'center',
  gap: '1rem',
  flexWrap: 'wrap'
};

const backButtonStyle = {
  background: 'none',
  border: '2px solid #3498db',
  color: '#3498db',
  padding: '0.5rem 1rem',
  borderRadius: '6px',
  cursor: 'pointer',
  fontWeight: '500',
  transition: 'all 0.3s ease'
};

const titleStyle = {
  color: '#2c3e50',
  margin: 0,
  fontSize: '2.2rem',
  flex: 1
};

const difficultyStyle = (difficulty) => {
  const colors = {
    easy: '#28a745',
    medium: '#ffc107',
    hard: '#dc3545'
  };
  
  return {
    padding: '0.5rem 1rem',
    borderRadius: '20px',
    backgroundColor: colors[difficulty] || '#6c757d',
    color: 'white',
    fontWeight: '600',
    fontSize: '0.9rem'
  };
};

const contentStyle = {
  display: 'grid',
  gridTemplateColumns: '1fr 1.5fr',
  gap: '2rem',
  alignItems: 'start'
};

const descriptionSectionStyle = {
  height: 'fit-content'
};

const solutionSectionStyle = {
  height: 'fit-content'
};

const cardStyle = {
  backgroundColor: '#fff',
  padding: '2rem',
  borderRadius: '12px',
  boxShadow: '0 4px 15px rgba(0,0,0,0.1)',
  marginBottom: '1.5rem'
};

const sectionTitleStyle = {
  color: '#2c3e50',
  marginBottom: '1.5rem',
  fontSize: '1.5rem',
  borderBottom: '3px solid #3498db',
  paddingBottom: '0.5rem'
};

const descriptionStyle = {
  color: '#495057',
  lineHeight: '1.6',
  fontSize: '1.1rem',
  marginBottom: '1.5rem'
};

const templateSectionStyle = {
  marginTop: '1.5rem'
};

const subTitleStyle = {
  color: '#2c3e50',
  marginBottom: '1rem',
  fontSize: '1.2rem'
};

const templateStyle = {
  backgroundColor: '#2d3748',
  borderRadius: '8px',
  overflow: 'hidden'
};

const buttonGroupStyle = {
  display: 'flex',
  gap: '1rem',
  marginBottom: '1.5rem'
};

const submitButtonStyle = (isSubmitting) => ({
  padding: '1rem 2rem',
  backgroundColor: isSubmitting ? '#6c757d' : '#28a745',
  color: 'white',
  border: 'none',
  borderRadius: '6px',
  fontSize: '1rem',
  fontWeight: '600',
  cursor: isSubmitting ? 'not-allowed' : 'pointer',
  transition: 'background-color 0.3s ease',
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem'
});

const resetButtonStyle = {
  padding: '1rem 1.5rem',
  backgroundColor: '#6c757d',
  color: 'white',
  border: 'none',
  borderRadius: '6px',
  fontSize: '1rem',
  fontWeight: '600',
  cursor: 'pointer',
  transition: 'background-color 0.3s ease'
};

const errorStyle = {
  backgroundColor: '#f8d7da',
  color: '#721c24',
  padding: '1rem',
  borderRadius: '6px',
  marginBottom: '1.5rem',
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem'
};

const errorIconStyle = {
  fontSize: '1.2rem'
};

const outputSectionStyle = {
  marginTop: '2rem'
};

const outputTitleStyle = {
  color: '#2c3e50',
  marginBottom: '1rem',
  fontSize: '1.3rem',
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem'
};

const outputStyle = {
  backgroundColor: '#f8f9fa',
  border: '2px solid #e9ecef',
  borderRadius: '8px',
  padding: '1rem',
  marginBottom: '1rem'
};

const outputTextStyle = {
  margin: 0,
  whiteSpace: 'pre-wrap',
  fontFamily: '"Monaco", "Menlo", "Ubuntu Mono", monospace',
  fontSize: '14px',
  lineHeight: '1.4'
};

const successIconStyle = {
  fontSize: '1.5rem'
};

const successMessageStyle = {
  backgroundColor: '#d4edda',
  color: '#155724',
  padding: '1.5rem',
  borderRadius: '8px',
  textAlign: 'center',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '1rem'
};

const successCheckStyle = {
  width: '40px',
  height: '40px',
  backgroundColor: '#28a745',
  color: 'white',
  borderRadius: '50%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontSize: '1.5rem',
  fontWeight: 'bold'
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

const buttonSpinnerStyle = {
  border: '2px solid #f3f3f3',
  borderTop: '2px solid #fff',
  borderRadius: '50%',
  width: '16px',
  height: '16px',
  animation: 'spin 1s linear infinite'
};

// Добавляем анимации
const style = document.createElement('style');
style.textContent = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .back-button:hover {
    background-color: #3498db;
    color: white;
  }
  
  .submit-button:hover:not(:disabled) {
    background-color: #218838;
  }
  
  .reset-button:hover:not(:disabled) {
    background-color: #5a6268;
  }
`;
document.head.appendChild(style);

export default Task;