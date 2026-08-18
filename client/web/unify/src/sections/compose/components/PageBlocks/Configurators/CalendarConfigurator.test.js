import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { compose } from '@planetcrust/human-js'

const MODULES = [{ moduleID: 'M1', name: 'Activity', fields: [] }]

const stub = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CInputColorPicker: { name: 'CInputColorPicker', props: ['modelValue'], template: '<div />' },
  },
  useModuleStore: () => ({
    set: MODULES,
    getByID: id => MODULES.find(m => m.moduleID === id),
  }),
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: {} }),
}))

import CalendarConfigurator from './CalendarConfigurator.vue'

const SelectStub = {
  name: 'Select',
  props: ['modelValue', 'options', 'placeholder', 'disabled'],
  template: '<div />',
}

// The feed's module and field mapping go through the shared pickers.
const PickerStub = {
  name: 'CInputModule',
  props: ['modelValue', 'placeholder', 'disabled'],
  template: '<div />',
}

const FieldPickerStub = {
  name: 'CInputModuleField',
  props: ['modelValue', 'moduleID', 'kinds', 'placeholder', 'disabled'],
  template: '<div />',
}

function mountConfigurator(feed = {}) {
  const blockDraft = ref(
    new compose.PageBlockCalendar({
      kind: 'Calendar',
      options: { feeds: [compose.PageBlockCalendar.makeFeed(feed)] },
    }),
  )

  const wrapper = mount(CalendarConfigurator, {
    props: { namespace: { namespaceID: 'N1' }, page: { pageID: 'P0' } },
    global: {
      provide: { blockDraft },
      mocks: { $t: k => k },
      stubs: {
        Select: SelectStub,
        Fieldset: stub('Fieldset', ['legend']),
        Divider: stub('Divider'),
        Button: stub('Button', ['label', 'icon']),
        Checkbox: { name: 'Checkbox', props: ['modelValue'], template: '<input />' },
        CFormGroup: stub('CFormGroup', ['label']),
        CInputModule: PickerStub,
        CInputModuleField: FieldPickerStub,
        CInputExpression: stub('CInputExpression', ['modelValue']),
        CExpressionHint: stub('CExpressionHint'),
      },
    },
  })

  const sourceSelect = wrapper
    .findAllComponents(SelectStub)
    .find(s => (s.props('options') || []).some(o => o.value === 'system:reminder'))

  return { blockDraft, wrapper, sourceSelect }
}

const draftFeed = draft => draft.value.options.feeds[0]

// The block is re-made from its draft on every open and every save
// (Builder.vue editBlock/commitEditingBlock), so a feed field the model does
// not declare does not survive being chosen.
const roundTrip = draft =>
  compose.PageBlockMaker(JSON.parse(JSON.stringify(draft.value))).options.feeds[0]

describe('CalendarConfigurator event source', () => {
  it('offers every resource the model knows, labelled from its own key', () => {
    const { sourceSelect } = mountConfigurator()

    expect(sourceSelect.props('options')).toEqual([
      { value: 'compose:record', label: 'block.calendar.recordFeed.optionLabel' },
      { value: 'system:reminder', label: 'block.calendar.reminderFeed.optionLabel' },
    ])
  })

  it('starts a new feed on records', () => {
    const { sourceSelect } = mountConfigurator()

    expect(sourceSelect.props('modelValue')).toBe('compose:record')
  })

  it('keeps a reminder source across the block round-trip', async () => {
    const { blockDraft, sourceSelect } = mountConfigurator()

    await sourceSelect.vm.$emit('update:modelValue', 'system:reminder')

    expect(draftFeed(blockDraft).resource).toBe('system:reminder')
    expect(roundTrip(blockDraft).resource).toBe('system:reminder')
  })

  it('shows a stored reminder source as reminders, not records', () => {
    const { sourceSelect } = mountConfigurator({ resource: 'system:reminder' })

    expect(sourceSelect.props('modelValue')).toBe('system:reminder')
  })

  it('asks for a module and field mapping only for a record feed', () => {
    const record = mountConfigurator()
    const reminder = mountConfigurator({ resource: 'system:reminder' })

    const placeholders = w =>
      [...w.findAllComponents(SelectStub), ...w.findAllComponents(PickerStub)]
        .map(s => s.props('placeholder'))
        .filter(Boolean)

    expect(placeholders(record.wrapper)).toContain('block.calendar.recordFeed.modulePlaceholder')
    expect(placeholders(reminder.wrapper)).not.toContain(
      'block.calendar.recordFeed.modulePlaceholder',
    )
  })
})
