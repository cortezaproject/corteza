import { ref, watch, nextTick } from 'vue'
import { useVueFlow } from '@vue-flow/core'

/**
 * Undo/redo stack for workflow nodes + edges — delta-compressed storage.
 *
 * History entries are either a full snapshot or a delta against the previous
 * entry. A full checkpoint is stored at entry 0 and every CHECKPOINT_EVERY
 * entries thereafter; all other entries are compact deltas (added/removed/
 * changed sets). Reconstructing state at any index costs at most
 * CHECKPOINT_EVERY delta applications, and unchanged nodes between snapshots
 * share the same frozen reference so they take no extra memory.
 *
 * On restore we reuse the existing VueFlow node/edge objects where the id
 * matches so VueFlow's measured fields (dimensions, handleBounds,
 * computedPosition) survive the round-trip. After each restore we call
 * `updateNodeInternals()` for nodes whose type or data changed so handle
 * positions stay in sync.
 *
 * Max MAX entries total.
 *
 * @param {import('vue').Ref<Array>} nodes reactive VueFlow nodes array
 * @param {import('vue').Ref<Array>} edges reactive VueFlow edges array
 * @param {Object} [opts]
 * @param {string} [opts.vfId] VueFlow instance id
 * @param {(ids?: string[]) => void} [opts.updateNodeInternals] explicit
 *   remeasure callback; overrides the auto-resolved VueFlow one.
 */
