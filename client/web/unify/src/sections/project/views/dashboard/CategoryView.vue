<template>
  <!-- Unknown category (bad deep-link / stale nav) — render a muted note only. -->
  <div v-if="!cfg" class="p-4">
    <p class="text-muted-color">{{ $t('project.dashboard.list.empty') }}</p>
  </div>

  <!-- The category page: a title bar + metrics, then the list. The whole page
       scrolls vertically (root is `overflow-y-auto`, not clipped) so the list
       can get a generous fixed height instead of fighting the metrics band for
       leftover space. `min-w-0` still keeps the wide table from pushing the
       page past the screen — its own horizontal scroll stays contained inside
       CResourceList (see `.category-list` below), it just isn't clipped
       vertically anymore. -->
  <div v-else class="flex flex-col h-full min-w-0 overflow-y-auto">
    <!-- Title bar — wizard-style: leading badge + title + description. -->
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <KindIcon :config="cfg.badge" size="xl" />
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">{{ $t(cfg.titleKey) }}</h2>
        <p class="text-sm text-muted-color">{{ $t(cfg.descKey) }}</p>
      </div>
    </header>

    <!-- Metrics — KPI row + compact breakdown grid + trend, fixed above the
         list. Breakdowns flow up to four across on wide screens so the band
         stays shallow and leaves room for the list. Driven by the report
         endpoint (see loadMetrics), so it stays accurate above the list's
         200-row cap; only this band gets skeleton/error states, the list below
         keeps its own store-driven loading. -->
    <div class="shrink-0 p-4 flex flex-col gap-3">
      <!-- First load (or a project/category switch): skeleton rather than a
           zeroed-out band. -->
      <div v-if="metricsLoading" class="flex flex-col gap-3">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div
            v-for="n in 3"
            :key="n"
            class="rounded-xl border border-surface bg-surface px-4 py-3 flex flex-col gap-2"
          >
            <span class="h-3 w-1/2 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
            <span class="h-6 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
          </div>
        </div>
        <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
          <div
            v-for="n in cfg.charts.length"
            :key="n"
            class="rounded-lg border border-surface bg-surface p-4"
          >
            <div class="h-44 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
          </div>
        </div>
        <div class="rounded-lg border border-surface bg-surface p-4">
          <div class="h-44 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
        </div>
      </div>

      <!-- Failed: a report call errored. Replaces the KPIs/donuts/trend with a
           visible retry rather than silently rendering them as "0 items". -->
      <section v-else-if="metricsFailed" class="flex flex-col items-center text-center gap-2 py-8">
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
        <CategoryKpiRow
          :kpis="kpiList"
          :spark="kpiSpark"
          :spark-labels="trend.labels"
          :accent="accentColor"
        />
        <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
          <!-- Severity/risk are ordinal (ranked, not parts-of-a-whole) — an
               ordered bar reads rank directly; status/type stay donuts. -->
          <template v-for="chart in cfg.charts" :key="chart.field">
            <CategoryRankBar
              v-if="chart.variant === 'severity' || chart.variant === 'risk'"
              :title-key="chart.titleKey"
              :data="breakdownFor(chart.variant)"
              :variant="chart.variant"
              :accent="accentColor"
              :height="180"
            />
            <CategoryDonutChart
              v-else
              :title-key="chart.titleKey"
              :data="breakdownFor(chart.variant)"
              :variant="chart.variant"
              :category="category"
              :accent="accentColor"
              :height="180"
            />
          </template>
        </div>
        <CategoryTrendChart
          title-key="project.dashboard.chart.createdOverTime"
          :labels="trend.labels"
          :series="trend.series"
          :accent="accentColor"
          :height="180"
        />
      </template>
    </div>

    <!-- The category's items — a generous fixed-height band (~70vh, see
         `.category-list` below) rather than "whatever's left" after the
         metrics; the list scrolls internally (both axes, via CResourceList's
         own `h-full`) so the wide table never overflows the page. -->
    <div class="category-list shrink-0 px-4 pb-4">
      <CResourceList
        ref="resourceListRef"
        class="h-full"
        primary-key="id"
        :fields="fields"
        :items="visibleItems"
        :filter="filter"
        @update:filter="Object.assign(filter, $event)"
        :sorting="sorting"
        :pagination="pagination"
        :loading="store.loading"
        :action-items="actionItemsFor"
        :translations="{
          searchPlaceholder: $t('project.dashboard.list.searchPlaceholder'),
          noItems: $t('project.dashboard.list.empty'),
        }"
        clickable
        @sort="onSort"
        @row-click="onRowClick"
      >
        <!-- New-item action lives in the list toolbar. -->
        <template #header>
          <Button
            icon="pi pi-plus"
            :label="$t('project.dashboard.newButton', { type: $t(cfg.singularKey) })"
            size="small"
            @click="dialogVisible = true"
          />
        </template>

        <!-- One dynamic #body-<col.key> slot per configured column. -->
        <template
          v-for="col in cfg.columns"
          :key="col.key"
          #[`body-${col.key}`]="{ data }"
        >
          <div v-if="col.kind === 'title'" class="flex flex-col">
            <span class="font-medium text-color">{{ data.title }}</span>
            <span v-if="data.description" class="text-xs text-muted-color truncate max-w-72">
              {{ data.description }}
            </span>
          </div>

          <EventBadge v-else-if="col.kind === 'type'" :value="data[col.key]" variant="type" />

          <EventBadge
            v-else-if="col.kind === 'severity'"
            :value="data[col.key]"
            variant="severity"
          />

          <EventBadge v-else-if="col.kind === 'status'" :value="data[col.key]" variant="status" />

          <RiskPips v-else-if="col.kind === 'risk'" :level="data[col.key]" />

          <UserCell v-else-if="col.kind === 'user'" :name="data[col.key]" />

          <span v-else-if="col.kind === 'date'" class="text-sm text-muted-color">
            {{ formatDate(data[col.key]) }}
          </span>

          <span v-else>{{ data[col.key] || '—' }}</span>
        </template>
      </CResourceList>
    </div>

    <!-- New-item dialog — reuses the per-category GovernanceForm schema. -->
    <NewEventDialog
      v-model:visible="dialogVisible"
      :category="category"
      :schema="cfg.formSchema"
      :user-options="store.ownerOptions"
      :on-create="onCreate"
    />

    <!-- Row-click detail drawer — read-only summary + backlog items; its own
         Edit button opens the edit dialog below without closing the drawer. -->
    <EventDetailDrawer
      v-model:visible="drawerVisible"
      :category="category"
      :record="selectedEvent"
      :user-options="store.ownerOptions"
      @edit="editVisible = true"
    />

    <!-- Edit dialog — same schema, pre-filled from the clicked record;
         Save/Delete persist via the store. Opened directly from the kebab
         menu (skipping the drawer), or from the drawer's own Edit button
         (stacks above it). -->
    <EventDetailDialog
      v-model:visible="editVisible"
      :category="category"
      :record="selectedEvent"
      :user-options="store.ownerOptions"
      :on-save="onUpdate"
      :on-delete="onDelete"
    />
  </div>
