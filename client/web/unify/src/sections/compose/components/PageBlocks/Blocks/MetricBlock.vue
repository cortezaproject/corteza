<template>
  <PageBlock :block="block" @refreshBlock="refresh">
    <div class="flex flex-col h-full">
      <template v-if="processing">
        <ProgressSpinner style="width: 24px; height: 24px" />
      </template>

      <template v-else-if="error">
        <div class="text-muted-color italic p-3">
          {{ error }}
        </div>
      </template>

      <template v-else>
        <div
          v-for="(m, mi) in configuredMetrics"
          :key="mi"
          class="flex w-full flex-1 min-h-0 overflow-hidden"
        >
          <div v-for="(v, vi) in formatResponse(m.metric, m.index)" :key="vi" class="w-full h-full">
            <MetricItem
              :metric="m.metric"
              :value="v"
              :change-value="m.metric.comparison?.enabled ? (changeValues[m.index] ?? null) : null"
              @drill-down="onMetricDrillDown(m.metric, v)"
            />
          </div>
        </div>

        <div v-if="!options.metrics.length" class="text-muted-color italic">
          {{ $t('block.metric.defaultMetricLabel') }}
        </div>
      </template>
    </div>

    <!-- Drill down modal -->
    <Dialog
      v-if="drillDownTargetBlock"
      v-model:visible="drillDownModalVisible"
      :header="drillDownModalTitle"
      modal
      dismissableMask
      :style="{ width: '90vw', height: '90vh' }"
      :pt="{
        content: { class: 'flex-1 flex flex-col overflow-hidden p-0 bg-surface h-full' },
      }"
    >
      <RecordListBlock
        :block="drillDownTargetBlock"
        :namespace="namespace"
        :page="page"
        :record="record"
        class="h-full rounded-none border-0 shadow-none border-t"
      />
    </Dialog>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, inject, defineAsyncComponent } from 'vue'
import { compose } from '@planetcrust/human-js'
import numeral from 'numeral'
import PageBlock from './PageBlock.vue'
import MetricItem from './Metric/MetricItem.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'
const RecordListBlock = defineAsyncComponent(() => import('./RecordListBlock.vue'))

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')
const $Auth = inject('$Auth', {})
const $eventBus = inject('$eventBus', null)

const processing = ref(false)
const error = ref(undefined)
const reports = ref([])
const changeValues = ref({})
const drillDownModalVisible = ref(false)
const drillDownTargetBlock = ref(null)
const drillDownModalTitle = ref('')

const options = computed(() => props.block.options || {})

const configuredMetrics = computed(() => {
  return (options.value.metrics || [])
    .map((metric, index) => ({ metric, index }))
    .filter(({ metric }) => metric.moduleID)
})

/**
 * Performs post-processing on the provided data — formats numbers.
 */
function formatResponse(m, i) {
  const vals = reports.value[i]
  if (!vals) return []

  return vals.map(({ label, value }) => {
    if (m.numberFormat) {
      value = numeral(value).format(m.numberFormat)
    }

    return { label, value }
  })
}

/**
 * Builds a date cutoff filter: records created before N periods ago.
 * This gives us "the total value as of N periods ago".
 */
function buildPreviousCutoffFilter(period) {
  const now = new Date()
  const cutoff = new Date(now)

  switch (period) {
    case 'day':
      cutoff.setDate(cutoff.getDate() - 1)
      break
    case 'week':
      cutoff.setDate(cutoff.getDate() - 7)
      break
    case 'month':
      cutoff.setMonth(cutoff.getMonth() - 1)
      break
    case 'quarter':
      cutoff.setMonth(cutoff.getMonth() - 3)
      break
    case 'year':
      cutoff.setFullYear(cutoff.getFullYear() - 1)
      break
    default:
      return ''
  }

  const fmt = d => d.toISOString().replace('T', ' ').substring(0, 19)
  return `createdAt < '${fmt(cutoff)}'`
}

/**
 * Combines a base filter with an additional date filter using AND.
 */
function combineFilters(baseFilter, dateFilter) {
  if (!dateFilter) return baseFilter || ''
  if (!baseFilter) return dateFilter
  return `(${baseFilter}) AND ${dateFilter}`
}

/**
 * Extracts the raw numeric value from a fetch result.
 */
function extractValue(fetchResult) {
  if (!fetchResult || !fetchResult.length) return 0
  return Number(fetchResult[0].value) || 0
}

/**
 * Computes the percentage change between current and previous period.
 * Always returns a number so the badge is visible when comparison is enabled.
 */
function computeChange(current, previous) {
  if (previous === 0 && current === 0) return 0
  if (previous === 0) return 100 // something from nothing → +100%

  const change = ((current - previous) / Math.abs(previous)) * 100

  if (!isFinite(change) || isNaN(change)) return 0

  return change
}

function usesRecordVars(filter) {
  return !!filter && (filter.includes('${record') || filter.includes('${ownerID}'))
}

