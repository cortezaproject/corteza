import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref } from 'vue'

// A condition compared against a Record field renders a record select, which
// needs a namespace to list against. The builder is outside compose, so there
// is no `$namespace` to inject — the namespace of the upstream step's module is
// the one in play.

const activities = {
  moduleID: '200',
  fields: [
    { name: 'opportunity', label: 'Opportunity', kind: 'Record', options: { moduleID: '300' } },
    { name: 'subject', label: 'Subject', kind: 'String', options: {} },
  ],
}

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => ({ findByID: vi.fn(async () => activities) }),
}))

vi.mock('@planetcrust/human-vue/src/components/field', () => ({
  CFieldEditor: {
    name: 'CFieldEditor',
    props: {
      field: { type: Object, default: () => ({}) },
      namespace: { type: Object, default: () => ({}) },
      modelValue: { type: [String, Array], default: '' },
    },
    template: '<div />',
  },
}))

vi.mock('../form/CReferenceChip.vue', () => ({
  default: { name: 'CReferenceChip', template: '<div />' },
}))

import ConditionRow from './ConditionRow.vue'

const upstream = [
  {
    handle: 'create_activity',
    results: [
      {
        sourceName: 'record',
        types: ['ComposeRecord'],
        namespaceID: '100',
        moduleID: '200',
      },
    ],
  },
]

let wrapper

async function mountRow(symbol) {
  wrapper = mount(ConditionRow, {
    props: {
      node: {
        ref: 'eq',
        args: [
          { symbol, meta: { scope: 'create_activity' } },
          { value: { '@type': 'String', '@value': '' } },
        ],
      },
      edgeId: 'e1',
      rowIndex: 0,
    },
    global: {
      provide: {
        'taq-upstream-results': ref(upstream),
        $ComposeAPI: {},
      },
      stubs: { Select: true, Button: true, CInputExpression: true, InputText: true },
    },
  })
  await flushPromises()
  await flushPromises()
  return wrapper
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('condition row — namespace passthrough', () => {
  it("hands the upstream module's namespace to the field editor", async () => {
    await mountRow('record.values.opportunity')

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.exists()).toBe(true)
    expect(editor.props('field').kind).toBe('Record')
    expect(editor.props('namespace')).toEqual({ namespaceID: '100' })
  })

  it('leaves the namespace empty when the field is not resolved from a module', async () => {
    await mountRow('record.recordID')

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.props('field').kind).toBe('ID')
    expect(editor.props('namespace')).toEqual({})
  })

  it('clears the namespace when the variable is repointed off the module field', async () => {
    await mountRow('record.values.opportunity')
    expect(wrapper.findComponent({ name: 'CFieldEditor' }).props('namespace')).toEqual({
      namespaceID: '100',
    })

    await wrapper.setProps({
      node: {
        ref: 'eq',
        args: [
          { symbol: 'record.createdAt', meta: { scope: 'create_activity' } },
          { value: { '@type': 'String', '@value': '' } },
        ],
      },
    })
    await flushPromises()
    await flushPromises()

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.props('field').kind).toBe('DateTime')
    expect(editor.props('namespace')).toEqual({})
  })
})
