import { describe, it, expect } from 'vitest'
import { defineComponent, reactive } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import CFieldNumberEditor from './CFieldNumberEditor.vue'

const field = { name: 'amount', kind: 'Number', options: {} }

// Mirror how every caller drives the editor: the record holds the value, the
// editor's emit writes it back, and the new value is handed straight back down.
const Parent = defineComponent({
  components: { CFieldNumberEditor },
  setup() {
    const record = reactive({ values: { amount: '1' } as Record<string, string> })
    return { record, field }
  },
  template: `
    <CFieldNumberEditor
      :field="field"
      :model-value="record.values.amount"
      @update:model-value="record.values.amount = $event"
    />
  `,
})

// PrimeVue's InputNumber writes typed text straight to the DOM input and only
// syncs its own d_value from the model-value prop. Vue re-applies the `value`
// binding on every patch of that input, so a re-render mid-typing overwrites the
// text with whatever the prop last carried.
async function typeDigit(wrapper: ReturnType<typeof mount>, digit: string) {
  const input = wrapper.find('input')
  const el = input.element as HTMLInputElement
  el.setSelectionRange(el.value.length, el.value.length)
  await input.trigger('keypress', { key: digit, code: `Digit${digit}` })
  await flushPromises()
}

describe('CFieldNumberEditor', () => {
  it('emits the typed value on every keystroke', async () => {
    const wrapper = mount(Parent)
    await flushPromises()
    await wrapper.find('input').trigger('focus')

    await typeDigit(wrapper, '2')

    const editor = wrapper.findComponent(CFieldNumberEditor)
    expect(editor.emitted('update:modelValue')?.at(-1)).toEqual(['12'])
    expect(wrapper.vm.record.values.amount).toBe('12')
  })

  it('survives a re-render while focused — model-value must not lag behind', async () => {
    const wrapper = mount(Parent)
    await flushPromises()
    await wrapper.find('input').trigger('focus')

    await typeDigit(wrapper, '2')
    expect((wrapper.find('input').element as HTMLInputElement).value).toBe('12')

    // Any re-render of the inner input re-applies its `value` binding. Freezing
    // model-value for the duration of the focus made that binding stale, which
    // wiped every keystroke and left number fields untypeable in a record list.
    wrapper.findComponent({ name: 'InputText' }).vm.$forceUpdate()
    await flushPromises()

    expect((wrapper.find('input').element as HTMLInputElement).value).toBe('12')
  })
})
