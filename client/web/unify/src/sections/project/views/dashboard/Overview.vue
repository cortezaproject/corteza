<template>
  <!-- Project overview: one card per category fusing its KPIs (total in the
       doughnut hole, open count in the header) with a status doughnut, plus a
       cross-category activity trend. Driven entirely by the server-side report
       endpoint (grouped counts, not full rowsets). Scrolls internally. -->
  <div class="h-full overflow-y-auto p-4 flex flex-col gap-6">
    <!-- Failed: a report call errored (e.g. a transient API failure). Replaces
         the cards + trend with a visible retry rather than silently rendering
         them as "0 items" (EventsActivityPanel below fails soft on its own). -->
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
      <!-- Category cards — KPI header + status doughnut, links to the list. -->
      <section>
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.dashboard.views.dashboard') }}
        </h2>

        <!-- First load (or a project switch): skeleton cards rather than a
             zeroed-out grid, so nothing reads as "this project has no data". -->
        <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-3">
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
                <span class="h-3 w-2/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
                <span class="h-2 w-1/3 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
              </span>
            </span>
            <span class="h-32 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
          </div>
        </div>

        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-3">
          <RouterLink
            v-for="c in cards"
            :key="c.key"
            :to="{ name: 'project.overview.category', params: { projectId, category: c.key } }"
            class="group relative overflow-hidden rounded-xl border border-surface bg-surface pl-5 pr-4 py-3 flex flex-col gap-1 hover:shadow-md transition-all"
          >
            <!-- Left accent rail in the category's colour. -->
            <span class="absolute inset-y-0 left-0 w-1" :style="{ background: c.accentColor }" />
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
          </RouterLink>
        </div>
      </section>

      <!-- Audit events — pulse + recent activity feed (admin-only; self-hides). -->
      <EventsActivityPanel :project-id="projectId" />

      <!-- Activity — new items created over time, stacked by category. The
           project's pulse in one glance. Range control is right-aligned next
           to the heading; changing it reloads only the trend below, not the
           category cards above. -->
      <section>
        <div class="flex items-center justify-between gap-3 mb-3">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color">
            {{ $t('project.dashboard.chart.activity') }}
          </h2>
          <TimeRangeSelect :model-value="trendRange" @update:model-value="onRangeChange" />
        </div>
        <div
          v-if="loading || trendLoading"
          class="rounded-lg border border-surface bg-surface p-4"
        >
          <div class="h-48 rounded bg-emphasis animate-pulse motion-reduce:animate-none" />
        </div>
        <CategoryTrendChart
          v-else
          title-key="project.dashboard.chart.trend"
          :labels="trend.labels"
          :series="trend.series"
        />
      </section>
    </template>
  </div>
</template>

<script setup>
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventsActivityPanel from '@/sections/project/components/dashboard/EventsActivityPanel.vue'
import TimeRangeSelect from '@/sections/project/components/dashboard/TimeRangeSelect.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { CATEGORY_COLORS } from '@/sections/project/config/chartColors'
import { RANGES, adaptiveWindow, rangeFrom } from '@/sections/project/config/trend'
import { useReportStore } from '@/sections/project/stores/report'
import { isOpenStatus } from '@/sections/project/stores/events'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const report = useReportStore()

const projectId = computed(() => route.params.projectId)

// Per-category aggregates: key -> { total, open, status: [{ label, value }] }.
const data = reactive({})

// Stacked created-over-time trend across all categories, windowed by the
// preset below. Default m6 (6 months): a full year of week/month-adaptive
// buckets read noisy for a landing-page glance, so 6 months is the starting
// point — the control lets you widen it in one click.
const trend = reactive({ labels: [], series: [] })
const trendRange = ref('m6')

// Loading covers the category cards + trend (EventsActivityPanel manages its
// own loading/failed state independently). `failed` surfaces a visible error
// with retry instead of the previous silent "0 items" on a failed report.
// `trendLoading` is the narrower one: a range-change reload touches only the
// trend chart, not the category cards, so it gets its own flag.
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
    return { key, titleKey: cfg.titleKey, badge: cfg.badge, accentColor: CATEGORY_COLORS[key], ...d }
  }),
)

// One report per category (grouped by status) yields the total (sum), the open
// count (sum of non-Completed — see stores/events.js#isOpenStatus, the single
// definition of that rule) and the chart series in a single call.
async function loadCategory(pid, key, mySeq) {
  const rows = await report.report(pid, key, { dimensions: ['status'] })
  if (mySeq !== loadSeq) return // stale response — a newer project load is in flight
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

// Created-over-time, stacked by category. One report per category (grouped by
// day), bucketed by the active range's adaptive window (day/week/month) and
// aligned to a shared axis.
async function loadTrend(pid, mySeq) {
  const def = RANGES.find(r => r.key === trendRange.value)
  const { fromISO, toISO, labels, bucket } = adaptiveWindow(rangeFrom(def), new Date())
  const series = await Promise.all(
    CATEGORY_ORDER.map(async key => {
      const points = await report.trend(pid, key, { from: fromISO, to: toISO })
      return {
        name: t(CATEGORY_CONFIG[key].titleKey),
        color: CATEGORY_COLORS[key],
        data: bucket(points),
      }
    }),
  )
  if (mySeq !== trendSeq) return // stale response
  trend.labels = labels
  trend.series = series
}

// Two sequence tokens: `loadSeq` guards the category cards, `trendSeq` guards
// the trend chart. They're separate because a range change reloads only the
// trend — bumping a single shared counter would also orphan an in-flight
// category-card response from the initial project load. Both still bump
// together on a project switch (loadAll), so a stale trend from the previous
// project never lands either (same pattern as AllEventsView's load()).
let loadSeq = 0
let trendSeq = 0

async function loadAll(pid) {
  if (!pid) return
  const mySeq = ++loadSeq
  const myTrendSeq = ++trendSeq
  loading.value = true
  failed.value = false
  try {
    await Promise.all([
      ...CATEGORY_ORDER.map(key => loadCategory(pid, key, mySeq)),
      loadTrend(pid, myTrendSeq),
    ])
  } catch (err) {
    if (mySeq !== loadSeq) return // superseded by a newer project switch
    console.error('Failed to load project overview', err)
    failed.value = true
  } finally {
    if (mySeq === loadSeq) loading.value = false
  }
}

function retry() {
  loadAll(projectId.value)
}

// Preset change: reload the trend only — the per-category status reports
// (cards above) are untouched.
async function onRangeChange(key) {
  trendRange.value = key
  const pid = projectId.value
  if (!pid) return
  const myTrendSeq = ++trendSeq
  trendLoading.value = true
  try {
    await loadTrend(pid, myTrendSeq)
  } catch (err) {
    if (myTrendSeq !== trendSeq) return
    console.error('Failed to load project activity trend', err)
  } finally {
    if (myTrendSeq === trendSeq) trendLoading.value = false
  }
}

watch(projectId, loadAll, { immediate: true })
</script>
