import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import CResourceList from './CResourceList.vue'

// PrimeVue overlays query the media list on mount; jsdom has no matchMedia.
window.matchMedia =
  window.matchMedia ||
  ((() => ({ matches: false, addEventListener() {}, removeEventListener() {} })) as never)

const items = [
  { id: '1', name: 'first' },
  { id: '2', name: 'second' },
]

function mountList(actionItems: (row: unknown) => unknown[]) {
  return mount(CResourceList, {
    attachTo: document.body,
    props: {
      primaryKey: 'id',
      fields: [{ key: 'name', header: 'Name' }],
      items,
      filter: {},
      sorting: {},
      pagination: { total: 2, limit: 10, page: 1 },
      actionItems,
    },
  })
}

async function openMenu(wrapper: ReturnType<typeof mountList>, row: number) {
  const buttons = wrapper.findAll('button.row-action-btn')
  await buttons[row].trigger('click')
  await new Promise(resolve => setTimeout(resolve, 20))
  return buttons[row].element
}

const menuOf = (wrapper: ReturnType<typeof mountList>) =>
  wrapper.findComponent({ name: 'Menu' }).vm as unknown as {
    target: HTMLElement
    hide: () => void
  }

describe('CResourceList row actions menu', () => {
  // One popup is shared by every row, so the anchor has to be reassigned on
  // each open. A popup that keeps its first anchor opens beside the wrong row —
  // and at the viewport corner once that row is gone.
  it('anchors to the row whose button was clicked, on every open', async () => {
    const wrapper = mountList(() => [{ label: 'Edit', icon: 'pi pi-pencil' }])
    const menu = menuOf(wrapper)

    const firstButton = await openMenu(wrapper, 0)
    expect(menu.target).toBe(firstButton)

    menu.hide()
    await wrapper.vm.$nextTick()

    const secondButton = await openMenu(wrapper, 1)
    expect(menu.target).toBe(secondButton)
    expect(menu.target).not.toBe(firstButton)

    wrapper.unmount()
  })

  it('rebuilds the model from the clicked row', async () => {
    const actionItems = vi.fn((row: { name: string }) => [{ label: `Edit ${row.name}` }])
    const wrapper = mountList(actionItems as never)

    await openMenu(wrapper, 1)
    expect(actionItems).toHaveBeenCalledWith(items[1])

    wrapper.unmount()
  })
})
