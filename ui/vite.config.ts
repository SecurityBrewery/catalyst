import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import tailwindcss from '@tailwindcss/vite'
import autoprefixer from "autoprefixer"

// https://vitejs.dev/config/
export default defineConfig({
  base: '/ui/',
  css: {
    postcss: {
      plugins: [autoprefixer()],
    },
  },
  plugins: [
    tailwindcss(),
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
