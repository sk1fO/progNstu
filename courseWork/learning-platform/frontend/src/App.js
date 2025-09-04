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

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userId = localStorage.getItem('user_id');
    const username = localStorage.getItem('username');
    
    if (token && userId && username) {
      setIsAuthenticated(true);
      setUser({ id: userId, username });
    }
  }, []);

  const handleLogin = (userData) => {
    setIsAuthenticated(true);
    setUser(userData);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('username');
    setIsAuthenticated(false);
    setUser(null);
  };

  return (
    <Router>
      <div className="App">
        <Header user={user} onLogout={handleLogout} />
        <div style={{ padding: '1rem', maxWidth: '1200px', margin: '0 auto' }}>
          <Routes>
            <Route path="/" element={<TaskList />} />
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

export default App;