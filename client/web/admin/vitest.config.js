import { fileURLToPath } from 'node:url'
import { defineConfig, configDefaults } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    exclude: [...configDefaults.exclude, 'e2e/**'],
    include: ['src/**/*.test.ts'],
    coverage: {
      provider: 'v8',
      include: ['src/stores/**', 'src/composables/**', 'src/components/**'],
      exclude: ['src/test/**'],
      thresholds: { lines: 60, functions: 60 },
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
  ssr: {
    noExternal: ['pinia'],
  },
})
