import React from 'react';

function ProgressBar({ completed, total }) {
  const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;

  return (
    <div style={containerStyle}>
      <div style={progressInfoStyle}>
        <span style={textStyle}>Прогресс: {completed}/{total} заданий</span>
        <span style={percentageStyle}>{percentage}%</span>
      </div>
      
      <div style={progressBarContainerStyle}>
        <div 
          style={{
            ...progressBarStyle,
            width: `${percentage}%`,
            backgroundColor: percentage === 100 ? '#27ae60' : '#3498db'
          }}
        />
      </div>
      
      <div style={statusStyle}>
        {percentage === 100 ? (
          <span style={successStyle}>🎉 Все задания выполнены!</span>
        ) : percentage >= 70 ? (
          <span style={goodStyle}>Хороший прогресс!</span>
        ) : percentage >= 30 ? (
          <span style={mediumStyle}>Продолжайте в том же духе!</span>
        ) : (
          <span style={startStyle}>Начните с первого задания!</span>
        )}
      </div>
    </div>
  );
}

const containerStyle = {
  backgroundColor: '#fff',
  padding: '1.5rem',
  borderRadius: '12px',
  boxShadow: '0 4px 15px rgba(0,0,0,0.1)',
  marginBottom: '2rem'
};

const progressInfoStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginBottom: '0.5rem'
};

const textStyle = {
  fontSize: '1.1rem',
  fontWeight: '600',
  color: '#2c3e50'
};

const percentageStyle = {
  fontSize: '1.2rem',
  fontWeight: 'bold',
  color: '#3498db'
};

const progressBarContainerStyle = {
  width: '100%',
  height: '12px',
  backgroundColor: '#ecf0f1',
  borderRadius: '6px',
  overflow: 'hidden',
  marginBottom: '0.5rem'
};

const progressBarStyle = {
  height: '100%',
  borderRadius: '6px',
  transition: 'width 0.5s ease, background-color 0.5s ease'
};

const statusStyle = {
  textAlign: 'center',
  fontSize: '0.9rem'
};

const successStyle = {
  color: '#27ae60',
  fontWeight: '600'
};

const goodStyle = {
  color: '#f39c12',
  fontWeight: '600'
};

const mediumStyle = {
  color: '#3498db',
  fontWeight: '600'
};

const startStyle = {
  color: '#7f8c8d',
  fontWeight: '600'
};

export default ProgressBar;