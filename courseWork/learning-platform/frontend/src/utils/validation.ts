export const validateEmail = (email: string): string => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!email) return 'Email обязателен';
  if (!emailRegex.test(email)) return 'Некорректный email';
  return '';
};

export const validatePassword = (password: string): string => {
  if (!password) return 'Пароль обязателен';
  if (password.length < 6) return 'Пароль должен быть не менее 6 символов';
  return '';
};

export const validateUsername = (username: string): string => {
  if (!username) return 'Имя пользователя обязательно';
  if (username.length < 3) return 'Имя пользователя должно быть не менее 3 символов';
  return '';
};

export const validateConfirmPassword = (password: string, confirmPassword: string): string => {
  if (!confirmPassword) return 'Подтверждение пароля обязательно';
  if (password !== confirmPassword) return 'Пароли не совпадают';
  return '';
};