import apiClient from './axiosConfig';

export const adminApi = {
  getStats: () => apiClient.get('/admin/stats'),
  getUsers: (page = 1, pageSize = 50) => 
    apiClient.get(`/admin/users?page=${page}&page_size=${pageSize}`),
  makeAdmin: (userId) => apiClient.post(`/admin/users/${userId}/make-admin`),
  addSource: (sourceData) => apiClient.post('/admin/sources', sourceData),
};