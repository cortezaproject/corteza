import { describe, it, expect } from 'vitest'
import { ref } from 'vue'
import { useWorkflowDnD } from './useWorkflowDnD'
import { CANVAS_GRID, getStyleFromKind } from '../lib/style'

// `%` keeps the sign, so a snapped -0 is on the grid as much as +0 is.
const onGrid = v => v % CANVAS_GRID === 0

// A drop event carrying the toolbar item on the dataTransfer, which is the
// branch onCanvasDrop takes when no drag-start ran. project() is the identity
// here, so (x, y) is the cursor in flow coordinates.
const dropEvent = (x, y, kind) => ({
  preventDefault() {},
  currentTarget: { getBoundingClientRect: () => ({ left: 0, top: 0 }) },
  clientX: x,
  clientY: y,
  dataTransfer: { getData: () => JSON.stringify({ kind }) },
})

function drop({ x, y, kind = 'function' }) {
  const nodes = ref([])
  const edges = ref([])
  const { onCanvasDrop } = useWorkflowDnD(
    nodes,
    edges,
    null,
    pos => pos, // project(): container coords are already flow coords
  )

  return onCanvasDrop(dropEvent(x, y, kind))
}

describe('CANVAS_GRID', () => {
  it("is half a step node's height, so a stacked step clears one whole cell", () => {
    expect(CANVAS_GRID).toBe(getStyleFromKind({ kind: 'function' }).height / 2)
  })
})

describe('onCanvasDrop grid snap', () => {
  it.each([
    [0, 0],
    [1, 1],
    [17, 47],
    [123, 456],
    [-5, -70],
    [999, 1001],
  ])('drop at (%s, %s) lands on the canvas grid', (x, y) => {
    const { position } = drop({ x, y })

    expect(onGrid(position.x)).toBe(true)
    expect(onGrid(position.y)).toBe(true)
  })

  it('rounds to the nearest grid line rather than the floor', () => {
    // 180x64 node, so the drop point is offset by (-90, -32) before snapping.
    // Cursor at (90, 32) → raw (0, 0); at (90 + 20, 32) → raw (20, 0) → 32.
    expect(drop({ x: 90, y: 32 }).position.x).toBe(0)
    expect(drop({ x: 90 + 20, y: 32 }).position.x).toBe(CANVAS_GRID)
    expect(drop({ x: 90 + 12, y: 32 }).position.x).toBe(0)
  })

  it('snaps every step kind the toolbar can drop', () => {
    for (const kind of ['function', 'expressions', 'iterator', 'gateway', 'trigger', 'visual']) {
      const { position } = drop({ x: 137, y: 251, kind })

      expect(onGrid(position.x), kind).toBe(true)
      expect(onGrid(position.y), kind).toBe(true)
    }
  })
})
