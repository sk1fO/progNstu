import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../services/api';

function Register() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    try {
      await api.post('/register', {
        username,
        password
      });
      
      setSuccess(true);
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка регистрации');
    }
  };

  if (success) {
    return (
      <div style={{ maxWidth: '400px', margin: '2rem auto', textAlign: 'center' }}>
        <h2>Регистрация успешна!</h2>
        <p>Вы будете перенаправлены на страницу входа.</p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '400px', margin: '2rem auto' }}>
      <h2>Регистрация</h2>
      
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

        <div style={{ marginBottom: '1rem' }}>
          <input
            type="password"
            placeholder="Подтвердите пароль"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            style={{ width: '100%', padding: '0.5rem' }}
            required
          />
        </div>
        
        <button 
          type="submit"
          style={{ width: '100%', padding: '0.5rem', backgroundColor: '#007bff', color: 'white', border: 'none' }}
        >
          Зарегистрироваться
        </button>
      </form>
      
      <p style={{ marginTop: '1rem' }}>
        Уже есть аккаунт? <Link to="/login">Войдите</Link>
      </p>
    </div>
  );
}

export default Register;