import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { buildWorkflowScope, components } from '@planetcrust/human-vue'
import ExpressionEditor from './ExpressionEditor.vue'
import CInputExpression from '@planetcrust/human-vue/src/components/expression/CInputExpression.vue'
import { currentCompletions, startCompletion } from '@codemirror/autocomplete'
import { EditorView } from '@codemirror/view'

// WorkflowEditor provides the scope; this component is nested several levels
// under it in five different configurators, so the provide/inject contract is
// the thing worth pinning.
const scope = buildWorkflowScope([
  { name: 'record', types: ['ComposeRecord'] },
  { name: 'oldRecord', types: ['ComposeRecord'] },
])

const mountEditor = (props = {}, provide = {}) =>
  mount(ExpressionEditor, {
    props: { modelValue: '', ...props },
    global: { components: { CInputExpression }, provide },
  })

describe('workflow ExpressionEditor', () => {
  it('renders the shared expression input, not a plain textarea', () => {
    const w = mountEditor()
    expect(w.find('.c-expression[data-dialect="expr"]').exists()).toBe(true)
    expect(w.find('textarea').exists()).toBe(false)
  })

  it('shows the value it is given', () => {
    expect(mountEditor({ modelValue: 'record != oldRecord' }).text()).toContain(
      'record != oldRecord',
    )
  })

  it('emits on both v-model shapes, since callers use each', async () => {
    const w = mountEditor({ modelValue: 'a' })
    const view = EditorView.findFromDOM(w.element)
    view.dispatch({ changes: { from: 1, insert: 'b' } })
    await nextTick()

    expect(w.emitted('update:modelValue')?.at(-1)).toEqual(['ab'])
    expect(w.emitted('update:value')?.at(-1)).toEqual(['ab'])
  })

  it('offers the workflow variables the editor provided', async () => {
    const w = mountEditor({}, { workflowScope: scope })
    const view = EditorView.findFromDOM(w.element)

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))

    const labels = (currentCompletions(view.state) || []).map(c => c.label)
    expect(labels).toEqual(expect.arrayContaining(['record', 'oldRecord']))
  })

  it('works with no workflow above it, offering only the language functions', async () => {
    const w = mountEditor()
    const view = EditorView.findFromDOM(w.element)

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))

    const labels = (currentCompletions(view.state) || []).map(c => c.label)
    expect(labels).toContain('coalesce')
    expect(labels).not.toContain('record')
  })

  it('accepts the scope as a ref, which is how it is provided', async () => {
    const w = mountEditor({}, { workflowScope: { value: scope } })
    const view = EditorView.findFromDOM(w.element)

    startCompletion(view)
    await new Promise(r => setTimeout(r, 60))

    expect((currentCompletions(view.state) || []).map(c => c.label)).toContain('record')
  })

  it('is exported by the shared library it now depends on', () => {
    expect(components).toBeTypeOf('object')
  })
})
