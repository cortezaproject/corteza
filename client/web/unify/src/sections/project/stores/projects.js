import { sections, stepsForTab } from '@/sections/project/config/pipeline'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// API-backed projects store. Everything here persists on the backend; the
// store grows alongside the pipeline, one verified step at a time. Current
// surface: project CRUD, members (read-only, for capability resolution), the
// per-step governance workflow, and modules (real compose modules in the
// project's namespace).
export const useProjectsStore = defineStore('projects', () => {
  const $SystemAPI = inject('$SystemAPI')
  const $ComposeAPI = inject('$ComposeAPI')

  const projects = ref([])
  const loaded = ref(false)
  let loading = null

  const findById = computed(() => id => projects.value.find(p => p.id === String(id)))

  // --- payload mapping --------------------------------------------------------

  // Approved gates, derived from governance state against the gated pipeline.
  function countGatesApproved(governance = {}) {
    return sections(stepsForTab('governance'))
      .filter(sec => sec.gateKey)
      .filter(sec => sec.steps.every(s => governance[s.key]?.status === 'approved')).length
  }

  function unmarshalProject(raw, prev = {}) {
    const cfg = raw.config || {}
    const meta = raw.meta || {}
    const governance = raw.governance || {}

    return {
      id: String(raw.projectID),
      handle: raw.handle,
      name: meta.short || raw.handle,
      description: meta.description || '',
      mode: cfg.mode || 'free',
      status: raw.status || 'draft',
      namespaceID: cfg.namespaceID ? String(cfg.namespaceID) : null,
      governance,
      gatesApproved: countGatesApproved(governance),
      createdBy: String(raw.createdBy || ''),
      createdAt: raw.createdAt,
      updatedAt: raw.updatedAt || raw.createdAt,
      canGrant: !!raw.canGrant,
      canUpdateProject: !!raw.canUpdateProject,
      canDeleteProject: !!raw.canDeleteProject,
      canManageMembers: !!raw.canManageMembers,

      // Loaded separately (fetchProject); preserved across re-unmarshals.
      members: prev.members || [],
      resources: prev.resources || [],

      // Backend bookkeeping for optimistic-lock updates. Config is carried
      // whole so updates round-trip fields the UI doesn't surface yet
      // (mode, namespaceID, deployer answers, …).
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

  // --- module (compose) mapping ------------------------------------------------

  const localId = prefix =>
    `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`

  // Compose field names must be handles; the human label carries the FE name.
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

  // Full fetch for the wizard: project + members (capability resolution) +
  // the real compose modules from the project namespace.
  async function fetchProject(id) {
    const raw = await $SystemAPI.projectRead({ projectID: id })
    const p = absorb(raw)

    await loadSensitivityLevels()

    const [members, modules] = await Promise.all([
      $SystemAPI.projectListMembers({ projectID: p.id }).catch(() => ({ set: [] })),
      p.namespaceID
        ? $ComposeAPI
            .moduleList({ namespaceID: p.namespaceID, limit: 500 })
            .catch(() => ({ set: [] }))
        : { set: [] },
    ])

    p.members = (members.set || []).map(unmarshalMember)
    p.resources = (modules.set || []).map(unmarshalModule)

    return p
  }

  // --- project CRUD --------------------------------------------------------------

  // Create a draft project. The backend generates the internal handle, creates
  // the compose namespace and adds the creator as a developer; `deployer`
  // carries the AI Act deployer answers.
  async function create({ name, description = '', mode = 'free', deployer = {} } = {}) {
    const raw = await $SystemAPI.projectCreate({
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

  // Push the cached project state (name/description/status) to the API.
  async function pushProject(p) {
    const raw = await $SystemAPI.projectUpdate({
      projectID: p.id,
      handle: p.handle,
      status: p.status,
      config: p._raw.config,
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

  // --- resources (modules) -------------------------------------------------------
  // All wizard resources are real compose modules in the project namespace.

  async function addResource(projectId, { kind, name } = {}) {
    const p = findById.value(projectId)
    if (!p || kind !== 'module' || !p.namespaceID) return null

    const raw = await $ComposeAPI.moduleCreate({
      namespaceID: p.namespaceID,
      name: (name || '').trim() || 'Untitled',
      handle: fieldName(name).toLowerCase() + '_' + Date.now().toString(36),
      fields: [],
      meta: {},
    })
    const mod = unmarshalModule(raw)
    p.resources.push(mod)
    return mod.id
  }

  // Remove a module and scrub record-field references to it.
  async function removeResource(projectId, resourceId) {
    const p = findById.value(projectId)
    const r = p?.resources?.find(x => x.id === resourceId)
    if (!r) return

    await $ComposeAPI.moduleDelete({ namespaceID: p.namespaceID, moduleID: r.id })

    p.resources = p.resources.filter(x => x.id !== resourceId)
    for (const o of p.resources) {
      for (const f of o.fields || []) if (f.targetModuleId === resourceId) f.targetModuleId = null
    }
  }

  async function updateResource(projectId, resourceId, patch = {}) {
    const p = findById.value(projectId)
    const r = p?.resources?.find(x => x.id === resourceId)
    if (!r) return
    if (typeof patch.name === 'string') r.name = patch.name
    if (typeof patch.sensitivity === 'string') r.sensitivity = patch.sensitivity
    await pushModule(p, r)
  }

  // --- module fields ---------------------------------------------------------------

  function moduleOf(projectId, moduleId) {
    const p = findById.value(projectId)
    return p?.resources?.find(r => r.id === moduleId && r.kind === 'module') || null
  }

  async function updateField(projectId, moduleId, fieldId, patch = {}) {
    const m = moduleOf(projectId, moduleId)
    const f = m?.fields?.find(x => x.id === fieldId)
    if (!f) return
    Object.assign(f, patch)
    // A non-record type can't carry a target.
    if (f.type !== 'Record') f.targetModuleId = null
    await pushModule(findById.value(projectId), m)
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
    await pushModule(findById.value(projectId), m)
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

  return {
    projects,
    loaded,
    load,
    fetchProject,
    findById,
    create,
    updateProject,
    removeProject,
    addResource,
    removeResource,
    updateResource,
    updateField,
    setFields,
    saveStepForm,
    submitSection,
    transitionStep,
  }
})
