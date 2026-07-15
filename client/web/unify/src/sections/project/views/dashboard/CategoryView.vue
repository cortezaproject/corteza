<template>
  <!-- Unknown category (bad deep-link / stale nav) — render a muted note only. -->
  <div v-if="!cfg" class="p-4">
    <p class="text-muted-color">{{ $t('project.dashboard.list.empty') }}</p>
  </div>

  <!-- The category page: a fixed title bar + metrics, then the list fills the
       rest. Structure mirrors the admin list views (e.g. DataSource/List.vue):
       an `overflow-hidden min-w-0` root so the wide table can't push the layout
       past the screen, with the list in a plain `flex-1 min-h-0` wrapper and
       `class="h-full"` on CResourceList (which owns its own internal scroll). -->
  <div v-else class="flex flex-col h-full min-h-0 min-w-0 overflow-hidden">
    <!-- Title bar — wizard-style: leading badge + title + description. -->
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 shrink-0"
        :class="[cfg.badge.bg, cfg.badge.ring]"
      >
        <i :class="[cfg.badge.icon, cfg.badge.text]" />
      </span>
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">{{ $t(cfg.titleKey) }}</h2>
        <p class="text-sm text-muted-color">{{ $t(cfg.descKey) }}</p>
      </div>
    </header>

    <!-- Metrics — KPI row + compact breakdown grid + trend, fixed above the
         list. Breakdowns flow up to four across on wide screens so the band
         stays shallow and leaves room for the list. -->
    <div class="shrink-0 p-4 flex flex-col gap-3">
      <CategoryKpiRow :kpis="kpiList" />
      <div class="grid grid-cols-2 xl:grid-cols-4 gap-3">
        <CategoryDonutChart
          v-for="chart in cfg.charts"
          :key="chart.field"
          :title-key="chart.titleKey"
          :data="breakdownFor(chart.field, chart.variant)"
          :variant="chart.variant"
          :accent="accentColor"
          :height="180"
        />
      </div>
      <CategoryTrendChart
        title-key="project.dashboard.chart.createdOverTime"
        :labels="trend.labels"
        :series="trend.series"
        :accent="accentColor"
        :height="180"
      />
    </div>

    <!-- The category's items — fills the remaining space; the list scrolls
         internally (both axes) so the wide table never overflows the page. -->
    <div class="flex-1 min-h-0 px-4 pb-4">
      <CResourceList
        class="h-full"
        primary-key="id"
        :fields="fields"
        :items="visibleItems"
        :filter="filter"
        @update:filter="Object.assign(filter, $event)"
        :sorting="sorting"
        :pagination="pagination"
        :loading="store.loading"
        :translations="{
          searchPlaceholder: $t('project.dashboard.list.searchPlaceholder'),
          noItems: $t('project.dashboard.list.empty'),
        }"
        @sort="onSort"
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
          <span v-if="col.kind === 'id'" class="font-mono text-xs text-muted-color">
            {{ data[col.key] }}
          </span>

          <div v-else-if="col.kind === 'title'" class="flex flex-col">
            <span class="font-medium text-color">{{ data.title }}</span>
            <span v-if="data.description" class="text-xs text-muted-color truncate max-w-[18rem]">
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

          <div v-else-if="col.kind === 'backlog'" class="flex flex-wrap gap-1">
            <span
              v-for="bl in data.backlog || []"
              :key="bl"
              class="font-mono text-[10px] px-1.5 py-0.5 rounded border border-surface text-muted-color"
            >
              {{ bl }}
            </span>
            <span v-if="!(data.backlog && data.backlog.length)" class="text-muted-color">—</span>
          </div>

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
      @create="onCreate"
    />
  </div>
</template>

<script setup>
import CategoryDonutChart from '@/sections/project/components/dashboard/CategoryDonutChart.vue'
import CategoryKpiRow from '@/sections/project/components/dashboard/CategoryKpiRow.vue'
import CategoryTrendChart from '@/sections/project/components/dashboard/CategoryTrendChart.vue'
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import NewEventDialog from '@/sections/project/components/dashboard/NewEventDialog.vue'
import RiskPips from '@/sections/project/components/dashboard/RiskPips.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { CATEGORY_COLORS, colorFor, orderIndex } from '@/sections/project/config/chartColors'
import { bucketWeekly, trendWindow, weekLabel } from '@/sections/project/config/trend'
import { useEventsStore } from '@/sections/project/stores/events'
import { useReportStore } from '@/sections/project/stores/report'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CResourceList } = components

const { t } = useI18n()
const route = useRoute()
const store = useEventsStore()
const $toast = inject('$toast')

// Active category comes straight from the route; cfg is null for unknown keys.
const category = computed(() => route.params.category)
const cfg = computed(() => CATEGORY_CONFIG[category.value] || null)
const accentColor = computed(() => CATEGORY_COLORS[category.value] || '')

// Created-over-time trend for this category, from the report endpoint (accurate,
// not subject to the events store's per-category list cap).
const reportStore = useReportStore()
const trend = reactive({ labels: [], series: [] })

async function loadTrend() {
  const pid = route.params.projectId
  if (!pid || !cfg.value) {
    trend.labels = []
    trend.series = []
    return
  }
  const { fromISO, toISO, starts } = trendWindow(12)
  const groupBy = cfg.value.trendGroupBy
  const points = await reportStore
    .trend(pid, category.value, { from: fromISO, to: toISO, groupBy })
    .catch(() => [])
  trend.labels = starts.map(weekLabel)

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
  trend.series = keys.map(g => ({
    name: g,
    color: colorFor(groupBy, g),
    data: bucketWeekly(byGroup.get(g), starts),
  }))
}

watch([category, () => route.params.projectId], loadTrend, { immediate: true })

const dialogVisible = ref(false)
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

// KPI cards: pull each configured metric off the store's live kpis(cat) result.
const kpiList = computed(() => {
  if (!cfg.value) return []
  const k = store.kpis(category.value)
  return cfg.value.kpis.map(({ key, labelKey }) => ({ labelKey, value: k[key] }))
})

// Chart data helper — grouped counts for a field, ordered canonically for
// ranked variants (severity/risk/status) so bars read worst→best.
const breakdownFor = (field, variant) => {
  const rows = store.breakdown(category.value, field)
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

// Client-side sort (byCategory is already newest-first as the default order).
const visibleItems = computed(() => {
  const list = [...filteredItems.value]
  const { sortBy, sortDesc } = sorting
  if (sortBy) {
    list.sort((a, b) => {
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
  sorting.sortBy = 'dateDue'
  sorting.sortDesc = true
})

// Create handler — persist via the store (list/KPIs/charts/nav badge all react
// off the returned record), toast, and close. Errors surface as a toast and
// keep the dialog open so the user can retry.
const onCreate = async payload => {
  try {
    await store.add(category.value, payload)
    $toast.toastSuccess(
      t('project.dashboard.newButton', { type: t(cfg.value.singularKey) }),
      t('project.dashboard.event.toast.created'),
    )
    dialogVisible.value = false
  } catch (err) {
    console.error('Failed to create event', err)
    $toast.toastErrorHandler(t('project.dashboard.event.toast.createFailed'))(err)
  }
}
</script>
