import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import './Navbar.css';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/dashboard">
          <h2>Learning Platform</h2>
        </Link>
      </div>
      
      <div className="navbar-menu">
        <Link to="/dashboard" className="navbar-item">
          Главная
        </Link>
        <Link to="/courses" className="navbar-item">
          Курсы
        </Link>
      </div>
      
      <div className="navbar-user">
        <span className="user-info">
          {user?.username} ({user?.role === 'teacher' ? 'Преподаватель' : 'Студент'})
        </span>
        <button onClick={handleLogout} className="logout-button">
          Выйти
        </button>
      </div>
    </nav>
  );
};

export default Navbar;