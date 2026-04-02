import type { automation } from '@cortezaproject/corteza-js-next'
import type { IconDef } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import {
  DEFAULT_ICONS,
  normalizeIcon,
} from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import type { Edge, Node } from '@vue-flow/core'
import dagre from 'dagre'

import type { AutomationFunction, AutomationTrigger } from '@/stores/automation'
import { getTriggerMeta, DEFAULT_TRIGGER_ICON } from '@/utils/flow-constants'

type NgAutomation = automation.NgAutomation
type NgAutomationTrigger = automation.NgAutomationTrigger
type NgAutomationStep = automation.NgAutomationStep
type NgAutomationPath = automation.NgAutomationPath
type Expr = automation.Expr
type TriggerConstraint = automation.TriggerConstraint

/**
 * Optional catalog data for resolving icons during conversion
 */
export interface ConversionCatalog {
  functions?: AutomationFunction[]
  triggers?: AutomationTrigger[]
}

/**
 * Node data structure for VueFlow nodes
 */
export interface FlowNodeData {
  label: string
  description?: string
  icon?: IconDef
  nodeType: string // The step/trigger type ID (e.g., 'webhook', 'http-request')
  config: Record<string, unknown> // Used by triggers (maps to trigger.input)
  arguments: Expr[] // Used by steps (maps to step.arguments)
  constraints?: TriggerConstraint[] // Used by triggers (maps to trigger.constraints)
  resourceType?: string // Trigger resource type (e.g., 'compose:record')
  ref: string // Reference back to automation step/trigger ref
  stepID?: string // Backend step ID
  triggerID?: string // Backend trigger ID
}

/**
 * VueFlow node types used in the builder
 */
export type FlowNodeType = 'trigger' | 'step' | 'branch' | 'iterator' | 'end' | 'loop'

/**
 * VueFlow state structure
 */
export interface VueFlowState {
  nodes: Node<FlowNodeData>[]
  edges: Edge[]
}

// Import centralized dimensions
import { NODE_DIMENSIONS } from './flow-constants'

// Layout constants (from centralized config)
const NODE_WIDTH = NODE_DIMENSIONS.WIDTH
const NODE_HEIGHT = NODE_DIMENSIONS.HEIGHT
const NODE_SEP = NODE_DIMENSIONS.HORIZONTAL_SEP
const RANK_SEP = NODE_DIMENSIONS.VERTICAL_SEP

/**
 * Convert NgAutomation (API format) to VueFlow state (frontend format)
 * Pass catalog to resolve icons from the construct library inline.
 */
