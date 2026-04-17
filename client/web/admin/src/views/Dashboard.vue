<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('dashboard.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full p-3 sm:p-6 gap-4 sm:gap-10 overflow-y-auto min-w-0">
    <!-- Stat cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-4">
      <div
        v-for="card in statCards"
        :key="card.key"
        class="flex flex-col gap-2 p-4 sm:p-5 rounded-border border border-surface bg-surface cursor-pointer transition-all duration-200 hover:shadow-lg hover:-translate-y-0.5 hover:border-primary min-w-0"
        :class="{ '!border-primary': activeCard === card.key }"
        @click="switchCard(card.key)"
      >
        <div class="flex items-center justify-between">
          <span class="text-muted-color text-sm font-medium truncate">{{ card.label }}</span>
          <div class="flex items-center justify-center w-8 h-8 rounded-full bg-highlight shrink-0">
            <i :class="card.icon" class="text-primary" />
          </div>
        </div>

        <span class="text-3xl font-bold text-color">
          {{ card.loading ? '—' : card.value.toLocaleString() }}
        </span>

        <!-- Status breakdown pills -->
        <div v-if="!card.loading && card.statuses.length" class="flex flex-wrap gap-1.5 mt-1">
          <span
            v-for="status in card.statuses"
            :key="status.label"
            class="inline-flex items-center gap-1 text-xs font-medium px-2 py-0.5 rounded-full"
            :class="status.class"
          >
            <span class="w-1.5 h-1.5 rounded-full shrink-0" :class="status.dotClass" />
            {{ status.count }} {{ status.label }}
          </span>
        </div>
      </div>
    </div>

    <!-- Monthly stacked chart -->
    <div class="rounded-border p-3 sm:p-5 min-w-0">
      <div class="h-80 sm:h-100">
        <CChart v-if="monthlyChartOptions" :chart="monthlyChartOptions" :key="activeCard" />
        <div v-else class="flex items-center justify-center h-full">
          <ProgressSpinner style="width: 40px; height: 40px" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CChart } = components
const { t } = useI18n()

const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')

// ── Reactive state ──────────────────────────────────────────
const loading = ref(true)

// Stat card totals
const totalUsers = ref(0)
const totalRoles = ref(0)
const totalWorkflows = ref(0)
const totalNamespaces = ref(0)

// User status breakdown
const activeUsers = ref(0)
const suspendedUsers = ref(0)
const deletedUsers = ref(0)

// Role status breakdown
const archivedRoles = ref(0)
const deletedRoles = ref(0)
const activeRoles = ref(0)

// Workflow status breakdown
const disabledWorkflows = ref(0)
const deletedWorkflows = ref(0)
const enabledWorkflows = ref(0)

// Active card & chart
const activeCard = ref('users')
const monthlyChartOptions = ref(null)

