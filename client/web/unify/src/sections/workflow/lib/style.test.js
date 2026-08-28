import { describe, it, expect } from 'vitest'
import { CANVAS_GRID, STEP_NODE, VISUAL_NODE, getStyleFromKind } from './style'
import toolbar from './toolbar'

const cells = px => px / CANVAS_GRID
const wholeCells = px => Number.isInteger(cells(px))

describe('the canvas grid', () => {
  it('is half a step node, so a stacked step clears exactly one cell', () => {
    expect(CANVAS_GRID).toBe(STEP_NODE.height / 2)
  })

  it('measures the step box in whole cells', () => {
    expect(cells(STEP_NODE.width)).toBe(6)
    expect(cells(STEP_NODE.height)).toBe(2)
  })

  it('measures every visual box in whole cells', () => {
    expect(cells(VISUAL_NODE.swimlane.width)).toBe(10)
    expect(cells(VISUAL_NODE.swimlane.height)).toBe(5)
    expect(cells(VISUAL_NODE.content.width)).toBe(12)
    expect(cells(VISUAL_NODE.content.height)).toBe(7)
    expect(cells(VISUAL_NODE.min.width)).toBe(5)
    expect(cells(VISUAL_NODE.min.height)).toBe(2)
  })
})

describe('every kind the palette offers', () => {
  // The toolbar is the set a user can actually drop, so a kind added there
  // with an off-grid box fails here rather than on the canvas.
  const droppable = toolbar.filter(i => i.kind !== 'hr')

  it.each(droppable.map(i => [i.ref ? `${i.kind}/${i.ref}` : i.kind, i]))(
    '%s has a box of whole cells',
    (_label, item) => {
      const { width, height } = getStyleFromKind(item)

      expect(width, 'width').toBeGreaterThan(0)
      expect(height, 'height').toBeGreaterThan(0)
      expect(wholeCells(width), `${width}px is ${cells(width)} cells`).toBe(true)
      expect(wholeCells(height), `${height}px is ${cells(height)} cells`).toBe(true)
    },
  )
})
