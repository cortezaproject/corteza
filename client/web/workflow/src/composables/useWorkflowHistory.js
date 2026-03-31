import { ref, watch } from 'vue'

/**
 * Undo/redo history stack for workflow nodes and edges.
 * Takes a snapshot of the full nodes + edges arrays after every call to saveToHistory().
 * Max 100 entries.
 */
export function useWorkflowHistory (nodes, edges) {
  const history = ref([])
  const pointer = ref(-1)
  const skipwatch = ref(false)
  const MAX = 100

  function snapshot () {
    return JSON.stringify({ nodes: nodes.value, edges: edges.value })
  }

  function restore (json) {
    skipwatch.value = true
    const data = JSON.parse(json)
    nodes.value = data.nodes
    edges.value = data.edges
    skipwatch.value = false
  }

  function saveToHistory () {
    const snap = snapshot()

    // Drop any redo entries after current pointer
    if (pointer.value < history.value.length - 1) {
      history.value = history.value.slice(0, pointer.value + 1)
    }

    history.value.push(snap)

    // Enforce max
    if (history.value.length > MAX) {
      history.value = history.value.slice(history.value.length - MAX)
    }

    pointer.value = history.value.length - 1
  }

  function undo () {
    if (pointer.value > 0) {
      pointer.value--
      restore(history.value[pointer.value])
    }
  }

  function redo () {
    if (pointer.value < history.value.length - 1) {
      pointer.value++
      restore(history.value[pointer.value])
    }
  }

  const canUndo = ref(false)
  const canRedo = ref(false)

  watch(pointer, () => {
    canUndo.value = pointer.value > 0
    canRedo.value = pointer.value < history.value.length - 1
  })

  // Take initial snapshot
  saveToHistory()

  return {
    saveToHistory,
    undo,
    redo,
    canUndo,
    canRedo,
  }
}
