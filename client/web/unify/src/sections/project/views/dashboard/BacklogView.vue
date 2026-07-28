<template>
  <!-- Backlog: every follow-up work item across every category, as a flat,
       sortable/searchable table (reference demo's Backlog page) — replaces
       the old tag-grouped view now that backlog items are their own linked
       resource (see stores/backlogItems.js) rather than a comma-tag field on
       the category records. Title-bar idiom matches CategoryView. -->
  <div class="flex flex-col h-full min-w-0 overflow-y-auto">
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <KindIcon :config="BADGE" size="xl" />
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.backlog.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.backlog.desc') }}</p>
      </div>
    </header>

    <div class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-4">
      <!-- First load (or a project switch): skeleton rather than a
           zeroed-out page. -->
      <div v-if="store.loading" class="flex flex-col gap-2">
        <span
          v-for="n in 6"
          :key="n"
          class="h-10 rounded-lg bg-emphasis block animate-pulse motion-reduce:animate-none"
        />
      </div>

      <!-- Nothing created yet — a CTA gets the first item created without a
           trip through the (currently empty) table toolbar below. -->
      <CEmptyState v-else-if="!store.items.length">
        <p class="font-medium text-color">{{ $t('project.dashboard.backlog.empty') }}</p>
        <p class="text-xs mt-1">{{ $t('project.dashboard.backlog.emptyHint') }}</p>
        <Button
          icon="pi pi-plus"
          :label="$t('project.dashboard.backlog.newButton')"
          size="small"
          class="mt-3"
          @click="openCreate"
        />
      </CEmptyState>

      <template v-else>
        <!-- KPI trio + chart band — report-endpoint-driven (see loadMetrics),
             own skeleton/failed state independent of the store/list above
             (same idiom as CategoryPanel.vue's metrics band). -->
        <div v-if="reportLoading" class="flex flex-col gap-3">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div
              v-for="n in 3"
              :key="n"
              class="rounded-xl border border-surface bg-surface px-4 py-3 flex flex-col gap-2"
            >
              <span
                class="h-3 w-1/2 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
              />
              <span
                class="h-6 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
            <div v-for="n in 4" :key="n" class="rounded-lg border border-surface bg-surface p-4">
              <div class="h-44 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
            </div>
          </div>
        </div>

        <section v-else-if="reportFailed" class="flex flex-col items-center text-center gap-2 py-8">
          <span
            class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emphasis text-red-500"
          >
            <i class="pi pi-exclamation-triangle text-xl" />
          </span>
          <p class="text-sm font-medium text-color">{{ $t('project.dashboard.overview.error') }}</p>
          <Button
            type="button"
            size="small"
            severity="secondary"
            outlined
            :label="$t('project.dashboard.overview.retry')"
            @click="retryMetrics"
          />
        </section>

        <template v-else>
          <CategoryKpiRow :kpis="kpiList" />

          <!-- Chart band — same shape as CategoryView's, report-endpoint-
               driven now (see loadReport): 'backlog-item' is one of the
               report endpoint's six registered resources, so this is no
               longer a client-side trade against the store's 200-row cap
               either. The created-over-time trend stays omitted — this view
               has no per-item detail page to windows a trend against, same
               as before. -->
          <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
            <CategoryDonutChart
              title-key="project.dashboard.chart.byStatus"
              :data="statusBreakdown"
              variant="status"
              :height="180"
            />
            <CategoryDonutChart
              title-key="project.dashboard.chart.byCategory"
              :data="categoryBreakdown"
              variant="category"
              :height="180"
            />
            <CategoryRankBar
              title-key="project.dashboard.chart.byPriority"
              :data="priorityBreakdown"
              variant="priority"
              :height="180"
            />
          </div>
        </template>

        <div class="backlog-list shrink-0">
          <CResourceList
            class="h-full"
            primary-key="id"
            :fields="fields"
            :items="listItems"
            :filter="filter"
            @update:filter="Object.assign(filter, $event)"
            :sorting="sorting"
            :pagination="pagination"
            :loading="listLoading"
            :translations="{
              searchPlaceholder: $t('project.dashboard.backlog.searchPlaceholder'),
              noItems: $t('project.dashboard.list.empty'),
            }"
            clickable
            @sort="handleSort"
            @row-click="onRowClick"
            @page-change="handlePageChange"
          >
            <!-- There WAS an "Unassigned" filter chip/checkbox in this
                 toolbar — dropped, see the module doc comment's UNASSIGNED
                 FILTER note: a server page can't honestly claim "N
                 unassigned" when it only means "N unassigned on this page". -->
            <template #header>
              <Button
                icon="pi pi-plus"
                :label="$t('project.dashboard.backlog.newButton')"
                size="small"
                @click="openCreate"
              />
            </template>

            <template #body-title="{ data }">
              <span class="font-medium text-color">{{ data.title || '—' }}</span>
            </template>

            <template #body-eventID="{ data }">
              <div v-if="linkedEventFor(data)" class="flex items-center gap-2 min-w-0">
                <KindIcon :config="CATEGORY_CONFIG[data.category]?.badge" size="sm" />
                <span class="truncate max-w-40 text-sm text-color">
                  {{ linkedEventFor(data).title || $t('project.dashboard.event.untitled') }}
                </span>
              </div>
              <span v-else class="font-mono text-xs text-muted-color">#{{ data.eventID }}</span>
            </template>

            <!-- Chain-wide scope (project.intent.md "Dashboards"): this list
                 spans every revision, so each row says which one it belongs
                 to — unassigned items (no revisionID) render distinctly
                 rather than as a blank cell (see useRevisionLabel). -->
            <template #body-revisionID="{ data }">
              <span
                v-if="revisionInfo(data.revisionID).unassigned"
                class="inline-flex items-center gap-1 text-xs text-muted-color italic"
              >
                <i class="pi pi-question-circle" />
                {{ revisionInfo(data.revisionID).label }}
              </span>
              <Tag
                v-else
                :value="revisionInfo(data.revisionID).label"
                :severity="revisionInfo(data.revisionID).severity"
                class="!text-xs"
              />
            </template>

            <template #body-assignee="{ data }">
              <UserCell :name="data.assignee" />
            </template>

            <template #body-priority="{ data }">
              <EventBadge :value="data.priority" variant="priority" />
            </template>

            <template #body-status="{ data }">
              <EventBadge :value="data.status" variant="status" />
            </template>

            <template #body-dateDue="{ data }">
              <span class="text-sm text-muted-color">{{ formatDate(data.dateDue) }}</span>
            </template>
          </CResourceList>
        </div>
      </template>
    </div>

    <!-- Create + row-click edit share one dialog — see BacklogItemDialog.
         `allow-revision-select` offers the revisionID field on create too —
         this view is always chain-wide (no revisionId/board scope exists
         here), same as CategoryPanel's dashboard-mode "+ New". -->
    <BacklogItemDialog
      v-model:visible="dialogVisible"
      :record="selectedItem"
      :user-options="eventsStore.ownerOptions"
      allow-revision-select
      :on-save="onSave"
      :on-delete="onDelete"
    />

    <!-- Row-click detail drawer — read-only summary; its Edit button opens the
         edit dialog above without closing the drawer, and clicking the linked
         event opens that event's edit dialog below. -->
    <BacklogItemDrawer
      v-model:visible="drawerVisible"
      :record="selectedItem"
      @edit="dialogVisible = true"
      @open-event="openLinkedEvent"
    />

    <!-- Linked event's edit dialog — opened from the drawer's linked-event
         row; saves/deletes go through the events store (same handlers idiom
         as CategoryView's). -->
    <EventDetailDialog
      v-if="selectedItem"
      v-model:visible="eventDialogVisible"
      :category="selectedItem.category"
      :record="linkedEventRecord"
      :user-options="eventsStore.ownerOptions"
      :on-save="onEventSave"
      :on-delete="onEventDelete"
    />
  </div>
</template>

<script setup>
// THE LIST is server-paged directly off the backlog-item resource endpoint
// (GET /project-backlog-items/, see the useResourceList call below) via the same
// useResourceList idiom views/ProjectList.vue uses, with `incTotal` for a
// true count — NOT store.items (capped at 200 rows — that store's own load()
// comment).
//
// THE KPI TRIO + CHART BAND above the list are report-endpoint-driven now too
// (see loadMetrics below) — 'backlog-item' is one of the report endpoint's
// six registered resources (server/system/service/project_report.go's
// projectReportSources), so this no longer needs a client-side trade against
// store.items' own 200-row cap either (mirrors CategoryPanel.vue's own move).
// `store` (useBacklogItemsStore) stays only for the actual create/update/
// remove calls (onSave/onDelete below), the top-level loading/empty-state
// gate above, and — via `eventsStore` — for linkedEventFor's lookup into the
// events store.
//
// UNASSIGNED FILTER — DROPPED, not just hidden (see the old template's
// Popover/Chip, now gone): the generated resource filters treat
// `revisionID = 0` as "no constraint", not "unassigned only" (see
// server/store/adapters/rdbms/filters.gen.go's `if f.RevisionID > 0` guard,
// applied identically across every one of these six resources) — there is no
// server-side way to ask for "only unassigned rows" today. Client-filtering
// one fetched PAGE to fake it would misreport a page's own unassigned count
// as the whole backlog's, which is worse than not offering the filter.
import BacklogItemDialog from '@/sections/project/components/dashboard/BacklogItemDialog.vue'
import BacklogItemDrawer from '@/sections/project/components/dashboard/BacklogItemDrawer.vue'
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryRankBar from '@/sections/project/components/dashboard/CategoryRankBar.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { useRevisionLabel } from '@/sections/project/composables/useRevisionLabel'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { CATEGORY_COLORS, orderIndex } from '@/sections/project/config/chartColors'
import { mapBacklogRow, useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { useReportStore } from '@/sections/project/stores/report'
import { components, useResourceList, useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CResourceList, CEmptyState } = components

defineOptions({ name: 'BacklogView' })

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const route = useRoute()
const store = useBacklogItemsStore()
const eventsStore = useEventsStore()
const { revisionInfo } = useRevisionLabel()
const $toast = inject('$toast')

// Title-bar badge — neutral (this page spans every category, so it doesn't
// borrow any single category's colour), same icon-square shape as
// CategoryView's per-category badge.
const BADGE = {
  icon: 'pi pi-th-large',
  bg: 'bg-emphasis',
  ring: 'ring-surface',
  text: 'text-color',
}

// Resolve a backlog item's linked event via the events store (already loaded
// alongside this store by DashboardLayout) — null when the event no longer
// exists (deleted) or hasn't loaded, in which case the column falls back to
// the raw `#<eventID>` reference.
function linkedEventFor(item) {
  return eventsStore.byCategory(item.category).find(e => e.id === String(item.eventID)) || null
}

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}

// --- Server-paged list -------------------------------------------------------
// Same idiom views/ProjectList.vue uses: useResourceList owns filter/sorting/
// pagination state and the actual fetch, calling straight through to
// projectBacklogItemList (with `incTotal`) rather than a new endpoint. This
// view is always chain-wide (no revisionId scope — see the class-level
// comment), so `projectID` is just `route.params.projectId`, read fresh on
// every call so a project switch (this component's instance is reused across
// one, see the route.params.projectId watcher below) picks it up.
const {
  items: listItems,
  loading: listLoading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  fetchItems: refetchList,
  filterList,
} = useResourceList(
  params => {
    const { response, cancel } = $SystemAPI.projectBacklogItemListCancellable({
      ...params,
      projectID: route.params.projectId,
    })
    return {
      cancel,
      response: async () => {
        const result = await response()
        return { ...result, set: (result.set || []).map(mapBacklogRow) }
      },
    }
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'dateDue', sortDesc: true },
    pagination: { limit: 50 },
  },
)

// --- KPI trio + chart band, report-endpoint-driven -------------------------
// This view is always chain-wide (no revisionId scope exists here — see the
// class-level comment), so every call below just takes route.params.projectId
// directly. `reportLoading`/`reportFailed` mirror the idiom in
// CategoryPanel.vue/OverviewPanel.vue; `loadSeq` guards a late response from
// a superseded project switch overwriting the current one's data.
const reportStore = useReportStore()
const reportTotals = reactive({ total: 0, open: 0, overdue: 0 })
const reportBreakdowns = reactive({ status: [], priority: [], category: [] })
const reportLoading = ref(false)
const reportFailed = ref(false)
let loadSeq = 0

// One report call per breakdown dimension the three charts need, plus one
// grand-total call for the KPI trio — count/open/overdue are computed
// server-side (the "open ⇔ not Completed" rule lives in
// system/service/project_report.go, not re-derived here).
async function loadReport(pid, mySeq) {
  const [totals, statusRows, priorityRows, categoryRows] = await Promise.all([
    reportStore.report(pid, 'backlog-item', { metrics: ['count', 'open', 'overdue'] }),
    reportStore.report(pid, 'backlog-item', { dimensions: ['status'] }),
    reportStore.report(pid, 'backlog-item', { dimensions: ['priority'] }),
    reportStore.report(pid, 'backlog-item', { dimensions: ['category'] }),
  ])
  if (mySeq !== loadSeq) return // stale response — a newer project switch is in flight

  const g = (totals && totals[0]?.metrics) || {}
  reportTotals.total = Number(g.count || 0)
  reportTotals.open = Number(g.open || 0)
  reportTotals.overdue = Number(g.overdue || 0)

  reportBreakdowns.status = statusRows.map(r => ({
    label: String(r.dimensions?.status || '—'),
    value: Number(r.metrics?.count || 0),
  }))
  reportBreakdowns.priority = priorityRows.map(r => ({
    label: String(r.dimensions?.priority || '—'),
    value: Number(r.metrics?.count || 0),
  }))
  reportBreakdowns.category = categoryRows.map(r => ({
    key: String(r.dimensions?.category || ''),
    value: Number(r.metrics?.count || 0),
  }))
}

// `silent` refreshes the numbers in place (e.g. after a create/update/delete
// below) without flashing the skeleton — same idiom as CategoryPanel.vue's
// loadMetrics.
async function loadMetrics(pid, { silent = false } = {}) {
  if (!pid) return
  const mySeq = ++loadSeq
  if (!silent) reportLoading.value = true
  reportFailed.value = false
  try {
    await loadReport(pid, mySeq)
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer switch
    console.error('Failed to load backlog metrics', err)
    reportFailed.value = true
  } finally {
    if (mySeq === loadSeq) reportLoading.value = false
  }
}

function retryMetrics() {
  loadMetrics(route.params.projectId)
}

watch(
  () => route.params.projectId,
  pid => loadMetrics(pid),
  { immediate: true },
)

// KPI trio above the table — count/open/overdue straight off the report
// endpoint's own totals (no more client-side open/closed re-derivation).
const kpiList = computed(() => [
  { labelKey: 'project.dashboard.kpi.total', value: reportTotals.total },
  { labelKey: 'project.dashboard.kpi.open', value: reportTotals.open },
  { labelKey: 'project.dashboard.kpi.overdue', value: reportTotals.overdue },
])

// Chart data helper — grouped counts for a report dimension, ordered
// canonically so bars/donuts read consistently (mirrors CategoryPanel.vue's
// own sortBreakdown). No zero-filtering needed: aggregateProjectReport only
// ever emits buckets that actually occurred.
function sortBreakdown(variant, rows) {
  return [...rows].sort(
    (a, b) =>
      orderIndex(variant, a.label) - orderIndex(variant, b.label) || a.label.localeCompare(b.label),
  )
}

const statusBreakdown = computed(() => sortBreakdown('status', reportBreakdowns.status))
const priorityBreakdown = computed(() => sortBreakdown('priority', reportBreakdowns.priority))

// Labels are the translated category titles (colorFor can't key on those, so
// each row pins its colour — see CategoryDonutChart's data prop).
const categoryBreakdown = computed(() =>
  CATEGORY_ORDER.filter(key => reportBreakdowns.category.some(r => r.key === key)).map(key => ({
    label: t(CATEGORY_CONFIG[key].titleKey),
    value: reportBreakdowns.category.find(r => r.key === key)?.value || 0,
    color: CATEGORY_COLORS[key],
    key,
  })),
)

// Completed — a stat, not a chart (ruled 2026-07-28): plain Completed ÷ total
// off the status column, from the SAME count/open metrics the KPI trio above
// already fetched (no extra report call, no client-side re-derivation of the
// open/closed rule). Deliberately the SAME formula
// components/wizard/RevisionCompletenessBar.vue uses for its header stat —
// see CategoryPanel.vue's identical completedPercent comment for the full
// reasoning (a completed-over-time series was ruled out: backlog items carry
// no CompletedDate at all, only status, and the report endpoint can't bucket
// by anything but CreatedAt today).

// `sortable` is an EXPLICIT per-column declaration, verified against
// server/store/adapters/rdbms/rdbms.gen.go's sortableProjectBacklogItemFields()
// map as of 2026-07-28 — NOT a blanket assumption: clicking an unsortable
// header sends an unrecognised sort column and the store rejects the query,
// blanking the whole list. title/status/dateDue are sortable there; assignee
// and priority are not (yet — a concurrent backend change is adding
// sortability to those, but "planned" isn't "true" until verified again
// here). revisionID is the one deliberate exception: marked sortable AHEAD of
// that same concurrent change (needs codegen + a server restart before it
// actually works) because a sortable Revision column is the agreed
// replacement for the removed unassigned-only filter — see this file's
// UNASSIGNED FILTER comment above.
const fields = computed(() => [
  { key: 'title', header: t('project.dashboard.columns.title'), sortable: true },
  { key: 'eventID', header: t('project.dashboard.backlog.columns.linkedEvent'), sortable: false },
  { key: 'revisionID', header: t('project.dashboard.columns.revision'), sortable: true },
  { key: 'assignee', header: t('project.dashboard.backlog.f.assignee'), sortable: true },
  { key: 'priority', header: t('project.dashboard.backlog.f.priority'), sortable: true },
  { key: 'status', header: t('project.dashboard.event.f.status'), sortable: true },
  { key: 'dateDue', header: t('project.dashboard.event.f.dateDue'), sortable: true },
])

const dialogVisible = ref(false)
const selectedItem = ref(null)

// The drawer registers with the shared right-sidebar store so it behaves
// like the app's other right panels (exclusive with TAQ/notifications/the
// event drawer — opening one closes the others).
const rightSidebar = useRightSidebarStore()
const drawerVisible = computed({
  get: () => rightSidebar.isOpen('project-backlog-item-detail'),
  set: v =>
    v
      ? rightSidebar.open('project-backlog-item-detail')
      : rightSidebar.close('project-backlog-item-detail'),
})

// Leaving the view with the drawer open would strand the shared store's
// active panel — close it on unmount.
onUnmounted(() => rightSidebar.close('project-backlog-item-detail'))

function openCreate() {
  selectedItem.value = null
  dialogVisible.value = true
}

function onRowClick({ data }) {
  selectedItem.value = data
  drawerVisible.value = true
}

// Save handler — passed to BacklogItemDialog's `onSave` prop; branches
// create/update off whether `id` is set (see the dialog's own contract
// comment). Same truthy/falsy-or-throw idiom as the event dialogs. Refreshes
// the server-paged list on success — it's no longer store-driven (see the
// module doc comment), so a create/update wouldn't otherwise show up:
// refetchList() (keep the current page) for an in-place edit, filterList()
// (back to page 1) for a create, mirroring CategoryPanel's own onUpdate/
// onCreate split. Also refreshes the report-driven KPI/chart band in place
// (silent — no skeleton flash), same reason: that band no longer reacts to
// store.items automatically now that it's report-driven (see loadMetrics).
const onSave = async (id, payload) => {
  try {
    if (id) {
      await store.update(id, payload)
      if (selectedItem.value?.id === String(id)) {
        selectedItem.value = store.items.find(i => i.id === String(id)) || selectedItem.value
      }
      $toast.toastSuccess(
        t('project.dashboard.backlog.singular'),
        t('project.dashboard.backlog.toast.updated'),
      )
      refetchList()
    } else {
      await store.add(payload)
      $toast.toastSuccess(
        t('project.dashboard.backlog.singular'),
        t('project.dashboard.backlog.toast.created'),
      )
      filterList()
    }
    loadMetrics(route.params.projectId, { silent: true })
    return true
  } catch (err) {
    console.error('Failed to save backlog item', err)
    const key = id
      ? 'project.dashboard.backlog.toast.updateFailed'
      : 'project.dashboard.backlog.toast.createFailed'
    $toast.toastErrorHandler(t(key))(err)
    return false
  }
}

const onDelete = async id => {
  try {
    await store.remove(id)
    if (selectedItem.value?.id === String(id)) {
      drawerVisible.value = false
      selectedItem.value = null
    }
    $toast.toastSuccess(
      t('project.dashboard.backlog.singular'),
      t('project.dashboard.backlog.toast.deleted'),
    )
    refetchList()
    loadMetrics(route.params.projectId, { silent: true })
    return true
  } catch (err) {
    console.error('Failed to delete backlog item', err)
    $toast.toastErrorHandler(t('project.dashboard.backlog.toast.deleteFailed'))(err)
    return false
  }
}

// Linked event's edit dialog — opened from the drawer. `linkedEventRecord` is
// resolved live from the events store so a save repoints it automatically
// (store.update splices a fresh object into the list).
const eventDialogVisible = ref(false)
const linkedEventRecord = computed(() =>
  selectedItem.value ? linkedEventFor(selectedItem.value) : null,
)

function openLinkedEvent() {
  if (!linkedEventRecord.value) return
  eventDialogVisible.value = true
}

// Save/delete for the linked event — events store + toasts, same contract as
// CategoryView's onUpdate/onDelete (no metrics reload: this page's charts
// count backlog items, not events).
const onEventSave = async (id, payload) => {
  const category = selectedItem.value?.category
  try {
    await eventsStore.update(category, id, payload)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[category].singularKey),
      t('project.dashboard.event.toast.updated'),
    )
    return true
  } catch (err) {
    console.error('Failed to update event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.updateFailed'))(err)
    return false
  }
}

const onEventDelete = async id => {
  const category = selectedItem.value?.category
  try {
    await eventsStore.remove(category, id)
    $toast.toastSuccess(
      t(CATEGORY_CONFIG[category].singularKey),
      t('project.dashboard.event.toast.deleted'),
    )
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}

// Switching projects resets transient list/dialog state (mirrors
// CategoryPanel's category-switch watch) and re-runs the server-paged list
// from its first page (filterList) — store data stays live (reloaded by
// DashboardLayout).
watch(
  () => route.params.projectId,
  () => {
    filter.query = ''
    dialogVisible.value = false
    drawerVisible.value = false
    eventDialogVisible.value = false
    selectedItem.value = null
    sorting.sortBy = 'dateDue'
    sorting.sortDesc = true
    filterList()
  },
)
</script>

<style scoped>
/* Same fixed-height idiom as CategoryView's `.category-list` — Tailwind's
   scale has no 70vh utility and arbitrary bracket values are off the table,
   so this lives in plain CSS. */
.backlog-list {
  height: 70vh;
}
</style>
