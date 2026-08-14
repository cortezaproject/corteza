import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive, ref } from 'vue'

// The editor is reachable without a screen before it — a deep link, or the
// bounce a disabled namespace performs. Where Back lands then is decided by
// useHistoryBack (covered in lib/vue); what this screen owns is which fallback
// it names, and the namespace list is the only one that always exists.

const route = reactive({ name: 'namespace.edit', params: { slug: 'ns' }, query: {} })

const router = {
  push: vi.fn(),
  replace: vi.fn(),
  back: vi.fn(),
}

const goBack = vi.fn()

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
  useHistoryBack: () => goBack,
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
  router.push.mockClear()
  router.back.mockClear()
  goBack.mockClear()
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Namespace Edit back button', () => {
  it('names the namespace list as the fallback', async () => {
    await mountEditor()
    await wrapper.find('.editor-actions').trigger('click')

    expect(goBack).toHaveBeenCalledWith({ name: 'namespace.list' })
  })

  it('never navigates on its own', async () => {
    await mountEditor()
    await wrapper.find('.editor-actions').trigger('click')

    expect(router.back).not.toHaveBeenCalled()
    expect(router.push).not.toHaveBeenCalled()
  })
})
