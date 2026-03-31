/**
 * Copy / Cut / Paste composable for workflow nodes and edges.
 * Copies selected nodes (and their connecting edges), remaps IDs on paste,
 * and offsets positions by 160px.
 */
export function useWorkflowClipboard (nodes, edges, saveToHistory) {
  let clipboard = null

  function getSelected () {
    return nodes.value.filter(n => n.selected)
  }

  function copySelected () {
    const selectedNodes = getSelected()
    if (!selectedNodes.length) return

    const selectedIds = new Set(selectedNodes.map(n => n.id))
    const connectedEdges = edges.value.filter(
      e => selectedIds.has(e.source) && selectedIds.has(e.target),
    )

    clipboard = {
      nodes: JSON.parse(JSON.stringify(selectedNodes)),
      edges: JSON.parse(JSON.stringify(connectedEdges)),
    }
  }

  function cutSelected () {
    const selectedNodes = getSelected()
    if (!selectedNodes.length) return

    copySelected()

    const selectedIds = new Set(selectedNodes.map(n => n.id))
    nodes.value = nodes.value.filter(n => !selectedIds.has(n.id))
    edges.value = edges.value.filter(
      e => !selectedIds.has(e.source) && !selectedIds.has(e.target),
    )

    if (saveToHistory) saveToHistory()
  }

  function pasteClipboard () {
    if (!clipboard) return

    const idMap = {}
    const offset = 160

    // Create new nodes with remapped IDs
    const newNodes = clipboard.nodes.map(n => {
      const newId = String(Date.now() + Math.random())
      idMap[n.id] = newId

      return {
        ...n,
        id: newId,
        selected: true,
        position: {
          x: (n.position?.x || 0) + offset,
          y: (n.position?.y || 0) + offset,
        },
        data: {
          ...n.data,
          stepID: newId,
        },
      }
    })

    // Create new edges with remapped IDs
    const newEdges = clipboard.edges
      .filter(e => idMap[e.source] && idMap[e.target])
      .map(e => ({
        ...e,
        id: `e-${Date.now()}-${Math.random()}`,
        source: idMap[e.source],
        target: idMap[e.target],
        data: {
          ...e.data,
          parentID: idMap[e.source],
          childID: idMap[e.target],
        },
      }))

    // Deselect existing nodes
    nodes.value = nodes.value.map(n => ({ ...n, selected: false }))

    // Add new
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
