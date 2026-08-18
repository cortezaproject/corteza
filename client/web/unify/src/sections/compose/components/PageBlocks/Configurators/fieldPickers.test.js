import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { compose } from '@planetcrust/human-js'

// Block models cast moduleID through HumanID, so it has to look like one.
const MODULE_ID = '100001'

// One module with a field of every kind these configurators narrow to, so a
// picker handed the wrong `kinds` shows up as the wrong set, not an empty one.
const MODULE = {
  moduleID: MODULE_ID,
  name: 'Activity',
  fields: [
    { name: 'title', label: 'Title', kind: 'String' },
    { name: 'weight', label: 'Weight', kind: 'Number' },
    { name: 'due', label: 'Due', kind: 'DateTime' },
    { name: 'where', label: 'Where', kind: 'Geometry' },
    { name: 'doc', label: 'Doc', kind: 'File' },
    { name: 'parent', label: 'Parent', kind: 'Record' },
    { name: 'link', label: 'Link', kind: 'Url' },
  ],
  systemFields: () => [],
}

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CInputColorPicker: { name: 'CInputColorPicker', props: ['modelValue'], template: '<div />' },
  },
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

import RecordConfigurator from './RecordConfigurator.vue'
import RecordOrganizerConfigurator from './RecordOrganizerConfigurator.vue'
import CommentConfigurator from './CommentConfigurator.vue'
import IFrameConfigurator from './IFrameConfigurator.vue'

const passthrough = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

const FieldPicker = {
  name: 'CInputModuleField',
  props: [
    'modelValue',
    'module',
    'moduleID',
    'kinds',
    'excludeMulti',
    'valueKey',
    'includeSystem',
    'filter',
    'size',
  ],
  template: '<div />',
}

const ModulePicker = {
  name: 'CInputModule',
  props: ['modelValue', 'namespaceID', 'placeholder'],
  template: '<div />',
}

// The block model is the options schema — a draft made with the wrong kind
// silently drops the very options these pickers edit.
function mountConfigurator(Component, kind, options, page = { pageID: 'P0' }) {
  const blockDraft = ref(compose.PageBlockMaker({ kind, options }))

  return mount(Component, {
    props: { namespace: { namespaceID: 'N1' }, page },
    global: {
      provide: { blockDraft, $ComposeAPI: {} },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      stubs: {
        CInputModuleField: FieldPicker,
        CInputModule: ModulePicker,
        CFormGroup: passthrough('CFormGroup', ['label', 'description']),
        Fieldset: passthrough('Fieldset', ['legend']),
        Divider: passthrough('Divider'),
        Button: passthrough('Button', ['label', 'icon']),
        InputText: { name: 'InputText', props: ['modelValue'], template: '<input />' },
        Select: { name: 'Select', props: ['modelValue', 'options'], template: '<div />' },
        CInputExpression: passthrough('CInputExpression', ['modelValue', 'minLines', 'size']),
        CExpressionHint: passthrough('CExpressionHint'),
        CInputToggleCard: passthrough('CInputToggleCard', ['modelValue', 'label', 'description']),
        CFieldPicker: passthrough('CFieldPicker', ['allFields', 'modelValue']),
        CFormList: {
          name: 'CFormList',
          props: ['modelValue', 'columns', 'emptyMessage'],
          template:
            '<div><slot name="row" v-for="(item, index) in modelValue" :item="item" :index="index" /></div>',
        },
        Checkbox: { name: 'Checkbox', props: ['modelValue'], template: '<input />' },
      },
    },
  })
}

// Every picker in the wrapper, as { the value it edits → the kinds it offers }.
// An unbound `kinds` reads as 'all' — that is what the picker does with it.
const wiring = w =>
  Object.fromEntries(
    w
      .findAllComponents(FieldPicker)
      .map(p => [p.props('modelValue') ?? '', p.props('kinds') ?? 'all']),
  )

describe('configurator field pickers', () => {
  it('narrows each record organizer picker to what the option can hold', () => {
    const w = mountConfigurator(RecordOrganizerConfigurator, 'RecordOrganizer', {
      moduleID: MODULE_ID,
      labelField: 'a',
      descriptionField: 'b',
      positionField: 'c',
      groupField: 'd',
    })

    expect(wiring(w)).toEqual({ a: 'all', b: 'all', c: ['Number'], d: 'all' })
  })

  it('narrows each comment picker to the field kind it maps', () => {
    const w = mountConfigurator(CommentConfigurator, 'Comment', {
      moduleID: MODULE_ID,
      titleField: 'title',
      contentField: 'content',
      replyField: 'reply',
      referenceField: 'reference',
      attachmentField: 'attachment',
      reactionsField: 'reactions',
    })

    expect(wiring(w)).toEqual({
      title: ['String'],
      content: ['String'],
      reply: ['Record'],
      reference: ['Record'],
      attachment: ['File'],
      reactions: ['String'],
    })
  })

  it('offers only Url fields as an iframe source, and only on a record page', () => {
    const recordPage = mountConfigurator(
      IFrameConfigurator,
      'IFrame',
      { srcField: 'link' },
      { moduleID: MODULE_ID },
    )
    const plainPage = mountConfigurator(
      IFrameConfigurator,
      'IFrame',
      { srcField: '' },
      { moduleID: '0' },
    )

    expect(wiring(recordPage)).toEqual({ link: ['Url'] })
    expect(plainPage.findAllComponents(FieldPicker)).toHaveLength(0)
  })

  it('picks a conditioned field by ID, from the fields the block draws', () => {
    const configured = mountConfigurator(
      RecordConfigurator,
      'Record',
      {
        fields: ['title', 'weight'],
        fieldConditions: [{ field: 'F1', condition: 'true' }],
      },
      { moduleID: MODULE_ID },
    )

    const picker = configured.findAllComponents(FieldPicker).at(-1)
    expect(picker.props('valueKey')).toBe('fieldID')
    expect(MODULE.fields.filter(picker.props('filter')).map(f => f.name)).toEqual([
      'title',
      'weight',
    ])
  })

  it('offers every module field when the block configures none', () => {
    const all = mountConfigurator(
      RecordConfigurator,
      'Record',
      { fields: [], fieldConditions: [{ field: '', condition: '' }] },
      { moduleID: MODULE_ID },
    )

    expect(all.findAllComponents(FieldPicker).at(-1).props('filter')).toBe(null)
  })

  it('gives every picker a module to read its fields from', () => {
    const w = mountConfigurator(CommentConfigurator, 'Comment', {
      moduleID: MODULE_ID,
      titleField: 'title',
    })

    for (const p of w.findAllComponents(FieldPicker)) {
      expect(p.props('module') || p.props('moduleID')).toBeTruthy()
    }
  })
})
