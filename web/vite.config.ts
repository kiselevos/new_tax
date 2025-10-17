import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: Number(process.env.VITE_PORT) || 8080,
    proxy: {
      // 🔹 Прокси для /tax.TaxService (ConnectRPC)
      '/tax.TaxService': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      // 🔹 Универсальный прокси для всех /api/* маршрутов
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''), // ← убирает /api перед отправкой
      },
    },
  },
})