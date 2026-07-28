<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 p-4 pb-3 border-b border-surface">
      <h2 class="text-lg font-medium mb-1">{{ $t('project.manage.views.metrics') }}</h2>
      <p class="text-sm text-muted-color">{{ $t('project.manage.metrics.blurb') }}</p>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <!-- Empty state matters here for the same reason as ManageBoard.vue's:
         every work item defaults to unassigned until revision-assignment UI
         exists, so an open revision legitimately has nothing to plot — explain
         that instead of leaving what looks like broken/blank charts. -->
    <div
      v-else-if="!totalItems"
      class="flex-1 flex flex-col items-center justify-center text-center gap-2 px-6"
    >
      <span
        class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emphasis text-muted-color"
      >
        <i class="pi pi-chart-bar text-xl" />
      </span>
      <p class="text-sm font-medium text-color">{{ $t('project.manage.metrics.empty.title') }}</p>
      <p class="text-sm text-muted-color max-w-md">{{ $t('project.manage.metrics.empty.hint') }}</p>
    </div>

    <div v-else class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-3">
      <CategoryKpiRow :kpis="kpiList" />

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <CategoryDonutChart
          title-key="project.dashboard.chart.byStatus"
          :data="statusBreakdown"
          variant="status"
          :height="200"
        />
        <!-- Priority (High/Medium/Low) only exists on backlog items — the
             five event categories carry `severity` instead, a different
             ordinal scale (see config/chartColors.js) — so this mirrors
             BacklogView.vue's own priority breakdown rather than forcing
             severity values through the priority palette. -->
        <CategoryRankBar
          title-key="project.dashboard.chart.byPriority"
          :data="priorityBreakdown"
          variant="priority"
          :height="200"
        />
      </div>

      <CategoryTrendChart
        title-key="project.dashboard.chart.completedOverTime"
        :labels="throughput.labels"
        :range-labels="throughput.rangeLabels"
        :series="throughput.series"
        :height="200"
      />
    </div>
  </div>
</template>

<script setup>
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryRankBar from '@/sections/project/components/dashboard/CategoryRankBar.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import { colorFor, PRIORITY_ORDER, STATUS_ORDER } from '@/sections/project/config/chartColors'
import { adaptiveWindow } from '@/sections/project/config/trend'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { toISODate } from '@/sections/project/stores/dateUtils'
import { isOpenStatus, useEventsStore } from '@/sections/project/stores/events'
import { computed, watch } from 'vue'

// Manage & Monitor > Monitor > Metrics — a read-only snapshot of the ONE
// revision open in the wizard (route.params.projectId), same revision-scoping
// contract as ManageBoard.vue: eventsStore/backlogStore are loaded with (root
// project, this revision), never the whole project (that's the live
// dashboard's scope; see project.intent.md's "Dashboards" locked contract).
// Covers the same six item types as the board: the five event categories
// (stores/events.js) plus backlog items (stores/backlogItems.js).
const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  // Accepted for the shared Wizard.vue dispatch prop set (see
  // ManagePlaceholder.vue); this section is read-only so there is nothing to
  // disable.
  disabled: { type: Boolean, default: false },
})

const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

const rootProjectId = computed(() => props.project?.rootProjectID || props.project?.projectID)
const revisionId = computed(() => props.project?.projectID)
const loading = computed(() => eventsStore.loading || backlogStore.loading)

// Scope both dashboard stores to (root project, this revision), exactly like
// ManageBoard.vue — re-fetches whenever the open revision changes (e.g.
// switching revisions via the wizard-header revision switcher).
watch(
  [rootProjectId, revisionId],
  ([pid, rid]) => {
    if (!pid || !rid) return
    eventsStore.load(pid, rid)
    backlogStore.load(pid, rid)
  },
  { immediate: true },
)

const totalItems = computed(() => eventsStore.events.length + backlogStore.items.length)

