<template>
  <!-- Project overview: one card per category fusing its KPIs (total in the
       doughnut hole, open count in the header) with a status doughnut, plus a
       cross-category activity trend. Driven entirely by the server-side report
       endpoint (grouped counts, not full rowsets). Scrolls internally. -->
  <div class="h-full overflow-y-auto p-4 flex flex-col gap-6">
    <!-- Failed: a report call errored (e.g. a transient API failure). Replaces
         the cards + trend (+ KPI trio, revision-scoped) with a visible retry
         rather than silently rendering them as "0 items"
         (EventsActivityPanel below fails soft on its own). -->
    <section v-if="failed" class="flex flex-col items-center text-center gap-2 py-12">
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
        @click="retry"
      />
    </section>

    <template v-else>
      <!-- Category cards — KPI header + status doughnut, links to the list.
           Chain-wide and revision-scoped alike (the report call underneath
           just narrows to `revisionId` when set — see loadCategory). -->
      <section>
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.dashboard.nav.categories') }}
        </h2>

        <!-- First load (or a project/revision switch): skeleton cards rather
             than a zeroed-out grid, so nothing reads as "this project has no
             data". -->
        <div
          v-if="loading"
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-3"
        >
          <div
            v-for="n in CATEGORY_ORDER.length"
            :key="n"
            class="rounded-xl border border-surface bg-surface pl-5 pr-4 py-3 flex flex-col gap-3"
          >
            <span class="flex items-center gap-2 min-w-0">
              <span
                class="inline-flex w-8 h-8 rounded-md bg-emphasis shrink-0 animate-pulse motion-reduce:animate-none"
              />
              <span class="flex flex-col gap-1 flex-1 min-w-0">
                <span
                  class="h-3 w-2/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
                />
                <span
                  class="h-2 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none"
                />
              </span>
            </span>
            <span class="h-32 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
          </div>
        </div>

        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-3">
          <!-- Chain-wide: RouterLinks to the dashboard's own category route.
               Revision-scoped: no child route to link to (the wizard's
               category sections live behind ManageNav's `section` query
               param, not real routes — see config/manageNav.js), so the card
               renders as a button emitting category-selected instead, picked
               up by components/wizard/manage/ManageOverview.vue (re-emits)
               and Wizard.vue (sets activeSection) — see cardTag/cardProps/
               cardListeners below. -->
          <component
            :is="cardTag"
            v-for="c in cards"
            :key="c.key"
            v-bind="cardProps(c.key)"
            v-on="cardListeners(c.key)"
            class="group relative overflow-hidden rounded-xl border border-surface bg-surface pl-5 pr-4 py-3 flex flex-col gap-1 hover:shadow-md transition-all"
            :class="
              isRevisionScoped
                ? 'w-full text-left cursor-pointer focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
                : ''
            "
          >
            <!-- Left accent rail takes the badge icon's exact text colour via
                 bg-current — light+dark for free. -->
            <span class="absolute inset-y-0 left-0 w-1 bg-current" :class="c.badge.text" />
            <span class="flex items-center gap-2 min-w-0">
              <KindIcon :config="c.badge" size="lg" plain-icon />
              <span class="flex flex-col min-w-0">
                <span class="text-sm font-medium text-color truncate">{{ $t(c.titleKey) }}</span>
                <span class="text-xs text-muted-color">
                  {{ c.open }} {{ $t('project.dashboard.kpi.open').toLowerCase() }}
                </span>
              </span>
            </span>
            <!-- Status doughnut (total in the hole) fused into the card. -->
            <CategoryDonutChart bare variant="status" :data="c.status" :height="150" />
          </component>
        </div>
      </section>

      <!-- Audit events — pulse + recent activity feed (admin-only; self-hides).
           Chain-wide only: this is the project-wide action log, which has no
           per-revision scope at all (EventsActivityPanel takes only
           `projectId`) — dropped entirely revision-scoped rather than shown
           unscoped and misleadingly labelled. -->
      <EventsActivityPanel v-if="!isRevisionScoped" :project-id="props.projectId" />

      <!-- KPI trio — total/open/overdue across every category PLUS backlog
           items, revision-scoped only (ruled 2026-07-28): chain-wide Overview
           has never carried a single cross-category trio (each card already
           states its own open count), and adding one there was out of scope
           for this change. Revision-scoped, this is what used to be
           ManageMetrics.vue's KPI row — its status donut and priority
           rank-bar are deliberately NOT reproduced here: redundant with the
           category cards above (each already breaks its own category down by
           status), doubly so at a single revision's typically small volume. -->
      <section v-if="isRevisionScoped">
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.dashboard.overview.kpiHeading') }}
        </h2>
        <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
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
        <CategoryKpiRow v-else :kpis="kpiList" />
      </section>

      <!-- Activity — new items created over time, stacked by category. The
           project's (or, revision-scoped, this revision's) pulse in one
           glance. Range control is right-aligned next to the heading,
           chain-wide only (see the TimeRangeSelect v-if below); changing it
           reloads only the trend below, not the category cards/KPI trio
           above. -->
      <section>
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.dashboard.chart.activity') }}
        </h2>
        <div v-if="loading || trendLoading" class="rounded-lg border border-surface bg-surface p-4">
          <div class="h-48 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
        </div>
        <CategoryTrendChart
          v-else
          title-key="project.dashboard.chart.trend"
          :labels="trend.labels"
          :range-labels="trend.rangeLabels"
          :series="trend.series"
          legend-variant="category"
        >
          <!-- Windows this chart only — lives in the card it acts on.
               Chain-wide only: a single revision's item count is small
               (often zero pre-assignment), so the revision-scoped path skips
               the selector and stays on the fixed ~6-month adaptiveWindow
               default (mirrors CategoryPanel.vue's own revision-scoped trend
               / the old ManageMetrics.vue throughput chart, same
               rationale). -->
          <template v-if="!isRevisionScoped" #actions>
            <TimeRangeSelect :model-value="trendRange" @update:model-value="onRangeChange" />
          </template>
        </CategoryTrendChart>
      </section>
    </template>
  </div>
