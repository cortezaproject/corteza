import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive, h } from 'vue'

// A chart's palette comes from one of two places: the tables compiled into
// lib/js, or the schemes an admin defined for this instance and stored in
// ui.charts.colorSchemes. These pin the second half — that the picker offers
// them, that managing them is gated on settings.manage, and that saving one
// writes the setting and nothing else. The chart itself stays unsaved: adding a
// palette must not commit whatever else the author has in flight.

const route = reactive({
  name: 'admin.charts.edit',
  params: { slug: 'ns', chartID: 'C1' },
  query: {},
})

const router = { push: vi.fn(), replace: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

vi.mock('@planetcrust/human-js', () => {
  class Chart {
    constructor(opts = {}) {
      Object.assign(this, opts)
    }
  }
  return {
    compose: {
      Chart,
      GaugeChart: class extends Chart {},
      FunnelChart: class extends Chart {},
      RadarChart: class extends Chart {},
    },
    shared: { colorschemes: { tableau: { Tableau10: ['#a', '#b'] } } },
  }
})

let can

const chartStore = {
  findByID: vi.fn(() =>
    Promise.resolve({
      chartID: 'C1',
      name: 'Chart',
      handle: 'chart',
      namespaceID: 'N1',
      meta: {},
      canUpdateChart: true,
      config: { colorScheme: 'custom-1', reports: [{ moduleID: 'M1' }] },
    }),
  ),
  update: vi.fn(),
  create: vi.fn(),
}

vi.mock('@planetcrust/human-vue', () => ({
  useChartStore: () => chartStore,
  useModuleStore: () => ({ set: [], getByID: () => undefined }),
  useHistoryBack: () => vi.fn(),
  useDraftGuard: () => ({ capture: vi.fn(), markSaved: vi.fn() }),
  useRBACStore: () => ({ can: (...a) => can(...a) }),
  components: {
    CInputDelete: {
      name: 'CInputDelete',
      emits: ['confirm'],
      template: '<button data-testid="scheme-delete" @click="$emit(\'confirm\')" />',
    },
    CInputColorPicker: {
      name: 'CInputColorPicker',
      props: { modelValue: { type: String, default: '' } },
      template: '<span class="color-picker" />',
    },
  },
}))

vi.mock('../../../lib/charts', () => ({
  chartConstructor: c => c,
}))

vi.mock('../../../components/Chart/ChartRenderer.vue', () => ({
  default: {
    name: 'ChartRenderer',
    template: '<div />',
    methods: { updateChart: vi.fn() },
  },
}))

vi.mock('../../../components/Admin/Chart/ChartTranslator.vue', () => ({
  default: { name: 'ChartTranslator', template: '<div />' },
}))

vi.mock('../../../components/Chart/Report/index.js', () => ({
  GenericChart: { name: 'GenericChart', template: '<div />' },
  FunnelChart: { name: 'FunnelChart', template: '<div />' },
  GaugeChart: { name: 'GaugeChart', template: '<div />' },
  RadarChart: { name: 'RadarChart', template: '<div />' },
}))

import Edit from './Edit.vue'

// A container that is only in the way: it renders whatever it wraps.
const passthrough = name => ({ name, template: '<div><slot /></div>' })

// The picker, cut down to the two things these tests read off it — the option
// list it was handed, and the header slot the add control lives in.
const SelectStub = {
  name: 'Select',
  props: { options: { type: Array, default: () => [] }, modelValue: { default: undefined } },
  template: '<div class="select"><slot name="header" /></div>',
}

const DialogStub = {
  name: 'Dialog',
  props: { visible: { type: Boolean, default: false } },
  setup:
    (props, { slots }) =>
    () =>
      props.visible ? h('div', { class: 'dialog' }, [slots.default?.(), slots.footer?.()]) : null,
}

const ButtonStub = {
  name: 'Button',
  props: { label: { type: String, default: '' }, disabled: { type: Boolean, default: false } },
  emits: ['click'],
  template: '<button :disabled="disabled" @click="$emit(\'click\')">{{ label }}</button>',
}

const InputStub = {
  name: 'InputText',
  props: { modelValue: { type: String, default: '' } },
  emits: ['update:modelValue'],
  template:
    '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)">',
}

let settings
let settingsUpdate

function mountEditor() {
  return mount(Edit, {
    props: { namespace: { namespaceID: 'N1', slug: 'ns' } },
    global: {
      stubs: {
        teleport: true,
        Teleport: true,
        Form: passthrough('Form'),
        Panel: passthrough('Panel'),
        Card: passthrough('Card'),
        CFormGroup: passthrough('CFormGroup'),
        CEditorActions: passthrough('CEditorActions'),
        ButtonGroup: passthrough('ButtonGroup'),
        Select: SelectStub,
        Dialog: DialogStub,
        Button: ButtonStub,
        InputText: InputStub,
        Textarea: true,
        ToggleSwitch: true,
        ProgressSpinner: true,
        CPermissionsButton: true,
      },
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: {},
        $SystemAPI: { settingsUpdate },
        $Settings: settings,
        $toast: { toastSuccess: vi.fn(), toastWarning: vi.fn(), toastErrorHandler: () => vi.fn() },
        $Auth: { user: {} },
      },
    },
  })
}

