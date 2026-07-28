<template>
  <!-- Real event log, as an activity timeline: who did what, to which
       resource, and — for updates — what it actually became. Wears the
       CategoryPanel layout idiom (fixed badge/title header, fixed filter row,
       list fills the rest and scrolls internally). Route-free: extracted from
       views/dashboard/AllEventsView.vue so the wizard's Manage & Monitor
       "Activity" section (components/wizard/manage/ManageActivity.vue) can
       mount the exact same screen, revision-scoped. Whoever mounts this owns
       the topbar — DashboardLayout.vue chain-wide, the wizard header
       revision-scoped — so there is NO Teleport here, same as CategoryPanel.

       bg-surface is load-bearing: DashboardLayout's content <section> is a
       rounded, bordered panel with NO background of its own. Sibling views get
       theirs from CResourceList's <Card>; this one does not use CResourceList,
       so without an explicit surface the whole list renders transparent.

       Note bg-surface is NOT a Tailwind utility — tailwindcss-primeui ships
       bg-emphasis/border-surface but no bg-surface, and the surface palette has
       no DEFAULT. It is declared by the theme itself (lib/vue/src/composables/
       useTheme.ts, `.bg-surface { background-color: var(--p-content-background) }`),
       which is the same token PrimeVue's Card uses. Don't go looking for it in
       tailwind.config. -->
  <div class="flex flex-col h-full min-h-0 min-w-0 overflow-hidden bg-surface">
    <!-- Title bar — wizard-style leading badge + title + description. -->
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 shrink-0 bg-emphasis ring-surface"
      >
        <i class="pi pi-list text-primary" />
      </span>
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.activity.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.activity.desc') }}</p>
      </div>
    </header>

    <!-- Toolbar — deliberately mirrors CResourceList's header so this view reads
         like every other list in the app: same bar shape, actions and search
         right-aligned, filter button in the same place and same icon-only shape.
         Active filters live on the left as removable chips, so the six selects
         stay folded into the popover instead of eating the viewport. -->
    <div
      class="shrink-0 flex flex-wrap items-center justify-between gap-3 px-4 py-3 border-b border-surface"
    >
      <div class="flex-1 min-w-0 flex flex-wrap items-center gap-1.5">
        <template v-if="activeFilters.length">
          <Chip
            v-for="f in activeFilters"
            :key="f.field"
            :label="f.label"
            removable
            class="text-xs"
            @remove="clearFilter(f.field)"
          />
          <Button
            type="button"
            size="small"
            severity="secondary"
            text
            :label="$t('project.dashboard.activity.clearAll')"
            @click="clearAllFilters"
          />
        </template>
      </div>

      <div class="flex-1 flex items-center justify-end gap-2">
        <!-- Quick time ranges (stock-chart style) — window the list + metrics
             in one click; the filter popover's from/to stay the precise
             instrument, and editing those (or removing a chip) releases the
             active preset. -->
        <TimeRangeSelect :model-value="activeRange" @update:model-value="applyRange" />
        <Button
          type="button"
          icon="pi pi-filter"
          severity="secondary"
          size="small"
          text
          :aria-label="$t('project.dashboard.activity.filters')"
          v-tooltip.top="$t('project.dashboard.activity.filters')"
          @click="toggleFilters"
        />
        <Button
          type="button"
          icon="pi pi-refresh"
          severity="secondary"
          size="small"
          text
          :loading="loading"
          :aria-label="$t('project.dashboard.activity.refresh')"
          v-tooltip.top="$t('project.dashboard.activity.refresh')"
          @click="reload"
        />
        <CInputSearch
          v-model="search"
          size="small"
          class="flex-1 min-w-0 max-w-xl"
          :placeholder="$t('project.dashboard.activity.searchPlaceholder')"
        />
      </div>

      <Popover ref="filterPanel">
        <div class="flex flex-col gap-3 w-72">
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.from') }}
            </span>
            <DatePicker
              v-model="filter.from"
              showTime
              hourFormat="24"
              showButtonBar
              size="small"
              fluid
              @update:modelValue="reload"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.to') }}
            </span>
            <DatePicker
              v-model="filter.to"
              showTime
              hourFormat="24"
              showButtonBar
              size="small"
              fluid
              @update:modelValue="reload"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.actor') }}
            </span>
            <CInputUser v-model="filter.actorID" size="small" @update:modelValue="reload" />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.origin') }}
            </span>
            <Select
              v-model="filter.origin"
              :options="originOptions"
              option-label="label"
              option-value="value"
              filter
              show-clear
              size="small"
              @update:modelValue="reload"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.resource') }}
            </span>
            <Select
              v-model="filter.resource"
              :options="resourceOptions"
              option-label="label"
              option-value="value"
              filter
              show-clear
              size="small"
              @update:modelValue="reload"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-muted-color">
              {{ $t('system.actionlog.list.filter.action') }}
            </span>
            <Select
              v-model="filter.action"
              :options="actionOptions"
              option-label="label"
              option-value="value"
              filter
              show-clear
              size="small"
              @update:modelValue="reload"
            />
          </label>
        </div>
      </Popover>
    </div>

    <!-- Metrics band — event volume + at-a-glance stats over the effective
         range, reflecting the active filters. Hidden when the report is
         unavailable (e.g. the viewer lacks action-log.read), so non-admins just
         see the timeline. -->
    <div
      v-if="!metrics.failed"
      class="shrink-0 border-b border-surface px-4 py-3 flex flex-col items-stretch gap-3"
    >
      <div class="flex justify-start gap-8">
        <div class="flex flex-col">
          <span class="text-xl font-semibold text-color leading-none">{{ metrics.total }}</span>
          <span class="text-xs text-muted-color">
            {{ $t('project.dashboard.activity.metrics.events') }}
          </span>
        </div>
        <div class="flex flex-col">
          <span class="text-xl font-semibold text-color leading-none">{{ metrics.actors }}</span>
          <span class="text-xs text-muted-color">
            {{ $t('project.dashboard.activity.metrics.people') }}
          </span>
        </div>
        <div class="flex flex-col">
          <span
            class="text-xl font-semibold leading-none"
            :class="metrics.errors ? 'text-red-500' : 'text-color'"
          >
            {{ metrics.errors }}
          </span>
          <span class="text-xs text-muted-color">
            {{ $t('project.dashboard.activity.metrics.errors') }}
          </span>
        </div>
      </div>
      <div class="w-full min-w-0">
        <CategoryTrendChart
          bare
          :height="128"
          :labels="metrics.labels"
          :range-labels="metrics.rangeLabels"
          :series="metrics.series"
        />
      </div>
    </div>

    <!-- The timeline — fills the remaining space full-width and owns its
         scroll. No horizontal padding: rows carry their own px, and the
         sticky day headings must span edge to edge without -mx hacks. -->
    <div ref="scroller" class="flex-1 min-h-0 overflow-y-auto pb-4" @scroll.passive="onScroll">
      <!-- First load: skeleton rows rather than a spinner, so the layout does
           not jump when the events land. -->
      <div v-if="loading && !items.length" class="flex flex-col gap-4 px-4 pt-4">
        <div v-for="n in 5" :key="n" class="flex gap-3">
          <div
            class="w-7 h-7 rounded-full bg-emphasis shrink-0 animate-pulse motion-reduce:animate-none"
          />
          <div class="flex-1 flex flex-col gap-2 pt-1">
            <div class="h-3 w-2/5 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
            <div class="h-2 w-1/5 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
          </div>
        </div>
      </div>

      <!-- Empty. This view is empty far more often than most, and for reasons
           that are not obvious (only create/update/delete are recorded, nothing
           is backfilled), so say why instead of leaving a dead end. -->
      <div v-else-if="!groups.length" class="flex flex-col items-center text-center gap-2 py-12">
        <span
          class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emphasis text-muted-color"
        >
          <i class="pi pi-inbox text-xl" />
        </span>
        <p class="text-sm font-medium text-color">
          {{
            hasQuery
              ? $t('project.dashboard.activity.emptyFiltered')
              : $t('project.dashboard.activity.empty')
          }}
        </p>
        <p class="text-sm text-muted-color max-w-md">
          {{
            hasQuery
              ? $t('project.dashboard.activity.emptyFilteredHint')
              : $t('project.dashboard.activity.emptyHint')
          }}
        </p>
        <Button
          v-if="hasQuery"
          type="button"
          size="small"
          severity="secondary"
          outlined
          :label="$t('project.dashboard.activity.clearAll')"
          @click="clearAllFilters"
        />
      </div>

      <!-- Day groups. Sticky headers keep "which day am I in?" answered while
           scrolling a long log. -->
      <div v-else class="flex flex-col">
        <section v-for="group in groups" :key="group.key">
          <!-- Must carry the SAME surface as the view root, or rows show through
               it while scrolling underneath. -->
          <!-- bg-emphasis (--p-content-hover-background: surface.100 / surface.800)
               reads as a real band against the view's bg-surface, so a pinned day
               is visible rather than text floating over the rows. It must stay a
               solid token: a translucent one would let rows show through while
               they scroll underneath. -->
          <h3
            class="sticky top-0 z-10 px-4 py-1.5 bg-emphasis text-xs font-semibold uppercase tracking-wide text-muted-color"
          >
            {{ group.label }}
          </h3>
          <div class="pt-2 pb-2 px-2 flex flex-col gap-1">
            <EventTimelineItem
              v-for="item in group.items"
              :key="item.actionID"
              :data="item"
              :actor-name="actorLabel(item)"
              :active-resource="filter.resource"
              @filter="drillDown"
            />
          </div>
        </section>

        <!-- Infinite scroll, with an explicit fallback so the log is still
             reachable when the scroll handler cannot fire (short viewport). -->
        <div class="flex justify-center py-2">
          <Button
            v-if="!exhausted"
            type="button"
            :label="$t('system.actionlog.list.loadOlder')"
            :loading="loading"
            severity="secondary"
            text
            size="small"
            @click="loadOlder"
          />
          <span v-else class="text-xs text-muted-color">
            {{ $t('project.dashboard.activity.end') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, inject, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters, useUserStore } from '@planetcrust/human-vue'
// Cross-section reuse: the ActionLog vocab (resource/action/origin option lists
// and label helpers) is the single source of truth — do NOT duplicate it here.
import {
  RESOURCE_TYPES,
  COMMON_ACTIONS,
  ORIGINS,
  actionLabel,
  actionVerb,
  originLabel,
  resourceLabel,
  resourceTypeLabel,
} from '@/sections/admin/views/system/ActionLog/vocab'
import EventTimelineItem from '@/sections/project/components/dashboard/EventTimelineItem.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import TimeRangeSelect from '@/sections/project/components/dashboard/TimeRangeSelect.vue'
import { useEventActivity } from '@/sections/project/composables/useEventActivity'
import { RANGES, rangeFrom } from '@/sections/project/config/trend'

const { CInputUser, CInputSearch } = components
const { locDate } = filters

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

// Shared user store: batch-resolve actor IDs and look up labels (no bespoke cache).
const userStore = useUserStore()

const props = defineProps({
  // Chain-root project id — scopes the log to this project's events (see
  // buildParams). The backend tags each event with the project owning the
  // affected resource (Action.ResourceProjectID), so every project-scoped
  // resource touched here (compose, agents, chatbots, roles, users…) is
  // included. views/dashboard/AllEventsView.vue passes route.params.projectId
  // (a chain HEAD, which is the chain ROOT for an original, never-revised
  // project); components/wizard/manage/ManageActivity.vue passes
  // project.rootProjectID.
  //
  // Caveats, by design:
  // - List/search actions carry no single resource, so they are NOT attributed
  //   and will not appear here.
  // - Events recorded before resource attribution shipped are not backfilled.
  // - Reading the log needs the global `action-log.read` permission, so this is
  //   admin-only until RBAC becomes scope-aware.
  projectId: { type: [String, Number], default: '' },
  // The OPEN revision (a lib system.Project instance's projectID) — absent
  // (null, the default) means the dashboard's chain-wide scope; set, it's the
  // wizard's Manage & Monitor single-revision scope (project.intent.md
  // "Dashboards" / wizard.intent.md). Forwarded to the actionlog list/report
  // calls as `revisionID` below — NOTE: the actionlog endpoints do not accept
  // that filter yet (it is landing separately, server-side); until it does,
  // this prop is wired through but has no effect and the revision-scoped
  // panel shows the SAME chain-wide events as the dashboard.
  //
  // Even once the filter lands, this view will not show everything for a
  // revision: work items (incidents, features, privacy requests, tasks,
  // reviews) are filed against the chain ROOT, not a revision, so their audit
  // events are not revision-attributable and will keep appearing only
  // chain-wide. Only activity on this revision's OWN build artifacts
  // (namespaces, workflows, agents, etc. — resources that belong to a
  // revision's own project row) will actually narrow.
  revisionId: { type: [String, Number], default: null },
})

const items = ref([])
const loading = ref(false)
const exhausted = ref(false)
const search = ref('')
const filterPanel = ref(null)
const scroller = ref(null)

const filter = reactive({
  from: null,
  to: null,
  resource: '',
  action: '',
  origin: '',
  actorID: '',
})

// --- Quick time-range presets ------------------------------------------------
// RANGES/rangeFrom are shared (config/trend.js) with the Overview/CategoryPanel
// trend charts; this view's own logic is just the filter-window wiring below:
// applying a preset sets filter.from (to stays open = "until now"), and
// editing the window by hand releases the active preset (see the watch below).
const activeRange = ref('all')
// applyRange edits filter.from/to itself; this flag keeps the release-watcher
// below from immediately clearing the preset it just set.
let applyingRange = false

function applyRange(key) {
  const r = RANGES.find(x => x.key === key)
  if (!r) return
  activeRange.value = key
  applyingRange = true
  filter.from = rangeFrom(r)
  filter.to = null
  applyingRange = false
  reload()
}

// Manually editing the window (filter popover, chip removal, clear-all) means
// the preset no longer describes it — release the selection.
watch(
  () => [filter.from, filter.to],
  () => {
    if (!applyingRange) activeRange.value = null
  },
)

const PAGE_SIZE = 50

function toggleFilters(e) {
  filterPanel.value?.toggle(e)
}

// Authoritative option lists from server enums; actions also merge any unusual
// values present in loaded items so service-specific names stay discoverable.
// If user drills down by clicking a row, filter.resource may carry an ID
// suffix that's not in RESOURCE_TYPES. Append it as a virtual option so the
// Select can render it instead of going blank with only the clear button.
const resourceOptions = computed(() => {
  const v = filter.resource
  if (!v) return RESOURCE_TYPES
  if (RESOURCE_TYPES.some(o => o.value === v)) return RESOURCE_TYPES
  return [...RESOURCE_TYPES, { value: v, label: resourceLabel(v) }]
})

const originOptions = ORIGINS

const actionOptions = computed(() => {
  const known = new Set(COMMON_ACTIONS.map(o => o.value))
  const extras = [...new Set(items.value.map(i => i.action).filter(v => v && !known.has(v)))].map(
    value => ({ value, label: actionLabel(value) }),
  )
  return [...COMMON_ACTIONS, ...extras].sort((a, b) => a.label.localeCompare(b.label))
})

// Active filters as chips. `from`/`to` are Dates; the rest are plain values.
const activeFilters = computed(() => {
  const out = []
  if (filter.from)
    out.push({
      field: 'from',
      label: `${t('system.actionlog.list.filter.from')}: ${locDate(filter.from)}`,
    })
  if (filter.to)
    out.push({
      field: 'to',
      label: `${t('system.actionlog.list.filter.to')}: ${locDate(filter.to)}`,
    })
  if (filter.resource) out.push({ field: 'resource', label: resourceTypeLabel(filter.resource) })
  if (filter.action) out.push({ field: 'action', label: actionLabel(filter.action) })
  if (filter.origin) out.push({ field: 'origin', label: originLabel(filter.origin) })
  if (filter.actorID) {
    const u = userStore.findByID(filter.actorID)
    out.push({
      field: 'actorID',
      label: u ? u.name || u.handle || u.username || u.email || filter.actorID : filter.actorID,
    })
  }
  return out
})

const hasQuery = computed(() => activeFilters.value.length > 0 || !!search.value.trim())

function clearFilter(field) {
  filter[field] = field === 'from' || field === 'to' ? null : ''
  reload()
}

function clearAllFilters() {
  filter.from = null
  filter.to = null
  filter.resource = ''
  filter.action = ''
  filter.origin = ''
  filter.actorID = ''
  search.value = ''
  reload()
}

// Click a value to filter by it; click the active filter value to clear it.
function drillDown(field, value) {
  filter[field] = filter[field] === value ? '' : value
  reload()
}

// Free-text search is client-side over what's loaded: the backend has no search
// param, and adding one is a server change. Scoped to the fields you'd actually
// search by — actor, resource, action, description.
const visible = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return items.value
  return items.value.filter(e =>
    [
      actorLabel(e),
      resourceTypeLabel(e.resource),
      e.resource,
      actionVerb(e.action),
      e.action,
      e.description,
    ]
      .filter(Boolean)
      .some(v => String(v).toLowerCase().includes(q)),
  )
})