export function automationToVueFlow(
  automation: NgAutomation,
  catalog?: ConversionCatalog,
): VueFlowState {
  const nodes: Node<FlowNodeData>[] = []
  const edges: Edge[] = []

  // Map clean IDs to VueFlow prefixed IDs for path resolution
  const idMap = new Map<string, string>()

  // Convert triggers to nodes
  automation.triggers?.forEach(trigger => {
    const cleanId = trigger.triggerID! || trigger.handle!
    const vueFlowId = `trigger_${cleanId}`
    idMap.set(cleanId, vueFlowId)

    nodes.push({
      id: vueFlowId,
      type: 'trigger',
      position: { x: 0, y: 0 }, // Will be set by layout
      data: {
        label: trigger.meta?.short || trigger.eventType || 'Trigger',
        description: trigger.meta?.description || '',
        icon:
          normalizeIcon(trigger.meta?.icon) ||
          getTriggerIcon(trigger.eventType, catalog, trigger.resourceType),
        nodeType: trigger.eventType || 'trigger',
        config: trigger.input || {},
        arguments: [],
        constraints: trigger.constraints || [],
        resourceType: trigger.resourceType || '',
        ref: trigger.handle || '',
        triggerID: trigger.triggerID,
      },
    })
  })

  // Convert steps to nodes
  automation.steps?.forEach(step => {
    const cleanId = step.stepID! || step.handle!
    const vueFlowId = `step_${cleanId}`
    idMap.set(cleanId, vueFlowId)

    const isCondition = step.kind?.startsWith('gateway')
    const isIterator = step.kind === 'iterator'
    const isTermination = step.kind === 'termination'

    let nodeType: FlowNodeType = 'step'
    if (isTermination) nodeType = 'end'
    else if (isCondition) nodeType = 'branch'
    else if (isIterator) nodeType = 'iterator'

    nodes.push({
      id: vueFlowId,
      type: nodeType,
      position: { x: 0, y: 0 }, // Will be set by layout
      selectable: !isTermination,
      data: {
        // Termination steps always display as "End" to user
        label: isTermination ? 'End' : step.meta?.short || step.ref || 'Step',
        description: step.meta?.description || '',
        icon: isTermination
          ? DEFAULT_ICONS.END
          : normalizeIcon(step.meta?.icon) || getStepIcon(step.ref, isCondition, catalog, isIterator),
        nodeType: isTermination ? 'termination' : step.ref,
        config: {},
        arguments: step.arguments || [],
        ref: step.handle || '',
        stepID: step.stepID,
      },
    })
  })

  const skippedBackEdges = new Set<string>()

  // Convert paths to edges (map clean IDs to VueFlow IDs)
  // No sourceHandle needed - edges are ordered by array position
  automation.paths?.forEach(path => {
    const sourceId = idMap.get(path.parentID) || path.parentID
    const targetId = idMap.get(path.childID) || path.childID

    // Detect back-edges (body chain leaf → iterator) and skip them.
    // A back-edge creates a cycle: the source is a descendant of the iterator
    // (reachable via the iterator's outgoing edges). Only these cyclic paths
    // are skipped; legitimate incoming edges (e.g., step → iterator) are kept.
    const targetNode = nodes.find(n => n.id === targetId)
    if (targetNode?.type === 'iterator') {
      // Check if source is reachable from the iterator (i.e., source is in the body chain)
      const isDescendant = (startId: string, searchId: string, visited = new Set<string>()): boolean => {
        if (startId === searchId) return true
        if (visited.has(startId)) return false
        visited.add(startId)
        // Follow already-added edges from startId
        return edges.some(e => e.source === startId && isDescendant(e.target, searchId, visited))
      }
      if (isDescendant(targetId, sourceId)) {
        // This is a genuine back-edge (cycle) — skip it for dagre
        skippedBackEdges.add(sourceId)
        return
      }
    }

    edges.push({
      id: path.handle || `${sourceId}_${targetId}`,
      source: sourceId,
      target: targetId,
      type: 'addable',
      data: {
        condition: path.condition || null,
      },
    })
  })

  // Count outgoing edges per node
  const outgoingEdgeCount = new Map<string, number>()
  edges.forEach(e => {
    outgoingEdgeCount.set(e.source, (outgoingEdgeCount.get(e.source) || 0) + 1)
  })

  // Add end nodes for leaf nodes
  nodes.forEach(node => {
    if (node.type === 'branch') {
      // Branches need at least 2 outputs - add end nodes as needed
      const currentCount = outgoingEdgeCount.get(node.id) || 0
      const endsNeeded = Math.max(0, 2 - currentCount)
      for (let i = 0; i < endsNeeded; i++) {
        const endId = `end_${node.id}_${i}`
        nodes.push({
          id: endId,
          type: 'end',
          position: { x: 0, y: 0 },
          selectable: false,
          data: {
            label: 'End',
            nodeType: 'termination',
            icon: DEFAULT_ICONS.END,
            config: {},
            arguments: [],
            ref: endId,
          },
        })
        edges.push({
          id: `${node.id}_${endId}`,
          source: node.id,
          target: endId,
          type: 'addable',
        })
      }
    } else if (node.type === 'iterator') {
      // Iterators need at least 2 outputs: body (first) + exit (second)
      const currentCount = outgoingEdgeCount.get(node.id) || 0
      const endsNeeded = Math.max(0, 2 - currentCount)
      for (let i = 0; i < endsNeeded; i++) {
        const isBodyEdge = currentCount === 0 && i === 0
        const endId = isBodyEdge ? `loop_${node.id}_${i}` : `end_${node.id}_${i}`
        const endType = isBodyEdge ? 'loop' : 'end'
        nodes.push({
          id: endId,
          type: endType,
          position: { x: 0, y: 0 },
          selectable: false,
          data: {
            label: isBodyEdge ? 'Loop' : 'End',
            nodeType: isBodyEdge ? 'loop' : 'termination',
            icon: isBodyEdge ? undefined : DEFAULT_ICONS.END,
            config: {},
            arguments: [],
            ref: endId,
          },
        })
        edges.push({
          id: `${node.id}_${endId}`,
          source: node.id,
          target: endId,
          type: 'addable',
        })
      }
    } else if (!outgoingEdgeCount.has(node.id) && node.type !== 'end' && node.type !== 'loop') {
      // Regular leaf nodes get one end node, but iterator body tips get a loop node
      const isBodyTip = skippedBackEdges.has(node.id)
      const endId = isBodyTip ? `loop_${node.id}` : `end_${node.id}`
      const type = isBodyTip ? 'loop' : 'end'
      nodes.push({
        id: endId,
        type: type,
        position: { x: 0, y: 0 },
        selectable: false,
        data: {
          label: isBodyTip ? 'Loop' : 'End',
          nodeType: isBodyTip ? 'loop' : 'termination',
          icon: isBodyTip ? undefined : DEFAULT_ICONS.END,
          config: {},
          ref: endId,
        },
      })
      edges.push({
        id: `${node.id}_${endId}`,
        source: node.id,
        target: endId,
        type: 'addable',
      })
    }
  })

  // Apply automatic layout
  return applyDagreLayout({ nodes, edges })
}

