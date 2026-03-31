/**
 * Workflow codec — serializes between VueFlow nodes/edges and the API model.
 *
 * decodeWorkflow  : (workflow, triggers) → { nodes, edges }
 * encodeWorkflow  : (nodes, edges)       → { steps, paths, triggers }
 */
import { getStyleFromKind, getKindFromStyle } from './style'

/**
 * Map mxGraph edge style string to VueFlow sourceHandle / targetHandle IDs.
 *
 * mxGraph stores exit/entry positions like:
 *   "exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;"
 *
 * Available VueFlow handles per node:
 *   Source: source-bottom (50%), source-bottom-left (25%), source-bottom-right (75%),
 *           source-right, source-left
 *   Target: target-top (50%), target-top-left (25%), target-top-right (75%),
 *           target-left, target-bottom
 */
function mxStyleToHandles (style) {
  if (!style) return {}
  const exitX = parseFloat(style.match(/exitX=([0-9.]+)/)?.[1] ?? -1)
  const exitY = parseFloat(style.match(/exitY=([0-9.]+)/)?.[1] ?? -1)
  const entryX = parseFloat(style.match(/entryX=([0-9.]+)/)?.[1] ?? -1)
  const entryY = parseFloat(style.match(/entryY=([0-9.]+)/)?.[1] ?? -1)
  const result = {}

  // Source handle mapping (exit position)
  if (exitY >= 0) {
    if (exitY <= 0.25) {
      // Top exit → treat as right (VueFlow has no top source handle)
      result.sourceHandle = 'source-right'
    } else if (exitY >= 0.75) {
      // Bottom exit
      if (exitX <= 0.35) result.sourceHandle = 'source-bottom-left'
      else if (exitX >= 0.65) result.sourceHandle = 'source-bottom-right'
      else result.sourceHandle = 'source-bottom'
    } else {
      // Side exit (exitY ≈ 0.5)
      if (exitX >= 0.5) result.sourceHandle = 'source-right'
      else result.sourceHandle = 'source-left'
    }
  }

  // Target handle mapping (entry position)
  if (entryY >= 0) {
    if (entryY <= 0.25) {
      // Top entry
      if (entryX <= 0.35) result.targetHandle = 'target-top-left'
      else if (entryX >= 0.65) result.targetHandle = 'target-top-right'
      else result.targetHandle = 'target-top'
    } else if (entryY >= 0.75) {
      // Bottom entry
      result.targetHandle = 'target-bottom'
    } else {
      // Side entry (entryY ≈ 0.5)
      result.targetHandle = 'target-left'
    }
  }

  return result
}

/** Handle ID → mxGraph exit/entry positions */
const HANDLE_TO_MX = {
  'source-bottom': { exitX: 0.5, exitY: 1 },
  'source-bottom-left': { exitX: 0.25, exitY: 1 },
  'source-bottom-right': { exitX: 0.75, exitY: 1 },
  'source-right': { exitX: 1, exitY: 0.5 },
  'source-left': { exitX: 0, exitY: 0.5 },
  'target-top': { entryX: 0.5, entryY: 0 },
  'target-top-left': { entryX: 0.25, entryY: 0 },
  'target-top-right': { entryX: 0.75, entryY: 0 },
  'target-left': { entryX: 0, entryY: 0.5 },
  'target-bottom': { entryX: 0.5, entryY: 1 },
}

/**
 * Reverse mapping: convert VueFlow handle IDs to mxGraph style string.
 * Used when encoding newly-drawn edges that have no original mxGraph style.
 */
function handlesToMxStyle (sourceHandle, targetHandle) {
  const src = HANDLE_TO_MX[sourceHandle] || HANDLE_TO_MX['source-bottom']
  const tgt = HANDLE_TO_MX[targetHandle] || HANDLE_TO_MX['target-top']
  const exitX = src.exitX ?? 0.5
  const exitY = src.exitY ?? 1
  const entryX = tgt.entryX ?? 0.5
  const entryY = tgt.entryY ?? 0
  return `exitX=${exitX};exitY=${exitY};exitDx=0;exitDy=0;entryX=${entryX};entryY=${entryY};entryDx=0;entryDy=0;`
}

