import { ACCESS_KINDS, NODE_LAYER_KINDS } from '@/sections/project/config/kinds'
import { sections, stepsForTab } from '@/sections/project/config/pipeline'
import { SENSITIVITY_LEVELS } from '@/sections/project/config/sensitivity'
import { fieldName } from '@/sections/project/utils/fields'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// API-backed projects store. Everything here persists on the backend; the
// store grows alongside the pipeline, one verified step at a time. Current
// surface: project CRUD, members (role-preset CRUD + capability resolution),
// the per-step governance workflow, and modules (real compose modules in the
// project's namespace).
export const useProjectsStore = defineStore('projects', () => {
  const $SystemAPI = inject('$SystemAPI')
  const $ComposeAPI = inject('$ComposeAPI')
  const $AutomationAPI = inject('$AutomationAPI')

  const projects = ref([])
  const loaded = ref(false)
  let loading = null

  const findById = computed(() => id => projects.value.find(p => p.id === String(id)))

  // Bumped after every persisting mutation; the resource graph watches it and
  // refetches, so the panel always reflects current state.
  const graphVersion = ref(0)
  const touch = () => {
    graphVersion.value++
  }

  // --- resource graph view (layer selector) --------------------------------------
  // Which kinds the graph currently shows. One global selection shared across
  // projects (it resets on a hard reload — it's deliberately session state, not
  // persisted). Seeded with every node-layer kind so all layers start on; the
  // access overlay (roles/users) is a separate flag, off by default.
  const graphVisibleKinds = ref(new Set(NODE_LAYER_KINDS))
  const graphShowAccess = ref(false)

  const graphKindVisible = computed(() => kind => {
    if (ACCESS_KINDS.includes(kind)) return graphShowAccess.value
    return graphVisibleKinds.value.has(kind)
  })

  function graphToggleKind(kind) {
    // Reassign the Set so the ref's dependents re-run (Set mutation alone won't).
    const next = new Set(graphVisibleKinds.value)
    next.has(kind) ? next.delete(kind) : next.add(kind)
    graphVisibleKinds.value = next
  }

  function graphToggleAccess() {
    graphShowAccess.value = !graphShowAccess.value
  }

  // Resources (compose modules) are deliberately NOT stored on the project
  // object. They live here keyed by projectID and are (re)fetched on demand,
  // filtered by projectID — fetching is the only way to get a project's
  // resources, and a mutation always refetches rather than patching in place.
  const resourcesByProject = ref({})
  const resourcesFor = computed(() => projectId => resourcesByProject.value[String(projectId)] || [])

  // The connection library (catalog + already-configured connections) — the
  // same set the Admin connection screen lists. Loaded once and shared across
  // projects; the Connections step filters it down to the Resource Management
  // whitelist. A project's own connections are kept separately, keyed by id.
  const connectionLibrary = ref([])
  const connectionsByProject = ref({})
  const connectionsFor = computed(() => projectId => connectionsByProject.value[String(projectId)] || [])

  // A project's automations (TAQs) — NgAutomation records stamped with its
  // projectID. Kept separately, keyed by projectID, and (re)fetched on demand
  // like connections.
  const automationsByProject = ref({})
  const automationsFor = computed(() => projectId => automationsByProject.value[String(projectId)] || [])

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

      // Members are loaded separately (fetchProject) and preserved across
      // re-unmarshals. Resources are NOT held here — they live in the
      // store's resourcesByProject cache, fetched by projectID.
      members: prev.members || [],

      // Backend bookkeeping for optimistic-lock updates. Config is carried
      // whole so updates round-trip fields the UI doesn't surface yet
      // (mode, namespaceID, deployer answers, …).
      _raw: { config: cfg, meta, status: raw.status, updatedAt: raw.updatedAt },
    }
  }

  // Merge a fresh backend payload into the cached project (or insert it).
  // Always returns the reactive instance from the array (never the raw object
  // that was pushed) so later mutations on it are seen by watchers.
  function absorb(raw) {
    const prev = projects.value.find(p => p.id === String(raw.projectID))
    const next = unmarshalProject(raw, prev || {})
    if (prev) {
      Object.assign(prev, next)
      return prev
    }
    projects.value.push(next)
    return projects.value[projects.value.length - 1]
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

  // Sensitivity levels: FE works with handles, compose stores level IDs.
  // The standard scheme (config/sensitivity.js) is real DAL sensitivity-level
  // resources; if any are missing we create them on demand so field
  // classification persists (see ensureStandardLevels).
  const sensitivityLevels = ref([]) // [{ id, handle, level }]
  let sensitivityLoading = null
  async function loadSensitivityLevels() {
    if (sensitivityLevels.value.length) return
    if (sensitivityLoading) return sensitivityLoading
    sensitivityLoading = (async () => {
      try {
        const { set = [] } = await $SystemAPI.dalSensitivityLevelList({})
        sensitivityLevels.value = set.map(l => ({
          id: String(l.sensitivityLevelID),
          handle: l.handle,
          level: l.level,
        }))
        await ensureStandardLevels()
      } catch (err) {
        console.error('Failed to load sensitivity levels', err)
      } finally {
        sensitivityLoading = null
      }
    })()
    return sensitivityLoading
  }

  // Create any standard levels that don't exist yet. The backend requires a
  // non-empty meta.name and a unique `level` int (it doubles as the rank), so
  // missing levels get appended with the next free ascending ints, in scheme
  // order. Best-effort and per-level: creating needs the manage grant, so a
  // failure just leaves that level absent rather than breaking the load.
  async function ensureStandardLevels() {
    const have = new Set(sensitivityLevels.value.map(l => l.handle))
    const missing = SENSITIVITY_LEVELS.filter(s => !have.has(s.id))
    if (!missing.length) return
    let nextLevel = sensitivityLevels.value.reduce((m, l) => Math.max(m, l.level || 0), 0) + 1
    for (const s of missing) {
      try {
        const raw = await $SystemAPI.dalSensitivityLevelCreate({
          handle: s.id,
          level: nextLevel++,
          meta: { name: s.label, description: '' },
        })
        sensitivityLevels.value.push({
          id: String(raw.sensitivityLevelID),
          handle: raw.handle,
          level: raw.level,
        })
      } catch (err) {
        console.error(`Failed to create sensitivity level "${s.id}"`, err)
      }
    }
  }
  const sensitivityID = handle =>
    sensitivityLevels.value.find(l => l.handle === handle)?.id || undefined
  const sensitivityHandle = id =>
    sensitivityLevels.value.find(l => l.id === String(id))?.handle || null

  const unmarshalModule = m => ({
    id: String(m.moduleID),
    kind: 'module',
    name: m.name || m.handle,
    description: m.meta?.description || '',
    // Modules carry no sensitivity level — only their fields are classified.
    fields: (m.fields || []).map(f => ({
      id: String(f.fieldID),
      name: f.label || f.name,
      type: f.kind,
      required: !!f.isRequired,
      multi: !!f.isMulti,
      // Universal help text (Compose stores it under options.description).
      helpText: f.options?.description?.edit ?? f.options?.description?.view ?? '',
      // Text: render across multiple lines.
      multiLine: f.kind === 'String' ? !!f.options?.multiLine : false,
      // Number: digits after the decimal point (defaults to 0 = whole numbers).
      precision: f.kind === 'Number' ? (Number.isFinite(f.options?.precision) ? f.options.precision : 0) : null,
      // DateTime: 'date' | 'time' | 'datetime' (derived from the two flags).
      dateMode:
        f.kind === 'DateTime'
          ? f.options?.onlyDate
            ? 'date'
            : f.options?.onlyTime
              ? 'time'
              : 'datetime'
          : 'datetime',
      targetModuleId: f.kind === 'Record' ? String(f.options?.moduleID || '') || null : null,
      labelField: f.kind === 'Record' ? f.options?.labelField || null : null,
      // Select options arrive as [{ value, text }] (or bare strings); kept as
      // { value, text } pairs so the editor can show a value + label table.
      selectOptions:
        f.kind === 'Select'
          ? (f.options?.options || []).map(o =>
              typeof o === 'string'
                ? { value: o, text: o }
                : { value: o.value ?? '', text: o.text ?? o.value ?? '' },
            )
          : [],
      sensitivity: sensitivityHandle(f.config?.privacy?.sensitivityLevelID),
    })),
    _raw: m,
  })

  // Field options for the compose payload: universal help text plus the
  // kind-specific settings the UI surfaces.
  const fieldOptions = f => {
    const opts = {}
    const help = (f.helpText || '').trim()
    if (help) opts.description = { view: help, edit: help }
    if (f.type === 'String' && f.multiLine) opts.multiLine = true
    if (f.type === 'Number' && Number.isFinite(f.precision)) opts.precision = f.precision
    if (f.type === 'DateTime') {
      if (f.dateMode === 'date') opts.onlyDate = true
      else if (f.dateMode === 'time') opts.onlyTime = true
    }
    if (f.type === 'Record' && f.targetModuleId) {
      opts.moduleID = f.targetModuleId
      if (f.labelField) opts.labelField = f.labelField
    }
    if (f.type === 'Select' && f.selectOptions?.length) {
      const options = f.selectOptions
        .filter(o => (o.value ?? '').toString().trim())
        .map(o => ({ value: o.value, text: o.text || o.value }))
      if (options.length) opts.options = options
    }
    return opts
  }

  const marshalFields = (fields = []) =>
    fields.map((f, i) => ({
      name: fieldName(f.name),
      label: f.name || '',
      kind: f.type || 'String',
      place: i,
      isRequired: !!f.required,
      isMulti: !!f.multi,
      options: fieldOptions(f),
      config: f.sensitivity
        ? { privacy: { sensitivityLevelID: sensitivityID(f.sensitivity) } }
        : {},
    }))

  async function pushModule(p, mod) {
    const raw = mod._raw
    await $ComposeAPI.moduleUpdate({
      namespaceID: p.namespaceID,
      moduleID: mod.id,
      name: mod.name,
      handle: raw.handle,
      fields: marshalFields(mod.fields),
      meta: { ...(raw.meta || {}), description: mod.description || '' },
      // Module config is passed through untouched — modules carry no
      // sensitivity level (only fields are classified, via marshalFields).
      config: { ...(raw.config || {}) },
      updatedAt: raw.updatedAt,
    })
    // Resync from the server (by projectID) rather than patching in place.
    await loadResources(p.id)
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

  // Full fetch for the wizard: project + members (capability resolution). The
  // project's resources are loaded separately into the resourcesByProject
  // cache (fetched by projectID), never attached to the project object.
  async function fetchProject(id) {
    const raw = await $SystemAPI.projectRead({ projectID: id })
    const p = absorb(raw)

    await loadSensitivityLevels()

    const [members] = await Promise.all([
      $SystemAPI.projectListMembers({ projectID: p.id }).catch(() => ({ set: [] })),
      loadResources(p.id),
    ])

    p.members = (members.set || []).map(unmarshalMember)

    return p
  }

  // Fetch a project's resources (compose modules) filtered by projectID and
  // cache them. This is the single source of resources; callers read them via
  // resourcesFor(projectId), and every resource mutation calls this to resync.
  async function loadResources(projectId) {
    const p = findById.value(projectId)
    if (!p?.namespaceID) {
      resourcesByProject.value[String(projectId)] = []
      return []
    }
    const { set = [] } = await $ComposeAPI
      .moduleList({ namespaceID: p.namespaceID, projectID: p.id, limit: 500 })
      .catch(() => ({ set: [] }))
    resourcesByProject.value[String(projectId)] = set.map(unmarshalModule)
    touch()
    return resourcesByProject.value[String(projectId)]
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
    touch()
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

  // --- members --------------------------------------------------------------
  // Memberships are keyed by user (one record per user per project); the role
  // is a fixed named preset and capabilities come back derived from it.
  // Mutations are RBAC-checked server-side (project members.manage).

  async function addMember(projectId, { userId, role } = {}) {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectAddMember({
      projectID: p.id,
      userID: userId,
      rolePreset: role,
    })
    p.members.push(unmarshalMember(raw))
    touch()
  }

  async function updateMember(projectId, userId, role) {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectUpdateMember({
      projectID: p.id,
      userID: userId,
      rolePreset: role,
    })
    const m = p.members.find(x => x.userId === String(userId))
    if (m) Object.assign(m, unmarshalMember(raw))
    touch()
  }

  async function removeMember(projectId, userId) {
    const p = findById.value(projectId)
    if (!p) return
    await $SystemAPI.projectRemoveMember({ projectID: p.id, userID: userId })
    p.members = p.members.filter(m => m.userId !== String(userId))
    touch()
  }

  // --- resources (modules) -------------------------------------------------------
  // All wizard resources are real compose modules in the project namespace.

  async function addResource(projectId, { kind, name } = {}) {
    const p = findById.value(projectId)
    if (!p || kind !== 'module' || !p.namespaceID) return null

    const raw = await $ComposeAPI.moduleCreate({
      namespaceID: p.namespaceID,
      name: (name || '').trim() || 'Untitled',
      // Handle is the slugified title (no hash). Duplicate titles collide on
      // the unique handle; the create dialog validates against that first.
      handle: fieldName(name).toLowerCase(),
      fields: [],
      meta: {},
    })
    await loadResources(projectId)
    return String(raw.moduleID)
  }

  async function removeResource(projectId, resourceId) {
    const p = findById.value(projectId)
    const r = moduleOf(projectId, resourceId)
    if (!r) return

    await $ComposeAPI.moduleDelete({ namespaceID: p.namespaceID, moduleID: r.id })
    await loadResources(projectId)
  }

  async function updateResource(projectId, resourceId, patch = {}) {
    const r = moduleOf(projectId, resourceId)
    if (!r) return
    if (typeof patch.name === 'string') r.name = patch.name
    if (typeof patch.description === 'string') r.description = patch.description
    await pushModule(findById.value(projectId), r)
  }

  // --- module fields ---------------------------------------------------------------

  function moduleOf(projectId, moduleId) {
    return (
      resourcesByProject.value[String(projectId)]?.find(
        r => r.id === moduleId && r.kind === 'module',
      ) || null
    )
  }

  // Canonical field shape: keeps only the settings that apply to the field's
  // type, so switching kinds never leaves stale options behind.
  const normalizeField = f => ({
    id: f.id || localId('fld'),
    name: f.name || '',
    type: f.type || 'String',
    required: !!f.required,
    multi: !!f.multi,
    helpText: f.helpText || '',
    multiLine: f.type === 'String' ? !!f.multiLine : false,
    precision: f.type === 'Number' ? (Number.isFinite(f.precision) ? f.precision : 0) : null,
    dateMode: f.type === 'DateTime' ? f.dateMode || 'datetime' : 'datetime',
    targetModuleId: f.type === 'Record' ? f.targetModuleId || null : null,
    labelField: f.type === 'Record' ? f.labelField || null : null,
    selectOptions:
      f.type === 'Select'
        ? (f.selectOptions || []).map(o =>
            typeof o === 'string' ? { value: o, text: o } : { value: o.value ?? '', text: o.text ?? '' },
          )
        : [],
    sensitivity: f.sensitivity || null,
  })

  async function updateField(projectId, moduleId, fieldId, patch = {}) {
    const m = moduleOf(projectId, moduleId)
    const f = m?.fields?.find(x => x.id === fieldId)
    if (!f) return
    Object.assign(f, normalizeField({ ...f, ...patch, id: f.id }))
    await pushModule(findById.value(projectId), m)
  }

  // Append a single field to a module; returns the new field's id (or null if
  // the module is gone). Used by the quick "Add field" flow on the data model.
  async function addField(projectId, moduleId, patch = {}) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return null
    const id = localId('fld')
    const field = normalizeField({ ...patch, id })
    m.fields = [...(m.fields || []), field]
    await pushModule(findById.value(projectId), m)
    return id
  }

  // Remove a single field from a module.
  async function removeField(projectId, moduleId, fieldId) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return
    m.fields = (m.fields || []).filter(f => f.id !== fieldId)
    await pushModule(findById.value(projectId), m)
  }

  // Replace a module's whole field list (used to apply a dialog's staged draft).
  async function setFields(projectId, moduleId, fields = []) {
    const m = moduleOf(projectId, moduleId)
    if (!m) return
    m.fields = fields.map(f => normalizeField({ ...f, id: f.id || localId('fld') }))
    await pushModule(findById.value(projectId), m)
  }

  // --- per-step governance ----------------------------------------------------------
  // Governance is kept in memory ONLY for now — nothing is persisted to the
  // backend. The project object carries the working governance state
  // (project.governance[stepKey] = { values, status, reviewNote }) but it is
  // never saved, so it resets on reload. (Backend governance endpoints still
  // exist; the FE just doesn't call them yet.)

  function ensureGovStep(p, stepKey) {
    if (!p.governance) p.governance = {}
    if (!p.governance[stepKey]) {
      p.governance[stepKey] = { values: {}, status: 'draft', reviewNote: '' }
    }
    return p.governance[stepKey]
  }

  async function saveStepForm(projectId, stepKey, values) {
    const p = findById.value(projectId)
    if (!p) return
    const step = ensureGovStep(p, stepKey)
    step.values = { ...values }
    p.gatesApproved = countGatesApproved(p.governance)
    touch()
  }

  // Local mirror of the (server-side) state machine: submit
  // (draft|changes-requested → submitted), approve (submitted → approved),
  // request-changes (submitted → changes-requested), reopen (approved → draft),
  // recall (submitted → draft).
  const GOV_TRANSITIONS = {
    submit: { from: ['draft', 'changes-requested'], to: 'submitted', clearNote: true },
    approve: { from: ['submitted'], to: 'approved', clearNote: true },
    'request-changes': { from: ['submitted'], to: 'changes-requested' },
    reopen: { from: ['approved'], to: 'draft' },
    recall: { from: ['submitted'], to: 'draft', clearNote: true },
  }

  async function transitionStep(projectId, stepKey, action, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const step = ensureGovStep(p, stepKey)
    const t = GOV_TRANSITIONS[action]
    if (!t || !t.from.includes(step.status)) return
    step.status = t.to
    step.reviewNote = t.clearNote ? '' : note
    p.gatesApproved = countGatesApproved(p.governance)
    touch()
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

  // --- connections ---------------------------------------------------------------
  // The Connections step instantiates connections the project is permitted to
  // use. The catalogue is the shared connection library; the whitelist comes
  // from the Resource Management governance step.

  // Load the connection library (catalog + configured), normalised for the
  // picker. Same source as the Admin connection screen.
  async function loadConnectionLibrary() {
    const { set = [] } = await $SystemAPI.connectionList({})
    connectionLibrary.value = set.map(c => ({
      catalogID: c.catalogID || '',
      connectionID: c.connectionID ? String(c.connectionID) : null,
      label: c.meta?.short || c.handle || c.catalogID || 'Connection',
      description: c.meta?.description || '',
      handle: c.handle || '',
      status: c.status || '',
      source: c.source || '',
    }))
    return connectionLibrary.value
  }

  // Permitted connection catalogIDs from the Resource Management whitelist, or
  // null when no whitelist exists (e.g. Free mode) — meaning "no constraint".
  function allowedConnectorIds(projectId) {
    const p = findById.value(projectId)
    const wl = p?.governance?.['resource-management']?.values?.connections
    if (!wl || !wl.length) return null
    return new Set(wl.map(c => c.connector).filter(Boolean))
  }

  // A project's connections — the configured connections stamped with its
  // projectID. Fetched from the backend filtered by projectID.
  async function loadConnections(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI.configuredConnectionList({
      projectID: String(projectId),
      limit: 100,
    })
    connectionsByProject.value[key] = set.map(c => ({
      id: String(c.configurationID),
      configuredConnectionID: String(c.configurationID),
      connectionID: String(c.connectionID),
      catalogID: c.catalogID || '',
      name: c.name,
      status: c.status || 'active',
      config: c.config,
    }))
    touch()
    return connectionsByProject.value[key]
  }

  // Picking a connector "does the actual connection": import the catalog entry
  // (when it isn't already a real connection) and read it back so we have its
  // connectionID and auth-field schema (derivedParams) for the configure
  // dialog.
  async function prepareConnection(item) {
    let connectionID = item.connectionID
    if (!connectionID) {
      const imported = await $SystemAPI.connectionImport({ catalogID: item.catalogID })
      connectionID = imported.connectionID
    }
    return $SystemAPI.connectionRead({ connectionID })
  }

  // Create or update a configured connection (the auth/params) on a base
  // connection, scoped to the project. `projectID` is sent so the configured
  // connection is stamped to the project once the backend honours it; extra
  // params are ignored server-side until then.
  async function saveConnection(projectId, { connection, configuredConnectionID, name, config }) {
    const projectID = String(projectId)
    let saved
    if (configuredConnectionID) {
      saved = await $SystemAPI.connectionUpdateConfiguration({
        connectionID: connection.connectionID,
        configuredConnectionID,
        name,
        config,
        labels: {},
        projectID,
      })
    } else {
      saved = await $SystemAPI.connectionConfigure({
        connectionID: connection.connectionID,
        catalogID: connection.catalogID,
        name,
        config,
        labels: {},
        projectID,
      })
    }

    const entry = {
      id: String(saved?.configurationID || configuredConnectionID),
      configuredConnectionID: String(saved?.configurationID || configuredConnectionID),
      connectionID: String(connection.connectionID),
      catalogID: connection.catalogID || '',
      name,
      status: saved?.status || 'active',
      // kept for in-session edit prefill (full listing-by-project is a backend dep)
      config,
    }
    const key = String(projectId)
    const list = connectionsByProject.value[key] || []
    const i = list.findIndex(c => c.id === entry.id)
    connectionsByProject.value[key] =
      i === -1 ? [...list, entry] : list.map(c => (c.id === entry.id ? entry : c))
    touch()
    return entry
  }

  async function removeConnection(projectId, connectionId) {
    const key = String(projectId)
    const entry = (connectionsByProject.value[key] || []).find(c => c.id === connectionId)
    if (entry?.configuredConnectionID) {
      await $SystemAPI.configuredConnectionDelete({ connectionID: entry.configuredConnectionID })
    }
    connectionsByProject.value[key] = (connectionsByProject.value[key] || []).filter(
      c => c.id !== connectionId,
    )
    touch()
  }

  // --- automations (TAQs) --------------------------------------------------------
  // A project's automations — NgAutomation records filtered by projectID. TAQs
  // are created disabled (they have no logic yet), so `disabled: 1` is required
  // to include them in the listing.
  async function loadAutomations(projectId) {
    const key = String(projectId)
    const { set = [] } = await $AutomationAPI
      .ngAutomationList({ projectID: key, disabled: 1, limit: 100 })
      .catch(() => ({ set: [] }))
    automationsByProject.value[key] = set.map(a => ({
      id: String(a.automationID),
      name: a.meta?.short || a.handle || '',
      description: a.meta?.description || '',
      enabled: !!a.enabled,
    }))
    touch()
    return automationsByProject.value[key]
  }

  // Create a project-scoped TAQ. Unlike compose modules (which inherit the
  // project from their namespace) a TAQ has no parent, so the projectID is
  // passed explicitly. Created disabled — the user builds its logic in the TAQ
  // builder afterwards. Returns the new automationID for deep-linking.
  async function addAutomation(projectId, { name, description } = {}) {
    const raw = await $AutomationAPI.ngAutomationCreate({
      projectID: String(projectId),
      meta: {
        short: (name || '').trim() || 'Untitled',
        description: (description || '').trim() || undefined,
      },
      enabled: false,
      triggers: [],
      steps: [],
      paths: [],
    })
    await loadAutomations(projectId)
    return String(raw.automationID)
  }

  // Update a TAQ's name/description. The update endpoint REPLACES triggers,
  // steps and paths, so we re-fetch the full definition first and resend it
  // untouched — editing the name here must never wipe the automation's logic.
  async function updateAutomation(projectId, automationId, { name, description } = {}) {
    const full = await $AutomationAPI.ngAutomationRead({ automationID: automationId })
    await $AutomationAPI.ngAutomationUpdate({
      automationID: automationId,
      handle: full.handle,
      labels: full.labels || {},
      meta: {
        ...(full.meta || {}),
        short: (name || '').trim() || 'Untitled',
        description: (description || '').trim() || undefined,
      },
      enabled: full.enabled,
      scope: full.scope,
      triggers: full.triggers || [],
      steps: full.steps || [],
      paths: full.paths || [],
      runAs: full.runAs,
      ownedBy: full.ownedBy,
      updatedAt: full.updatedAt,
    })
    await loadAutomations(projectId)
  }

  async function removeAutomation(projectId, automationId) {
    const key = String(projectId)
    await $AutomationAPI.ngAutomationDelete({ automationID: automationId })
    automationsByProject.value[key] = (automationsByProject.value[key] || []).filter(
      a => a.id !== automationId,
    )
    touch()
  }

  // --- resource graph ------------------------------------------------------------
  // Single backend endpoint re-derives the project's dependency graph from saved
  // state. We pass nodes through untouched and only rename edge endpoints to the
  // source/target shape the graph component (ECharts) expects.
  async function graph(projectID) {
    const { nodes = [], edges = [] } = await $SystemAPI.projectGraph({ projectID })
    return {
      nodes,
      edges: edges.map(e => ({
        source: e.sourceID,
        target: e.targetID,
        reason: e.reason,
      })),
    }
  }

  return {
    projects,
    loaded,
    graphVersion,
    load,
    fetchProject,
    findById,
    resourcesFor,
    loadResources,
    create,
    updateProject,
    removeProject,
    addMember,
    updateMember,
    removeMember,
    addResource,
    removeResource,
    updateResource,
    connectionLibrary,
    connectionsFor,
    loadConnectionLibrary,
    allowedConnectorIds,
    loadConnections,
    prepareConnection,
    saveConnection,
    removeConnection,
    automationsFor,
    loadAutomations,
    updateAutomation,
    addAutomation,
    removeAutomation,
    updateField,
    addField,
    removeField,
    setFields,
    saveStepForm,
    submitSection,
    transitionStep,
    graph,
    graphVisibleKinds,
    graphShowAccess,
    graphKindVisible,
    graphToggleKind,
    graphToggleAccess,
  }
})
