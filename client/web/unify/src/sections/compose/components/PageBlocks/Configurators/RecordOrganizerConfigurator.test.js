import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref, nextTick } from 'vue'
import { compose } from '@planetcrust/human-js'

// Block models cast moduleID through HumanID, so it has to look like one.
const MODULE_ID = '100001'

const MODULE = {
  moduleID: MODULE_ID,
  name: 'Activity',
  fields: [
    { name: 'title', label: 'Title', kind: 'String' },
    {
      name: 'status',
      label: 'Status',
      kind: 'Select',
      options: {
        selectType: 'multiple',
        options: [
          { value: 'todo', text: 'To do' },
          { value: 'done', text: 'Done' },
        ],
      },
    },
    { name: 'due', label: 'Due', kind: 'DateTime' },
  ],
  systemFields: () => [],
}

// Hoisted: vi.mock's factory is lifted above the file, so the stub it hands
// back has to exist by then.
const { FieldEditor } = vi.hoisted(() => ({
  FieldEditor: {
    name: 'CFieldEditor',
    props: ['field', 'modelValue', 'namespace'],
    template: '<div />',
  },
}))

vi.mock('@planetcrust/human-vue', () => ({
  components: { CFieldEditor: FieldEditor },
  useModuleStore: () => ({
    set: [MODULE],
    getByID: id => (id === MODULE_ID ? MODULE : undefined),
    findByID: () => Promise.resolve(MODULE),
  }),
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: {}, queryFields: [] }),
}))

import RecordOrganizerConfigurator from './RecordOrganizerConfigurator.vue'

const passthrough = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

const FieldPicker = {
  name: 'CInputModuleField',
  // excludeMulti typed, the way the real picker declares it: an array prop
  // reads a valueless attribute as '' and never as true.
  props: {
    modelValue: {},
    module: {},
    moduleID: {},
    kinds: {},
    excludeMulti: { type: Boolean, default: false },
    placeholder: {},
  },
  emits: ['update:modelValue'],
  template: '<div />',
}

// Renders the legend as an attribute so a test can say which section a control
// sits in — the point of the key field/key value pairing.
const FieldsetStub = {
  name: 'Fieldset',
  props: ['legend'],
  template: '<section :data-legend="legend"><slot /></section>',
}

function mountConfigurator(options) {
  const blockDraft = ref(compose.PageBlockMaker({ kind: 'RecordOrganizer', options }))

  const w = mount(RecordOrganizerConfigurator, {
    props: { namespace: { namespaceID: 'N1' }, page: { pageID: 'P0' } },
    global: {
      provide: { blockDraft, $ComposeAPI: {} },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      stubs: {
        CInputModuleField: FieldPicker,
        CInputModule: passthrough('CInputModule', ['modelValue', 'namespaceID', 'placeholder']),
        CFormGroup: passthrough('CFormGroup', ['label', 'description']),
        Fieldset: FieldsetStub,
        Divider: passthrough('Divider'),
        InputText: {
          name: 'InputText',
          props: { modelValue: {}, disabled: { type: Boolean, default: false } },
          template: '<input />',
        },
        Select: { name: 'Select', props: ['modelValue', 'options'], template: '<div />' },
        CInputExpression: passthrough('CInputExpression', ['modelValue']),
        CExpressionHint: passthrough('CExpressionHint'),
      },
    },
  })

  return { w, blockDraft }
}

const section = (w, legend) => w.find(`[data-legend="${legend}"]`)

// The picker bound to a given option value, so a test names it the way the
// configurator does rather than by position.
const pickerFor = (root, value) =>
  root.findAllComponents(FieldPicker).find(p => p.props('modelValue') === value)

describe('record organizer configurator', () => {
  it('keeps the key field and key value together under behaviour', () => {
    const { w } = mountConfigurator({ moduleID: MODULE_ID, groupField: 'status', group: 'todo' })

    const behaviour = section(w, 'block.recordOrganizer.behaviour')
    const cardFields = section(w, 'block.recordOrganizer.cardFields')

    expect(pickerFor(behaviour, 'status')).toBeTruthy()
    expect(behaviour.findComponent(FieldEditor).exists()).toBe(true)
    expect(pickerFor(cardFields, 'status')).toBeUndefined()
    expect(cardFields.findComponent(FieldEditor).exists()).toBe(false)
  })

  it('separates the controls in every section', () => {
    const { w } = mountConfigurator({ moduleID: MODULE_ID })

    for (const legend of ['block.recordOrganizer.cardFields', 'block.recordOrganizer.behaviour']) {
      const wrapper = section(w, legend).find('div')
      expect(wrapper.classes()).toContain('gap-3')
    }
  })

  it('edits the key value with the key field own editor, single-valued', () => {
    const { w } = mountConfigurator({ moduleID: MODULE_ID, groupField: 'status', group: 'todo' })

    const editor = w.findComponent(FieldEditor)
    const field = editor.props('field')

    expect(editor.props('modelValue')).toBe('todo')
    expect(field.kind).toBe('Select')
    expect(field.options.options).toHaveLength(2)
    // One column holds one value; the server refuses a multi-value key field.
    expect(field.isMulti).toBe(false)
    expect(field.options.selectType).toBe('default')
  })

  it('writes what the editor emits back as the key value', async () => {
    const { w, blockDraft } = mountConfigurator({
      moduleID: MODULE_ID,
      groupField: 'status',
      group: 'todo',
    })

    await w.findComponent(FieldEditor).vm.$emit('update:modelValue', 'done')
    expect(blockDraft.value.options.group).toBe('done')

    // An editor that hands back an array contributes its first entry: the
    // block sends the key value as a string.
    await w.findComponent(FieldEditor).vm.$emit('update:modelValue', ['todo'])
    expect(blockDraft.value.options.group).toBe('todo')
  })

  it('offers a plain disabled input until a key field is picked', () => {
    const { w } = mountConfigurator({ moduleID: MODULE_ID })

    expect(w.findComponent(FieldEditor).exists()).toBe(false)
    const behaviour = section(w, 'block.recordOrganizer.behaviour')
    expect(behaviour.findComponent({ name: 'InputText' }).props('disabled')).toBe(true)
  })

  it('clears the key value when the key field changes', async () => {
    const { w, blockDraft } = mountConfigurator({
      moduleID: MODULE_ID,
      groupField: 'status',
      group: 'todo',
    })

    await pickerFor(w, 'status').vm.$emit('update:modelValue', 'due')
    await nextTick()

    expect(blockDraft.value.options.groupField).toBe('due')
    expect(blockDraft.value.options.group).toBe('')
  })

  it('keeps a saved key value when the block is opened', async () => {
    const { blockDraft } = mountConfigurator({
      moduleID: MODULE_ID,
      groupField: 'status',
      group: 'todo',
    })

    await nextTick()
    expect(blockDraft.value.options.group).toBe('todo')
  })

  it('does not offer multi-value fields as the key or sort field', () => {
    const { w } = mountConfigurator({
      moduleID: MODULE_ID,
      groupField: 'status',
      positionField: 'weight',
    })

    expect(pickerFor(w, 'status').props('excludeMulti')).toBe(true)
    expect(pickerFor(w, 'weight').props('excludeMulti')).toBe(true)
  })
})
