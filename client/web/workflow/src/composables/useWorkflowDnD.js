import { ref } from 'vue'
import { getStyleFromKind } from '../lib/style'

/**
 * Drag-from-toolbar-to-canvas composable.
 * On drag start, stores the toolbar item configuration.
 * On canvas drop, creates a VueFlow node at the drop position.
 *
 * @param {Ref} nodes - reactive nodes array
 * @param {Function} saveToHistory - history snapshot function
 * @param {Function} screenToFlowPosition - from useVueFlow() in the parent component
 */
export function useWorkflowDnD (nodes, saveToHistory, screenToFlowPosition) {

  const draggedItem = ref(null)

  function onToolbarDragStart (event, toolbarItem) {
    draggedItem.value = toolbarItem
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('application/workflow-node', JSON.stringify(toolbarItem))
  }

  function onCanvasDragOver (event) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
  }

  function onCanvasDrop (event) {
    event.preventDefault()

    let item = draggedItem.value
    if (!item) {
      try {
        item = JSON.parse(event.dataTransfer.getData('application/workflow-node'))
      } catch {
        return
      }
    }

    const styleInfo = getStyleFromKind(item) || {}

    // Convert screen coordinates to flow canvas coordinates
    let position
    const toFlowPos = typeof screenToFlowPosition === 'function'
      ? screenToFlowPosition
      : (screenToFlowPosition?.value || null)

    if (toFlowPos) {
      position = toFlowPos({
        x: event.clientX,
        y: event.clientY,
      })
    } else {
      // Fallback: use client coordinates relative to the canvas container
      const container = event.currentTarget || event.target
      const rect = container.getBoundingClientRect()
      position = {
        x: event.clientX - rect.left,
        y: event.clientY - rect.top,
      }
    }

    // 8px grid snap
    position.x = Math.round(position.x / 8) * 8
    position.y = Math.round(position.y / 8) * 8

    // Generate unique ID
    const id = String(Date.now())

    let nodeType = 'workflow'
    if (item.kind === 'trigger') nodeType = 'trigger'
    else if (item.kind === 'termination') nodeType = 'termination'
    else if (item.kind === 'visual') nodeType = 'visual'

    // Determine default label
    let label = ''
    if (['break', 'continue'].includes(item.kind)) {
      label = item.kind === 'break' ? 'Stop iterator execution' : 'Skip current iteration'
    } else if (item.kind === 'gateway') {
      label = item.ref || ''
    } else if (item.kind === 'visual' && item.ref === 'content') {
      label = 'Text here'
    } else if (item.kind === 'expressions') {
      label = 'Define and mutate scope variables'
    }

    const newNode = {
      id,
      type: nodeType,
      position,
      zIndex: nodeType === 'visual' ? -1 : undefined,
      connectable: nodeType === 'visual' ? false : undefined,
      data: {
        stepID: id,
        kind: item.kind || '',
        ref: item.ref || '',
        label,
        description: '',
        arguments: [],
        results: [],
        defaultName: true,
        highlighted: false,
        traceState: null,
        traceLog: null,
        width: styleInfo.width || 200,
        height: styleInfo.height || 80,
        // Trigger-specific
        ...(item.kind === 'trigger'
          ? {
              triggers: {
                resourceType: null,
                eventType: null,
                constraints: [],
                enabled: true,
              },
            }
          : {}),
      },
    }

    nodes.value = [...nodes.value, newNode]

    if (saveToHistory) {
      saveToHistory()
    }

    draggedItem.value = null

    return newNode
  }

  return {
    draggedItem,
    onToolbarDragStart,
    onCanvasDragOver,
    onCanvasDrop,
  }
}
