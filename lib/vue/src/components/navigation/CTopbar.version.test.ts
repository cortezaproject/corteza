import { describe, it, expect, vi, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import PrimeVue from 'primevue/config'
import CTopbar from './CTopbar.vue'

// The build version is the profile menu's last entry; the topbar carries no
// other menu for it.

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
  window.matchMedia = (q: string) =>
    ({
      matches: false,
      media: q,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
    }) as unknown as MediaQueryList
  vi.stubGlobal('VERSION', '2026.9.0-test')
})

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/:p(.*)*', component: { template: '<div />' } }],
})

async function mountTopbar() {
  const pinia = createPinia()
  setActivePinia(pinia)

  const wrapper = mount(CTopbar, {
    attachTo: document.body,
    props: {
      settings: { hideAgentSidebar: true, hideNotifications: true },
      labels: {
        userSettingsProfile: 'Profile',
        userSettingsChangePassword: 'Change password',
        userSettingsLogout: 'Logout',
        userSettingsTheme: 'Theme',
        lightTheme: 'Light',
        darkTheme: 'Dark',
        version: 'Version:',
      },
    },
    global: {
      plugins: [pinia, router, PrimeVue],
      provide: {
        $SystemAPI: { baseURL: '' },
        $Settings: { get: (_k: string, d?: unknown) => d },
        $Auth: {
          user: { name: 'Dev Agent', email: 'agent@local.dev', meta: {} },
          authURL: '',
          logout: vi.fn(),
        },
      },
    },
  })
  await flushPromises()
  return wrapper
}

describe('CTopbar version entry', () => {
  it('has no help menu button', async () => {
    const wrapper = await mountTopbar()
    expect(wrapper.find('[data-test-id="dropdown-helper"]').exists()).toBe(false)
    expect(wrapper.find('.pi-dollar').exists()).toBe(false)
    wrapper.unmount()
  })

  it('states the build version last in the profile menu', async () => {
    const wrapper = await mountTopbar()
    await wrapper.find('[data-test-id="dropdown-profile"]').trigger('click')
    await flushPromises()
    const labels = Array.from(document.body.querySelectorAll('[data-pc-section="itemlabel"]')).map(
      el => el.textContent?.trim(),
    )
    expect(labels.at(-1)).toBe('Version: 2026.9.0-test')
    expect(labels.indexOf('Logout')).toBe(labels.length - 2)
    wrapper.unmount()
  })
})
