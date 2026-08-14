import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import CInputExpression from './CInputExpression.vue'
import CExpressionHint from './CExpressionHint.vue'
import { buildScope } from './catalog'
import { currentCompletions, startCompletion } from '@codemirror/autocomplete'
import { EditorView } from '@codemirror/view'

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

const mountInput = ({ attachTo, ...props }: any = {}) =>
  mount(CInputExpression, { props: { modelValue: '', scope, ...props }, attachTo })

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

// Drives CodeMirror's own completion state, so the wiring between
// completionAt's options and what lands in the document is covered — not just
// the option data, which syntax.test.ts already pins.
describe('CInputExpression completion wiring', () => {
  const typeInto = async (w, text: string) => {
    const view = EditorView.findFromDOM(w.element)
    view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } })
    view.dispatch({ selection: { anchor: text.length } })
    await nextTick()
    return view
  }

  const labelsOf = (view: any) => (currentCompletions(view.state) || []).map((c: any) => c.label)

  it('offers the scope after a bare $ and writes the braces on accept', async () => {
    const w = mountInput({ modelValue: '' })
    const view = await typeInto(w, 'x = $')

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))
    expect(labelsOf(view)).toContain('recordID')

    // Accept the recordID option explicitly rather than relying on ordering.
    const opt = currentCompletions(view.state).find((c: any) => c.label === 'recordID')
    opt.apply(view, opt, 4, 5)
    await nextTick()
    expect(view.state.doc.toString()).toBe('x = ${recordID}')
  })

  it('lands the caret inside the braces when an object is accepted', async () => {
    const w = mountInput({ modelValue: '' })
    const view = await typeInto(w, '$')

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))
    const opt = currentCompletions(view.state).find((c: any) => c.label === 'record')
    opt.apply(view, opt, 0, 1)
    await nextTick()

    expect(view.state.doc.toString()).toBe('${record.}')
    // Caret sits before the closing brace, ready for the member name.
    expect(view.state.selection.main.head).toBe(9)
  })

  it('keeps Escape from reaching the dialog while the suggestion list is open', async () => {
    const w = mountInput({ modelValue: '', queryFields: orders.fields, attachTo: document.body })
    const view = await typeInto(w, 'sta')

    view.focus()
    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))
    expect(currentCompletions(view.state).length).toBeGreaterThan(0)

    // What a PrimeVue dialog listens on: Escape arriving at the document.
    const reachedDialog = vi.fn()
    document.addEventListener('keydown', reachedDialog)

    view.contentDOM.dispatchEvent(
      new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }),
    )
    await nextTick()

    expect(reachedDialog).not.toHaveBeenCalled()
    document.removeEventListener('keydown', reachedDialog)
    w.unmount()
  })

  it('lets Escape through once no suggestions are open', async () => {
    const w = mountInput({ modelValue: 'x', attachTo: document.body })
    const view = await typeInto(w, 'x')

    const reachedDialog = vi.fn()
    document.addEventListener('keydown', reachedDialog)

    view.contentDOM.dispatchEvent(
      new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }),
    )
    await nextTick()

    // Escape must still close the dialog when it is not dismissing a popup.
    expect(reachedDialog).toHaveBeenCalled()
    document.removeEventListener('keydown', reachedDialog)
    w.unmount()
  })

  it('offers module fields, record fields and keywords on explicit request', async () => {
    const w = mountInput({ modelValue: '', queryFields: orders.fields })
    const view = await typeInto(w, '')

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))
    const labels = labelsOf(view)
    expect(labels).toEqual(expect.arrayContaining(['status', 'recordID', 'AND']))
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
