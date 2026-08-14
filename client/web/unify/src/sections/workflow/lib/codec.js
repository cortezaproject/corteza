/**
 * Workflow codec — serializes between VueFlow nodes/edges and the API model.
 *
 * decodeWorkflow  : (workflow, triggers) → { nodes, edges }
 * encodeWorkflow  : (nodes, edges)       → { steps, paths, triggers }
 */
import { getStyleFromKind } from './style'

/**
 * Handle ↔ mxGraph exit/entry bidirectional map.
 *
 * mxGraph stores edge anchors as exit/entry fractional coordinates in the
 * edge style string:
 *   "exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;"
 *
 * VueFlow uses named handle IDs on each node. The current FlowNode components
 * expose a fixed set of handles (same IDs across workflow/trigger/termination
 * — gateway/iterator are not distinguished at the DOM level, their "slots"
 * are purely positional).
 *
 * Source handles:
 *   source-bottom        (50%, bottom)
 *   source-bottom-left   (25%, bottom)
 *   source-bottom-right  (75%, bottom)
 *   source-left          (0%,   mid)
 *   source-right         (100%, mid)
 *
 * Target handles:
 *   target-top           (50%, top)
 *   target-top-left      (25%, top)
 *   target-top-right     (75%, top)
 *   target-left          (0%,   mid)
 *   target-bottom        (50%, bottom)
 *
 * The map below is canonical: every handle ID has one exit/entry coord pair,
 * and every coord pair decodes back to the same handle. Round-trip stability
 * is enforced by encoding directly from the map (no coord tolerance logic).
 */
const HANDLE_TO_MX = {
  'source-top': { exitX: 0.5, exitY: 0 },
  'source-top-left': { exitX: 0.25, exitY: 0 },
  'source-top-right': { exitX: 0.75, exitY: 0 },
  'source-bottom': { exitX: 0.5, exitY: 1 },
  'source-bottom-left': { exitX: 0.25, exitY: 1 },
  'source-bottom-right': { exitX: 0.75, exitY: 1 },
  'source-right': { exitX: 1, exitY: 0.5 },
  'source-left': { exitX: 0, exitY: 0.5 },
  'target-top': { entryX: 0.5, entryY: 0 },
  'target-top-left': { entryX: 0.25, entryY: 0 },
  'target-top-right': { entryX: 0.75, entryY: 0 },
  'target-left': { entryX: 0, entryY: 0.5 },
  'target-right': { entryX: 1, entryY: 0.5 },
  'target-bottom': { entryX: 0.5, entryY: 1 },
  'target-bottom-left': { entryX: 0.25, entryY: 1 },
  'target-bottom-right': { entryX: 0.75, entryY: 1 },
}

// Reverse lookup built from HANDLE_TO_MX; the key is the rounded coord pair.
// Using a canonical key guarantees decode → encode → decode returns the
// same handle ID (Fix #1).
const MX_TO_SOURCE_HANDLE = {}
const MX_TO_TARGET_HANDLE = {}
for (const [id, c] of Object.entries(HANDLE_TO_MX)) {
  if (id.startsWith('source-')) {
    MX_TO_SOURCE_HANDLE[`${c.exitX}|${c.exitY}`] = id
  } else {
    MX_TO_TARGET_HANDLE[`${c.entryX}|${c.entryY}`] = id
  }
}

// Snap a fractional coord coming from a Human-written style string to the
// nearest canonical value we know about. Without this an edge dragged in
// mxGraph to, e.g., exitX=0.247 would decode as an unknown handle.
function snapCoord(value, candidates) {
  let best = candidates[0]
  let bestDiff = Math.abs(value - best)
  for (let i = 1; i < candidates.length; i++) {
    const diff = Math.abs(value - candidates[i])
    if (diff < bestDiff) {
      best = candidates[i]
      bestDiff = diff
    }
  }
  return best
}

const SOURCE_X_CANDIDATES = [0, 0.25, 0.5, 0.75, 1]
const SOURCE_Y_CANDIDATES = [0, 0.5, 1]
const TARGET_X_CANDIDATES = [0, 0.25, 0.5, 0.75, 1]
const TARGET_Y_CANDIDATES = [0, 0.5, 1]

