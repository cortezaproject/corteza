import { isOpenStatus } from '@/sections/project/stores/events'
import { normalizeDates } from '@/sections/project/stores/dateUtils'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// Project-dashboard backlog store — backed by the ProjectBacklogItem system
// resource (projectBacklogItem{List,Create,Update,Delete} via $SystemAPI).
// Each item is a "sub-issue" linked to one category event (category +
// eventID). Mirrors stores/events.js's shape/idiom: a flat reactive list for
// the active project, assignee resolved to a display name on read and sent
// back as an ID on write, dates normalized via the shared helper.

const DATE_KEYS = ['dateDue']

// Map one raw backlog-item record (from this store's own load()/add()/
// update(), or a list/read call made directly by another consumer — see
// BacklogView.vue's server-paged list) to the loose-bag shape the dashboard
// reads everywhere: a stable `id`, and assignee resolved to a display name
// (raw id kept under `assigneeId` so edit flows can round-trip it).
// Module-level (not a store action) and exported, same reasoning as
// stores/events.js#mapEventRow — resolves the shared user directory fresh
// each call via useProjectUsersStore() rather than a captured closure.
export function mapBacklogRow(row) {
  const users = useProjectUsersStore()
  const item = { ...row, id: String(row.backlogItemID ?? '') }
  if (row.assignee) {
    item.assigneeId = String(row.assignee)
    item.assignee = users.userName(row.assignee)
  } else {
    item.assignee = ''
  }
  return item
}

