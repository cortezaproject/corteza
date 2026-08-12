/**
 * Ctrl+click path highlighting composable.
 * Sets node.data.highlighted and edge.data.highlighted on connected paths.
 */
export function useWorkflowHighlight(nodes, edges) {
  function highlightConnected(nodeId) {
    // Find the node
    const node = nodes.value.find(n => n.id === nodeId)
    if (!node) return

    // Toggle highlight on the node
    node.data = { ...node.data, highlighted: !node.data.highlighted }

    // Highlight connected edges
    edges.value.forEach(edge => {
      if (edge.source === nodeId || edge.target === nodeId) {
        edge.data = { ...edge.data, highlighted: node.data.highlighted }
      }
    })

    // Highlight connected nodes via edges
    edges.value
      .filter(e => e.source === nodeId || e.target === nodeId)
      .forEach(edge => {
        const connectedId = edge.source === nodeId ? edge.target : edge.source
        const connectedNode = nodes.value.find(n => n.id === connectedId)
        if (connectedNode) {
          connectedNode.data = { ...connectedNode.data, highlighted: node.data.highlighted }
        }
      })
  }

  function clearHighlights() {
    nodes.value.forEach(node => {
      if (node.data?.highlighted) {
        node.data = { ...node.data, highlighted: false }
      }
    })
    edges.value.forEach(edge => {
      if (edge.data?.highlighted) {
        edge.data = { ...edge.data, highlighted: false }
      }
    })
  }

  return {
    highlightConnected,
    clearHighlights,
  }
}
