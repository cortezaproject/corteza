import { ref, watch } from 'vue'

/**
 * Undo/redo stack for workflow nodes + edges.
 *
 * We snapshot ONLY the persistent, user-meaningful fields — position, data,
 * type, handles, labels, selection — and strip VueFlow's runtime-measured
 * fields (dimensions, handleBounds, computedPosition, etc.). On restore we
 * reuse the existing VueFlow node/edge objects where the id matches so those
 * measured fields survive the round-trip; only the snapshotted fields are
 * overwritten. That prevents the "nodes teleport / handles misalign after
 * undo" drift we used to see when we reassigned whole arrays built from JSON.
 *
 * Max 100 entries.
 */
export function useWorkflowHistory (nodes, edges) {
  const history = ref([])
  const pointer = ref(-1)
  const skipwatch = ref(false)
  const MAX = 100

  function snapshotNode (n) {
    return {
      id: n.id,
      type: n.type,
      position: { x: n.position?.x ?? 0, y: n.position?.y ?? 0 },
      zIndex: n.zIndex,
      connectable: n.connectable,
      parentNode: n.parentNode,
      extent: n.extent,
      selected: !!n.selected,
      data: JSON.parse(JSON.stringify(n.data || {})),
    }
  }

  function snapshotEdge (e) {
    return {
      id: e.id,
      source: e.source,
      target: e.target,
      sourceHandle: e.sourceHandle ?? null,
      targetHandle: e.targetHandle ?? null,
      type: e.type,
      label: e.label ?? '',
      selected: !!e.selected,
      data: JSON.parse(JSON.stringify(e.data || {})),
    }
  }

  function snapshot () {
    return {
      nodes: nodes.value.map(snapshotNode),
      edges: edges.value.map(snapshotEdge),
    }
  }

  function restore (snap) {
    skipwatch.value = true

    const currentNodeById = new Map(nodes.value.map(n => [n.id, n]))
    const currentEdgeById = new Map(edges.value.map(e => [e.id, e]))

    nodes.value = snap.nodes.map(s => {
      const existing = currentNodeById.get(s.id)
      if (existing) {
        // Preserve VueFlow's measured fields (dimensions, handleBounds, etc.)
        // and just overwrite the snapshot-managed ones.
        return Object.assign(existing, {
          type: s.type,
          position: { ...s.position },
          zIndex: s.zIndex,
          connectable: s.connectable,
          parentNode: s.parentNode,
          extent: s.extent,
          selected: s.selected,
          data: s.data,
        })
      }
      return { ...s }
    })

    edges.value = snap.edges.map(s => {
      const existing = currentEdgeById.get(s.id)
      if (existing) {
        return Object.assign(existing, {
          source: s.source,
          target: s.target,
          sourceHandle: s.sourceHandle,
          targetHandle: s.targetHandle,
          type: s.type,
          label: s.label,
          selected: s.selected,
          data: s.data,
        })
      }
      return { ...s }
    })

    skipwatch.value = false
  }

  function saveToHistory () {
    if (skipwatch.value) return
    const snap = snapshot()

    if (pointer.value < history.value.length - 1) {
      history.value = history.value.slice(0, pointer.value + 1)
    }

    history.value.push(snap)

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

  saveToHistory()

  return {
    saveToHistory,
    undo,
    redo,
    canUndo,
    canRedo,
  }
}
