import apiClient from './axiosConfig';

export const newsApi = {
  // Получаем новости пользователя с пагинацией
  getUserNews: (page = 1, pageSize = 20) => 
    apiClient.get(`/news/?page=${page}&page_size=${pageSize}`),
  
  // Получаем новости конкретного источника с пагинацией
  getNewsBySource: (sourceId, page = 1, pageSize = 20) => 
    apiClient.get(`/news/source/${sourceId}?page=${page}&page_size=${pageSize}`),
  
  // Получаем источники с пагинацией
  getSources: (page = 1, pageSize = 100) => 
    apiClient.get(`/news/sources?page=${page}&page_size=${pageSize}`),
  
  // Получаем подписки пользователя
  getUserSubscriptions: () =>
    apiClient.get('/user/subscriptions/'),
  
  // Подписаться на источник
  subscribe: (sourceId) =>
    apiClient.post('/user/subscriptions/', { source_id: sourceId }),
  
  // Отписаться от источника
  unsubscribe: (sourceId) =>
    apiClient.delete(`/user/subscriptions/${sourceId}`),
  
  // Получить категории
  getCategories: () =>
    apiClient.get('/news/categories'),
  
  // Добавить новый источник
  addSource: (sourceData) =>
    apiClient.post('/news/sources', sourceData),

  refreshNews: () => 
    apiClient.post('/user/refresh'),

  getAllSources: (page = 1, pageSize = 20) => 
    apiClient.get(`/news/all-sources?page=${page}&page_size=${pageSize}`),
  
  // Получение всех пользователей с пагинацией (админский эндпоинт)
  getAllUsers: (page = 1, pageSize = 20) => 
    apiClient.get(`/admin/users?page=${page}&page_size=${pageSize}`),
  
  makeAdmin: (userId) => 
    apiClient.post(`/admin/users/${userId}/make-admin`),
  
  removeAdmin: (userId) => 
    apiClient.post(`/admin/users/${userId}/remove-admin`),
  
  updateSource: (sourceId, data) => 
    apiClient.put(`/admin/sources/${sourceId}`, data),
  
  deleteSource: (sourceId) => 
    apiClient.delete(`/admin/sources/${sourceId}`),
  
  addCategory: (data) => 
    apiClient.post('/admin/categories', data),
};