export function useWorkflowHistory(nodes, edges, opts = {}) {
  // Each entry:
  //   { type: 'full',  nodes: Map<id, frozenSnap>, edges: Map<id, frozenSnap> }
  //   { type: 'delta', nodes: Delta, edges: Delta }
  //   Delta = { added: frozenSnap[], removed: string[], changed: frozenSnap[] }
  const history = ref([])
  const pointer = ref(-1)
  const skipwatch = ref(false)
  const MAX = 100
  const CHECKPOINT_EVERY = 10

  // After truncation (undo + new edit) we force a full checkpoint so the new
  // branch never depends on the discarded delta chain.
  let needsCheckpoint = false

  const vfId = opts.vfId ?? 'workflow-editor-flow'
  let resolvedUpdate =
    typeof opts.updateNodeInternals === 'function' ? opts.updateNodeInternals : null

  function remeasure(ids) {
    if (!resolvedUpdate) {
      try {
        const vf = useVueFlow(vfId)
        if (vf && typeof vf.updateNodeInternals === 'function') {
          resolvedUpdate = vf.updateNodeInternals.bind(vf)
        }
      } catch (_e) {
        // VueFlow instance not ready yet — skip silently.
      }
    }
    if (resolvedUpdate) {
      resolvedUpdate(ids ? [...ids] : undefined)
    }
  }

  function snapshotNode(n) {
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

  function snapshotEdge(e) {
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

  // Build a Map<id, frozenSnap> from a reactive array.
  function toMap(arr, snapshotFn) {
    return new Map(arr.map(x => [x.id, Object.freeze(snapshotFn(x))]))
  }

  // Compute what changed between two Maps. Returns a Delta.
  function computeDelta(prevMap, currMap) {
    const added = []
    const removed = []
    const changed = []
    for (const [id, snap] of currMap) {
      if (!prevMap.has(id)) {
        added.push(snap)
      } else if (JSON.stringify(prevMap.get(id)) !== JSON.stringify(snap)) {
        changed.push(snap)
      }
    }
    for (const id of prevMap.keys()) {
      if (!currMap.has(id)) removed.push(id)
    }
    return { added, removed, changed }
  }

  // Apply a Delta to a Map, returning a new Map (non-mutating).
  function applyDelta(map, delta) {
    const next = new Map(map)
    for (const id of delta.removed) next.delete(id)
    for (const snap of delta.added) next.set(snap.id, snap)
    for (const snap of delta.changed) next.set(snap.id, snap)
    return next
  }

  // Reconstruct full node/edge Maps at the given history index.
  // Walks back to the nearest full checkpoint, then replays deltas forward.
  function reconstructAt(index) {
    let base = index
    while (base > 0 && history.value[base].type !== 'full') base--
    let nodeMap = new Map(history.value[base].nodes)
    let edgeMap = new Map(history.value[base].edges)
    for (let i = base + 1; i <= index; i++) {
      nodeMap = applyDelta(nodeMap, history.value[i].nodes)
      edgeMap = applyDelta(edgeMap, history.value[i].edges)
    }
    return { nodeMap, edgeMap }
  }

  function restore({ nodeMap, edgeMap }) {
    skipwatch.value = true

    const currentNodeById = new Map(nodes.value.map(n => [n.id, n]))
    const currentEdgeById = new Map(edges.value.map(e => [e.id, e]))

    const toRemeasure = []

    nodes.value = [...nodeMap.values()].map(s => {
      const existing = currentNodeById.get(s.id)
      if (existing) {
        const typeChanged = existing.type !== s.type
        const dataChanged = JSON.stringify(existing.data || {}) !== JSON.stringify(s.data || {})
        if (typeChanged || dataChanged) toRemeasure.push(s.id)
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
      toRemeasure.push(s.id)
      return { ...s }
    })

    edges.value = [...edgeMap.values()].map(s => {
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

    nextTick(() => {
      remeasure(toRemeasure.length > 0 ? toRemeasure : undefined)
    })
  }

  function saveToHistory() {
    if (skipwatch.value) return

    const currNodeMap = toMap(nodes.value, snapshotNode)
    const currEdgeMap = toMap(edges.value, snapshotEdge)

    // Discard any future entries when branching from an undo point.
    if (pointer.value < history.value.length - 1) {
      history.value = history.value.slice(0, pointer.value + 1)
      // Force a full checkpoint so the new branch never references the
      // discarded chain.
      needsCheckpoint = true
    }

    const newIndex = history.value.length

    if (newIndex === 0 || needsCheckpoint || newIndex % CHECKPOINT_EVERY === 0) {
      history.value.push({ type: 'full', nodes: currNodeMap, edges: currEdgeMap })
      needsCheckpoint = false
    } else {
      const { nodeMap: prevNodeMap, edgeMap: prevEdgeMap } = reconstructAt(newIndex - 1)
      history.value.push({
        type: 'delta',
        nodes: computeDelta(prevNodeMap, currNodeMap),
        edges: computeDelta(prevEdgeMap, currEdgeMap),
      })
    }

    // Trim oldest entries while maintaining the invariant that history[0] is
    // always a full snapshot.
    if (history.value.length > MAX) {
      const excess = history.value.length - MAX
      const { nodeMap, edgeMap } = reconstructAt(excess)
      history.value = [
        { type: 'full', nodes: nodeMap, edges: edgeMap },
        ...history.value.slice(excess + 1),
      ]
    }

    pointer.value = history.value.length - 1
  }

  // Reset to a clean single-entry history using the current nodes/edges as the
  // new base. Called after the workflow is initially loaded so that the loaded
  // state is the floor — undo is disabled until the user makes a real change.
  function resetHistory() {
    const currNodeMap = toMap(nodes.value, snapshotNode)
    const currEdgeMap = toMap(edges.value, snapshotEdge)
    history.value = [{ type: 'full', nodes: currNodeMap, edges: currEdgeMap }]
    needsCheckpoint = false
    pointer.value = 0
  }

  function undo() {
    if (pointer.value > 0) {
      pointer.value--
      restore(reconstructAt(pointer.value))
    }
  }

  function redo() {
    if (pointer.value < history.value.length - 1) {
      pointer.value++
      restore(reconstructAt(pointer.value))
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
    resetHistory,
    undo,
    redo,
    canUndo,
    canRedo,
  }
}
