import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    
    try {
      const response = await api.post('/login', {
        username,
        password
      });
      
      localStorage.setItem('token', response.data.token);
      onLogin();
      navigate('/');
    } catch (err) {
      setError('Ошибка входа. Проверьте имя пользователя и пароль.');
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '2rem auto' }}>
      <h2>Вход</h2>
      
      {error && <div style={{ color: 'red', marginBottom: '1rem' }}>{error}</div>}
      
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '1rem' }}>
          <input
            type="text"
            placeholder="Имя пользователя"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={{ width: '100%', padding: '0.5rem' }}
            required
          />
        </div>
        
        <div style={{ marginBottom: '1rem' }}>
          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ width: '100%', padding: '0.5rem' }}
            required
          />
        </div>
        
        <button 
          type="submit"
          style={{ width: '100%', padding: '0.5rem', backgroundColor: '#007bff', color: 'white', border: 'none' }}
        >
          Войти
        </button>
      </form>
      
      <p style={{ marginTop: '1rem' }}>
        Нет аккаунта? <a href="/register">Зарегистрируйтесь</a>
      </p>
    </div>
  );
}

export default Login;