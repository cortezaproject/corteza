import { useProjectsStore } from '@/sections/project/stores/projects'
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
// field, and the user-reference keys to resolve for display.
const CATS = {
  incident: {
    list: 'projectIncidentList',
    create: 'projectIncidentCreate',
    idKey: 'incidentID',
    userKeys: ['issueOwner', 'changeOwner', 'changeApprovedBy'],
  },
  feature: {
    list: 'projectFeatureList',
    create: 'projectFeatureCreate',
    idKey: 'featureID',
    userKeys: ['featureOwner', 'changeOwner', 'changeApprovedBy'],
  },
  privacy: {
    list: 'projectPrivacyList',
    create: 'projectPrivacyCreate',
    idKey: 'privacyID',
    userKeys: ['requestOwner', 'changeOwner', 'changeApprovedBy'],
  },
  task: {
    list: 'projectTaskList',
    create: 'projectTaskCreate',
    idKey: 'taskID',
    userKeys: ['owner', 'changeOwner'],
  },
  review: {
    list: 'projectReviewList',
    create: 'projectReviewCreate',
    idKey: 'reviewID',
    userKeys: ['reviewer', 'approvedBy'],
  },
}

const CATEGORIES = Object.keys(CATS)

// Backlog is stored comma-separated server-side; the list renders it as pills.
const splitBacklog = v =>
  String(v || '')
    .split(',')
    .map(s => s.trim())
    .filter(Boolean)
const joinBacklog = v => (Array.isArray(v) ? v.join(',') : String(v || ''))

export const useEventsStore = defineStore('events', () => {
  const $SystemAPI = inject('$SystemAPI')
  const users = useProjectUsersStore()
  const projects = useProjectsStore()

  // Flattened events for the active project (all categories). Keyed loads
  // replace this wholesale so getters always reflect the current project.
  const events = ref([])
  const currentProjectId = ref('')
  const loading = ref(false)

  // Map one raw resource record to the loose-bag event shape the dashboard
  // reads: a stable `id`, its `category`, backlog as an array, and owner refs
  // resolved to display names. Raw json keys already match the UI field keys
  // (incidentType, issueOwner, …) so the rest is a straight spread.
  function mapRow(cat, row) {
    const cfg = CATS[cat]
    const e = { ...row, id: String(row[cfg.idKey] ?? ''), category: cat }
    e.backlog = splitBacklog(row.backlog)
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

  // Load every category for a project, flatten, and cache as the active list.
  async function load(projectId) {
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
          $SystemAPI[CATS[cat].list]({ projectID: pid, limit: 200 })
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
    const p = projects.findById(pid)
    for (const m of p?.members || []) if (m.userId) ids.add(String(m.userId))
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
      const isOpen = e.status !== 'Completed'
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

  // Create an event via the category's resource. Owner fields in the payload are
  // user IDs (from the picker); backlog is an array. Returns the mapped event.
  async function add(cat, payload = {}) {
    const cfg = CATS[cat]
    if (!cfg) throw new Error(`Unknown category: ${cat}`)
    const pid = currentProjectId.value
    const body = {
      ...payload,
      projectID: pid,
      status: payload.status || 'Open',
      backlog: joinBacklog(payload.backlog),
    }
    const raw = await $SystemAPI[cfg.create](body)
    const event = mapRow(cat, raw || {})
    events.value.unshift(event)
    return event
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
    add,
    categories: CATEGORIES,
  }
})