</template>

<script setup>
// The project overview, extracted from views/dashboard/Overview.vue (which
// now just resolves route.params into props below) so the wizard's Manage &
// Monitor tab can mount the exact same screen, revision-scoped, absorbing
// what used to be the separate ManageMetrics.vue panel (ruled 2026-07-28 —
// see views/views.intent.md / DashboardLayout.intent.md: the dashboard and
// the wizard's Manage & Monitor tab are ONE surface differing only in scope).
// Deliberately free of route assumptions — `projectId`/`revisionId` are
// props, not read off route.params — so it works the same whether mounted as
// a routed view's body (views/dashboard/Overview.vue) or inline inside the
// wizard tab (components/wizard/manage/ManageOverview.vue). Mirrors
// components/dashboard/CategoryPanel.vue's own extraction exactly; see
// dashboard.intent.md for why both live here rather than under the locked
// views/dashboard/ set.
//
// DATA LOADING: unlike CategoryPanel, this component has NO events/backlog
// store dependency at all, chain-wide or revision-scoped — every number here
// (category cards, KPI trio, trend) comes from the report endpoint, which now
// accepts an optional revisionID that scopes aggregation server-side (see
// stores/report.js). That is a deliberate, EXACT alternative to the
// events/backlog stores' 200-row-per-category cap: ManageOverview.vue does
// not need to watch/load those stores the way ManageBoard.vue and the
// Manage<Category>.vue files do — it only derives (rootProjectId, revisionId)
// from the `project` prop and passes them straight through as props here.
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventsActivityPanel from '@/sections/project/components/dashboard/EventsActivityPanel.vue'
import TimeRangeSelect from '@/sections/project/components/dashboard/TimeRangeSelect.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { CATEGORY_COLORS } from '@/sections/project/config/chartColors'
import {
  RANGES,
  adaptiveWindow,
  earliestPointDate,
  rangeFrom,
} from '@/sections/project/config/trend'
import { useReportStore } from '@/sections/project/stores/report'
import { isOpenStatus } from '@/sections/project/stores/events'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

const { t } = useI18n()
const report = useReportStore()

const props = defineProps({
  // Chain-root project id — every report call's required ProjectID scope,
  // regardless of whether `revisionId` narrows it further. views/dashboard/
  // Overview.vue passes route.params.projectId (a chain HEAD, which is the
  // chain ROOT for an original, never-revised project — see
  // system.Project#rootProjectID); components/wizard/manage/ManageOverview.vue
  // passes project.rootProjectID.
  projectId: { type: [String, Number], default: '' },
  // The OPEN revision (a lib system.Project instance's projectID) — absent
  // (null, the default) means the dashboard's chain-wide scope (every
  // revision in the chain); set, it's the wizard's Manage & Monitor
  // single-revision scope (project.intent.md "Dashboards" / wizard.intent.md).
  // Drives: whether report calls narrow to this revision (now supported
  // server-side — server/system/service/project_report.go, commit
  // d3e89ab07), whether the KPI trio section renders, whether the audit
  // EventsActivityPanel renders (chain-wide only), whether the category cards
  // link to the dashboard's per-category route, and whether the trend chart's
  // range selector renders.
  revisionId: { type: [String, Number], default: null },
})