</template>

<script setup>
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryRankBar from '@/sections/project/components/dashboard/CategoryRankBar.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import EventDetailDrawer from '@/sections/project/components/dashboard/EventDetailDrawer.vue'
import NewEventDialog from '@/sections/project/components/dashboard/NewEventDialog.vue'
import RiskPips from '@/sections/project/components/dashboard/RiskPips.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { CATEGORY_COLORS, colorFor, orderIndex } from '@/sections/project/config/chartColors'
import { bucketWeekly, trendWindow, weekLabel } from '@/sections/project/config/trend'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { useReportStore } from '@/sections/project/stores/report'
import { components, useConfirmDelete, useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CResourceList } = components

const { t } = useI18n()
const route = useRoute()
const store = useEventsStore()
const backlogStore = useBacklogItemsStore()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// Active category comes straight from the route; cfg is null for unknown keys.
const category = computed(() => route.params.category)
const cfg = computed(() => CATEGORY_CONFIG[category.value] || null)
const accentColor = computed(() => CATEGORY_COLORS[category.value] || '')

// Metrics band (KPIs + donuts + trend) — all from the report endpoint
// (accurate, not subject to the events store's per-category 200-row list
// cap). `metricsLoading`/`metricsFailed` mirror the idiom in Overview.vue;
// `loadSeq` guards against a late response from a superseded project/category
// switch overwriting the current one's data (same pattern as Overview's
// loadAll/AllEventsView's load()).
const reportStore = useReportStore()
const trend = reactive({ labels: [], series: [] })
const reportBreakdowns = reactive({ total: 0, open: 0, overdue: 0, byDim: {} })
const metricsLoading = ref(false)
const metricsFailed = ref(false)
let loadSeq = 0

