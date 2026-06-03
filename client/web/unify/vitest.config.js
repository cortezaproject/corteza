import { fileURLToPath } from 'node:url'
import { defineConfig, configDefaults } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    exclude: [...configDefaults.exclude, 'e2e/**'],
    root: fileURLToPath(new URL('./', import.meta.url)),
    // Inline pinia + its devtools-api dep so Vite resolves their package
    // export maps (raw Node ESM resolution picks a non-existent entry).
    server: {
      deps: {
        inline: ['pinia', '@vue/devtools-api'],
      },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  define: {
    VERSION: JSON.stringify('test'),
    BUILD_TIME: JSON.stringify('test'),
  },
})
