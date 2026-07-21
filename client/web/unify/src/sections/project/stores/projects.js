import { ACCESS_KINDS, NODE_LAYER_KINDS } from '@/sections/project/config/kinds'
import { SENSITIVITY_LEVELS } from '@/sections/project/config/sensitivity'
import { fieldName } from '@/sections/project/utils/fields'
import { compose, NoID, system } from '@planetcrust/human-js'
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
  let loading = null

  const findById = computed(() => id => projects.value.find(p => p.projectID === String(id)))

  // A project has a compose namespace only once created server-side; the class
  // getter returns NoID ('0') until then, so guards test against NoID.
  const hasNamespace = p => !!p?.namespaceID && p.namespaceID !== NoID

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
  // access overlay (roles/users) rides its own set so each chip toggles
  // independently, both off by default.
  const graphVisibleKinds = ref(new Set(NODE_LAYER_KINDS))
  const graphVisibleAccessKinds = ref(new Set())

  const graphKindVisible = computed(() => kind => {
    if (ACCESS_KINDS.includes(kind)) return graphVisibleAccessKinds.value.has(kind)
    return graphVisibleKinds.value.has(kind)
  })

  function graphToggleKind(kind) {
    // Access kinds ride their own set; node-layer kinds the main one. Reassign
    // the Set so the ref's dependents re-run (Set mutation alone won't).
    const target = ACCESS_KINDS.includes(kind) ? graphVisibleAccessKinds : graphVisibleKinds
    const next = new Set(target.value)
    next.has(kind) ? next.delete(kind) : next.add(kind)
    target.value = next
  }

  // Force the access overlay on/off — used to auto-reveal it on the access-kind
  // steps (roles/users), the same way setGraphVisibleKinds re-seeds node layers.
  // Reveals both role and user chips; each can then be toggled independently.
  function setGraphShowAccess(on) {
    graphVisibleAccessKinds.value = on ? new Set(ACCESS_KINDS) : new Set()
  }

  // Reset the visible node-layer kinds to exactly `kinds` (a Set/iterable) —
  // used to gate the graph to the resources built up to the active step. Access
  // kinds ride their own overlay and are filtered out here. Manual per-kind
  // toggles refine the view within the current step; this re-seeds it on a step
  // change, so a kind from a later step can still be toggled on to peek ahead.
  function setGraphVisibleKinds(kinds) {
    graphVisibleKinds.value = new Set([...kinds].filter(k => NODE_LAYER_KINDS.includes(k)))
  }

  // Resources (compose modules) are deliberately NOT stored on the project
  // object. They live here keyed by projectID and are (re)fetched on demand,
  // filtered by projectID — fetching is the only way to get a project's
  // resources, and a mutation always refetches rather than patching in place.
  const resourcesByProject = ref({})
  const resourcesFor = computed(
    () => projectId => resourcesByProject.value[String(projectId)] || [],
  )

  // The connection library (catalog + already-configured connections) — the
  // same set the Admin connection screen lists. Loaded once and shared across
  // projects; the Connections step filters it down to the Resource Management
  // whitelist. A project's own connections are kept separately, keyed by id.
  const connectionLibrary = ref([])
  const connectionsByProject = ref({})
  const connectionsFor = computed(
    () => projectId => connectionsByProject.value[String(projectId)] || [],
  )

  // A project's automations (TAQs) — NgAutomation records stamped with its
  // projectID. Kept separately, keyed by projectID, and (re)fetched on demand
  // like connections.
  const automationsByProject = ref({})
  const automationsFor = computed(
    () => projectId => automationsByProject.value[String(projectId)] || [],
  )

  // A project's agents — system Agent records stamped with its projectID. Kept
  // separately, keyed by projectID, and (re)fetched on demand like automations.
  const agentsByProject = ref({})
  const agentsFor = computed(() => projectId => agentsByProject.value[String(projectId)] || [])

  // A project's chatbots — system Chatbot records stamped with its projectID.
  // Chatbots carry a top-level name and an `enabled` flag (no description).
  const chatbotsByProject = ref({})
  const chatbotsFor = computed(() => projectId => chatbotsByProject.value[String(projectId)] || [])

  // A project's end-user access roles — system Role records stamped with its
  // projectID (scoped server-side). Distinct from project *members* (the build
  // team): these are the roles end-users get on the deployed app's pages and
  // records. Description lives in role meta.
  const rolesByProject = ref({})
  const rolesFor = computed(() => projectId => rolesByProject.value[String(projectId)] || [])

  // A project's end-users — system users who are members of one or more of the
  // project's access roles. Derived from role membership (there is no separate
  // project-user record); each entry carries the role IDs the user holds.
  const projectUsersByProject = ref({})
  const projectUsersFor = computed(
    () => projectId => projectUsersByProject.value[String(projectId)] || [],
  )

  // A project's compose pages, keyed by projectID. These live in the project's
  // namespace (like modules): record pages (bound to a module, auto-created with
  // each module) and standalone pages the user adds in the Pages step.
  const pagesByProject = ref({})
  const pagesFor = computed(() => projectId => pagesByProject.value[String(projectId)] || [])

  // A project's build-team members. The Project class doesn't hold members
  // (they're a separate resource), so they're cached per-project here — like
  // resources/connections/roles above — and read via membersFor(projectId).
  const membersByProject = ref({})
  const membersFor = computed(() => projectId => membersByProject.value[String(projectId)] || [])

  // --- payload mapping --------------------------------------------------------

  // Merge a fresh backend payload into the cached project (or insert it).
  // Projects are lib `system.Project` instances; `apply` re-runs the class
  // mapping in place so later mutations stay reactive. Members are NOT carried
  // on the instance — they live in membersByProject, fetched by projectID.
  function absorb(raw) {
    const prev = projects.value.find(p => p.projectID === String(raw.projectID))
    if (prev) {
      prev.apply(raw)
      return prev
    }
    projects.value.push(new system.Project(raw))
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
      precision:
        f.kind === 'Number'
          ? Number.isFinite(f.options?.precision)
            ? f.options.precision
            : 0
          : null,
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
    await loadResources(p.projectID)
  }

  // --- loading -----------------------------------------------------------------

  async function load() {
    if (loading) return loading
    loading = $SystemAPI
      .projectList({ limit: 500, sort: 'createdAt DESC' })
      .then(({ set = [] } = {}) => {
        for (const raw of set) absorb(raw)
      })
      .catch(err => console.error('Failed to load projects', err))
      .finally(() => {
        loading = null
      })
    return loading
  }

  // Full fetch for a project view: project + members (capability resolution).
  // The project's resources are loaded separately into the resourcesByProject
  // cache (fetched by projectID), never attached to the project object.
  //
  // Always re-reads (callers want fresh state on entry), but concurrent calls
  // for the same id share one in-flight request so navigating between a
  // project's views doesn't fire duplicate read+members+resources waterfalls.
  const fetchInFlight = new Map()
  async function fetchProject(id) {
    const key = String(id)
    if (fetchInFlight.has(key)) return fetchInFlight.get(key)

    const promise = (async () => {
      const raw = await $SystemAPI.projectRead({ projectID: id })
      const p = absorb(raw)

      await loadSensitivityLevels()

      const [members] = await Promise.all([
        $SystemAPI.projectListMembers({ projectID: p.projectID }).catch(() => ({ set: [] })),
        loadResources(p.projectID),
      ])

      membersByProject.value[p.projectID] = (members.set || []).map(unmarshalMember)

      return p
    })().finally(() => fetchInFlight.delete(key))

    fetchInFlight.set(key, promise)
    return promise
  }

  // Fetch a project's resources (compose modules) filtered by projectID and
  // cache them. This is the single source of resources; callers read them via
  // resourcesFor(projectId), and every resource mutation calls this to resync.
  async function loadResources(projectId) {
    const p = findById.value(projectId)
    if (!hasNamespace(p)) {
      resourcesByProject.value[String(projectId)] = []
      return []
    }
    const { set = [] } = await $ComposeAPI
      .moduleList({ namespaceID: p.namespaceID, projectID: p.projectID, limit: 500 })
      .catch(() => ({ set: [] }))
    resourcesByProject.value[String(projectId)] = set.map(unmarshalModule)
    touch()
    return resourcesByProject.value[String(projectId)]
  }

  // --- project CRUD --------------------------------------------------------------

  // Create a draft project. The backend generates the internal handle, creates
  // the compose namespace and adds the creator as a developer; `deployer`
  // carries the AI Act deployer-category answers collected in
  // NewProjectDialog.vue, which drive the backend's FriaRequired derivation.
  async function create({ name, description = '', deployer = {} } = {}) {
    const raw = await $SystemAPI.projectCreate({
      status: 'draft',
      config: {
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

  // Push the cached project state to the API. The Project instance carries the
  // whole config/meta, so updates round-trip fields the UI doesn't surface yet.
  async function pushProject(p) {
    const raw = await $SystemAPI.projectUpdate({
      projectID: p.projectID,
      handle: p.handle,
      status: p.status,
      config: p.config,
      meta: p.meta,
      labels: p.labels,
      updatedAt: p.updatedAt,
    })
    touch()
    return absorb(raw)
  }

  // Patch a project. Name/description live in meta.
  async function updateProject(id, patch = {}) {
    const p = findById.value(id)
    if (!p) return
    // Snapshot the fields we touch so a failed push rolls back the optimistic
    // mutation instead of leaving the instance showing unsaved state.
    const prev = { status: p.status, meta: { ...p.meta } }
    if (typeof patch.status === 'string') p.status = patch.status
    if (typeof patch.name === 'string') p.meta = { ...p.meta, short: patch.name }
    if (typeof patch.description === 'string')
      p.meta = { ...p.meta, description: patch.description }
    try {
      return await pushProject(p)
    } catch (err) {
      p.status = prev.status
      p.meta = prev.meta
      throw err
    }
  }

  async function removeProject(id) {
    await $SystemAPI.projectDelete({ projectID: id })
    const i = projects.value.findIndex(p => p.projectID === String(id))
    if (i !== -1) projects.value.splice(i, 1)
    delete membersByProject.value[String(id)]
  }

  // Publish a draft project. For an original (no parent revision) this simply
  // promotes it draft→active with no record migration; the returned project
  // carries its new status, which we absorb into the cache. `mappings` stays
  // empty until the revision/migration flow is wired.
  //
  // Every project also requires the `'publish'` governance step to be
  // approved (enforced server-side, unconditionally); on success the backend
  // resets that step back to `draft` and returns it that way, so absorbing
  // the response here already leaves the UI showing the fresh, unapproved
  // cycle for next time — no separate governance refetch needed.
  async function publishProject(id) {
    const p = findById.value(id)
    if (!p) return
    const raw = await $SystemAPI.projectPublish({
      projectID: p.projectID,
      confirm: true,
      mappings: [],
    })
    touch()
    return absorb(raw)
  }

  // --- members --------------------------------------------------------------
  // Memberships are keyed by user (one record per user per project); the role
  // is a fixed named preset and capabilities come back derived from it.
  // Mutations are RBAC-checked server-side (project members.manage).

  async function addMember(projectId, { userId, role } = {}) {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectAddMember({
      projectID: p.projectID,
      userID: userId,
      rolePreset: role,
    })
    const key = String(projectId)
    membersByProject.value[key] = [...(membersByProject.value[key] || []), unmarshalMember(raw)]
    touch()
  }

  async function updateMember(projectId, userId, role) {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectUpdateMember({
      projectID: p.projectID,
      userID: userId,
      rolePreset: role,
    })
    const key = String(projectId)
    const next = unmarshalMember(raw)
    membersByProject.value[key] = (membersByProject.value[key] || []).map(m =>
      m.userId === String(userId) ? { ...m, ...next } : m,
    )
    touch()
  }

  async function removeMember(projectId, userId) {
    const p = findById.value(projectId)
    if (!p) return
    await $SystemAPI.projectRemoveMember({ projectID: p.projectID, userID: userId })
    const key = String(projectId)
    membersByProject.value[key] = (membersByProject.value[key] || []).filter(
      m => m.userId !== String(userId),
    )
    touch()
  }

  // --- resources (modules) -------------------------------------------------------
  // All wizard resources are real compose modules in the project namespace.

  async function addResource(projectId, { kind, name } = {}) {
    const p = findById.value(projectId)
    if (!p || kind !== 'module' || !hasNamespace(p)) return null

    const raw = await $ComposeAPI.moduleCreate({
      namespaceID: p.namespaceID,
      projectID: String(p.projectID),
      name: (name || '').trim() || 'Untitled',
      // Handle is the slugified title (no hash). Duplicate titles collide on
      // the unique handle; the create dialog validates against that first.
      handle: fieldName(name).toLowerCase(),
      fields: [],
      meta: {},
    })

    // Mirror the compose module editor: every new module gets a record page (+
    // default layout) so it surfaces in the Pages step. Best-effort and
    // intentionally non-atomic — the module already exists, so a page failure
    // must not fail the module create.
    try {
      await createRecordPageForModule(
        p.namespaceID,
        String(raw.moduleID),
        (name || '').trim() || 'Untitled',
        String(p.projectID),
      )
    } catch (err) {
      console.error('Failed to create record page for module:', err)
    }

    await loadResources(projectId)
    return String(raw.moduleID)
  }

  async function removeResource(projectId, resourceId) {
    const p = findById.value(projectId)
    const r = moduleOf(projectId, resourceId)
    if (!r) return

    await $ComposeAPI.moduleDelete({ namespaceID: p.namespaceID, moduleID: r.id })

    // Drop any record pages bound to this module so the Pages step shows no
    // orphans. Best-effort — the module is already gone.
    try {
      const { set = [] } = await $ComposeAPI
        .pageList({ namespaceID: p.namespaceID, moduleID: r.id })
        .catch(() => ({ set: [] }))
      for (const pg of set) {
        await $ComposeAPI.pageDelete({ namespaceID: p.namespaceID, pageID: pg.pageID })
      }
    } catch (err) {
      console.error('Failed to remove record pages for module:', err)
    }

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
            typeof o === 'string'
              ? { value: o, text: o }
              : { value: o.value ?? '', text: o.text ?? '' },
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
  // Form-type steps (summary, resource-management) use saveStepForm to persist
  // FORM VALUES, kept in memory ONLY (never sent to the backend as part of
  // governance) — see the form/sensitivity steps. That's a separate concern
  // from *governance status*, which transitionStep() below drives and which IS
  // always backend-persisted, for every step key:
  //
  //   - the well-known `'publish'` step key (PUBLISH_GOVERNANCE_STEP_KEY) runs
  //     the full submit -> approve/request-changes cycle that gates
  //     publishing itself, unconditionally (server/system/service/
  //     project_revision.go's Publish() requires it approved, and resets it
  //     back to draft on success — see publishProject() below, whose
  //     response already carries that reset);
  //   - every OTHER (Build/Govern) step key has NO submit stage: a member with
  //     grant-approval capability can send either 'approve' or
  //     'request-changes' directly, from any current status, at any time (see
  //     Wizard.vue's per-step Approve / "Request changes" toolbar actions).
  //     'request-changes' flags that step to changes-requested and sends
  //     'publish' back for review too if it was submitted/approved
  //     (server/system/service/project_governance.go); 'approve' clears a
  //     changes-requested flag (or simply marks a draft step approved).
  //
  // Every transition, regardless of stepKey or action, round-trips through the
  // real governance API and absorbs the response, keeping every viewer in
  // sync with the same state — there is no in-memory mock of the state
  // machine here.

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
    touch()
  }

  async function transitionStep(projectId, stepKey, action, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectGovernanceTransition({
      projectID: p.projectID,
      stepKey,
      action,
      note,
    })
    touch()
    return absorb(raw)
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
  // null when no whitelist has been declared yet — meaning "no constraint".
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
  // connection, scoped to the project. `projectID` is sent and the backend
  // stamps it onto the configured connection (rel_project), so it shows in the
  // project-scoped resource graph.
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

  // Update a TAQ's name/description/enabled. The update endpoint REPLACES
  // triggers, steps and paths, so we re-fetch the full definition first and
  // resend it untouched — editing here must never wipe the automation's logic.
  // `enabled` falls back to the stored value when not supplied.
  async function updateAutomation(projectId, automationId, { name, description, enabled } = {}) {
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
      enabled: enabled === undefined ? full.enabled : enabled,
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

  // --- agents --------------------------------------------------------------------
  // A project's agents — system Agent records filtered by projectID. Agents
  // carry a status (`active`/`inactive`) rather than an enabled flag; the
  // listing returns both, so no extra filter is needed (unlike TAQs).
  async function loadAgents(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI
      .agentList({ projectID: key, limit: 100 })
      .catch(() => ({ set: [] }))
    agentsByProject.value[key] = set.map(a => ({
      id: String(a.agentID),
      name: a.meta?.short || a.handle || '',
      description: a.meta?.description || '',
      status: a.status || 'inactive',
    }))
    touch()
    return agentsByProject.value[key]
  }

  // Create a project-scoped agent. Like a TAQ it has no parent, so the projectID
  // is passed explicitly. Created inactive — the user builds its behaviour in the
  // agent builder afterwards. Returns the new agentID for deep-linking.
  async function addAgent(projectId, { name, description } = {}) {
    const raw = await $SystemAPI.agentCreate({
      projectID: String(projectId),
      status: 'inactive',
      meta: {
        short: (name || '').trim() || 'Untitled',
        description: (description || '').trim() || undefined,
      },
    })
    await loadAgents(projectId)
    return String(raw.agentID)
  }

  // Update an agent's name/description/status. agentUpdate REPLACES the whole
  // definition, so we re-fetch the full agent first and resend it untouched —
  // editing here must never wipe the agent's behaviour, execution or access.
  // `status` falls back to the stored value when not supplied.
  async function updateAgent(projectId, agentId, { name, description, status } = {}) {
    const full = await $SystemAPI.agentRead({ agentID: agentId })
    await $SystemAPI.agentUpdate({
      agentID: agentId,
      handle: full.handle,
      labels: full.labels || {},
      status: status === undefined ? full.status : status,
      meta: {
        ...(full.meta || {}),
        short: (name || '').trim() || 'Untitled',
        description: (description || '').trim() || undefined,
      },
      behavior: full.behavior,
      execution: full.execution,
      access: full.access,
      invocation: full.invocation,
      updatedAt: full.updatedAt,
    })
    await loadAgents(projectId)
  }

  async function removeAgent(projectId, agentId) {
    const key = String(projectId)
    await $SystemAPI.agentDelete({ agentID: agentId })
    agentsByProject.value[key] = (agentsByProject.value[key] || []).filter(a => a.id !== agentId)
    touch()
  }

  // --- chatbots ------------------------------------------------------------------
  // A project's chatbots — system Chatbot records filtered by projectID. Chatbots
  // have a name and an `enabled` flag; the listing returns both regardless of
  // state, so no extra filter is needed.
  async function loadChatbots(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI
      .chatbotList({ projectID: key, limit: 100 })
      .catch(() => ({ set: [] }))
    chatbotsByProject.value[key] = set.map(c => ({
      id: String(c.chatbotID),
      name: c.name || c.handle || '',
      enabled: !!c.enabled,
    }))
    touch()
    return chatbotsByProject.value[key]
  }

  // Create a project-scoped chatbot. Created disabled — the user configures its
  // scenarios, styling and channels in the chatbot builder afterwards. Returns
  // the new chatbotID for deep-linking.
  async function addChatbot(projectId, { name } = {}) {
    const raw = await $SystemAPI.chatbotCreate({
      projectID: String(projectId),
      name: (name || '').trim() || 'Untitled',
      enabled: false,
    })
    await loadChatbots(projectId)
    return String(raw.chatbotID)
  }

  // Update a chatbot's name/enabled. chatbotUpdate REPLACES the whole definition,
  // so we re-fetch the full chatbot first and resend it untouched — editing here
  // must never wipe its scenarios, styling or handoff config. `enabled` falls
  // back to the stored value when not supplied.
  async function updateChatbot(projectId, chatbotId, { name, enabled } = {}) {
    const full = await $SystemAPI.chatbotRead({ chatbotID: chatbotId })
    await $SystemAPI.chatbotUpdate({
      chatbotID: chatbotId,
      handle: full.handle,
      labels: full.labels || {},
      name: (name || '').trim() || 'Untitled',
      enabled: enabled === undefined ? full.enabled : enabled,
      sessionTTL: full.sessionTTL,
      allowedOrigins: full.allowedOrigins,
      handoff: full.handoff,
      styling: full.styling,
      scenarios: full.scenarios,
      updatedAt: full.updatedAt,
    })
    await loadChatbots(projectId)
  }

  async function removeChatbot(projectId, chatbotId) {
    const key = String(projectId)
    await $SystemAPI.chatbotDelete({ chatbotID: chatbotId })
    chatbotsByProject.value[key] = (chatbotsByProject.value[key] || []).filter(
      c => c.id !== chatbotId,
    )
    touch()
  }

  // --- access roles --------------------------------------------------------------
  // End-user access roles scoped to the project via projectID. roleList filters
  // server-side; description is carried in role meta.
  async function loadRoles(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI
      .roleList({ projectID: key, limit: 100 })
      .catch(() => ({ set: [] }))
    rolesByProject.value[key] = set.map(r => ({
      id: String(r.roleID),
      name: r.name || r.handle || '',
      description: r.meta?.description || '',
    }))
    touch()
    return rolesByProject.value[key]
  }

  // Create a project-scoped access role. The handle is project-prefixed so two
  // projects can each define e.g. "Customer" without colliding on the global
  // unique handle.
  async function addRole(projectId, { name, description = '' } = {}) {
    const title = (name || '').trim() || 'Untitled'
    const raw = await $SystemAPI.roleCreate({
      projectID: String(projectId),
      name: title,
      handle: `proj_${projectId}_${fieldName(title).toLowerCase()}`,
      meta: { description: description.trim() },
    })
    await loadRoles(projectId)
    return String(raw.roleID)
  }

  // Update an access role's name/description. roleUpdate REPLACES the role, so we
  // re-fetch it first and resend it untouched — editing here must never wipe the
  // role's members, handle or archived state. The handle is kept as-is (renaming
  // a role never re-slugs it, so it can't collide with a sibling).
  async function updateRole(projectId, roleId, { name, description } = {}) {
    const full = await $SystemAPI.roleRead({ roleID: roleId })
    await $SystemAPI.roleUpdate({
      ...full,
      roleID: roleId,
      handle: full.handle,
      name: name === undefined ? full.name : (name || '').trim() || 'Untitled',
      meta: {
        ...(full.meta || {}),
        description:
          description === undefined ? full.meta?.description : (description || '').trim(),
      },
      updatedAt: full.updatedAt,
    })
    await loadRoles(projectId)
  }

  async function removeRole(projectId, roleId) {
    const key = String(projectId)
    await $SystemAPI.roleDelete({ roleID: roleId })
    rolesByProject.value[key] = (rolesByProject.value[key] || []).filter(r => r.id !== roleId)
    touch()
  }

  // Access permissions (compose/system/automation RBAC) are edited through the
  // shared permissions modal (CPermissionsDialog), which reads/writes the rules
  // directly against the matching API. The Permissions step opens it per
  // resource+role and passes `onSaved: touch` so the graph refetches the updated
  // role→resource edges — so no role-permission read/write lives on the store.

  // --- effective access evaluation -----------------------------------------------
  // Per-project cache of each role's effective access on each resource, across
  // ALL operations (read/create/update/delete/…), computed via permissionsTrace,
  // so the Permissions matrix can tint each cell for the selected op. Keyed
  // `${roleId}::${resource}` → { [operation]: 'allow' | 'deny' | 'unknown' }.
  // Versioned to graphVersion, so touch() (fired on every permission save +
  // resource/role mutation) invalidates it and the matrix recomputes.
  const effectiveAccessByProject = ref({}) // { [pid]: { version, map: Map } }
  const effectiveAccessLoading = ref({}) // { [pid]: bool }

  // Route a resource to the API client that owns it (mirrors CPermissionsDialog).
  // Handles both object resources (corteza::compose:module/…) and component-level
  // resources (corteza::automation/), where the service is followed by "/" not ":".
  function apiForResource(resource) {
    switch (resource.match(/^corteza::(\w+)[:/]/)?.[1]) {
      case 'compose':
        return $ComposeAPI
      case 'automation':
        return $AutomationAPI
      default:
        return $SystemAPI
    }
  }

  // Effective access for one (role, resource, operation): 'allow' | 'deny' |
  // 'unknown', or undefined when not evaluated / the op wasn't returned.
  function effectiveAccess(projectId, roleId, resource, operation) {
    return effectiveAccessByProject.value[String(projectId)]?.map.get(`${roleId}::${resource}`)?.[
      operation
    ]
  }

  function isEffectiveAccessLoading(projectId) {
    return !!effectiveAccessLoading.value[String(projectId)]
  }

  // Quick-toggle: set a single role's direct rule for one (resource, operation)
  // and refresh. Optimistically patches the cache so the cell flips instantly;
  // touch() then re-traces so the tint reflects the true effective access.
  async function setAccess(projectId, roleId, resource, operation, access) {
    const key = String(projectId)
    const entry = effectiveAccessByProject.value[key]
    if (entry) {
      const cell = `${roleId}::${resource}`
      const map = new Map(entry.map)
      map.set(cell, { ...(map.get(cell) || {}), [operation]: access })
      effectiveAccessByProject.value = {
        ...effectiveAccessByProject.value,
        [key]: { version: entry.version, map },
      }
    }
    try {
      await apiForResource(resource).permissionsUpdate({
        roleID: roleId,
        rules: [{ resource, operation, access }],
      })
    } finally {
      touch() // re-trace effective access (and refresh the resource graph)
    }
  }

  // Set a *capability* — a user-facing action that may map to several RBAC ops on
  // possibly different resources (e.g. "view records" = record.read on the record
  // wildcard + records.search on the module). Writes them all to one `access` in a
  // single pass: optimistically patches every affected cell, groups the rules by
  // the API component that owns each resource, and issues one permissionsUpdate
  // per component, then re-traces. `ops` is [{ resource, op }].
  async function setCapabilityAccess(projectId, roleId, ops = [], access) {
    if (!ops.length) return
    const key = String(projectId)
    const entry = effectiveAccessByProject.value[key]
    if (entry) {
      const map = new Map(entry.map)
      for (const { resource, op } of ops) {
        const cell = `${roleId}::${resource}`
        map.set(cell, { ...(map.get(cell) || {}), [op]: access })
      }
      effectiveAccessByProject.value = {
        ...effectiveAccessByProject.value,
        [key]: { version: entry.version, map },
      }
    }
    try {
      // permissionsUpdate is per API component; a component's rules may span
      // several of its resources, so group by client and send one call each.
      const byApi = new Map()
      for (const { resource, op } of ops) {
        const api = apiForResource(resource)
        if (!byApi.has(api)) byApi.set(api, [])
        byApi.get(api).push({ resource, operation: op, access })
      }
      await Promise.all(
        [...byApi.entries()].map(([api, rules]) =>
          api.permissionsUpdate({ roleID: roleId, rules }),
        ),
      )
    } finally {
      touch() // re-trace effective access (and refresh the resource graph)
    }
  }

  const accessLoadingVersion = {} // pid → the graphVersion currently being fetched
  const accessSeqByProject = {} // pid → monotonic run token; per-project stale-write guard

  // Compute effective access for every (role, resource). One permissionsTrace
  // call per (role, API-component), passing that component's whole resource list
  // as a `resource[]` array (chunked to keep the GET query string sane) — the
  // server returns one trace per resource×op, so we keep ALL operations and N
  // resources cost O(roles × components) calls, not O(roles × N). Results survive
  // expand/collapse and axis switches (cached by pid, all ops present).
  async function loadEffectiveAccess(projectId, resources = []) {
    const key = String(projectId)
    const version = graphVersion.value
    const cached = effectiveAccessByProject.value[key]
    if (cached && cached.version === version) return
    if (accessLoadingVersion[key] === version) return // already fetching this version
    accessLoadingVersion[key] = version

    // Per-project run token: a newer run for THIS project supersedes this one.
    // Keyed by pid so a concurrent load for another project can't invalidate it.
    const seq = (accessSeqByProject[key] = (accessSeqByProject[key] || 0) + 1)
    const isLatest = () => seq === accessSeqByProject[key]
    effectiveAccessLoading.value = { ...effectiveAccessLoading.value, [key]: true }

    const map = new Map()
    try {
      const roles = rolesFor.value(projectId)
      if (roles.length && resources.length) {
        // Group resources by API client, chunked (~60 → GET query stays well
        // under proxy/URL limits at ~50 chars/resource).
        const byApi = new Map()
        for (const r of resources) {
          const api = apiForResource(r)
          if (!byApi.has(api)) byApi.set(api, [])
          byApi.get(api).push(r)
        }
        const CHUNK = 60
        await Promise.all(
          roles.flatMap(role =>
            [...byApi.entries()].flatMap(([api, list]) => {
              const chunks = []
              for (let i = 0; i < list.length; i += CHUNK) chunks.push(list.slice(i, i + CHUNK))
              return chunks.map(chunk =>
                api
                  .permissionsTrace({ resource: chunk, roleID: [role.id] })
                  .then(traces => {
                    for (const tr of traces || []) {
                      if (!tr?.resource || !tr.operation) continue
                      const state =
                        tr.resolution === 'unknown-context'
                          ? 'unknown'
                          : tr.access === 'allow'
                            ? 'allow'
                            : 'deny'
                      const cell = `${role.id}::${tr.resource}`
                      let ops = map.get(cell)
                      if (!ops) map.set(cell, (ops = {}))
                      ops[tr.operation] = state
                    }
                  })
                  // No grant (4xx) or transient failure → leave those cells unset
                  // so the matrix dims them rather than erroring.
                  .catch(() => {}),
              )
            }),
          ),
        )
      }
      // Only the latest run for this project writes (empty map = no roles/resources).
      if (isLatest()) {
        effectiveAccessByProject.value = {
          ...effectiveAccessByProject.value,
          [key]: { version, map },
        }
      }
    } finally {
      // Only the latest run clears the flags — a superseded run must not hide the
      // spinner or unpin the version while a newer fetch is still in flight.
      if (isLatest()) {
        if (accessLoadingVersion[key] === version) accessLoadingVersion[key] = null
        effectiveAccessLoading.value = { ...effectiveAccessLoading.value, [key]: false }
      }
    }
  }

  // --- per-user effective access -------------------------------------------------
  // A user's *resolved* effective access across ALL the roles they hold, traced
  // per resource (all ops). Distinct from the per-role cache above: the server
  // resolves the user's whole role set, so this is the read-only "Evaluated"
  // column in the per-user permissions view. Keyed `${pid}::${userId}`, mapped
  // `resource` → { [operation]: 'allow' | 'deny' | 'unknown' }. Versioned to
  // graphVersion so touch() (fired on every permission/membership change)
  // invalidates it and the column re-traces.
  const userEffectiveAccessByPU = ref({}) // { [`${pid}::${uid}`]: { version, map: Map } }
  const userEffectiveAccessLoading = ref({}) // { [`${pid}::${uid}`]: bool }

  function userEffectiveAccess(projectId, userId, resource, operation) {
    return userEffectiveAccessByPU.value[`${projectId}::${userId}`]?.map.get(resource)?.[operation]
  }

  function isUserEffectiveAccessLoading(projectId, userId) {
    return !!userEffectiveAccessLoading.value[`${projectId}::${userId}`]
  }

  const userAccessLoadingVersion = {} // `${pid}::${uid}` → graphVersion being fetched
  const userAccessSeq = {} // `${pid}::${uid}` → monotonic run token; stale-write guard

  // Trace one user's effective access over every resource (mirrors
  // loadEffectiveAccess but passes userID instead of roleID, so the platform
  // resolves the user's full role set into one allow/deny per resource×op).
  async function loadUserEffectiveAccess(projectId, userId, resources = []) {
    if (!userId) return
    const key = `${projectId}::${userId}`
    const version = graphVersion.value
    const cached = userEffectiveAccessByPU.value[key]
    if (cached && cached.version === version) return
    if (userAccessLoadingVersion[key] === version) return // already fetching this version
    userAccessLoadingVersion[key] = version

    const seq = (userAccessSeq[key] = (userAccessSeq[key] || 0) + 1)
    const isLatest = () => seq === userAccessSeq[key]
    userEffectiveAccessLoading.value = { ...userEffectiveAccessLoading.value, [key]: true }

    const map = new Map()
    try {
      if (resources.length) {
        const byApi = new Map()
        for (const r of resources) {
          const api = apiForResource(r)
          if (!byApi.has(api)) byApi.set(api, [])
          byApi.get(api).push(r)
        }
        const CHUNK = 60
        await Promise.all(
          [...byApi.entries()].flatMap(([api, list]) => {
            const chunks = []
            for (let i = 0; i < list.length; i += CHUNK) chunks.push(list.slice(i, i + CHUNK))
            return chunks.map(chunk =>
              api
                // userID is a scalar here (roleID stays an empty array) — the
                // trace endpoint resolves the user's whole role set server-side.
                // Mirrors the admin permissions grid; passing userID as an array
                // yields no resolution.
                .permissionsTrace({ resource: chunk, userID: userId, roleID: [] })
                .then(traces => {
                  for (const tr of traces || []) {
                    if (!tr?.resource || !tr.operation) continue
                    const state =
                      tr.resolution === 'unknown-context'
                        ? 'unknown'
                        : tr.access === 'allow'
                          ? 'allow'
                          : 'deny'
                    let ops = map.get(tr.resource)
                    if (!ops) map.set(tr.resource, (ops = {}))
                    ops[tr.operation] = state
                  }
                })
                .catch(() => {}),
            )
          }),
        )
      }
      if (isLatest()) {
        userEffectiveAccessByPU.value = {
          ...userEffectiveAccessByPU.value,
          [key]: { version, map },
        }
      }
    } finally {
      if (isLatest()) {
        if (userAccessLoadingVersion[key] === version) userAccessLoadingVersion[key] = null
        userEffectiveAccessLoading.value = { ...userEffectiveAccessLoading.value, [key]: false }
      }
    }
  }

  // --- project users -------------------------------------------------------------
  // Build the project-user list by collecting each access role's members and
  // grouping by user (a user may hold several roles).
  async function loadProjectUsers(projectId) {
    const key = String(projectId)
    await loadRoles(projectId)
    const roles = rolesFor.value(projectId)
    const byUser = {}
    await Promise.all(
      roles.map(async role => {
        const ids = await $SystemAPI.roleMemberList({ roleID: role.id }).catch(() => [])
        for (const uid of ids || []) {
          const u = String(uid)
          ;(byUser[u] || (byUser[u] = new Set())).add(role.id)
        }
      }),
    )
    projectUsersByProject.value[key] = Object.entries(byUser).map(([userId, set]) => ({
      userId,
      roleIds: [...set],
    }))
    touch()
    return projectUsersByProject.value[key]
  }

  // Add/remove a user's membership in one access role, then refresh the list.
  async function setProjectUserRole(projectId, userId, roleId, on) {
    if (on) await $SystemAPI.roleMemberAdd({ roleID: roleId, userID: userId })
    else await $SystemAPI.roleMemberRemove({ roleID: roleId, userID: userId })
    await loadProjectUsers(projectId)
  }

  // Assign a user to several access roles at once (used when adding a user), then
  // refresh once rather than per role.
  async function assignProjectUserRoles(projectId, userId, roleIds = []) {
    await Promise.all(
      roleIds.map(roleId => $SystemAPI.roleMemberAdd({ roleID: roleId, userID: userId })),
    )
    await loadProjectUsers(projectId)
  }

  // Remove a user from the project entirely — drop them from every project role.
  async function removeProjectUser(projectId, userId, roleIds = []) {
    await Promise.all(
      roleIds.map(roleId =>
        $SystemAPI.roleMemberRemove({ roleID: roleId, userID: userId }).catch(() => {}),
      ),
    )
    await loadProjectUsers(projectId)
  }

  // The default user group new users join — userCreate requires one. Resolved
  // once (preferring the well-known default handles) and cached.
  let defaultUserGroupID = null
  async function ensureDefaultUserGroup() {
    if (defaultUserGroupID) return defaultUserGroupID
    const { set = [] } = await $SystemAPI.userGroupList({ limit: 100 }).catch(() => ({ set: [] }))
    const g =
      set.find(
        x =>
          x.handle === 'default-root' || x.handle === 'users' || x.meta?.short?.includes('Default'),
      ) || set[0]
    defaultUserGroupID = g ? String(g.userGroupID) : null
    return defaultUserGroupID
  }

  // Create a brand-new platform user (email + name) to assign to project roles.
  // Returns the new userID; the caller assigns the role membership.
  async function addProjectUser(_projectId, { email, name } = {}) {
    const userGroupID = await ensureDefaultUserGroup()
    const raw = await $SystemAPI.userCreate({
      email: (email || '').trim(),
      name: (name || '').trim(),
      userGroupID,
    })
    return String(raw.userID)
  }

  // --- pages ---------------------------------------------------------------------
  // Compose pages in the project's namespace. Record pages are bound to a module
  // (auto-created with it); standalone pages are added in the Pages step. Pages
  // are stamped with the project via the namespace (backend page onCreate).

  async function loadPages(projectId) {
    const p = findById.value(projectId)
    const key = String(projectId)
    if (!hasNamespace(p)) {
      pagesByProject.value[key] = []
      return []
    }
    const { set = [] } = await $ComposeAPI
      .pageList({ namespaceID: p.namespaceID, limit: 500 })
      .catch(() => ({ set: [] }))
    // Layouts live as a separate compose resource (pageList doesn't return them).
    // One namespace-wide call — every layout carries its pageID — then group by
    // page, so the permissions matrix can nest layout rows the way modules nest
    // fields. limit:500 doesn't page the cursor (same cap as pageList above).
    const { set: layoutSet = [] } = await $ComposeAPI
      .pageLayoutListNamespace({ namespaceID: p.namespaceID, limit: 500 })
      .catch(() => ({ set: [] }))
    const layoutsByPage = new Map()
    for (const pl of layoutSet) {
      const pid = String(pl.pageID)
      const list = layoutsByPage.get(pid) || []
      list.push({
        id: String(pl.pageLayoutID),
        handle: pl.handle || '',
        title: pl.meta?.title || '',
        primary: pl.handle === 'primary',
      })
      layoutsByPage.set(pid, list)
    }
    const capitalize = s => (s ? s.charAt(0).toUpperCase() + s.slice(1) : s)
    pagesByProject.value[key] = set.map(pg => {
      const moduleID = pg.moduleID && String(pg.moduleID) !== '0' ? String(pg.moduleID) : null
      const pageName = pg.title || pg.handle || ''
      const rawLayouts = layoutsByPage.get(String(pg.pageID)) || []
      const soleLayout = rawLayouts.length === 1
      return {
        id: String(pg.pageID),
        name: pageName,
        description: pg.description || '',
        visible: !!pg.visible,
        moduleID,
        // Record pages are bound to a module; standalone pages are not.
        isRecordPage: !!moduleID,
        // Parent page for the nav hierarchy; '0'/NoID means a root page.
        selfID: pg.selfID && String(pg.selfID) !== '0' ? String(pg.selfID) : null,
        // Sibling ordering within a parent (compose page weight).
        weight: Number(pg.weight) || 0,
        // Matrix row label: a custom title that isn't just a copy of the page
        // name always wins. Otherwise, when the page has a single layout whose
        // title merely mirrors the page (the auto-created default via
        // createDefaultPageLayout), show "Primary" (capitalized handle) rather
        // than a child row identical to its parent page.
        layouts: rawLayouts.map(l => {
          let name
          if (l.title && l.title !== pageName) {
            name = l.title
          } else if (soleLayout) {
            name = capitalize(l.handle) || 'Primary'
          } else {
            name = l.handle || l.title || 'layout'
          }
          return { ...l, name }
        }),
      }
    })
    touch()
    return pagesByProject.value[key]
  }

  // Seed the layout from the page's blocks (blockID + xywh) so the Builder shows
  // them — an empty layout would hide every block until the user saves. Mirrors
  // the compose module editor's default-layout creation.
  async function createDefaultPageLayout(namespaceID, page) {
    if (!page?.pageID) return
    const blocks = (page.blocks || []).map(b => ({ blockID: b.blockID, xywh: b.xywh }))
    await $ComposeAPI.pageLayoutCreate(
      new compose.PageLayout({
        namespaceID,
        pageID: page.pageID,
        handle: 'primary',
        meta: { title: page.title },
        blocks,
      }),
    )
  }

  // Replicates the compose module editor's "Create record page" button: a page
  // bound to the module with a single Record block, plus a default layout. Named
  // "<Module> Details" per the detail-page convention.
  async function createRecordPageForModule(namespaceID, moduleID, name, projectID) {
    const page = new compose.Page({
      namespaceID,
      projectID,
      moduleID,
      selfID: '0',
      title: `${name} Details`,
      blocks: [new compose.PageBlockRecord({ xywh: [0, 0, 48, 36] })],
    })
    const created = await $ComposeAPI.pageCreate(page)
    await createDefaultPageLayout(namespaceID, created)
    return created
  }

  // Create a standalone page (no module binding). Visible by default so it shows
  // in the namespace navigation; the user builds its blocks in the page builder.
  async function addPage(projectId, { name } = {}) {
    const p = findById.value(projectId)
    if (!hasNamespace(p)) return null
    const page = new compose.Page({
      namespaceID: p.namespaceID,
      projectID: String(p.projectID),
      title: (name || '').trim() || 'Untitled',
      visible: true,
      blocks: [],
    })
    const created = await $ComposeAPI.pageCreate(page)
    await createDefaultPageLayout(p.namespaceID, created)
    await loadPages(projectId)
    return String(created.pageID)
  }

  // Update a page's title/visibility/parent. pageUpdate REPLACES the page, so we
  // spread the full record and override only the given fields — never wiping
  // blocks, layout binding or module link. `selfID` reparents the page in the nav
  // tree (used by the Pages step drag-to-nest); pass '0' for a root page.
  async function updatePage(projectId, pageId, { name, visible, selfID, description } = {}) {
    const p = findById.value(projectId)
    if (!hasNamespace(p)) return
    const full = await $ComposeAPI.pageRead({ namespaceID: p.namespaceID, pageID: pageId })
    await $ComposeAPI.pageUpdate({
      ...full,
      namespaceID: p.namespaceID,
      pageID: pageId,
      title: name === undefined ? full.title : (name || '').trim() || 'Untitled',
      description: description === undefined ? full.description : (description || '').trim(),
      visible: visible === undefined ? full.visible : visible,
      selfID: selfID === undefined ? full.selfID : selfID,
    })
    await loadPages(projectId)
  }

  // Reorder sibling pages under one parent (selfID '0' = root). pageReorder
  // stamps each page's weight from the array order, which drives the compose
  // namespace sidebar. The caller reparents (updatePage selfID) before reordering
  // so pageIDs are all true siblings under selfID.
  async function reorderPages(projectId, selfID, pageIDs) {
    const p = findById.value(projectId)
    if (!hasNamespace(p) || !pageIDs?.length) return
    await $ComposeAPI.pageReorder({
      namespaceID: p.namespaceID,
      selfID: selfID || '0',
      pageIDs,
    })
  }

  async function removePage(projectId, pageId) {
    const p = findById.value(projectId)
    const key = String(projectId)
    await $ComposeAPI.pageDelete({ namespaceID: p.namespaceID, pageID: pageId })
    pagesByProject.value[key] = (pagesByProject.value[key] || []).filter(pg => pg.id !== pageId)
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
    graphVersion,
    touch,
    load,
    fetchProject,
    findById,
    membersFor,
    resourcesFor,
    loadResources,
    create,
    updateProject,
    removeProject,
    publishProject,
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
    agentsFor,
    loadAgents,
    addAgent,
    updateAgent,
    removeAgent,
    chatbotsFor,
    loadChatbots,
    addChatbot,
    updateChatbot,
    removeChatbot,
    rolesFor,
    loadRoles,
    addRole,
    updateRole,
    removeRole,
    effectiveAccess,
    isEffectiveAccessLoading,
    loadEffectiveAccess,
    userEffectiveAccess,
    isUserEffectiveAccessLoading,
    loadUserEffectiveAccess,
    setAccess,
    setCapabilityAccess,
    projectUsersFor,
    loadProjectUsers,
    setProjectUserRole,
    assignProjectUserRoles,
    removeProjectUser,
    addProjectUser,
    pagesFor,
    loadPages,
    addPage,
    updatePage,
    reorderPages,
    removePage,
    updateField,
    addField,
    removeField,
    setFields,
    saveStepForm,
    transitionStep,
    graph,
    graphVisibleKinds,
    graphVisibleAccessKinds,
    graphKindVisible,
    graphToggleKind,
    setGraphShowAccess,
    setGraphVisibleKinds,
  }
})
