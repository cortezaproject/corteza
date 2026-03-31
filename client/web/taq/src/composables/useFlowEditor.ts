import { automation } from '@cortezaproject/corteza-js-next'
import type { StackFrame, ExecutionResult, TraceStatus } from '@cortezaproject/corteza-js-next/src/automation/types/trace'
import { withMinDuration } from '@cortezaproject/corteza-vue-next'
import type { IconDef } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import type { Edge, Node } from '@vue-flow/core'
import { computed, inject, nextTick, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useAutomationStore } from '@/stores/automation'
import { applyDagreLayout, automationToVueFlow, type FlowNodeData } from '@/utils/taq-parser'

const { NgAutomation } = automation
type NgAutomationInstance = InstanceType<typeof NgAutomation>

interface InsertionPoint {
  first?: boolean
  edgeId?: string
  source?: string
  target?: string
}

interface NodeType {
  id: string
  type: string
  label: string
  description?: string
  icon?: IconDef
  ref?: string
  kind?: string
  eventType?: string
  resourceType?: string
}

/**
 * Composable for managing the flow editor state and operations.
 * Handles automation loading/saving, node manipulation, and undo/redo.
 */
export function useFlowEditor() {
  const $AutomationAPI = inject<any>('$AutomationAPI')
  const $toast = inject<any>('$toast')
  const { t } = useI18n()
  const router = useRouter()

  // Core automation state
  const automation = ref<NgAutomationInstance>(new NgAutomation())
  const loading = ref(false)
  const saving = ref(false)
  const running = ref(false)

  // Trace state
  const traceFrames = ref<StackFrame[]>([])
  const traceExecution = ref<ExecutionResult | null>(null)
  const traceStatus = ref<TraceStatus>('idle')

  // VueFlow state
  const nodes = ref<Node<FlowNodeData>[]>([])
  const edges = ref<Edge[]>([])

  // History for undo/redo
  const history = ref<string[]>([])
  const historyIndex = ref(-1)

  // Computed
  const canUndo = computed(() => historyIndex.value > 0)
  const canRedo = computed(() => historyIndex.value < history.value.length - 1)
  const automationId = computed(() => automation.value.automationID)
  const isEmpty = computed(() => nodes.value.length === 0)

  const name = computed({
    get: () => automation.value.meta?.short || t('builder.newTaq'),
    set: (val: string) => {
      automation.value.meta = { ...automation.value.meta, short: val }
    },
  })

  const enabled = computed({
    get: () => automation.value.enabled ?? false,
    set: (val: boolean) => {
      automation.value.enabled = val
    },
  })

  // History management
  function saveToHistory() {
    const state = JSON.stringify({ nodes: nodes.value, edges: edges.value })
    history.value = history.value.slice(0, historyIndex.value + 1)
    history.value.push(state)
    historyIndex.value = history.value.length - 1
  }

  function undo() {
    if (!canUndo.value) return
    historyIndex.value--
    const state = JSON.parse(history.value[historyIndex.value])
    nodes.value = state.nodes
    edges.value = state.edges
  }

  function redo() {
    if (!canRedo.value) return
    historyIndex.value++
    const state = JSON.parse(history.value[historyIndex.value])
    nodes.value = state.nodes
    edges.value = state.edges
  }

  // Load automation from API
  async function load(id: string) {
    loading.value = true
    try {
      const response = await withMinDuration($AutomationAPI.ngAutomationRead({ automationID: id }))
      automation.value = new NgAutomation(response)

      // Convert to VueFlow format (catalog resolves icons inline)
      const store = useAutomationStore()
      const state = automationToVueFlow(automation.value, {
        functions: store.functions,
        triggers: store.triggers,
      })
      nodes.value = state.nodes
      edges.value = state.edges

      nextTick(() => saveToHistory())
    } catch (e) {
      console.error('Failed to load automation:', e)
      $toast?.toastDanger(t('builder.toast.loadError.detail'), t('builder.toast.loadError.summary'))
      throw e
    } finally {
      loading.value = false
    }
  }

  // Save automation to API
  async function save() {
    saving.value = true
    try {
      // Build triggers and steps from nodes while preserving existing data
      const triggers: any[] = []
      const steps: any[] = []

      // Use a SINGLE counter for all IDs to avoid collisions in paths
      // Triggers and steps share the ID space for path parentID/childID
      // Start at 0, increment before use so IDs are 1, 2, 3, ...
      let globalIndex = 0

      // Map to track VueFlow node ID -> backend ID for path generation
      const nodeIdToBackendId = new Map<string, string>()

      // Helper to process a node into trigger or step
      const processNode = (node: Node<FlowNodeData>) => {
        const data = node.data as FlowNodeData
        globalIndex++
        const newID = String(globalIndex)
        nodeIdToBackendId.set(node.id, newID)

        if (node.type === 'trigger') {
          const existing = automation.value.triggers?.find(
            (t: any) => t.triggerID === data.triggerID,
          )
          triggers.push({
            triggerID: newID,
            handle: `trigger_${newID}`,
            enabled: existing?.enabled ?? true,
            resourceType: data.resourceType || existing?.resourceType || '',
            eventType: data.nodeType || existing?.eventType || '',
            constraints: data.constraints || existing?.constraints || [],
            meta: {
              short: data.label,
              description: data.description || '',
            },
            input: data.config || existing?.input || {},
          })
        } else {
          const existing = automation.value.steps?.find((s: any) => s.stepID === data.stepID)

          let kind = 'function'
          if (node.type === 'end') {
            kind = 'termination'
          } else if (node.type === 'branch') {
            kind = data.nodeType || 'gatewayExclusive'
          } else if (node.type === 'iterator') {
            kind = 'iterator'
          }

          steps.push({
            stepID: newID,
            handle: `step_${newID}`,
            kind,
            ref: node.type === 'end' ? 'termination' : data.nodeType || existing?.ref || '',
            meta: {
              short: node.type === 'end' ? t('builder.nodes.end') : data.label,
              description: data.description || '',
            },
            arguments: data.arguments || existing?.arguments || [],
          })
        }
      }

      // Process in order: 1) triggers, 2) terminations, 3) other steps
      nodes.value.filter(n => n.type === 'trigger').forEach(processNode)
      nodes.value.filter(n => n.type === 'end').forEach(processNode)
      nodes.value.filter(n => n.type !== 'trigger' && n.type !== 'end').forEach(processNode)

      // Derive paths from edges (including paths to termination steps)
      // Filter out any invalid or duplicate edges
      const seenEdges = new Set<string>()
      const paths = edges.value
        .map(edge => {
          const sourceId = nodeIdToBackendId.get(edge.source)
          const targetId = nodeIdToBackendId.get(edge.target)

          if (!sourceId || !targetId) return null

          // Skip duplicates
          const edgeKey = `${sourceId}_${targetId}_${edge.sourceHandle || ''}`
          if (seenEdges.has(edgeKey)) return null
          seenEdges.add(edgeKey)

          // Include condition data for gateway paths
          const condition = edge.data?.condition || null
          return {
            parentID: sourceId,
            childID: targetId,
            handle: `path_${sourceId}_${targetId}`,
            ...(condition ? { condition } : {}),
            meta: {},
          }
        })
        .filter(Boolean)

      // Generate back-edges for iterator nodes
      // For each iterator, find its body path (first outgoing edge) and
      // trace down to the last node whose only outgoing edge would be a back-edge.
      // The last step in the body chain gets a path back to the iterator.
      nodes.value.filter(n => n.type === 'iterator').forEach(iterNode => {
        const iterEdges = edges.value.filter(e => e.source === iterNode.id)
        const bodyEdge = iterEdges[0] // First edge is body
        if (!bodyEdge) return

        // Walk body chain to find the leaf (node with no outgoing edges, or an end node)
        let current = bodyEdge.target
        const visited = new Set<string>()
        while (current && !visited.has(current)) {
          visited.add(current)
          const currentNode = nodes.value.find(n => n.id === current)
          // Stop at end nodes - the step BEFORE the end is the body's last real step
          if (currentNode?.type === 'end') break
          const nextEdge = edges.value.find(e => e.source === current)
          if (!nextEdge) break
          current = nextEdge.target
        }

        // Find the last non-end step in the body chain
        let lastBodyStep: string | null = null
        for (const nodeId of visited) {
          const node = nodes.value.find(n => n.id === nodeId)
          if (node && node.type !== 'end') {
            lastBodyStep = nodeId
          }
        }

        if (lastBodyStep) {
          const iterBackendId = nodeIdToBackendId.get(iterNode.id)
          const bodyBackendId = nodeIdToBackendId.get(lastBodyStep)
          if (iterBackendId && bodyBackendId) {
            paths.push({
              parentID: bodyBackendId,
              childID: iterBackendId,
              handle: `path_${bodyBackendId}_${iterBackendId}`,
              meta: {},
            })
          }
        }
      })

      const dataToSave = {
        automationID: automation.value.automationID,
        handle: automation.value.handle,
        meta: automation.value.meta,
        enabled: automation.value.enabled,
        triggers,
        steps,
        paths,
        runAs: automation.value.runAs,
        ownedBy: automation.value.ownedBy,
      }

      let response
      if (automation.value.automationID && automation.value.automationID !== '0') {
        response = await $AutomationAPI.ngAutomationUpdate(dataToSave)
      } else {
        response = await $AutomationAPI.ngAutomationCreate(dataToSave)
      }

      automation.value = new NgAutomation(response)

      // Reload the flow to get proper IDs from backend
      const store = useAutomationStore()
      const state = automationToVueFlow(automation.value, {
        functions: store.functions,
        triggers: store.triggers,
      })
      nodes.value = state.nodes
      edges.value = state.edges

      // Update URL if this was a new automation
      if (router.currentRoute.value.params.id !== automation.value.automationID) {
        router.replace(`/builder/${automation.value.automationID}`)
      }

      $toast?.toastSuccess(t('builder.toast.saved.detail'), t('builder.toast.saved.summary'))
      return automation.value
    } catch (e) {
      console.error('Failed to save automation:', e)
      $toast?.toastDanger(t('builder.toast.error.detail'), t('builder.toast.error.summary'))
      throw e
    } finally {
      saving.value = false
    }
  }

  // Reset to empty state
  function reset() {
    automation.value = new NgAutomation()
    nodes.value = []
    edges.value = []
    history.value = []
    historyIndex.value = -1
    nextTick(() => saveToHistory())
  }

  // Add a new node to the flow
  function addNode(
    nodeType: NodeType,
    insertionPoint: InsertionPoint | null,
  ): Node<FlowNodeData> | null {
    const isBranch =
      nodeType.id === 'branch' || nodeType.type === 'condition' || nodeType.ref === 'gateway' || nodeType.ref?.startsWith('gateway')
    const isIterator = nodeType.kind === 'iterator'
    const gatewayRef = nodeType.ref?.startsWith('gateway') ? nodeType.ref : 'gatewayExclusive'
    const isEnd = nodeType.type === 'end'
    const isTrigger = nodeType.type === 'trigger'

    // Find max ID from BOTH triggers and steps (single shared counter)
    const maxTriggerID =
      automation.value.triggers?.reduce((max: number, t: any) => {
        const num = parseInt(t.triggerID, 10)
        return isNaN(num) ? max : Math.max(max, num)
      }, 0) || 0
    const maxStepID =
      automation.value.steps?.reduce((max: number, s: any) => {
        const num = parseInt(s.stepID, 10)
        return isNaN(num) ? max : Math.max(max, num)
      }, 0) || 0
    const newId = String(Math.max(maxTriggerID, maxStepID) + 1)

    let newHandle: string

    if (isTrigger) {
      newHandle = `trigger_${newId}`
      automation.value.triggers = automation.value.triggers || []
      automation.value.triggers.push({
        triggerID: newId,
        handle: newHandle,
        enabled: true,
        resourceType: nodeType.resourceType || '',
        eventType: nodeType.eventType || nodeType.ref || '',
        constraints: [],
        meta: {
          short: nodeType.label,
          description: nodeType.description || '',
        },
        input: {},
      })
    } else if (!isEnd) {
      newHandle = `step_${newId}`
      automation.value.steps = automation.value.steps || []
      automation.value.steps.push({
        stepID: newId,
        handle: newHandle,
        kind: isIterator ? 'iterator' : isBranch ? gatewayRef : 'function',
        ref: isIterator ? nodeType.ref || '' : isBranch ? gatewayRef : nodeType.ref || '',
        meta: {
          short: nodeType.label,
          description: nodeType.description || '',
        },
        arguments: [],
      })
    } else {
      newHandle = `end_${newId}`
    }

    // VueFlow node ID must be unique - prefix with type
    const vueFlowNodeId = isTrigger ? `trigger_${newId}` : isEnd ? `end_${newId}` : `step_${newId}`

    // Determine VueFlow node type
    let vueFlowType: string = 'step'
    if (isBranch) vueFlowType = 'branch'
    else if (isIterator) vueFlowType = 'iterator'
    else if (isTrigger) vueFlowType = 'trigger'
    else if (isEnd) vueFlowType = 'end'

    // Create VueFlow node
    const newNode: Node<FlowNodeData> = {
      id: vueFlowNodeId,
      type: vueFlowType,
      position: { x: 0, y: 0 },
      selectable: !isEnd,
      data: {
        label: nodeType.label,
        description: nodeType.description,
        icon: nodeType.icon,
        nodeType: isBranch
          ? gatewayRef
          : isTrigger
            ? nodeType.eventType || nodeType.ref || ''
            : nodeType.ref || (isEnd ? 'end' : ''),
        config: {},
        arguments: [],
        constraints: isTrigger ? [] : undefined,
        resourceType: isTrigger ? nodeType.resourceType || '' : undefined,
        ref: newHandle,
        stepID: isTrigger ? undefined : newId,
        triggerID: isTrigger ? newId : undefined,
      },
    }

    if (insertionPoint?.first) {
      // First node (trigger)
      nodes.value = [newNode]

      // Add end node (will be saved as termination step)
      const endId = String(parseInt(newId) + 1)
      const endVueId = `end_${endId}`

      nodes.value.push({
        id: endVueId,
        type: 'end',
        position: { x: 0, y: 0 },
        selectable: false,
        data: {
          label: t('builder.nodes.end'),
          nodeType: 'termination',
          icon: DEFAULT_ICONS.END,
          config: {},
          arguments: [],
          ref: endVueId,
        },
      })
      edges.value = [
        {
          id: `${vueFlowNodeId}_${endVueId}`,
          source: vueFlowNodeId,
          target: endVueId,
          type: 'addable',
        },
      ]
    } else if (insertionPoint?.edgeId) {
      // Insert on an edge
      const edgeIndex = edges.value.findIndex(e => e.id === insertionPoint.edgeId)
      const edge = edges.value[edgeIndex]
      if (edge && edgeIndex !== -1) {
        // Collect new nodes and edges first, then apply in batch
        const newNodes: typeof nodes.value = [newNode]
        const newEdges: typeof edges.value = []
        let insertEdge: (typeof edges.value)[0] | null = null

        if (isBranch || isIterator) {
          // Edge from source to branch/iterator (will be inserted at original position)
          insertEdge = {
            id: `${edge.source}_${vueFlowNodeId}`,
            source: edge.source,
            target: vueFlowNodeId,
            type: 'addable',
          }

          // First output: body (for iterator) or first branch path
          // For iterators: body starts with a new end node (user adds steps into it)
          // For branches: connect to original target (leftmost)
          const bodyEndId = String(parseInt(newId) + 1)
          const bodyEndVueId = `end_${bodyEndId}`

          newNodes.push({
            id: bodyEndVueId,
            type: 'end',
            position: { x: 0, y: 0 },
            selectable: false,
            data: {
              label: t('builder.nodes.end'),
              nodeType: 'termination',
              icon: DEFAULT_ICONS.END,
              config: {},
              arguments: [],
              ref: bodyEndVueId,
            },
          })

          if (isIterator) {
            // Iterator: first output (body) to new end node, second output (done) to original target
            newEdges.push({
              id: `${vueFlowNodeId}_0_${bodyEndVueId}`,
              source: vueFlowNodeId,
              target: bodyEndVueId,
              type: 'addable',
            })
            newEdges.push({
              id: `${vueFlowNodeId}_1_${edge.target}`,
              source: vueFlowNodeId,
              target: edge.target,
              type: 'addable',
            })
          } else {
            // Branch: first output to original target, second to new end node
            newEdges.push({
              id: `${vueFlowNodeId}_0_${edge.target}`,
              source: vueFlowNodeId,
              target: edge.target,
              type: 'addable',
            })
            newEdges.push({
              id: `${vueFlowNodeId}_1_${bodyEndVueId}`,
              source: vueFlowNodeId,
              target: bodyEndVueId,
              type: 'addable',
            })
          }
        } else {
          // Regular node: source -> new -> target
          insertEdge = {
            id: `${edge.source}_${vueFlowNodeId}`,
            source: edge.source,
            target: vueFlowNodeId,
            type: 'addable',
          }
          newEdges.push({
            id: `${vueFlowNodeId}_${edge.target}`,
            source: vueFlowNodeId,
            target: edge.target,
            type: 'addable',
          })
        }

        // Apply changes: NODES FIRST, then edges
        nodes.value.push(...newNodes)

        // Remove old edge and insert replacement at same position, then add rest
        const filteredEdges = edges.value.filter(e => e.id !== insertionPoint.edgeId)
        if (insertEdge) {
          filteredEdges.splice(edgeIndex, 0, insertEdge)
        }
        edges.value = [...filteredEdges, ...newEdges]
      }
    }

    // Re-layout nodes
    const layouted = applyDagreLayout({ nodes: nodes.value, edges: edges.value })
    nodes.value = layouted.nodes
    edges.value = layouted.edges

    saveToHistory()

    // Return the layouted node for centering
    return layouted.nodes.find(n => n.id === vueFlowNodeId) || null
  }

  // Delete a node and reconnect edges
  // Termination nodes stay at leaves - when deleting a node, its termination children
  // are reconnected to the parent (not orphaned or duplicated)
  // Orphaned nodes (no incoming edges, except triggers) are cleaned up
  function deleteNode(nodeToDelete: Node<FlowNodeData>) {
    const nodeId = nodeToDelete.id

    // Remove from automation model
    if (nodeToDelete.data?.stepID) {
      automation.value.steps =
        automation.value.steps?.filter((s: any) => s.stepID !== nodeToDelete.data.stepID) || []
    }
    if (nodeToDelete.data?.triggerID) {
      automation.value.triggers =
        automation.value.triggers?.filter(
          (t: any) => t.triggerID !== nodeToDelete.data.triggerID,
        ) || []
    }

    // Handle reconnection & dynamic End nodes
    const incomingEdges = edges.value.filter(e => e.target === nodeId)
    const outgoingEdges = edges.value.filter(e => e.source === nodeId)

    // Remove node and its incident edges
    nodes.value = nodes.value.filter(n => n.id !== nodeId)
    edges.value = edges.value.filter(e => e.source !== nodeId && e.target !== nodeId)

    // If this was an End node, the parent becomes a leaf - add new termination
    if (nodeToDelete.type === 'end') {
      incomingEdges.forEach(incoming => {
        const parentNode = nodes.value.find(n => n.id === incoming.source)
        // Check if parent still has other outgoing edges
        const parentHasOtherChildren = edges.value.some(e => e.source === incoming.source)
        if (parentNode && !parentHasOtherChildren) {
          // Parent is now a leaf, add termination
          const newEndId = `end_${incoming.source}_${incoming.sourceHandle || 'default'}`
          nodes.value.push({
            id: newEndId,
            type: 'end',
            position: { x: 0, y: 0 },
            selectable: false,
            data: {
              label: t('builder.nodes.end'),
              nodeType: 'termination',
              icon: DEFAULT_ICONS.END,
              config: {},
              arguments: [],
              ref: newEndId,
            },
          })
          edges.value.push({
            id: `${incoming.source}_${newEndId}`,
            source: incoming.source,
            sourceHandle: incoming.sourceHandle,
            target: newEndId,
            type: 'addable',
          })
        }
      })
      cleanupOrphanedNodes()
      relayout()
      return
    }

    // For non-end nodes: reconnect parent to ALL children (including terminations)
    // EXCEPT for branches, where we want to delete the subtrees instead
    incomingEdges.forEach(incoming => {
      if (outgoingEdges.length > 0 && nodeToDelete.type !== 'branch' && nodeToDelete.type !== 'iterator') {
        // Reconnect parent to all children of deleted node
        outgoingEdges.forEach(outgoing => {
          if (incoming.source !== outgoing.target) {
            edges.value.push({
              id: `${incoming.source}_${outgoing.target}`,
              source: incoming.source,
              sourceHandle: incoming.sourceHandle,
              target: outgoing.target,
              type: 'addable',
            })
          }
        })
      } else {
        // Deleted node had no children OR was a branch: parent needs termination
        const parentNode = nodes.value.find(n => n.id === incoming.source)
        if (parentNode) {
          const newEndId = `end_${incoming.source}_${incoming.sourceHandle || 'default'}`
          nodes.value.push({
            id: newEndId,
            type: 'end',
            position: { x: 0, y: 0 },
            selectable: false,
            data: {
              label: t('builder.nodes.end'),
              nodeType: 'termination',
              icon: DEFAULT_ICONS.END,
              config: {},
              arguments: [],
              ref: newEndId,
            },
          })
          edges.value.push({
            id: `${incoming.source}_${newEndId}`,
            source: incoming.source,
            sourceHandle: incoming.sourceHandle,
            target: newEndId,
            type: 'addable',
          })
        }
      }
    })

    // Clean up any orphaned nodes (nodes with no incoming edges, except triggers)
    cleanupOrphanedNodes()
    relayout()
  }

  // Remove orphaned nodes - nodes with no incoming edges (except triggers which are roots)
  function cleanupOrphanedNodes() {
    let changed = true
    while (changed) {
      changed = false
      const orphanedNodes = nodes.value.filter(node => {
        // Triggers are roots, they don't need incoming edges
        if (node.type === 'trigger') return false
        // Check if this node has any incoming edges
        const hasIncoming = edges.value.some(e => e.target === node.id)
        return !hasIncoming
      })

      if (orphanedNodes.length > 0) {
        changed = true
        orphanedNodes.forEach(orphan => {
          // Remove orphan's outgoing edges
          edges.value = edges.value.filter(e => e.source !== orphan.id)
          // Remove orphan node
          nodes.value = nodes.value.filter(n => n.id !== orphan.id)
          // Remove from automation model if it has an ID
          if (orphan.data?.stepID) {
            automation.value.steps =
              automation.value.steps?.filter((s: any) => s.stepID !== orphan.data.stepID) || []
          }
        })
      }
    }
  }

  // Update node data and save to history
  function updateNodeData(nodeId: string, dataUpdate: Partial<FlowNodeData>) {
    const nodeIndex = nodes.value.findIndex(n => n.id === nodeId)
    if (nodeIndex === -1) return

    const node = nodes.value[nodeIndex]
    const newData = { ...node.data, ...dataUpdate }

    // Replace node in array to trigger Vue reactivity
    nodes.value[nodeIndex] = { ...node, data: newData }

    // Also update the automation model if this is a step with arguments
    if (dataUpdate.arguments && newData.stepID) {
      const step = automation.value.steps?.find((s: any) => s.stepID === newData.stepID)
      if (step) {
        step.arguments = dataUpdate.arguments
      }
    }

    // Update trigger constraints in automation model
    if (dataUpdate.constraints && newData.triggerID) {
      const trigger = automation.value.triggers?.find((t: any) => t.triggerID === newData.triggerID)
      if (trigger) {
        trigger.constraints = dataUpdate.constraints
      }
    }

    saveToHistory()
  }

  // Re-layout nodes after changes
  function relayout() {
    if (nodes.value.length > 0) {
      const layouted = applyDagreLayout({ nodes: nodes.value, edges: edges.value })
      nodes.value = layouted.nodes
      edges.value = layouted.edges
    }
    saveToHistory()
  }

  // Add a new output branch to an existing branch node (for Else If)
  function addBranchOutput(branchNode: Node<FlowNodeData>) {
    if (branchNode.type !== 'branch') return

    // Generate new end node ID
    const newEndId = `end_${branchNode.id}_${Date.now()}`

    // Add new end node
    nodes.value.push({
      id: newEndId,
      type: 'end',
      position: { x: 0, y: 0 },
      selectable: false,
      data: {
        label: t('builder.nodes.end'),
        nodeType: 'termination',
        icon: DEFAULT_ICONS.END,
        config: {},
        arguments: [],
        ref: newEndId,
      },
    })

    // Find existing edges from this branch to determine insertion position
    // Insert new edge BEFORE the last edge (the "Else" branch)
    const branchEdges = edges.value.filter(e => e.source === branchNode.id)
    const lastBranchEdgeIndex = edges.value.findIndex(
      e => e.id === branchEdges[branchEdges.length - 1]?.id,
    )

    const newEdge = {
      id: `${branchNode.id}_${newEndId}`,
      source: branchNode.id,
      target: newEndId,
      type: 'addable',
    }

    // Insert before the last branch edge (Else stays last)
    if (lastBranchEdgeIndex !== -1) {
      edges.value.splice(lastBranchEdgeIndex, 0, newEdge)
    } else {
      edges.value.push(newEdge)
    }

    // Re-layout
    const layouted = applyDagreLayout({ nodes: nodes.value, edges: edges.value })
    nodes.value = layouted.nodes
    edges.value = layouted.edges

    saveToHistory()
  }

  // Reorder edges for a branch node based on new order
  function reorderBranchEdges(branchNodeId: string, newEdgeOrder: string[]) {
    // Get all edges from this branch
    const branchEdgeIds = new Set(edges.value.filter(e => e.source === branchNodeId).map(e => e.id))

    // Get branch edges
    const branchEdges = edges.value.filter(e => branchEdgeIds.has(e.id))

    // Reorder branch edges according to newEdgeOrder
    const reorderedBranchEdges = newEdgeOrder
      .map(id => branchEdges.find(e => e.id === id))
      .filter(Boolean) as typeof edges.value

    // Find where the first branch edge was in the original array
    const firstBranchIndex = edges.value.findIndex(e => branchEdgeIds.has(e.id))

    // Reconstruct edges array: non-branch edges with reordered branch edges inserted at original position
    const beforeBranch = edges.value
      .slice(0, firstBranchIndex)
      .filter(e => !branchEdgeIds.has(e.id))
    const afterBranch = edges.value.slice(firstBranchIndex).filter(e => !branchEdgeIds.has(e.id))

    edges.value = [...beforeBranch, ...reorderedBranchEdges, ...afterBranch]

    // Re-layout
    const layouted = applyDagreLayout({ nodes: nodes.value, edges: edges.value })
    nodes.value = layouted.nodes
    edges.value = layouted.edges

    saveToHistory()
  }

  async function exec() {
    const id = automation.value.automationID
    if (!id || id === '0') return

    running.value = true
    traceStatus.value = 'running'
    traceFrames.value = []
    traceExecution.value = null

    try {
      // ExecAndWait blocks until the automation completes
      const result = await $AutomationAPI.ngAutomationExec({
        automationID: id,
        trace: true,
      })

      // Extract execution result from response
      const execResult: ExecutionResult = result?.response ?? result
      traceExecution.value = execResult

      // Fetch execution trace frames
      if (execResult?.executionID) {
        try {
          const traceResponse = await $AutomationAPI.ngAutomationExecutionTrace({
            automationID: id,
            executionID: execResult.executionID,
          })
          traceFrames.value = traceResponse?.response ?? traceResponse ?? []
        } catch (traceErr) {
          console.error('Failed to fetch trace:', traceErr)
        }
      }

      if (execResult?.status === 'failed') {
        traceStatus.value = 'failed'
        $toast?.toastDanger(
          execResult.error || t('builder.toast.runError.detail'),
          t('builder.toast.runError.summary'),
        )
      } else {
        traceStatus.value = 'completed'
        $toast?.toastSuccess(t('builder.toast.run.detail'), t('builder.toast.run.summary'))
      }
    } catch (e) {
      console.error('Failed to execute automation:', e)
      traceStatus.value = 'failed'
      $toast?.toastDanger(t('builder.toast.runError.detail'), t('builder.toast.runError.summary'))
    } finally {
      running.value = false
    }
  }

  /**
   * Clear active trace overlay.
   */
  function clearTrace() {
    traceFrames.value = []
    traceExecution.value = null
    traceStatus.value = 'idle'
  }

  /**
   * Map step handle → StackFrame for quick lookup by flow nodes.
   * The backend StackFrame.handle matches the step handle (e.g., "step_3").
   */
  const traceByHandle = computed(() => {
    const map = new Map<string, StackFrame>()
    for (const frame of traceFrames.value) {
      if (frame.handle) {
        map.set(frame.handle, frame)
      }
    }
    return map
  })

  /**
   * Walk edges backward from a node to find all upstream step nodes,
   * then look up each step's function definition to get its results.
   */
  function getUpstreamResults(nodeId: string) {
    const store = useAutomationStore()
    const upstream: Array<{
      handle: string
      label: string
      icon?: IconDef
      results: Array<{
        name: string
        sourceName: string
        types: string[]
        expandable?: boolean
        namespaceID?: string
        moduleID?: string
      }>
    }> = []

    // Walk backward through edges to find all ancestor nodes
    const visited = new Set<string>()
    const queue = [nodeId]

    while (queue.length > 0) {
      const currentId = queue.shift()!
      if (visited.has(currentId)) continue
      visited.add(currentId)

      const incomingEdges = edges.value.filter(e => e.target === currentId)
      for (const edge of incomingEdges) {
        queue.push(edge.source)
      }
    }

    // Remove the node itself from visited
    visited.delete(nodeId)

    // Helper: resolve a trigger constraint value by property name
    function getTriggerConstraintValue(nodeData: any, propName: string): string | null {
      const constraints = nodeData?.constraints || []
      const c = constraints.find((cc: any) => cc.name === propName)
      return c?.values?.[0]?.['@value'] ?? null
    }

    // For each ancestor node, look up results (functions) or properties (triggers)
    for (const ancestorId of visited) {
      const node = nodes.value.find(n => n.id === ancestorId)
      if (!node || node.type === 'end') continue

      if (node.type === 'trigger') {
        // Look up trigger definition for properties
        const eventType = node.data?.nodeType
        const resourceType = node.data?.resourceType
        const triggerDef = store.triggers.find(
          t => t.eventType === eventType && (!resourceType || t.resourceType === resourceType),
        )
        if (!triggerDef?.properties?.length) continue

        upstream.push({
          handle: node.data?.ref || ancestorId,
          label: node.data?.label || triggerDef.meta?.short || triggerDef.eventType,
          icon: (triggerDef.meta?.icon || node.data?.icon) as IconDef | undefined,
          results: triggerDef.properties.map(p => {
            const result: any = {
              name: p.meta?.short || p.name,
              sourceName: p.name,
              types: p.type ? [p.type] : [],
            }

            // Mark record-type properties as expandable and attach constraint IDs
            if (p.type === 'ComposeRecord') {
              result.expandable = true
              result.namespaceID = getTriggerConstraintValue(node.data, 'namespace')
              result.moduleID = getTriggerConstraintValue(node.data, 'module')
            }

            return result
          }),
        })
      } else {
        // Look up function definition for results
        const funcDef = store.functions.find(f => f.ref === node.data?.nodeType)
        if (!funcDef?.results?.length) continue

        upstream.push({
          handle: node.data?.ref || ancestorId,
          label: node.data?.label || funcDef.meta?.short || funcDef.ref,
          icon: (funcDef.meta?.icon || node.data?.icon) as IconDef | undefined,
          results: funcDef.results.map(r => ({
            name: r.argumentName,
            sourceName: r.argumentName,
            types: r.types || [],
          })),
        })
      }
    }

    return upstream
  }

  // Update the condition on an edge (for branch conditions)
  function updateEdgeCondition(edgeId: string, condition: Record<string, unknown> | null) {
    const edgeIndex = edges.value.findIndex(e => e.id === edgeId)
    if (edgeIndex === -1) return

    edges.value[edgeIndex] = {
      ...edges.value[edgeIndex],
      data: { ...edges.value[edgeIndex].data, condition },
    }

    saveToHistory()
  }

  // Update the gateway type (gatewayExclusive/gatewayInclusive) for a branch node
  /**
   * Replace a node with a different node type, preserving its position in the flow (edges).
   * Resets arguments and updates the automation model accordingly.
   */
  function replaceNode(nodeId: string, newNodeType: NodeType) {
    const nodeIndex = nodes.value.findIndex(n => n.id === nodeId)
    if (nodeIndex === -1) return

    const oldNode = nodes.value[nodeIndex]
    const isBranch =
      newNodeType.id === 'branch' || newNodeType.type === 'condition' || newNodeType.ref === 'gateway'
    const isIterator = newNodeType.kind === 'iterator'
    const isTrigger = newNodeType.type === 'trigger'
    const hasTwoOutputs = isBranch || isIterator

    let newVueFlowType = 'step'
    if (isBranch) newVueFlowType = 'branch'
    else if (isIterator) newVueFlowType = 'iterator'
    else if (isTrigger) newVueFlowType = 'trigger'

    const gatewayRef = newNodeType.ref?.startsWith('gateway') ? newNodeType.ref : 'gatewayExclusive'

    // Determine new nodeType value for data
    const newDataNodeType = isBranch
      ? gatewayRef
      : isTrigger
        ? newNodeType.eventType || newNodeType.ref || ''
        : newNodeType.ref || ''

    // Update automation model
    if (oldNode.type === 'trigger' && oldNode.data?.triggerID) {
      // Update existing trigger in automation model
      const trigger = automation.value.triggers?.find(
        (t: any) => t.triggerID === oldNode.data.triggerID,
      )
      if (trigger) {
        trigger.eventType = newNodeType.eventType || newNodeType.ref || ''
        trigger.resourceType = newNodeType.resourceType || ''
        trigger.constraints = []
        trigger.input = {}
        trigger.meta = {
          short: newNodeType.label,
          description: newNodeType.description || '',
        }
      }
    } else if (oldNode.data?.stepID) {
      // Update existing step in automation model
      const step = automation.value.steps?.find(
        (s: any) => s.stepID === oldNode.data.stepID,
      )
      if (step) {
        step.kind = isIterator ? 'iterator' : isBranch ? gatewayRef : 'function'
        step.ref = isIterator ? newNodeType.ref || '' : isBranch ? gatewayRef : newNodeType.ref || ''
        step.arguments = []
        step.meta = {
          short: newNodeType.label,
          description: newNodeType.description || '',
        }
      }
    }

    // Build updated node data
    const newData: FlowNodeData = {
      ...oldNode.data,
      label: newNodeType.label,
      description: newNodeType.description,
      icon: newNodeType.icon,
      nodeType: newDataNodeType,
      arguments: [],
      config: {},
      constraints: isTrigger ? [] : undefined,
      resourceType: isTrigger ? newNodeType.resourceType || '' : undefined,
    }

    // Replace the node in the array (type may change, e.g. step→branch)
    nodes.value[nodeIndex] = {
      ...oldNode,
      type: newVueFlowType,
      data: newData,
    }

    const oldHasTwoOutputs = oldNode.type === 'branch' || oldNode.type === 'iterator'

    // If switching to a two-output type and the node currently has a single outgoing edge,
    // add a second output (end node)
    if (hasTwoOutputs && !oldHasTwoOutputs) {
      const outgoingEdges = edges.value.filter(e => e.source === nodeId)
      if (outgoingEdges.length === 1) {
        const newEndId = `end_${nodeId}_${Date.now()}`
        nodes.value.push({
          id: newEndId,
          type: 'end',
          position: { x: 0, y: 0 },
          selectable: false,
          data: {
            label: t('builder.nodes.end'),
            nodeType: 'termination',
            icon: DEFAULT_ICONS.END,
            config: {},
            arguments: [],
            ref: newEndId,
          },
        })
        edges.value.push({
          id: `${nodeId}_${newEndId}`,
          source: nodeId,
          target: newEndId,
          type: 'addable',
        })
      }
    }

    // If switching FROM a two-output type to a single-output type, collapse extra outputs:
    // keep only the first outgoing edge and orphan-cleanup the rest
    if (!hasTwoOutputs && oldHasTwoOutputs) {
      const outgoingEdges = edges.value.filter(e => e.source === nodeId)
      if (outgoingEdges.length > 1) {
        // Keep first edge, remove the rest
        const edgesToRemove = outgoingEdges.slice(1).map(e => e.id)
        edges.value = edges.value.filter(e => !edgesToRemove.includes(e.id))
        // Also remove any condition data from the kept edge
        const keptEdge = edges.value.find(e => e.source === nodeId)
        if (keptEdge?.data?.condition) {
          keptEdge.data = { ...keptEdge.data, condition: null }
        }
        cleanupOrphanedNodes()
      }
    }

    relayout()
  }

  function updateGatewayType(nodeId: string, newGatewayRef: string) {
    const node = nodes.value.find(n => n.id === nodeId)
    if (!node) return

    // Resolve label and description for the new gateway type
    const isExclusive = newGatewayRef === 'gatewayExclusive'
    const label = isExclusive
      ? t('builder.nodePicker.nodes.branches.exclusive.label')
      : t('builder.nodePicker.nodes.branches.inclusive.label')
    const description = isExclusive
      ? t('builder.nodePicker.nodes.branches.exclusive.description')
      : t('builder.nodePicker.nodes.branches.inclusive.description')

    // Directly mutate node data properties for VueFlow reactivity
    node.data = { ...node.data, nodeType: newGatewayRef, label, description }

    // Also update the automation model
    if (node.data?.stepID) {
      const step = automation.value.steps?.find((s: any) => s.stepID === node.data.stepID)
      if (step) {
        step.ref = newGatewayRef
        step.kind = newGatewayRef
        step.meta = { ...step.meta, short: label, description }
      }
    }

    saveToHistory()
  }

  return {
    // State
    automation,
    nodes,
    edges,
    loading,
    saving,
    running,

    // Trace state
    traceFrames,
    traceExecution,
    traceStatus,
    traceByHandle,

    // Computed
    name,
    enabled,
    automationId,
    isEmpty,
    canUndo,
    canRedo,

    // Actions
    load,
    save,
    exec,
    clearTrace,
    reset,
    undo,
    redo,
    addNode,
    addBranchOutput,
    reorderBranchEdges,
    deleteNode,
    replaceNode,
    updateNodeData,
    updateEdgeCondition,
    updateGatewayType,
    saveToHistory,
    getUpstreamResults,
  }
}
