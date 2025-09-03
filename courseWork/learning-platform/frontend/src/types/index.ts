export interface User {
  id: number;
  username: string;
  email: string;
  role: 'student' | 'teacher';
}

export interface Course {
  id: number;
  title: string;
  description: string;
  teacher_id: number;
  created_at: string;
}

export interface Lesson {
  id: number;
  course_id: number;
  title: string;
  order: number;
  theory_content: string;
}

export interface Assignment {
  id: number;
  lesson_id: number;
  title: string;
  description: string;
  starter_code: string;
  language: string;
}

export interface Solution {
  id: number;
  assignment_id: number;
  user_id: number;
  code: string;
  status: 'sent' | 'tested' | 'checked';
  passed_autotests: boolean;
  teacher_comment: string;
  score: number;
  submitted_at: string;
}

export interface TestResult {
  input: string;
  expected: string;
  actual: string;
  passed: boolean;
  error_message?: string;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<AuthResponse>;
  register: (userData: RegisterData) => Promise<AuthResponse>;
  logout: () => void;
  isLoading: boolean;
}

export interface RegisterData {
  username: string;
  email: string;
  password: string;
  role: 'student' | 'teacher';
}

export interface LoginData {
  username: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}