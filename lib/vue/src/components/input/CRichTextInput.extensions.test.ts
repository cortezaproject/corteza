import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CRichTextInput from './CRichTextInput.vue'

// Every toolbar action goes through a tiptap extension. The extensions were
// regrouped in tiptap 3 (underline and link in the starter kit, tables and
// task lists in one package each), so make sure each command still exists
// and round-trips through the HTML the editor emits.

const content = [
  '<h2>Title</h2>',
  '<p><strong>bold</strong> <em>italic</em> <u>under</u> <s>strike</s> ',
  '<a href="https://example.tld">link</a> <mark data-color="#ff0">hl</mark> ',
  '<span style="color: #ff0000">red</span></p>',
  '<ul data-type="taskList"><li data-type="taskItem" data-checked="true"><label><input type="checkbox" checked="checked"><span></span></label><div><p>done</p></div></li></ul>',
  '<table><tbody><tr><th><p>h</p></th></tr><tr><td><p>c</p></td></tr></tbody></table>',
  '<p style="text-align: center">centered</p>',
  '<p>:smile:</p>',
].join('')

const mountEditor = async (modelValue = '') => {
  const w = mount(CRichTextInput, { props: { modelValue }, attachTo: document.body })
  await flushPromises()
  const editor = (w.vm as any).editor
  expect(editor, 'editor instance').toBeTruthy()
  return { w, editor }
}

describe('CRichTextInput extensions', () => {
  it('keeps every mark, node and style when loading content', async () => {
    const { editor } = await mountEditor(content)
    const html = editor.getHTML()

    for (const needle of ['<h2>', '<strong>', '<em>', '<u>', '<s>', 'href="https://example.tld"', '<mark', 'color: rgb(255, 0, 0)', 'data-type="taskList"', 'data-checked="true"', '<table', '<th', '<td', 'text-align: center']) {
      expect(html, needle).toContain(needle)
    }
  })

  it('offers every command the toolbar uses', async () => {
    const { editor } = await mountEditor('<p>text</p>')
    const c = editor.commands

    for (const cmd of ['toggleBold', 'toggleItalic', 'toggleUnderline', 'toggleStrike', 'toggleHeading', 'setParagraph',
      'toggleBulletList', 'toggleOrderedList', 'toggleTaskList', 'toggleBlockquote', 'toggleCodeBlock', 'setHorizontalRule',
      'setLink', 'unsetLink', 'setColor', 'toggleHighlight', 'unsetHighlight', 'setTextAlign', 'unsetAllMarks',
      'insertTable', 'deleteTable', 'addRowAfter', 'addColumnAfter']) {
      expect(typeof c[cmd], cmd).toBe('function')
    }

    editor.commands.selectAll()
    expect(editor.chain().focus().toggleBold().run()).toBe(true)
    expect(editor.isActive('bold')).toBe(true)
    expect(editor.chain().focus().toggleUnderline().run()).toBe(true)
    expect(editor.isActive('underline')).toBe(true)
    expect(editor.chain().focus().setLink({ href: 'https://x.tld' }).run()).toBe(true)
    expect(editor.isActive('link')).toBe(true)
    expect(editor.chain().focus().setColor('#00ff00').run()).toBe(true)
    expect(editor.getHTML()).toContain('color: rgb(0, 255, 0)')
    editor.commands.setTextSelection(2)
    expect(editor.chain().focus().toggleTaskList().run()).toBe(true)
    expect(editor.isActive('taskList')).toBe(true)
    expect(editor.getHTML()).toContain('data-type="taskList"')
  })

  it('inserts emoji', async () => {
    const { editor } = await mountEditor('<p></p>')
    expect(editor.commands.setEmoji('smile')).toBe(true)
    expect(editor.getHTML()).toContain('data-type="emoji"')
    expect(editor.getHTML()).toContain('data-name="smile"')
  })
})
