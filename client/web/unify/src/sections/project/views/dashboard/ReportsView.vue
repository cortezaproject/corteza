<template>
  <!-- Cross-category analytics: a time-range filter drives a KPI band (with
       trend deltas vs the equal-length prior window), severity/open-overdue
       breakdowns by category, and a created-over-time trend — all from the
       server-side report endpoint (grouped counts, windowed on created-at).
       Scrolls internally, same idiom as Overview.vue. -->
  <div class="h-full overflow-y-auto p-4 flex flex-col gap-6">
    <!-- Filter row — right-aligned like the section's list toolbars. Always
         interactive (even mid-load) so switching ranges never feels stuck. -->
    <div class="flex justify-end">
      <SelectButton
        v-model="range"
        :options="rangeOptions"
        option-label="label"
        option-value="value"
        :allow-empty="false"
        size="small"
      />
    </div>

    <!-- Failed: a report call errored. Same visible-retry idiom as Overview/
         CategoryView, reusing their error/retry copy. -->
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

    <template v-else-if="loading">
      <!-- First load (or a project/range switch): skeleton rather than a
           zeroed-out page. -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        <div
          v-for="n in 4"
          :key="n"
          class="rounded-xl border border-surface bg-surface px-4 py-3 flex flex-col gap-2"
        >
          <span class="h-3 w-1/2 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
          <span class="h-6 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
        </div>
      </div>
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <div v-for="n in 2" :key="n" class="rounded-lg border border-surface bg-surface p-4">
          <div class="h-48 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
        </div>
      </div>
      <div class="rounded-lg border border-surface bg-surface p-4">
        <div class="h-48 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
      </div>
    </template>

    <template v-else>
      <!-- KPI band. Created/Completion rate are straightforward window
           aggregates; Open/Overdue are point-in-time states (a row can be
           "open" long after the window it was created in), only counted here
           over rows created in the window — not a live cross-project total.
           Each tile's delta compares against the equal-length prior window
           (omitted entirely for "All time", which has no prior window). -->
      <section class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        <div
          v-for="tile in kpiTiles"
          :key="tile.key"
          class="rounded-xl border border-surface bg-surface px-4 py-3 flex flex-col gap-1"
        >
          <div class="text-[11px] font-semibold uppercase tracking-wide text-muted-color">
            {{ $t(tile.labelKey) }}
          </div>
          <div class="text-2xl font-semibold text-color leading-none">{{ tile.value }}</div>
          <!-- Delta is informational, not good/bad — muted text, never
               green/red. -->
          <div
            v-if="tile.delta !== null"
            class="text-xs text-muted-color"
            :title="deltaCopy(tile.delta)"
            :aria-label="deltaCopy(tile.delta)"
          >
            {{ deltaArrow(tile.delta) }} {{ Math.abs(tile.delta) }}%
          </div>
        </div>
      </section>

      <!-- Severity by category (composition — "where is the risk"): each
           row's bar is 100%-stacked to ITS OWN total, review excluded (no
           severity dimension). -->
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <ReportStackedBar
          title-key="project.dashboard.reports.chart.severityByCategory"
          :rows="severityRows"
          scale="share"
          :legend="severityLegend"
          :height="180"
        />

        <!-- Open vs overdue by category (volume — bar length IS the open
             count, comparable across categories on one shared scale;
             overdue renders as its darker leading segment). -->
        <ReportStackedBar
          title-key="project.dashboard.reports.chart.openOverdueByCategory"
          :rows="openOverdueRows"
          scale="value"
          :legend="openOverdueLegend"
          :height="200"
        />
      </div>

      <!-- Created over time, stacked by category — same chart + weekly
           bucketing as Overview, windowed by the active range ("All time"
           defaults to the last 12 months of weeks). -->
      <CategoryTrendChart
        title-key="project.dashboard.chart.trend"
        :labels="trend.labels"
        :series="trend.series"
        :height="220"
      />
    </template>
  </div>
</template>

<script setup>
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import ReportStackedBar from '@/sections/project/components/dashboard/ReportStackedBar.vue'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { CATEGORY_COLORS, colorFor, orderIndex, RISK_COLORS, SEVERITY_COLORS, SEVERITY_ORDER, STATUS_COLORS } from '@/sections/project/config/chartColors'
import { bucketWeekly, weekLabel, weekStarts } from '@/sections/project/config/trend'
import { useReportStore } from '@/sections/project/stores/report'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const report = useReportStore()

const projectId = computed(() => route.params.projectId)

// Categories charted in the severity breakdown — review has no severity
// dimension (see server/system/service/project_report.go's source registry).
const SEVERITY_CATEGORIES = CATEGORY_ORDER.filter(key => key !== 'review')

// Time-range filter. Default 90 days per spec; 'all' omits from/to entirely
// (no windowing, no prior-window delta).
const range = ref('90d')
const rangeOptions = computed(() => [
  { value: '30d', label: t('project.dashboard.reports.range.d30') },
  { value: '90d', label: t('project.dashboard.reports.range.d90') },
  { value: '12m', label: t('project.dashboard.reports.range.m12') },
  { value: 'all', label: t('project.dashboard.reports.range.all') },
])

