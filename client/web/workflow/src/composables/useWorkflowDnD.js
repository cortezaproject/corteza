import { ref } from 'vue'
import { getStyleFromKind } from '../lib/style'
import { nextId } from '../lib/id'

/**
 * Drag-from-toolbar-to-canvas composable.
 * On drag start, stores the toolbar item configuration.
 * On canvas drop, creates a VueFlow node at the drop position.
 *
 * @param {Ref} nodes - reactive nodes array
 * @param {Ref} edges - reactive edges array
 * @param {Function} saveToHistory - history snapshot function
 * @param {Function} projectPosition - VueFlow's project() fn (container-relative → flow coords)
 */
export function useWorkflowDnD (nodes, edges, saveToHistory, projectPosition) {

  const draggedItem = ref(null)

  function onToolbarDragStart (event, toolbarItem) {
    draggedItem.value = toolbarItem
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('application/workflow-node', JSON.stringify(toolbarItem))

    // Dashed node-sized preview — parity with Corteza's mxGraph drag ghost.
    // The element must be attached to the DOM for setDragImage to rasterise it,
    // so we append it off-screen, call setDragImage, then remove it on the next
    // tick (by which point the browser has snapshotted it).
    const styleInfo = getStyleFromKind(toolbarItem) || {}
    const w = styleInfo.width || 200
    const h = styleInfo.height || 80
    const ghost = document.createElement('div')
    ghost.style.position = 'absolute'
    ghost.style.top = '-1000px'
    ghost.style.left = '-1000px'
    ghost.style.width = `${w}px`
    ghost.style.height = `${h}px`
    ghost.style.border = '2px dashed var(--p-primary-color, #3b82f6)'
    ghost.style.borderRadius = '6px'
    ghost.style.background = 'color-mix(in srgb, var(--p-primary-color, #3b82f6) 10%, transparent)'
    ghost.style.boxSizing = 'border-box'
    ghost.style.pointerEvents = 'none'
    document.body.appendChild(ghost)
    try {
      event.dataTransfer.setDragImage(ghost, w / 2, h / 2)
    } catch {
      // setDragImage unsupported; fall back silently to the default preview
    }
    setTimeout(() => ghost.remove(), 0)
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
    const nodeWidth = styleInfo.width || 200
    const nodeHeight = styleInfo.height || 80

    // VueFlow's project() expects coordinates relative to the flow container,
    // not raw screen/client coordinates. Subtract the container's bounding rect.
    const flowContainer = event.currentTarget
    const rect = flowContainer ? flowContainer.getBoundingClientRect() : { left: 0, top: 0 }

    let position
    try {
      if (typeof projectPosition === 'function') {
        position = { ...projectPosition({
          x: event.clientX - rect.left,
          y: event.clientY - rect.top,
        }) }
      } else {
        throw new Error('project not available')
      }
    } catch {
      position = {
        x: event.clientX - rect.left,
        y: event.clientY - rect.top,
      }
    }

    // Center the node on the cursor
    position.x -= nodeWidth / 2
    position.y -= nodeHeight / 2

    // 8px grid snap
    position.x = Math.round(position.x / 8) * 8
    position.y = Math.round(position.y / 8) * 8

    // Generate incrementing integer ID (like Corteza's mxGraph)
    const id = String(nextId(nodes, edges))

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
        width: nodeWidth,
        height: nodeHeight,
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
