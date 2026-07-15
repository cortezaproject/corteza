<template>
  <!-- Project overview: one card per category fusing its KPIs (total in the
       doughnut hole, open count in the header) with a status doughnut, plus a
       cross-category activity trend. Driven entirely by the server-side report
       endpoint (grouped counts, not full rowsets). Scrolls internally. -->
  <div class="h-full overflow-y-auto p-4 flex flex-col gap-6">
    <!-- Category cards — KPI header + status doughnut, links to the list. -->
    <section>
      <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
        {{ $t('project.dashboard.views.dashboard') }}
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-5 gap-3">
        <RouterLink
          v-for="c in cards"
          :key="c.key"
          :to="{ name: 'project.overview.category', params: { projectId, category: c.key } }"
          class="group relative overflow-hidden rounded-xl border border-surface bg-surface-50 dark:bg-surface-950 pl-5 pr-4 py-3 flex flex-col gap-1 hover:shadow-md transition-all"
        >
          <!-- Left accent rail in the category's colour. -->
          <span class="absolute inset-y-0 left-0 w-1" :style="{ background: c.accentColor }" />
          <span class="flex items-center gap-2 min-w-0">
            <span
              class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
              :class="[c.badge.bg, c.badge.ring]"
            >
              <i :class="[c.badge.icon, c.badge.text]" />
            </span>
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
         project's pulse in one glance. -->
    <section>
      <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-color mb-3">
        {{ $t('project.dashboard.chart.activity') }}
      </h2>
      <CategoryTrendChart
        title-key="project.dashboard.chart.trend"
        :labels="trend.labels"
        :series="trend.series"
      />
    </section>
  </div>
</template>

<script setup>
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventsActivityPanel from '@/sections/project/components/dashboard/EventsActivityPanel.vue'
import { CATEGORY_CONFIG, CATEGORY_ORDER } from '@/sections/project/config/categories'
import { CATEGORY_COLORS } from '@/sections/project/config/chartColors'
import { bucketWeekly, trendWindow, weekLabel } from '@/sections/project/config/trend'
import { useReportStore } from '@/sections/project/stores/report'
import { computed, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const report = useReportStore()

const projectId = computed(() => route.params.projectId)

// Per-category aggregates: key -> { total, open, status: [{ label, value }] }.
const data = reactive({})

// Stacked created-over-time trend across all categories (last 12 weeks).
const trend = reactive({ labels: [], series: [] })

// Cards join the live aggregates with each category's static visual config,
// in the fixed display order. Falls back to zeros before a load resolves.
const cards = computed(() =>
  CATEGORY_ORDER.map(key => {
    const cfg = CATEGORY_CONFIG[key]
    const d = data[key] || { total: 0, open: 0, status: [] }
    return { key, titleKey: cfg.titleKey, badge: cfg.badge, accentColor: CATEGORY_COLORS[key], ...d }
  }),
)

// One report per category (grouped by status) yields the total (sum), the open
// count (sum of non-Completed — same rule as the events store) and the chart
// series in a single call.
async function loadCategory(pid, key) {
  const rows = await report.report(pid, key, { dimensions: ['status'] })
  let total = 0
  let open = 0
  const status = rows.map(r => {
    const label = String(r.dimensions?.status || '—')
    const value = Number(r.metrics?.count || 0)
    total += value
    if (label !== 'Completed') open += value
    return { label, value }
  })
  data[key] = { total, open, status }
}

// Created-over-time, stacked by category. One report per category (grouped by
// day) bucketed into weeks and aligned to a shared week axis.
async function loadTrend(pid) {
  const { fromISO, toISO, starts } = trendWindow(12)
  trend.labels = starts.map(weekLabel)
  trend.series = await Promise.all(
    CATEGORY_ORDER.map(async key => {
      const points = await report.trend(pid, key, { from: fromISO, to: toISO }).catch(() => [])
      return {
        name: t(CATEGORY_CONFIG[key].titleKey),
        color: CATEGORY_COLORS[key],
        data: bucketWeekly(points, starts),
      }
    }),
  )
}

function loadAll(pid) {
  if (!pid) return
  for (const key of CATEGORY_ORDER) {
    data[key] = { total: 0, open: 0, status: [] }
    loadCategory(pid, key).catch(err => console.error(`Failed to load ${key} report`, err))
  }
  loadTrend(pid).catch(err => console.error('Failed to load trend', err))
}

watch(projectId, loadAll, { immediate: true })
</script>
