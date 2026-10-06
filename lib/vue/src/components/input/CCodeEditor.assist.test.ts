import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { EditorState } from '@codemirror/state'
import { EditorView } from '@codemirror/view'
import { currentCompletions, snippet, startCompletion } from '@codemirror/autocomplete'
import CCodeEditor from './CCodeEditor.vue'

const mountEditor = (props: any = {}) =>
  mount(CCodeEditor, { props: { modelValue: '', ...props }, attachTo: document.body })

const typeInto = async (w: any, text: string) => {
  const view = EditorView.findFromDOM(w.element)!
  view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } })
  view.dispatch({ selection: { anchor: text.length } })
  await nextTick()
  return view
}

const suggest = async (view: EditorView) => {
  view.focus()
  startCompletion(view)
  await new Promise(r => setTimeout(r, 60))
  return (currentCompletions(view.state) || []).map((c: any) => c.label)
}

// The source is made once: CodeMirror tells sources apart by identity, so one
// made afresh on every call is a new source each keystroke and never answers.
const everywhere = (labels: string[]) => {
  const autocomplete = (cx: any) => ({
    from: cx.matchBefore(/\w*/).from,
    options: labels.map(label => ({ label })),
  })
  return EditorState.languageData.of(() => [{ autocomplete }])
}

describe('CCodeEditor assist', () => {
  it('suggests nothing without assist, so the template editor stays as it was', async () => {
    const w = mountEditor()
    const view = await typeInto(w, '<di')

    expect(await suggest(view)).toEqual([])
    w.unmount()
  })

  it('suggests HTML tags, and what the assist adds, when given one', async () => {
    const w = mountEditor({ assist: [everywhere(['divine'])] })
    const view = await typeInto(w, '<di')

    const labels = await suggest(view)
    expect(labels).toContain('div')
    expect(labels).toContain('divine')
    w.unmount()
  })

  it('keeps Escape from reaching the dialog while the suggestion list is open', async () => {
    const w = mountEditor({ assist: [] })
    const view = await typeInto(w, '<di')
    expect((await suggest(view)).length).toBeGreaterThan(0)

    const reachedDialog = vi.fn()
    document.addEventListener('keydown', reachedDialog)
    view.contentDOM.dispatchEvent(
      new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }),
    )
    await nextTick()

    expect(reachedDialog).not.toHaveBeenCalled()
    expect(currentCompletions(view.state)).toEqual([])
    document.removeEventListener('keydown', reachedDialog)
    w.unmount()
  })

  it('keeps Escape from reaching the dialog while leaving the fields of a call', async () => {
    const w = mountEditor({ assist: [] })
    const view = EditorView.findFromDOM(w.element)!
    view.focus()
    // Two fields, the way a call such as records.read is written; a call with
    // one field only selects it and leaves nothing to leave.
    snippet("f('${a}', ${b})")(view, null, 0, 0)

    const reachedDialog = vi.fn()
    document.addEventListener('keydown', reachedDialog)
    const escape = () =>
      view.contentDOM.dispatchEvent(
        new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }),
      )

    escape()
    expect(reachedDialog).not.toHaveBeenCalled()

    escape()
    expect(reachedDialog).toHaveBeenCalledTimes(1)
    document.removeEventListener('keydown', reachedDialog)
    w.unmount()
  })

  it('inserts at the cursor, indenting later lines like the line it lands on', async () => {
    const w = mountEditor({ modelValue: '<script>\n    \n</script>' })
    const view = EditorView.findFromDOM(w.element)!
    view.dispatch({ selection: { anchor: '<script>\n    '.length } })
    ;(w.vm as any).insert('a()\nb()\n')
    await nextTick()

    expect(view.state.doc.toString()).toBe('<script>\n    a()\n    b()\n    \n</script>')
    expect(w.emitted('update:modelValue')?.at(-1)).toEqual([view.state.doc.toString()])
    w.unmount()
  })
})