export const useBacklogItemsStore = defineStore('backlog-items', () => {
  const $SystemAPI = inject('$SystemAPI')
  const users = useProjectUsersStore()

  const items = ref([])
  const currentProjectId = ref('')
  const loading = ref(false)

  const mapRow = mapBacklogRow

  // Load every backlog item for a project. Same 200-row cap idiom as
  // stores/events.js#load — acceptable for v1, no fix needed. `revisionId` is
  // optional and mirrors events.js#load exactly: omitted, every revision in
  // the project comes back (the live dashboard's whole-project scope);
  // passed, only items filed against that one revision do (the wizard's
  // Manage & Monitor board's scope). `projectId` is always the chain ROOT in
  // the revision-scoped case (work items are filed against the root project).
  async function load(projectId, revisionId) {
    const pid = String(projectId || '')
    if (!pid) return
    currentProjectId.value = pid
    loading.value = true
    try {
      await users.load()
      const { set = [] } = await $SystemAPI.projectBacklogItemList({
        projectID: pid,
        revisionID: revisionId || undefined,
        limit: 200,
      })
      items.value = (set || []).map(mapRow)
    } catch (err) {
      console.error('Failed to load backlog items', err)
      items.value = []
    } finally {
      loading.value = false
    }
  }

  // Items linked to one category event, in load order.
  const byEvent = computed(() => (category, eventID) => {
    const eid = String(eventID ?? '')
    return items.value.filter(i => i.category === category && String(i.eventID) === eid)
  })

  // Open (not Completed) items across the whole project — feeds the Backlog
  // nav item's live badge (see DashboardNav.vue#badgeValue).
  const openCount = computed(() => items.value.filter(i => isOpenStatus(i.status)).length)

  // Create a backlog item. `payload.assignee` is a user ID (or blank/null for
  // unassigned); dates are Date objects from the picker. `revisionId` is
  // optional, same contract as events.js#add — the Manage & Monitor board's
  // quick-add (ManageBoard.vue) passes the open revision so items queued via
  // NewEventDialog's inline backlog widget land assigned too, taking
  // precedence over `payload.revisionID` when set. When absent,
  // `payload.revisionID` is used instead — BacklogItemDialog's own
  // revisionID field (config/eventForm.js REVISION_FIELD, shown whenever
  // `allowRevisionSelect` is set — see BacklogView.vue/EventDetailDrawer.vue)
  // lets the dashboard's chain-wide create flows pre-assign the item,
  // including explicitly to null/unassigned, hence reading it even when
  // falsy. Returns the mapped item and prepends it to the local list so
  // callers relying on the flat list (byEvent/openCount) react without a
  // refetch.
  async function add(payload = {}, revisionId) {
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
    const raw = await $SystemAPI.projectBacklogItemCreate(body)
    const item = mapRow(raw || {})
    items.value.unshift(item)
    return item
  }

  // Update a backlog item. `id` is the record's backlogItemID.
  async function update(id, payload = {}) {
    const body = normalizeDates({ ...payload, backlogItemID: id }, DATE_KEYS)
    const raw = await $SystemAPI.projectBacklogItemUpdate(body)
    const item = mapRow(raw || {})
    const idx = items.value.findIndex(i => i.id === String(id))
    if (idx !== -1) items.value.splice(idx, 1, item)
    else items.value.unshift(item)
    return item
  }

  // Fetch and map a single record fresh off the resource's own `read` call —
  // mirrors stores/events.js#fetchOne exactly, for a backlog item id that
  // isn't (or isn't known to be) in this store's own capped `items` list
  // (see load()'s 200-row cap). Used by updateStatus below, and by
  // BoardPanel.vue to resolve a clicked backlog card's full detail when it
  // came off the board endpoint rather than this store.
  async function fetchOne(id) {
    const raw = await $SystemAPI.projectBacklogItemRead({ backlogItemID: id })
    return mapRow(raw || {})
  }

  // Status-only update — the board's drag interaction (BoardPanel.vue,
  // mounted chain-wide by views/dashboard/BoardView.vue and revision-scoped
  // by components/wizard/manage/ManageBoard.vue) moving a card to another
  // column. Mirrors stores/events.js#updateStatus exactly: when the item is
  // already cached, this patches its status in place FIRST (so any other
  // view reading this store's reactive list reflects the move right away),
  // then persists it; a failed push rolls the status back and rethrows so
  // the caller can toast. BoardPanel.vue does NOT rely on this in-place
  // patch for its own rendering — it pages the board endpoint directly (past
  // this store's 200-row cap) and keeps its own optimistic column state,
  // rolling that back independently on the same failure.
  //
  // Not cached: fetched fresh via fetchOne rather than silently no-op'ing
  // (the old cache-only behaviour) — a card the board endpoint can show but
  // this store never loaded (beyond its 200-row cap) must still be able to
  // change status. Either way the push carries every OTHER field off the
  // base record unchanged, remapping `assignee` back to its raw
  // `assigneeId` — the update endpoint is a full PUT, not a merge, so a
  // status-only body would blank out the rest.
  async function updateStatus(id, status) {
    const cachedItem = items.value.find(i => i.id === String(id))
    const prevStatus = cachedItem?.status
    if (cachedItem) cachedItem.status = status
    try {
      const base = cachedItem || (await fetchOne(id))
      const body = { ...base, status, assignee: base.assigneeId || null }
      const raw = await $SystemAPI.projectBacklogItemUpdate(
        normalizeDates({ ...body, backlogItemID: id }, DATE_KEYS),
      )
      const updated = mapRow(raw || {})
      const idx = items.value.findIndex(i => i.id === String(id))
      if (idx !== -1) items.value.splice(idx, 1, updated)
      return updated
    } catch (err) {
      if (cachedItem) cachedItem.status = prevStatus
      throw err
    }
  }

  // Delete a backlog item, then drop it from the local list so the
  // list/badges react without a refetch.
  async function remove(id) {
    await $SystemAPI.projectBacklogItemDelete({ backlogItemID: id })
    const idx = items.value.findIndex(i => i.id === String(id))
    if (idx !== -1) items.value.splice(idx, 1)
  }

  return {
    items,
    loading,
    currentProjectId,
    load,
    byEvent,
    openCount,
    fetchOne,
    add,
    update,
    updateStatus,
    remove,
  }
})
