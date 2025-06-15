import { defineStore } from 'pinia';
import router from '@/router';
import api from '@/services/api'; // Importa nosso serviço de API refatorado

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
  actions: {
    async login(credentials) {
      try {
        // AQUI ESTÁ A MUDANÇA: Usamos a função explícita 'api.login'
        const response = await api.login(credentials); 
        
        const token = response.data.token;
        this.token = token;
        localStorage.setItem('token', token);
        router.push('/');
      } catch (error) {
        console.error("Falha no login:", error);
        throw error;
      }
    },
    logout() {
      this.token = null;
      localStorage.removeItem('token');
      router.push('/login');
    },
  },
});