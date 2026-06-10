import { sections, stepsForTab } from '@/sections/project-poc/config/pipeline'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// API-backed projects store. Projects, members, governance, groups and modules
// persist on the backend; the remaining planning-level concepts (non-module
// resources, manual links, role bindings) are session-local until they gain
// backend modeling.
//
// FE project shape mirrors the wizard's needs; `unmarshalProject` maps the
// backend payload (config/meta/governance) onto it and preserves the
// session-local parts across refreshes of the same object.
export const useProjectsStore = defineStore('projects', () => {
  const $SystemAPI = inject('$SystemAPI')
  const $ComposeAPI = inject('$ComposeAPI')

  const projects = ref([])
  const loaded = ref(false)
  let loading = null

  const findById = computed(() => id => projects.value.find(p => p.id === String(id)))

  // Latest published version, or null. Publish flow isn't backend-modeled yet,
  // so this stays empty until it lands.
  const currentVersion = p => (p?.versions?.length ? p.versions[p.versions.length - 1] : null)

  const slugify = name =>
    (name || 'project')
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '') || 'project'

  const localId = prefix =>
    `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`

  // --- payload mapping --------------------------------------------------------

  // Approved gates, derived from governance state against the gated pipeline.
  function countGatesApproved(governance = {}, friaRequired = false) {
    return sections(stepsForTab('governance', { friaRequired }))
      .filter(sec => sec.gateKey)
      .filter(sec => sec.steps.every(s => governance[s.key]?.status === 'approved')).length
  }

  function unmarshalProject(raw, prev = {}) {
    const cfg = raw.config || {}
    const meta = raw.meta || {}
    const governance = raw.governance || {}
    const friaRequired = !!cfg.friaRequired

    return {
      id: String(raw.projectID),
      handle: raw.handle,
      name: meta.short || raw.handle,
      description: meta.description || '',
      mode: cfg.mode || 'free',
      status: raw.status || 'draft',
      namespaceID: cfg.namespaceID ? String(cfg.namespaceID) : null,
      deployerCategories: cfg.deployerCategories || {},
      friaRequired,
      resourceManagement: {
        ai: {},
        infra: {},
        connections: [],
        ...(cfg.resourceManagement || {}),
      },
      governance,
      gatesApproved: countGatesApproved(governance, friaRequired),
      createdBy: String(raw.createdBy || ''),
      createdAt: raw.createdAt,
      updatedAt: raw.updatedAt || raw.createdAt,
      canGrant: !!raw.canGrant,
      canUpdateProject: !!raw.canUpdateProject,
      canDeleteProject: !!raw.canDeleteProject,
      canManageMembers: !!raw.canManageMembers,

      // Loaded separately (fetchProject); preserved across re-unmarshals.
      members: prev.members || [],
      groups: prev.groups || [],
      resources: prev.resources || [],

      // Session-local planning state, not backend-modeled yet.
      roleBindings: prev.roleBindings || [],
      versions: prev.versions || [],

      // Backend bookkeeping for optimistic-lock updates.
      _raw: { config: cfg, meta, status: raw.status, updatedAt: raw.updatedAt },
    }
  }

  // Merge a fresh backend payload into the cached project (or insert it).
  function absorb(raw) {
    const prev = projects.value.find(p => p.id === String(raw.projectID))
    const next = unmarshalProject(raw, prev || {})
    if (prev) Object.assign(prev, next)
    else projects.value.push(next)
    return prev || next
  }

  const unmarshalMember = m => ({
    id: String(m.projectMemberID),
    userId: String(m.userID),
    role: m.rolePreset,
    capabilities: m.capabilities || {},
  })

  // Group entries carry refs to roles vs other resources with a prefix so the
  // two lists survive the single resourceRef field.
  const roleRef = id => `role:${id}`
  const resRef = id => `res:${id}`
  const unmarshalGroup = g => {
    const refs = (g.entries || []).map(e => e.resourceRef)
    return {
      id: String(g.projectGroupID),
      handle: g.handle,
      name: g.meta?.short || g.handle,
      description: g.meta?.description || '',
      roleIds: refs.filter(r => r.startsWith('role:')).map(r => r.slice(5)),
      resourceIds: refs.filter(r => r.startsWith('res:')).map(r => r.slice(4)),
      _updatedAt: g.updatedAt,
    }
  }

  // --- module (compose) mapping ------------------------------------------------

  const fieldName = label =>
    (label || 'field')
      .trim()
      .replace(/[^a-zA-Z0-9_]+/g, '_')
      .replace(/^([0-9])/, 'f$1')
      .replace(/(^_+|_+$)/g, '') || 'field'

  // Sensitivity levels: FE works with handles, compose stores level IDs.
  const sensitivityLevels = ref([]) // [{ id, handle }]
  let sensitivityLoading = null
  async function loadSensitivityLevels() {
    if (sensitivityLevels.value.length) return
    if (sensitivityLoading) return sensitivityLoading
    sensitivityLoading = $SystemAPI
      .dalSensitivityLevelList({})
      .then(({ set = [] } = {}) => {
        sensitivityLevels.value = set.map(l => ({
          id: String(l.sensitivityLevelID),
          handle: l.handle,
        }))
      })
      .catch(err => console.error('Failed to load sensitivity levels', err))
      .finally(() => {
        sensitivityLoading = null
      })
    return sensitivityLoading
  }
  const sensitivityID = handle =>
    sensitivityLevels.value.find(l => l.handle === handle)?.id || undefined
  const sensitivityHandle = id =>
    sensitivityLevels.value.find(l => l.id === String(id))?.handle || null

  const unmarshalModule = m => ({
    id: String(m.moduleID),
    kind: 'module',
    name: m.name || m.handle,
    sensitivity: sensitivityHandle(m.config?.privacy?.sensitivityLevelID),
    fields: (m.fields || []).map(f => ({
      id: String(f.fieldID),
      name: f.label || f.name,
      type: f.kind,
      required: !!f.isRequired,
      targetModuleId: f.kind === 'Record' ? String(f.options?.moduleID || '') || null : null,
      sensitivity: sensitivityHandle(f.config?.privacy?.sensitivityLevelID),
    })),
    links: [],
    _raw: m,
  })

  const marshalFields = (fields = []) =>
    fields.map((f, i) => ({
      name: fieldName(f.name),
      label: f.name || '',
      kind: f.type || 'String',
      place: i,
      isRequired: !!f.required,
      options: f.type === 'Record' && f.targetModuleId ? { moduleID: f.targetModuleId } : {},
      config: f.sensitivity
        ? { privacy: { sensitivityLevelID: sensitivityID(f.sensitivity) } }
        : {},
    }))

  async function pushModule(p, mod) {
    const raw = mod._raw
    const updated = await $ComposeAPI.moduleUpdate({
      namespaceID: p.namespaceID,
      moduleID: mod.id,
      name: mod.name,
      handle: raw.handle,
      fields: marshalFields(mod.fields),
      meta: raw.meta || {},
      config: {
        ...(raw.config || {}),
        privacy: {
          ...(raw.config?.privacy || {}),
          sensitivityLevelID: mod.sensitivity ? sensitivityID(mod.sensitivity) : undefined,
        },
      },
      updatedAt: raw.updatedAt,
    })
    Object.assign(mod, unmarshalModule(updated))
  }

  // --- loading -----------------------------------------------------------------

  async function load() {
    if (loading) return loading
    loading = $SystemAPI
      .projectList({ limit: 500, sort: 'createdAt DESC' })
      .then(({ set = [] } = {}) => {
        for (const raw of set) absorb(raw)
        loaded.value = true
      })
      .catch(err => console.error('Failed to load projects', err))
      .finally(() => {
        loading = null
      })
    return loading
  }

  // Full fetch for the wizard: project + members + groups + real modules.
  async function fetchProject(id) {
    const raw = await $SystemAPI.projectRead({ projectID: id })
    const p = absorb(raw)

    await loadSensitivityLevels()

    const [members, groups, modules] = await Promise.all([
      $SystemAPI.projectListMembers({ projectID: p.id }).catch(() => ({ set: [] })),
      $SystemAPI.projectGroupList({ projectID: p.id }).catch(() => ({ set: [] })),
      p.namespaceID
        ? $ComposeAPI
            .moduleList({ namespaceID: p.namespaceID, limit: 500 })
            .catch(() => ({ set: [] }))
        : { set: [] },
    ])

    p.members = (members.set || []).map(unmarshalMember)
    p.groups = (groups.set || []).map(unmarshalGroup)

    // Real modules replace cached module entries; session-local planning
    // resources of other kinds are kept.
    const planning = (p.resources || []).filter(r => r.kind !== 'module')
    p.resources = [...(modules.set || []).map(unmarshalModule), ...planning]

    return p
  }

  // --- project CRUD --------------------------------------------------------------

  // Create a draft project. The backend creates the compose namespace and adds
  // the creator as a developer; `deployer` carries the AI Act deployer answers.
  async function create({ name, description = '', mode = 'free', deployer = {} } = {}) {
    const raw = await $SystemAPI.projectCreate({
      handle: slugify(name),
      status: 'draft',
      config: {
        mode,
        deployerCategories: {
          publicAuthorityAnnex3: !!deployer.publicAuthorityAnnex3,
          privateEssentialServices: !!deployer.privateEssentialServices,
          insuranceBanking: !!deployer.insuranceBanking,
        },
      },
      meta: { short: (name || '').trim() || 'Untitled project', description: description.trim() },
    })
    return fetchProject(raw.projectID)
  }

  // Push the cached project state (name/description/status/config) to the API.
  async function pushProject(p) {
    const raw = await $SystemAPI.projectUpdate({
      projectID: p.id,
      handle: p.handle,
      status: p.status,
      config: { ...p._raw.config, resourceManagement: p.resourceManagement },
      meta: { ...p._raw.meta, short: p.name, description: p.description },
      updatedAt: p._raw.updatedAt,
    })
    return absorb(raw)
  }

  // Patch a project. `mode` is immutable (also enforced server-side).
  async function updateProject(id, patch = {}) {
    const p = findById.value(id)
    if (!p) return
    const { mode: _ignored, ...rest } = patch
    Object.assign(p, rest)
    return pushProject(p)
  }

  async function removeProject(id) {
    await $SystemAPI.projectDelete({ projectID: id })
    const i = projects.value.findIndex(p => p.id === String(id))
    if (i !== -1) projects.value.splice(i, 1)
  }

  // --- members -------------------------------------------------------------------

  async function addMember(projectId, { userId, role } = {}) {
    const p = findById.value(projectId)
    if (!p || !userId || !role) return null
    if (p.members.some(m => m.userId === String(userId))) return null
    const m = await $SystemAPI.projectAddMember({
      projectID: p.id,
      userID: userId,
      rolePreset: role,
    })
    const member = unmarshalMember(m)
    p.members.push(member)
    return member.id
  }

  async function removeMember(projectId, memberId) {
    const p = findById.value(projectId)
    const m = p?.members?.find(x => x.id === memberId)
    if (!m) return
    await $SystemAPI.projectRemoveMember({ projectID: p.id, userID: m.userId })
    p.members = p.members.filter(x => x.id !== memberId)
  }

  // --- resources -------------------------------------------------------------------
  // Modules are real compose modules in the project namespace. Other kinds
  // (pages, connections, automations, agents, chatbots, roles) are still
  // session-local planning entries until they gain project scoping.

  async function addResource(projectId, { kind, name, connector } = {}) {
    const p = findById.value(projectId)
    if (!p || !kind) return null

    if (kind === 'module' && p.namespaceID) {
      const raw = await $ComposeAPI.moduleCreate({
        namespaceID: p.namespaceID,
        name: (name || '').trim() || 'Untitled',
        handle: slugify(name) + '-' + Date.now().toString(36),
        fields: [],
        meta: {},
      })
      const mod = unmarshalModule(raw)
      p.resources.push(mod)
      return mod.id
    }

    const id = localId(kind)
    p.resources.push({
      id,
      kind,
      name: (name || '').trim() || 'Untitled',
      links: [],
      ...(connector ? { connector } : {}),
    })
    return id
  }

  // Remove a resource and scrub references to it (links, fields, group
  // membership, role bindings).
  async function removeResource(projectId, resourceId) {
    const p = findById.value(projectId)
    if (!p?.resources) return
    const r = p.resources.find(x => x.id === resourceId)
    if (!r) return

    if (r.kind === 'module' && r._raw) {
      await $ComposeAPI.moduleDelete({ namespaceID: p.namespaceID, moduleID: r.id })
    }

    p.resources = p.resources.filter(x => x.id !== resourceId)
    for (const o of p.resources) {
      if (o.links) o.links = o.links.filter(l => l.targetId !== resourceId)
      for (const f of o.fields || []) if (f.targetModuleId === resourceId) f.targetModuleId = null
    }
    for (const g of p.groups || []) {
      const patch = {}
      if ((g.roleIds || []).includes(resourceId)) {
        patch.roleIds = g.roleIds.filter(id => id !== resourceId)
      }
      if ((g.resourceIds || []).includes(resourceId)) {
        patch.resourceIds = g.resourceIds.filter(id => id !== resourceId)
      }
      if (Object.keys(patch).length) await updateGroup(p.id, g.id, patch)
    }
    if (p.roleBindings) {
      p.roleBindings = p.roleBindings.filter(
        b => b.roleId !== resourceId && b.resourceId !== resourceId,
      )
    }
  }

  async function updateResource(projectId, resourceId, patch = {}) {
    const p = findById.value(projectId)
    const r = p?.resources?.find(x => x.id === resourceId)
    if (!r) return
    if (typeof patch.name === 'string') r.name = patch.name
    if (typeof patch.sensitivity === 'string') r.sensitivity = patch.sensitivity
    if (r.kind === 'module' && r._raw) await pushModule(p, r)
  }

  // --- module fields ---------------------------------------------------------------

  function moduleOf(projectId, moduleId) {
    const p = findById.value(projectId)
    return p?.resources?.find(r => r.id === moduleId && r.kind === 'module') || null
  }

  async function addField(projectId, moduleId, field = {}) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return null
    if (!m.fields) m.fields = []
    m.fields.push({
      id: localId('fld'),
      name: field.name || '',
      type: field.type || 'String',
      required: !!field.required,
      targetModuleId: field.targetModuleId || null,
    })
    if (m._raw) await pushModule(findById.value(projectId), m)
    return m.fields[m.fields.length - 1].id
  }

  async function updateField(projectId, moduleId, fieldId, patch = {}) {
    const m = moduleOf(projectId, moduleId)
    const f = m?.fields?.find(x => x.id === fieldId)
    if (!f) return
    Object.assign(f, patch)
    // A non-record type can't carry a target.
    if (f.type !== 'Record') f.targetModuleId = null
    if (m._raw) await pushModule(findById.value(projectId), m)
  }

  async function removeField(projectId, moduleId, fieldId) {
    const m = moduleOf(projectId, moduleId)
    if (!m?.fields) return
    m.fields = m.fields.filter(f => f.id !== fieldId)
    if (m._raw) await pushModule(findById.value(projectId), m)
  }

  // Replace a module's whole field list (used to apply a dialog's staged draft).
  async function setFields(projectId, moduleId, fields = []) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return
    m.fields = fields.map(f => ({
      id: f.id || localId('fld'),
      name: f.name || '',
      type: f.type || 'String',
      required: !!f.required,
      targetModuleId: f.type === 'Record' ? f.targetModuleId || null : null,
      sensitivity: f.sensitivity || null,
    }))
    if (m._raw) await pushModule(findById.value(projectId), m)
  }

  // --- resource links (session-local, undirected, untyped) ---------------------
  // Manual planning links between resources. The persisted relationship graph
  // comes from the backend (/projects/{id}/graph); these only augment the
  // session view until links gain backend modeling.

  const pairLinked = (p, a, b) => {
    for (const r of p.resources || []) {
      if (r.id !== a && r.id !== b) continue
      const other = r.id === a ? b : a
      if ((r.links || []).some(l => l.targetId === other)) return true
    }
    return false
  }

  function addLink(projectId, resourceId, targetId) {
    const p = findById.value(projectId)
    if (!p || !resourceId || !targetId || resourceId === targetId) return
    if (pairLinked(p, resourceId, targetId)) return
    const r = p.resources?.find(x => x.id === resourceId)
    if (!r) return
    if (!r.links) r.links = []
    r.links.push({ id: localId('link'), targetId })
  }

  function removeLink(projectId, a, b) {
    const p = findById.value(projectId)
    if (!p) return
    for (const r of p.resources || []) {
      if (r.id !== a && r.id !== b) continue
      const other = r.id === a ? b : a
      if (r.links) r.links = r.links.filter(l => l.targetId !== other)
    }
  }

  // --- role permissions (session-local RBAC overrides) --------------------------
  // Explicit overrides at project.roleBindings = [{ id, roleId, resourceId,
  // level: 'none' | 'read' | 'write' }]. Effective level falls back to a
  // group-membership baseline (see utils/rbac.js). Not backend-modeled yet.

  function setBinding(projectId, roleId, resourceId, level) {
    const p = findById.value(projectId)
    if (!p || !roleId || !resourceId) return
    if (!p.roleBindings) p.roleBindings = []
    const existing = p.roleBindings.find(b => b.roleId === roleId && b.resourceId === resourceId)
    if (existing) existing.level = level
    else p.roleBindings.push({ id: localId('rb'), roleId, resourceId, level })
  }

  // Restore session-local state from a snapshot — used by the edit dialog's
  // Cancel to discard live changes (links, permissions). Renames and field
  // edits on real modules persist through the API and are not restored.
  function restoreState(projectId, { resources, roleBindings } = {}) {
    const p = findById.value(projectId)
    if (!p) return
    if (resources) p.resources = resources
    if (roleBindings) p.roleBindings = roleBindings
  }

  // --- groups (Technical Architecture) ------------------------------------------

  async function syncGroupEntries(p, group, roleIds, resourceIds) {
    const want = new Set([...(roleIds || []).map(roleRef), ...(resourceIds || []).map(resRef)])
    const have = new Set([
      ...(group.roleIds || []).map(roleRef),
      ...(group.resourceIds || []).map(resRef),
    ])
    for (const ref of want) {
      if (!have.has(ref)) {
        await $SystemAPI.projectGroupEntryAdd({
          projectID: p.id,
          projectGroupID: group.id,
          resourceRef: ref,
        })
      }
    }
    for (const ref of have) {
      if (!want.has(ref)) {
        await $SystemAPI.projectGroupEntryRemove({
          projectID: p.id,
          projectGroupID: group.id,
          resourceRef: ref,
        })
      }
    }
  }

  async function addGroup(projectId, group = {}) {
    const p = findById.value(projectId)
    if (!p) return null
    const name = (group.name || '').trim() || 'Untitled'
    const raw = await $SystemAPI.projectGroupCreate({
      projectID: p.id,
      handle: slugify(name) + '-' + Date.now().toString(36),
      name,
      description: (group.description || '').trim(),
    })
    const g = unmarshalGroup(raw)
    p.groups.push(g)
    await syncGroupEntries(p, g, group.roleIds || [], group.resourceIds || [])
    g.roleIds = group.roleIds || []
    g.resourceIds = group.resourceIds || []
    return g.id
  }

  async function updateGroup(projectId, groupId, patch = {}) {
    const p = findById.value(projectId)
    const g = p?.groups?.find(x => x.id === groupId)
    if (!g) return

    if (patch.roleIds || patch.resourceIds) {
      await syncGroupEntries(p, g, patch.roleIds || g.roleIds, patch.resourceIds || g.resourceIds)
    }

    if (patch.name !== undefined || patch.description !== undefined) {
      const raw = await $SystemAPI.projectGroupUpdate({
        projectID: p.id,
        projectGroupID: g.id,
        handle: g.handle,
        name: patch.name ?? g.name,
        description: patch.description ?? g.description,
        updatedAt: g._updatedAt,
      })
      g._updatedAt = raw.updatedAt
    }

    Object.assign(g, patch)
  }

  async function removeGroup(projectId, groupId) {
    const p = findById.value(projectId)
    if (!p?.groups) return
    await $SystemAPI.projectGroupDelete({ projectID: p.id, projectGroupID: groupId })
    p.groups = p.groups.filter(g => g.id !== groupId)
  }

  // --- resource management --------------------------------------------------------
  // project.resourceManagement = { ai, infra, connections[] } persists in the
  // project config; the connections list is the whitelist/catalogue the later
  // Connections step picks from.

  async function updateResourceManagement(projectId, section, patch = {}) {
    const p = findById.value(projectId)
    if (!p) return
    p.resourceManagement[section] = { ...(p.resourceManagement[section] || {}), ...patch }
    return pushProject(p)
  }

  async function addPermittedConnection(projectId, conn = {}) {
    const p = findById.value(projectId)
    if (!p) return null
    const id = localId('pconn')
    p.resourceManagement.connections.push({
      id,
      name: conn.name || '',
      connector: conn.connector || null,
      actionIfUnavailable: conn.actionIfUnavailable || 'Deactivate',
      replacement: conn.replacement || '',
      type: conn.type || null,
      description: conn.description || '',
      isAiSystem: conn.isAiSystem || 'No',
    })
    await pushProject(p)
    return id
  }

  async function updatePermittedConnection(projectId, connId, patch = {}) {
    const p = findById.value(projectId)
    const c = p?.resourceManagement?.connections?.find(x => x.id === connId)
    if (!c) return
    Object.assign(c, patch)
    return pushProject(p)
  }

  async function removePermittedConnection(projectId, connId) {
    const p = findById.value(projectId)
    const list = p?.resourceManagement?.connections
    if (!list) return
    p.resourceManagement.connections = list.filter(c => c.id !== connId)
    return pushProject(p)
  }

  // --- per-step governance ----------------------------------------------------------
  // Step state lives on the backend (project.governance[stepKey] = { values,
  // status, reviewNote }); transitions are capability-checked server-side.

  async function saveStepForm(projectId, stepKey, values) {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectGovernanceSave({
      projectID: p.id,
      stepKey,
      values: { ...values },
    })
    absorb(raw)
  }

  // State machine (server-side): submit (draft|changes-requested → submitted),
  // approve (submitted → approved), request-changes (submitted →
  // changes-requested), reopen (approved → draft), recall (submitted → draft).
  async function transitionStep(projectId, stepKey, action, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectGovernanceTransition({
      projectID: p.id,
      stepKey,
      action,
      note,
    })
    absorb(raw)
  }

  // Submit a whole gate section: every editable step (draft/changes-requested)
  // in the section is sent for approval and locked.
  async function submitSection(projectId, stepKeys = []) {
    const p = findById.value(projectId)
    if (!p) return
    for (const key of stepKeys) {
      const status = p.governance?.[key]?.status || 'draft'
      if (status === 'draft' || status === 'changes-requested') {
        await transitionStep(projectId, key, 'submit')
      }
    }
  }

  // --- graph ------------------------------------------------------------------------

  // Backend-composed resource graph (real resources + derived relations).
  async function graph(projectId) {
    const p = findById.value(projectId)
    if (!p) return { nodes: [], edges: [] }
    try {
      const g = await $SystemAPI.projectGraph({ projectID: p.id })
      return { nodes: g?.nodes || [], edges: g?.edges || [] }
    } catch (err) {
      console.error('Failed to load project graph', err)
      return { nodes: [], edges: [] }
    }
  }

  return {
    projects,
    loaded,
    load,
    fetchProject,
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
    graph,
  }
})
