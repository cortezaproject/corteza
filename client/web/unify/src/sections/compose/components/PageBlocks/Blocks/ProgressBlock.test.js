import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('./PageBlock.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

import ProgressBlock from './ProgressBlock.vue'

// The value comes off an aggregate, so an averaged field arrives as
// 42.77777777777778. The relative label has always rounded its percentage; the
// absolute one printed the float in full, straight into a progress bar.

function mountBlock(options, fetched) {
  return mount(ProgressBlock, {
    props: {
      block: {
        options,
        // the class method the block calls for live values
        fetch: () => Promise.resolve(fetched),
      },
      namespace: { namespaceID: 'N1' },
    },
    global: {
      stubs: { ProgressBar: true, ProgressSpinner: true },
      provide: { $ComposeAPI: {}, $Auth: { user: { userID: 'U1' } }, $eventBus: null },
    },
  })
}

const label = w => w.get('span').text()

const staticValue = v => ({ default: v, moduleID: '', filter: '', field: '', operation: '' })

describe('ProgressBlock label', () => {
  it('rounds an absolute value to two decimals', async () => {
    const w = mountBlock({
      value: staticValue(42.77777777777778),
      minValue: staticValue(0),
      maxValue: staticValue(100),
      display: { showValue: true, showRelative: false, showProgress: true },
    })
    await flushPromises()

    expect(label(w)).toBe('42.78 / 100')
  })

  it('leaves a value that needs no rounding alone', async () => {
    const w = mountBlock({
      value: staticValue(42),
      minValue: staticValue(0),
      maxValue: staticValue(60),
      display: { showValue: true, showRelative: false, showProgress: false },
    })
    await flushPromises()

    expect(label(w)).toBe('42')
  })

  it('still rounds the percentage on the relative label', async () => {
    const w = mountBlock({
      value: staticValue(42.77777777777778),
      minValue: staticValue(0),
      maxValue: staticValue(100),
      display: { showValue: true, showRelative: true, showProgress: true },
    })
    await flushPromises()

    expect(label(w)).toBe('43% / 100%')
  })
})