// One report call per breakdown dimension the category's donut charts need
// (see cfg.charts / config/categories.js — the backend's generic dimension
// keys are status/severity/risk/type/day, or status/type/scope/day for
// review), plus one grand-total call for the KPI trio — count/open/overdue
// are computed server-side (the "open ⇔ not Completed" rule and the
// best-effort DateDue parsing live in system/service/project_report.go).
async function loadReport(pid, key, mySeq) {
  const dims = [...new Set(CATEGORY_CONFIG[key].charts.map(c => c.variant))]
  const [totals, ...results] = await Promise.all([
    reportStore.report(pid, key, { metrics: ['count', 'open', 'overdue'] }),
    ...dims.map(dim => reportStore.report(pid, key, { dimensions: [dim] })),
  ])
  if (mySeq !== loadSeq) return // stale response — a newer switch is in flight

  const byDim = {}
  dims.forEach((dim, i) => {
    byDim[dim] = results[i].map(r => ({
      label: String(r.dimensions?.[dim] || '—'),
      value: Number(r.metrics?.count || 0),
    }))
  })

  const g = (totals && totals[0]?.metrics) || {}
  reportBreakdowns.byDim = byDim
  reportBreakdowns.total = Number(g.count || 0)
  reportBreakdowns.open = Number(g.open || 0)
  reportBreakdowns.overdue = Number(g.overdue || 0)
}

// Created-over-time trend for this category, from the report endpoint.
async function loadTrend(pid, key, mySeq) {
  const { fromISO, toISO, starts } = trendWindow(12)
  const groupBy = CATEGORY_CONFIG[key].trendGroupBy
  const points = await reportStore.trend(pid, key, { from: fromISO, to: toISO, groupBy })
  if (mySeq !== loadSeq) return // stale response

  // Pivot the day×group points into one stacked series per group value
  // (e.g. per severity), ordered canonically and coloured to match its badge.
  const byGroup = new Map()
  for (const p of points) {
    const g = p.group || '—'
    if (!byGroup.has(g)) byGroup.set(g, [])
    byGroup.get(g).push({ date: p.date, value: p.value })
  }
  const keys = [...byGroup.keys()].sort(
    (a, b) => orderIndex(groupBy, a) - orderIndex(groupBy, b) || a.localeCompare(b),
  )
  trend.labels = starts.map(weekLabel)
  trend.series = keys.map(g => ({
    name: g,
    color: colorFor(groupBy, g),
    data: bucketWeekly(byGroup.get(g), starts),
  }))
}

// `silent` refreshes the numbers in place (e.g. after a create) without
// flashing the skeleton.
async function loadMetrics(pid, key, { silent = false } = {}) {
  if (!pid || !key || !CATEGORY_CONFIG[key]) {
    reportBreakdowns.total = 0
    reportBreakdowns.open = 0
    reportBreakdowns.overdue = 0
    reportBreakdowns.byDim = {}
    trend.labels = []
    trend.series = []
    return
  }
  const mySeq = ++loadSeq
  if (!silent) metricsLoading.value = true
  metricsFailed.value = false
  try {
    await Promise.all([loadReport(pid, key, mySeq), loadTrend(pid, key, mySeq)])
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer switch
    console.error('Failed to load category metrics', err)
    metricsFailed.value = true
  } finally {
    if (mySeq === loadSeq) metricsLoading.value = false
  }
}

function retryMetrics() {
  loadMetrics(route.params.projectId, category.value)
}

watch(
  [category, () => route.params.projectId],
  () => loadMetrics(route.params.projectId, category.value),
  { immediate: true },
)

const dialogVisible = ref(false)

// Row-click opens the read-only detail drawer; the kebab's "Edit" opens the
// edit dialog directly (skipping the drawer). `selectedEvent` is the clicked
// row (a store-mapped event) and backs both — cleared alongside
// `drawerVisible`/`editVisible` on category switch.
// The drawer registers with the shared right-sidebar store so it behaves
// like every other right panel: opening one (TAQ config, notifications,
// agent…) closes the rest, and vice versa.
const rightSidebar = useRightSidebarStore()
const drawerVisible = computed({
  get: () => rightSidebar.isOpen('project-event-detail'),
  set: v => (v ? rightSidebar.open('project-event-detail') : rightSidebar.close('project-event-detail')),
})
const editVisible = ref(false)
const selectedEvent = ref(null)

// Leaving the view with the drawer open would strand the shared store's
// activePanel on our name — release it.
onUnmounted(() => rightSidebar.close('project-event-detail'))

