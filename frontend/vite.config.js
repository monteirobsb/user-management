import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],

  // Esta seção é a CHAVE para resolver o erro.
  // Ela diz ao Vite (e ao Rollup) que o atalho "@" aponta para a pasta "/src".
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },

  // Esta seção é para o servidor de desenvolvimento local.
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Aponta para o backend Go
        changeOrigin: true,
      }
    }
  }
})