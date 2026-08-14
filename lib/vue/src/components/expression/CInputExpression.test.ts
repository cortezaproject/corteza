import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import CInputExpression from './CInputExpression.vue'
import CExpressionHint from './CExpressionHint.vue'
import { buildScope } from './catalog'

const orders = {
  moduleID: '1',
  fields: [
    { name: 'status', label: 'Status', kind: 'Select' },
    { name: 'quantity', label: 'Quantity', kind: 'Number' },
    { name: 'owner', kind: 'User' },
    { name: 'due', kind: 'DateTime' },
    { name: 'note', kind: 'String' },
  ],
}

const scope = buildScope({ recordModule: orders, hasRecord: true })

const mountInput = (props = {}) =>
  mount(CInputExpression, { props: { modelValue: '', scope, ...props } })

describe('CInputExpression', () => {
  it('renders the model value into the editor', () => {
    const w = mountInput({ modelValue: "status = 'Open'" })
    expect(w.text()).toContain("status = 'Open'")
  })

  it('emits on edit and does not echo the value back into the document', async () => {
    const w = mountInput({ modelValue: 'a' })
    w.vm.insert('b')
    await nextTick()

    const emitted = w.emitted('update:modelValue')
    expect(emitted?.at(-1)).toEqual(['ab'])
    expect(w.text()).toContain('ab')
  })

  it('applies an external value change without re-emitting it', async () => {
    const w = mountInput({ modelValue: 'a' })
    await w.setProps({ modelValue: 'zzz' })
    await nextTick()

    expect(w.text()).toContain('zzz')
    expect(w.emitted('update:modelValue')).toBeUndefined()
  })

  it('marks ${} holes so they are visually distinct', () => {
    const w = mountInput({ modelValue: 'id = ${recordID}' })
    expect(w.find('.c-expression__hole').exists()).toBe(true)
  })

  it('marks QL keywords in the ql dialect only', () => {
    expect(
      mountInput({ modelValue: 'a = 1 AND b = 2' }).find('.c-expression__keyword').exists(),
    ).toBe(true)
    expect(
      mountInput({ modelValue: 'a = 1 AND b = 2', dialect: 'interpolation' })
        .find('.c-expression__keyword')
        .exists(),
    ).toBe(false)
  })

  it('inserts at the cursor rather than appending', async () => {
    const w = mountInput({ modelValue: '' })
    w.vm.insert('${recordID}')
    await nextTick()
    expect(w.emitted('update:modelValue')?.at(-1)).toEqual(['${recordID}'])
  })

  it('is not editable when disabled', () => {
    const w = mountInput({ modelValue: 'x', disabled: true })
    expect(w.find('.c-expression--disabled').exists()).toBe(true)
    expect(w.find('[contenteditable="true"]').exists()).toBe(false)
  })
})

describe('CExpressionHint', () => {
  const mountHint = (props = {}) => mount(CExpressionHint, { props: { scope, ...props } })

  it('names the record variables on a record page', () => {
    const text = mountHint().text()
    expect(text).toContain('${recordID}')
    expect(text).toContain('${ownerID}')
    expect(text).toContain('${user.email}')
  })

  it('uses the real module field names', () => {
    const text = mountHint().text()
    expect(text).toContain('${record.values.status}')
    expect(text).toContain('${record.values.quantity}')
  })

  it('caps the field list rather than listing every field', () => {
    const text = mountHint({ maxFields: 2 }).text()
    expect(text).toContain('${record.values.status}')
    expect(text).not.toContain('${record.values.note}')
    expect(text).toContain('${record.values.…}')
  })

  it('names the shape when the module has not loaded', () => {
    const text = mountHint({ scope: buildScope({ recordModule: null, hasRecord: true }) }).text()
    expect(text).toContain('${record.values.fieldName}')
  })

  it('hides record variables when the page has no record', () => {
    const text = mountHint({ scope: buildScope({ hasRecord: false }) }).text()
    expect(text).not.toContain('${recordID}')
    expect(text).not.toContain('${record.values')
    expect(text).toContain('${userID}')
  })

  it('emits the chip text for insertion when clicked', async () => {
    const w = mountHint()
    await w.findAll('.c-expression-hint__chip')[0].trigger('click')
    expect(w.emitted('insert')?.[0]).toEqual(['${recordID}'])
  })
})