function onRowClick({ data }) {
  selectedEvent.value = data
  drawerVisible.value = true
}

// Per-row kebab menu — same mechanism as ProjectList.vue (CResourceList's
// built-in actionItems column/TieredMenu handles the trigger button, popup
// positioning and stopPropagation so it never also fires row-click).
const resourceListRef = ref()
const closeMenu = () => resourceListRef.value?.hideActionsMenu?.()

// "Edit" opens the edit dialog directly — no drawer detour.
function openEdit(row) {
  closeMenu()
  selectedEvent.value = row
  editVisible.value = true
}

// "Delete" confirms first, then reuses the same onDelete handler/toasts the
// dialog's own Delete button calls.
function confirmDeleteRow(row) {
  closeMenu()
  confirmDelete({
    header: t('project.dashboard.event.confirmDelete.header'),
    message: t('project.dashboard.event.confirmDelete.message', {
      name: row.title || t('project.dashboard.event.untitled'),
    }),
    onConfirm: () => onDelete(row.id),
  })
}

const actionItemsFor = row => [
  { label: t('general.label.edit'), icon: 'pi pi-pencil', command: () => openEdit(row) },
  { separator: true },
  {
    label: t('general.label.delete'),
    icon: 'pi pi-trash',
    class: 'text-red-500',
    command: () => confirmDeleteRow(row),
  },
]

const filter = reactive({ query: '' })
const sorting = reactive({ sortBy: 'dateDue', sortDesc: true })
const pagination = reactive({
  limit: 50,
  pageCursor: undefined,
  prevPage: '',
  nextPage: '',
  total: 0,
  page: 1,
})

// Columns → CResourceList fields (all sortable, headers resolved via i18n).
const fields = computed(() =>
  (cfg.value?.columns ?? []).map(c => ({ key: c.key, header: t(c.headerKey), sortable: true })),
)

// KPI cards: the whole trio (total/open/overdue) comes from the report
// endpoint's grand-total call, so none of it is subject to the events
// store's 200-row-per-category fetch cap.
const kpiList = computed(() => {
  if (!cfg.value) return []
  return cfg.value.kpis.map(({ key, labelKey }) => ({ labelKey, value: reportBreakdowns[key] }))
})

// KPI tiles' sparkline — a single weekly-total series derived from the
// already-loaded trend (no new report call): sum each group's stacked value
// at every week index. Empty trend (not loaded yet, or an invalid category)
// ⇒ empty array, so CategoryKpiRow renders its tiles without a strip.
const kpiSpark = computed(() =>
  trend.labels.map((_, i) => trend.series.reduce((sum, s) => sum + (s.data[i] || 0), 0)),
)

// Chart data helper — grouped counts for a report dimension, ordered
// canonically for ranked variants (severity/risk/status) so bars read
// worst→best.
const breakdownFor = variant => {
  const rows = reportBreakdowns.byDim[variant] || []
  if (!variant || variant === 'type') return rows
  return [...rows].sort(
    (a, b) => orderIndex(variant, a.label) - orderIndex(variant, b.label) || a.label.localeCompare(b.label),
  )
}

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}

// Client-side query filter across the item's string values.
const filteredItems = computed(() => {
  const list = cfg.value ? store.byCategory(category.value) : []
  const q = (filter.query || '').trim().toLowerCase()
  if (!q) return list
  return list.filter(item =>
    Object.values(item).some(v => typeof v === 'string' && v.toLowerCase().includes(q)),
  )
})

// Columns whose values are ranked enums (see config/chartColors) rather than
// free text — sort these by canonical rank, not alphabetically, so e.g.
// severity reads Critical…Informational instead of A→Z.
const RANKED_COLUMNS = new Set(['severity', 'risk', 'status'])

// Client-side sort (byCategory is already newest-first as the default order).
const visibleItems = computed(() => {
  const list = [...filteredItems.value]
  const { sortBy, sortDesc } = sorting
  if (sortBy) {
    const ranked = RANKED_COLUMNS.has(sortBy)
    list.sort((a, b) => {
      if (ranked) {
        const ai = orderIndex(sortBy, a[sortBy])
        const bi = orderIndex(sortBy, b[sortBy])
        if (ai !== bi) return sortDesc ? bi - ai : ai - bi
        return 0
      }
      const av = a[sortBy] ?? ''
      const bv = b[sortBy] ?? ''
      if (av < bv) return sortDesc ? 1 : -1
      if (av > bv) return sortDesc ? -1 : 1
      return 0
    })
  }
  return list
})

