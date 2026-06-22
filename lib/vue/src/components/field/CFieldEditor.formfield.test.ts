import { describe, it, expect, vi } from 'vitest'
import { defineComponent, reactive } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import { compose } from '@planetcrust/human-js'
import CFieldEditor from './CFieldEditor.vue'
import CFieldStringEditor from './editors/CFieldStringEditor.vue'

// Real String editor -> real PrimeVue InputText.
vi.mock('./registry', () => ({
  resolveFieldEditor: () => CFieldStringEditor,
}))

const mod = new compose.Module({
  fields: [{ name: 'tags', kind: 'String', isMulti: true }],
})

// Mirror Create.vue + RecordBlock: <Form> wrapping <FormField :name> around
// CFieldEditor, value read from record.values and written via setValue.
const Parent = defineComponent({
  components: { CFieldEditor },
  setup() {
    const record = reactive(new compose.Record(mod))
    const field = mod.fields.find(f => f.name === 'tags')
    const resolver = () => ({ errors: {} })
    const getFieldValue = () => {
      const v = record.values['tags']
      return v === undefined || v === null ? [] : v
    }
    const setFieldValue = (v: unknown) => record.setValue('tags', v as string[])
    return { record, field, resolver, getFieldValue, setFieldValue }
  },
  template: `
    <Form :resolver="resolver">
      <FormField name="tags" v-slot="{ invalid, error }">
        <CFieldEditor
          :field="field"
          :model-value="getFieldValue()"
          @update:model-value="setFieldValue($event)"
        />
      </FormField>
    </Form>
  `,
})

describe('multivalue inside @primevue/forms Form/FormField', () => {
  it('typing in one entry must not change the others', async () => {
    const wrapper = mount(Parent)
    await flushPromises()

    const addBtn = () => wrapper.findAll('button').at(-1)!
    await addBtn().trigger('click')
    await flushPromises()

    let inputs = wrapper.findAll('input')
    expect(inputs).toHaveLength(2)

    await inputs[0].setValue('hello')
    await flushPromises()

    inputs = wrapper.findAll('input')
    expect(inputs.map(i => (i.element as HTMLInputElement).value)).toEqual(['hello', ''])
  })
})
