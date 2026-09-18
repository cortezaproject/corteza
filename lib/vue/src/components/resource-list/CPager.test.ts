import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CPager from './CPager.vue'

const button = (wrapper: ReturnType<typeof mount>, id: string) =>
  wrapper.get(`[data-testid="pager-${id}"]`)

describe('CPager', () => {
  it('labels Previous and Next with the shared strings when none are passed', () => {
    const wrapper = mount(CPager)

    expect(button(wrapper, 'prev').text()).toBe('general.resourceList.pagination.prev')
    expect(button(wrapper, 'next').text()).toBe('general.resourceList.pagination.next')
    expect(button(wrapper, 'first').attributes('aria-label')).toBe(
      'general.resourceList.pagination.first',
    )
  })

  it('prefers the labels it is given', () => {
    const wrapper = mount(CPager, { props: { prevLabel: 'Back', nextLabel: 'Forward' } })

    expect(button(wrapper, 'prev').text()).toBe('Back')
    expect(button(wrapper, 'next').text()).toBe('Forward')
  })

  it('disables First and Previous without hasPrev, Next without hasNext', async () => {
    const wrapper = mount(CPager, { props: { hasPrev: false, hasNext: true } })

    expect(button(wrapper, 'first').attributes('disabled')).toBeDefined()
    expect(button(wrapper, 'prev').attributes('disabled')).toBeDefined()
    expect(button(wrapper, 'next').attributes('disabled')).toBeUndefined()

    await wrapper.setProps({ hasPrev: true, hasNext: false })

    expect(button(wrapper, 'first').attributes('disabled')).toBeUndefined()
    expect(button(wrapper, 'prev').attributes('disabled')).toBeUndefined()
    expect(button(wrapper, 'next').attributes('disabled')).toBeDefined()
  })

  it('emits one event per button', async () => {
    const wrapper = mount(CPager, { props: { hasPrev: true, hasNext: true } })

    await button(wrapper, 'first').trigger('click')
    await button(wrapper, 'prev').trigger('click')
    await button(wrapper, 'next').trigger('click')

    expect(Object.keys(wrapper.emitted())).toEqual(
      expect.arrayContaining(['first', 'prev', 'next']),
    )
    expect(wrapper.emitted('first')).toHaveLength(1)
    expect(wrapper.emitted('prev')).toHaveLength(1)
    expect(wrapper.emitted('next')).toHaveLength(1)
  })
})