// ── Theme helpers ───────────────────────────────────────────
function getCssVar(name) {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

function themeColors() {
  return {
    primary: getCssVar('--p-primary-color') || '#6366f1',
    text: getCssVar('--p-text-color') || '#334155',
    textMuted: getCssVar('--p-text-muted-color') || '#94a3b8',
    border: getCssVar('--p-content-border-color') || '#e2e8f0',
    bg: getCssVar('--p-content-background') || '#ffffff',
    green: getCssVar('--p-green-500') || '#22c55e',
    yellow: getCssVar('--p-yellow-500') || '#eab308',
    red: getCssVar('--p-red-500') || '#ef4444',
    blue: getCssVar('--p-blue-500') || '#3b82f6',
  }
}

// ── Status style helpers ────────────────────────────────────
const greenPill = {
  class: 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-300',
  dotClass: 'bg-green-500',
}
const yellowPill = {
  class: 'bg-yellow-50 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-300',
  dotClass: 'bg-yellow-500',
}
const redPill = {
  class: 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300',
  dotClass: 'bg-red-500',
}
// ── Stat cards ──────────────────────────────────────────────
const statCards = computed(() => [
  {
    key: 'users',
    label: t('dashboard.stats.users'),
    value: totalUsers.value,
    icon: 'pi pi-users',
    loading: loading.value,
    statuses: [
      { count: activeUsers.value, label: t('dashboard.status.active'), ...greenPill },
      { count: suspendedUsers.value, label: t('dashboard.status.suspended'), ...yellowPill },
      { count: deletedUsers.value, label: t('dashboard.status.deleted'), ...redPill },
    ],
  },
  {
    key: 'roles',
    label: t('dashboard.stats.roles'),
    value: totalRoles.value,
    icon: 'pi pi-shield',
    loading: loading.value,
    statuses: [
      { count: activeRoles.value, label: t('dashboard.status.active'), ...greenPill },
      { count: archivedRoles.value, label: t('dashboard.status.archived'), ...yellowPill },
      { count: deletedRoles.value, label: t('dashboard.status.deleted'), ...redPill },
    ],
  },
  {
    key: 'workflows',
    label: t('dashboard.stats.workflows'),
    value: totalWorkflows.value,
    icon: 'pi pi-sitemap',
    loading: loading.value,
    statuses: [
      { count: enabledWorkflows.value, label: t('dashboard.status.enabled'), ...greenPill },
      { count: disabledWorkflows.value, label: t('dashboard.status.disabled'), ...yellowPill },
      { count: deletedWorkflows.value, label: t('dashboard.status.deleted'), ...redPill },
    ],
  },
  {
    key: 'namespaces',
    label: t('dashboard.stats.namespaces'),
    value: totalNamespaces.value,
    icon: 'pi pi-database',
    loading: loading.value,
    statuses: [],
  },
])

// ── Cheap count helper ──────────────────────────────────────
async function fetchCount(apiCall) {
  try {
    const { response } = apiCall
    const result = await response()
    return result.filter?.total ?? result.filter?.count ?? 0
  } catch {
    return 0
  }
}

// ── Fetch all items with pagination ─────────────────────────
async function fetchAllItems(listFn, params = {}) {
  const allItems = []
  let cursor = ''

  do {
    try {
      const result = await listFn({
        ...params,
        limit: 200,
        sort: 'createdAt ASC',
        incTotal: true,
        pageCursor: cursor || undefined,
      })

      const items = result.set || []
      allItems.push(...items)
      cursor = result.filter?.nextPage || ''
    } catch {
      break
    }
  } while (cursor)

  return allItems
}

// ── Status classifiers per resource type ────────────────────
const statusClassifiers = {
  users: item => {
    if (item.deletedAt) return 'deleted'
    if (item.suspendedAt) return 'suspended'
    return 'active'
  },
  roles: item => {
    if (item.deletedAt) return 'deleted'
    if (item.archivedAt) return 'archived'
    return 'active'
  },

  workflows: item => {
    if (item.deletedAt) return 'deleted'
    if (!item.enabled) return 'disabled'
    return 'enabled'
  },
  namespaces: () => 'active',
}

// Status display config per resource type
const statusConfig = {
  users: [
    { key: 'active', label: () => t('dashboard.status.active'), color: 'green' },
    { key: 'suspended', label: () => t('dashboard.status.suspended'), color: 'yellow' },
    { key: 'deleted', label: () => t('dashboard.status.deleted'), color: 'red' },
  ],
  roles: [
    { key: 'active', label: () => t('dashboard.status.active'), color: 'green' },
    { key: 'archived', label: () => t('dashboard.status.archived'), color: 'yellow' },
    { key: 'deleted', label: () => t('dashboard.status.deleted'), color: 'red' },
  ],

  workflows: [
    { key: 'enabled', label: () => t('dashboard.status.enabled'), color: 'green' },
    { key: 'disabled', label: () => t('dashboard.status.disabled'), color: 'yellow' },
    { key: 'deleted', label: () => t('dashboard.status.deleted'), color: 'red' },
  ],
  namespaces: [{ key: 'active', label: () => t('dashboard.status.active'), color: 'green' }],
}

// ── Group items by month + status ───────────────────────────
function groupByMonthAndStatus(items, classifier) {
  // { "YYYY-MM": { active: N, suspended: N, ... } }
  const months = {}

  for (const item of items) {
    const created = item.createdAt
    if (!created) continue
    const monthKey = created.slice(0, 7)
    const status = classifier(item)

    if (!months[monthKey]) months[monthKey] = {}
    months[monthKey][status] = (months[monthKey][status] || 0) + 1
  }

  // Fill in gaps
  const sortedKeys = Object.keys(months).sort()
  if (sortedKeys.length > 1) {
    const [fy, fm] = sortedKeys[0].split('-').map(Number)
    const [ly, lm] = sortedKeys[sortedKeys.length - 1].split('-').map(Number)

    const d = new Date(fy, fm - 1)
    const end = new Date(ly, lm - 1)

    while (d <= end) {
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
      if (!months[key]) months[key] = {}
      d.setMonth(d.getMonth() + 1)
    }
  }

  const finalKeys = Object.keys(months).sort()
  const labels = finalKeys.map(k => {
    const [y, m] = k.split('-')
    const date = new Date(Number(y), Number(m) - 1)
    return date.toLocaleDateString(undefined, { month: 'short', year: 'numeric' })
  })

  return { monthKeys: finalKeys, labels, data: months }
}

// ── Build stacked monthly chart ─────────────────────────────
function buildStackedChart(resourceKey, items) {
  const c = themeColors()
  const classifier = statusClassifiers[resourceKey]
  const config = statusConfig[resourceKey]
  const { monthKeys, labels, data } = groupByMonthAndStatus(items, classifier)

  const colorMap = {
    green: c.green,
    yellow: c.yellow,
    red: c.red,
  }

  // Build cumulative total line
  const cumulative = []
  let sum = 0
  for (const mk of monthKeys) {
    const monthData = data[mk] || {}
    const monthTotal = Object.values(monthData).reduce((a, b) => a + b, 0)
    sum += monthTotal
    cumulative.push(sum)
  }

  // For each month, find which status is the topmost non-zero series
  // so we can round only that segment's top corners
  const topmostPerMonth = monthKeys.map(mk => {
    const monthData = data[mk] || {}
    // Walk from last (topmost) to first, find first non-zero
    for (let i = config.length - 1; i >= 0; i--) {
      if ((monthData[config[i].key] || 0) > 0) return i
    }
    return -1
  })

  // Build stacked bar series for each status
  const barSeries = config.map((st, idx) => ({
    name: st.label(),
    type: 'bar',
    stack: 'statuses',
    data: monthKeys.map((mk, mIdx) => ({
      value: (data[mk] || {})[st.key] || 0,
      itemStyle: {
        borderRadius: topmostPerMonth[mIdx] === idx ? [6, 6, 0, 0] : 0,
      },
    })),
    barMaxWidth: 40,
    itemStyle: {
      color: colorMap[st.color],
    },
    emphasis: {
      itemStyle: { opacity: 0.85 },
    },
  }))

  monthlyChartOptions.value = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: c.bg,
      borderColor: c.border,
      textStyle: { color: c.text },
      axisPointer: { type: 'shadow' },
    },
    legend: {
      bottom: 0,
      textStyle: { color: c.textMuted },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      top: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: {
        color: c.textMuted,
        fontSize: 11,
        rotate: labels.length > 12 ? 45 : 0,
      },
      axisLine: { lineStyle: { color: c.border } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: c.textMuted, fontSize: 11 },
      splitLine: { lineStyle: { color: c.border, type: 'dashed' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series: [
      ...barSeries,
      {
        name: t('dashboard.charts.monthly.cumulative'),
        type: 'line',
        data: cumulative,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { color: c.primary, width: 2 },
        itemStyle: { color: c.primary },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: c.primary + '20' },
              { offset: 1, color: c.primary + '05' },
            ],
          },
        },
      },
    ],
  }
}

