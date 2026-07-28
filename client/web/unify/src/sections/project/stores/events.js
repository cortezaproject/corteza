import { useProjectsStore } from '@/sections/project/stores/projects'
import { normalizeDates } from '@/sections/project/stores/dateUtils'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// Project-dashboard event store — backed by the real per-project system
// resources (project_incident / project_feature / project_privacy /
// project_task / project_review) via $SystemAPI. Each category maps to its own
// resource; the store loads all five for the active project, flattens them into
// one reactive list, and exposes the same getters the dashboard already relies
// on (byCategory / countByCategory / kpis / breakdown) so the views are
// unchanged. Owner fields are user references (uint64); they are resolved to
// display names via the project users directory on read, and sent back as IDs
// on create.

// Per-category API binding: which $SystemAPI methods to call, the record's id
// field, and the user-reference keys to resolve for display. Exported (as
// EVENT_RESOURCES) so callers that page a category's list directly against
// the resource endpoint (CategoryPanel.vue, via useResourceList — see its own
// module comment) can build the same list/read calls this store uses,
// without a second, drifting copy of the per-category API binding.
export const EVENT_RESOURCES = {
  incident: {
    list: 'projectIncidentList',
    create: 'projectIncidentCreate',
    update: 'projectIncidentUpdate',
    delete: 'projectIncidentDelete',
    read: 'projectIncidentRead',
    idKey: 'incidentID',
    userKeys: ['issueOwner', 'changeOwner', 'changeApprovedBy'],
  },
  feature: {
    list: 'projectFeatureList',
    create: 'projectFeatureCreate',
    update: 'projectFeatureUpdate',
    delete: 'projectFeatureDelete',
    read: 'projectFeatureRead',
    idKey: 'featureID',
    userKeys: ['featureOwner', 'changeOwner', 'changeApprovedBy'],
  },
  privacy: {
    list: 'projectPrivacyList',
    create: 'projectPrivacyCreate',
    update: 'projectPrivacyUpdate',
    delete: 'projectPrivacyDelete',
    read: 'projectPrivacyRead',
    idKey: 'privacyID',
    userKeys: ['requestOwner', 'changeOwner', 'changeApprovedBy'],
  },
  task: {
    list: 'projectTaskList',
    create: 'projectTaskCreate',
    update: 'projectTaskUpdate',
    delete: 'projectTaskDelete',
    read: 'projectTaskRead',
    idKey: 'taskID',
    userKeys: ['owner', 'changeOwner'],
  },
  review: {
    list: 'projectReviewList',
    create: 'projectReviewCreate',
    update: 'projectReviewUpdate',
    delete: 'projectReviewDelete',
    read: 'projectReviewRead',
    idKey: 'reviewID',
    userKeys: ['reviewer', 'approvedBy'],
  },
}
const CATS = EVENT_RESOURCES

const CATEGORIES = Object.keys(CATS)

// The single "open" rule for events/reports: a record is open iff its status
// is not Completed. Shared by this store's KPI trio and the Overview report
// cards (views/dashboard/Overview.vue) — do not re-derive it locally.
export function isOpenStatus(status) {
  return status !== 'Completed'
}

// GovernanceForm's date fields bind PrimeVue date pickers, which surface
// `dateDue`/`completedDate` as JS Date objects; normalizeDates (stores/
// dateUtils.js) serializes them to the plain YYYY-MM-DD strings the backend
// stores (the report endpoint's overdue metric parses them as ISO).
const DATE_KEYS = ['dateDue', 'completedDate']

// Map one raw resource record (from this store's own load()/add()/update(),
// or a category list/read call made directly by another consumer — see
// CategoryPanel.vue's server-paged list) to the loose-bag event shape the
// dashboard reads everywhere: a stable `id`, its `category`, and owner refs
// resolved to display names. Raw json keys already match the UI field keys
// (incidentType, issueOwner, …) so the rest is a straight spread. Module-level
// (not a store action) and exported so it can be called outside this store's
// setup — it only needs the shared user directory, resolved fresh each call
// via useProjectUsersStore() rather than a captured closure.
export function mapEventRow(cat, row) {
  const users = useProjectUsersStore()
  const cfg = CATS[cat]
  const e = { ...row, id: String(row[cfg.idKey] ?? ''), category: cat }
  for (const k of cfg.userKeys) {
    // Keep the raw id under <key>Id so edit flows can round-trip; show name.
    if (row[k]) {
      e[`${k}Id`] = String(row[k])
      e[k] = users.userName(row[k])
    } else {
      e[k] = ''
    }
  }
  return e
}

