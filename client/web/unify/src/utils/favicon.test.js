import { describe, it, expect, vi, beforeEach } from 'vitest'
import { effectScope, nextTick, ref } from 'vue'
import { flushPromises } from '@vue/test-utils'

const renderIcon = vi.hoisted(() => vi.fn())
vi.mock('@planetcrust/human-vue', () => ({ renderIcon }))

const { BUILT_IN_ICON, useFavicon } = await import('./favicon')

const CUSTOM = 'http://api.test/attachment/settings/1/original/logo.png'
const link = () => document.querySelector('link[rel="icon"]')

function start(iconUrl, hasUnread) {
  effectScope().run(() => useFavicon(iconUrl, hasUnread))
  return flushPromises()
}

beforeEach(() => {
  document.head.innerHTML = '<link rel="icon" type="image/svg+xml" href="/icon.svg" />'
  renderIcon.mockReset()
  renderIcon.mockImplementation(async (src, { dot }) => `png:${src}:${dot ? 'dot' : 'plain'}`)
})

describe('useFavicon', () => {
  it('shows the configured icon as a PNG, with no dot while nothing is unread', async () => {
    await start(ref(CUSTOM), ref(false))

    expect(link().getAttribute('href')).toBe(`png:${CUSTOM}:plain`)
    expect(link().type).toBe('image/png')
  })

  it('adds the dot while something is unread, and takes it away again', async () => {
    const unread = ref(true)
    await start(ref(CUSTOM), unread)
    expect(link().getAttribute('href')).toBe(`png:${CUSTOM}:dot`)

    unread.value = false
    await nextTick()
    await flushPromises()

    expect(link().getAttribute('href')).toBe(`png:${CUSTOM}:plain`)
  })

  it('falls back to the built-in icon when the configured one cannot be drawn', async () => {
    renderIcon.mockImplementation(async (src, { dot }) => {
      if (src === CUSTOM) throw new Error('icon did not load')
      return `png:${src}:${dot ? 'dot' : 'plain'}`
    })

    await start(ref(CUSTOM), ref(true))

    expect(link().getAttribute('href')).toBe(`png:${BUILT_IN_ICON}:dot`)
  })

  it('uses the built-in icon when none is configured', async () => {
    await start(ref(undefined), ref(false))

    expect(link().getAttribute('href')).toBe(`png:${BUILT_IN_ICON}:plain`)
  })

  it('keeps the newest state when an earlier, slower draw finishes last', async () => {
    let finishSlow
    renderIcon.mockImplementationOnce(
      (src, { dot }) => new Promise(resolve => (finishSlow = () => resolve(`slow:${dot}`))),
    )
    const unread = ref(true)
    await start(ref(CUSTOM), unread)

    unread.value = false
    await nextTick()
    await flushPromises()
    finishSlow()
    await flushPromises()

    expect(link().getAttribute('href')).toBe(`png:${CUSTOM}:plain`)
  })
})