/**
 * Convert VueFlow state (frontend) to NgAutomation data (for API)
 * End nodes are saved as termination steps
 */
export function vueFlowToAutomation(
  state: VueFlowState,
): Pick<NgAutomation, 'triggers' | 'steps' | 'paths'> {
  const triggers: NgAutomationTrigger[] = []
  const steps: NgAutomationStep[] = []
  const paths: NgAutomationPath[] = []

  let triggerIndex = 0
  let stepIndex = 0
  let pathIndex = 0

  // Map node IDs to their new stepID/triggerID for path generation
  const nodeIdToStepId = new Map<string, string>()

  // Convert nodes to triggers/steps
  state.nodes.forEach(node => {
    if (node.type === 'loop') return

    if (node.type === 'trigger') {
      triggerIndex++
      const data = node.data as FlowNodeData
      const newTriggerID = String(triggerIndex)
      nodeIdToStepId.set(node.id, newTriggerID)
      triggers.push({
        triggerID: newTriggerID,
        handle: `trigger_${triggerIndex}`,
        enabled: true,
        resourceType: data.resourceType || '',
        eventType: data.nodeType,
        constraints: data.constraints || [],
        meta: {
          short: data.label,
          description: data.description || '',
          icon: data.icon,
        },
        input: data.config || {},
      })
    } else {
      // All non-trigger nodes become steps (including end nodes as termination)
      stepIndex++
      const data = node.data as FlowNodeData
      const newStepID = String(stepIndex)
      nodeIdToStepId.set(node.id, newStepID)

      // Determine step kind
      let kind: string = 'function'
      if (node.type === 'end') {
        kind = 'termination'
      } else if (node.type === 'branch') {
        kind = 'gateway'
      } else if (node.type === 'iterator') {
        kind = 'iterator'
      }

      steps.push({
        stepID: newStepID,
        handle: `step_${stepIndex}`,
        kind,
        ref: node.type === 'end' ? 'termination' : data.nodeType,
        meta: {
          short: data.label,
          description: data.description || '',
          icon: data.icon,
        },
        arguments: data.arguments || [],
      })
    }
  })

  // Convert edges to paths (including paths to termination steps)
  state.edges.forEach(edge => {
    const sourceId = nodeIdToStepId.get(edge.source)
    const targetId = nodeIdToStepId.get(edge.target)

    if (!sourceId || !targetId) return

    pathIndex++
    const condition = edge.data?.condition || null
    paths.push({
      parentID: sourceId,
      childID: targetId,
      handle: `path_${pathIndex}`,
      ...(condition ? { condition } : {}),
      meta: {
        short: condition ? conditionToShort(condition) : '',
      },
    })
  })

  return { triggers, steps, paths }
}

/**
 * Apply dagre layout to position nodes automatically
 */
