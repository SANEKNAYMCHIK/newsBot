import apiClient from './axiosConfig';

export const sourcesApi = {
  // Получить все источники
  getSources: () =>
    apiClient.get(`/news/sources`),
  
  // Получить категории
  getCategories: () =>
    apiClient.get('/news/categories'),
  
  // Добавить источник (админ)
  addSource: (sourceData) =>
    apiClient.post('/admin/sources', sourceData),
  
  // Обновить источник (админ)
  updateSource: (sourceId, data) =>
    apiClient.put(`/admin/sources/${sourceId}`, data),
  
  // Удалить источник (админ)
  deleteSource: (sourceId) =>
    apiClient.delete(`/admin/sources/${sourceId}`),
};