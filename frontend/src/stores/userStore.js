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
                const response = await apiService.createUser(user);
                // Assume que response.data contém o usuário criado, incluindo o ID do servidor
                if (response && response.data) {
                    this.users.push(response.data);
                } else {
                    // Se a API não retornar o usuário criado, pode ser necessário recarregar
                    // ou lidar com isso de outra forma, mas idealmente a API retorna o novo recurso.
                    await this.fetchUsers(); // Fallback, mas idealmente não necessário
                }
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
                const response = await apiService.updateUser(id, user);
                 // Assume que response.data contém o usuário atualizado
                if (response && response.data) {
                    const index = this.users.findIndex(u => u.id === id);
                    if (index !== -1) {
                        this.users[index] = response.data;
                    } else {
                        // Usuário não encontrado na lista local, talvez recarregar?
                        await this.fetchUsers(); // Fallback
                    }
                } else {
                    await this.fetchUsers(); // Fallback se a API não retornar o usuário atualizado
                }
            } catch (error) {
                this.error = error.response?.data?.error || 'Falha ao atualizar usuário.';
                console.error(error);
                throw error; // Propaga o erro para o componente
            } finally {
                this.loading = false;
            }
        },
        async removeUser(id) {
            this.loading = true;
            this.error = null;
            try {
                await apiService.deleteUser(id);
                // Remove o usuário da lista local
                const initialLength = this.users.length;
                this.users = this.users.filter(u => u.id !== id);
                if (this.users.length === initialLength) {
                    // Se o usuário não foi encontrado para filtro (improvável se a UI está sincronizada),
                    // ou se a deleção falhou silenciosamente na API antes do catch (também improvável com await),
                    // pode ser necessário um fallback.
                    // No entanto, filter é geralmente seguro.
                }
            } catch (error) {
                this.error = error.response?.data?.error || 'Falha ao remover usuário.';
                console.error(error);
                throw error; // Propaga o erro para o componente
            } finally {
                this.loading = false;
            }
        },
    },
});