// Group by calendar day, newest first (the API already returns newest-first).
const groups = computed(() => {
  const out = []
  let current = null

  for (const e of visible.value) {
    const d = e.timestamp ? new Date(e.timestamp) : null
    const key = d && !Number.isNaN(d.getTime()) ? d.toDateString() : 'unknown'

    if (!current || current.key !== key) {
      current = { key, label: dayLabel(d), items: [] }
      out.push(current)
    }
    current.items.push(e)
  }
  return out
})

function dayLabel(d) {
  if (!d || Number.isNaN(d.getTime())) return t('project.dashboard.activity.unknownDate')

  const today = new Date()
  const yesterday = new Date()
  yesterday.setDate(today.getDate() - 1)

  if (d.toDateString() === today.toDateString()) return t('project.dashboard.activity.today')
  if (d.toDateString() === yesterday.toDateString())
    return t('project.dashboard.activity.yesterday')
  return locDate(d)
}

function buildParams(beforeActionID) {
  return {
    from: filter.from ? filter.from.toISOString() : undefined,
    to: filter.to ? filter.to.toISOString() : undefined,
    resource: filter.resource || undefined,
    // Scope to this project so the log shows only its events. Filters on the
    // project owning the affected resource (rel_resource_project), NOT on the
    // request scope — see the projectId prop comment above.
    resourceProjectID: props.projectId || undefined,
    // Narrows to one revision once the backend filter lands (see the
    // revisionId prop comment above) — forwarded unconditionally so this
    // panel needs no further change when it does.
    revisionID: props.revisionId ? String(props.revisionId) : undefined,
    action: filter.action || undefined,
    origin: filter.origin || undefined,
    // actorID can be null when cleared via Select's clear button
    actorID: filter.actorID ? [filter.actorID] : undefined,
    beforeActionID: beforeActionID || undefined,
    limit: PAGE_SIZE,
  }
}

