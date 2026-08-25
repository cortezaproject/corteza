import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// A field row carries one translate button and it covers the whole field.
// Select and Bool used to get a second, near-empty kebab menu beside it whose
// only entries were their own translation keys; the row's own button filtered
// exactly those keys out, so neither reached the other's.

const route = reactive({
  name: 'admin.modules.edit',
  params: { slug: 'ns', moduleID: 'M1' },
  query: {},
})

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  onBeforeRouteLeave: () => {},
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('primevue/useconfirm', () => ({ useConfirm: () => ({ require: vi.fn() }) }))

vi.mock('@/sections/compose/composables/useResourceTranslations', () => ({
  useResourceTranslations: () => ({
    showTranslatorButton: ref(true),
    currentLanguage: ref('en'),
    resourceTranslationsEnabled: ref(true),
    canManageResourceTranslations: ref(true),
  }),
}))

let module_

const moduleStore = {
  findByID: vi.fn(() => Promise.resolve(module_)),
  update: vi.fn(),
  create: vi.fn(),
  delete: vi.fn(),
}

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => moduleStore,
  usePageStore: () => ({ set: [], create: vi.fn(), update: vi.fn() }),
  useHistoryBack: () => vi.fn(),
  useConfirmDelete: () => ({ confirmDelete: vi.fn() }),
  usePermissions: () => ({ open: vi.fn() }),
  useRBACStore: () => ({ can: () => true, canGlobal: () => true }),
  useDraftGuard: () => ({ capture: vi.fn(), markSaved: vi.fn(), isDirty: { value: false } }),
  components: {
    CInputDelete: { template: '<div />' },
    CRouterLinkButton: { props: ['to', 'label', 'icon'], template: '<div />' },
  },
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: {
    Module: class {
      constructor(m = {}) {
        this.moduleID = 'M1'
        this.namespaceID = 'N1'
        this.name = 'Probe'
        this.handle = 'probe'
        this.fields = []
        this.issues = []
        this.config = {}
        Object.assign(this, m)
      }
      systemFields() {
        return []
      }
    },
    ModuleFieldMaker: ({ kind }) => ({
      kind,
      cap: { multi: kind !== 'Bool', required: true },
    }),
    ModuleFieldString: class {
      constructor() {
        this.kind = 'String'
        this.name = ''
        this.label = ''
        this.options = {}
      }
    },
  },
}))

import Edit from './Edit.vue'
import { useTranslatorStore } from '@/sections/compose/stores/translator'

const field = (name, kind) => ({
  fieldID: `F-${name}`,
  name,
  label: name,
  kind,
  options: {},
  isSystem: false,
  cap: { multi: true },
})

const FormListStub = {
  name: 'CFormList',
  // Typed, not an array: a valueless attribute reaches an array-props stub as
  // '' rather than true, and `fit-width` would read as absent.
  props: {
    modelValue: Array,
    columns: Array,
    emptyMessage: String,
    draggable: Boolean,
    stickyHeader: Boolean,
    fitWidth: Boolean,
    hideRemove: Boolean,
  },
  template: `<div>
    <div v-for="(item, index) in modelValue" :key="index" class="field-row">
      <slot name="row" :item="item" :index="index" />
    </div>
  </div>`,
}

const GLOBAL_COMPONENTS = Object.fromEntries(
  [
    'Button',
    'ButtonGroup',
    'Checkbox',
    'Dialog',
    'Divider',
    'InputText',
    'Menu',
    'Message',
    'ProgressSpinner',
    'Select',
    'Tab',
    'TabList',
  ].map(name => [name, true]),
)

// The fields list is several slots deep; a default stub renders none of them.
const passthrough = (name, slot = 'default') => [
  name,
  { name, template: `<div><slot name="${slot}" /></div>` },
]

const SLOT_COMPONENTS = Object.fromEntries([
  passthrough('Card', 'content'),
  passthrough('Tabs'),
  passthrough('TabPanels'),
  passthrough('TabPanel'),
  passthrough('InputGroup'),
  passthrough('InputGroupAddon'),
])

// A real element, so the row's buttons can be told apart by icon and clicked.
const ButtonStub = {
  name: 'Button',
  props: ['icon', 'label', 'severity', 'size', 'text', 'rounded', 'outlined', 'disabled'],
  emits: ['click'],
  template: '<button :data-icon="icon" @click="$emit(\'click\', $event)"><slot /></button>',
}

