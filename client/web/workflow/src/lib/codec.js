export function encodeGraph (model, vertices, edges) {
  const steps = []
  const paths = {}
  const triggers = []

  Object.values(model.cells)
    .filter(cell => {
      return !!cell.vertex
    }).forEach(cell => {
      const triggerEdges = []

      const defaultName = vertices[cell.id].config.defaultName || false

      cell = {
        id: cell.id,
        value: cell.value,
        defaultName,
        xywh: [
          cell.geometry.x || 0,
          cell.geometry.y || 0,
          cell.geometry.width || 0,
          cell.geometry.height || 0,
        ],
        parent: cell.parent.id,
        edges: (cell.edges || []).forEach(({ id, value, parent, source, target, geometry, style }) => {
          const edge = {
            ...((edges[id] || {}).config || {}),
            parentID: source.id,
            childID: target.id,
            meta: {
              label: value || '',
              description: '',
              visual: {
                id,
                value,
                parent: parent.id,
                points: geometry.points,
                style,
              },
            },
          }

          if (vertices[source.id].triggers || vertices[target.id].triggers) {
            triggerEdges.push(edge)
          } else if (!paths[id]) {
            paths[id] = edge
          }
        }),
      }

      if (vertices[cell.id].triggers) {
        cell.edges = triggerEdges
        triggers.push({
          ...vertices[cell.id].triggers,
          stepID: (triggerEdges[0] || { childID: '0' }).childID,
          enabled: vertices[cell.id].triggers.enabled,
          constraints: vertices[cell.id].triggers.constraints,
          meta: {
            name: cell.value || '',
            description: '',
            visual: cell,
          },
        })
      } else {
        steps.push({
          ...vertices[cell.id].config,
          meta: {
            label: cell.value || '',
            description: '',
            visual: cell,
          },
        })
      }
    })

  return { steps, paths: Object.values(paths), triggers }
}
