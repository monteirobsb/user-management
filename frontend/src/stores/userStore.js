import { defineStore } from 'pinia';
import apiService from '../services/api.js';

export const useUserStore = defineStore('user', {
    state: () => ({
        users: [],
        loading: false,
        error: null,
    }),
    actions: {
        async fetchUsers() {
            this.loading = true;
            this.error = null;
            try {
                const response = await apiService.getUsers();
                this.users = response.data;
            } catch (error) {
                this.error = 'Falha ao buscar usuários.';
                console.error(error);
            } finally {
                this.loading = false;
            }
        },
        async addUser(user) {
            this.loading = true;
            this.error = null;
            try {
                await apiService.createUser(user);
                await this.fetchUsers(); // Re-carrega a lista após adicionar
            } catch (error) {
              // Captura a mensagem de erro da API
                this.error = error.response?.data?.error || 'Falha ao adicionar usuário.';
                console.error(error);
                throw error; // Propaga o erro para o componente, se necessário
            } finally {
                this.loading = false;
            }
        },
        async updateUser(id, user) {
            this.loading = true;
            this.error = null;
            try {
                await apiService.updateUser(id, user);
                await this.fetchUsers(); // Re-carrega a lista após atualizar
            } catch (error) {
                this.error = 'Falha ao atualizar usuário.';
                console.error(error);
            } finally {
                this.loading = false;
            }
        },
        async removeUser(id) {
            this.loading = true;
            this.error = null;
            try {
                await apiService.deleteUser(id);
                await this.fetchUsers(); // Re-carrega a lista após remover
            } catch (error) {
                this.error = 'Falha ao remover usuário.';
                console.error(error);
            } finally {
                this.loading = false;
            }
        },
    },
});