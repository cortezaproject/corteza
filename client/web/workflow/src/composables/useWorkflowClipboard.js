/**
 * Copy / Cut / Paste composable for workflow nodes and edges.
 *
 * Behaviour parity with Corteza's mxGraph clipboard:
 *   - selecting a swimlane implicitly copies every node nested inside it
 *   - edges between copied nodes are preserved; edges that cross the
 *     selection boundary are dropped
 *   - paste remaps IDs (incrementing integers, matching Corteza's scheme)
 *   - paste offsets positions by 160px so copies don't stack exactly on the
 *     original
 *   - clipboard payload is also written to the system clipboard as JSON with
 *     a magic prefix, so paste works across tabs
 */
import { nextId } from '../lib/id'

const CLIPBOARD_PREFIX = '__human_workflow_clipboard__:'

export function useWorkflowClipboard (nodes, edges, saveToHistory) {
  let localClipboard = null

  function getSelected () {
    return nodes.value.filter(n => n.selected)
  }

  function expandSelectionWithChildren (selectedIds) {
    // When a swimlane is selected, its nested children come along.
    const expanded = new Set(selectedIds)
    let changed = true
    while (changed) {
      changed = false
      nodes.value.forEach(n => {
        if (n.parentNode && expanded.has(n.parentNode) && !expanded.has(n.id)) {
          expanded.add(n.id)
          changed = true
        }
      })
    }
    return expanded
  }

  function buildPayload () {
    const selectedNodes = getSelected()
    if (!selectedNodes.length) return null

    const rootIds = new Set(selectedNodes.map(n => n.id))
    const allIds = expandSelectionWithChildren(rootIds)
    const allNodes = nodes.value.filter(n => allIds.has(n.id))
    const connectedEdges = edges.value.filter(
      e => allIds.has(e.source) && allIds.has(e.target),
    )

    return {
      nodes: JSON.parse(JSON.stringify(allNodes)),
      edges: JSON.parse(JSON.stringify(connectedEdges)),
    }
  }

  async function writeSystemClipboard (payload) {
    if (!navigator.clipboard?.writeText) return
    try {
      await navigator.clipboard.writeText(CLIPBOARD_PREFIX + JSON.stringify(payload))
    } catch {
      // Clipboard permissions can be revoked; silently fall back to local-only.
    }
  }

  async function readSystemClipboard () {
    if (!navigator.clipboard?.readText) return null
    try {
      const text = await navigator.clipboard.readText()
      if (!text || !text.startsWith(CLIPBOARD_PREFIX)) return null
      return JSON.parse(text.slice(CLIPBOARD_PREFIX.length))
    } catch {
      return null
    }
  }

  function copySelected () {
    const payload = buildPayload()
    if (!payload) return
    localClipboard = payload
    writeSystemClipboard(payload)
  }

  function cutSelected () {
    const payload = buildPayload()
    if (!payload) return
    localClipboard = payload
    writeSystemClipboard(payload)

    const removeIds = new Set(payload.nodes.map(n => n.id))
    nodes.value = nodes.value.filter(n => !removeIds.has(n.id))
    edges.value = edges.value.filter(
      e => !removeIds.has(e.source) && !removeIds.has(e.target),
    )

    if (saveToHistory) saveToHistory()
  }

  async function pasteClipboard () {
    const payload = (await readSystemClipboard()) || localClipboard
    if (!payload) return

    const idMap = {}
    const offset = 160

    let currentId = nextId(nodes, edges)

    const newNodes = payload.nodes.map(n => {
      const newId = String(currentId++)
      idMap[n.id] = newId
      return n
    }).map(n => {
      // Root-level nodes (those whose parent is not part of the payload) get
      // the paste offset. Nested children keep their relative position so the
      // swimlane layout is preserved.
      const parentMapped = n.parentNode && idMap[n.parentNode] ? idMap[n.parentNode] : undefined
      const applyOffset = !parentMapped
      return {
        ...n,
        id: idMap[n.id],
        selected: true,
        parentNode: parentMapped,
        extent: parentMapped ? 'parent' : undefined,
        position: {
          x: (n.position?.x || 0) + (applyOffset ? offset : 0),
          y: (n.position?.y || 0) + (applyOffset ? offset : 0),
        },
        data: {
          ...n.data,
          stepID: idMap[n.id],
        },
      }
    })

    const newEdges = payload.edges
      .filter(e => idMap[e.source] && idMap[e.target])
      .map(e => ({
        ...e,
        id: String(currentId++),
        source: idMap[e.source],
        target: idMap[e.target],
        data: {
          ...e.data,
          parentID: idMap[e.source],
          childID: idMap[e.target],
        },
      }))

    nodes.value = nodes.value.map(n => ({ ...n, selected: false }))
    nodes.value = [...nodes.value, ...newNodes]
    edges.value = [...edges.value, ...newEdges]

    if (saveToHistory) saveToHistory()
  }

  return {
    copySelected,
    cutSelected,
    pasteClipboard,
  }
}