export function applyDagreLayout(state: VueFlowState): VueFlowState {
  const g = new dagre.graphlib.Graph()
  g.setDefaultEdgeLabel(() => ({}))
  // Configure graph: TB = top-to-bottom, align undefined = center alignment
  g.setGraph({ rankdir: 'TB', nodesep: NODE_SEP, ranksep: RANK_SEP, align: undefined })

  // Add nodes to graph
  state.nodes.forEach(node => {
    const height = node.type === 'end' ? 40 : NODE_HEIGHT
    g.setNode(node.id, { width: NODE_WIDTH, height })
  })

  // Add edges to graph
  state.edges.forEach(edge => {
    g.setEdge(edge.source, edge.target)
  })

  // Run layout
  dagre.layout(g)

  // Apply positions
  const layoutedNodes = state.nodes.map(node => {
    const nodeWithPosition = g.node(node.id)
    const height = node.type === 'end' ? 40 : NODE_HEIGHT
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - NODE_WIDTH / 2,
        y: nodeWithPosition.y - height / 2,
      },
    }
  })

  // Post-process: Align branch/iterator children's Y positions and ensure minimum End node distance
  const forkedNodes = layoutedNodes.filter(n => n.type === 'branch' || n.type === 'iterator')

  // Helper to get all descendants of a node
  const getDescendants = (nodeId: string, visited = new Set<string>()): string[] => {
    if (visited.has(nodeId)) return []
    visited.add(nodeId)

    const childEdges = state.edges.filter(e => e.source === nodeId)
    const descendants: string[] = []

    childEdges.forEach(edge => {
      descendants.push(edge.target)
      descendants.push(...getDescendants(edge.target, visited))
    })

    return descendants
  }

  forkedNodes.forEach(forkedNode => {
    const forkedEdges = state.edges.filter(e => e.source === forkedNode.id)
    if (forkedEdges.length < 2) return

    const childNodes = forkedEdges
      .map(e => layoutedNodes.find(n => n.id === e.target))
      .filter(Boolean) as typeof layoutedNodes

    if (childNodes.length < 2) return

    // REORDER X POSITIONS: Ensure children are positioned left-to-right matching edge array order
    const xPositions = childNodes.map(n => n.position.x).sort((a, b) => a - b)

    forkedEdges.forEach((edge, index) => {
      const childNode = childNodes[index]
      if (!childNode) return

      const targetX = xPositions[index]
      const xDiff = targetX - childNode.position.x

      if (xDiff !== 0) {
        const descendants = [edge.target, ...getDescendants(edge.target)]
        descendants.forEach(id => {
          const node = layoutedNodes.find(n => n.id === id)
          if (node) node.position.x += xDiff
        })
      }
    })

    // Align all children to the same Y level (use the deepest one)
    const maxY = Math.max(...childNodes.map(n => n.position.y))
    forkedEdges.forEach((edge, index) => {
      const childNode = childNodes[index]
      if (!childNode) return

      const yDiff = maxY - childNode.position.y
      if (yDiff > 0) {
        const descendants = [edge.target, ...getDescendants(edge.target)]
        descendants.forEach(id => {
          const node = layoutedNodes.find(n => n.id === id)
          if (node) node.position.y += yDiff
        })
      }
    })
  })

  // Ensure End nodes directly connected to branches/iterators have consistent minimum distance
  const MIN_END_DISTANCE = 120
  forkedNodes.forEach(forkedNode => {
    const forkedEdges = state.edges.filter(e => e.source === forkedNode.id)

    forkedEdges.forEach(edge => {
      const targetNode = layoutedNodes.find(n => n.id === edge.target)
      if (targetNode && targetNode.type === 'end') {
        const currentDistance = targetNode.position.y - forkedNode.position.y
        if (currentDistance < MIN_END_DISTANCE) {
          const adjustment = MIN_END_DISTANCE - currentDistance
          targetNode.position.y += adjustment
        }
      }
    })
  })

  return { nodes: layoutedNodes, edges: state.edges }
}

/**
 * Generate a unique ref for new nodes
 */
export function generateRef(_nodeType: string, existingRefs: string[]): string {
  let maxId = 0
  for (const ref of existingRefs) {
    const num = parseInt(ref, 10)
    if (!isNaN(num) && num > maxId) {
      maxId = num
    }
  }
  return String(maxId + 1)
}

/**
 * Get all refs from VueFlow state
 */
export function getAllRefs(state: VueFlowState): string[] {
  return state.nodes.filter(n => n.type !== 'end').map(n => (n.data as FlowNodeData).ref || n.id)
}

