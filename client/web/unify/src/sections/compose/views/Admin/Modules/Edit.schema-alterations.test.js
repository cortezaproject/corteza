import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// A module with issues loads its pending schema alterations by resource ident.

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
let systemAPI

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
        $SystemAPI: systemAPI,
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
        DalSchemaAlterations: {
          name: 'DalSchemaAlterations',
          props: ['batch', 'modal', 'module'],
          template: '<div />',
        },
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

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    name: 'Probe',
    handle: 'probe',
    fields: [],
    issues: [{ kind: 'model', issue: 'requires schema alterations', meta: { batchID: 'B1' } }],
  }
  systemAPI = {
    dalSchemaAlterationList: vi.fn(() =>
      Promise.resolve({ set: [{ alterationID: 'A1', batchID: 'B1' }] }),
    ),
  }
})

describe('module schema alterations', () => {
  it('asks for the alterations of this module by its resource ident', async () => {
    await mountEdit()

    expect(systemAPI.dalSchemaAlterationList).toHaveBeenCalledWith({
      resource: ['corteza::compose:module/N1/M1'],
    })
  })

  it('hands the found batch to the alterations dialog', async () => {
    const wrapper = await mountEdit()

    expect(wrapper.findComponent({ name: 'DalSchemaAlterations' }).props('batch')).toEqual(['B1'])
  })
})
