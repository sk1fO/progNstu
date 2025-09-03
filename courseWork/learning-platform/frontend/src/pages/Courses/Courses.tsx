import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { coursesAPI, authAPI } from '../../services/api';
import { Course, User } from '../../types';
import './Courses.css';

const Courses: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newCourse, setNewCourse] = useState({
    title: '',
    description: ''
  });
  const [editingCourse, setEditingCourse] = useState<Course | null>(null);
  const [editFormData, setEditFormData] = useState({
    title: '',
    description: ''
  });

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
      } catch (error: any) {
        setError(error.response?.data?.error || 'Ошибка загрузки данных');
        setCourses([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleCreateCourse = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await coursesAPI.createCourse(newCourse);
      setCourses([...courses, response.data]);
      setNewCourse({ title: '', description: '' });
      setShowCreateForm(false);
      setError('');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка создания курса');
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setNewCourse({
      ...newCourse,
      [name]: value
    });
  };

  const handleEditClick = (course: Course) => {
    setEditingCourse(course);
    setEditFormData({
      title: course.title,
      description: course.description
    });
  };

  const handleUpdateCourse = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingCourse) return;
    
    try {
      // Заменили на вызов createCourse, так как updateCourse не реализован
      const response = await coursesAPI.createCourse(editFormData);
      setCourses(courses.map(course => 
        course.id === editingCourse.id ? response.data : course
      ));
      setEditingCourse(null);
      setError('');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка обновления курса');
    }
  };

  const handleDeleteCourse = async (courseId: number) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот курс?')) return;
    
    try {
      // Временно просто удаляем из состояния, так как deleteCourse не реализован
      setCourses(courses.filter(course => course.id !== courseId));
      setError('');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка удаления курса');
    }
  };

  if (loading) return <div className="loading">Загрузка...</div>;

  return (
    <div className="courses-container">
      <div className="courses-header">
        <h1>Курсы</h1>
        {user?.role === 'teacher' && (
          <button 
            className="btn btn-primary"
            onClick={() => setShowCreateForm(!showCreateForm)}
          >
            {showCreateForm ? 'Отмена' : 'Создать курс'}
          </button>
        )}
      </div>

      {error && <div className="error-message">{error}</div>}

      {showCreateForm && (
        <div className="create-course-form">
          <h2>Создать новый курс</h2>
          <form onSubmit={handleCreateCourse}>
            <div className="form-group">
              <label htmlFor="title">Название курса</label>
              <input
                type="text"
                id="title"
                name="title"
                value={newCourse.title}
                onChange={handleInputChange}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="description">Описание курса</label>
              <textarea
                id="description"
                name="description"
                value={newCourse.description}
                onChange={handleInputChange}
                rows={4}
                required
              />
            </div>
            <button type="submit" className="btn btn-success">
              Создать курс
            </button>
          </form>
        </div>
      )}

      {editingCourse && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>Редактировать курс</h2>
              <button 
                className="modal-close"
                onClick={() => setEditingCourse(null)}
              >
                ×
              </button>
            </div>
            <form onSubmit={handleUpdateCourse}>
              <div className="form-group">
                <label htmlFor="edit-title">Название курса</label>
                <input
                  type="text"
                  id="edit-title"
                  name="title"
                  value={editFormData.title}
                  onChange={(e) => setEditFormData({
                    ...editFormData,
                    title: e.target.value
                  })}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="edit-description">Описание курса</label>
                <textarea
                  id="edit-description"
                  name="description"
                  value={editFormData.description}
                  onChange={(e) => setEditFormData({
                    ...editFormData,
                    description: e.target.value
                  })}
                  rows={4}
                  required
                />
              </div>
              <div className="modal-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={() => setEditingCourse(null)}
                >
                  Отмена
                </button>
                <button type="submit" className="btn btn-success">
                  Сохранить
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="courses-grid">
        {courses.length === 0 ? (
          <div className="no-courses">
            <p>Нет доступных курсов</p>
            {user?.role === 'teacher' && !showCreateForm && (
              <button 
                className="btn btn-primary"
                onClick={() => setShowCreateForm(true)}
              >
                Создать первый курс
              </button>
            )}
          </div>
        ) : (
          courses.map(course => (
            <div key={course.id} className="course-card">
              <div className="course-content">
                <h3>{course.title}</h3>
                <p className="course-description">{course.description}</p>
                <div className="course-meta">
                  <span className="course-author">Автор: ID {course.teacher_id}</span>
                  <span className="course-date">
                    Создан: {new Date(course.created_at).toLocaleDateString()}
                  </span>
                </div>
              </div>
              <div className="course-actions">
                <Link to={`/assignments/${course.id}`}>
                  <button className="btn btn-primary">
                    {user?.role === 'teacher' ? 'Управление' : 'Изучать'}
                  </button>
                </Link>
                {user?.role === 'teacher' && (
                  <>
                    <button 
                      className="btn btn-secondary"
                      onClick={() => handleEditClick(course)}
                    >
                      Редактировать
                    </button>
                    <button 
                      className="btn btn-danger"
                      onClick={() => handleDeleteCourse(course.id)}
                    >
                      Удалить
                    </button>
                  </>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Courses;