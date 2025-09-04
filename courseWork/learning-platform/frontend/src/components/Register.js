import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../services/api';

function Register() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    if (password.length < 6) {
      setError('Пароль должен содержать не менее 6 символов');
      return;
    }

    setLoading(true);

    try {
      const response = await api.post('/register', {
        username,
        password
      });
      
      setSuccess(true);
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка регистрации. Возможно, пользователь уже существует.');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div style={{ maxWidth: '400px', margin: '2rem auto' }}>
        <div style={cardStyle}>
          <div style={successStyle}>
            <div style={successIconStyle}>✅</div>
            <h2 style={successTitleStyle}>Регистрация успешна!</h2>
            <p style={successTextStyle}>Вы будете перенаправлены на страницу входа.</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '400px', margin: '2rem auto' }}>
      <div style={cardStyle}>
        <h2 style={titleStyle}>Регистрация</h2>
        
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
              placeholder="Придумайте имя пользователя"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={inputStyle}
              required
              disabled={loading}
              minLength={3}
              maxLength={20}
            />
            <div style={hintStyle}>От 3 до 20 символов</div>
          </div>
          
          <div style={inputGroupStyle}>
            <label style={labelStyle}>
              Пароль
            </label>
            <input
              type="password"
              placeholder="Придумайте пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              style={inputStyle}
              required
              disabled={loading}
              minLength={6}
            />
            <div style={hintStyle}>Не менее 6 символов</div>
          </div>

          <div style={inputGroupStyle}>
            <label style={labelStyle}>
              Подтверждение пароля
            </label>
            <input
              type="password"
              placeholder="Повторите пароль"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
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
            {loading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>
        </form>
        
        <div style={loginLinkStyle}>
          <p>Уже есть аккаунт? <Link to="/login" style={linkStyle}>Войдите</Link></p>
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

const hintStyle = {
  fontSize: '0.8rem',
  color: '#6c757d',
  marginTop: '0.25rem'
};

const buttonStyle = (loading) => ({
  width: '100%',
  padding: '0.75rem',
  backgroundColor: loading ? '#6c757d' : '#27ae60',
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

const loginLinkStyle = {
  textAlign: 'center',
  color: '#6c757d'
};

const linkStyle = {
  color: '#3498db',
  textDecoration: 'none',
  fontWeight: '600'
};

const successStyle = {
  textAlign: 'center',
  padding: '2rem'
};

const successIconStyle = {
  fontSize: '3rem',
  marginBottom: '1rem'
};

const successTitleStyle = {
  color: '#27ae60',
  marginBottom: '0.5rem'
};

const successTextStyle = {
  color: '#6c757d'
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
    backgroundColor: '#219a52'
  }
});

Object.assign(linkStyle, {
  ':hover': {
    textDecoration: 'underline'
  }
});

export default Register;