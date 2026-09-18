import { describe, it, expect, vi, beforeAll, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive } from 'vue'
import PrimeVue from 'primevue/config'
import Select from 'primevue/select'

const route = reactive({ name: 'pages', params: { slug: 'alpha' } })
const router = { push: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

const namespaceStore = {
  set: [
    { namespaceID: '1', slug: 'alpha', name: 'Alpha CRM', enabled: true },
    { namespaceID: '2', slug: 'beta', name: 'Beta Sales', enabled: true },
    { namespaceID: '3', slug: 'crm-legacy', name: 'Gamma', enabled: true },
  ],
}

vi.mock('@planetcrust/human-vue', () => ({
  useNamespaceStore: () => namespaceStore,
}))

import CSidebarNamespaceSwitcher from './CSidebarNamespaceSwitcher.vue'

beforeAll(() => {
  window.matchMedia = q => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

let wrapper

afterEach(() => {
  wrapper?.unmount()
  document.body.innerHTML = ''
})

const mountSwitcher = () =>
  mount(CSidebarNamespaceSwitcher, {
    attachTo: document.body,
    global: {
      plugins: [[PrimeVue, { unstyled: true }]],
      components: { Select },
      provide: { $Settings: { get: () => ({}) } },
      mocks: { $t: k => k },
      directives: { tooltip: {} },
      stubs: {
        FloatLabel: { template: '<div><slot /></div>' },
        Button: true,
        RouterLink: true,
      },
    },
  })

const optionLabels = () =>
  [...document.body.querySelectorAll('[role="option"]')].map(li => li.textContent.trim())

const open = async () => {
  await wrapper.find('[role="combobox"]').trigger('click')
  await flushPromises()
}

const search = async text => {
  const input = document.body.querySelector('input[role="searchbox"]')
  input.value = text
  input.dispatchEvent(new Event('input'))
  await flushPromises()
}

describe('CSidebarNamespaceSwitcher search', () => {
  it('narrows the list by name', async () => {
    wrapper = mountSwitcher()
    await open()
    expect(optionLabels()).toEqual(['Alpha CRM', 'Beta Sales', 'Gamma'])

    await search('sales')
    expect(optionLabels()).toEqual(['Beta Sales'])
  })

  it('matches the short name too', async () => {
    wrapper = mountSwitcher()
    await open()

    await search('crm')
    expect(optionLabels()).toEqual(['Alpha CRM', 'Gamma'])
  })

  it('keeps naming the current namespace while the search hides it', async () => {
    wrapper = mountSwitcher()
    await open()

    await search('sales')
    expect(wrapper.find('[data-pc-section="label"]').text()).toBe('Alpha CRM')
  })
})
