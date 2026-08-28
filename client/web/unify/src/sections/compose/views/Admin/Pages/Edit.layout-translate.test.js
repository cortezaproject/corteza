import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// A page's translations are one set covering the page, its blocks and every
// layout. The topbar opens it at the top; a layout row opens the same dialog
// aimed at its own row, which is the only caller `highlight` has.

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

let translationsOn = true

vi.mock('@/sections/compose/composables/useResourceTranslations', () => ({
  useResourceTranslations: () => ({
    showTranslatorButton: ref(translationsOn),
    currentLanguage: ref('en'),
    resourceTranslationsEnabled: ref(translationsOn),
    canManageResourceTranslations: ref(translationsOn),
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

// Stands in for the one PageTranslator in the topbar. The row button reaches it
// by ref and asks it to open, so the stub only has to record the ask.
const opened = []
const PageTranslatorStub = {
  name: 'PageTranslator',
  props: ['page', 'namespace', 'layouts', 'highlight', 'disabled'],
  methods: {
    open(highlight) {
      opened.push(highlight)
    },
  },
  template: '<div class="page-translator" />',
}

const ButtonStub = {
  name: 'Button',
  props: ['icon', 'label', 'severity', 'size', 'text', 'rounded', 'outlined', 'disabled'],
  emits: ['click'],
  template: '<button :data-icon="icon" @click="$emit(\'click\', $event)"><slot /></button>',
}

const passthrough = name => [name, { name, template: '<div><slot /></div>' }]

const GLOBAL_COMPONENTS = Object.fromEntries(
  [
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
  ].map(name => [name, true]),
)

const layoutWith = (pageLayoutID = 'L1') => [
  {
    pageLayoutID,
    pageID: 'P1',
    namespaceID: 'N1',
    meta: { title: 'Account', description: '' },
    config: { useTitle: false, visibility: { expression: '', roles: [] } },
    blocks: [],
  },
]

async function mountEdit() {
  const wrapper = mount(Edit, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      plugins: [createPinia()],
      stubs: {
        teleport: true,
        Teleport: true,
        ...GLOBAL_COMPONENTS,
        ...Object.fromEntries([passthrough('Panel'), passthrough('InputGroup')]),
        InputGroupAddon: { name: 'InputGroupAddon', template: '<div><slot /></div>' },
        Form: {
          name: 'Form',
          props: ['resolver', 'initialValues'],
          template: '<div><slot /></div>',
        },
        Button: ButtonStub,
        CFormList: FormListStub,
        PageTranslator: PageTranslatorStub,
      },
      // The topbar tools live in a Teleport; without this the PageTranslator
      // the row button reaches by ref never mounts.
      renderStubDefaultSlot: true,
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: { iconList: () => Promise.resolve({ set: [] }), baseURL: '' },
        $SystemAPI: {},
        $Settings: { get: () => undefined },
        $Auth: { user: { userID: 'U1', roles: [] } },
        $eventBus: null,
      },
    },
  })
  await flushPromises()
  return wrapper
}

const translateButtons = w => w.findAll('.layout-row [data-icon="pi pi-language"]')
const titleCell = w => w.find('[data-layout-title]')
const actionCell = w => w.find('[data-layout-actions]')

beforeEach(() => {
  translationsOn = true
  opened.length = 0
  page = { pageID: 'P1', namespaceID: 'N1', title: 'Account', handle: 'account', blocks: [] }
  layouts = layoutWith()
})

describe('layout row translate button', () => {
  it('offers one per saved layout', async () => {
    const w = await mountEdit()
    expect(translateButtons(w)).toHaveLength(1)
  })

  it('withholds it while resource translations are off', async () => {
    translationsOn = false
    const w = await mountEdit()
    expect(translateButtons(w)).toHaveLength(0)
  })

  it('withholds it from a layout that has never been saved', async () => {
    layouts = layoutWith('0')
    const w = await mountEdit()
    expect(translateButtons(w)).toHaveLength(0)
  })

  it("opens the page's translator on that layout's own title row", async () => {
    const w = await mountEdit()
    await translateButtons(w)[0].trigger('click')
    expect(opened).toEqual([{ resource: 'compose:page-layout/N1/P1/L1', key: 'meta.title' }])
  })

  it('sits in the title cell, since meta.title is what it translates', async () => {
    const w = await mountEdit()
    expect(titleCell(w).find('[data-icon="pi pi-language"]').exists()).toBe(true)
    expect(actionCell(w).find('[data-icon="pi pi-language"]').exists()).toBe(false)
  })

  it("leaves configure, build and delete as the row's own actions", async () => {
    const w = await mountEdit()
    expect(
      actionCell(w)
        .findAll('button')
        .map(b => b.attributes('data-icon')),
    ).toEqual(['pi pi-cog', 'pi pi-wrench', 'pi pi-trash'])
  })
})
