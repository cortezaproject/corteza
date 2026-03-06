<template>
  <PageBlock :block="block">
    <div class="flex items-center justify-center h-full">
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
          v-for="(m, mi) in options.metrics"
          :key="mi"
          class="flex items-center justify-center overflow-hidden h-full w-full"
        >
          <div
            v-for="(v, vi) in formatResponse(m, mi)"
            :key="vi"
            class="w-full h-full px-2 py-1"
          >
            <MetricItem
              :metric="m"
              :value="v"
            />
          </div>
        </div>

        <div
          v-if="!options.metrics.length"
          class="text-muted-color italic"
        >
          {{ $t('block.metric.defaultMetricLabel') }}
        </div>
      </template>
    </div>
  </PageBlock>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, inject } from 'vue'
import numeral from 'numeral'
import PageBlock from './PageBlock.vue'
import MetricItem from './Metric/MetricItem.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')

const processing = ref(false)
const error = ref(undefined)
const reports = ref([])

const options = computed(() => props.block.options || {})

/**
 * Performs post-processing on the provided data — formats numbers.
 */
function formatResponse (m, i) {
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
 * Pulls fresh data from the API using PageBlockMetric.fetch().
 */
async function refresh () {
  if (!$ComposeAPI) return

  error.value = undefined
  processing.value = true

  try {
    const rtr = []
    const namespaceID = props.namespace.namespaceID

    const reporter = (r) => {
      return $ComposeAPI.recordReport({ ...r, namespaceID })
        .then(({ report }) => report || [])
    }

    for (const m of options.value.metrics || []) {
      if (m.moduleID) {
        const vals = await props.block.fetch({ m }, reporter)
        rtr.push(vals)
      }
    }

    reports.value = rtr

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

onMounted(() => {
  refresh()
})

onBeforeUnmount(() => {
  reports.value = []
})
</script>
