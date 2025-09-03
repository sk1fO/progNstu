import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { coursesAPI, authAPI } from '../../services/api';
import { Course, User } from '../../types';
import './Dashboard.css';

const Dashboard: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [userResponse, coursesResponse] = await Promise.all([
          authAPI.getCurrentUser(),
          coursesAPI.getCourses()
        ]);
        
        setUser(userResponse);
        
        // Проверяем, что coursesResponse.data является массивом
        if (Array.isArray(coursesResponse.data)) {
          setCourses(coursesResponse.data);
        } else {
          console.error('Ожидался массив курсов, но получено:', coursesResponse.data);
          setCourses([]);
          setError('Некорректный формат данных курсов');
        }
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
        setError('Ошибка загрузки данных');
        setCourses([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) return <div className="loading">Загрузка...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div style={{ padding: '20px' }}>
      <h1>Добро пожаловать, {user?.username}!</h1>
      <p>Ваша роль: {user?.role === 'teacher' ? 'Преподаватель' : 'Студент'}</p>
      
      <div style={{ marginTop: '30px' }}>
        <h2>Доступные курсы</h2>
        {courses.length === 0 ? (
          <p>Нет доступных курсов</p>
        ) : (
          <div style={{ display: 'grid', gap: '20px', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
            {courses.map(course => (
              <div key={course.id} style={{ 
                border: '1px solid #ddd', 
                padding: '20px', 
                borderRadius: '8px' 
              }}>
                <h3>{course.title}</h3>
                <p>{course.description}</p>
                <Link to={`/assignments/${course.id}`}>
                  <button style={{ 
                    padding: '10px 20px', 
                    backgroundColor: '#007bff', 
                    color: 'white', 
                    border: 'none', 
                    borderRadius: '4px' 
                  }}>
                    Перейти к занятиям
                  </button>
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;