// Sequence token: discard responses from filter changes that have been superseded.
let loadSeq = 0

async function load(reset = false) {
  // Allow reset to bypass the in-flight guard (filter change should always reload).
  if (loading.value && !reset) return

  if (reset) {
    items.value = []
    exhausted.value = false
  }

  const beforeActionID = reset
    ? undefined
    : items.value.length > 0
      ? items.value[items.value.length - 1].actionID
      : undefined

  loading.value = true
  const mySeq = ++loadSeq
  try {
    const { set } = await $SystemAPI.actionlogList(buildParams(beforeActionID))
    if (mySeq !== loadSeq) return // stale response

    const batch = set ?? []
    // A short page means the server has nothing older left; stop asking.
    exhausted.value = batch.length < PAGE_SIZE

    if (reset) {
      items.value = batch
    } else {
      items.value = [...items.value, ...batch]
    }
    resolveActors(batch)
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

function reload() {
  load(true)
  loadMetrics()
}

// Event metrics over the effective range, reflecting the active filters. The
// report requires an explicit window, so an unset from/to defaults to 30 days.
// The actual report calls + day/week bucketing live in useEventActivity,
// shared with the Overview activity band — this just supplies the filters and
// owns the fail-soft/failed-band behaviour specific to this view.
const activity = useEventActivity()
const metrics = reactive({
  labels: [],
  rangeLabels: [],
  series: [],
  total: 0,
  actors: 0,
  errors: 0,
  failed: false,
})

function metricsRange() {
  const to = filter.to || new Date()
  const from = filter.from || new Date(to.getTime() - 30 * 86400000)
  return { from, to }
}

async function loadMetrics() {
  const { from, to } = metricsRange()
  try {
    const m = await activity.loadMetrics(props.projectId, {
      from,
      to,
      resource: filter.resource,
      action: filter.action,
      origin: filter.origin,
      actorID: filter.actorID,
      revisionId: props.revisionId,
    })
    Object.assign(metrics, m)
    metrics.failed = false
  } catch (e) {
    // Fail soft: hide the band (e.g. the viewer lacks action-log.read) rather
    // than showing a broken zero-strip.
    console.error('Failed to load event metrics', e)
    metrics.failed = true
  }
}

function loadOlder() {
  load(false)
}

// Pull the next page as the bottom comes into view.
function onScroll() {
  const el = scroller.value
  if (!el || loading.value || exhausted.value) return
  if (el.scrollHeight - el.scrollTop - el.clientHeight < 240) loadOlder()
}

// Actor resolution via the shared user store — batch-resolve encountered IDs.
function resolveActors(rows) {
  const ids = rows.map(r => r.actorID).filter(id => id && id !== '0')
  // Fire-and-forget; resolveUsers re-throws on API failure, so swallow it — a
  // failed name lookup must not surface as an unhandled rejection.
  userStore.resolveUsers(ids).catch(() => {})
}

function actorLabel(data) {
  const user = data.actorID && userStore.findByID(data.actorID)
  if (user) return user.name || user.handle || user.username || user.email || user.userID
  return data.actor || data.actorID || ''
}

// Initial load, and reload on a project/revision switch (e.g. the wizard's
// revision switcher, or a dashboard route param change) — `immediate: true`
// covers the mount case, mirroring CategoryPanel/OverviewPanel's own
// (projectId, revisionId) watchers rather than a separate onMounted.
watch(
  () => [props.projectId, props.revisionId],
  () => {
    load(true)
    loadMetrics()
  },
  { immediate: true },
)
</script>