// --- KPI trio ----------------------------------------------------------
// total/open/overdue across all six item types. Events: sum the store's own
// per-category `kpis` getter (never re-derive the open rule) across every
// category. Backlog items have no per-category kpis getter of their own, so
// they're tallied directly here the same way BacklogView.vue's kpiList does
// — via the shared `isOpenStatus` rule, not a re-derived one.
const kpiList = computed(() => {
  const totals = { total: 0, open: 0, overdue: 0 }
  for (const cat of eventsStore.categories) {
    const k = eventsStore.kpis(cat)
    totals.total += k.total
    totals.open += k.open
    totals.overdue += k.overdue
  }
  const todayISO = toISODate(new Date())
  const blOpen = backlogStore.items.filter(i => isOpenStatus(i.status))
  totals.total += backlogStore.items.length
  totals.open += blOpen.length
  totals.overdue += blOpen.filter(i => i.dateDue && i.dateDue < todayISO).length
  return [
    { labelKey: 'project.dashboard.kpi.total', value: totals.total },
    { labelKey: 'project.dashboard.kpi.open', value: totals.open },
    { labelKey: 'project.dashboard.kpi.overdue', value: totals.overdue },
  ]
})

// --- Status distribution -------------------------------------------------
// All six item types share the same four-value status set (config/
// eventForm.js EVENT_STATUS), so this spans events + backlog items, same as
// the board's shared column set.
const statusBreakdown = computed(() => {
  const counts = {}
  const bump = status => {
    const key = String(status || '')
    if (key) counts[key] = (counts[key] || 0) + 1
  }
  eventsStore.events.forEach(e => bump(e.status))
  backlogStore.items.forEach(i => bump(i.status))
  return STATUS_ORDER.filter(l => counts[l]).map(label => ({ label, value: counts[label] }))
})

// --- Priority distribution -------------------------------------------------
// Backlog items only — see the template comment above.
const priorityBreakdown = computed(() => {
  const counts = {}
  for (const i of backlogStore.items) {
    const key = String(i.priority || '')
    if (key) counts[key] = (counts[key] || 0) + 1
  }
  return PRIORITY_ORDER.filter(l => counts[l]).map(label => ({ label, value: counts[label] }))
})

// --- Throughput ------------------------------------------------------------
// Items completed per day/week, adaptively bucketed (config/trend.js). Only
// incident/task carry a `completedDate` field (config/eventForm.js); every
// other category and backlog items don't, so `updatedAt` (the audit
// timestamp every system resource carries, refreshed by updateStatus's PUT
// whenever an item is moved to Completed) is the realistic proxy completion
// date for those. No range selector: a single revision's item count is small
// (often zero pre-assignment), so the default ~6-month adaptiveWindow is
// enough — matches its own "sane empty chart" rationale.
function completedDateOf(item) {
  const raw = item.completedDate || item.updatedAt
  if (!raw) return ''
  // completedDate already arrives as a plain YYYY-MM-DD string (see
  // stores/events.js DATE_KEYS); updatedAt is a full audit timestamp, so only
  // that branch needs reformatting through toISODate.
  if (/^\d{4}-\d{2}-\d{2}$/.test(raw)) return raw
  const d = new Date(raw)
  return Number.isNaN(d.getTime()) ? '' : toISODate(d)
}

const throughput = computed(() => {
  const points = []
  const collect = item => {
    if (item.status !== 'Completed') return
    const date = completedDateOf(item)
    if (date) points.push({ date, value: 1 })
  }
  eventsStore.events.forEach(collect)
  backlogStore.items.forEach(collect)

  const { labels, rangeLabels, bucket } = adaptiveWindow()
  return {
    labels,
    rangeLabels,
    // Status enum values are literal, not translated copy, throughout this
    // section (see ManageBoard.vue's board.blurb comment) — same here.
    series: [{ name: 'Completed', color: colorFor('status', 'Completed'), data: bucket(points) }],
  }
})
</script>