// ── Card click handlers ─────────────────────────────────────
const listFnMap = {
  users: () => fetchAllItems(p => $SystemAPI.userList(p), { incSuspended: true, incDeleted: true }),
  roles: () => fetchAllItems(p => $SystemAPI.roleList(p), { deleted: 1, archived: 1 }),

  workflows: () => fetchAllItems(p => $AutomationAPI.workflowList(p), { deleted: 1, disabled: 1 }),
  namespaces: () => fetchAllItems(p => $ComposeAPI.namespaceList(p)),
}

// Cache fetched items per key to avoid re-fetching
const itemsCache = {}

async function switchCard(key) {
  if (activeCard.value === key) return

  activeCard.value = key
  monthlyChartOptions.value = null

  try {
    if (!itemsCache[key]) {
      itemsCache[key] = await listFnMap[key]()
    }
    buildStackedChart(key, itemsCache[key])
  } catch (e) {
    console.error('Dashboard: failed to fetch monthly data for', key, e)
  }
}

async function loadChart(key) {
  monthlyChartOptions.value = null

  try {
    if (!itemsCache[key]) {
      itemsCache[key] = await listFnMap[key]()
    }
    buildStackedChart(key, itemsCache[key])
  } catch (e) {
    console.error('Dashboard: failed to fetch chart data for', key, e)
  }
}

