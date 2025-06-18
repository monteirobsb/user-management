<template>
  <div class="users-view">
    <header>
      <h1>Gerenciamento de Usuários</h1>
      <div>
        <button @click="openForm(null)" class="btn-primary">Adicionar Usuário</button>
        <button @click="handleLogout" class="btn-logout">Sair</button>
      </div>
    </header>

    <div v-if="store.loading">Carregando...</div>
    <div v-if="store.error" class="error">{{ store.error }}</div>

    <button @click="openForm(null)" v-if="!isFormVisible">Adicionar Novo Usuário</button>

    <UserForm
      v-if="isFormVisible"
      :user-to-edit="userToEdit"
      @submit="handleFormSubmit"
      @close="closeForm"
    />

    <UserTable
      :users="store.users"
      @edit="openForm"
      @delete="handleDelete"
    />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useUserStore } from '../stores/userStore';
import UserTable from '../components/UserTable.vue';
import UserForm from '../components/UserForm.vue';
import { useAuthStore } from '@/stores/authStore';

const authStore = useAuthStore();
const store = useUserStore();

const isFormVisible = ref(false);
const userToEdit = ref(null);

onMounted(() => {
  store.fetchUsers();
});

const openForm = (user) => {
  userToEdit.value = user;
  isFormVisible.value = true;
};

const closeForm = () => {
  isFormVisible.value = false;
  userToEdit.value = null;
};

const handleFormSubmit = async (userData) => {
  if (userData.id) {
    // Atualizar
    await store.updateUser(userData.id, { name: userData.name, email: userData.email });
  } else {
    // Criar
    await store.addUser({ name: userData.name, email: userData.email, password: userData.password });
  }
  if (!store.error) {
    closeForm();
  }
};

const handleDelete = async (id) => {
  if (confirm('Tem certeza de que deseja remover este usuário?')) {
    await store.removeUser(id);
  }
};

const handleLogout = () => {
  authStore.logout();
  router.push('/login'); // Redirecionamento feito no componente
};
</script>

<style scoped>
.users-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.error {
  color: red;
  margin-bottom: 15px;
}
button {
  margin-bottom: 15px;
}
/* ... (outros estilos) */
.btn-logout {
  background-color: #6c757d;
  color: white;
  margin-left: 10px;
}
</style>