import { mount, flushPromises } from '@vue/test-utils'
import { defineComponent, h, ref, Teleport } from 'vue'
import { describe, it, expect, beforeEach } from 'vitest'

import { BASE_TITLE, composeTitle, readTitleText, useDocumentTitle } from './documentTitle'

// The tab title is whatever the active view teleported into `#topbar-title`.
// What this file pins is the reading of that node — what counts as the heading,
// what is chrome sitting beside it — and that a live change to it reaches
// document.title without a navigation.

const el = html => {
  const node = document.createElement('div')
  node.innerHTML = html
  return node
}

describe('readTitleText', () => {
  it('reads the heading text', () => {
    expect(readTitleText(el('<span>Users</span>'))).toBe('Users')
  })

  it('condenses the whitespace a multi-line template leaves behind', () => {
    expect(readTitleText(el('<span>\n  Nightly   sync\n</span>'))).toBe('Nightly sync')
  })

  it('keeps sibling elements apart', () => {
    expect(readTitleText(el('<span>Sales</span><span>Pipeline</span>'))).toBe('Sales Pipeline')
  })

  it('leaves out a subtree marked data-title-exclude', () => {
    const node = el(
      '<span><span>My project</span><span data-title-exclude>Revision 2 <b>Draft</b></span></span>',
    )
    expect(readTitleText(node)).toBe('My project')
  })

  it('is empty for a node with no heading in it', () => {
    expect(readTitleText(el(''))).toBe('')
    expect(readTitleText(null)).toBe('')
  })
})

describe('composeTitle', () => {
  it('is the heading alone', () => {
    expect(composeTitle('Users', 0)).toBe('Users')
  })

  it('falls back to the app name when there is no heading', () => {
    expect(composeTitle('', 0)).toBe(BASE_TITLE)
  })

  it('puts the unread count in front', () => {
    expect(composeTitle('Users', 3)).toBe('(3) Users')
    expect(composeTitle('', 3)).toBe(`(3) ${BASE_TITLE}`)
  })
})

describe('useDocumentTitle', () => {
  const unread = ref(0)
  const heading = ref('')

  // The shell: owns the teleport target and the title, exactly as App.vue does.
  const Shell = defineComponent({
    setup() {
      useDocumentTitle(() => unread.value)
      return () => [
        h('div', { id: 'topbar-title' }),
        heading.value
          ? h(Teleport, { to: '#topbar-title', defer: true }, [h('span', heading.value)])
          : null,
      ]
    },
  })

  beforeEach(() => {
    unread.value = 0
    heading.value = ''
    document.title = ''
  })

  it('takes the app name while no view has stated a heading', async () => {
    const wrapper = mount(Shell, { attachTo: document.body })
    await flushPromises()

    expect(document.title).toBe(BASE_TITLE)
    wrapper.unmount()
  })

  it('follows the heading a view teleports in, and its later changes', async () => {
    const wrapper = mount(Shell, { attachTo: document.body })
    await flushPromises()

    heading.value = 'Users'
    await flushPromises()
    expect(document.title).toBe('Users')

    heading.value = 'Invoice INV-0042'
    await flushPromises()
    expect(document.title).toBe('Invoice INV-0042')

    // Navigating away from a view that states no heading of its own.
    heading.value = ''
    await flushPromises()
    expect(document.title).toBe(BASE_TITLE)

    wrapper.unmount()
  })

  it('re-prefixes the unread count without losing the heading', async () => {
    const wrapper = mount(Shell, { attachTo: document.body })
    heading.value = 'Users'
    await flushPromises()

    unread.value = 3
    await flushPromises()
    expect(document.title).toBe('(3) Users')

    unread.value = 0
    await flushPromises()
    expect(document.title).toBe('Users')

    wrapper.unmount()
  })

  it('carries a standing count onto the next heading', async () => {
    const wrapper = mount(Shell, { attachTo: document.body })
    unread.value = 2
    heading.value = 'Users'
    await flushPromises()
    expect(document.title).toBe('(2) Users')

    // Navigation, with the count untouched.
    heading.value = 'TAQ Automations'
    await flushPromises()
    expect(document.title).toBe('(2) TAQ Automations')

    wrapper.unmount()
  })

  it('stops writing the title once the shell is gone', async () => {
    const wrapper = mount(Shell, { attachTo: document.body })
    heading.value = 'Users'
    await flushPromises()

    wrapper.unmount()
    document.title = 'untouched'
    heading.value = 'Roles'
    await flushPromises()

    expect(document.title).toBe('untouched')
  })
})
