export const handleApiError = (error: any): string => {
  if (error.response) {
    switch (error.response.status) {
      case 400:
        return 'Неверный запрос. Проверьте введенные данные';
      case 401:
        return 'Необходима авторизация';
      case 403:
        return 'Доступ запрещен';
      case 404:
        return 'Ресурс не найден';
      case 500:
        return 'Внутренняя ошибка сервера';
      default:
        return `Ошибка: ${error.response.status}`;
    }
  } else if (error.request) {
    return 'Нет ответа от сервера. Проверьте подключение к интернету';
  } else {
    return 'Неизвестная ошибка';
  }
};