// Revision-scoped only: a category card switches the wizard's Manage &
// Monitor section instead of navigating (see cardTag/cardListeners below).
// Chain-wide never emits this — its cards stay plain RouterLinks.
const emit = defineEmits(['category-selected'])

const isRevisionScoped = computed(() => !!props.revisionId)

// Per-category aggregates: key -> { total, open, status: [{ label, value }] }.
const data = reactive({})

// Stacked created-over-time trend across all categories, windowed by the
// preset below (chain-wide) or a fixed ~6-month lookback (revision-scoped —
// see loadTrend). Default m6 (6 months): a full year of week/month-adaptive
// buckets read noisy for a landing-page glance, so 6 months is the starting
// point — the control lets you widen it in one click.
const trend = reactive({ labels: [], rangeLabels: [], series: [] })
const trendRange = ref('m6')

// total/open/overdue across every category plus backlog items — revision-
// scoped KPI trio only (see the template's isRevisionScoped v-if).
const kpiTotals = reactive({ total: 0, open: 0, overdue: 0 })
const kpiList = computed(() => [
  { labelKey: 'project.dashboard.kpi.total', value: kpiTotals.total },
  { labelKey: 'project.dashboard.kpi.open', value: kpiTotals.open },
  { labelKey: 'project.dashboard.kpi.overdue', value: kpiTotals.overdue },
])
// The report endpoint's resource registry (server/system/service/
// project_report.go projectReportSources) covers the five categories plus
// 'backlog-item' — the same six item types ManageMetrics.vue's old KPI trio
// tallied off the events/backlog stores, now via one report call per
// resource instead.
const KPI_RESOURCES = [...CATEGORY_ORDER, 'backlog-item']

// Loading covers the category cards + trend (+ KPI trio, revision-scoped) —
// EventsActivityPanel manages its own loading/failed state independently.
// `failed` surfaces a visible error with retry instead of the previous silent
// "0 items" on a failed report. `trendLoading` is the narrower one: a
// range-change reload touches only the trend chart, not the cards/KPI trio.
const loading = ref(false)
const trendLoading = ref(false)
const failed = ref(false)

// Cards join the live aggregates with each category's static visual config,
// in the fixed display order. Only rendered once `loading` clears (see
// template), so no zeroed placeholder data is needed here.
const cards = computed(() =>
  CATEGORY_ORDER.map(key => {
    const cfg = CATEGORY_CONFIG[key]
    const d = data[key] || { total: 0, open: 0, status: [] }
    return { key, titleKey: cfg.titleKey, badge: cfg.badge, ...d }
  }),
)

// Card element/props/listeners — RouterLink chain-wide (the dashboard's
// per-category route, untouched), a `button` revision-scoped (see the
// template comment above for why there's nothing to route to there — it
// emits category-selected instead, picked up the same way ManageNav's own
// item clicks are).
const cardTag = computed(() => (isRevisionScoped.value ? 'button' : RouterLink))
const cardProps = key =>
  isRevisionScoped.value
    ? { type: 'button' }
    : {
        to: {
          name: 'project.overview.category',
          params: { projectId: props.projectId, category: key },
        },
      }
const cardListeners = key =>
  isRevisionScoped.value ? { click: () => emit('category-selected', key) } : {}

// One report per category (grouped by status) yields the total (sum), the open
// count (sum of non-Completed — see stores/events.js#isOpenStatus, the single
// definition of that rule) and the chart series in a single call. `revId`
// narrows the aggregation server-side when set (see stores/report.js).
async function loadCategory(pid, key, revId, mySeq) {
  const rows = await report.report(pid, key, { dimensions: ['status'], revisionId: revId })
  if (mySeq !== loadSeq) return // stale response — a newer project/revision switch is in flight
  let total = 0
  let open = 0
  const status = rows.map(r => {
    const label = String(r.dimensions?.status || '—')
    const value = Number(r.metrics?.count || 0)
    total += value
    if (isOpenStatus(label)) open += value
    return { label, value }
  })
  data[key] = { total, open, status }
}

