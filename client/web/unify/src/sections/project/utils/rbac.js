// Role → resource permissions.
//
// project.roleBindings holds *explicit* overrides: [{ roleId, resourceId,
// level: 'none' | 'read' | 'write' }]. The effective level falls back to a
// group-membership baseline: a role that shares a group with a resource has
// `read` unless an explicit binding says otherwise. Write implies read.

export const LEVELS = ['read', 'write']
export const RANK = { none: 0, read: 1, write: 2 }

export const PERM_COLORS = {
  read: { stroke: '#94a3b8', label: 'Read' }, // slate-400
  write: { stroke: '#2563eb', label: 'Write' }, // blue-600
}

// Resource kinds that can be permission targets (not roles/users).
export const isPermissionTarget = kind => kind !== 'role' && kind !== 'user'

export function sharesGroup(project, roleId, resourceId) {
  for (const g of project?.groups || []) {
    if ((g.roleIds || []).includes(roleId) && (g.resourceIds || []).includes(resourceId)) {
      return true
    }
  }
  return false
}

function explicitLevel(project, roleId, resourceId) {
  const b = (project?.roleBindings || []).find(
    x => x.roleId === roleId && x.resourceId === resourceId,
  )
  return b ? b.level : undefined
}

// 'none' | 'read' | 'write'. Explicit binding wins; otherwise the group baseline.
export function effectiveLevel(project, roleId, resourceId) {
  const ex = explicitLevel(project, roleId, resourceId)
  if (ex !== undefined) return ex
  return sharesGroup(project, roleId, resourceId) ? 'read' : 'none'
}

// True when the effective level comes from the group baseline, not an explicit
// binding — used to label the selector ("Read · from group").
export function isBaseline(project, roleId, resourceId) {
  return (
    explicitLevel(project, roleId, resourceId) === undefined &&
    sharesGroup(project, roleId, resourceId)
  )
}

// Every role→resource edge with a non-none effective level, for the graph.
export function permissionEdges(project, includeId = () => true) {
  const resources = project?.resources || []
  const roles = resources.filter(r => r.kind === 'role')
  const targets = resources.filter(r => isPermissionTarget(r.kind))
  const out = []
  for (const role of roles) {
    if (!includeId(role.id)) continue
    for (const res of targets) {
      if (!includeId(res.id)) continue
      const level = effectiveLevel(project, role.id, res.id)
      if (level === 'none') continue
      out.push({ source: role.id, target: res.id, level })
    }
  }
  return out
}
