import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

function Header({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    navigate('/login');
  };

  return (
    <header style={headerStyle}>
      <div style={logoStyle}>
        <h1 style={titleStyle}>C++ Learner</h1>
        <span style={subtitleStyle}>Платформа для изучения C++</span>
      </div>
      
      <nav style={navStyle}>
        {user ? (
          <>
            <Link to="/" style={linkStyle} onMouseOver={(e) => e.target.style.backgroundColor = '#34495e'} onMouseOut={(e) => e.target.style.backgroundColor = 'transparent'}>
              Задания
            </Link>
            <Link to="/profile" style={linkStyle} onMouseOver={(e) => e.target.style.backgroundColor = '#34495e'} onMouseOut={(e) => e.target.style.backgroundColor = 'transparent'}>
              Профиль
            </Link>
            <div style={userInfoStyle}>
              <span style={usernameStyle}>👤 {user.username}</span>
              <button 
                onClick={handleLogout} 
                style={logoutButtonStyle}
                onMouseOver={(e) => e.target.style.backgroundColor = '#c0392b'} 
                onMouseOut={(e) => e.target.style.backgroundColor = '#e74c3c'}
              >
                Выйти
              </button>
            </div>
          </>
        ) : (
          <Link 
            to="/login" 
            style={loginButtonStyle}
            onMouseOver={(e) => e.target.style.backgroundColor = '#2980b9'} 
            onMouseOut={(e) => e.target.style.backgroundColor = '#3498db'}
          >
            Войти
          </Link>
        )}
      </nav>
    </header>
  );
}

const headerStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  padding: '1rem 2rem',
  backgroundColor: '#2c3e50',
  color: 'white',
  boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
  position: 'sticky',
  top: 0,
  zIndex: 1000
};

const logoStyle = {
  display: 'flex',
  flexDirection: 'column'
};

const titleStyle = {
  margin: 0,
  fontSize: '1.8rem',
  fontWeight: 'bold',
  color: '#3498db'
};

const subtitleStyle = {
  fontSize: '0.9rem',
  color: '#bdc3c7',
  marginTop: '0.2rem'
};

const navStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '2rem'
};

const linkStyle = {
  color: 'white',
  textDecoration: 'none',
  padding: '0.5rem 1rem',
  borderRadius: '4px',
  transition: 'background-color 0.3s ease',
  fontWeight: '500'
};

const userInfoStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '1rem'
};

const usernameStyle = {
  color: '#ecf0f1',
  fontWeight: '500',
  fontSize: '0.9rem'
};

const logoutButtonStyle = {
  background: '#e74c3c',
  color: 'white',
  border: 'none',
  padding: '0.5rem 1rem',
  borderRadius: '4px',
  cursor: 'pointer',
  fontWeight: '500',
  transition: 'background-color 0.3s ease'
};

const loginButtonStyle = {
  background: '#3498db',
  color: 'white',
  textDecoration: 'none',
  padding: '0.75rem 1.5rem',
  borderRadius: '4px',
  fontWeight: '500',
  transition: 'background-color 0.3s ease'
};

export default Header;