watch(visibleItems, list => { pagination.total = list.length }, { immediate: true })

const onSort = ({ sortField, sortOrder }) => {
  if (!sortField) return
  sorting.sortBy = sortField
  sorting.sortDesc = sortOrder === -1
}

// Switching categories via the nav: reset the transient list/dialog state so the
// new page starts clean (store data stays live).
watch(category, () => {
  filter.query = ''
  dialogVisible.value = false
  drawerVisible.value = false
  editVisible.value = false
  selectedEvent.value = null
  sorting.sortBy = 'dateDue'
  sorting.sortDesc = true
})

// Create handler — persist via the store (list/KPIs/charts/nav badge all react
// off the returned record) and toast. Passed down to NewEventDialog as its
// `onCreate` prop: the dialog awaits this and only closes (dropping the draft)
// when it resolves truthy, so a failed create keeps the dialog open with the
// user's input intact for a retry. `backlogTitles` is the dialog's queued
// Step-3 widget list (may be empty) — each becomes a backlog item linked to
// the newly created record. The event itself already exists by the time
// these run, so a backlog-create failure only toasts; it never reopens the
// dialog (see NewEventDialog's onSubmit contract).
const onCreate = async (payload, backlogTitles = []) => {
  let event
  try {
    event = await store.add(category.value, payload)
  } catch (err) {
    console.error('Failed to create event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.createFailed'))(err)
    return false
  }
  // The metrics band is report-driven, so the new record isn't in it yet —
  // refresh in place (no skeleton flash).
  loadMetrics(route.params.projectId, category.value, { silent: true })
  $toast.toastSuccess(
    t('project.dashboard.newButton', { type: t(cfg.value.singularKey) }),
    t('project.dashboard.event.toast.created'),
  )
  if (backlogTitles.length) {
    const results = await Promise.allSettled(
      backlogTitles.map(title =>
        backlogStore.add({
          category: category.value,
          eventID: event.id,
          title,
          priority: 'Medium',
          status: 'Open',
        }),
      ),
    )
    const failed = results.find(r => r.status === 'rejected')
    if (failed) {
      console.error('Failed to create backlog item', failed.reason)
      $toast.toastErrorHandler(t('project.dashboard.backlog.toast.createFailed'))(failed.reason)
    }
  }
  return true
}

// Update handler — passed to EventDetailDialog's `onSave` prop; same
// truthy/falsy-or-throw contract as onCreate. store.update() splices a fresh
// mapped object into the store's list rather than mutating the old one in
// place, so `selectedEvent` (captured at row-click/edit-open time) is
// repointed at that fresh row — otherwise the drawer behind the edit dialog
// would keep showing stale values after a save.
const onUpdate = async (id, payload) => {
  try {
    await store.update(category.value, id, payload)
    loadMetrics(route.params.projectId, category.value, { silent: true })
    $toast.toastSuccess(t(cfg.value.singularKey), t('project.dashboard.event.toast.updated'))
    if (selectedEvent.value?.id === String(id)) {
      selectedEvent.value = store.byCategory(category.value).find(e => e.id === String(id)) || selectedEvent.value
    }
    return true
  } catch (err) {
    console.error('Failed to update event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.updateFailed'))(err)
    return false
  }
}

// Delete handler — passed to EventDetailDialog's `onDelete` prop (the dialog
// confirms first) and to the kebab's confirmDeleteRow. If the deleted record
// is the one the drawer is showing, close the drawer along with it.
const onDelete = async id => {
  try {
    await store.remove(category.value, id)
    loadMetrics(route.params.projectId, category.value, { silent: true })
    $toast.toastSuccess(t(cfg.value.singularKey), t('project.dashboard.event.toast.deleted'))
    if (selectedEvent.value?.id === String(id)) {
      drawerVisible.value = false
      selectedEvent.value = null
    }
    return true
  } catch (err) {
    console.error('Failed to delete event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.deleteFailed'))(err)
    return false
  }
}
</script>

<style scoped>
/* Tailwind's scale has no 70vh utility (and arbitrary bracket values are off
   the table here), so the list's fixed height lives in plain CSS. Generous
   enough that the list reads as a real list rather than a sliver of leftover
   space, while still leaving the metrics band and page chrome visible above
   the fold on typical viewports. */
.category-list {
  height: 70vh;
}
</style>