function interpolateFilter(filter) {
  if (!filter) return filter

  const record = props.record
  const user = $Auth?.user || {}

  return evaluatePrefilter(filter, {
    record,
    user,
    recordID: record?.recordID || '0',
    ownerID: record?.ownedBy || '0',
    userID: user?.userID || '0',
  })
}

/**
 * Pulls fresh data from the API using PageBlockMetric.fetch().
 */
async function refresh() {
  if (!$ComposeAPI) return

  error.value = undefined
  processing.value = true

  try {
    const rtr = []
    const changes = {}
    const namespaceID = props.namespace.namespaceID

    const reporter = r => {
      return $ComposeAPI.recordReport({ ...r, namespaceID }).then(result => result || [])
    }

    const metrics = options.value.metrics || []

    for (let mi = 0; mi < metrics.length; mi++) {
      const m = metrics[mi]
      if (!m.moduleID) continue

      const customFilter = m.comparison?.customFilter

      if (!props.record && (usesRecordVars(m.filter) || usesRecordVars(customFilter))) {
        console.warn('Skipping metric: filter uses record variables outside a record page')
        continue
      }

      // Evaluated copies — reused for both the current and comparison fetch so
      // interpolation (e.g. ${recordID}) runs once against the same values.
      let evaluatedMetric
      let evaluatedCustomFilter
      try {
        evaluatedMetric = { ...m, filter: interpolateFilter(m.filter) }
        evaluatedCustomFilter = interpolateFilter(customFilter)
      } catch (e) {
        console.warn('Skipping metric: filter interpolation failed', mi, e)
        continue
      }

      // Fetch current value
      const vals = await props.block.fetch({ m: evaluatedMetric }, reporter)
      // Keyed by the configured metric index — formatResponse() looks values up
      // by that index, and skipped metrics must not shift the ones after them.
      rtr[mi] = vals

      // Fetch comparison if enabled
      if (m.comparison?.enabled && m.comparison?.period) {
        try {
          // Current value is the all-time total (already fetched above)
          const currentValue = extractValue(vals)

          // Previous value = total as of N periods ago (records before the cutoff)
          const cutoffFilter = buildPreviousCutoffFilter(m.comparison.period)
          const prevFilter = evaluatedCustomFilter
            ? combineFilters(evaluatedCustomFilter, cutoffFilter)
            : combineFilters(evaluatedMetric.filter, cutoffFilter)

          const prevMetric = { ...evaluatedMetric, filter: prevFilter }
          const prevVals = await props.block.fetch({ m: prevMetric }, reporter)
          const previousValue = extractValue(prevVals)

          console.debug('[Metric comparison]', {
            period: m.comparison.period,
            currentValue,
            previousValue,
            prevFilter,
          })

          changes[mi] = computeChange(currentValue, previousValue)
        } catch (e) {
          console.warn('Comparison fetch failed for metric', mi, e)
          changes[mi] = null
        }
      }
    }

    reports.value = rtr
    changeValues.value = changes

    setTimeout(() => {
      processing.value = false
    }, 300)
  } catch (e) {
    console.error('Metric fetch error:', e)
    error.value = e.message || String(e)

    setTimeout(() => {
      processing.value = false
    }, 300)
  }
}

watch(
  () => props.block.options,
  () => refresh(),
  { deep: true },
)

onMounted(() => {
  refresh()
})

onBeforeUnmount(() => {
  reports.value = []
  changeValues.value = {}
  offRefetch?.()
})

const offRefetch = $eventBus?.on('refetch-records', () => refresh())

function onMetricDrillDown(metric, value) {
  drillDown(metric, value)
}

function drillDown(metric, value) {
  const drillDownOpts = metric.drillDown || {}
  if (!drillDownOpts.enabled) return

  const drillDownValue = metric.label || value?.label || value?.value || ''
  // Raw filter — the child RecordListBlock evaluates the prefilter itself.
  const prefilter = metric.filter ? `(${metric.filter})` : ''

  if (drillDownOpts.blockID) {
    $eventBus?.emit(`drill-down-recordList:${drillDownOpts.blockID}`, {
      prefilter,
      name: drillDownValue,
      value: drillDownValue,
    })
    return
  }

  const title = props.block.title || metric.label
  drillDownModalTitle.value = title ? `${title} - "${drillDownValue}"` : drillDownValue

  const fields = drillDownOpts.recordListOptions?.fields || []

  drillDownTargetBlock.value = new compose.PageBlockRecordList({
    blockID: `drillDown-${props.block.blockID}-${metric.moduleID}`,
    options: {
      moduleID: metric.moduleID,
      fields,
      prefilter,
      presort: 'createdAt DESC',
      hideRecordReminderButton: true,
      hideRecordViewButton: false,
      hideConfigureFieldsButton: false,
      hideImportButton: true,
      enableRecordPageNavigation: true,
      selectable: true,
      allowExport: true,
      perPage: 14,
      showTotalCount: true,
      recordDisplayOption: 'modal',
    },
    style: {
      wrap: { kind: 'plain' },
    },
  })

  drillDownModalVisible.value = true
}
</script>