export const useEventsStore = defineStore('events', () => {
  const $SystemAPI = inject('$SystemAPI')
  const users = useProjectUsersStore()
  const projects = useProjectsStore()

  // Flattened events for the active project (all categories). Keyed loads
  // replace this wholesale so getters always reflect the current project.
  const events = ref([])
  const currentProjectId = ref('')
  const loading = ref(false)

  const mapRow = mapEventRow

  // Load every category for a project, flatten, and cache as the active list.
  // `revisionId` is optional: omitted, the live dashboard's usual call reads
  // every revision in the project (its locked scope is the whole project
  // across revisions); passed, only items filed against that revision come
  // back — the wizard's Manage & Monitor board's scope (one revision, like an
  // issue filed against a milestone). `projectId` is always the chain ROOT in
  // the revision-scoped case (work items are filed against the root project).
  async function load(projectId, revisionId) {
    const pid = String(projectId || '')
    if (!pid) return
    currentProjectId.value = pid
    loading.value = true
    try {
      // User directory + the project's user set (members + assigned) first, so
      // owner refs resolve to names and the owner picker is populated.
      await Promise.all([users.load(), projects.loadProjectUsers(pid).catch(() => {})])
      const results = await Promise.all(
        CATEGORIES.map(cat =>
          $SystemAPI[CATS[cat].list]({
            projectID: pid,
            revisionID: revisionId || undefined,
            limit: 200,
          })
            .then(({ set = [] } = {}) => (set || []).map(row => mapRow(cat, row)))
            .catch(err => {
              console.error(`Failed to load ${cat} events`, err)
              return []
            }),
        ),
      )
      events.value = results.flat()
    } finally {
      loading.value = false
    }
  }

  // Owner-picker options: the project's related users (build members + access
  // users), resolved to display names. Falls back to the full directory if the
  // project sets aren't loaded yet. Shape: [{ label, value }] (value = user ID).
  const ownerOptions = computed(() => {
    const pid = currentProjectId.value
    const ids = new Set()
    for (const m of projects.membersFor(pid)) if (m.userId) ids.add(String(m.userId))
    for (const u of projects.projectUsersFor(pid)) if (u.userId) ids.add(String(u.userId))
    let opts = [...ids].map(id => ({ label: users.userName(id), value: id }))
    if (!opts.length) opts = users.users.map(u => ({ label: u.name, value: u.id }))
    return opts.sort((a, b) => a.label.localeCompare(b.label))
  })

  // Items for one category, newest first (records come back newest-first).
  const byCategory = computed(() => cat => events.value.filter(e => e.category === cat))

  const countByCategory = computed(() => cat => byCategory.value(cat).length)

  // KPI trio for a category. `open` = not Completed; `overdue` = due in the past
  // and not Completed (bad/blank dates are guarded and never count as overdue).
  const kpis = computed(() => cat => {
    const items = byCategory.value(cat)
    const now = new Date()
    let open = 0
    let overdue = 0
    for (const e of items) {
      const isOpen = isOpenStatus(e.status)
      if (isOpen) open++
      const due = e.dateDue ? new Date(e.dateDue) : null
      if (isOpen && due && !Number.isNaN(due.getTime()) && due < now) overdue++
    }
    return { total: items.length, open, overdue }
  })

  // Group a category's items by one field and count each bucket. Missing/blank
  // values collapse into an em-dash bucket. Shaped for the bar charts.
  const breakdown = computed(() => (cat, field) => {
    const counts = new Map()
    for (const e of byCategory.value(cat)) {
      const label = String(e[field] || '—')
      counts.set(label, (counts.get(label) || 0) + 1)
    }
    return [...counts.entries()].map(([label, value]) => ({ label, value }))
  })

  // Create an event via the category's resource. Owner fields in the payload
  // are user IDs (from the picker). `revisionId` is optional: the Manage &
  // Monitor board's quick-add (components/wizard/manage/ManageBoard.vue)
  // passes the open revision so the new card doesn't vanish the instant it's
  // created (items default unassigned otherwise — see load()'s comment). It
  // takes precedence over anything in `payload.revisionID` — a board-scoped
  // create always lands on the board's own revision regardless of what a
  // (normally hidden, see NewEventDialog's allowRevisionSelect) form field
  // would say. When `revisionId` is absent, `payload.revisionID` is used
  // instead — the dashboard's chain-wide "+ New" flow lets the revision field
  // itself (config/eventForm.js REVISION_FIELD) pre-assign the item, INCLUDING
  // explicitly to null/unassigned (its own default), which is why this reads
  // `payload.revisionID` even when falsy rather than requiring it truthy.
  // Returns the mapped event.
  async function add(cat, payload = {}, revisionId) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    const pid = currentProjectId.value
    const body = normalizeDates(
      {
        ...payload,
        projectID: pid,
        revisionID: revisionId || payload.revisionID || undefined,
        status: payload.status || 'Open',
      },
      DATE_KEYS,
    )
    const raw = await $SystemAPI[cfg.create](body)
    const event = mapRow(cat, raw || {})
    events.value.unshift(event)
    return event
  }

  // Update an event via the category's resource. `id` is the record's category
  // id (the mapped event's `id`, e.g. incidentID); payload has the same shape
  // as add() — owner fields as user IDs, dates as Date objects from the
  // picker. Patches the record in place so the list/KPIs/nav badges reflect
  // the change without a refetch.
  async function update(cat, id, payload = {}) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    const body = normalizeDates({ ...payload, [cfg.idKey]: id }, DATE_KEYS)
    const raw = await $SystemAPI[cfg.update](body)
    const event = mapRow(cat, raw || {})
    const idx = events.value.findIndex(e => e.category === cat && e.id === String(id))
    if (idx !== -1) events.value.splice(idx, 1, event)
    else events.value.unshift(event)
    return event
  }

  // Fetch and map a single record fresh off the resource's own `read` call —
  // for a category id that isn't (or isn't known to be) in this store's own
  // capped `events` list (see load()'s 200-row-per-category cap). Used by
  // updateStatus below, and by BoardPanel.vue to resolve a clicked card's
  // full detail when the card came off the board endpoint rather than this
  // store (a board past its own 200-row cache reaches items this store never
  // loaded).
  async function fetchOne(cat, id) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    const raw = await $SystemAPI[cfg.read]({ [cfg.idKey]: id })
    return mapRow(cat, raw || {})
  }

  // Status-only update — the board's drag interaction (BoardPanel.vue,
  // mounted chain-wide by views/dashboard/BoardView.vue and revision-scoped
  // by components/wizard/manage/ManageBoard.vue) moving a card to another
  // column. When the item is already cached (the common case), this patches
  // its status in place FIRST so any other view reading this store's
  // reactive list reflects the move right away, then persists it; a failed
  // push rolls the status back and rethrows so the caller can toast (mirrors
  // stores/projects.js#updateProject's snapshot/mutate/rollback idiom).
  // BoardPanel.vue does NOT rely on this in-place patch for its own
  // rendering though — it now pages the board endpoint directly (past this
  // store's 200-row cap) and keeps its own optimistic column state, rolling
  // that back independently on the same failure.
  //
  // Not cached: fetched fresh via fetchOne rather than silently no-op'ing
  // (the old cache-only behaviour) — a card the board endpoint can show but
  // this store never loaded (beyond its 200-row cap) must still be able to
  // change status. Either way the push carries every OTHER field off the
  // base record unchanged, remapping owner fields back to their raw
  // `<key>Id` — the update endpoint is a full PUT (see
  // server/system/service/project_incident.go#Update, a total field
  // replace, not a merge), so a status-only body would blank out the rest.
  async function updateStatus(cat, id, status) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    const cachedItem = events.value.find(e => e.category === cat && e.id === String(id))
    const prevStatus = cachedItem?.status
    if (cachedItem) cachedItem.status = status
    try {
      const base = cachedItem || (await fetchOne(cat, id))
      const body = { ...base, status }
      for (const k of cfg.userKeys) body[k] = base[`${k}Id`] || null
      const raw = await $SystemAPI[cfg.update](
        normalizeDates({ ...body, [cfg.idKey]: id }, DATE_KEYS),
      )
      const event = mapRow(cat, raw || {})
      const idx = events.value.findIndex(e => e.category === cat && e.id === String(id))
      if (idx !== -1) events.value.splice(idx, 1, event)
      return event
    } catch (err) {
      if (cachedItem) cachedItem.status = prevStatus
      throw err
    }
  }

  // Delete an event via the category's resource, then drop it from the local
  // list so the list/KPIs/nav badges react without a refetch.
  async function remove(cat, id) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    await $SystemAPI[cfg.delete]({ [cfg.idKey]: id })
    const idx = events.value.findIndex(e => e.category === cat && e.id === String(id))
    if (idx !== -1) events.value.splice(idx, 1)
  }

  return {
    events,
    loading,
    currentProjectId,
    load,
    ownerOptions,
    byCategory,
    countByCategory,
    kpis,
    breakdown,
    fetchOne,
    add,
    update,
    updateStatus,
    remove,
    categories: CATEGORIES,
  }
})
