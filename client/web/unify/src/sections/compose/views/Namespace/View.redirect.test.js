import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive, ref } from 'vue'

// A namespace this view refuses to render bounces the user elsewhere. The bounce
// has to replace the current entry: pushing leaves the refused URL in history, so
// Back returns to it, the view runs again and bounces again — the editor becomes
// a screen you cannot leave.

const route = reactive({ name: 'pages', params: { slug: 'ns' }, query: {} })

const router = { push: vi.fn(), replace: vi.fn(), back: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

let namespaces = []

const namespaceStore = {
  getByUrlPart: urlPart => namespaces.find(ns => ns.slug === urlPart) || null,
  load: vi.fn(() => Promise.resolve(namespaces)),
}

const emptyStore = () => ({ clearSet: vi.fn(), load: vi.fn(() => Promise.resolve()), set: [] })
const moduleStore = emptyStore()
const pageStore = emptyStore()
const chartStore = emptyStore()
const pageLayoutStore = emptyStore()

vi.mock('@planetcrust/human-vue', () => ({
  useNamespaceStore: () => namespaceStore,
  useModuleStore: () => moduleStore,
  usePageStore: () => pageStore,
  useChartStore: () => chartStore,
  usePageLayoutStore: () => pageLayoutStore,
  // The real one delays for the spinner floor; tests want the work, not the wait.
  useMinDuration: () => ({ loading: ref(false), run: fn => fn() }),
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: {
    Namespace: class {
      constructor(opts = {}) {
        Object.assign(this, opts)
      }
    },
  },
}))

vi.mock('@/sections/compose/components/Record/RecordModal.vue', () => ({
  default: { template: '<div />' },
}))

import View from './View.vue'

let wrapper

async function mountView(slug = 'ns') {
  wrapper = mount(View, {
    props: { slug },
    global: {
      stubs: { teleport: true, Teleport: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastDanger: vi.fn(), toastWarning: vi.fn() },
      },
      renderStubDefaultSlot: false,
    },
    shallow: true,
  })

  await flushPromises()
  return wrapper
}

beforeEach(() => {
  route.name = 'pages'
  route.params = { slug: 'ns' }
  namespaces = []
  router.push.mockClear()
  router.replace.mockClear()
  namespaceStore.load.mockClear()
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Namespace View redirects', () => {
  it('replaces the refused URL when sending an admin to the editor', async () => {
    namespaces = [{ namespaceID: 'N1', slug: 'ns', enabled: false, canUpdateNamespace: true }]

    await mountView()

    expect(router.push).not.toHaveBeenCalled()
    expect(router.replace).toHaveBeenCalledWith({
      name: 'namespace.edit',
      params: { slug: 'ns' },
    })
  })

  it('replaces the refused URL when sending a non-admin away', async () => {
    namespaces = [{ namespaceID: 'N1', slug: 'ns', enabled: false, canUpdateNamespace: false }]

    await mountView()

    expect(router.push).not.toHaveBeenCalled()
    expect(router.replace).toHaveBeenCalledWith({ name: 'root' })
  })

  it('replaces the refused URL when the namespace cannot be found', async () => {
    namespaces = []

    await mountView('nope')

    expect(router.push).not.toHaveBeenCalled()
    expect(router.replace).toHaveBeenCalledWith({ name: 'root' })
  })

  it('lets an admin through to a disabled namespace on an admin route', async () => {
    route.name = 'admin.modules'
    namespaces = [{ namespaceID: 'N1', slug: 'ns', enabled: false, canUpdateNamespace: true }]

    await mountView()

    expect(router.replace).not.toHaveBeenCalled()
    expect(moduleStore.load).toHaveBeenCalled()
  })
})
