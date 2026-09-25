import { describe, expect, it, vi } from 'vitest'
import { defineComponent, h } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import { RANGES, sum, useSystemStats } from './useSystemStats'

function mountWith(api) {
  let exposed
  const Host = defineComponent({
    setup() {
      exposed = useSystemStats()
      return () => h('div')
    },
  })
  const wrapper = mount(Host, { global: { provide: { $SystemAPI: api } } })
  return { wrapper, ...exposed }
}

function fakeApi(payload) {
  const calls = []
  const api = {
    statsListCancellable(params) {
      calls.push(params)
      return { response: () => Promise.resolve(payload), cancel: () => {} }
    },
  }
  return { api, calls }
}

describe('useSystemStats', () => {
  it('requests the preset range with its bucket, once per range change', async () => {
    localStorage.removeItem('admin.dashboard.range')
    const { api, calls } = fakeApi({
      range: { bucket: 'day', buckets: ['2026-09-20'] },
      resources: {},
    })
    const { range, stats } = mountWith(api)
    await flushPromises()

    expect(calls).toHaveLength(1)
    expect(calls[0].bucket).toBe('day')
    // from is local midnight 30 days back, so the span is 30 days plus today so far
    const days = (new Date(calls[0].to) - new Date(calls[0].from)) / 86400000
    expect(days).toBeGreaterThanOrEqual(30)
    expect(days).toBeLessThan(31)
    expect(stats.value.range.bucket).toBe('day')

    range.value = '90d'
    await flushPromises()
    expect(calls).toHaveLength(2)
    expect(calls[1].bucket).toBe(RANGES.find(r => r.key === '90d').bucket)
  })

  it('labels day, week and month buckets', async () => {
    const { api } = fakeApi({
      range: { bucket: 'week', buckets: ['2026-09-14', '2026-09-21'], to: '2026-09-25T10:00:00Z' },
    })
    const { bucketLabels, rangeLabels } = mountWith(api)
    await flushPromises()
    expect(bucketLabels.value).toHaveLength(2)
    expect(rangeLabels.value[0]).toContain('–')
    expect(rangeLabels.value[1]).toContain('–')
  })

  it('sums a series', () => {
    expect(sum([1, 2, 3])).toBe(6)
    expect(sum(undefined)).toBe(0)
  })
})