const DAY_MS = 86400000

function daysAgo(from, days) {
  return new Date(from.getTime() - days * DAY_MS)
}

function monthsAgo(from, months) {
  const d = new Date(from)
  d.setMonth(d.getMonth() - months)
  return d
}

// Resolve the active range key into the current window (Dates, for the
// trend's week bucketing, plus their ISO strings for the report calls), the
// equal-length prior window (null for 'all' — there is no delta), and the
// trend's own from/to (same as current, except 'all' defaults to 12 months).
function computeRange(key) {
  const now = new Date()
  if (key === 'all') {
    return { current: {}, previous: null, trendFrom: monthsAgo(now, 12), trendTo: now }
  }

  let from
  if (key === '30d') from = daysAgo(now, 30)
  else if (key === '12m') from = monthsAgo(now, 12)
  else from = daysAgo(now, 90) // '90d' and any unknown value

  const lengthMs = now.getTime() - from.getTime()
  const prevTo = from
  const prevFrom = new Date(from.getTime() - lengthMs)

  return {
    current: { fromISO: from.toISOString(), toISO: now.toISOString() },
    previous: { fromISO: prevFrom.toISOString(), toISO: prevTo.toISOString() },
    trendFrom: from,
    trendTo: now,
  }
}

// Per-category metrics for the current window (key -> { count, open, overdue }),
// mirrored for the equal-length prior window (null entries when 'all').
const perCategory = reactive({})
const perCategoryPrev = reactive({})
// Per-category severity breakdown for the current window (key -> [{ label, value }]).
const severityByCategory = reactive({})
const trend = reactive({ labels: [], series: [] })

const loading = ref(false)
const failed = ref(false)
let loadSeq = 0

async function loadCategoryMetrics(pid, key, from, to, mySeq) {
  const rows = await report.report(pid, key, { metrics: ['count', 'open', 'overdue'], from, to })
  if (mySeq !== loadSeq) return // stale response
  const g = rows[0]?.metrics || {}
  return { count: Number(g.count || 0), open: Number(g.open || 0), overdue: Number(g.overdue || 0) }
}

async function loadSeverity(pid, key, from, to, mySeq) {
  const rows = await report.report(pid, key, { dimensions: ['severity'], from, to })
  if (mySeq !== loadSeq) return // stale response
  return rows.map(r => ({ label: String(r.dimensions?.severity || '—'), value: Number(r.metrics?.count || 0) }))
}

async function loadTrend(pid, rangeInfo, mySeq) {
  const starts = weekStarts(rangeInfo.trendFrom, rangeInfo.trendTo)
  const from = rangeInfo.trendFrom.toISOString()
  const to = rangeInfo.trendTo.toISOString()
  const series = await Promise.all(
    CATEGORY_ORDER.map(async key => {
      const points = await report.trend(pid, key, { from, to })
      return {
        name: t(CATEGORY_CONFIG[key].titleKey),
        color: CATEGORY_COLORS[key],
        data: bucketWeekly(points, starts),
      }
    }),
  )
  if (mySeq !== loadSeq) return // stale response
  trend.labels = starts.map(weekLabel)
  trend.series = series
}

