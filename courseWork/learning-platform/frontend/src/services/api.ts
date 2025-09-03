import axios from 'axios';
import { RegisterData, LoginData, AuthResponse, User, Course } from '../types';

const API_BASE_URL = 'http://localhost:8000/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Интерцептор для добавления токена к запросам
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Интерцептор для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Не перенаправляем на страницу входа если мы уже на странице аутентификации
      if (!window.location.pathname.includes('/login') && 
          !window.location.pathname.includes('/register')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: async (userData: RegisterData): Promise<AuthResponse> => {
    const response = await api.post('/register', userData);
    
    // Сохраняем токен для будущих запросов
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      api.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`;
    }
    
    return response.data;
  },

  login: async (loginData: LoginData): Promise<AuthResponse> => {
    const response = await api.post('/login', loginData);
    
    // Сохраняем токен для будущих запросов
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      api.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`;
    }
    
    return response.data;
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/user');
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    delete api.defaults.headers.common['Authorization'];
  }
};

export const coursesAPI = {
  getCourses: () => api.get('/courses'),
  createCourse: (courseData: Omit<Course, 'id' | 'created_at' | 'teacher_id'>) => 
    api.post('/courses', courseData),
  updateCourse: (courseId: number, courseData: Partial<Course>) => 
    api.put(`/courses/${courseId}`, courseData),
  deleteCourse: (courseId: number) => api.delete(`/courses/${courseId}`),
  getLessons: (courseId: number) => api.get(`/courses/${courseId}/lessons`),
};

export const assignmentsAPI = {
  getAssignments: (lessonId: number) => api.get(`/lessons/${lessonId}/assignments`),
  runCode: (assignmentId: number, code: string) =>
    api.post(`/assignments/${assignmentId}/run`, { code }),
};

export default api;