<template>
  <!-- Unknown category (bad deep-link / stale nav, or a bad prop from a
       future caller) — render a muted note only. -->
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
         stays shallow and leaves room for the list. Chain-wide (no
         `revisionId`): driven by the report endpoint (see loadMetrics), so it
         stays accurate above the events store's 200-row list cap; only this
         band gets skeleton/error states then. Revision-scoped (`revisionId`
         set): the report endpoint has no per-revision filter (see
         server/system/rest/request/projectReport.go), so these are computed
         client-side off the already revision-scoped events/backlog store rows
         instead (see localKpis/localBreakdown/localTrend/localSpark below) —
         no separate load, no failure state, just the store's own `loading`. -->
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
            <span
              class="h-3 w-1/2 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
            />
            <span
              class="h-6 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
            />
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

      <!-- Failed: a report call errored. Chain-wide only — the revision-scoped
           path never sets this (see metricsFailed below). Replaces the
           KPIs/donuts/trend with a visible retry rather than silently
           rendering them as "0 items". -->
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
          :spark="spark.values"
          :spark-labels="spark.labels"
          :accent="accentClass"
        />
        <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
          <!-- Severity/risk are ordinal (ranked, not parts-of-a-whole) — an
               ordered bar reads rank directly; status/type stay donuts. -->
          <template v-for="chart in cfg.charts" :key="chart.field">
            <CategoryRankBar
              v-if="chart.variant === 'severity' || chart.variant === 'risk'"
              :title-key="chart.titleKey"
              :data="breakdownFor(chart)"
              :variant="chart.variant"
              :height="180"
            />
            <CategoryDonutChart
              v-else
              :title-key="chart.titleKey"
              :data="breakdownFor(chart)"
              :variant="chart.variant"
              :category="category"
              :height="180"
            />
          </template>
        </div>
        <CategoryTrendChart
          title-key="project.dashboard.chart.createdOverTime"
          :labels="trend.labels"
          :range-labels="trend.rangeLabels"
          :series="trend.series"
          :legend-variant="cfg.trendGroupBy"
          :height="180"
        >
          <!-- Windows this chart (and the KPI sparkline derived from it) only
               — the KPIs/donuts above stay on the report loads. Chain-wide
               only: a single revision's item count is small (often zero
               pre-assignment), so the revision-scoped path skips the
               selector and stays on the fixed ~6-month adaptiveWindow default
               (mirrors ManageMetrics.vue's own trend chart, same
               rationale). -->
          <template v-if="!isRevisionScoped" #actions>
            <TimeRangeSelect :model-value="trendRange" @update:model-value="onRangeChange" />
          </template>
        </CategoryTrendChart>
      </template>
    </div>

    <!-- The category's items — a generous fixed-height band (~70vh, see
         `.category-list` below) rather than "whatever's left" after the
         metrics; the list scrolls internally (both axes, via CResourceList's
         own `h-full`) so the wide table never overflows the page. `items`
         comes straight off the events store's `byCategory` getter — whichever
         scope the mounting host loaded it with (chain-wide for the dashboard,
         one revision for the wizard — see the module doc comment below), so
         this component never re-filters by revision itself. -->
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
        <!-- New-item action + active-filter chips live in the list toolbar's
             left side (CResourceList toolbar convention: chips left; filter →
             refresh → search right). The Unassigned chip/filter only makes
             sense chain-wide — revision-scoped mode (the wizard board) never
             loads an unassigned row in the first place (its list call is
             scoped to one revisionID), same reason the revisionID column
             itself is hidden there (see `fields` below). -->
        <template #header>
          <div class="flex items-center gap-2 flex-wrap">
            <Button
              icon="pi pi-plus"
              :label="$t('project.dashboard.newButton', { type: $t(cfg.singularKey) })"
              size="small"
              @click="dialogVisible = true"
            />
            <Chip
              v-if="!isRevisionScoped && filter.unassignedOnly"
              :label="$t('project.dashboard.revision.unassigned')"
              removable
              class="text-xs"
              @remove="filter.unassignedOnly = false"
            />
          </div>
        </template>

        <!-- One dynamic #body-<col.key> slot per configured column. -->
        <template v-for="col in cfg.columns" :key="col.key" #[`body-${col.key}`]="{ data }">
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

        <!-- Chain-wide scope only (project.intent.md "Dashboards"): this
             fixed extra column says which revision a row belongs to, for the
             dashboard's whole-chain list (unassigned items render distinctly
             — see useRevisionLabel). Not in `fields` at all when
             `revisionId` is set (the wizard's single-revision scope) — every
             row already shares that one revision there, so the column would
             be redundant on every row instead of informative. -->
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

        <!-- Filter button — same icon-only shape + Popover idiom as every
             other CResourceList consumer (see ProjectList.vue). Chain-wide
             only, same reasoning as the chip above. -->
        <template v-if="!isRevisionScoped" #filter>
          <Button
            type="button"
            icon="pi pi-filter"
            severity="secondary"
            size="small"
            text
            :aria-label="$t('project.dashboard.list.filters')"
            v-tooltip.top="$t('project.dashboard.list.filters')"
            @click="toggleFilterMenu"
          />
        </template>
      </CResourceList>
    </div>

    <Popover v-if="!isRevisionScoped" ref="filterMenu">
      <div class="flex items-center gap-2 p-2 w-56">
        <Checkbox v-model="filter.unassignedOnly" inputId="unassignedOnlyFilter" binary />
        <label for="unassignedOnlyFilter" class="text-sm cursor-pointer">
          {{ $t('project.dashboard.list.unassignedFilter') }}
        </label>
      </div>
    </Popover>

    <!-- New-item dialog — reuses the per-category GovernanceForm schema.
         `allow-revision-select` offers the revisionID field only chain-wide —
         revision-scoped creation (the wizard board) already assigns the new
         item to the open revision implicitly (see onCreate's `props.revisionId`
         below); see NewEventDialog's own prop comment for the full reasoning. -->
    <NewEventDialog
      v-model:visible="dialogVisible"
      :category="category"
      :schema="cfg.formSchema"
      :user-options="store.ownerOptions"
      :allow-revision-select="!isRevisionScoped"
      :on-create="onCreate"
    />

    <!-- Row-click detail drawer — read-only summary + backlog items; its own
         Edit button opens the edit dialog below without closing the drawer.
         `revision-id` is null/omitted chain-wide (the dashboard's own
         behaviour — backlog items queued from here stay unassigned, same as
         the create dialog); the wizard's revision-scoped panels pass the
         open revision through, mirroring ManageBoard.vue. -->
    <EventDetailDrawer
      v-model:visible="drawerVisible"
      :category="category"
      :record="selectedEvent"
      :user-options="store.ownerOptions"
      :revision-id="props.revisionId"
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
// The dashboard's per-category screen, extracted from views/dashboard/
// CategoryView.vue (which now just resolves route params into props below)
// so the wizard's Manage & Monitor category sections
// (components/wizard/manage/Manage{Incident,Feature,Privacy,Task,Review}.vue)
// can mount the exact same screen, revision-scoped. Deliberately free of
// route assumptions — `category` is a prop, not read off route.params — so
// it works the same whether mounted as a routed view's body or inline inside
// the wizard tab. See dashboard.intent.md for why this lives here rather than
// under the locked views/dashboard/ set.
//
// DATA LOADING: this component does NOT call eventsStore.load()/
// backlogStore.load() itself — it only reads their already-loaded state
// (byCategory/kpis/breakdown getters). Whoever mounts it owns the load, at
// whatever scope that host needs:
//   - views/dashboard/DashboardLayout.vue loads both stores chain-wide (no
//     revisionId) for every child route, CategoryView.vue included — that
//     mirrors this component's pre-extraction behaviour exactly.
//   - Each components/wizard/manage/Manage<Category>.vue file loads both
//     stores scoped to (root project, open revision) itself, watching
//     [rootProjectId, revisionId] — the same pattern ManageBoard.vue and
//     ManageMetrics.vue already use, copied rather than shared so each
//     category section stays independently editable (see wizard.intent.md).
// This keeps exactly one load per (stores, scope) pair per screen instead of
// this component re-loading on top of a host that already did.
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryRankBar from '@/sections/project/components/dashboard/CategoryRankBar.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import EventDetailDialog from '@/sections/project/components/dashboard/EventDetailDialog.vue'
import EventDetailDrawer from '@/sections/project/components/dashboard/EventDetailDrawer.vue'
import NewEventDialog from '@/sections/project/components/dashboard/NewEventDialog.vue'
import RiskPips from '@/sections/project/components/dashboard/RiskPips.vue'
import TimeRangeSelect from '@/sections/project/components/dashboard/TimeRangeSelect.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { useRevisionLabel } from '@/sections/project/composables/useRevisionLabel'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { colorFor, orderIndex } from '@/sections/project/config/chartColors'
import {
  RANGES,
  adaptiveWindow,
  earliestPointDate,
  rangeFrom,
} from '@/sections/project/config/trend'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { toISODate } from '@/sections/project/stores/dateUtils'
import { useEventsStore } from '@/sections/project/stores/events'
import { useReportStore } from '@/sections/project/stores/report'
import { components, useConfirmDelete, useRightSidebarStore } from '@planetcrust/human-vue'
import { computed, inject, onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CResourceList } = components

const props = defineProps({
  // Which of the five categories to show (config/categories.js CATEGORY_CONFIG
  // key) — replaces the old route.params.category read.
  category: { type: String, required: true },
  // Chain-wide project id fed to the report endpoint for the KPI/donut/trend
  // metrics band. Only meaningful (and only read) when `revisionId` is unset
  // — see isRevisionScoped below. CategoryView.vue passes route.params.projectId
  // here, matching this component's pre-extraction behaviour exactly.
  projectId: { type: [String, Number], default: '' },
  // The OPEN revision (a lib system.Project instance's projectID) — absent
  // (null, the default) means the dashboard's chain-wide scope; set, it's the
  // wizard's Manage & Monitor single-revision scope (project.intent.md
  // "Dashboards" / wizard.intent.md). Drives three things: which metrics path
  // runs (see isRevisionScoped), whether the table's revision column renders,
  // and — the AGREED behaviour — that an item created from this panel is
  // assigned to it (events/backlogItems stores' add() trailing revisionId
  // arg), exactly like ManageBoard.vue's quick-add.
  revisionId: { type: [String, Number], default: null },
})

const { t } = useI18n()
const store = useEventsStore()
const backlogStore = useBacklogItemsStore()
const { revisionInfo } = useRevisionLabel()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

const category = computed(() => props.category)
const cfg = computed(() => CATEGORY_CONFIG[category.value] || null)
// KPI sparkline accent — the badge icon's text classes (via bg-current in
// CategoryKpiRow), NOT the chart-palette hex: identity surfaces track the
// icon colour exactly (same move as Overview's card rail), while multi-series
// charts keep the CVD-validated palette (where e.g. privacy is teal).
const accentClass = computed(() => cfg.value?.badge.text || '')

// Which metrics path is active: the report endpoint has no per-revision
// filter at all (server/system/rest/request/projectReport.go carries no
// revisionID field), so a `revisionId` prop switches the whole metrics band
// to client-side computation off the events/backlog stores' already
// revision-scoped rows instead of trying to (impossibly) scope the report
// calls themselves.
const isRevisionScoped = computed(() => !!props.revisionId)

// --- Chain-wide metrics band (KPIs + donuts + trend), report-endpoint-driven —
// all from the report endpoint (accurate, not subject to the events store's
// per-category 200-row list cap). `reportLoading`/`reportFailed` mirror the
// idiom in Overview.vue; `loadSeq` guards against a late response from a
// superseded project/category switch overwriting the current one's data
// (same pattern as Overview's loadAll/ActivityPanel's load()).
const reportStore = useReportStore()
const reportTrend = reactive({ labels: [], rangeLabels: [], series: [] })
// Default m6 (6 months) — matches Overview's default; see that view's comment
// for why.
const trendRange = ref('m6')
const reportBreakdowns = reactive({ total: 0, open: 0, overdue: 0, byDim: {} })
const reportLoading = ref(false)
const reportFailed = ref(false)
// `loadSeq` guards the KPI/donut report calls (loadReport); `trendSeq` guards
// the trend chart separately so a range-change reload (trend only) can't
// orphan an in-flight loadReport response from a category/project switch,
// and vice versa — same split as Overview.vue's loadSeq/trendSeq.
let loadSeq = 0
let trendSeq = 0

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
// Bucketed by the active range's adaptive window (day/week/month).
async function loadTrend(pid, key, mySeq) {
  const def = RANGES.find(r => r.key === trendRange.value)
  const from = rangeFrom(def)
  const now = new Date()
  // Fetch first: a bounded preset windows the query itself; 'all' (from =
  // null) fetches unbounded and sizes the axis from the earliest point below.
  const groupBy = CATEGORY_CONFIG[key].trendGroupBy
  const points = await reportStore.trend(pid, key, {
    from: from?.toISOString(),
    to: now.toISOString(),
    groupBy,
  })
  if (mySeq !== trendSeq) return // stale response
  // `points` is the flat, date-ascending array before the group pivot below —
  // its first element is the earliest overall point, which is all
  // earliestPointDate needs for the 'all' preset.
  const { labels, rangeLabels, bucket } = adaptiveWindow(from || earliestPointDate(points), now)

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
  reportTrend.labels = labels
  reportTrend.rangeLabels = rangeLabels
  reportTrend.series = keys.map(g => ({
    name: g,
    color: colorFor(groupBy, g),
    data: bucket(byGroup.get(g)),
  }))
}

// KPI tiles' sparkline — its own FIXED 12-week weekly pulse, deliberately
// decoupled from the trend chart's range control (the tiles give steady
// temporal context; the chart is the thing you window). One extra ungrouped
// report call per category load. `labels` holds the tooltip-facing range
// strings (e.g. "Jul 13 – 19, 2026"), not the short axis labels — the bars
// have no axis, so the fuller string is what CategoryKpiRow's tooltip shows.
const reportSpark = reactive({ labels: [], values: [] })

async function loadSpark(pid, key, mySeq) {
  const from = new Date()
  from.setDate(from.getDate() - 12 * 7)
  const { fromISO, toISO, rangeLabels, bucket } = adaptiveWindow(from, new Date())
  const points = await reportStore.trend(pid, key, { from: fromISO, to: toISO })
  if (mySeq !== loadSeq) return // stale response
  reportSpark.labels = rangeLabels
  reportSpark.values = bucket(points)
}

// `silent` refreshes the numbers in place (e.g. after a create) without
// flashing the skeleton.
async function loadMetrics(pid, key, { silent = false } = {}) {
  if (!pid || !key || !CATEGORY_CONFIG[key]) {
    ++trendSeq // invalidate any in-flight trend call (e.g. a pending range change)
    reportBreakdowns.total = 0
    reportBreakdowns.open = 0
    reportBreakdowns.overdue = 0
    reportBreakdowns.byDim = {}
    reportTrend.labels = []
    reportTrend.rangeLabels = []
    reportTrend.series = []
    reportSpark.labels = []
    reportSpark.values = []
    return
  }
  const mySeq = ++loadSeq
  const myTrendSeq = ++trendSeq
  if (!silent) reportLoading.value = true
  reportFailed.value = false
  try {
    await Promise.all([
      loadReport(pid, key, mySeq),
      loadSpark(pid, key, mySeq),
      loadTrend(pid, key, myTrendSeq),
    ])
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer switch
    console.error('Failed to load category metrics', err)
    reportFailed.value = true
  } finally {
    if (mySeq === loadSeq) reportLoading.value = false
  }
}

function retryMetrics() {
  loadMetrics(props.projectId, category.value)
}

// Chain-wide only — revision-scoped metrics are computed reactively below,
// no fetch to (re)run on a category/project switch.
watch(
  [category, () => props.projectId, isRevisionScoped],
  () => {
    if (!isRevisionScoped.value) loadMetrics(props.projectId, category.value)
  },
  { immediate: true },
)

// Preset change: reload the trend only — loadReport (KPIs/donuts) and the
// KPI sparkline (its own fixed 12-week pulse, see loadSpark) are untouched.
// Chain-wide only (see the TimeRangeSelect's v-if in the template).
async function onRangeChange(key) {
  trendRange.value = key
  const pid = props.projectId
  if (!pid || !category.value || !CATEGORY_CONFIG[category.value]) return
  const myTrendSeq = ++trendSeq
  try {
    await loadTrend(pid, category.value, myTrendSeq)
  } catch (err) {
    if (myTrendSeq !== trendSeq) return
    console.error('Failed to load category activity trend', err)
  }
}

// --- Revision-scoped metrics band — computed client-side off the
// events/backlog stores' already-loaded rows (the mounting Manage<Category>.vue
// loaded them scoped to (root project, this revision) — see the module doc
// comment above), mirroring components/wizard/manage/ManageMetrics.vue's own
// aggregate panel exactly, narrowed to one category. No separate load, no
// failure state (a store load failure is already logged by the store itself;
// there is nothing here to retry).
const localKpis = computed(() => store.kpis(category.value))

// Chart data helper — grouped counts for a report dimension (chain-wide) or a
// raw event field (revision-scoped), ordered canonically for ranked variants
// (severity/risk/status) so bars read worst→best.
function sortBreakdown(variant, rows) {
  if (!variant || variant === 'type') return rows
  return [...rows].sort(
    (a, b) =>
      orderIndex(variant, a.label) - orderIndex(variant, b.label) || a.label.localeCompare(b.label),
  )
}
const breakdownFor = chart => {
  const rows = isRevisionScoped.value
    ? store.breakdown(category.value, chart.field)
    : reportBreakdowns.byDim[chart.variant] || []
  return sortBreakdown(chart.variant, rows)
}

// Created-over-time trend, computed from the category's already-scoped rows'
// own `createdAt` (every system resource carries one) — the report-endpoint
// equivalent (loadTrend above) isn't available per-revision. Same day/week/
// month adaptiveWindow as the chain-wide path, fixed to its ~6-month default
// (no range selector — see the template's TimeRangeSelect v-if).
const localTrend = computed(() => {
  if (!cfg.value) return { labels: [], rangeLabels: [], series: [] }
  const groupBy = cfg.value.trendGroupBy
  const byGroup = new Map()
  for (const e of store.byCategory(category.value)) {
    if (!e.createdAt) continue
    const d = new Date(e.createdAt)
    if (Number.isNaN(d.getTime())) continue
    const g = e[groupBy] || '—'
    if (!byGroup.has(g)) byGroup.set(g, [])
    byGroup.get(g).push({ date: toISODate(d), value: 1 })
  }
  const { labels, rangeLabels, bucket } = adaptiveWindow()
  const keys = [...byGroup.keys()].sort(
    (a, b) => orderIndex(groupBy, a) - orderIndex(groupBy, b) || a.localeCompare(b),
  )
  return {
    labels,
    rangeLabels,
    series: keys.map(g => ({ name: g, color: colorFor(groupBy, g), data: bucket(byGroup.get(g)) })),
  }
})

// KPI tiles' sparkline, revision-scoped equivalent of loadSpark above — same
// fixed 12-week weekly pulse, built from the category's already-loaded rows.
const localSpark = computed(() => {
  const from = new Date()
  from.setDate(from.getDate() - 12 * 7)
  const { rangeLabels, bucket } = adaptiveWindow(from, new Date())
  const points = []
  for (const e of store.byCategory(category.value)) {
    if (!e.createdAt) continue
    const d = new Date(e.createdAt)
    if (Number.isNaN(d.getTime())) continue
    points.push({ date: toISODate(d), value: 1 })
  }
  return { labels: rangeLabels, values: bucket(points) }
})

// --- Metrics band, unified across both paths ------------------------------
const metricsLoading = computed(() =>
  isRevisionScoped.value ? store.loading || backlogStore.loading : reportLoading.value,
)
// The revision-scoped path never fails independently of the store load
// itself (which is already handled/logged by the host that triggered it) —
// there is nothing here to retry.
const metricsFailed = computed(() => (isRevisionScoped.value ? false : reportFailed.value))

const kpiList = computed(() => {
  if (!cfg.value) return []
  const source = isRevisionScoped.value ? localKpis.value : reportBreakdowns
  return cfg.value.kpis.map(({ key, labelKey }) => ({ labelKey, value: source[key] }))
})

const trend = computed(() => (isRevisionScoped.value ? localTrend.value : reportTrend))
const spark = computed(() => (isRevisionScoped.value ? localSpark.value : reportSpark))

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
  set: v =>
    v ? rightSidebar.open('project-event-detail') : rightSidebar.close('project-event-detail'),
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

// `unassignedOnly` only applies chain-wide (see the #filter slot's v-if) —
// revision-scoped mode never has an unassigned row to begin with (its store
// load is itself scoped to one revisionID).
const filter = reactive({ query: '', unassignedOnly: false })
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value?.toggle(event)
}
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
// The trailing revisionID column is fixed (every category gets it, see the
// #body-revisionID slot above) rather than part of cfg.columns, and only
// added chain-wide — see that slot's comment for why the wizard's
// revision-scoped panels drop it instead.
const fields = computed(() => {
  const cols = (cfg.value?.columns ?? []).map(c => ({
    key: c.key,
    header: t(c.headerKey),
    sortable: true,
  }))
  if (isRevisionScoped.value) return cols
  return [
    ...cols,
    { key: 'revisionID', header: t('project.dashboard.columns.revision'), sortable: false },
  ]
})

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}

