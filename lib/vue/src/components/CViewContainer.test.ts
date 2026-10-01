import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CViewContainer from './CViewContainer.vue'

const column = (w: ReturnType<typeof mount>) => w.find('[data-col]')

describe('CViewContainer', () => {
  it('fixed shape: one element, list cap, page padding, own overflow', () => {
    const w = mount(CViewContainer, { slots: { default: '<p data-col>x</p>' } })
    const root = w.element as HTMLElement
    expect(root.className).toContain('max-w-[1440px]')
    expect(root.className).toContain('p-4 md:p-6')
    expect(root.className).toContain('h-full overflow-hidden')
    expect(column(w).element.parentElement).toBe(root)
  })

  it('scroll shape: full-width scroller around a capped form column', () => {
    const w = mount(CViewContainer, {
      props: { scroll: true },
      slots: { default: '<p data-col>x</p>' },
    })
    const scroller = w.element as HTMLElement
    const col = column(w).element.parentElement as HTMLElement
    expect(scroller.className).toContain('overflow-y-auto')
    expect(scroller.className).not.toContain('max-w-')
    expect(col.parentElement).toBe(scroller)
    expect(col.className).toContain('max-w-[1024px]')
    expect(col.className).toContain('flex flex-col gap-4')
  })

  it('width overrides the shape default', () => {
    const list = mount(CViewContainer, {
      props: { scroll: true, width: 'list' },
      slots: { default: '<p data-col>x</p>' },
    })
    expect((column(list).element.parentElement as HTMLElement).className).toContain(
      'max-w-[1440px]',
    )
    const form = mount(CViewContainer, { props: { width: 'form' } })
    expect((form.element as HTMLElement).className).toContain('max-w-[1024px]')
  })

  it('gap="5" widens the scroll column gap', () => {
    const w = mount(CViewContainer, {
      props: { scroll: true, gap: '5' },
      slots: { default: '<p data-col>x</p>' },
    })
    expect((column(w).element.parentElement as HTMLElement).className).toContain('gap-5')
  })

  it('a class passed in lands on the column, not the scroller', () => {
    const w = mount(CViewContainer, {
      props: { scroll: true },
      attrs: { class: 'extra' },
      slots: { default: '<p data-col>x</p>' },
    })
    expect((w.element as HTMLElement).className).not.toContain('extra')
    expect((column(w).element.parentElement as HTMLElement).className).toContain('extra')
  })
})
