// Если хотите оставить файл, но убрать ошибку
import { useEffect } from 'react';
import { useAuth } from './useAuth';
import { authAPI } from '../services/api';
import api from '../services/api';

export const useAuthCheck = () => {
  const { user, isLoading } = useAuth();

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      const savedUser = localStorage.getItem('user');
      
      if (token && !user && !isLoading) {
        try {
          // Проверяем валидность токена
          const userData = await authAPI.getCurrentUser();
          // Обновляем состояние через существующие методы
          // (этот хук может не понадобиться, если AuthProvider уже обрабатывает это)
        } catch (error) {
          console.error('Auth check failed:', error);
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
      }
    };

    checkAuth();
  }, [user, isLoading]);
};