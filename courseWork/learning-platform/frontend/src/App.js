import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Header from './components/Header';
import TaskList from './components/TaskList';
import Task from './components/Task';
import Profile from './components/Profile';
import Login from './components/Login';
import Register from './components/Register';
import './App.css';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = () => {
    const token = localStorage.getItem('token');
    const userId = localStorage.getItem('user_id');
    const username = localStorage.getItem('username');
    
    if (token && userId && username) {
      setIsAuthenticated(true);
      setUser({ id: userId, username });
    } else {
      setIsAuthenticated(false);
      setUser(null);
    }
    setLoading(false);
  };

  const handleLogin = (userData) => {
    localStorage.setItem('token', userData.token);
    localStorage.setItem('user_id', userData.user_id);
    localStorage.setItem('username', userData.username);
    
    setIsAuthenticated(true);
    setUser({ id: userData.user_id, username: userData.username });
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('username');
    setIsAuthenticated(false);
    setUser(null);
  };

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        backgroundColor: '#f8f9fa'
      }}>
        <div style={spinnerStyle}></div>
      </div>
    );
  }

  return (
    <Router>
      <div className="App">
        <Header user={user} onLogout={handleLogout} />
        <div style={{ 
          padding: '2rem', 
          maxWidth: '1200px', 
          margin: '0 auto',
          minHeight: 'calc(100vh - 80px)',
          backgroundColor: '#f8f9fa'
        }}>
          <Routes>
            <Route path="/" element={isAuthenticated ? <TaskList /> : <Navigate to="/login" />} />
            <Route path="/task/:id" element={isAuthenticated ? <Task /> : <Navigate to="/login" />} />
            <Route 
              path="/profile" 
              element={isAuthenticated ? <Profile /> : <Navigate to="/login" />} 
            />
            <Route 
              path="/login" 
              element={isAuthenticated ? <Navigate to="/" /> : <Login onLogin={handleLogin} />} 
            />
            <Route 
              path="/register" 
              element={isAuthenticated ? <Navigate to="/" /> : <Register />} 
            />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

const spinnerStyle = {
  border: '4px solid #f3f3f3',
  borderTop: '4px solid #3498db',
  borderRadius: '50%',
  width: '40px',
  height: '40px',
  animation: 'spin 1s linear infinite'
};

// Добавляем стили для анимации
const style = document.createElement('style');
style.textContent = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;
document.head.appendChild(style);

export default App;