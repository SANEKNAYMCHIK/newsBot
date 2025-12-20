import axios from 'axios';

// Создаем экземпляр axios с настройками
const apiClient = axios.create({
  baseURL: 'https://localhost:8443', // Изменено на HTTP и порт 8080
  timeout: 15000, // Увеличено время ожидания
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: false, // Установлено false, так как используем Bearer токен
});

// Интерцептор для логирования и добавления токена
apiClient.interceptors.request.use(
  (config) => {
    // Получаем токен из localStorage (главный источник)
    const token = localStorage.getItem('token');
    
    console.log('🔑 Axios Request:', {
      url: config.url,
      method: config.method,
      hasToken: !!token,
      tokenPreview: token ? token.substring(0, 30) + '...' : 'No token'
    });
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    } else {
      console.warn('⚠️ No token found for request to:', config.url);
    }
    
    return config;
  },
  (error) => {
    console.error('❌ Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// Интерцептор для обработки ответов
apiClient.interceptors.response.use(
  (response) => {
    console.log('✅ Axios Response:', {
      url: response.config.url,
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error) => {
    console.error('❌ Axios Response Error:', {
      url: error.config?.url,
      status: error.response?.status,
      message: error.message,
      data: error.response?.data
    });
    
    if (error.response?.status === 401) {
      console.log('🚫 Unauthorized (401) - clearing auth data');
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // Редирект только если не на странице логина
      if (!window.location.pathname.includes('/login') && 
          !window.location.pathname.includes('/register')) {
        setTimeout(() => {
          window.location.href = '/login?sessionExpired=true';
        }, 100);
      }
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;