// ── Data fetching ───────────────────────────────────────────
async function fetchAllData() {
  loading.value = true

  try {
    const [
      usersTotal,
      usersSuspended,
      usersDeleted,
      rolesTotal,
      rolesArchived,
      rolesDeleted,
      workflowsTotal,
      workflowsDisabled,
      workflowsDeleted,
      namespacesTotal,
    ] = await Promise.all([
      // Users: default = active only (excludes suspended/deleted)
      fetchCount($SystemAPI.userListCancellable({ limit: 1, incTotal: true })),
      fetchCount($SystemAPI.userListCancellable({ limit: 1, incTotal: true, suspended: 2 })),
      fetchCount($SystemAPI.userListCancellable({ limit: 1, incTotal: true, deleted: 2 })),
      // Roles: default = active only (excludes archived/deleted)
      fetchCount($SystemAPI.roleListCancellable({ limit: 1, incTotal: true })),
      fetchCount($SystemAPI.roleListCancellable({ limit: 1, incTotal: true, archived: 2 })),
      fetchCount($SystemAPI.roleListCancellable({ limit: 1, incTotal: true, deleted: 2 })),
      // Workflows: default = active only (excludes disabled/deleted)
      fetchCount($AutomationAPI.workflowListCancellable({ limit: 1, incTotal: true })),
      fetchCount($AutomationAPI.workflowListCancellable({ limit: 1, incTotal: true, disabled: 2 })),
      fetchCount($AutomationAPI.workflowListCancellable({ limit: 1, incTotal: true, deleted: 2 })),
      // Namespaces
      fetchCount($ComposeAPI.namespaceListCancellable({ limit: 1, incTotal: true })),
    ])

    // Default API returns active-only, so sum up for true total
    activeUsers.value = usersTotal
    suspendedUsers.value = usersSuspended
    deletedUsers.value = usersDeleted
    totalUsers.value = usersTotal + usersSuspended + usersDeleted

    activeRoles.value = rolesTotal
    archivedRoles.value = rolesArchived
    deletedRoles.value = rolesDeleted
    totalRoles.value = rolesTotal + rolesArchived + rolesDeleted

    enabledWorkflows.value = workflowsTotal
    disabledWorkflows.value = workflowsDisabled
    deletedWorkflows.value = workflowsDeleted
    totalWorkflows.value = workflowsTotal + workflowsDisabled + workflowsDeleted

    totalNamespaces.value = namespacesTotal
  } catch (e) {
    console.error('Dashboard: failed to fetch counts', e)
  }

  loading.value = false

  // Load user chart by default
  loadChart('users')
}

// ── Init ────────────────────────────────────────────────────
onMounted(() => {
  fetchAllData()
})
</script>
