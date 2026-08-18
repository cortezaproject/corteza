import type { Edge, Node } from '@vue-flow/core'

import type { FlowNodeData } from '@/sections/taq/utils/taq-parser'

/**
 * One frame of an execution trace, as `ngAutomationExecutionTrace` returns it.
 * Only the fields the canvas colours itself from are modelled here.
 */
export interface TraceFrame {
  stepID?: string
  handle?: string
  kind?: string
  error?: unknown
}

export type FlowNode = Node<FlowNodeData>

/**
 * The trace as the canvas reads it: every frame, plus a handle index. An
 * iterator body records one frame per pass under the same handle, so the index
 * holds the last of them.
 */
export interface Trace {
  frames: TraceFrame[]
  byHandle: Map<string, TraceFrame>
}

export function indexTrace(frames: TraceFrame[]): Trace {
  const byHandle = new Map<string, TraceFrame>()
  for (const frame of frames) {
    if (frame.handle) byHandle.set(frame.handle, frame)
  }
  return { frames, byHandle }
}

/**
 * The frame a node ran under, matched on its backend handle and falling back to
 * its step ID. Null for a node the run never reached.
 */
export function traceFrameFor(node: FlowNode | null | undefined, trace: Trace): TraceFrame | null {
  if (!node) return null

  const handle = node.data?.ref
  if (handle) {
    const frame = trace.byHandle.get(handle)
    if (frame) return frame
  }

  const stepID = node.data?.stepID
  if (stepID) {
    return trace.frames.find(f => f.stepID === stepID) || null
  }

  return null
}

/** Whether a node ran. Triggers always count: the run started at one. */
export function nodeHasTrace(node: FlowNode | null | undefined, trace: Trace): boolean {
  if (!node) return false
  if (node.type === 'trigger') return true
  return traceFrameFor(node, trace) !== null
}

/**
 * Whether control flowed along this edge — what paints an edge as executed.
 *
 * Control leaves a step only when it ran and did not fail there. An End or loop
 * marker the parser invents to pad a dangling arm has no backend step behind
 * it, carries no `stepID`, and so can never have a frame: reaching one follows
 * from its source, except off a `branch`, where the trace does not record which
 * arm ran. A persisted marker carries a `stepID` and resolves by frame like any
 * other step.
 */
export function isEdgeTraversed(
  edge: Pick<Edge, 'source' | 'target'>,
  nodes: FlowNode[],
  trace: Trace,
): boolean {
  const sourceNode = nodes.find(n => n.id === edge.source)
  const targetNode = nodes.find(n => n.id === edge.target)
  if (!sourceNode || !targetNode) return false

  if (!nodeHasTrace(sourceNode, trace)) return false
  if (traceFrameFor(sourceNode, trace)?.error) return false

  const isMarker = targetNode.type === 'end' || targetNode.type === 'loop'
  if (isMarker && !targetNode.data?.stepID) return sourceNode.type !== 'branch'

  return nodeHasTrace(targetNode, trace)
}
