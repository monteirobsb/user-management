import { createRouter, createWebHistory } from 'vue-router'
import UsersView from '../views/UsersView.vue'
import LoginView from '../views/LoginView.vue'
import { useAuthStore } from '@/stores/authStore'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'users',
      component: UsersView,
      meta: { requiresAuth: true } // Marca a rota como protegida
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    }
  ]
})

// Guarda de navegação global
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // Se a rota requer autenticação e o usuário não está logado, redireciona para o login
    next({ name: 'login' });
  } else if (to.name === 'login' && authStore.isAuthenticated) {
    // Se o usuário já está logado e tenta acessar a página de login, redireciona para a home
    next({ name: 'users' });
  } else {
    // Caso contrário, permite a navegação
    next();
  }
});

export default router