async function mountEdit() {
  const wrapper = mount(Edit, {
    props: { namespace: { namespaceID: 'N1', canManageNamespace: false } },
    global: {
      plugins: [createPinia()],
      provide: {
        $ComposeAPI: composeAPI,
        $SystemAPI: {},
        $Settings: { get: () => false },
        $toast: { toastDanger: vi.fn(), toastSuccess: vi.fn() },
      },
      mocks: { $t: k => k },
      stubs: {
        teleport: true,
        Teleport: true,
        ...GLOBAL_COMPONENTS,
        ...SLOT_COMPONENTS,
        Button: ButtonStub,
        CFormList: FormListStub,
        CFormGroup: { template: '<div><slot /></div>' },
        CFieldConfigurator: true,
        CEditorActions: true,
        DalSettings: true,
        UniqueValues: true,
        RecordRevisionsSettings: true,
        ModuleIssues: true,
        ModuleTranslator: true,
        DalSchemaAlterations: true,
        FederationSettings: true,
        DiscoverySettings: true,
        Form: { props: ['resolver', 'initialValues'], template: '<div><slot /></div>' },
      },
      directives: { tooltip: () => {} },
    },
  })
  await flushPromises()
  return wrapper
}

// Every translation row the server returns for this module.
const TRANSLATIONS = [
  { resource: 'compose:module/N1/M1', key: 'name', lang: 'en', message: 'Probe' },
  { resource: 'compose:module-field/N1/M1/F-status', key: 'label', lang: 'en', message: 'status' },
  {
    resource: 'compose:module-field/N1/M1/F-status',
    key: 'meta.description.view',
    lang: 'en',
    message: '',
  },
  {
    resource: 'compose:module-field/N1/M1/F-status',
    key: 'meta.options.new.text',
    lang: 'en',
    message: 'New',
  },
  {
    resource: 'compose:module-field/N1/M1/F-flag',
    key: 'meta.bool.true.label',
    lang: 'en',
    message: 'Yes',
  },
]

const composeAPI = {
  moduleListTranslations: vi.fn(() => Promise.resolve(TRANSLATIONS)),
  moduleUpdateTranslations: vi.fn(() => Promise.resolve()),
}

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    name: 'Probe',
    handle: 'probe',
    issues: [],
    config: {},
    fields: [field('status', 'Select'), field('flag', 'Bool'), field('title', 'String')],
  }
  composeAPI.moduleListTranslations.mockClear()
})

describe('module field translation', () => {
  it('gives no field row a second actions menu', async () => {
    const wrapper = await mountEdit()
    expect(wrapper.findAll('.field-row')).toHaveLength(3)
    expect(wrapper.html()).not.toContain('pi-ellipsis-v')
  })

  it('lets the field list share the width it is given', async () => {
    // Sized to content, the six columns are wider than any laptop viewport and
    // the row's trailing controls sit off screen behind a sideways scroll.
    const wrapper = await mountEdit()
    const list = wrapper.findComponent({ name: 'CFormList' })
    expect(list.props('fitWidth')).toBe(true)

    // Sharing the width only works if the flexible columns carry a floor —
    // a bare `1fr` collapses to nothing once the fixed columns are paid for.
    const flexible = list.props('columns').filter(c => (c.width || '').includes('fr'))
    expect(flexible.length).toBeGreaterThan(0)
    expect(flexible.every(c => c.width.startsWith('minmax('))).toBe(true)
  })

  it('gives every field the same one translate button', async () => {
    const wrapper = await mountEdit()
    const perRow = wrapper
      .findAll('.field-row')
      .map(row => row.findAll('[data-icon="pi pi-language"]').length)
    expect(perRow).toEqual([1, 1, 1])
  })

  it('opens one dialog holding every key the select field has', async () => {
    const wrapper = await mountEdit()
    await wrapper.findAll('.field-row')[0].find('[data-icon="pi pi-language"]').trigger('click')

    const cfg = useTranslatorStore().config
    expect(cfg.resource).toBe('compose:module-field/N1/M1/F-status')

    const set = await cfg.fetcher()
    expect(set.map(t => t.key)).toEqual(['label', 'meta.description.view', 'meta.options.new.text'])
  })

  it('keeps a bool field to its own keys', async () => {
    const wrapper = await mountEdit()
    await wrapper.findAll('.field-row')[1].find('[data-icon="pi pi-language"]').trigger('click')

    const set = await useTranslatorStore().config.fetcher()
    expect(set.map(t => t.key)).toEqual(['meta.bool.true.label'])
  })

  it('writes a saved option translation back onto the field', async () => {
    const wrapper = await mountEdit()
    const selectField = module_.fields[0]
    selectField.options.options = [{ value: 'new', text: 'New' }]

    await wrapper.findAll('.field-row')[0].find('[data-icon="pi pi-language"]').trigger('click')
    const cfg = useTranslatorStore().config

    composeAPI.moduleListTranslations.mockResolvedValueOnce([
      ...TRANSLATIONS.filter(t => t.key !== 'meta.options.new.text'),
      {
        resource: 'compose:module-field/N1/M1/F-status',
        key: 'meta.options.new.text',
        lang: 'en',
        message: 'Nouveau',
      },
    ])

    await cfg.updater([])
    expect(selectField.options.options[0].text).toBe('Nouveau')
  })
})
