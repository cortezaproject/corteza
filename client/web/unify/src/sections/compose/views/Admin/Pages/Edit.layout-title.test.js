import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { reactive } from 'vue'

// A layout's `meta.title` is edited in two places — the list row and the config
// dialog — and with `config.useTitle` on it is the page's own title, read as a
// `${}` template. Both places author it the same way.

const route = reactive({
  name: 'admin.pages.edit',
  params: { slug: 'ns', pageID: 'P1' },
  query: {},
})

const router = { push: vi.fn(), replace: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
  onBeforeRouteLeave: () => {},
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('primevue/useconfirm', () => ({ useConfirm: () => ({ require: vi.fn() }) }))

const SCOPE = [{ name: 'recordID' }]

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: SCOPE, exprScope: SCOPE }),
}))

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

// Renders each item through the `#row` slot, which is where the title input is.
const FormListStub = {
  name: 'CFormList',
  props: ['modelValue', 'columns', 'emptyMessage', 'draggable', 'hideRemove'],
  template: `<div>
    <div v-for="(item, index) in modelValue" :key="index" class="layout-row">
      <slot name="row" :item="item" :index="index" />
    </div>
  </div>`,
}

const ExpressionStub = {
  name: 'CInputExpression',
  props: ['modelValue', 'dialect', 'scope', 'minLines', 'invalid', 'placeholder'],
  methods: {
    insert(text) {
      this.$emit('update:modelValue', (this.modelValue || '') + text)
    },
  },
  template: '<div class="expr" />',
}

const HintStub = { name: 'CExpressionHint', props: ['scope'], template: '<div class="hint" />' }

const InputTextStub = {
  name: 'InputText',
  props: ['modelValue', 'invalid'],
  template: '<input class="plain" />',
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
    'InputGroup',
    'InputGroupAddon',
    'InputNumber',
    'Message',
    'ProgressSpinner',
    'Select',
    'Textarea',
    'ToggleSwitch',
  ].map(name => [name, true]),
)

let page
let wrapper

const layoutWith = (useTitle, title = 'Layout name') => [
  {
    pageLayoutID: 'L1',
    pageID: 'P1',
    meta: { title, description: '' },
    config: { useTitle, visibility: { expression: '', roles: [] } },
    blocks: [],
  },
]

async function mountEdit() {
  wrapper = mount(Edit, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      stubs: {
        teleport: true,
        Teleport: true,
        ...GLOBAL_COMPONENTS,
        // Render their slots: the layouts list lives inside both.
        Form: {
          name: 'Form',
          props: ['resolver', 'initialValues'],
          template: '<div><slot /></div>',
        },
        Panel: { name: 'Panel', props: ['header', 'toggleable'], template: '<div><slot /></div>' },
        CFormList: FormListStub,
        CInputExpression: ExpressionStub,
        CExpressionHint: HintStub,
        InputText: InputTextStub,
      },
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
      renderStubDefaultSlot: false,
    },
  })
  await flushPromises()
  return wrapper
}

beforeEach(() => {
  createTestPinia({ $ComposeAPI: {}, $SystemAPI: {} })
  page = { pageID: 'P1', title: 'A page', moduleID: 'M1', blocks: [], config: {} }
  layouts = layoutWith(false)
  vi.clearAllMocks()
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

const rowInputs = w => w.findAll('.layout-row .expr')
const rowPlain = w => w.findAll('.layout-row .plain')

describe('layout list title input', () => {
  it('authors the title as a template once the layout titles the page', async () => {
    layouts = layoutWith(true, '${record.values.name}')
    await mountEdit()

    expect(rowInputs(wrapper)).toHaveLength(1)
    expect(rowPlain(wrapper).length).toBeLessThan(2) // the handle input only

    const expr = wrapper.findAllComponents(ExpressionStub)[0]
    expect(expr.props('dialect')).toBe('interpolation')
    expect(expr.props('minLines')).toBe(1)
    expect(expr.props('scope')).toEqual(SCOPE)
    expect(expr.props('modelValue')).toBe('${record.values.name}')
  })

  it('leaves it a plain input while the layout does not title the page', async () => {
    layouts = layoutWith(false)
    await mountEdit()

    expect(rowInputs(wrapper)).toHaveLength(0)
  })

  it('offers the hint only once some layout titles the page', async () => {
    layouts = layoutWith(false)
    await mountEdit()
    expect(wrapper.findComponent(HintStub).exists()).toBe(false)

    wrapper.unmount()
    layouts = layoutWith(true)
    await mountEdit()
    expect(wrapper.findComponent(HintStub).exists()).toBe(true)
  })

  it('marks the layout changed when the template is edited', async () => {
    layouts = layoutWith(true)
    await mountEdit()

    await wrapper.findAllComponents(ExpressionStub)[0].vm.$emit('update:modelValue', 'x')
    expect(wrapper.vm.layouts[0]._updated).toBe(true)
  })

  it('inserts a hint chip into the row that was last focused', async () => {
    layouts = [...layoutWith(true, 'one'), { ...layoutWith(true, 'two')[0], pageLayoutID: 'L2' }]
    await mountEdit()

    await wrapper.findAll('.layout-row .expr')[1].trigger('focusin')
    await wrapper.findComponent(HintStub).vm.$emit('insert', '${recordID}')

    expect(wrapper.vm.layouts[1].meta.title).toBe('two${recordID}')
    expect(wrapper.vm.layouts[0].meta.title).toBe('one')
  })
})
