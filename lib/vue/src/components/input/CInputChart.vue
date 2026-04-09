<template>
  <Select
    :model-value="selectedChart"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    data-key="chartID"
    :placeholder="placeholder"
    :disabled="disabled || !namespaceID"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['name', 'handle', 'chartID']"
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      <span>{{ option.name || option.handle || option.chartID }}</span>
    </template>
  </Select>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useComposeResourceStore } from '../../stores/useComposeResourceStore'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: null,
  },
  namespaceID: {
    type: [String, Number],
    default: null,
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const store = useComposeResourceStore()

const options = ref([])
const selectedChart = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(chart) {
  if (!chart) return ''
  return chart.name || chart.handle || chart.chartID
}

async function fetchCharts() {
  if (!props.namespaceID) {
    options.value = []
    return
  }

  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = store.searchCharts(props.namespaceID, {
      query: '',
      limit: 100,
      sort: 'name ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = result.set || []
  } catch (e) {
    // Ignore cancelled requests
    if (e?.message !== 'canceled') {
      options.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0 && props.namespaceID) {
    fetchCharts()
  }
}

function onSelect(value) {
  selectedChart.value = value
  emit('update:modelValue', value?.chartID || null)
}

async function loadChartById(chartID) {
  if (!chartID || !props.namespaceID) return

  // First check if already in options
  const existing = options.value.find(m => m.chartID === chartID)
  if (existing) {
    selectedChart.value = existing
    return
  }

  // Resolve through store (cache-first)
  loading.value = true
  try {
    const chart = await store.resolveChart(props.namespaceID, chartID)
    if (chart) {
      selectedChart.value = chart
      if (!options.value.find(m => m.chartID === chartID)) {
        options.value = [...options.value, chart]
      }
    }
  } catch (_e) {
    // Chart not found or API error
  } finally {
    loading.value = false
  }
}

// Watch for namespace changes - clear selection and reload
watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      selectedChart.value = null
      options.value = []
      emit('update:modelValue', null)
    }
    if (newVal) {
      fetchCharts()
    }
  },
)

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedChart.value || selectedChart.value.chartID !== newVal)) {
      loadChartById(newVal)
    } else if (!newVal) {
      selectedChart.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.namespaceID) {
    fetchCharts()
  }
  if (props.modelValue && props.namespaceID) {
    loadChartById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
