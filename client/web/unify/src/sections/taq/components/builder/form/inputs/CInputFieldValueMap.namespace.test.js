import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// The values map is the only namespace a Record field's editor can get: the TAQ
// builder is outside compose, so there is no `$namespace` to inject and the
// record select would sit disabled on "Select a namespace first".

const activities = {
  moduleID: '200',
  fields: [
    { name: 'opportunity', label: 'Opportunity', kind: 'Record', options: { moduleID: '300' } },
    { name: 'subject', label: 'Subject', kind: 'String', options: {} },
  ],
}

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@planetcrust/human-js/src/compose/types/module', () => ({ systemFields: [] }))

vi.mock('@planetcrust/human-vue', () => ({
  components: { CInputDelete: { name: 'CInputDelete', template: '<button />' } },
  useModuleStore: () => ({ findByID: vi.fn(async () => activities) }),
}))

vi.mock('@planetcrust/human-vue/src/components/field', () => ({
  CFieldEditor: {
    name: 'CFieldEditor',
    props: {
      field: { type: Object, default: () => ({}) },
      namespace: { type: Object, default: () => ({}) },
      modelValue: { type: [String, Array], default: '' },
      disabled: { type: Boolean, default: false },
      addLabel: { type: String, default: '' },
    },
    template: '<div />',
  },
}))

import CInputFieldValueMap from './CInputFieldValueMap.vue'

let wrapper

async function mountMap(props = {}) {
  wrapper = mount(CInputFieldValueMap, {
    props: {
      namespaceID: '100',
      moduleID: '200',
      modelValue: { opportunity: { value: '' } },
      ...props,
    },
    global: {
      stubs: { Select: true, Button: true, CReferenceChip: true },
    },
  })
  await flushPromises()
  return wrapper
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('field value map — namespace passthrough', () => {
  it("hands the step's namespace to the field editor", async () => {
    await mountMap()

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.exists()).toBe(true)
    expect(editor.props('namespace')).toEqual({ namespaceID: '100' })
  })

  it('resolves the row to the real Record field def, so the editor knows its module', async () => {
    await mountMap()

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.props('field').kind).toBe('Record')
    expect(editor.props('field').options.moduleID).toBe('300')
  })

  it('follows the namespace when the step is repointed', async () => {
    await mountMap()

    // Repointing the namespace clears the rows; the parent re-supplies them.
    await wrapper.setProps({ namespaceID: '101' })
    await flushPromises()
    await wrapper.setProps({ modelValue: { opportunity: { value: '' } } })
    await flushPromises()

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.props('namespace')).toEqual({ namespaceID: '101' })
  })

  it('keeps one namespace object across re-renders that do not move it', async () => {
    await mountMap()

    const before = wrapper.findComponent({ name: 'CFieldEditor' }).props('namespace')

    await wrapper.setProps({ modelValue: { subject: { value: 'hi' } } })
    await flushPromises()

    expect(wrapper.findComponent({ name: 'CFieldEditor' }).props('namespace')).toBe(before)
  })

  it('passes a multi-value array through to the editor untouched', async () => {
    await mountMap({ modelValue: { opportunity: { value: ['1', '2'] } } })

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    expect(editor.props('modelValue')).toEqual(['1', '2'])
  })

  it('emits a multi-value array back as the row value', async () => {
    await mountMap()

    const editor = wrapper.findComponent({ name: 'CFieldEditor' })
    await editor.vm.$emit('update:modelValue', ['1', '2'])

    const emitted = wrapper.emitted('update:modelValue')
    expect(emitted.at(-1)[0]).toEqual({ opportunity: { value: ['1', '2'] } })
  })
})
