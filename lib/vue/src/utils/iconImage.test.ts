import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderIcon } from './iconImage'

// jsdom neither loads images nor draws, so both are recorded instead.
let natural = { width: 0, height: 0 }

class FakeImage {
  crossOrigin = ''
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  naturalWidth = natural.width
  naturalHeight = natural.height

  set src(value: string) {
    if (value.includes('stalled')) return
    queueMicrotask(() => (value.includes('broken') ? this.onerror?.() : this.onload?.()))
  }
}

function fakeContext() {
  const calls: Array<[string, ...unknown[]]> = []
  const ctx: any = {
    globalCompositeOperation: 'source-over',
    fillStyle: '',
    drawImage: (...a: unknown[]) => calls.push(['drawImage', ...a.slice(1)]),
    beginPath: () => {},
    arc: (...a: unknown[]) => calls.push(['arc', ...a.slice(0, 3)]),
    fill: () => calls.push(['fill', ctx.globalCompositeOperation, ctx.fillStyle]),
  }
  return { ctx, calls }
}

let context: ReturnType<typeof fakeContext>
const createElement = document.createElement.bind(document)

beforeEach(() => {
  vi.stubGlobal('Image', FakeImage)
  context = fakeContext()
  vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
    if (tag !== 'canvas') return createElement(tag)
    return {
      width: 0,
      height: 0,
      getContext: () => context.ctx,
      toDataURL: () => 'data:image/png;base64,DRAWN',
    } as unknown as HTMLCanvasElement
  })
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('renderIcon', () => {
  it('contains and centres a non-square image in the square', async () => {
    natural = { width: 200, height: 100 }

    expect(await renderIcon('/wide.png')).toBe('data:image/png;base64,DRAWN')
    expect(context.calls).toEqual([['drawImage', 0, 16, 64, 32]])
  })

  it('fills the square with an image that reports no size, as an SVG can', async () => {
    natural = { width: 0, height: 0 }

    await renderIcon('/icon.svg')

    expect(context.calls).toEqual([['drawImage', 0, 0, 64, 64]])
  })

  it('draws a dot in the top-right corner only when asked, cut out from the icon', async () => {
    natural = { width: 64, height: 64 }

    await renderIcon('/icon.svg', { dot: true, dotColor: 'rgb(1, 2, 3)' })

    const r = 64 * 0.22
    expect(context.calls.slice(1)).toEqual([
      ['arc', 64 - r, r, r + 64 * 0.06],
      ['fill', 'destination-out', ''],
      ['arc', 64 - r, r, r],
      ['fill', 'source-over', 'rgb(1, 2, 3)'],
    ])
  })

  it('gives up on an image that never finishes loading', async () => {
    vi.useFakeTimers()
    const drawing = renderIcon('/stalled.png')
    const settled = expect(drawing).rejects.toThrow('icon took too long: /stalled.png')

    vi.advanceTimersByTime(5_000)

    await settled
    vi.useRealTimers()
  })

  it('rejects when the image does not load', async () => {
    await expect(renderIcon('/broken.png')).rejects.toThrow('icon did not load: /broken.png')
  })
})
