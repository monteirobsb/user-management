<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin" class="login-form">
      <h2>Login</h2>
      <div class="form-group">
        <label for="email">Email</label>
        <input type="email" v-model="email" required />
      </div>
      <div class="form-group">
        <label for="password">Senha</label>
        <input type="password" v-model="password" required />
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit">Entrar</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/authStore';
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const error = ref(null);
const authStore = useAuthStore();
// O router é inicializado aqui, no escopo correto do setup.
const router = useRouter();

const handleLogin = async () => {
  error.value = null; // Limpa erros anteriores
  try {
    // Espera a action da store ser concluída
    await authStore.login({ email: email.value, password: password.value });
    
    // Se a linha acima não gerou erro, o login foi bem-sucedido.
    // AGORA sim fazemos o redirecionamento.
    router.push('/');

  } catch (err) {
    // Se a action 'login' lançou um erro, ele será capturado aqui.
    error.value = 'Email ou senha inválidos.';
  }
};
</script>

<style scoped>
/* ... seus estilos ... */
.login-container {
  max-width: 400px;
  margin: 5rem auto;
  padding: 2rem;
  border: 1px solid #ccc;
  border-radius: 8px;
}
.login-form h2 {
  text-align: center;
}
.error {
  color: red;
  text-align: center;
  margin-top: 1rem;
}
</style>