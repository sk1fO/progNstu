import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { AuthContextType, User, RegisterData, LoginData, AuthResponse } from '../../types';
import { authAPI } from '../../services/api';
import api from '../../services/api';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Проверяем аутентификацию при загрузке приложения
  useEffect(() => {
    const initAuth = async () => {
      const savedToken = localStorage.getItem('token');
      const savedUser = localStorage.getItem('user');

      if (savedToken) {
        try {
          setToken(savedToken);
          
          // Устанавливаем токен для API
          api.defaults.headers.common['Authorization'] = `Bearer ${savedToken}`;
          
          // Проверяем валидность токена
          const userData = await authAPI.getCurrentUser();
          setUser(userData);
          
          // Сохраняем пользователя в localStorage
          localStorage.setItem('user', JSON.stringify(userData));
        } catch (error) {
          console.error('Invalid token:', error);
          logout();
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = async (username: string, password: string): Promise<AuthResponse> => {
    try {
      const response = await authAPI.login({ username, password });
      
      setToken(response.token);
      setUser(response.user);
      
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      return response;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed');
    }
  };

  const register = async (userData: RegisterData): Promise<AuthResponse> => {
    try {
      const response = await authAPI.register(userData);
      
      setToken(response.token);
      setUser(response.user);
      
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      return response;
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Registration failed');
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    authAPI.logout();
  };

  const value: AuthContextType = {
    user,
    token,
    login,
    register,
    logout,
    isLoading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};