import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// System fields sit in the field list as a row like any other, told apart by
// its disabled controls and its lock, not by a colour laid over the whole row.
// Required and multi are not theirs to state, so those cells stay empty.

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
    showTranslatorButton: ref(false),
    currentLanguage: ref('en'),
    resourceTranslationsEnabled: ref(false),
    canManageResourceTranslations: ref(false),
  }),
}))

let module_
let sysFields

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => ({
    findByID: vi.fn(() => Promise.resolve(module_)),
    update: vi.fn(),
    create: vi.fn(),
    delete: vi.fn(),
  }),
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
        this.fields = []
        this.issues = []
        this.config = {}
        Object.assign(this, m)
      }
      systemFields() {
        return sysFields
      }
    },
    ModuleFieldMaker: ({ kind }) => ({ kind, cap: { multi: true, required: true } }),
    ModuleFieldString: class {
      constructor() {
        this.fieldID = '0'
        this.kind = 'String'
        this.name = ''
        this.label = ''
        this.options = {}
      }
    },
  },
}))

import Edit from './Edit.vue'

// Renders the footer slot the system rows live in, with the same grid style
// argument the real list passes.
const FormListStub = {
  name: 'CFormList',
  props: { modelValue: Array, columns: Array },
  template: `<div>
    <div v-for="(item, index) in modelValue" :key="index" class="field-row">
      <slot name="row" :item="item" :index="index" />
    </div>
    <slot name="footer" :grid-style="{}" />
  </div>`,
}

const passthrough = (name, slot = 'default') => [
  name,
  { name, template: `<div><slot name="${slot}" /></div>` },
]

const tooltip = (el, { value }) => {
  el.dataset.tooltip = value || ''
}

async function mountEdit() {
  const wrapper = mount(Edit, {
    props: { namespace: { namespaceID: 'N1', canManageNamespace: false } },
    global: {
      plugins: [createPinia()],
      provide: {
        $ComposeAPI: {},
        $SystemAPI: {},
        $Settings: { get: () => false },
        $toast: { toastDanger: vi.fn(), toastSuccess: vi.fn() },
      },
      mocks: { $t: k => k },
      stubs: {
        teleport: true,
        Teleport: true,
        ...Object.fromEntries(
          [
            'Button',
            'ButtonGroup',
            'Dialog',
            'Divider',
            'Menu',
            'Message',
            'ProgressSpinner',
            'Tab',
            'TabList',
          ].map(n => [n, true]),
        ),
        ...Object.fromEntries([
          passthrough('Card', 'content'),
          passthrough('Tabs'),
          passthrough('TabPanels'),
          passthrough('TabPanel'),
          passthrough('InputGroupAddon'),
        ]),
        InputGroup: { name: 'InputGroup', template: '<div class="input-group"><slot /></div>' },
        InputText: {
          props: ['modelValue', 'size', 'disabled'],
          template: '<input class="text-input" :value="modelValue" :disabled="disabled" />',
        },
        Checkbox: {
          props: ['modelValue', 'binary', 'disabled'],
          template:
            '<input type="checkbox" class="flag-box" :checked="modelValue" :disabled="disabled" />',
        },
        Select: {
          props: ['modelValue', 'options', 'size', 'disabled'],
          template: '<select class="kind-select" :disabled="disabled" />',
        },
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
      directives: { tooltip },
    },
  })
  await flushPromises()
  return wrapper
}

const systemRows = wrapper => wrapper.findAll('.system-field-row')

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    name: 'Probe',
    handle: 'probe',
    fields: [{ fieldID: 'F1', name: 'amount', label: 'Amount', kind: 'String', options: {} }],
  }
  sysFields = [
    { name: 'ownedBy', label: 'Owner', kind: 'User', isRequired: true, isMulti: false },
    { name: 'createdAt', label: 'Created', kind: 'DateTime', isRequired: false, isMulti: false },
  ]
})

describe('module system field rows', () => {
  it('lays no highlight over the row', async () => {
    const rows = systemRows(await mountEdit())

    expect(rows).toHaveLength(2)
    rows.forEach(row => {
      expect(row.classes()).not.toContain('bg-highlight')
      expect(row.classes()).not.toContain('text-muted-color')
    })
  })

  it('disables every control in the row', async () => {
    const [row] = systemRows(await mountEdit())

    const inputs = row.findAll('.text-input')
    expect(inputs.map(i => i.element.value)).toEqual(['ownedBy', 'Owner', 'User'])
    inputs.forEach(i => expect(i.attributes('disabled')).toBeDefined())
  })

  it('leaves the required and multi cells empty', async () => {
    systemRows(await mountEdit()).forEach(row => {
      expect(row.find('.flag-box').exists()).toBe(false)
      expect(row.find('i.pi-check').exists()).toBe(false)
      expect(row.find('i.pi-minus').exists()).toBe(false)
    })
  })
})
