/**
 * Incrementing ID generator for workflow nodes and edges.
 *
 * Mirrors Human's mxGraph behavior where IDs are simple incrementing
 * integers: 2, 3, 4, ... (1 is reserved for the root/default parent).
 *
 * @param {Ref|Array} nodes - reactive nodes array (or .value)
 * @param {Ref|Array} [edges] - reactive edges array (or .value)
 * @returns {number} next available integer ID
 */
export function nextId (nodes, edges) {
  const nodeArr = Array.isArray(nodes) ? nodes : (nodes?.value || [])
  const edgeArr = Array.isArray(edges) ? edges : (edges?.value || [])

  let max = 1 // mxGraph reserves 1 for root

  const parse = (id) => {
    const n = parseInt(id, 10)
    if (!isNaN(n) && n > max) {
      max = n
    }
  }

  nodeArr.forEach(n => parse(n.id))
  edgeArr.forEach(e => parse(e.id))

  return max + 1
}
