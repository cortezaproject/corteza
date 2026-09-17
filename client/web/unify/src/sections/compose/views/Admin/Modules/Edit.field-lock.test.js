import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// Once a module has records, a saved field's name and kind are locked; fields
// not yet saved stay editable.

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
        return []
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

const field = (fieldID, name) => ({ fieldID, name, label: name, kind: 'String', options: {} })

const FormListStub = {
  name: 'CFormList',
  props: { modelValue: Array, columns: Array },
  template: `<div>
    <div v-for="(item, index) in modelValue" :key="index" class="field-row">
      <slot name="row" :item="item" :index="index" />
    </div>
  </div>`,
}

const passthrough = (name, slot = 'default') => [
  name,
  { name, template: `<div><slot name="${slot}" /></div>` },
]

// The tooltip text lands on the element, so a test can read what it says.
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
            'Checkbox',
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
          template: '<input class="name-input" :disabled="disabled" />',
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

// Name input, kind select, and the tooltips on their cells, for one row.
function rowState(row) {
  return {
    name: row.find('.name-input').attributes('disabled') !== undefined,
    kind: row.find('.kind-select').attributes('disabled') !== undefined,
    nameTip: row.find('.name-input').element.parentElement.dataset.tooltip,
    kindTip: row.find('.kind-select').element.closest('.input-group').dataset.tooltip,
  }
}

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    name: 'Probe',
    handle: 'probe',
    fields: [field('F1', 'amount'), field('0', 'draft')],
  }
})

describe('module field lock', () => {
  it('locks a saved field of a module with records and says why', async () => {
    module_.hasRecords = true
    const [saved] = (await mountEdit()).findAll('.field-row')

    expect(rowState(saved)).toEqual({
      name: true,
      kind: true,
      nameTip: 'field.disabled.lockedByRecords',
      kindTip: 'field.disabled.lockedByRecords',
    })
  })

  it('leaves a field that is not saved yet editable', async () => {
    module_.hasRecords = true
    const [, draft] = (await mountEdit()).findAll('.field-row')

    expect(rowState(draft)).toEqual({ name: false, kind: false, nameTip: '', kindTip: '' })
  })

  it('leaves every field editable on a module without records', async () => {
    const rows = (await mountEdit()).findAll('.field-row')

    expect(rows.map(rowState)).toEqual([
      { name: false, kind: false, nameTip: '', kindTip: '' },
      { name: false, kind: false, nameTip: '', kindTip: '' },
    ])
  })
})