function mxStyleToHandles(style) {
  if (!style) return {}
  const exitXm = style.match(/exitX=([0-9.]+)/)
  const exitYm = style.match(/exitY=([0-9.]+)/)
  const entryXm = style.match(/entryX=([0-9.]+)/)
  const entryYm = style.match(/entryY=([0-9.]+)/)
  const result = {}

  if (exitXm && exitYm) {
    const x = snapCoord(parseFloat(exitXm[1]), SOURCE_X_CANDIDATES)
    const y = snapCoord(parseFloat(exitYm[1]), SOURCE_Y_CANDIDATES)
    const key = `${x}|${y}`
    result.sourceHandle = MX_TO_SOURCE_HANDLE[key] || 'source-bottom'
  }

  if (entryXm && entryYm) {
    const x = snapCoord(parseFloat(entryXm[1]), TARGET_X_CANDIDATES)
    const y = snapCoord(parseFloat(entryYm[1]), TARGET_Y_CANDIDATES)
    const key = `${x}|${y}`
    result.targetHandle = MX_TO_TARGET_HANDLE[key] || 'target-top'
  }

  return result
}

/**
 * Encode VueFlow handle IDs back into an mxGraph style string.
 *
 * If `baseStyle` is provided (an original Human-written style string like
 * "edgeStyle=orthogonalEdgeStyle;rounded=1;exitX=1;exitY=0.5;..."), we splice
 * our canonical exit/entry values in, preserving any unknown properties so
 * data written by Human mxGraph clients survives the round-trip.
 */
function handlesToMxStyle(sourceHandle, targetHandle, baseStyle = '') {
  const src = HANDLE_TO_MX[sourceHandle] || HANDLE_TO_MX['source-bottom']
  const tgt = HANDLE_TO_MX[targetHandle] || HANDLE_TO_MX['target-top']
  const fresh = {
    exitX: src.exitX ?? 0.5,
    exitY: src.exitY ?? 1,
    exitDx: 0,
    exitDy: 0,
    entryX: tgt.entryX ?? 0.5,
    entryY: tgt.entryY ?? 0,
    entryDx: 0,
    entryDy: 0,
  }

  // Parse baseStyle into ordered key/value pairs so we can keep foreign
  // properties in their original position.
  const parts = []
  const seen = new Set()
  if (baseStyle) {
    baseStyle.split(';').forEach(seg => {
      if (!seg) return
      const eq = seg.indexOf('=')
      if (eq < 0) {
        // bare token like "orthogonalEdgeStyle" — keep as-is
        if (!seen.has(seg)) {
          parts.push({ bare: seg })
          seen.add(seg)
        }
        return
      }
      const key = seg.slice(0, eq)
      const val = seg.slice(eq + 1)
      if (Object.prototype.hasOwnProperty.call(fresh, key)) {
        if (!seen.has(key)) {
          parts.push({ key, value: String(fresh[key]) })
          seen.add(key)
        }
      } else if (!seen.has(key)) {
        parts.push({ key, value: val })
        seen.add(key)
      }
    })
  }

  // Append any handle coords that weren't already present in baseStyle.
  for (const [key, value] of Object.entries(fresh)) {
    if (!seen.has(key)) {
      parts.push({ key, value: String(value) })
      seen.add(key)
    }
  }

  return parts.map(p => (p.bare ? p.bare : `${p.key}=${p.value}`)).join(';') + ';'
}

/* ───────── DECODE ───────── */

/**
 * Convert API workflow + triggers into VueFlow nodes and edges
 */
