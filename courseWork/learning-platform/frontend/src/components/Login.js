import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    
    try {
      const response = await api.post('/login', {
        username,
        password
      });
      
      // Вызываем колбэк с данными пользователя
      onLogin({
        token: response.data.token,
        user_id: response.data.user_id,
        username: response.data.username
      });
      
      // Перенаправляем на главную страницу
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка входа. Проверьте имя пользователя и пароль.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '2rem auto' }}>
      <div style={cardStyle}>
        <h2 style={titleStyle}>Вход в систему</h2>
        
        {error && (
          <div style={errorStyle}>
            <div style={errorIconStyle}>⚠️</div>
            <p>{error}</p>
          </div>
        )}
        
        <form onSubmit={handleSubmit} style={formStyle}>
          <div style={inputGroupStyle}>
            <label style={labelStyle}>
              Имя пользователя
            </label>
            <input
              type="text"
              placeholder="Введите имя пользователя"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={inputStyle}
              required
              disabled={loading}
            />
          </div>
          
          <div style={inputGroupStyle}>
            <label style={labelStyle}>
              Пароль
            </label>
            <input
              type="password"
              placeholder="Введите пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              style={inputStyle}
              required
              disabled={loading}
            />
          </div>
          
          <button 
            type="submit"
            disabled={loading}
            style={buttonStyle(loading)}
          >
            {loading ? 'Вход...' : 'Войти'}
          </button>
        </form>
        
        <div style={registerLinkStyle}>
          <p>Нет аккаунта? <a href="/register" style={linkStyle}>Зарегистрируйтесь</a></p>
        </div>
      </div>
    </div>
  );
}

const cardStyle = {
  backgroundColor: '#fff',
  padding: '2rem',
  borderRadius: '12px',
  boxShadow: '0 4px 15px rgba(0,0,0,0.1)'
};

const titleStyle = {
  textAlign: 'center',
  color: '#2c3e50',
  marginBottom: '1.5rem',
  fontSize: '1.8rem'
};

const formStyle = {
  marginBottom: '1.5rem'
};

const inputGroupStyle = {
  marginBottom: '1.5rem'
};

const labelStyle = {
  display: 'block',
  marginBottom: '0.5rem',
  fontWeight: '600',
  color: '#2c3e50'
};

const inputStyle = {
  width: '100%',
  padding: '0.75rem',
  border: '2px solid #e9ecef',
  borderRadius: '6px',
  fontSize: '1rem',
  transition: 'border-color 0.3s'
};

const buttonStyle = (loading) => ({
  width: '100%',
  padding: '0.75rem',
  backgroundColor: loading ? '#6c757d' : '#3498db',
  color: 'white',
  border: 'none',
  borderRadius: '6px',
  fontSize: '1rem',
  fontWeight: '600',
  cursor: loading ? 'not-allowed' : 'pointer',
  transition: 'background-color 0.3s'
});

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

const registerLinkStyle = {
  textAlign: 'center',
  color: '#6c757d'
};

const linkStyle = {
  color: '#3498db',
  textDecoration: 'none',
  fontWeight: '600'
};

// Добавляем hover эффекты
Object.assign(inputStyle, {
  ':focus': {
    borderColor: '#3498db',
    outline: 'none'
  }
});

Object.assign(buttonStyle(false), {
  ':hover': {
    backgroundColor: '#2980b9'
  }
});

Object.assign(linkStyle, {
  ':hover': {
    textDecoration: 'underline'
  }
});

export default Login;