import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// The layout list's second column reads back the rule the public views apply:
// layouts are tried in weight order and the first whose condition and roles
// both pass is the one shown (composables/usePageVisibility.ts).

const route = reactive({
  name: 'admin.pages.edit',
  params: { slug: 'ns', pageID: 'P1' },
  query: {},
})

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  onBeforeRouteLeave: () => {},
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('primevue/useconfirm', () => ({ useConfirm: () => ({ require: vi.fn() }) }))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: [], exprScope: [] }),
}))

vi.mock('@/sections/compose/composables/useResourceTranslations', () => ({
  useResourceTranslations: () => ({
    showTranslatorButton: ref(false),
    currentLanguage: ref('en'),
    resourceTranslationsEnabled: ref(false),
    canManageResourceTranslations: ref(false),
  }),
}))

let page
let layouts

const pageStore = {
  getByID: vi.fn(() => page),
  findByID: vi.fn(() => Promise.resolve(page)),
  update: vi.fn(p => Promise.resolve(p)),
}
const pageLayoutStore = {
  getByPageID: vi.fn(() => layouts),
  findByPageID: vi.fn(() => Promise.resolve(layouts)),
  create: vi.fn(),
  update: vi.fn(),
  delete: vi.fn(),
}

vi.mock('@planetcrust/human-vue', () => ({
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useModuleStore: () => ({ getByID: () => null }),
  useHistoryBack: () => vi.fn(),
  useDraftGuard: () => ({
    isDirty: { value: false },
    capture: vi.fn(),
    reset: vi.fn(),
    markSaved: vi.fn(),
  }),
  useRBACStore: () => ({ can: () => true, canGlobal: () => true }),
  useUnsavedGuard: () => ({ markSaved: vi.fn() }),
  useFileUpload: () => ({
    uploading: { value: false },
    uploadError: { value: '' },
    uploadFileRaw: vi.fn(),
    reset: vi.fn(),
  }),
  components: {
    CInputDelete: { template: '<div />' },
    CInputToggleCard: { props: ['modelValue'], template: '<div />' },
    CFileDropZone: { template: '<div />' },
  },
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: {
    Page: class {
      constructor(p = {}) {
        Object.assign(this, p)
      }
    },
    PageLayout: class {
      constructor(l = {}) {
        this.meta = { title: '', description: '' }
        this.config = { useTitle: false, visibility: { expression: '', roles: [] } }
        Object.assign(this, l)
      }
    },
  },
}))

import Edit from './Edit.vue'

const FormListStub = {
  name: 'CFormList',
  props: {
    modelValue: Array,
    columns: Array,
    emptyMessage: String,
    draggable: Boolean,
    hideRemove: Boolean,
    fitWidth: Boolean,
  },
  template: `<div>
    <div v-for="(item, index) in modelValue" :key="index" class="layout-row">
      <slot name="row" :item="item" :index="index" />
    </div>
  </div>`,
}

const GLOBAL_COMPONENTS = Object.fromEntries(
  [
    'Button',
    'ButtonGroup',
    'CEditorActions',
    'CFormGroup',
    'CInputRole',
    'CPermissionsButton',
    'TieredMenu',
    'CInputToggleCard',
    'Checkbox',
    'Dialog',
    'Divider',
    'Fieldset',
    'InputNumber',
    'Message',
    'ProgressSpinner',
    'Select',
    'Textarea',
    'ToggleSwitch',
    'CInputExpression',
    'CExpressionHint',
    'InputText',
    'CInputModuleField',
    'PageTranslator',
  ].map(name => [name, true]),
)

const passthrough = name => [name, { name, template: '<div><slot /></div>' }]

const layoutWith = visibility => [
  {
    pageLayoutID: 'L1',
    pageID: 'P1',
    namespaceID: 'N1',
    meta: { title: 'Account', description: '' },
    config: { useTitle: false, visibility },
    blocks: [],
  },
]

// `$t` here has to carry its arguments through: the roles line is the only
// place the resolved names reach the DOM.
const t = (key, args) => (args ? `${key}:${args.join(',')}` : key)

let roleList

async function mountEdit() {
  const wrapper = mount(Edit, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      plugins: [createPinia()],
      stubs: {
        teleport: true,
        Teleport: true,
        ...GLOBAL_COMPONENTS,
        ...Object.fromEntries([
          passthrough('Panel'),
          passthrough('InputGroup'),
          passthrough('InputGroupAddon'),
        ]),
        Form: {
          name: 'Form',
          props: ['resolver', 'initialValues'],
          template: '<div><slot /></div>',
        },
        CFormList: FormListStub,
      },
      renderStubDefaultSlot: true,
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: t },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: { iconList: () => Promise.resolve({ set: [] }), baseURL: '' },
        $SystemAPI: { roleList },
        $Settings: { get: () => undefined },
        $Auth: { user: { userID: 'U1', roles: [] } },
        $eventBus: null,
      },
    },
  })
  await flushPromises()
  await flushPromises()
  return wrapper
}

// The applies-when cell is the row's second column.
const summary = w => w.findAll('.layout-row > div')[1].text()

beforeEach(() => {
  page = { pageID: 'P1', namespaceID: 'N1', title: 'Account', handle: 'account', blocks: [] }
  roleList = vi.fn(() => Promise.resolve({ set: [{ roleID: 'R1', name: 'Manager' }] }))
})

describe('layout applies-when column', () => {
  it('reads Always for a layout nothing restricts', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    expect(summary(w)).toBe('page.page-layout.appliesWhen.always')
    expect(roleList).not.toHaveBeenCalled()
  })

  it('shows the condition expression', async () => {
    layouts = layoutWith({ expression: 'screen.width < 1024', roles: [] })
    const w = await mountEdit()
    expect(summary(w)).toContain('screen.width < 1024')
    expect(summary(w)).not.toContain('appliesWhen.always')
  })

  it('names the roles it is restricted to', async () => {
    layouts = layoutWith({ expression: '', roles: ['R1'] })
    const w = await mountEdit()
    expect(roleList).toHaveBeenCalledWith({ roleID: ['R1'], limit: 1 })
    expect(summary(w)).toContain('Manager')
  })

  it('falls back to the id when the name never arrives', async () => {
    roleList = vi.fn(() => Promise.reject(new Error('nope')))
    layouts = layoutWith({ expression: '', roles: ['R9'] })
    const w = await mountEdit()
    expect(summary(w)).toContain('R9')
  })

  it('shows the condition and the roles together', async () => {
    layouts = layoutWith({ expression: 'isEdit', roles: ['R1'] })
    const w = await mountEdit()
    expect(summary(w)).toContain('isEdit')
    expect(summary(w)).toContain('Manager')
  })
})