// Client-side query filter across the item's string values.
const filteredItems = computed(() => {
  let list = cfg.value ? store.byCategory(category.value) : []
  if (!isRevisionScoped.value && filter.unassignedOnly) {
    list = list.filter(item => revisionInfo(item.revisionID).unassigned)
  }
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

watch(
  visibleItems,
  list => {
    pagination.total = list.length
  },
  { immediate: true },
)

const onSort = ({ sortField, sortOrder }) => {
  if (!sortField) return
  sorting.sortBy = sortField
  sorting.sortDesc = sortOrder === -1
}

// Switching categories (the dashboard route can swap `category` on a live
// instance without remounting — the wizard's dispatch remounts instead, so
// this is a no-op there): reset the transient list/dialog state so the new
// page starts clean (store data stays live).
watch(category, () => {
  filter.query = ''
  filter.unassignedOnly = false
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
// dialog (see NewEventDialog's onSubmit contract). `props.revisionId` is
// forwarded to both add() calls — AGREED BEHAVIOUR: null/omitted chain-wide
// (the dashboard's "+ New" stays unassigned), the open revision when this
// panel is revision-scoped (the wizard's Manage & Monitor panels), mirroring
// ManageBoard.vue's quick-add exactly.
const onCreate = async (payload, backlogTitles = []) => {
  let event
  try {
    event = await store.add(category.value, payload, props.revisionId)
  } catch (err) {
    console.error('Failed to create event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.createFailed'))(err)
    return false
  }
  // The chain-wide metrics band is report-driven, so the new record isn't in
  // it yet — refresh in place (no skeleton flash). The revision-scoped band
  // is a computed over the store's own rows, so it already reflects the new
  // record with no refresh needed.
  if (!isRevisionScoped.value) {
    loadMetrics(props.projectId, category.value, { silent: true })
  }
  $toast.toastSuccess(
    t('project.dashboard.newButton', { type: t(cfg.value.singularKey) }),
    t('project.dashboard.event.toast.created'),
  )
  if (backlogTitles.length) {
    const results = await Promise.allSettled(
      backlogTitles.map(title =>
        backlogStore.add(
          {
            category: category.value,
            eventID: event.id,
            title,
            priority: 'Medium',
            status: 'Open',
          },
          props.revisionId,
        ),
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
    if (!isRevisionScoped.value) {
      loadMetrics(props.projectId, category.value, { silent: true })
    }
    $toast.toastSuccess(t(cfg.value.singularKey), t('project.dashboard.event.toast.updated'))
    if (selectedEvent.value?.id === String(id)) {
      selectedEvent.value =
        store.byCategory(category.value).find(e => e.id === String(id)) || selectedEvent.value
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
    if (!isRevisionScoped.value) {
      loadMetrics(props.projectId, category.value, { silent: true })
    }
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