// Helper: Resolve trigger icon from catalog, TRIGGER_META, or fallback
function getTriggerIcon(
  eventType?: string,
  catalog?: ConversionCatalog,
  resourceType?: string,
): IconDef {
  if (catalog?.triggers && eventType) {
    const catalogTrigger = catalog.triggers.find(t => t.eventType === eventType)
    const catalogIcon = normalizeIcon(catalogTrigger?.meta?.icon)
    if (catalogIcon) return catalogIcon
  }

  if (eventType) {
    const meta = getTriggerMeta(eventType, resourceType)
    if (meta?.icon) return meta.icon
  }

  return DEFAULT_TRIGGER_ICON
}

// Helper: Resolve step icon from catalog or fallback
function getStepIcon(ref?: string, isCondition?: boolean, catalog?: ConversionCatalog, isIterator?: boolean): IconDef {
  if (isCondition) return DEFAULT_ICONS.BRANCH

  // Check catalog for function-specific icon (both regular functions and iterators)
  if (catalog?.functions && ref) {
    const catalogFn = catalog.functions.find(f => f.ref === ref)
    const catalogIcon = normalizeIcon(catalogFn?.meta?.icon)
    if (catalogIcon) return catalogIcon
  }

  if (isIterator) return DEFAULT_ICONS.ITERATOR
  return DEFAULT_ICONS.ACTION
}

// Operator labels for human-readable condition summaries
const OP_LABELS: Record<string, string> = {
  eq: 'equals',
  ne: 'not equals',
  lt: 'less than',
  lte: 'at most',
  gt: 'greater than',
  gte: 'at least',
  isNull: 'is empty',
  isNotNull: 'is not empty',
  and: 'AND',
  or: 'OR',
}

/**
 * Convert an ASTNode condition to a human-readable short string for meta.short
 */
export function conditionToShort(node: Record<string, unknown>): string {
  if (!node) return ''

  // Leaf: value
  if (node.value && typeof node.value === 'object') {
    const v = node.value as Record<string, unknown>
    return String(v['@value'] ?? '')
  }

  // Leaf: symbol
  if (node.symbol) return String(node.symbol)

  const ref = String(node.ref || '')
  const args = (node.args || []) as Record<string, unknown>[]

  // Unary operator (isNull/isNotNull)
  if (ref === 'isNull' || ref === 'isNotNull') {
    return `${conditionToShort(args[0])} ${OP_LABELS[ref] || ref}`
  }

  // Combinator (and/or)
  if (ref === 'and' || ref === 'or') {
    const sep = ` ${OP_LABELS[ref]} `
    return args.map(a => conditionToShort(a)).join(sep)
  }

  // Binary operator
  if (args.length >= 2) {
    const left = conditionToShort(args[0])
    const right = conditionToShort(args[1])
    return `${left} ${OP_LABELS[ref] || ref} ${right}`
  }

  return ref
}

export type ConditionSegment =
  | { type: 'text'; value: string }
  | { type: 'ref'; value: string; scope: string }

/**
 * Convert an ASTNode condition to an array of segments for rich rendering.
 * Symbol nodes become { type: 'ref' }, everything else becomes { type: 'text' }.
 */
export function conditionToSegments(node: Record<string, unknown>): ConditionSegment[] {
  if (!node) return []

  // Leaf: value
  if (node.value && typeof node.value === 'object') {
    const v = node.value as Record<string, unknown>
    return [{ type: 'text', value: String(v['@value'] ?? '') }]
  }

  // Leaf: symbol
  if (node.symbol) {
    const meta = (node.meta || {}) as Record<string, string>
    return [{ type: 'ref', value: String(node.symbol), scope: meta.scope || '' }]
  }

  const ref = String(node.ref || '')
  const args = (node.args || []) as Record<string, unknown>[]

  // Unary operator
  if (ref === 'isNull' || ref === 'isNotNull') {
    return [...conditionToSegments(args[0]), { type: 'text', value: ` ${OP_LABELS[ref] || ref}` }]
  }

  // Combinator (and/or)
  if (ref === 'and' || ref === 'or') {
    const sep: ConditionSegment = { type: 'text', value: ` ${OP_LABELS[ref]} ` }
    const result: ConditionSegment[] = []
    args.forEach((a, i) => {
      if (i > 0) result.push(sep)
      result.push(...conditionToSegments(a))
    })
    return result
  }

  // Binary operator
  if (args.length >= 2) {
    return [
      ...conditionToSegments(args[0]),
      { type: 'text', value: ` ${OP_LABELS[ref] || ref} ` },
      ...conditionToSegments(args[1]),
    ]
  }

  return [{ type: 'text', value: ref }]
}
