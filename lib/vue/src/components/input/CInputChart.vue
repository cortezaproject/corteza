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
import { computed, onMounted, ref, watch } from 'vue'
import { useChartStore } from '../../stores/useChartStore'

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

const store = useChartStore()
const selectedChart = ref(null)
const loading = ref(false)

const options = computed(() => store.set)

function getOptionLabel(chart) {
  if (!chart) return ''
  return chart.name || chart.handle || chart.chartID
}

async function fetchCharts() {
  if (!props.namespaceID) return

  loading.value = true
  try {
    await store.loadFor(props.namespaceID)
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function onShow() {
  if (props.namespaceID) {
    fetchCharts()
  }
}

function onSelect(value) {
  selectedChart.value = value
  emit('update:modelValue', value?.chartID || null)
}

async function loadChartById(chartID) {
  if (!chartID || !props.namespaceID) return

  const existing = options.value.find(c => c.chartID === chartID)
  if (existing) {
    selectedChart.value = existing
    return
  }

  loading.value = true
  try {
    const chart = await store.findByID({ namespaceID: props.namespaceID, chartID })
    if (chart) {
      selectedChart.value = store.set.find(c => c.chartID === chartID) || chart
    }
  } catch {
    // chart not found
  } finally {
    loading.value = false
  }
}

watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      selectedChart.value = null
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

watch(
  () => options.value.length,
  () => {
    if (props.modelValue && options.value.length) {
      const found = options.value.find(
        c => c.chartID === props.modelValue || c.chartID === String(props.modelValue),
      )
      if (found) selectedChart.value = found
    }
  },
)

onMounted(() => {
  if (props.namespaceID) {
    fetchCharts()
  }
  if (props.modelValue && props.namespaceID) {
    loadChartById(props.modelValue)
  }
})
</script>
