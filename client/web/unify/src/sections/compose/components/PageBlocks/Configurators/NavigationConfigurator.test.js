import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'

const PAGES = [
  { pageID: 'P1', title: 'Target Page', handle: 'target', selfID: '0', moduleID: '0' },
  { pageID: 'P2', title: 'Sub One', handle: 'sub-one', selfID: 'P1', moduleID: '0' },
  { pageID: 'P9', title: 'Thing Record', handle: 'thing', selfID: '0', moduleID: 'M1' },
]

const stub = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CInputColorPicker: { name: 'CInputColorPicker', props: ['modelValue'], template: '<div />' },
    CInputToggleCard: {
      name: 'CInputToggleCard',
      props: ['modelValue', 'label'],
      template: '<div />',
    },
  },
  usePageStore: () => ({ set: PAGES, getByID: id => PAGES.find(p => p.pageID === id) }),
  usePageLayoutStore: () => ({ getByPageID: () => [] }),
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: {} }),
}))

import NavigationConfigurator from './NavigationConfigurator.vue'

const SelectStub = {
  name: 'Select',
  props: ['modelValue', 'options', 'placeholder', 'disabled'],
  template: '<div />',
}

function mountConfigurator(navItem) {
  const blockDraft = ref({ options: { navigationItems: [navItem] } })

  const wrapper = mount(NavigationConfigurator, {
    props: { namespace: { namespaceID: 'N1' }, page: { pageID: 'P0' } },
    global: {
      provide: { blockDraft },
      mocks: { $t: k => k },
      stubs: {
        Select: SelectStub,
        Fieldset: stub('Fieldset', ['legend']),
        Divider: stub('Divider'),
        Panel: stub('Panel', ['header', 'toggleable']),
        Button: stub('Button', ['label', 'icon']),
        InputText: { name: 'InputText', props: ['modelValue'], template: '<input />' },
        CFormGroup: stub('CFormGroup', ['label']),
        CInputExpression: stub('CInputExpression', ['modelValue']),
        CExpressionHint: stub('CExpressionHint'),
      },
    },
  })

  // The page chooser is the one carrying the choose-a-page placeholder.
  const pageSelect = wrapper
    .findAllComponents(SelectStub)
    .find(s => s.props('placeholder') === 'block.navigation.selectPage')

  return { blockDraft, wrapper, pageSelect }
}

const composeItem = (label = '') => ({
  type: 'compose',
  options: { enabled: true, item: { label } },
})

const item = draft => draft.value.options.navigationItems[0].options.item

describe('NavigationConfigurator compose page', () => {
  it('names an unnamed item after the page picked for it', async () => {
    const { blockDraft, pageSelect } = mountConfigurator(composeItem(''))

    await pageSelect.vm.$emit('update:modelValue', 'P1')

    expect(item(blockDraft).label).toBe('Target Page')
    expect(item(blockDraft).pageID).toBe('P1')
  })

  it('leaves a label the author typed alone', async () => {
    const { blockDraft, pageSelect } = mountConfigurator(composeItem('My own name'))

    await pageSelect.vm.$emit('update:modelValue', 'P1')

    expect(item(blockDraft).label).toBe('My own name')
  })

  it('drops the layout, which belonged to the page being replaced', async () => {
    const navItem = composeItem('Old')
    navItem.options.item.pageLayoutID = 'L7'
    const { blockDraft, pageSelect } = mountConfigurator(navItem)

    await pageSelect.vm.$emit('update:modelValue', 'P9')

    expect(item(blockDraft).pageLayoutID).toBe('')
  })

  it('offers the sub-pages toggle only for a page that has sub-pages', () => {
    const withKids = mountConfigurator({ ...composeItem('A'), options: { item: { pageID: 'P1' } } })
    const without = mountConfigurator({ ...composeItem('B'), options: { item: { pageID: 'P9' } } })

    expect(withKids.wrapper.findComponent({ name: 'CInputToggleCard' }).exists()).toBe(true)
    expect(without.wrapper.findComponent({ name: 'CInputToggleCard' }).exists()).toBe(false)
  })

  it('offers a layout chooser for a compose item', () => {
    const { wrapper } = mountConfigurator({
      ...composeItem('A'),
      options: { item: { pageID: 'P1' } },
    })

    const layoutSelect = wrapper
      .findAllComponents(SelectStub)
      .find(s => s.props('placeholder') === 'block.navigation.defaultLayout')

    expect(layoutSelect).toBeTruthy()
  })
})