// One report call per category for the current window's totals, one more per
// category for the severity breakdown (review excluded), one per category for
// the prior window's totals (skipped entirely for "All time"), plus one
// per-category trend call — capped well under the ~25/load budget (19 at
// most, 14 for "All time").
async function loadAll(pid, rangeKey) {
  if (!pid) return
  const mySeq = ++loadSeq
  loading.value = true
  failed.value = false
  const rangeInfo = computeRange(rangeKey)
  try {
    const [currentResults, previousResults, severityResults] = await Promise.all([
      Promise.all(
        CATEGORY_ORDER.map(key =>
          loadCategoryMetrics(pid, key, rangeInfo.current.fromISO, rangeInfo.current.toISO, mySeq),
        ),
      ),
      rangeInfo.previous
        ? Promise.all(
            CATEGORY_ORDER.map(key =>
              loadCategoryMetrics(pid, key, rangeInfo.previous.fromISO, rangeInfo.previous.toISO, mySeq),
            ),
          )
        : Promise.resolve(null),
      Promise.all(
        SEVERITY_CATEGORIES.map(key =>
          loadSeverity(pid, key, rangeInfo.current.fromISO, rangeInfo.current.toISO, mySeq),
        ),
      ),
      loadTrend(pid, rangeInfo, mySeq),
    ])
    if (mySeq !== loadSeq) return // superseded by a newer switch

    CATEGORY_ORDER.forEach((key, i) => {
      perCategory[key] = currentResults[i]
      perCategoryPrev[key] = previousResults ? previousResults[i] : null
    })
    SEVERITY_CATEGORIES.forEach((key, i) => {
      severityByCategory[key] = severityResults[i]
    })
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer switch
    console.error('Failed to load project reports', err)
    failed.value = true
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

function retry() {
  loadAll(projectId.value, range.value)
}

watch([projectId, range], () => loadAll(projectId.value, range.value), { immediate: true })

// Cross-category KPI aggregates for the current window and (when available)
// the equal-length prior window.
function sumWindow(source) {
  let count = 0
  let open = 0
  let overdue = 0
  let any = false
  for (const key of CATEGORY_ORDER) {
    const d = source[key]
    if (!d) continue
    any = true
    count += d.count
    open += d.open
    overdue += d.overdue
  }
  const completion = count ? Math.round(((count - open) / count) * 100) : null
  return any ? { count, open, overdue, completion } : null
}

const kpiCurrent = computed(() => sumWindow(perCategory) || { count: 0, open: 0, overdue: 0, completion: null })
const kpiPrevious = computed(() => (range.value === 'all' ? null : sumWindow(perCategoryPrev)))

// Relative % change for count-like metrics; null (no delta shown) when there
// is no prior window or the prior value was zero (an undefined/∞ change).
function pctDelta(curr, prev) {
  if (!prev) return null
  return Math.round(((curr - prev) / prev) * 100)
}

// Percentage-POINT change for the completion-rate tile (it's already a %, so
// "up 3%" reads as +3 points, not a further ratio on the rate itself).
function ptsDelta(curr, prev) {
  if (curr == null || prev == null) return null
  return curr - prev
}

const kpiTiles = computed(() => {
  const cur = kpiCurrent.value
  const prev = kpiPrevious.value
  return [
    {
      key: 'created',
      labelKey: 'project.dashboard.reports.kpi.created',
      value: cur.count,
      delta: prev ? pctDelta(cur.count, prev.count) : null,
    },
    {
      key: 'open',
      labelKey: 'project.dashboard.reports.kpi.open',
      value: cur.open,
      delta: prev ? pctDelta(cur.open, prev.open) : null,
    },
    {
      key: 'overdue',
      labelKey: 'project.dashboard.reports.kpi.overdue',
      value: cur.overdue,
      delta: prev ? pctDelta(cur.overdue, prev.overdue) : null,
    },
    {
      key: 'completion',
      labelKey: 'project.dashboard.reports.kpi.completionRate',
      value: cur.completion == null ? '—' : `${cur.completion}%`,
      delta: prev ? ptsDelta(cur.completion, prev.completion) : null,
    },
  ]
})

function deltaArrow(delta) {
  if (delta > 0) return '▲'
  if (delta < 0) return '▼'
  return '•'
}

function deltaCopy(delta) {
  if (delta > 0) return t('project.dashboard.reports.kpi.deltaUp', { pct: delta })
  if (delta < 0) return t('project.dashboard.reports.kpi.deltaDown', { pct: Math.abs(delta) })
  return t('project.dashboard.reports.kpi.deltaFlat')
}

// Severity-by-category rows (composition): segments in canonical severity
// order, share-scaled (each row fills its own 100%). Review is excluded.
// Legends — segment identity must never rest on color+tooltip alone.
// Severity labels are the persisted enum values, same as the donut legends.
const severityLegend = SEVERITY_ORDER.map(label => ({ label, color: SEVERITY_COLORS[label] }))
const openOverdueLegend = computed(() => [
  { label: t('project.dashboard.reports.kpi.overdue'), color: RISK_COLORS.Critical },
  { label: t('project.dashboard.reports.kpi.open'), color: STATUS_COLORS.Open },
])

const severityRows = computed(() =>
  SEVERITY_CATEGORIES.map(key => {
    const raw = severityByCategory[key] || []
    const sorted = [...raw].sort(
      (a, b) => orderIndex('severity', a.label) - orderIndex('severity', b.label) || a.label.localeCompare(b.label),
    )
    const total = sorted.reduce((s, r) => s + r.value, 0)
    const segments = sorted.map(r => {
      const share = total ? Math.round((r.value / total) * 100) : 0
      return {
        label: r.label,
        value: r.value,
        color: colorFor('severity', r.label),
        title: `${r.label}: ${r.value} (${share}%)`,
      }
    })
    return { key, label: t(CATEGORY_CONFIG[key].titleKey), total, segments, rightText: String(total) }
  }),
)

// Open-vs-overdue rows (volume): overdue is the darker leading segment,
// remainder-open the neutral "Open" blue; value-scaled so bar length compares
// across categories on one shared axis.
const openOverdueRows = computed(() =>
  CATEGORY_ORDER.map(key => {
    const d = perCategory[key] || { open: 0, overdue: 0 }
    const remainder = Math.max(0, d.open - d.overdue)
    return {
      key,
      label: t(CATEGORY_CONFIG[key].titleKey),
      total: d.open,
      segments: [
        { label: 'overdue', value: d.overdue, color: RISK_COLORS.Critical },
        { label: 'open', value: remainder, color: STATUS_COLORS.Open },
      ],
      rightText: t('project.dashboard.reports.chart.openOverdueLabel', { open: d.open, overdue: d.overdue }),
    }
  }),
)
</script>
