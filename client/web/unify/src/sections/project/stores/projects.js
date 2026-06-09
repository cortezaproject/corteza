import { CURRENT_USER_ID } from '@/sections/project/mock/users'
import { MOCK_PROJECTS } from '@/sections/project/mock/projects'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

// Mock, in-memory store. Shape and mutations grow per-screen as we build.
export const useProjectsStore = defineStore('projects', () => {
  const projects = ref(MOCK_PROJECTS.map(p => ({ ...p })))

  const findById = computed(() => id => projects.value.find(p => p.id === id))

  // Latest published version, or null.
  const currentVersion = p => (p?.versions?.length ? p.versions[p.versions.length - 1] : null)

  const slugify = name =>
    (name || 'project')
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '') || 'project'

  const uniqueId = base => {
    let id = base
    let n = 2
    while (projects.value.some(p => p.id === id)) {
      id = `${base}-${n++}`
    }
    return id
  }

  // Create a fresh draft. The creator is automatically added as a Developer;
  // other members are configured in the Project Members step.
  // `deployer` holds the three AI Act Deployer-category answers gathered at
  // creation; any "yes" makes a Fundamental Rights Impact Assessment required.
  function create({ name, description = '', mode = 'free', deployer = {} } = {}) {
    const now = new Date().toISOString()

    const deployerCategories = {
      publicAuthorityAnnex3: !!deployer.publicAuthorityAnnex3,
      privateEssentialServices: !!deployer.privateEssentialServices,
      insuranceBanking: !!deployer.insuranceBanking,
    }
    const friaRequired = Object.values(deployerCategories).some(Boolean)

    const project = {
      id: uniqueId(slugify(name)),
      name: (name || '').trim() || 'Untitled project',
      description: description.trim(),
      mode,
      status: 'draft',
      createdBy: CURRENT_USER_ID,
      members: [{ id: localId('mbr'), userId: CURRENT_USER_ID, role: 'developer' }],
      versions: [],
      gatesApproved: 0,
      governance: {},
      deployerCategories,
      friaRequired,
      createdAt: now,
      updatedAt: now,
    }
    projects.value.push(project)
    return project
  }

  // Patch a project. `mode` is immutable and stripped from any patch.
  function updateProject(id, patch = {}) {
    const p = projects.value.find(x => x.id === id)
    if (!p) return
    const { mode: _ignored, ...rest } = patch
    Object.assign(p, rest, { updatedAt: new Date().toISOString() })
  }

  function removeProject(id) {
    const i = projects.value.findIndex(p => p.id === id)
    if (i !== -1) projects.value.splice(i, 1)
  }

  // --- Members ---------------------------------------------------------------
  // project.members = [{ id, userId, role }]. Capability flags come from the
  // role preset (see config/roles.js); they are not stored per member.

  function addMember(projectId, { userId, role } = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p || !userId || !role) return null
    if (!p.members) p.members = []
    if (p.members.some(m => m.userId === userId)) return null // one membership per user
    const id = localId('mbr')
    p.members.push({ id, userId, role })
    p.updatedAt = new Date().toISOString()
    return id
  }

  function removeMember(projectId, memberId) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p?.members) return
    const i = p.members.findIndex(m => m.id === memberId)
    if (i !== -1) p.members.splice(i, 1)
    p.updatedAt = new Date().toISOString()
  }

  // --- Resources -------------------------------------------------------------
  // project.resources is a flat list of { id, kind, name }. Linked into Groups
  // (Technical Architecture step) by id.

  const localId = prefix =>
    `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`

  function addResource(projectId, { kind, name, connector } = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p || !kind) return null
    if (!p.resources) p.resources = []
    const id = localId(kind)
    p.resources.push({
      id,
      kind,
      name: (name || '').trim() || 'Untitled',
      links: [],
      ...(connector ? { connector } : {}),
    })
    p.updatedAt = new Date().toISOString()
    return id
  }

  // Remove a resource and scrub references to it (links, fields, group
  // membership, role bindings).
  function removeResource(projectId, resourceId) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p?.resources) return
    const i = p.resources.findIndex(r => r.id === resourceId)
    if (i === -1) return
    p.resources.splice(i, 1)
    for (const r of p.resources) {
      if (r.links) r.links = r.links.filter(l => l.targetId !== resourceId)
      for (const f of r.fields || []) if (f.targetModuleId === resourceId) f.targetModuleId = null
    }
    for (const g of p.groups || []) {
      g.roleIds = (g.roleIds || []).filter(id => id !== resourceId)
      g.resourceIds = (g.resourceIds || []).filter(id => id !== resourceId)
    }
    if (p.roleBindings) {
      p.roleBindings = p.roleBindings.filter(b => b.roleId !== resourceId && b.resourceId !== resourceId)
    }
    p.updatedAt = new Date().toISOString()
  }

  function updateResource(projectId, resourceId, patch = {}) {
    const p = projects.value.find(x => x.id === projectId)
    const r = p?.resources?.find(x => x.id === resourceId)
    if (!r) return
    if (typeof patch.name === 'string') r.name = patch.name
    if (typeof patch.sensitivity === 'string') r.sensitivity = patch.sensitivity
    p.updatedAt = new Date().toISOString()
  }

  // --- Module fields ---------------------------------------------------------
  // A module's `fields` are planning-level: { id, name, type, required,
  // targetModuleId? }. A `record` field references another module (relationship).

  function moduleOf(projectId, moduleId) {
    const p = projects.value.find(x => x.id === projectId)
    return p?.resources?.find(r => r.id === moduleId && r.kind === 'module') || null
  }

  function addField(projectId, moduleId, field = {}) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return null
    if (!m.fields) m.fields = []
    const id = localId('fld')
    m.fields.push({
      id,
      name: field.name || '',
      type: field.type || 'text',
      required: !!field.required,
      targetModuleId: field.targetModuleId || null,
    })
    return id
  }

  function updateField(projectId, moduleId, fieldId, patch = {}) {
    const f = moduleOf(projectId, moduleId)?.fields?.find(x => x.id === fieldId)
    if (!f) return
    Object.assign(f, patch)
    // A non-record type can't carry a target.
    if (f.type !== 'Record') f.targetModuleId = null
  }

  function removeField(projectId, moduleId, fieldId) {
    const m = moduleOf(projectId, moduleId)
    if (!m?.fields) return
    const i = m.fields.findIndex(f => f.id === fieldId)
    if (i !== -1) m.fields.splice(i, 1)
  }

  // Replace a module's whole field list (used to apply a dialog's staged draft).
  function setFields(projectId, moduleId, fields = []) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return
    m.fields = fields.map(f => ({
      id: f.id || localId('fld'),
      name: f.name || '',
      type: f.type || 'text',
      required: !!f.required,
      targetModuleId: f.type === 'Record' ? f.targetModuleId || null : null,
      sensitivity: f.sensitivity || null,
    }))
  }

  // --- Resource links (undirected, untyped) ---------------------------------
  // Stored once per pair on the source resource's `links`. Helpers in
  // utils/links.js read them in both directions.

  const pairLinked = (p, a, b) => {
    for (const r of p.resources || []) {
      if (r.id !== a && r.id !== b) continue
      const other = r.id === a ? b : a
      if ((r.links || []).some(l => l.targetId === other)) return true
    }
    return false
  }

  function addLink(projectId, resourceId, targetId) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p || !resourceId || !targetId || resourceId === targetId) return
    if (pairLinked(p, resourceId, targetId)) return
    const r = p.resources?.find(x => x.id === resourceId)
    if (!r) return
    if (!r.links) r.links = []
    r.links.push({ id: localId('link'), targetId })
    p.updatedAt = new Date().toISOString()
  }

  function removeLink(projectId, a, b) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    for (const r of p.resources || []) {
      if (r.id !== a && r.id !== b) continue
      const other = r.id === a ? b : a
      if (r.links) r.links = r.links.filter(l => l.targetId !== other)
    }
    p.updatedAt = new Date().toISOString()
  }

  // --- Role permissions (RBAC bindings) -------------------------------------
  // Explicit overrides at project.roleBindings = [{ id, roleId, resourceId,
  // level: 'none' | 'read' | 'write' }]. Effective level falls back to a
  // group-membership baseline (see utils/rbac.js). Write implies read.

  // Restore resources + roleBindings from a snapshot — used by the edit dialog's
  // Cancel to discard live changes (links, permissions, renames).
  function restoreState(projectId, { resources, roleBindings } = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    if (resources) p.resources = resources
    if (roleBindings) p.roleBindings = roleBindings
    p.updatedAt = new Date().toISOString()
  }

  function setBinding(projectId, roleId, resourceId, level) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p || !roleId || !resourceId) return
    if (!p.roleBindings) p.roleBindings = []
    const existing = p.roleBindings.find(b => b.roleId === roleId && b.resourceId === resourceId)
    if (existing) existing.level = level
    else p.roleBindings.push({ id: localId('rb'), roleId, resourceId, level })
    p.updatedAt = new Date().toISOString()
  }

  // --- Groups (Technical Architecture) --------------------------------------
  // Formerly "stories". A group bundles roles + resources that work together.
  // project.groups is a list of { id, name, description, roleIds, resourceIds }.

  function addGroup(projectId, group = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return null
    if (!p.groups) p.groups = []
    const id = localId('group')
    p.groups.push({
      id,
      name: (group.name || '').trim() || 'Untitled',
      description: (group.description || '').trim(),
      roleIds: group.roleIds || [],
      resourceIds: group.resourceIds || [],
    })
    p.updatedAt = new Date().toISOString()
    return id
  }

  function updateGroup(projectId, groupId, patch = {}) {
    const p = projects.value.find(x => x.id === projectId)
    const g = p?.groups?.find(x => x.id === groupId)
    if (!g) return
    Object.assign(g, patch)
    p.updatedAt = new Date().toISOString()
  }

  function removeGroup(projectId, groupId) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p?.groups) return
    const i = p.groups.findIndex(g => g.id === groupId)
    if (i !== -1) p.groups.splice(i, 1)
    p.updatedAt = new Date().toISOString()
  }

  // --- Resource Management ---------------------------------------------------
  // project.resourceManagement = { ai, infra, connections[] }. The connections
  // list is the whitelist/catalogue the later Connections step picks from.
  // project.llmCatalog holds provider/model entries added from this project.

  function ensureResourceManagement(p) {
    if (!p.resourceManagement) {
      p.resourceManagement = { ai: {}, infra: {}, connections: [] }
    }
    return p.resourceManagement
  }

  // Patch a sub-section ('ai' | 'infra') of the resource management form.
  function updateResourceManagement(projectId, section, patch = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    const rm = ensureResourceManagement(p)
    rm[section] = { ...(rm[section] || {}), ...patch }
    p.updatedAt = new Date().toISOString()
  }

  function addPermittedConnection(projectId, conn = {}) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return null
    const rm = ensureResourceManagement(p)
    const id = localId('pconn')
    rm.connections.push({
      id,
      name: conn.name || '',
      connector: conn.connector || null,
      actionIfUnavailable: conn.actionIfUnavailable || 'Deactivate',
      replacement: conn.replacement || '',
      type: conn.type || null,
      description: conn.description || '',
      isAiSystem: conn.isAiSystem || 'No',
    })
    p.updatedAt = new Date().toISOString()
    return id
  }

  function updatePermittedConnection(projectId, connId, patch = {}) {
    const p = projects.value.find(x => x.id === projectId)
    const c = p?.resourceManagement?.connections?.find(x => x.id === connId)
    if (!c) return
    Object.assign(c, patch)
    p.updatedAt = new Date().toISOString()
  }

  function removePermittedConnection(projectId, connId) {
    const p = projects.value.find(x => x.id === projectId)
    const list = p?.resourceManagement?.connections
    if (!list) return
    const i = list.findIndex(c => c.id === connId)
    if (i !== -1) list.splice(i, 1)
    p.updatedAt = new Date().toISOString()
  }

  // --- Per-step governance ---------------------------------------------------
  // Each step's state lives at project.governance[stepKey] =
  //   { values, status: draft|submitted|approved|changes-requested, reviewNote }.

  function ensureStep(project, stepKey) {
    if (!project.governance) project.governance = {}
    if (!project.governance[stepKey]) {
      project.governance[stepKey] = { values: {}, status: 'draft', reviewNote: '' }
    }
    return project.governance[stepKey]
  }

  // Submit a whole gate section: every editable step (draft/changes-requested)
  // in the section is sent for approval and locked.
  function submitSection(projectId, stepKeys = []) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    for (const key of stepKeys) {
      const entry = ensureStep(p, key)
      if (entry.status === 'draft' || entry.status === 'changes-requested') {
        entry.status = 'submitted'
        entry.reviewNote = ''
      }
    }
    p.updatedAt = new Date().toISOString()
  }

  function saveStepForm(projectId, stepKey, values) {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    const entry = ensureStep(p, stepKey)
    entry.values = { ...values }
    p.updatedAt = new Date().toISOString()
  }

  // State machine: submit (draft|changes-requested → submitted), approve
  // (submitted → approved), request-changes (submitted → changes-requested),
  // reopen (approved → draft). No-op for invalid transitions.
  function transitionStep(projectId, stepKey, action, note = '') {
    const p = projects.value.find(x => x.id === projectId)
    if (!p) return
    const entry = ensureStep(p, stepKey)
    const s = entry.status
    if (action === 'submit' && (s === 'draft' || s === 'changes-requested')) {
      entry.status = 'submitted'
      entry.reviewNote = ''
    } else if (action === 'approve' && s === 'submitted') {
      entry.status = 'approved'
      entry.reviewNote = ''
    } else if (action === 'request-changes' && s === 'submitted') {
      entry.status = 'changes-requested'
      entry.reviewNote = note
    } else if (action === 'reopen' && s === 'approved') {
      entry.status = 'draft'
      entry.reviewNote = note
    } else if (action === 'recall' && s === 'submitted') {
      // Developer withdraws a pending submission back to draft.
      entry.status = 'draft'
      entry.reviewNote = ''
    } else {
      return
    }
    p.updatedAt = new Date().toISOString()
  }

  return {
    projects,
    findById,
    currentVersion,
    create,
    updateProject,
    removeProject,
    addMember,
    removeMember,
    addResource,
    removeResource,
    updateResource,
    addField,
    updateField,
    removeField,
    setFields,
    addLink,
    removeLink,
    setBinding,
    restoreState,
    addGroup,
    updateGroup,
    removeGroup,
    updateResourceManagement,
    addPermittedConnection,
    updatePermittedConnection,
    removePermittedConnection,
    saveStepForm,
    submitSection,
    transitionStep,
  }
})
