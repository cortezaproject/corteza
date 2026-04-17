import { createPinia, setActivePinia } from 'pinia'
import { createApp } from 'vue'

/**
 * Create a fresh Pinia instance with the given provides registered so stores
 * that call inject() inside defineStore() can resolve them.
 *
 * Call in beforeEach to ensure isolation between tests.
 */
export function createTestPinia(provides: Record<string, unknown> = {}) {
  const pinia = createPinia()

  // Mount a throwaway app to host the provides so inject() works inside stores
  const app = createApp({ template: '<div/>' })
  for (const [key, value] of Object.entries(provides)) {
    app.provide(key, value)
  }
  app.use(pinia)
  app.mount(document.createElement('div'))

  setActivePinia(pinia)
  return pinia
}
