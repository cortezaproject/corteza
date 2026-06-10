// Helpers for deriving the project graph from per-resource `links`.
//
// Every resource owns its outgoing links; there is no central `edges` array.
// These helpers produce edge-shaped objects ({ source, target, type, label? })
// that the relationship graph can feed straight into EchartsGraph.

export function allLinks(project) {
  if (!project) return []
  const out = []
  for (const r of project.resources || []) {
    for (const l of r.links || []) {
      out.push({
        linkId: l.id,
        source: r.id,
        target: l.targetId,
        type: l.type,
        label: l.label,
      })
    }
  }
  return out
}

// All resource ids linked to `resourceId`, in either direction (links are
// undirected; they're just stored once on one side).
export function linkedResourceIds(project, resourceId) {
  const ids = new Set()
  if (!project || !resourceId) return []
  const self = (project.resources || []).find(r => r.id === resourceId)
  for (const l of self?.links || []) ids.add(l.targetId)
  for (const r of project.resources || []) {
    if (r.id === resourceId) continue
    for (const l of r.links || []) {
      if (l.targetId === resourceId) ids.add(r.id)
    }
  }
  return [...ids]
}
