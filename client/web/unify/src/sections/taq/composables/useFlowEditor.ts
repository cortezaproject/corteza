import { automation } from '@planetcrust/human-js'
import type {
  StackFrame,
  ExecutionResult,
  TraceStatus,
} from '@planetcrust/human-js/src/automation/types/trace'
import { withMinDuration } from '@planetcrust/human-vue'
import type { IconDef } from '@planetcrust/human-js/src/automation/types/icon'
import { DEFAULT_ICONS } from '@planetcrust/human-js/src/automation/types/icon'
import type { Edge, Node } from '@vue-flow/core'
import { computed, inject, nextTick, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useAutomationStore } from '@planetcrust/human-vue'
import {
  applyDagreLayout,
  automationToVueFlow,
  type FlowNodeData,
} from '@/sections/taq/utils/taq-parser'

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
  const lastSavedHistoryIndex = ref(0)

  // Computed
  const isDirty = computed(() => historyIndex.value !== lastSavedHistoryIndex.value)
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
    // Cap history to prevent unbounded memory growth
    const MAX_HISTORY = 100
    if (history.value.length > MAX_HISTORY) {
      const excess = history.value.length - MAX_HISTORY
      history.value = history.value.slice(excess)
      historyIndex.value -= excess
      lastSavedHistoryIndex.value = Math.max(0, lastSavedHistoryIndex.value - excess)
    }
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

      nextTick(() => {
        saveToHistory()
        lastSavedHistoryIndex.value = historyIndex.value
      })
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

      // Map to track VueFlow node ID -> backend ID for path generation
      const nodeIdToBackendId = new Map<string, string>()

      // Compute the highest existing ID across all nodes to use as a baseline
      // for minting new IDs, so we never collide with existing scope references.
      let nextNewID = nodes.value.reduce((max, node) => {
        const data = node.data as FlowNodeData
        const id = parseInt((node.type === 'trigger' ? data.triggerID : data.stepID) || '0', 10)
        return isNaN(id) ? max : Math.max(max, id)
      }, 0)

      // Helper to process a node into trigger or step.
      // Preserves the existing backend ID for already-saved nodes; only assigns
      // a fresh ID (above the current max) to nodes that don't have one yet.
      // This ensures stored scope references (e.g. "step_2") stay valid across saves.
      const processNode = (node: Node<FlowNodeData>) => {
        const data = node.data as FlowNodeData

        let backendID: string
        if (node.type === 'trigger') {
          if (data.triggerID && data.triggerID !== '0') {
            backendID = data.triggerID
          } else {
            nextNewID++
            backendID = String(nextNewID)
          }
        } else {
          if (data.stepID && data.stepID !== '0') {
            backendID = data.stepID
          } else {
            nextNewID++
            backendID = String(nextNewID)
          }
        }

        nodeIdToBackendId.set(node.id, backendID)

        if (node.type === 'trigger') {
          const existing = automation.value.triggers?.find(
            (t: any) => t.triggerID === data.triggerID,
          )
          triggers.push({
            triggerID: backendID,
            handle: `trigger_${backendID}`,
            enabled: existing?.enabled ?? true,
            resourceType: data.resourceType || existing?.resourceType || '',
            eventType: data.nodeType || existing?.eventType || '',
            constraints: data.constraints || existing?.constraints || [],
            inputSchema: data.inputSchema || existing?.inputSchema || [],
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
            stepID: backendID,
            handle: `step_${backendID}`,
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

      const incomingTargets = new Set(edges.value.map(e => e.target))

      // Process in order: 1) triggers, 2) terminations, 3) other steps
      nodes.value.filter(n => n.type === 'trigger').forEach(processNode)
      nodes.value.filter(n => n.type === 'end' && incomingTargets.has(n.id)).forEach(processNode)
      nodes.value
        .filter(n => n.type !== 'trigger' && n.type !== 'end' && n.type !== 'loop')
        .forEach(processNode)

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

      const dataToSave = {
        automationID: automation.value.automationID,
        handle: automation.value.handle,
        labels: automation.value.labels,
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
      store.updateInList(automation.value)
      const state = automationToVueFlow(automation.value, {
        functions: store.functions,
        triggers: store.triggers,
      })
      nodes.value = state.nodes
      edges.value = state.edges

      // Update URL if this was a new automation
      if (router.currentRoute.value.params.id !== automation.value.automationID) {
        router.replace(`/taq/builder/${automation.value.automationID}`)
      }

      $toast?.toastSuccess(t('builder.toast.saved.detail'), t('builder.toast.saved.summary'))
      lastSavedHistoryIndex.value = historyIndex.value

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
    nextTick(() => {
      saveToHistory()
      lastSavedHistoryIndex.value = historyIndex.value
    })
  }

  // Add a new node to the flow
  function addNode(
    nodeType: NodeType,
    insertionPoint: InsertionPoint | null,
  ): Node<FlowNodeData> | null {
    const isBranch =
      nodeType.id === 'branch' ||
      nodeType.type === 'condition' ||
      nodeType.ref === 'gateway' ||
      nodeType.ref?.startsWith('gateway')
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

    const store = useAutomationStore()

    let defaultArguments: any[] = []
    let defaultConfig: Record<string, any> = {}

    if (isTrigger) {
      const dbTrigger = store.triggers.find(
        t => t.eventType === nodeType.eventType && t.resourceType === nodeType.resourceType,
      )
      if (dbTrigger?.segments) {
        dbTrigger.segments.forEach(seg => {
          seg.sections?.forEach(sec => {
            sec.elements?.forEach(el => {
              const elInput = el.input as any
              if (elInput && elInput.default !== undefined && elInput.argument) {
                defaultConfig[elInput.argument] = elInput.default
              }
            })
          })
        })
      }
    } else if (!isEnd && !isBranch) {
      const dbFunction = store.functions.find(f => f.ref === nodeType.ref)
      if (dbFunction?.segments) {
        dbFunction.segments.forEach(seg => {
          seg.sections?.forEach(sec => {
            sec.elements?.forEach(el => {
              const elInput = el.input as any
              if (elInput && elInput.default !== undefined && elInput.argument) {
                const param = dbFunction.parameters?.find(p => p.argumentName === elInput.argument)
                defaultArguments.push({
                  argumentName: elInput.argument,
                  type: param?.types?.[0] || 'Any',
                  value: elInput.default,
                })
              }
            })
          })
        })
      }
    }

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
        input: { ...defaultConfig },
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
        arguments: [...defaultArguments],
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
        config: { ...defaultConfig },
        arguments: [...defaultArguments],
        constraints: isTrigger ? [] : undefined,
        resourceType: isTrigger ? nodeType.resourceType || '' : undefined,
        ref: newHandle,
        stepID: isTrigger ? undefined : newId,
        triggerID: isTrigger ? newId : undefined,
      },
    }

    let nextNodes: typeof nodes.value = [...nodes.value]
    let nextEdges: typeof edges.value = [...edges.value]

    if (insertionPoint?.first) {
      // First node (trigger)
      nextNodes = [newNode]

      // Add end node (will be saved as termination step)
      const endId = String(parseInt(newId) + 1)
      const endVueId = `end_${endId}`

      nextNodes.push({
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
      nextEdges = [
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
          // For iterators: body starts with a new loop node (user adds steps into it)
          // For branches: connect to original target (leftmost)
          const bodyEndId = String(parseInt(newId) + 1)
          const bodyEndVueId = isIterator ? `loop_${bodyEndId}` : `end_${bodyEndId}`

          newNodes.push({
            id: bodyEndVueId,
            type: isIterator ? 'loop' : 'end',
            position: { x: 0, y: 0 },
            selectable: false,
            data: {
              label: isIterator ? t('builder.nodes.loop') : t('builder.nodes.end'),
              nodeType: isIterator ? 'loop' : 'termination',
              icon: isIterator ? undefined : DEFAULT_ICONS.END,
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

        // Prepare new state immutably: NODES FIRST, then edges
        nextNodes = [...nodes.value, ...newNodes]

        // Remove old edge and insert replacement at same position, then add rest
        const filteredEdges = edges.value.filter(e => e.id !== insertionPoint.edgeId)
        if (insertEdge) {
          filteredEdges.splice(edgeIndex, 0, insertEdge)
        }
        nextEdges = [...filteredEdges, ...newEdges]
      }
    }

    // Re-layout nodes
    const layouted = applyDagreLayout({ nodes: nextNodes, edges: nextEdges })
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

    // If this was an End or Loop node, the parent becomes a leaf - add new termination/loop
    if (nodeToDelete.type === 'end' || nodeToDelete.type === 'loop') {
      const type = nodeToDelete.type
      incomingEdges.forEach(incoming => {
        const parentNode = nodes.value.find(n => n.id === incoming.source)
        // Check if parent still has other outgoing edges
        const parentHasOtherChildren = edges.value.some(e => e.source === incoming.source)
        if (parentNode && !parentHasOtherChildren) {
          // Parent is now a leaf, add termination or loop
          const newEndId = `${type}_${incoming.source}_${incoming.sourceHandle || 'default'}`
          nodes.value.push({
            id: newEndId,
            type: type,
            position: { x: 0, y: 0 },
            selectable: false,
            data: {
              label: type === 'loop' ? t('builder.nodes.loop') : t('builder.nodes.end'),
              nodeType: type === 'loop' ? 'loop' : 'termination',
              icon: type === 'loop' ? undefined : DEFAULT_ICONS.END,
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
      if (
        outgoingEdges.length > 0 &&
        nodeToDelete.type !== 'branch' &&
        nodeToDelete.type !== 'iterator'
      ) {
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

    let nextNodes = [...nodes.value]
    let nextEdges = [...edges.value]

    // Add new end node
    nextNodes.push({
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
    const branchEdges = nextEdges.filter(e => e.source === branchNode.id)
    const lastBranchEdgeIndex = nextEdges.findIndex(
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
      nextEdges.splice(lastBranchEdgeIndex, 0, newEdge)
    } else {
      nextEdges.push(newEdge)
    }

    // Re-layout
    const layouted = applyDagreLayout({ nodes: nextNodes, edges: nextEdges })
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

  async function exec(input?: Record<string, unknown>) {
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
        ...(input ? { input } : {}),
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
      results?: Array<{
        name: string
        sourceName: string
        types: string[]
        expandable?: boolean
        namespaceID?: string
        moduleID?: string
        workflowID?: string
      }>
      properties?: Array<{
        name: string
        sourceName: string
        types: string[]
        expandable?: boolean
        namespaceID?: string
        moduleID?: string
        workflowID?: string
      }>
      description?: string
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

    const EXPANDABLE_TYPES = [
      'ComposeRecord',
      'SystemUser',
      'SystemRole',
      'SystemApplication',
      'ComposeNamespace',
      'ComposeModule',
      'ComposePage',
      'HttpRequest',
    ]

    // Helper: resolve a trigger constraint value by property name
    function getTriggerConstraintValue(nodeData: any, propName: string): string | null {
      const constraints = nodeData?.constraints || []
      const c = constraints.find((cc: any) => cc.name === propName)
      return c?.values?.[0]?.['@value'] ?? null
    }

    // Helper: resolve a function argument value by argument name
    function getFunctionArgumentValue(nodeData: any, argName: string): string | null {
      const args = nodeData?.arguments || []
      const a = args.find((aa: any) => aa.argumentName === argName)
      if (!a) return null

      // If entered manually in UI, expressions may be wrapped in quotes
      if (a.expr && typeof a.expr === 'string' && a.expr.startsWith('"') && a.expr.endsWith('"')) {
        return a.expr.slice(1, -1)
      }
      return a.expr || a.source || a.value || null
    }

    // Top-level scope vars injected by the server for every execution
    // (server/automation/service/ng_automation.go injectIdentities). Kept flat
    // so each click emits {scope:'invoker', source:'email'} which maps directly
    // to how the runtime resolves these names.
    const SYSTEM_USER_FIELDS = [
      { name: t('builder.referencePanel.systemUser.userID'), sourceName: 'userID', types: ['ID'] },
      {
        name: t('builder.referencePanel.systemUser.email'),
        sourceName: 'email',
        types: ['String'],
      },
      { name: t('builder.referencePanel.systemUser.name'), sourceName: 'name', types: ['String'] },
      {
        name: t('builder.referencePanel.systemUser.username'),
        sourceName: 'username',
        types: ['String'],
      },
      {
        name: t('builder.referencePanel.systemUser.handle'),
        sourceName: 'handle',
        types: ['Handle'],
      },
      {
        name: t('builder.referencePanel.systemUser.emailConfirmed'),
        sourceName: 'emailConfirmed',
        types: ['Boolean'],
      },
      {
        name: t('builder.referencePanel.systemUser.createdAt'),
        sourceName: 'createdAt',
        types: ['DateTime'],
      },
      {
        name: t('builder.referencePanel.systemUser.updatedAt'),
        sourceName: 'updatedAt',
        types: ['DateTime'],
      },
      {
        name: t('builder.referencePanel.systemUser.deletedAt'),
        sourceName: 'deletedAt',
        types: ['DateTime'],
      },
      {
        name: t('builder.referencePanel.systemUser.suspendedAt'),
        sourceName: 'suspendedAt',
        types: ['DateTime'],
      },
    ]

    upstream.push(
      {
        handle: 'invoker',
        label: t('builder.referencePanel.invoker'),
        description: t('builder.referencePanel.invokerDescription'),
        icon: { type: 'name', value: 'user' },
        properties: SYSTEM_USER_FIELDS,
      },
      {
        handle: 'runner',
        label: t('builder.referencePanel.runner'),
        description: t('builder.referencePanel.runnerDescription'),
        icon: { type: 'name', value: 'id-card' },
        properties: SYSTEM_USER_FIELDS,
      },
    )

    // For each ancestor node, look up results (functions) or properties (triggers)
    for (const ancestorId of visited) {
      const node = nodes.value.find(n => n.id === ancestorId)
      if (!node || node.type === 'end') continue

      if (node.type === 'trigger') {
        // Look up trigger definition for properties
        const eventType = node.data?.nodeType
        const resourceType = node.data?.resourceType

        // Agent trigger: properties come from the node's own inputSchema, not the catalog.
        // Each TAQ declares its own params, so the catalog definition is empty.
        if (eventType === 'onAgentic' || resourceType === 'automation:trigger:agentic') {
          const schema =
            (node.data?.inputSchema as Array<{ name: string; type: string }> | undefined) || []
          if (!schema.length) continue

          upstream.push({
            handle: node.data?.ref || ancestorId,
            label: node.data?.label || 'Agent Invoked',
            description: (node.data?.description as string | undefined) || undefined,
            icon: node.data?.icon as IconDef | undefined,
            properties: schema
              .filter(p => p.name)
              .map(p => ({
                name: p.name,
                sourceName: p.name,
                types: p.type ? [p.type] : [],
              })),
          })
          continue
        }

        const triggerDef = store.triggers.find(
          t => t.eventType === eventType && (!resourceType || t.resourceType === resourceType),
        )
        if (!triggerDef?.properties?.length) continue

        upstream.push({
          handle: node.data?.ref || ancestorId,
          label: node.data?.label || triggerDef.meta?.short || triggerDef.eventType,
          description:
            (node.data?.description as string | undefined) ||
            triggerDef.meta?.description ||
            undefined,
          icon: (triggerDef.meta?.icon || node.data?.icon) as IconDef | undefined,
          properties: triggerDef.properties.map(p => {
            const result: any = {
              name: p.meta?.short || p.name,
              sourceName: p.name,
              types: p.type ? [p.type] : [],
            }

            // Mark structurally complex types as expandable
            if (result.types.some((t: string) => EXPANDABLE_TYPES.includes(t))) {
              result.expandable = true

              if (result.types.includes('ComposeRecord')) {
                result.namespaceID = getTriggerConstraintValue(node.data, 'namespace')
                result.moduleID = getTriggerConstraintValue(node.data, 'module')
              }
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
          description:
            (node.data?.description as string | undefined) ||
            funcDef.meta?.description ||
            undefined,
          icon: (funcDef.meta?.icon || node.data?.icon) as IconDef | undefined,
          results: funcDef.results.map(r => {
            const result: any = {
              name: r.argumentName,
              sourceName: r.argumentName,
              types: r.types || [],
            }

            // Mark structurally complex types as expandable
            if (result.types.some((t: string) => EXPANDABLE_TYPES.includes(t))) {
              result.expandable = true

              if (result.types.includes('ComposeRecord')) {
                result.namespaceID = getFunctionArgumentValue(node.data, 'namespace')
                result.moduleID = getFunctionArgumentValue(node.data, 'module')
              }
            }

            // The Run Workflow result is an opaque Vars blob; expand it into the
            // selected workflow's declared outputs (results.<field>).
            if (funcDef.ref === 'workflowExec') {
              result.expandable = true
              result.workflowID = getFunctionArgumentValue(node.data, 'workflow')
            }

            return result
          }),
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
      newNodeType.id === 'branch' ||
      newNodeType.type === 'condition' ||
      newNodeType.ref === 'gateway'
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
      const step = automation.value.steps?.find((s: any) => s.stepID === oldNode.data.stepID)
      if (step) {
        step.kind = isIterator ? 'iterator' : isBranch ? gatewayRef : 'function'
        step.ref = isIterator
          ? newNodeType.ref || ''
          : isBranch
            ? gatewayRef
            : newNodeType.ref || ''
        step.arguments = []
        step.meta = {
          short: newNodeType.label,
          description: newNodeType.description || '',
        }
      }
    }

    // Seed default arguments from function/trigger segments (mirrors addNode logic)
    const store = useAutomationStore()
    let defaultArguments: any[] = []
    let defaultConfig: Record<string, any> = {}

    if (isTrigger) {
      const dbTrigger = store.triggers.find(
        t =>
          t.eventType === (newNodeType.eventType || newNodeType.ref) &&
          (!newNodeType.resourceType || t.resourceType === newNodeType.resourceType),
      )
      if (dbTrigger?.segments) {
        dbTrigger.segments.forEach(seg => {
          seg.sections?.forEach(sec => {
            sec.elements?.forEach(el => {
              const elInput = el.input as any
              if (elInput && elInput.default !== undefined && elInput.argument) {
                defaultConfig[elInput.argument] = elInput.default
              }
            })
          })
        })
      }
    } else if (!isBranch) {
      const dbFunction = store.functions.find(f => f.ref === newNodeType.ref)
      if (dbFunction?.segments) {
        dbFunction.segments.forEach(seg => {
          seg.sections?.forEach(sec => {
            sec.elements?.forEach(el => {
              const elInput = el.input as any
              if (elInput && elInput.default !== undefined && elInput.argument) {
                const param = dbFunction.parameters?.find(p => p.argumentName === elInput.argument)
                defaultArguments.push({
                  argumentName: elInput.argument,
                  type: param?.types?.[0] || 'Any',
                  value: elInput.default,
                })
              }
            })
          })
        })
      }
    }

    // Build updated node data
    const newData: FlowNodeData = {
      ...oldNode.data,
      label: newNodeType.label,
      description: newNodeType.description,
      icon: newNodeType.icon,
      nodeType: newDataNodeType,
      arguments: [...defaultArguments],
      config: { ...defaultConfig },
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

  // Get all expected properties from all triggers in the flow
  function getTriggerProperties() {
    const store = useAutomationStore()
    const allProperties = new Map<string, any>()

    // Helper: resolve a trigger constraint value by property name
    function getTriggerConstraintValue(nodeData: any, propName: string): string | null {
      const constraints = nodeData?.constraints || []
      const c = constraints.find((cc: any) => cc.name === propName)
      return c?.values?.[0]?.['@value'] ?? null
    }

    nodes.value
      .filter(n => n.type === 'trigger')
      .forEach(node => {
        const eventType = node.data?.nodeType
        const resourceType = node.data?.resourceType
        const triggerDef = store.triggers.find(
          t => t.eventType === eventType && (!resourceType || t.resourceType === resourceType),
        )

        // Triggers that declare their own input schema (agentic + manual schema triggers)
        if (
          eventType === 'onAgentic' ||
          resourceType === 'automation:trigger:agentic' ||
          resourceType?.startsWith('automation:trigger-definition:')
        ) {
          const schema: Array<{
            name: string
            type: string
            required?: boolean
            description?: string
          }> = node.data?.inputSchema || []
          schema
            .filter(p => p.name)
            .forEach(p => {
              if (!allProperties.has(p.name)) {
                allProperties.set(p.name, {
                  name: p.name,
                  type: p.type || 'String',
                  required: !!p.required,
                  meta: { short: p.name, description: p.description || '' },
                })
              }
            })
          return
        }

        if (triggerDef?.properties) {
          triggerDef.properties.forEach((p: any) => {
            // Store by name to deduplicate overlapping properties across triggers
            if (!allProperties.has(p.name)) {
              const prop: any = {
                name: p.name,
                type: p.type || 'String',
                defaultValue: getTriggerConstraintValue(node.data, p.name),
                meta: p.meta || {},
              }

              // Extract namespace and module if it's a ComposeRecord
              if (prop.type === 'ComposeRecord') {
                prop.namespaceID = getTriggerConstraintValue(node.data, 'namespace')
                prop.moduleID = getTriggerConstraintValue(node.data, 'module')
              }

              allProperties.set(p.name, prop)
            }
          })
        }
      })

    const propertyOrder = ['namespace', 'module', 'record', 'oldRecord', 'user', 'oldUser']
    return Array.from(allProperties.values()).sort((a, b) => {
      const indexA = propertyOrder.indexOf(a.name)
      const indexB = propertyOrder.indexOf(b.name)

      if (indexA !== -1 && indexB !== -1) return indexA - indexB
      if (indexA !== -1) return -1
      if (indexB !== -1) return 1

      return a.name.localeCompare(b.name)
    })
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
    isDirty,

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
    getTriggerProperties,
  }
}