/* ───────── DECODE ───────── */

/**
 * Convert API workflow + triggers into VueFlow nodes and edges
 */
export function decodeWorkflow (workflow, triggers = []) {
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
      data: {
        stepID: String(vis.id),
        kind: 'trigger',
        ref: '',
        label: vis.value || meta?.name || '',
        description: '',
        arguments: [],
        results: [],
        defaultName: vis.defaultName || false,
        triggers: {
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
      triggerEdges.forEach((edge) => {
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
  steps.forEach(({ stepID, kind, ref, meta, defaultName, arguments: args, results }) => {
    const vis = meta?.visual || {}
    const xywh = vis.xywh || [0, 0, 200, 80]

    let nodeType = 'workflow'
    if (kind === 'termination') nodeType = 'termination'
    else if (kind === 'visual') nodeType = 'visual'

    nodes.push({
      id: String(vis.id || stepID),
      type: nodeType,
      position: { x: xywh[0], y: xywh[1] },
      zIndex: nodeType === 'visual' ? -1 : undefined,
      parentNode: vis.parent && vis.parent !== '1' ? String(vis.parent) : undefined,
      data: {
        stepID: String(stepID),
        kind: kind || '',
        ref: ref || '',
        label: vis.value || meta?.label || '',
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
export function encodeWorkflow (nodes, edges) {
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

  // Encode nodes
  nodes.forEach(node => {
    const { data = {} } = node

    if (node.type === 'trigger' || data.kind === 'trigger') {
      // Trigger node → trigger entry
      const outEdges = edges.filter(e => e.source === node.id)
      const firstTarget = outEdges.length > 0 ? outEdges[0].target : '0'

      // Build trigger visual edges (same format the render() method expects)
      const triggerEdges = outEdges.map(e => ({
        parentID: e.source,
        childID: e.target,
        meta: {
          label: e.label || '',
          description: '',
          visual: {
            id: e.id,
            value: e.label || '',
            parent: node.parentNode || '1',
            points: e.data?.points || [],
            style: e.data?.style || '',
          },
        },
      }))

      triggers.push({
        ...(data.triggers || {}),
        stepID: firstTarget,
        enabled: data.triggers?.enabled !== false,
        constraints: data.triggers?.constraints || [],
        meta: {
          name: data.label || '',
          description: data.description || '',
          visual: {
            id: node.id,
            value: data.label || '',
            defaultName: data.defaultName || false,
            xywh: [
              node.position?.x || 0,
              node.position?.y || 0,
              data.width || 200,
              data.height || 80,
            ],
            parent: node.parentNode || '1',
            edges: triggerEdges,
          },
        },
      })
    } else {
      // Regular step (workflow / termination / visual)
      const styleInfo = getStyleFromKind(data)
      steps.push({
        stepID: data.stepID || node.id,
        kind: data.kind || '',
        ref: data.ref || '',
        defaultName: data.defaultName || false,
        arguments: data.arguments || [],
        results: data.results || [],
        meta: {
          label: data.label || '',
          description: data.description || '',
          visual: {
            id: node.id,
            value: data.label || '',
            defaultName: data.defaultName || false,
            xywh: [
              node.position?.x || 0,
              node.position?.y || 0,
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

    // Reconstruct mxGraph style from VueFlow handles + original style
    const edgeStyle = edge.data?.style || handlesToMxStyle(edge.sourceHandle, edge.targetHandle)

    paths.push({
      parentID: edge.source,
      childID: edge.target,
      expr: edge.data?.expr || '',
      meta: {
        label: edge.label || '',
        description: '',
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