// The editor holds its spinner for a beat after the fetch resolves, and the
// whole form is behind it. Every check below reads an element inside that form,
// so a wrapper handed back too early makes an absent control look like a hidden
// one — which is what three of these tests are meant to tell apart.
async function open() {
  const wrapper = mountEditor()
  await flushPromises()
  await new Promise(r => setTimeout(r, 350))
  await flushPromises()

  if (!wrapper.findComponent(SelectStub).exists()) {
    throw new Error('the editor never rendered its form')
  }

  return wrapper
}

const stored = () => settingsUpdate.mock.calls.at(-1)[0].values[0]

beforeEach(() => {
  can = () => true
  vi.clearAllMocks()

  let value = [{ id: 'custom-1', name: 'Brand', colors: ['#FF00AA', '#00FFD5'] }]
  settings = {
    get: (k, d) => (k === 'ui.charts.colorSchemes' ? value : d),
    fetch: vi.fn(() => Promise.resolve()),
  }
  settingsUpdate = vi.fn(({ values }) => {
    value = values[0].value
    return Promise.resolve(true)
  })
})

describe('the colour scheme picker', () => {
  it("offers the instance's own schemes ahead of the built-in tables", async () => {
    const wrapper = await open()
    const { options } = wrapper.findComponent(SelectStub).props()

    expect(options.map(o => o.id)).toEqual(['custom-1', 'tableau.Tableau10'])
    expect(options[0].colors).toEqual(['#FF00AA', '#00FFD5'])
  })

  it('offers to add one, and to edit the custom scheme in use', async () => {
    const wrapper = await open()

    expect(wrapper.find('[data-testid="chart-color-scheme-add"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="chart-color-scheme-edit"]').exists()).toBe(true)
  })

  it('offers neither without settings.manage', async () => {
    can = () => false
    const wrapper = await open()

    expect(wrapper.find('[data-testid="chart-color-scheme-add"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="chart-color-scheme-edit"]').exists()).toBe(false)
  })

  it('offers no pencil for a built-in scheme, which is not this instance to edit', async () => {
    chartStore.findByID.mockResolvedValueOnce({
      chartID: 'C1',
      name: 'Chart',
      namespaceID: 'N1',
      meta: {},
      canUpdateChart: true,
      config: { colorScheme: 'tableau.Tableau10', reports: [{ moduleID: 'M1' }] },
    })
    const wrapper = await open()

    expect(wrapper.find('[data-testid="chart-color-scheme-edit"]').exists()).toBe(false)
  })
})

describe('saving a custom scheme', () => {
  it('writes the new scheme to the setting and selects it', async () => {
    const wrapper = await open()

    await wrapper.find('[data-testid="chart-color-scheme-add"]').trigger('click')
    await wrapper.find('#colorSchemeName').setValue('Neon')
    await wrapper.find('[data-testid="chart-color-scheme-save"]').trigger('click')
    await flushPromises()

    const { name, value } = stored()
    expect(name).toBe('ui.charts.colorSchemes')
    expect(value.map(s => s.name)).toEqual(['Brand', 'Neon'])
    expect(wrapper.findComponent(SelectStub).props('modelValue')).toBe(value[1].id)
  })

  // The chart is the author's to save. Corteza committed it as a side effect of
  // adding a palette, which takes every other in-flight edit with it.
  it('leaves the chart itself unsaved', async () => {
    const wrapper = await open()

    await wrapper.find('[data-testid="chart-color-scheme-add"]').trigger('click')
    await wrapper.find('#colorSchemeName').setValue('Neon')
    await wrapper.find('[data-testid="chart-color-scheme-save"]').trigger('click')
    await flushPromises()

    expect(chartStore.update).not.toHaveBeenCalled()
    expect(chartStore.create).not.toHaveBeenCalled()
  })

  it('refuses to save a scheme with no name, which nothing could pick out', async () => {
    const wrapper = await open()

    await wrapper.find('[data-testid="chart-color-scheme-add"]').trigger('click')
    await wrapper.find('[data-testid="chart-color-scheme-save"]').trigger('click')
    await flushPromises()

    expect(settingsUpdate).not.toHaveBeenCalled()
  })

  it('replaces the scheme it was opened on rather than adding a second', async () => {
    const wrapper = await open()

    await wrapper.find('[data-testid="chart-color-scheme-edit"]').trigger('click')
    await wrapper.find('#colorSchemeName').setValue('Brand v2')
    await wrapper.find('[data-testid="chart-color-scheme-save"]').trigger('click')
    await flushPromises()

    expect(stored().value).toEqual([
      { id: 'custom-1', name: 'Brand v2', colors: ['#FF00AA', '#00FFD5'] },
    ])
  })
})

describe('deleting a custom scheme', () => {
  it('drops it from the setting and clears the chart that was using it', async () => {
    const wrapper = await open()

    await wrapper.find('[data-testid="chart-color-scheme-edit"]').trigger('click')
    await wrapper.find('[data-testid="scheme-delete"]').trigger('click')
    await flushPromises()

    expect(stored().value).toEqual([])
    expect(wrapper.findComponent(SelectStub).props('modelValue')).toBeUndefined()
  })
})
