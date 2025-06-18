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
      <!-- Display login error from authStore -->
      <p v-if="authStore.loginError" class="error-message" style="color: red;">{{ authStore.loginError }}</p>
      <p v-if="error && !authStore.loginError" class="error">{{ error }}</p> <!-- Keep local error for other potential issues if needed, but prioritize store error -->
      <button type="submit">Entrar</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/authStore'; // Already here, good.
import { useRouter } from 'vue-router';

const email = ref('');
const password = ref('');
const error = ref(null); // Local error state for non-auth related issues, if any.
const authStore = useAuthStore(); // Get instance of the auth store.
const router = useRouter();

const handleLogin = async () => {
  error.value = null; // Limpa erros locais anteriores
  // authStore.loginError is reset inside the action itself.
  try {
    await authStore.login({ email: email.value, password: password.value });
    // Navigation is handled by the store now on successful login.
    // router.push('/'); // This line can be removed if store always redirects. Kept for clarity if store might not.
    // For this task, the store handles redirection, so router.push('/') here is redundant.
  } catch (err) {
    // The authStore.loginError will be set by the action.
    // We can set a local error if there's a different kind of error
    // or if we want to display a generic message not from the store.
    // However, the primary error display is now through authStore.loginError.
    // If `throw error` was removed from store, this `err` would be from other issues.
    // Since `throw error` is kept, `err` here is the one from the store.
    // We don't need to set local `error.value` if `authStore.loginError` is the source of truth for login failures.
    // error.value = 'Email ou senha inválidos.'; // This can be removed or kept for non-store related errors.
  }
};
</script>

<style scoped>
/* ... seus estilos ... */
.error-message { /* Style for the new error message display */
  color: red;
  text-align: center;
  margin-top: 1rem;
  margin-bottom: 1rem; /* Added for spacing */
}
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
.error { /* This is for the local error.value, if still used */
  color: red;
  text-align: center;
  margin-top: 1rem;
}
</style>