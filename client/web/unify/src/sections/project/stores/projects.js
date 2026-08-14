import { OVERVIEW_KINDS } from '@/sections/project/config/kinds'
import { GRAPH_KIND_BY_RESOURCE_TYPE } from '@/sections/project/config/resourceRefs'
import { SENSITIVITY_LEVELS } from '@/sections/project/config/sensitivity'
import { fieldName } from '@/sections/project/utils/fields'
import { compose, NoID, system } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// Projects store. Most of this persists on the backend and grows alongside
// the pipeline, one verified step at a time. Current surface: project CRUD,
// members (role-preset CRUD + capability resolution), and modules (real
// compose modules in the project's namespace).
//
// The one exception is the per-step GOVERNANCE workflow (see the dedicated
// section below): it is intentionally session-local scaffolding, not
// backend-persisted — see that section's comment for why.
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

  // --- resource graph view (kind filters) ----------------------------------------
  // Which kinds the graph currently shows. One global selection shared across
  // projects (it resets on a hard reload — it's deliberately session state, not
  // persisted). Every kind starts visible, including roles/users: the graph is
  // the whole system, and what to leave out is the reader's call — not the
  // build step's, which is why nothing re-seeds this on navigation.
  const graphVisibleKinds = ref(new Set(OVERVIEW_KINDS))

  const graphKindVisible = computed(() => kind => graphVisibleKinds.value.has(kind))

  function graphToggleKind(kind) {
    // Reassign the Set so the ref's dependents re-run (Set mutation alone won't).
    const next = new Set(graphVisibleKinds.value)
    next.has(kind) ? next.delete(kind) : next.add(kind)
    graphVisibleKinds.value = next
  }

  // Show or hide every kind at once — the graph header's all/none buttons.
  function setGraphAllKindsVisible(on) {
    graphVisibleKinds.value = on ? new Set(OVERVIEW_KINDS) : new Set()
  }

  // Isolate one kind (double-click on its chip). Doing it again to the kind
  // that is already alone brings everything back, so solo reads as a two-state
  // zoom rather than a trap you need the all-button to escape.
  function graphSoloKind(kind) {
    const cur = graphVisibleKinds.value
    const alone = cur.size === 1 && cur.has(kind)
    graphVisibleKinds.value = alone ? new Set(OVERVIEW_KINDS) : new Set([kind])
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

  // A project's revision chain — the project rows that share its
  // rootProjectID (see the Project class' rootProjectID/parentRevisionID).
  // Each entry is a project row, so it's absorbed into the shared `projects`
  // cache like `load()`; the ordered chain itself is kept separately here.
  // Always keyed by the chain ROOT, never by the revision asked about, so
  // every member of a chain reads and writes one cache entry.
  const revisionsByProject = ref({})
  const revisionsFor = computed(
    () => projectId => revisionsByProject.value[String(rootIdFor(projectId))] || [],
  )

  // The chain root for any project in it; originals fall back to their own ID.
  function rootIdFor(projectId) {
    return findById.value(projectId)?.rootProjectID || projectId
  }

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
    // headsOnly: one row per revision chain — this is the section-wide seed
    // (sidebar tree, preload cache), not a chain-detail fetch. Non-head
    // revisions still enter the cache on demand via fetchProject/listRevisions.
    loading = $SystemAPI
      .projectList({ limit: 500, sort: 'createdAt DESC', headsOnly: true })
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
  // the compose namespace and adds the creator as a developer. Name and
  // description only: the AI Act deployer-category questions belong to the
  // Govern tab's 'fria-determination' step (see config/pipeline.js and
  // components/wizard/steps/FriaDeterminationStep.vue), so nothing here
  // populates config.deployerCategories. Wiring that step's answers to the
  // backend fields (ProjectDeployerCategories, FriaRequired) is a separate
  // backend slice.
  async function create({ name, description = '' } = {}) {
    const raw = await $SystemAPI.projectCreate({
      status: 'draft',
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
  // The endpoint is unguarded (see the governance section below), so we reset
  // the local `'publish'` governance step back to draft ourselves on success:
  // every publish, first or subsequent, needs its own submit → approve cycle.
  // `mappings` is how records reach the new revision: the backend's
  // migrateRecords no-ops on an empty set and publish then soft-deletes the old
  // namespace, so publishing without them silently drops every record in the
  // project. The Publish tab resolves them from the deployment plan and passes
  // them here; the default stays empty only for a first revision, which has no
  // parent to migrate from.
  async function publishProject(id, mappings = []) {
    const p = findById.value(id)
    if (!p) return
    const raw = await $SystemAPI.projectPublish({
      projectID: p.projectID,
      confirm: true,
      mappings,
    })
    touch()
    return absorb(raw)
  }

  // --- publish approval (REAL, PERSISTED backend state) ----------------------
  // Unlike the per-step governance further down, the revision's OWN approval
  // lives on the project row and is enforced by the server: publish refuses
  // anything that is not approved, refuses a decision from the person who
  // submitted it, and refuses an approval granted against a different version
  // of the revision (it fingerprints the deployment plan). All three actions
  // return the updated project, so the tab re-renders off the same row every
  // other screen reads.
  async function requestPublishApproval(projectId, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectRequestApproval({ projectID: p.projectID, note })
    touch()
    return absorb(raw)
  }

  async function grantPublishApproval(projectId, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectGrantApproval({ projectID: p.projectID, note })
    touch()
    return absorb(raw)
  }

  async function rejectPublishApproval(projectId, note = '') {
    const p = findById.value(projectId)
    if (!p) return
    const raw = await $SystemAPI.projectRejectApproval({ projectID: p.projectID, note })
    touch()
    return absorb(raw)
  }

  // The revision's approval standing, in the vocabulary the rest of the wizard
  // speaks. The server says 'rejected'; every status chip, tag and hint in this
  // section says 'changes-requested' for the same thing, and translating once
  // here is cheaper than teaching all of them a second word for it. An absent
  // value (a project row written before the column existed) reads as 'draft' —
  // unreviewed, which is exactly what it is.
  function publishApprovalStatus(projectId) {
    const status = findById.value(projectId)?.approvalStatus || 'draft'
    return status === 'rejected' ? GOVERNANCE_STATUS_CHANGES_REQUESTED : status
  }

  function publishApprovalNote(projectId) {
    return findById.value(projectId)?.approvalNote || ''
  }

  // Who asked for the approval. The server refuses a decision from this person
  // — two people minimum for anything going live — and the tab needs to say so
  // before the button is pressed rather than after.
  function publishApprovalSubmittedBy(projectId) {
    return findById.value(projectId)?.approvalSubmittedBy || NoID
  }

  // What publishing this draft would change, against its parent revision:
  // per-resource changes (op/kind/name/risk, plus the record count behind a
  // destructive one) and the suggested per-module record mapping. Not cached —
  // it is read once when the Publish tab opens and must reflect edits made
  // seconds earlier. A first revision legitimately has nothing to compare
  // against and comes back empty, not as an error.
  async function deploymentPlan(projectID) {
    const plan = (await $SystemAPI.projectGetDeploymentPlan({ projectID })) || {}
    // Coerced field by field, NOT by spreading over defaults: Go marshals a nil
    // slice as `null`, not `[]`, so a plan with nothing in it arrives as
    // {changes: null} and a spread would happily overwrite the default with it.
    // Every caller treats these as arrays.
    return {
      risk: plan.risk || 'safe',
      changes: plan.changes || [],
      suggestedMappings: plan.suggestedMappings || [],
    }
  }

  // --- revisions --------------------------------------------------------------
  // Revisions are project rows in a chain (root / parent / number), not a
  // field to increment — see the intent doc. `listRevisions` loads the whole
  // chain; `createRevision` branches a new draft off an active project.

  // Load a project's revision chain. Each entry is a project row, so it's
  // absorbed into the shared `projects` cache the same way `load()` absorbs
  // the plain project list; the ordered chain is then cached here.
  async function listRevisions(projectId) {
    const { set = [] } = await $SystemAPI.projectListRevisions({ projectID: projectId })
    const chain = set.map(raw => absorb(raw))
    // Key by the root, not by what was asked for: the chain is one thing, and
    // asking about revision 2 must not shadow the entry written for revision 1.
    const rootId = chain[0]?.rootProjectID || rootIdFor(projectId)
    revisionsByProject.value[String(rootId)] = chain
    return chain
  }

  // Branch a new draft revision off an active project. The backend enforces
  // its own rules (parent must be `active`; only one draft allowed per chain)
  // and rejects otherwise — those errors propagate to the caller rather than
  // being swallowed here. Resync the chain cache for the root afterwards,
  // same as other mutations resync via load rather than patching in place.
  async function createRevision(projectId) {
    const raw = await $SystemAPI.projectCreateRevision({ projectID: projectId })
    touch()
    const result = absorb(raw)
    await listRevisions(result.rootProjectID)
    return result
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

  // --- per-step governance (SESSION-LOCAL SCAFFOLDING) -------------------------
  // INTENTIONAL SCAFFOLDING, NOT A BUG: the real governance backend (the
  // SaveGovernanceStep/TransitionGovernanceStep RPCs and the project's
  // `governance` field) has been removed while the governance model is being
  // redesigned. Until that lands, every project's governance state — both
  // step status/note AND form-type steps' (summary, resource-management)
  // working values — lives ONLY here, in memory, keyed by projectID. Nothing
  // here is sent to the API or read from an API response; a fresh page load
  // starts every step at 'draft' with no note. This is deliberate — it lets
  // the FE UX keep iterating quickly — so don't "fix" the lack of
  // persistence.
  //
  // Per-step review only. The revision's OWN publish approval is real server
  // state on the project row — see publishApprovalStatus and the three actions
  // above.
  //
  // ONE cycle for everything: every step runs draft -> submitted -> approved |
  // changes-requested. Nothing is ever approved without having been submitted
  // first, so an approval always answers a request.
  //
  //   - submit: draft/changes-requested -> submitted, note cleared. Rejected
  //     from submitted/approved (there is nothing to ask for).
  //   - approve: submitted -> approved ONLY, note cleared. Rejected from any
  //     other status — this is what stops an Approve button from acting on a
  //     step nobody put up for review.
  //   - request-changes: ANY status -> changes-requested, note set. A reviewer
  //     may flag a step at any point, including one never submitted; the
  //     wizard only stops offering it once the revision is published.
  //   - "any change needs approval again": editing what a step CONTAINS
  //     invalidates its review — see invalidateReview and the
  //     `invalidating(...)` wrappers on the returned actions.
  //
  // The old "auto-clear on resubmit" rule is deliberately gone with this: a
  // flagged step now clears by being fixed and resubmitted through its own
  // cycle, so submitting the revision can no longer wipe a reviewer's note.

  const GOVERNANCE_STATUS_DRAFT = 'draft'
  const GOVERNANCE_STATUS_SUBMITTED = 'submitted'
  const GOVERNANCE_STATUS_APPROVED = 'approved'
  const GOVERNANCE_STATUS_CHANGES_REQUESTED = 'changes-requested'

  // { [projectID]: { [stepKey]: { status, note, values } } } — never persisted.
  const governanceByProject = ref({})

  // All step entries for a project, auto-vivifying the per-project bucket.
  function govSteps(projectId) {
    const key = String(projectId)
    return governanceByProject.value[key] || (governanceByProject.value[key] = {})
  }

  // One step's entry, auto-vivifying a fresh draft entry on first touch —
  // mirrors the removed backend's *ProjectGovernance.Step accessor.
  function govStep(projectId, stepKey) {
    const steps = govSteps(projectId)
    if (!steps[stepKey]) {
      steps[stepKey] = { status: GOVERNANCE_STATUS_DRAFT, note: '', values: {} }
    }
    return steps[stepKey]
  }

  function governanceStatus(projectId, stepKey) {
    return govSteps(projectId)[stepKey]?.status || GOVERNANCE_STATUS_DRAFT
  }

  function governanceNote(projectId, stepKey) {
    return govSteps(projectId)[stepKey]?.note || ''
  }

  function governanceValues(projectId, stepKey) {
    return govSteps(projectId)[stepKey]?.values || {}
  }

  // Whether ANY step currently has changes requested. The Publish tab refuses
  // to offer an approval over an open flag — a reviewer's note on any step is
  // a statement that the revision is not ready, whichever step it landed on.
  function hasFlaggedSteps(projectId) {
    return Object.values(govSteps(projectId)).some(
      step => step.status === GOVERNANCE_STATUS_CHANGES_REQUESTED,
    )
  }

  // "any change needs approval again": a review states that what was there
  // when it was granted is fit to ship, so changing that content retires it.
  // A submitted or approved step drops back to draft (its owner resubmits when
  // it's ready again). A changes-requested step is left alone on purpose: the
  // reviewer's note has to survive the work done to address it, and it is the
  // resubmit that clears it.
  //
  // The revision's OWN approval is deliberately NOT touched here any more. It
  // is server state now, and the server retires it by recomputing the
  // deployment plan's fingerprint at publish time — so nothing on either side
  // has to remember to call an invalidation hook, and an approval cannot be
  // silently kept alive by an edit path that forgot to.
  //
  // Every action that changes a step's content is wrapped with `invalidating`
  // at the return below rather than calling this itself — one list, so the
  // step each action belongs to is stated in one readable place. Loaders and
  // getters must never appear there.
  function invalidateReview(projectId, stepKeys) {
    const steps = govSteps(projectId)
    const stale = s => s === GOVERNANCE_STATUS_SUBMITTED || s === GOVERNANCE_STATUS_APPROVED
    let changed = false

    for (const key of [].concat(stepKeys || [])) {
      const step = steps[key]
      if (!step || !stale(step.status)) continue
      step.status = GOVERNANCE_STATUS_DRAFT
      step.note = ''
      changed = true
    }

    if (changed) touch()
  }

  // Wrap a mutating action so a SUCCESSFUL call retires the review of the
  // step(s) it changed. `stepKeys` is a key, an array of keys, or a function
  // of the action's own arguments for actions whose step depends on what was
  // changed (updateField). Sync actions stay sync — the FRIA scenario
  // mutators return their id to the caller directly.
  function invalidating(fn, stepKeys) {
    return (...args) => {
      const done = () =>
        invalidateReview(args[0], typeof stepKeys === 'function' ? stepKeys(...args) : stepKeys)
      const out = fn(...args)
      if (out instanceof Promise) return out.then(v => (done(), v))
      done()
      return out
    }
  }

  async function saveStepForm(projectId, stepKey, values) {
    const step = govStep(projectId, stepKey)
    step.values = { ...values }
    touch()
  }

  async function transitionStep(projectId, stepKey, action, note = '') {
    const step = govStep(projectId, stepKey)

    switch (action) {
      case 'submit':
        if (
          step.status !== GOVERNANCE_STATUS_DRAFT &&
          step.status !== GOVERNANCE_STATUS_CHANGES_REQUESTED
        ) {
          throw new Error(`Cannot submit "${stepKey}" from status "${step.status}"`)
        }
        step.status = GOVERNANCE_STATUS_SUBMITTED
        step.note = ''
        break
      case 'approve':
        if (step.status !== GOVERNANCE_STATUS_SUBMITTED) {
          throw new Error(`Cannot approve "${stepKey}" from status "${step.status}"`)
        }
        step.status = GOVERNANCE_STATUS_APPROVED
        step.note = ''
        break
      case 'request-changes':
        step.status = GOVERNANCE_STATUS_CHANGES_REQUESTED
        step.note = note
        break
      default:
        throw new Error(`Invalid governance action "${action}" for step "${stepKey}"`)
    }

    touch()
  }

  // --- AI systems -----------------------------------------------------------------
  // An AI system is the EU AI Act's regulated unit (Art. 3(1)): a named set of
  // THIS revision's resources serving one intended purpose, carrying its own
  // Art. 6 risk class. A project routinely holds several, which is exactly why
  // the FRIA cannot hang off the project itself -- every risk scenario names
  // one AI system (config/friaScenario.js's aiSystemID).
  //
  // UNLIKE the governance surface above, this is REAL, PERSISTED backend state
  // (ProjectAiSystem + ProjectAiSystemEntry). It survives reload; nothing here
  // is session-local scaffolding.
  //
  // Scoped to the revision, not the chain root: the backend stores whatever
  // projectID it is given with no root resolution (unlike members), and that is
  // correct -- a revision's resources are its own, so a boundary drawn around
  // them cannot span revisions.
  //
  // No touch(): AI systems are deliberately NOT a resource kind — no graph node,
  // no kinds.js entry, no permission-matrix row — so nothing here changes the
  // resource graph or the effective-access picture that touch() refreshes.
  const aiSystemsByProject = ref({})

  function aiSystemsFor(projectId) {
    return aiSystemsByProject.value[String(projectId)] || []
  }

  function aiSystem(projectId, aiSystemId) {
    return aiSystemsFor(projectId).find(s => s.id === String(aiSystemId)) || null
  }

  async function loadAiSystems(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI
      .projectAiSystemList({ projectID: key, limit: 100 })
      .catch(() => ({ set: [] }))

    aiSystemsByProject.value[key] = set.map(toAiSystem)
    return aiSystemsByProject.value[key]
  }

  // The API shape flattened to what the UI reads. riskClass is a real
  // filterable column; intendedPurpose is prose stored on meta but sent as a
  // TOP-LEVEL param (the BE hook writes it onto meta) -- so it is read from
  // meta and written flat, which is why the two are asymmetric below.
  //
  // Membership rides along inline: the REST payload embeds `entries` from the
  // service's MemberList, so listing systems already carries their resources
  // and no per-row fan-out is needed.
  function toAiSystem(raw) {
    return {
      id: String(raw.projectAiSystemID),
      handle: raw.handle || '',
      name: raw.meta?.short || raw.handle || '',
      description: raw.meta?.description || '',
      intendedPurpose: raw.meta?.intendedPurpose || '',
      riskClass: raw.riskClass || null,
      resourceRefs: (raw.entries || []).map(e => e.resourceRef).filter(Boolean),
    }
  }

  async function createAiSystem(projectId, { name, description = '', handle } = {}) {
    const key = String(projectId)
    // name/description are FLAT params, not nested under meta. The endpoint
    // takes them as genHook params and its beforeCreate writes them onto
    // res.Meta itself; a `meta` object is not in the request struct at all and
    // the generated client drops it silently, so nesting them persisted a row
    // with no name whatsoever.
    const raw = await $SystemAPI.projectAiSystemCreate({
      projectID: key,
      handle: handle || '',
      name: (name || '').trim() || 'Untitled AI system',
      description: description.trim(),
    })

    const created = toAiSystem(raw)
    aiSystemsByProject.value[key] = [...aiSystemsFor(key), created]
    return created.id
  }

  async function updateAiSystem(projectId, aiSystemId, patch = {}) {
    const key = String(projectId)
    const existing = aiSystem(key, aiSystemId)
    if (!existing) return

    const next = { ...existing, ...patch }
    const raw = await $SystemAPI.projectAiSystemUpdate({
      projectID: key,
      projectAiSystemID: String(aiSystemId),
      handle: next.handle,
      riskClass: next.riskClass || '',
      // All flat, none nested under meta — see createAiSystem. The BE takes
      // each as its own param and its hook writes them onto meta; anything
      // sent as `meta` is dropped by the generated client before the request
      // is even built.
      intendedPurpose: next.intendedPurpose || '',
      name: next.name || '',
      description: next.description || '',
    })

    aiSystemsByProject.value[key] = aiSystemsFor(key).map(s =>
      s.id === String(aiSystemId) ? toAiSystem(raw) : s,
    )
  }

  async function removeAiSystem(projectId, aiSystemId) {
    const key = String(projectId)
    await $SystemAPI.projectAiSystemDelete({
      projectID: key,
      projectAiSystemID: String(aiSystemId),
    })
    aiSystemsByProject.value[key] = aiSystemsFor(key).filter(s => s.id !== String(aiSystemId))
  }

  // --- AI system membership -------------------------------------------------------
  // Entries are (aiSystemID, resourceRef) pairs. `resourceRef` is the existing
  // RBAC/envoy reference string (e.g. corteza::compose:module/<id>), NOT a new
  // vocabulary -- so refs stay parseable by machinery that already exists.
  //
  // A ref whose resource has since been deleted is KEPT and rendered as a
  // tombstone rather than swept up: the FRIA claimed to cover that resource, so
  // its removal is material compliance information and a reassessment trigger,
  // not cleanup. Resolving a ref to a live resource is always allowed to fail.
  // Re-read ONE system (and therefore its membership) without refetching the
  // list. There is no entry-list endpoint -- entries only ever arrive embedded
  // in a system payload -- so the single-resource read is the way to resync.
  async function reloadAiSystem(projectId, aiSystemId) {
    const key = String(projectId)
    const raw = await $SystemAPI
      .projectAiSystemRead({ projectID: key, projectAiSystemID: String(aiSystemId) })
      .catch(() => null)
    if (!raw) return null

    const fresh = toAiSystem(raw)
    aiSystemsByProject.value[key] = aiSystemsFor(key).map(s =>
      s.id === String(aiSystemId) ? fresh : s,
    )
    return fresh
  }

  async function addAiSystemResource(projectId, aiSystemId, resourceRef) {
    const key = String(projectId)
    await $SystemAPI.projectAiSystemEntryAdd({
      projectID: key,
      projectAiSystemID: String(aiSystemId),
      resourceRef,
    })

    aiSystemsByProject.value[key] = aiSystemsFor(key).map(s =>
      s.id === String(aiSystemId)
        ? { ...s, resourceRefs: [...new Set([...(s.resourceRefs || []), resourceRef])] }
        : s,
    )
  }

  async function removeAiSystemResource(projectId, aiSystemId, resourceRef) {
    const key = String(projectId)
    await $SystemAPI.projectAiSystemEntryRemove({
      projectID: key,
      projectAiSystemID: String(aiSystemId),
      resourceRef,
    })

    aiSystemsByProject.value[key] = aiSystemsFor(key).map(s =>
      s.id === String(aiSystemId)
        ? { ...s, resourceRefs: (s.resourceRefs || []).filter(r => r !== resourceRef) }
        : s,
    )
  }

  // --- FRIA risk scenarios --------------------------------------------------------
  // PERSISTED (ProjectFriaScenario), unlike the session-local governance surface
  // above: an Art. 27 assessment has to survive a reload, be approvable against,
  // and give detection rules something durable to reference.
  //
  // Split of concerns on the backend: `aiSystemID`, `title` and `severity` are
  // REAL COLUMNS because a compliance product has to answer questions like
  // "every critical scenario on this AI system"; everything else — the
  // taxonomy key lists and free prose — rides in a meta JSON blob, because
  // those move with EU guidance and must not cost a migration each time.
  //
  // EXPLICIT SAVE, NOT KEYSTROKE WRITES. The editor holds a LOCAL DRAFT
  // (config/friaScenario.js's cloneFriaScenario) while editing; nothing below
  // runs until its Save commits the whole draft in one shot, and Cancel drops
  // the draft with no call at all, so a cancelled create leaves no trace.
  //
  // Still no touch(): scenarios are not resources and carry no graph, kind or
  // effective-access implications.
  const friaScenariosByProject = ref({})

  function friaScenariosFor(projectId) {
    return friaScenariosByProject.value[String(projectId)] || []
  }

  function friaScenario(projectId, scenarioId) {
    return friaScenariosFor(projectId).find(s => s.id === String(scenarioId)) || null
  }

  // API shape -> the flat shape config/friaScenario.js defines and every
  // section component reads. Arrays default to [] rather than undefined: the
  // section components index into them directly.
  function toFriaScenario(raw) {
    const m = raw.meta || {}
    return {
      id: String(raw.projectFriaScenarioID),
      aiSystemID: raw.aiSystemID && raw.aiSystemID !== '0' ? String(raw.aiSystemID) : null,
      title: raw.title || '',
      severity: raw.severity || null,
      description: m.description || '',
      triggerTypes: m.triggerTypes || [],
      triggerDescription: m.triggerDescription || '',
      impactedParties: m.impactedParties || [],
      vulnerableGroups: m.vulnerableGroups || [],
      vulnerableGroupsNotes: m.vulnerableGroupsNotes || '',
      rights: m.rights || [],
      harmVectors: m.harmVectors || [],
      harmVectorsDescription: m.harmVectorsDescription || '',
    }
  }

  // The flat shape -> the API's flat body. Note this sends EVERY field, never
  // a sparse patch: the generated PUT copies its update fields verbatim, so
  // omitting aiSystemID would zero it rather than leave it alone. The editor
  // always holds a full clone, so there is nothing to merge.
  function friaScenarioPayload(s) {
    return {
      aiSystemID: s.aiSystemID || '0',
      title: s.title || '',
      severity: s.severity || '',
      description: s.description || '',
      triggerTypes: s.triggerTypes || [],
      triggerDescription: s.triggerDescription || '',
      impactedParties: s.impactedParties || [],
      vulnerableGroups: s.vulnerableGroups || [],
      vulnerableGroupsNotes: s.vulnerableGroupsNotes || '',
      rights: s.rights || [],
      harmVectors: s.harmVectors || [],
      harmVectorsDescription: s.harmVectorsDescription || '',
    }
  }

  async function loadFriaScenarios(projectId) {
    const key = String(projectId)
    const { set = [] } = await $SystemAPI
      .projectFriaScenarioList({ projectID: key, limit: 100 })
      .catch(() => ({ set: [] }))

    friaScenariosByProject.value[key] = set.map(toFriaScenario)
    return friaScenariosByProject.value[key]
  }

  // Commit an already-built scenario (the editor's local draft). Takes the
  // full object rather than building one itself — see the header above for why
  // that matters for cancelled creates. The server mints the ID, so the
  // draft's client-side `fria-…` id is discarded here.
  async function createFriaScenario(projectId, scenario) {
    const key = String(projectId)
    const raw = await $SystemAPI.projectFriaScenarioCreate({
      projectID: key,
      ...friaScenarioPayload(scenario),
    })

    const created = toFriaScenario(raw)
    friaScenariosByProject.value[key] = [...friaScenariosFor(key), created]
    return created.id
  }

  async function updateFriaScenario(projectId, scenarioId, patch) {
    const key = String(projectId)
    const existing = friaScenario(key, scenarioId)
    const next = { ...(existing || {}), ...patch }

    const raw = await $SystemAPI.projectFriaScenarioUpdate({
      projectID: key,
      projectFriaScenarioID: String(scenarioId),
      ...friaScenarioPayload(next),
    })

    const saved = toFriaScenario(raw)
    friaScenariosByProject.value[key] = friaScenariosFor(key).map(s =>
      s.id === String(scenarioId) ? saved : s,
    )
    return saved.id
  }

  async function removeFriaScenario(projectId, scenarioId) {
    const key = String(projectId)
    await $SystemAPI.projectFriaScenarioDelete({
      projectID: key,
      projectFriaScenarioID: String(scenarioId),
    })
    friaScenariosByProject.value[key] = friaScenariosFor(key).filter(
      s => s.id !== String(scenarioId),
    )
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
    const wl = governanceValues(projectId, 'resource-management')?.connections
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
        // the page service creates), show "Primary" (capitalized handle) rather
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
  //
  // `missing` and `warnings` are the endpoint's account of what it could NOT
  // draw, and they matter more than the edges: a reference whose target is gone
  // (a branch copy that dropped a TAQ binding, a chatbot left pointing at
  // another revision's agent) otherwise shows up as nothing at all — one fewer
  // line on a canvas nobody counted. Both are `omitempty` on the wire and Go
  // marshals an empty slice as `null`, so they are normalised to arrays here.
  //
  // Their `kind` is a backend resource TYPE (`corteza::compose:module`), not one
  // of this section's node kinds; it is translated once, here, so nothing
  // downstream has to know both vocabularies. An untranslatable type keeps its
  // raw string in `resourceType` and leaves `kind` null — the graph then says
  // "a resource" rather than inventing a name for it.
  const graphRefKind = type => ({
    kind: GRAPH_KIND_BY_RESOURCE_TYPE[type] || null,
    resourceType: type || '',
  })

  async function graph(projectID) {
    const {
      nodes = [],
      edges = [],
      missing,
      warnings,
    } = await $SystemAPI.projectGraph({ projectID })

    return {
      nodes,
      edges: edges.map(e => ({
        source: e.sourceID,
        target: e.targetID,
        reason: e.reason,
      })),
      // A configured reference whose target could not be resolved at all.
      missing: (missing || []).map(m => ({
        sourceID: m.sourceID,
        ...graphRefKind(m.kind),
        targetID: m.targetID || '',
        targetIdent: m.targetIdent || '',
        reason: m.reason || '',
      })),
      // A reference that only resolves at run time (a computed step argument,
      // a scope variable) — not broken, but not checkable here either.
      warnings: (warnings || []).map(w => ({
        sourceID: w.sourceID,
        ...graphRefKind(w.kind),
        reason: w.reason || '',
        path: w.path || '',
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
    requestPublishApproval,
    grantPublishApproval,
    rejectPublishApproval,
    publishApprovalStatus,
    publishApprovalNote,
    publishApprovalSubmittedBy,
    deploymentPlan,
    revisionsFor,
    listRevisions,
    createRevision,
    addMember,
    updateMember,
    removeMember,
    // --- step content --------------------------------------------------------
    // Each of these changes what a step CONTAINS, so a successful call retires
    // that step's review and the revision's (see invalidateReview). The step
    // key is stated here, next to the action, rather than inside it: this list
    // IS the map of which step owns which resource, and a new mutating action
    // that is missing from it silently keeps a stale approval alive.
    // A module create/delete also makes/removes its record page, so it counts
    // against the Pages step too.
    addResource: invalidating(addResource, ['data-model', 'pages']),
    removeResource: invalidating(removeResource, ['data-model', 'pages']),
    updateResource: invalidating(updateResource, 'data-model'),
    connectionLibrary,
    connectionsFor,
    loadConnectionLibrary,
    allowedConnectorIds,
    loadConnections,
    prepareConnection,
    saveConnection: invalidating(saveConnection, 'connections'),
    removeConnection: invalidating(removeConnection, 'connections'),
    automationsFor,
    loadAutomations,
    updateAutomation: invalidating(updateAutomation, 'automations'),
    addAutomation: invalidating(addAutomation, 'automations'),
    removeAutomation: invalidating(removeAutomation, 'automations'),
    agentsFor,
    loadAgents,
    addAgent: invalidating(addAgent, 'agents'),
    updateAgent: invalidating(updateAgent, 'agents'),
    removeAgent: invalidating(removeAgent, 'agents'),
    chatbotsFor,
    loadChatbots,
    addChatbot: invalidating(addChatbot, 'chatbots'),
    updateChatbot: invalidating(updateChatbot, 'chatbots'),
    removeChatbot: invalidating(removeChatbot, 'chatbots'),
    rolesFor,
    loadRoles,
    addRole: invalidating(addRole, 'roles'),
    updateRole: invalidating(updateRole, 'roles'),
    removeRole: invalidating(removeRole, 'roles'),
    effectiveAccess,
    isEffectiveAccessLoading,
    loadEffectiveAccess,
    userEffectiveAccess,
    isUserEffectiveAccessLoading,
    loadUserEffectiveAccess,
    setAccess: invalidating(setAccess, 'permissions'),
    setCapabilityAccess: invalidating(setCapabilityAccess, 'permissions'),
    projectUsersFor,
    loadProjectUsers,
    setProjectUserRole: invalidating(setProjectUserRole, 'users'),
    assignProjectUserRoles: invalidating(assignProjectUserRoles, 'users'),
    removeProjectUser: invalidating(removeProjectUser, 'users'),
    addProjectUser: invalidating(addProjectUser, 'users'),
    pagesFor,
    loadPages,
    addPage: invalidating(addPage, 'pages'),
    updatePage: invalidating(updatePage, 'pages'),
    reorderPages: invalidating(reorderPages, 'pages'),
    removePage: invalidating(removePage, 'pages'),
    // A field's sensitivity level is the Data Sensitivity step's whole
    // subject, and that step classifies fields the Data Model step built — so
    // a sensitivity-only patch retires the classification, not the model.
    updateField: invalidating(updateField, (_projectId, _moduleId, _fieldId, patch = {}) =>
      Object.keys(patch).length === 1 && 'sensitivity' in patch ? 'data-sensitivity' : 'data-model',
    ),
    addField: invalidating(addField, 'data-model'),
    removeField: invalidating(removeField, 'data-model'),
    setFields: invalidating(setFields, 'data-model'),
    saveStepForm: invalidating(saveStepForm, (_projectId, stepKey) => stepKey),
    transitionStep,
    governanceStatus,
    governanceNote,
    governanceValues,
    hasFlaggedSteps,
    aiSystemsFor,
    aiSystem,
    loadAiSystems,
    // Wrapped like every other mutator family: reclassifying a system, or
    // changing which resources fall inside the assessed boundary, invalidates
    // the ai-systems review AND the revision's own approval. Leaving an
    // approval standing after the thing it approved was redefined is exactly
    // what this machinery exists to prevent.
    createAiSystem: invalidating(createAiSystem, 'ai-systems'),
    updateAiSystem: invalidating(updateAiSystem, 'ai-systems'),
    removeAiSystem: invalidating(removeAiSystem, 'ai-systems'),
    reloadAiSystem,
    addAiSystemResource: invalidating(addAiSystemResource, 'ai-systems'),
    removeAiSystemResource: invalidating(removeAiSystemResource, 'ai-systems'),

    friaScenariosFor,
    friaScenario,
    loadFriaScenarios,
    createFriaScenario: invalidating(createFriaScenario, 'fria-scenarios'),
    updateFriaScenario: invalidating(updateFriaScenario, 'fria-scenarios'),
    removeFriaScenario: invalidating(removeFriaScenario, 'fria-scenarios'),
    graph,
    graphVisibleKinds,
    graphKindVisible,
    graphToggleKind,
    graphSoloKind,
    setGraphAllKindsVisible,
  }
})
