import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

const apiClient = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor para adicionar o token de autenticação a cada requisição
apiClient.interceptors.request.use(config => {
    // É necessário instanciar a store dentro do interceptor
    const authStore = useAuthStore();
    const token = authStore.token;
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, error => {
    return Promise.reject(error);
});

// Exporta um objeto com métodos nomeados explicitamente
export default {
  // --- Auth ---
  login(credentials) {
    return apiClient.post('/login', credentials);
  },

  // --- Users ---
  getUsers() {
    return apiClient.get('/users');
  },
  getUser(id) {
    return apiClient.get(`/users/${id}`);
  },
  createUser(user) {
    return apiClient.post('/users', user);
  },
  updateUser(id, user) {
    return apiClient.put(`/users/${id}`, user);
  },
  deleteUser(id) {
    return apiClient.delete(`/users/${id}`);
  },
};