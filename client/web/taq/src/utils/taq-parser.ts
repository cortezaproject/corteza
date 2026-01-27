import type { automation } from '@cortezaproject/corteza-js-next'
import type { Edge, Node } from '@vue-flow/core'
import dagre from 'dagre'

type NgAutomation = automation.NgAutomation
type NgAutomationTrigger = automation.NgAutomationTrigger
type NgAutomationStep = automation.NgAutomationStep
type NgAutomationPath = automation.NgAutomationPath

/**
 * Node data structure for VueFlow nodes
 */
export interface FlowNodeData {
  label: string
  description?: string
  icon?: string
  nodeType: string // The step/trigger type ID (e.g., 'webhook', 'http-request')
  config: Record<string, unknown>
  ref: string // Reference back to automation step/trigger ref
  stepID?: string // Backend step ID
  triggerID?: string // Backend trigger ID
}

/**
 * VueFlow node types used in the builder
 */
export type FlowNodeType = 'trigger' | 'step' | 'branch' | 'end'

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
 */
export function automationToVueFlow(automation: NgAutomation): VueFlowState {
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
        icon: getTriggerIcon(trigger.eventType),
        nodeType: trigger.eventType || 'trigger',
        config: trigger.input || {},
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

    const isCondition = step.kind === 'gateway'
    const isTermination = step.kind === 'termination'
    nodes.push({
      id: vueFlowId,
      type: isTermination ? 'end' : isCondition ? 'branch' : 'step',
      position: { x: 0, y: 0 }, // Will be set by layout
      selectable: !isTermination,
      data: {
        // Termination steps always display as "End" to user
        label: isTermination ? 'End' : step.meta?.short || step.ref || 'Step',
        description: step.meta?.description || '',
        icon: isTermination ? 'pi pi-stop-circle' : getStepIcon(step.ref, isCondition),
        nodeType: isTermination ? 'termination' : step.ref,
        config: argumentsToConfig(step.arguments),
        ref: step.handle || '',
        stepID: step.stepID,
      },
    })
  })

  // Convert paths to edges (map clean IDs to VueFlow IDs)
  // No sourceHandle needed - edges are ordered by array position
  automation.paths?.forEach(path => {
    const sourceId = idMap.get(path.parentID) || path.parentID
    const targetId = idMap.get(path.childID) || path.childID

    edges.push({
      id: path.handle || `${sourceId}_${targetId}`,
      source: sourceId,
      target: targetId,
      type: 'addable',
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
            icon: 'pi pi-stop-circle',
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
    } else if (!outgoingEdgeCount.has(node.id) && node.type !== 'end') {
      // Regular leaf nodes get one end node
      const endId = `end_${node.id}`
      nodes.push({
        id: endId,
        type: 'end',
        position: { x: 0, y: 0 },
        selectable: false,
        data: {
          label: 'End',
          nodeType: 'termination',
          icon: 'pi pi-stop-circle',
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
    if (node.type === 'trigger') {
      triggerIndex++
      const data = node.data as FlowNodeData
      const newTriggerID = String(triggerIndex)
      nodeIdToStepId.set(node.id, newTriggerID)
      triggers.push({
        triggerID: newTriggerID,
        handle: `trigger_${triggerIndex}`,
        enabled: true,
        resourceType: '',
        eventType: data.nodeType,
        meta: {
          short: data.label,
          description: data.description || '',
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
      }

      steps.push({
        stepID: newStepID,
        handle: `step_${stepIndex}`,
        kind,
        ref: node.type === 'end' ? 'termination' : data.nodeType,
        meta: {
          short: data.label,
          description: data.description || '',
        },
        arguments: configToArguments(data.config || {}) as any,
      })
    }
  })

  // Convert edges to paths (including paths to termination steps)
  state.edges.forEach(edge => {
    const sourceId = nodeIdToStepId.get(edge.source)
    const targetId = nodeIdToStepId.get(edge.target)

    if (!sourceId || !targetId) return

    pathIndex++
    paths.push({
      parentID: sourceId,
      childID: targetId,
      handle: `path_${pathIndex}`,
      meta: {
        expr: edge.sourceHandle || undefined,
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

  // Post-process: Align branch children's Y positions and ensure minimum End node distance
  const branchNodes = layoutedNodes.filter(n => n.type === 'branch')

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

  branchNodes.forEach(branchNode => {
    // Get all edges from this branch (already ordered by array position)
    const branchEdges = state.edges.filter(e => e.source === branchNode.id)
    if (branchEdges.length < 2) return

    // Get direct children in edge array order
    const childNodes = branchEdges
      .map(e => layoutedNodes.find(n => n.id === e.target))
      .filter(Boolean) as typeof layoutedNodes

    if (childNodes.length < 2) return

    // REORDER X POSITIONS: Ensure children are positioned left-to-right matching edge array order
    // Sort children by their current X position to find the available slots
    const xPositions = childNodes.map(n => n.position.x).sort((a, b) => a - b)

    // Assign X positions based on edge order (first edge = leftmost position)
    branchEdges.forEach((edge, index) => {
      const childNode = childNodes[index]
      if (!childNode) return

      const targetX = xPositions[index]
      const xDiff = targetX - childNode.position.x

      if (xDiff !== 0) {
        // Move entire subtree
        const descendants = [edge.target, ...getDescendants(edge.target)]
        descendants.forEach(id => {
          const node = layoutedNodes.find(n => n.id === id)
          if (node) node.position.x += xDiff
        })
      }
    })

    // Align all children to the same Y level (use the deepest one)
    const maxY = Math.max(...childNodes.map(n => n.position.y))
    branchEdges.forEach((edge, index) => {
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

  // Ensure End nodes directly connected to branches have consistent minimum distance
  const MIN_END_DISTANCE = 120
  branchNodes.forEach(branchNode => {
    const branchEdges = state.edges.filter(e => e.source === branchNode.id)

    branchEdges.forEach(edge => {
      const targetNode = layoutedNodes.find(n => n.id === edge.target)
      if (targetNode && targetNode.type === 'end') {
        const currentDistance = targetNode.position.y - branchNode.position.y
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

// Helper: Convert config object to arguments array
function configToArguments(
  config: Record<string, unknown>,
): Array<{ target: string; type: string; value?: unknown }> {
  return Object.entries(config).map(([target, value]) => ({
    target,
    type: inferValueType(value),
    value,
  }))
}

// Helper: Convert arguments array to config object
function argumentsToConfig(
  args?: Array<{ target?: string; value?: unknown }>,
): Record<string, unknown> {
  if (!args) return {}
  const config: Record<string, unknown> = {}
  args.forEach(arg => {
    if (arg.target) {
      config[arg.target] = arg.value
    }
  })
  return config
}

// Helper: Infer value type from JavaScript value
function inferValueType(value: unknown): string {
  if (value === null || value === undefined) return 'Any'
  if (typeof value === 'string') return 'String'
  if (typeof value === 'number') return 'Number'
  if (typeof value === 'boolean') return 'Boolean'
  if (Array.isArray(value)) return 'Array'
  if (typeof value === 'object') return 'KV'
  return 'Any'
}

// Helper: Get icon for trigger type
function getTriggerIcon(eventType?: string): string {
  const icons: Record<string, string> = {
    webhook: 'pi pi-globe',
    schedule: 'pi pi-clock',
    manual: 'pi pi-play',
    onRecord: 'pi pi-database',
  }
  return icons[eventType || ''] || 'pi pi-bolt'
}

// Helper: Get icon for step type
function getStepIcon(ref?: string, isCondition?: boolean): string {
  if (isCondition) return 'pi pi-sitemap'
  const icons: Record<string, string> = {
    'http-request': 'pi pi-globe',
    'send-email': 'pi pi-envelope',
    delay: 'pi pi-clock',
    log: 'pi pi-file',
    javascript: 'pi pi-code',
  }
  return icons[ref || ''] || 'pi pi-cog'
}
