import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

function Header() {
  const navigate = useNavigate();
  const token = localStorage.getItem('token');

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <header style={headerStyle}>
      <h1>Обучение C++</h1>
      <nav>
        <Link to="/" style={linkStyle}>Задания</Link>
        <Link to="/profile" style={linkStyle}>Профиль</Link>
        {token ? (
          <button onClick={handleLogout} style={buttonStyle}>Выйти</button>
        ) : (
          <>
            <Link to="/login" style={linkStyle}>Войти</Link>
            <Link to="/register" style={linkStyle}>Регистрация</Link>
          </>
        )}
      </nav>
    </header>
  );
}

const headerStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  padding: '1rem',
  backgroundColor: '#f8f9fa',
  borderBottom: '1px solid #dee2e6'
};

const linkStyle = {
  margin: '0 10px',
  textDecoration: 'none',
  color: '#007bff'
};

const buttonStyle = {
  margin: '0 10px',
  background: 'none',
  border: 'none',
  color: '#007bff',
  cursor: 'pointer'
};

export default Header;