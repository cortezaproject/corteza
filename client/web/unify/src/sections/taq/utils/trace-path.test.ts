import { describe, expect, it } from 'vitest'

import {
  indexTrace,
  isEdgeTraversed,
  nodeHasTrace,
  traceFrameFor,
  type FlowNode,
  type TraceFrame,
} from '@/sections/taq/utils/trace-path'

// A node as taq-parser builds it: a persisted step carries its backend handle
// and stepID, a marker the parser invents to pad a dangling arm carries neither.
function step(id: string, type: string, stepID?: string): FlowNode {
  return {
    id,
    type,
    position: { x: 0, y: 0 },
    data: { ref: stepID ? id : id, stepID, label: id, nodeType: '', config: {}, arguments: [] },
  } as unknown as FlowNode
}

function marker(id: string, type: 'end' | 'loop'): FlowNode {
  return {
    id,
    type,
    position: { x: 0, y: 0 },
    data: { ref: id, label: 'End', nodeType: 'termination', config: {}, arguments: [] },
  } as unknown as FlowNode
}

function frame(handle: string, kind: string, error?: unknown): TraceFrame {
  return { handle, stepID: handle.replace(/^\w+_/, ''), kind, error }
}

const edge = (source: string, target: string) => ({ source, target })

describe('traceFrameFor', () => {
  it('matches on handle', () => {
    const trace = indexTrace([frame('step_2', 'function')])
    expect(traceFrameFor(step('step_2', 'step', '2'), trace)?.handle).toBe('step_2')
  })

  it('falls back to stepID when no handle matches', () => {
    const trace = indexTrace([{ stepID: '2', kind: 'function' }])
    expect(traceFrameFor(step('step_2', 'step', '2'), trace)).not.toBeNull()
  })

  it('keeps the last pass of an iterator body', () => {
    // The runtime records one frame per pass under the same handle.
    const trace = indexTrace([
      { handle: 'step_3', stepID: '3', kind: 'function', error: undefined },
      { handle: 'step_3', stepID: '3', kind: 'function', error: 'last' },
    ])
    expect(traceFrameFor(step('step_3', 'step', '3'), trace)?.error).toBe('last')
  })

  it('is null for a node the run never reached', () => {
    expect(traceFrameFor(step('step_9', 'step', '9'), indexTrace([]))).toBeNull()
  })
})

describe('nodeHasTrace', () => {
  it('always counts a trigger', () => {
    expect(nodeHasTrace(step('trigger_1', 'trigger', '1'), indexTrace([]))).toBe(true)
  })

  it('is false for a step with no frame', () => {
    expect(nodeHasTrace(step('step_4', 'step', '4'), indexTrace([]))).toBe(false)
  })
})

describe('isEdgeTraversed — a branch whose arms are steps', () => {
  // trigger_1 -> step_2 (branch) -> step_3 (taken) -> step_5 (End)
  //                              -> step_4 (skipped) -> step_6 (End)
  const nodes = [
    step('trigger_1', 'trigger', '1'),
    step('step_2', 'branch', '2'),
    step('step_3', 'step', '3'),
    step('step_4', 'step', '4'),
    marker('step_5', 'end'),
    marker('step_6', 'end'),
  ]
  nodes[4].data!.stepID = '5'
  nodes[5].data!.stepID = '6'

  const trace = indexTrace([
    frame('trigger_1', 'trigger'),
    frame('step_2', 'gatewayExclusive'),
    frame('step_3', 'function'),
    frame('step_5', 'termination'),
  ])

  it('lights the arm that ran', () => {
    expect(isEdgeTraversed(edge('step_2', 'step_3'), nodes, trace)).toBe(true)
  })

  it('leaves the arm that did not run unlit', () => {
    expect(isEdgeTraversed(edge('step_2', 'step_4'), nodes, trace)).toBe(false)
  })

  it('lights a reached End from its own frame', () => {
    expect(isEdgeTraversed(edge('step_3', 'step_5'), nodes, trace)).toBe(true)
  })

  it('leaves an unreached End unlit', () => {
    expect(isEdgeTraversed(edge('step_4', 'step_6'), nodes, trace)).toBe(false)
  })
})

describe('isEdgeTraversed — a branch arm going straight to End', () => {
  // trigger_1 -> step_2 (branch) -> step_3 (arm) ; -> step_5 (End, the else arm)
  const nodes = [
    step('trigger_1', 'trigger', '1'),
    step('step_2', 'branch', '2'),
    step('step_3', 'step', '3'),
    marker('step_5', 'end'),
  ]
  nodes[3].data!.stepID = '5'

  const base = [frame('trigger_1', 'trigger'), frame('step_2', 'gatewayExclusive')]

  it('lights the else arm when its termination has a frame', () => {
    const trace = indexTrace([...base, frame('step_5', 'termination')])
    expect(isEdgeTraversed(edge('step_2', 'step_5'), nodes, trace)).toBe(true)
    expect(isEdgeTraversed(edge('step_2', 'step_3'), nodes, trace)).toBe(false)
  })

  it('leaves the else arm unlit when it has none', () => {
    const trace = indexTrace([...base, frame('step_3', 'function')])
    expect(isEdgeTraversed(edge('step_2', 'step_5'), nodes, trace)).toBe(false)
  })
})

describe('isEdgeTraversed — a failed run', () => {
  // trigger_1 -> step_2 (ok) -> step_3 (failed) -> step_4 (End)
  const nodes = [
    step('trigger_1', 'trigger', '1'),
    step('step_2', 'step', '2'),
    step('step_3', 'step', '3'),
    marker('step_4', 'end'),
  ]
  nodes[3].data!.stepID = '4'

  const trace = indexTrace([
    frame('trigger_1', 'trigger'),
    frame('step_2', 'function'),
    frame('step_3', 'function', { message: 'user not found' }),
  ])

  it('lights the edge into the step that failed', () => {
    expect(isEdgeTraversed(edge('step_2', 'step_3'), nodes, trace)).toBe(true)
  })

  it('stops at the step that failed', () => {
    expect(isEdgeTraversed(edge('step_3', 'step_4'), nodes, trace)).toBe(false)
  })

  it('stops at the step that failed even where the marker is inferred', () => {
    // A TAQ authored through the API carries no termination step, so the
    // parser pads the leaf and the marker can never have a frame of its own —
    // only the source having failed keeps the edge unlit.
    const padded = [...nodes.slice(0, 3), marker('end_step_3', 'end')]
    expect(isEdgeTraversed(edge('step_3', 'end_step_3'), padded, trace)).toBe(false)
  })
})

describe('isEdgeTraversed — markers the parser invents', () => {
  const nodes = [
    step('step_2', 'step', '2'),
    step('step_7', 'branch', '7'),
    marker('end_step_2', 'end'),
    marker('loop_step_2', 'loop'),
    marker('end_step_7_0', 'end'),
  ]
  const trace = indexTrace([frame('step_2', 'function'), frame('step_7', 'gatewayExclusive')])

  it('infers a padded End from its source', () => {
    expect(isEdgeTraversed(edge('step_2', 'end_step_2'), nodes, trace)).toBe(true)
  })

  it('infers a loop marker from its source', () => {
    expect(isEdgeTraversed(edge('step_2', 'loop_step_2'), nodes, trace)).toBe(true)
  })

  it('refuses to infer one hanging off a branch', () => {
    expect(isEdgeTraversed(edge('step_7', 'end_step_7_0'), nodes, trace)).toBe(false)
  })
})
