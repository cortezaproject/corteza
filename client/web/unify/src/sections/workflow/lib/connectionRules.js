/**
 * Connection / outbound-cap rules for the workflow graph. Pure functions so
 * the editor canvas (`WorkflowEditor.vue`), the node component
 * (`WorkflowNode.vue`), and unit tests all share one source of truth.
 */

export function getMaxOutbound(node) {
  if (!node) return 0
  const kind = node.data?.kind
  const ref = node.data?.ref
  if (kind === 'gateway' && ['fork', 'excl', 'incl'].includes(ref)) return Infinity
  if (kind === 'iterator' || kind === 'error-handler') return 2
  return 1
}

export function isValidConnection(
  connection,
  { nodes = [], edges = [], edgeUpdatingId = null } = {},
) {
  const sourceNode = nodes.find(n => n.id === connection.source)
  const targetNode = nodes.find(n => n.id === connection.target)
  const updatingId = edgeUpdatingId

  if (sourceNode?.type === 'visual' || targetNode?.type === 'visual') return false
  if (targetNode?.type === 'trigger') return false
  if (sourceNode?.type === 'termination') return false

  if (targetNode?.type === 'termination') {
    const existingIn = edges.filter(
      e => e.target === connection.target && e.id !== connection.id && e.id !== updatingId,
    ).length
    if (existingIn >= 1) return false
  }

  if (sourceNode && sourceNode.type !== 'trigger') {
    const maxOut = getMaxOutbound(sourceNode)
    const existingOut = edges.filter(
      e => e.source === connection.source && e.id !== connection.id && e.id !== updatingId,
    ).length
    if (existingOut >= maxOut) return false
  }

  if (sourceNode?.type === 'trigger') {
    const existing = edges.some(
      e => e.source === connection.source && e.id !== connection.id && e.id !== updatingId,
    )
    if (existing) return false
  }

  if (connection.source && connection.source === connection.target) return false

  if (connection.source && connection.sourceHandle) {
    const existingFromHandle = edges.some(
      e =>
        e.source === connection.source &&
        e.sourceHandle === connection.sourceHandle &&
        e.id !== connection.id &&
        e.id !== updatingId,
    )
    if (existingFromHandle) return false
  }

  if (connection.target && connection.targetHandle) {
    const existingToHandle = edges.some(
      e =>
        e.target === connection.target &&
        e.targetHandle === connection.targetHandle &&
        e.id !== connection.id &&
        e.id !== updatingId,
    )
    if (existingToHandle) return false
  }

  return true
}
