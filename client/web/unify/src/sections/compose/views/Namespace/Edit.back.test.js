import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive, ref } from 'vue'

// The editor is reachable without a screen before it — a deep link, or the
// bounce a disabled namespace performs. There router.back() walks out of the
// app (or does nothing at all), so Back has to fall back to the namespace list.

const route = reactive({ name: 'namespace.edit', params: { slug: 'ns' }, query: {} })

const history = { state: {} }
const router = {
  push: vi.fn(),
  replace: vi.fn(),
  back: vi.fn(),
  options: { history },
}

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

const namespaceStore = {
  getByUrlPart: () => ({ namespaceID: 'N1', slug: 'ns', name: 'Ns', enabled: true, meta: {} }),
  load: vi.fn(() => Promise.resolve([])),
  update: vi.fn(),
  create: vi.fn(),
  clone: vi.fn(),
  delete: vi.fn(),
}

vi.mock('@planetcrust/human-vue', () => ({
  useNamespaceStore: () => namespaceStore,
  useUnsavedGuard: () => ({ markSaved: vi.fn() }),
  useFileUpload: () => ({
    uploading: ref(false),
    uploadError: ref(null),
    uploadFileRaw: vi.fn(),
    reset: vi.fn(),
  }),
  components: {
    CInputDelete: { template: '<div />' },
    CFileDropZone: { template: '<div />' },
    CInputLabel: { template: '<div />' },
  },
}))

vi.mock('@planetcrust/human-js', () => ({
  compose: {
    Namespace: class {
      constructor(opts = {}) {
        Object.assign(this, opts)
      }
    },
  },
}))

vi.mock('@/sections/compose/components/Namespaces/NamespaceTranslator.vue', () => ({
  default: { template: '<div />' },
}))

import Edit from './Edit.vue'

// Stands in for the real footer so the test can fire its `back` event without
// depending on PrimeVue's button being registered.
const EditorActions = {
  name: 'CEditorActions',
  props: ['backTo', 'backLabel'],
  emits: ['back'],
  template: '<div class="editor-actions" @click="$emit(\'back\')"><slot /></div>',
}

let wrapper

async function mountEditor() {
  wrapper = mount(Edit, {
    global: {
      components: { CEditorActions: EditorActions },
      stubs: { teleport: true, Teleport: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k, $router: router, $route: route },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastWarning: vi.fn() },
        $ComposeAPI: { baseURL: '', accessTokenFn: () => '' },
        $Settings: { attachment: () => '' },
      },
    },
  })

  await flushPromises()
  return wrapper
}

beforeEach(() => {
  history.state = {}
  router.push.mockClear()
  router.back.mockClear()
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Namespace Edit back button', () => {
  it('returns to the previous screen when there is one', async () => {
    history.state = { back: '/compose/namespaces' }

    await mountEditor()
    await wrapper.find('.editor-actions').trigger('click')

    expect(router.back).toHaveBeenCalled()
    expect(router.push).not.toHaveBeenCalled()
  })

  it('falls back to the namespace list when the editor is the first entry', async () => {
    history.state = { back: null }

    await mountEditor()
    await wrapper.find('.editor-actions').trigger('click')

    expect(router.push).toHaveBeenCalledWith({ name: 'namespace.list' })
    expect(router.back).not.toHaveBeenCalled()
  })
})
