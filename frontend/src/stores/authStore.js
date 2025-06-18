import { defineStore } from 'pinia';
import router from '@/router';
import api from '@/services/api'; // Importa nosso serviço de API refatorado

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    loginError: null, // Added loginError state
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
  actions: {
    async login(credentials) {
      try {
        this.loginError = null; // Reset login error
        // AQUI ESTÁ A MUDANÇA: Usamos a função explícita 'api.login'
        const response = await api.login(credentials); 
        
        const token = response.data.token;
        this.token = token;
        localStorage.setItem('token', token);
        router.push('/');
      } catch (error) {
        console.error("Falha no login:", error);
        // Set user-friendly error message
        this.loginError = "Falha no login. Verifique seu e-mail e senha.";
        // Optionally append backend error:
        // if (error.response && error.response.data && error.response.data.error) {
        //   this.loginError += ` Detalhe: ${error.response.data.error}`;
        // }
        throw error; // Re-throw for potential component-level handling
      }
    },
    logout() {
      this.token = null;
      this.loginError = null; // Reset login error on logout
      localStorage.removeItem('token');
      router.push('/login');
    },
  },
});