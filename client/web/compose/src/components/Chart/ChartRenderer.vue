<template>
  <div class="relative h-full w-full">
    <div v-if="processing" class="absolute inset-0 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <div v-else-if="error" class="absolute inset-0 p-3 text-red-500">
      {{ error }}
    </div>

    <CChart v-else-if="renderer" :chart="renderer" class="absolute inset-0 p-1" />
  </div>
</template>

<script setup>
import { ref, watch, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'
import { chartConstructor } from '../../lib/charts'
import { useModuleStore } from '../../stores/module'

const { CChart } = components

const props = defineProps({
  chart: {
    type: Object,
    required: true,
  },
  reporter: {
    type: Function,
    required: true,
  },
  record: {
    type: Object,
    required: false,
    default: undefined,
  },
})

const emit = defineEmits(['updated'])

const { t } = useI18n()
const moduleStore = useModuleStore()

const error = ref(undefined)
const processing = ref(false)
const renderer = ref(undefined)

async function updateChart() {
  error.value = undefined
  renderer.value = undefined

  const [report = {}] = props.chart.config.reports

  if (!report.moduleID) {
    return
  }

  processing.value = true

  const chart = chartConstructor(props.chart)

  try {
    chart.isValid()

    const data = await chart.fetchReports({ reporter: props.reporter })

    const module = moduleStore.getByID(report.moduleID)

    if (module) {
      const fields = [...module.fields, ...(module.systemFields ? module.systemFields() : [])]

      if (data.labels && Array.isArray(data.labels)) {
        const [dimension = {}] = report.dimensions
        let { field } = dimension

        const fieldObj = fields.find(({ name }) => name === field)

        if (fieldObj) {
          if (fieldObj.kind === 'Bool') {
            const { trueLabel, falseLabel } = fieldObj.options || {}
            data.labels = data.labels.map(value => {
              return value === '1'
                ? trueLabel || t('general.label.yes')
                : falseLabel || t('general.label.no')
            })
          } else if (fieldObj.kind === 'Select') {
            data.labels = data.labels.map(value => {
              const found = (fieldObj.options?.options || []).find(o => o.value === value)
              const text = found?.text
              return text || value
            })
          }
        }
      }

      // Map dataset labels to field labels
      data.datasets = (data.datasets || []).map(dataset => {
        if (!dataset) dataset = {}
        const { label } = dataset

        if (label === 'count') {
          dataset.label = t('chart.general.label.count')
        } else {
          const matchedField = fields.find(({ name }) => name === label)
          dataset.label = matchedField?.label || label
        }

        return dataset
      })

      // Replace 'undefined' labels
      if (data.labels) {
        data.labels = data.labels.map(l => (l === 'undefined' ? t('chart.undefined') : l))
      }
    }

    // Add theme variables for chart styling
    data.themeVariables = getThemeVariables()

    renderer.value = chart.makeOptions(data)
  } catch (e) {
    error.value = (e instanceof Error ? e.message : String(e)) || t('chart.notification.loadFailed')
    processing.value = false
    return
  }

  setTimeout(() => {
    processing.value = false
    emit('updated')
  }, 300)
}

function getThemeVariables() {
  const getCssVariable = variableName => {
    return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
  }

  return {
    white: getCssVariable('--p-surface-0') || '#ffffff',
    black: getCssVariable('--p-text-color') || '#333333',
    primary: getCssVariable('--p-primary-color') || '#3B82F6',
    secondary: getCssVariable('--p-text-muted-color') || '#6B7280',
    success: getCssVariable('--p-green-500') || '#22C55E',
    warning: getCssVariable('--p-yellow-500') || '#F59E0B',
    danger: getCssVariable('--p-red-500') || '#EF4444',
    light: getCssVariable('--p-surface-200') || '#E5E7EB',
    'extra-light': getCssVariable('--p-surface-100') || '#F3F4F6',
    dark: getCssVariable('--p-surface-700') || '#374151',
    'font-regular': getCssVariable('--p-font-family') || 'inherit',
  }
}

function setDefaultValues() {
  processing.value = false
  renderer.value = undefined
}

// Watch chart itself for changes (immediate to trigger first render)
watch(
  () => props.chart,
  () => {
    updateChart()
  },
  { immediate: true, deep: true },
)

// Watch for record changes to re-render (not immediate — chart watcher handles mount)
watch(
  () => props.record?.recordID,
  () => {
    updateChart()
  },
)

onBeforeUnmount(() => {
  setDefaultValues()
})

defineExpose({
  updateChart,
})
</script>
