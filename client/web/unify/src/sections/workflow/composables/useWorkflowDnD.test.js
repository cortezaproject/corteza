import { describe, it, expect } from 'vitest'
import { ref } from 'vue'
import { useWorkflowDnD } from './useWorkflowDnD'
import { CANVAS_GRID, STEP_NODE } from '../lib/style'

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
    // The drop centres the node on the cursor, so a cursor at half the node
    // puts its corner at the origin. From there the rounding is visible.
    const originX = STEP_NODE.width / 2
    const originY = STEP_NODE.height / 2

    expect(drop({ x: originX, y: originY }).position.x).toBe(0)
    expect(drop({ x: originX + CANVAS_GRID / 2 + 1, y: originY }).position.x).toBe(CANVAS_GRID)
    expect(drop({ x: originX + CANVAS_GRID / 2 - 1, y: originY }).position.x).toBe(0)
  })

  it('snaps every step kind the toolbar can drop', () => {
    for (const kind of ['function', 'expressions', 'iterator', 'gateway', 'trigger', 'visual']) {
      const { position } = drop({ x: 137, y: 251, kind })

      expect(onGrid(position.x), kind).toBe(true)
      expect(onGrid(position.y), kind).toBe(true)
    }
  })
})
