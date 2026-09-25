import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CResourceStatusFilter from './CResourceStatusFilter.vue'

window.matchMedia =
  window.matchMedia ||
  ((() => ({ matches: false, addEventListener() {}, removeEventListener() {} })) as never)

describe('CResourceStatusFilter', () => {
  const mountFilter = (filter: Record<string, unknown>) =>
    mount(CResourceStatusFilter, { props: { filter, states: ['archived', 'deleted'] } })

  it("offers Active and the list's own states, in order", () => {
    const select = mountFilter({ archived: '0', deleted: '0' }).findComponent({ name: 'Select' })
    expect(select.props('options').map((o: { value: string }) => o.value)).toEqual([
      'active',
      'archived',
      'deleted',
    ])
    expect(select.props('modelValue')).toBe('active')
  })

  it('writes the chosen status as the per-state filters', async () => {
    const w = mountFilter({ archived: '0', deleted: '0' })
    w.findComponent({ name: 'Select' }).vm.$emit('update:modelValue', 'archived')
    expect(w.emitted('update:filter')?.[0]).toEqual([{ archived: '2', deleted: '0' }])
  })
})