export function decodeWorkflow(workflow, triggers = []) {
  const nodes = []
  const edges = []

  // 1. Triggers → trigger nodes + trigger→step edges
  triggers.forEach(({ meta, stepID, ...triggerConfig }) => {
    const vis = meta?.visual || {}
    const xywh = vis.xywh || [0, 0, 200, 80]

    nodes.push({
      id: String(vis.id),
      type: 'trigger',
      position: { x: xywh[0], y: xywh[1] },
      parentNode: vis.parent && vis.parent !== '1' ? String(vis.parent) : undefined,
      data: {
        stepID: String(vis.id),
        kind: 'trigger',
        ref: '',
        label: vis.value || meta?.name || '',
        // Fix #5: trigger meta.description round-trips via node.data.description
        description: meta?.description || '',
        arguments: [],
        results: [],
        defaultName: vis.defaultName || false,
        triggers: {
          // Spread full server-side trigger config so fields required by the
          // JS SDK update call (workflowID, ownedBy, updatedAt, …) survive the
          // round-trip. Without these, triggerUpdate() throws synchronously on
          // the second save and the user sees "configure-triggers".
          ...triggerConfig,
          triggerID: triggerConfig.triggerID,
          resourceType: triggerConfig.resourceType || null,
          eventType: triggerConfig.eventType || null,
          constraints: triggerConfig.constraints || [],
          enabled: triggerConfig.enabled !== false,
          stepID: stepID || '0',
        },
        highlighted: false,
        traceState: null,
        traceLog: null,
      },
    })

    // Trigger→step edge (from trigger.stepID, NOT from paths)
    if (stepID && stepID !== '0') {
      // Use trigger visual edges if available
      const triggerEdges = vis.edges || []
      triggerEdges.forEach(edge => {
        const eVis = edge?.meta?.visual || {}
        const handles = mxStyleToHandles(eVis.style)
        edges.push({
          id: String(eVis.id || `te-${vis.id}-${stepID}`),
          source: String(vis.id),
          target: String(edge.childID || stepID),
          sourceHandle: handles.sourceHandle || null,
          targetHandle: handles.targetHandle || null,
          type: 'workflow',
          label: eVis.value || '',
          data: {
            expr: '',
            description: edge?.meta?.description || '',
            parentID: String(vis.id),
            childID: String(edge.childID || stepID),
            highlighted: false,
            traceState: null,
            style: eVis.style || '',
            points: eVis.points || [],
          },
        })
      })

      // If no visual edges but stepID exists, create a virtual edge
      if (triggerEdges.length === 0) {
        edges.push({
          id: `te-${vis.id}-${stepID}`,
          source: String(vis.id),
          target: String(stepID),
          type: 'workflow',
          label: '',
          data: {
            expr: '',
            parentID: String(vis.id),
            childID: String(stepID),
            highlighted: false,
            traceState: null,
          },
        })
      }
    }
  })

  // 2. Steps → workflow / termination / visual nodes
  const steps = workflow.steps || []
  steps.forEach(({ stepID, kind, ref, meta, defaultName, arguments: args, results, ...rest }) => {
    const vis = meta?.visual || {}
    const xywh = vis.xywh || [0, 0, 200, 80]

    // Fix #18: silently upgrade the deprecated `workflow` kind to `exec-workflow`.
    // Human wrote this kind before the rename; data shape is identical, only
    // the string changed. We warn once per occurrence so the operator sees it.
    let stepKind = kind || ''
    if (stepKind === 'workflow') {
      console.warn(
        `[workflow codec] upgrading legacy step kind "workflow" → "exec-workflow" (stepID=${stepID})`,
      )
      stepKind = 'exec-workflow'
    }

    let nodeType = 'workflow'
    if (stepKind === 'termination') nodeType = 'termination'
    else if (stepKind === 'visual') nodeType = 'visual'

    nodes.push({
      id: String(vis.id || stepID),
      type: nodeType,
      position: { x: xywh[0], y: xywh[1] },
      zIndex: nodeType === 'visual' ? -1 : undefined,
      parentNode: vis.parent && vis.parent !== '1' ? String(vis.parent) : undefined,
      extent: undefined,
      style:
        nodeType === 'visual'
          ? { width: `${xywh[2] || 400}px`, height: `${xywh[3] || 240}px` }
          : undefined,
      data: {
        ...rest,
        stepID: String(stepID),
        kind: stepKind,
        ref: ref || '',
        label: vis.value || meta?.label || '',
        // Fix #5: persist meta.description through sidebar edits.
        description: meta?.description || '',
        arguments: args || [],
        results: results || [],
        defaultName: defaultName || vis.defaultName || false,
        highlighted: false,
        traceState: null,
        traceLog: null,
        width: xywh[2] || undefined,
        height: xywh[3] || undefined,
      },
    })
  })

  // 2b. Convert absolute → relative positions for children of visual parents.
  // Human stores xywh in absolute canvas coordinates; VueFlow interprets a
  // child's `position` as relative to its `parentNode`. Without this pass,
  // nodes nested in a swimlane jump to the wrong spot on decode.
  const nodeById = new Map(nodes.map(n => [n.id, n]))
  nodes.forEach(n => {
    if (!n.parentNode) return
    const parent = nodeById.get(n.parentNode)
    if (!parent) {
      // Parent missing — treat as top-level to avoid orphaned relative coords.
      n.parentNode = undefined
      n.extent = undefined
      return
    }
    n.position = {
      x: (n.position?.x || 0) - (parent.position?.x || 0),
      y: (n.position?.y || 0) - (parent.position?.y || 0),
    }
  })

  // 3. Paths → workflow edges
  const paths = workflow.paths || []
  paths.forEach(({ parentID, childID, meta, expr, ...pathConfig }) => {
    const vis = meta?.visual || {}
    const handles = mxStyleToHandles(vis.style)

    edges.push({
      id: String(vis.id || `p-${parentID}-${childID}`),
      source: String(parentID),
      target: String(childID),
      sourceHandle: handles.sourceHandle || null,
      targetHandle: handles.targetHandle || null,
      type: 'workflow',
      label: vis.value || meta?.label || '',
      data: {
        expr: expr || pathConfig.expr || '',
        description: meta?.description || '',
        parentID: String(parentID),
        childID: String(childID),
        highlighted: false,
        traceState: null,
        style: vis.style || '',
        points: vis.points || [],
      },
    })
  })

  return { nodes, edges }
}

