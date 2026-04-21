import { ref, watch, nextTick } from 'vue'
import { useVueFlow } from '@vue-flow/core'

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
 * After each restore we also ask VueFlow to re-measure the affected nodes via
 * `updateNodeInternals()`. That catches the cases that a field-whitelist
 * cannot preserve — e.g. a previously-deleted node being re-added (where no
 * measured state exists yet), or a data/type change that shifts handle
 * positions and therefore invalidates cached handleBounds.
 *
 * Max 100 entries.
 *
 * @param {import('vue').Ref<Array>} nodes reactive VueFlow nodes array
 * @param {import('vue').Ref<Array>} edges reactive VueFlow edges array
 * @param {Object} [opts]
 * @param {string} [opts.vfId] VueFlow instance id (used to locate
 *   `updateNodeInternals`). Defaults to `'workflow-editor-flow'` to match the
 *   id the workflow editor registers. Ignored if `opts.updateNodeInternals`
 *   is provided directly.
 * @param {(ids?: string[]) => void} [opts.updateNodeInternals] explicit
 *   remeasure callback; overrides the auto-resolved VueFlow one. Useful for
 *   tests.
 */
export function useWorkflowHistory (nodes, edges, opts = {}) {
  const history = ref([])
  const pointer = ref(-1)
  const skipwatch = ref(false)
  const MAX = 100

  // Lazily resolve VueFlow's updateNodeInternals. We don't want to throw if
  // the flow instance isn't registered yet at composable-init time (it often
  // isn't — the <VueFlow> component mounts after its parent's setup runs).
  const vfId = opts.vfId ?? 'workflow-editor-flow'
  let resolvedUpdate = typeof opts.updateNodeInternals === 'function'
    ? opts.updateNodeInternals
    : null

  function remeasure (ids) {
    if (!resolvedUpdate) {
      try {
        const vf = useVueFlow(vfId)
        if (vf && typeof vf.updateNodeInternals === 'function') {
          resolvedUpdate = vf.updateNodeInternals.bind(vf)
        }
      } catch (_e) {
        // VueFlow instance not ready yet — skip silently, next restore will
        // try again.
      }
    }
    if (resolvedUpdate) {
      // Pass a defensive copy so downstream mutations can't feed back.
      resolvedUpdate(ids ? [...ids] : undefined)
    }
  }

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

    // Track which nodes need a VueFlow re-measure after restore:
    //   - brand new (no existing instance to preserve)
    //   - type changed (handle layout may differ)
    //   - data changed in a way that can alter handle count/position
    const toRemeasure = []

    nodes.value = snap.nodes.map(s => {
      const existing = currentNodeById.get(s.id)
      if (existing) {
        const typeChanged = existing.type !== s.type
        const dataChanged = JSON.stringify(existing.data || {}) !== JSON.stringify(s.data || {})
        if (typeChanged || dataChanged) {
          toRemeasure.push(s.id)
        }
        // Preserve VueFlow's measured fields (dimensions, handleBounds,
        // computedPosition, etc.) and just overwrite the snapshot-managed
        // ones.
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
      // Brand-new (or previously-deleted) node — VueFlow has no measured
      // state for it and needs to run its internal measure pass.
      toRemeasure.push(s.id)
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

    // Ask VueFlow to re-measure after Vue flushes the reactive updates so the
    // DOM reflects the restored state. Without this, handleBounds /
    // computedPosition can stay stale and subsequent interactions (dragging a
    // connection, moving the node) exhibit the "jitter" the editor used to
    // show after undo.
    nextTick(() => {
      remeasure(toRemeasure.length > 0 ? toRemeasure : undefined)
    })
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
