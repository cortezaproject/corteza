import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { reactive, ref } from 'vue'

// The layout list's second column EDITS the rule the public views apply: layouts
// are tried in weight order and the first whose condition and roles both pass is
// the one shown (composables/usePageVisibility.ts). Both halves bind straight
// into the layout, the way the title does.

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
    'CExpressionHint',
    'InputText',
    'CInputModuleField',
    'PageTranslator',
  ].map(name => [name, true]),
)

// Both halves of the cell are two-way bound, so their stubs have to emit.
const ExpressionStub = {
  name: 'CInputExpression',
  props: ['modelValue', 'dialect', 'scope', 'minLines', 'size', 'placeholder', 'invalid'],
  emits: ['update:modelValue'],
  template:
    '<input class="expr" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
}

const RoleStub = {
  name: 'CInputRole',
  props: ['modelValue', 'placeholder', 'multiple', 'size'],
  emits: ['update:modelValue'],
  template: '<div class="roles">{{ (modelValue || []).join(\',\') }}</div>',
}

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
        CInputExpression: ExpressionStub,
        CInputRole: RoleStub,
      },
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
  await flushPromises()
  return wrapper
}

// Condition and roles are columns of their own, between the title and the
// actions; each row cell carries its own hook.
const condition = w => w.find('[data-layout-condition] .expr')
const roles = w => w.find('[data-layout-roles]')

beforeEach(() => {
  page = { pageID: 'P1', namespaceID: 'N1', title: 'Account', handle: 'account', blocks: [] }
})

describe('layout condition and roles columns', () => {
  it('gives the condition and the roles a column each', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()

    expect(condition(w).exists()).toBe(true)
    expect(roles(w).exists()).toBe(true)
    // Neither may sit inside the other's column, nor inside the title's.
    expect(w.find('[data-layout-condition] [data-layout-roles]').exists()).toBe(false)
    expect(w.find('[data-layout-roles] .expr').exists()).toBe(false)
    expect(w.find('[data-layout-title] [data-layout-condition]').exists()).toBe(false)
  })

  it('carries the stored condition and roles into them', async () => {
    layouts = layoutWith({ expression: 'screen.width < 1024', roles: ['R1', 'R2'] })
    const w = await mountEdit()
    expect(condition(w).element.value).toBe('screen.width < 1024')
    expect(roles(w).text()).toBe('R1,R2')
  })

  it('authors the condition as an expr expression', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    const expr = w.find('[data-layout-condition]').findComponent(ExpressionStub)
    expect(expr.props('dialect')).toBe('expr')
    expect(expr.props('minLines')).toBe(1)
    expect(expr.props('size')).toBe('small')
  })

  it('writes an edited condition back and marks the layout changed', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    await condition(w).setValue('isEdit')
    expect(w.vm.layouts[0].config.visibility.expression).toBe('isEdit')
    expect(w.vm.layouts[0]._updated).toBe(true)
  })

  it('writes edited roles back and marks the layout changed', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    await w.findComponent(RoleStub).vm.$emit('update:modelValue', ['R7'])
    expect(w.vm.layouts[0].config.visibility.roles).toEqual(['R7'])
    expect(w.vm.layouts[0]._updated).toBe(true)
  })
})

describe('adding a layout', () => {
  it('opens the dialog and leaves the list alone until Save', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    expect(w.vm.layouts).toHaveLength(1)

    w.vm.addLayout()
    await flushPromises()
    expect(w.vm.layoutConfigVisible).toBe(true)
    expect(w.vm.layouts).toHaveLength(1)
  })

  it('titles the dialog as a create, not a configure', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    w.vm.addLayout()
    expect(w.vm.layoutConfigTitle).toBe('page.page-layout.create')

    w.vm.onLayoutConfigClose()
    w.vm.openLayoutConfig(w.vm.layouts[0])
    expect(w.vm.layoutConfigTitle).toBe('page.page-layout.configure')
  })

  it('adds the row once the dialog is saved', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()

    w.vm.addLayout()
    w.vm.configLayout.meta.title = 'Mobile'
    w.vm.saveLayoutConfig()
    await flushPromises()

    expect(w.vm.layouts).toHaveLength(2)
    expect(w.vm.layouts[1].meta.title).toBe('Mobile')
    expect(w.vm.layoutConfigVisible).toBe(false)
  })

  it('adds nothing when the dialog is dismissed', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()

    w.vm.addLayout()
    w.vm.onLayoutConfigClose()
    await flushPromises()

    expect(w.vm.layouts).toHaveLength(1)
  })

  it('gives a new layout the visibility shape both editors bind into', async () => {
    layouts = layoutWith({ expression: '', roles: [] })
    const w = await mountEdit()
    w.vm.addLayout()
    expect(w.vm.configLayout.config.visibility).toEqual({ expression: '', roles: [] })
  })
})