// Grand-total (count/open/overdue) across every KPI_RESOURCES entry —
// revision-scoped KPI trio only (see loadAll below).
async function loadKpiTotals(pid, revId, mySeq) {
  const rowsets = await Promise.all(
    KPI_RESOURCES.map(key =>
      report.report(pid, key, { metrics: ['count', 'open', 'overdue'], revisionId: revId }),
    ),
  )
  if (mySeq !== loadSeq) return // stale response
  let total = 0
  let open = 0
  let overdue = 0
  for (const rows of rowsets) {
    const g = rows?.[0]?.metrics || {}
    total += Number(g.count || 0)
    open += Number(g.open || 0)
    overdue += Number(g.overdue || 0)
  }
  kpiTotals.total = total
  kpiTotals.open = open
  kpiTotals.overdue = overdue
}

// Created-over-time, stacked by category. One report per category (grouped by
// day), bucketed by the active range's adaptive window (day/week/month) and
// aligned to a shared axis.
async function loadTrend(pid, revId, mySeq) {
  const now = new Date()
  // Revision-scoped: no range selector (see the template's TimeRangeSelect
  // v-if) — fixed ~6-month lookback, mirroring CategoryPanel.vue's own
  // revision-scoped trend / the old ManageMetrics.vue throughput chart, same
  // rationale (a single revision's item count is small, often zero
  // pre-assignment). Chain-wide: the active preset (trendRange); 'all'
  // resolves to an unbounded fetch sized from the earliest point below.
  let from
  if (revId) {
    from = new Date(now)
    from.setMonth(from.getMonth() - 6)
  } else {
    from = rangeFrom(RANGES.find(r => r.key === trendRange.value))
  }
  const toISO = now.toISOString()
  const perCategory = await Promise.all(
    CATEGORY_ORDER.map(key =>
      report.trend(pid, key, { from: from?.toISOString(), to: toISO, revisionId: revId }),
    ),
  )
  if (mySeq !== trendSeq) return // stale response
  const { labels, rangeLabels, bucket } = adaptiveWindow(
    from || earliestPointDate(...perCategory),
    now,
  )
  trend.labels = labels
  trend.rangeLabels = rangeLabels
  trend.series = CATEGORY_ORDER.map((key, i) => ({
    name: t(CATEGORY_CONFIG[key].titleKey),
    key,
    color: CATEGORY_COLORS[key],
    data: bucket(perCategory[i]),
  }))
}

// Two sequence tokens: `loadSeq` guards the category cards (+ KPI trio),
// `trendSeq` guards the trend chart. They're separate because a range change
// reloads only the trend — bumping a single shared counter would also orphan
// an in-flight category-card response from the initial load. Both still bump
// together on a project/revision switch (loadAll), so a stale trend from the
// previous scope never lands either (same pattern as AllEventsView's load()).
let loadSeq = 0
let trendSeq = 0

async function loadAll(pid, revId) {
  if (!pid) return
  const mySeq = ++loadSeq
  const myTrendSeq = ++trendSeq
  loading.value = true
  failed.value = false
  try {
    const tasks = [
      ...CATEGORY_ORDER.map(key => loadCategory(pid, key, revId, mySeq)),
      loadTrend(pid, revId, myTrendSeq),
    ]
    // KPI trio only fetched revision-scoped (see the template's v-if) — no
    // point spending six extra report calls on a chain-wide load that never
    // renders them.
    if (revId) tasks.push(loadKpiTotals(pid, revId, mySeq))
    await Promise.all(tasks)
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer switch
    console.error('Failed to load project overview', err)
    failed.value = true
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

function retry() {
  loadAll(props.projectId, props.revisionId)
}

// Preset change: reload the trend only — the per-category status reports
// (cards above) and the KPI trio are untouched. Chain-wide only (see the
// TimeRangeSelect's v-if in the template).
async function onRangeChange(key) {
  trendRange.value = key
  const pid = props.projectId
  if (!pid) return
  const myTrendSeq = ++trendSeq
  trendLoading.value = true
  try {
    await loadTrend(pid, props.revisionId, myTrendSeq)
  } catch (err) {
    if (myTrendSeq !== trendSeq) return
    console.error('Failed to load project activity trend', err)
  } finally {
    if (myTrendSeq === trendSeq) trendLoading.value = false
  }
}

watch([() => props.projectId, () => props.revisionId], ([pid, revId]) => loadAll(pid, revId), {
  immediate: true,
})
</script>
