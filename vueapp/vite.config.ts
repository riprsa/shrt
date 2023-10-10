import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: "./../view",
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': {
        target: 'https://s.x16.me', // Your API server's URL
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, 'api'),
        headers: {
          'Access-Control-Allow-Origin': '*', // Set CORS to allow all origins
          // You can also specify other CORS headers if needed:
          // 'Access-Control-Allow-Headers': 'Content-Type, Authorization',
          // 'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE',
        },
      },
    },
  }
})