/* ───────── ENCODE ───────── */

/**
 * Convert VueFlow nodes and edges back into API-shaped steps, paths, triggers
 */
export function encodeWorkflow(nodes, edges) {
  const steps = []
  const paths = []
  const triggers = []

  // Build a set of trigger node IDs for edge filtering
  const triggerNodeIds = new Set()
  nodes.forEach(n => {
    if (n.type === 'trigger' || n.data?.kind === 'trigger') {
      triggerNodeIds.add(n.id)
    }
  })

  // Resolve absolute positions: a child of a visual parent carries a position
  // relative to the parent, but Human expects absolute xywh. Walk parents
  // up to the root, summing offsets, and expose the result via `absPos`.
  const nodeById = new Map(nodes.map(n => [n.id, n]))
  function absolutePosition(node) {
    let x = node.position?.x || 0
    let y = node.position?.y || 0
    let cur = node.parentNode ? nodeById.get(node.parentNode) : null
    const guard = new Set()
    while (cur && !guard.has(cur.id)) {
      guard.add(cur.id)
      x += cur.position?.x || 0
      y += cur.position?.y || 0
      cur = cur.parentNode ? nodeById.get(cur.parentNode) : null
    }
    return { x, y }
  }

  // Encode nodes
  nodes.forEach(node => {
    const { data = {} } = node

    const abs = absolutePosition(node)

    if (node.type === 'trigger' || data.kind === 'trigger') {
      // Trigger node → trigger entry
      const outEdges = edges.filter(e => e.source === node.id)
      const firstTarget = outEdges.length > 0 ? outEdges[0].target : '0'

      // Build trigger visual edges (same format the render() method expects).
      // Fix #1: always re-emit canonical exit/entry coords derived from the
      // current handle IDs, but preserve any other properties present in the
      // original Human style string (orthogonalEdgeStyle, strokeColor, …).
      const triggerEdges = outEdges.map(e => {
        const baseStyle = e.data?.style || ''
        const style =
          e.sourceHandle || e.targetHandle
            ? handlesToMxStyle(e.sourceHandle, e.targetHandle, baseStyle)
            : baseStyle || handlesToMxStyle(null, null)
        return {
          parentID: e.source,
          childID: e.target,
          meta: {
            label: e.label || '',
            description: e.data?.description || '',
            visual: {
              id: e.id,
              value: e.label || '',
              parent: node.parentNode || '1',
              points: e.data?.points || [],
              style,
            },
          },
        }
      })

      triggers.push({
        ...(data.triggers || {}),
        stepID: firstTarget,
        enabled: data.triggers?.enabled !== false,
        constraints: data.triggers?.constraints || [],
        meta: {
          name: data.label || '',
          // Fix #5: trigger description persisted to server meta.description,
          // matches the WorkflowStepMeta.Description field on the Go side.
          description: data.description || '',
          visual: {
            id: node.id,
            value: data.label || '',
            defaultName: data.defaultName || false,
            xywh: [abs.x, abs.y, data.width || 200, data.height || 80],
            parent: node.parentNode || '1',
            edges: triggerEdges,
          },
        },
      })
    } else {
      // Regular step (workflow / termination / visual)
      const styleInfo = getStyleFromKind(data)
      const {
        highlighted: _highlighted,
        traceState: _traceState,
        traceLog: _traceLog,
        width: _width,
        height: _height,
        label: _label,
        description: _description,
        ...stepConfig
      } = data

      // Fix #18: encode only `exec-workflow`. If the in-memory model still
      // carries the legacy `workflow` kind (e.g. after an import that bypassed
      // decode), normalize it on the way out as well.
      let outKind = data.kind || ''
      if (outKind === 'workflow') {
        console.warn(
          `[workflow codec] encoding legacy step kind "workflow" as "exec-workflow" (stepID=${data.stepID || node.id})`,
        )
        outKind = 'exec-workflow'
      }

      steps.push({
        ...stepConfig,
        stepID: data.stepID || node.id,
        kind: outKind,
        ref: data.ref || '',
        defaultName: data.defaultName || false,
        arguments: data.arguments || [],
        results: data.results || [],
        meta: {
          label: data.label || '',
          // Fix #5: step description persisted to WorkflowStepMeta.Description
          description: data.description || '',
          visual: {
            id: node.id,
            value: data.label || '',
            defaultName: data.defaultName || false,
            xywh: [
              abs.x,
              abs.y,
              data.width || styleInfo.width || 200,
              data.height || styleInfo.height || 80,
            ],
            parent: node.parentNode || '1',
          },
        },
      })
    }
  })

  // Encode edges (exclude trigger edges, they are stored in trigger.meta.visual.edges)
  edges.forEach(edge => {
    // Skip edges that originate from a trigger node
    if (triggerNodeIds.has(edge.source)) return

    // Rebuild mxGraph style so anchor changes (user dragged an edge to a
    // different handle) survive the round-trip. When handles are known, re-emit
    // canonical exit/entry coords while preserving the loaded style string's
    // non-handle properties, so Human-written edges (orthogonalEdgeStyle,
    // strokeColor, rounded, …) aren't clobbered. With no handles set, keep the
    // original style so server-side geometry isn't lost.
    const baseStyle = edge.data?.style || ''
    const edgeStyle =
      edge.sourceHandle || edge.targetHandle
        ? handlesToMxStyle(edge.sourceHandle, edge.targetHandle, baseStyle)
        : baseStyle || handlesToMxStyle(null, null)

    paths.push({
      parentID: edge.source,
      childID: edge.target,
      expr: edge.data?.expr || '',
      meta: {
        label: edge.label || '',
        description: edge.data?.description || '',
        visual: {
          id: edge.id,
          value: edge.label || '',
          parent: '1',
          points: edge.data?.points || [],
          style: edgeStyle,
        },
      },
    })
  })

  return { steps, paths, triggers }
}

/* ─── Legacy compat (for WorkflowEditor.c3.js and any other callers) ─── */
export const encodeGraph = encodeWorkflow
