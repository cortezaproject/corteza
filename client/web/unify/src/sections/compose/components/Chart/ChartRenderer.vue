<template>
  <div class="relative h-full w-full">
    <!-- The chart outlives a refetch: it stays mounted while new data is on the
         way, so editing a setting updates the plot in place instead of tearing
         it down and building it again. The spinner is only for having nothing
         to show yet. -->
    <CChart
      v-if="renderer"
      ref="chartRef"
      :chart="renderer"
      class="absolute inset-0 p-1"
      @click="handleChartClick"
      @rendered="fitLabels"
    />

    <div v-if="processing && !renderer" class="absolute inset-0 flex items-center justify-center">
      <ProgressSpinner />
    </div>

    <div v-if="error" class="absolute inset-0 p-3 text-red-500 bg-surface-0 dark:bg-surface-900">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, watch, inject, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { getInstanceByDom } from 'echarts/core'
import { components } from '@planetcrust/human-vue'
import { chartConstructor, wrapLabel } from '../../lib/charts'
import { readColorSchemes } from '../../lib/chart-color-schemes'
import { useModuleStore, useRecordStore } from '@planetcrust/human-vue'

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

const emit = defineEmits(['updated', 'drill-down'])

const { t } = useI18n()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()
const $Settings = inject('$Settings', undefined)

const error = ref(undefined)
const processing = ref(false)
const renderer = ref(undefined)
const chartRef = ref(null)
// Per category axis: the width its labels were wrapped to, and the formatter doing it
const fittedLabels = new Map()
let measureContext
const valueMap = ref(new Map())
// The chart whose configured animation has already played. Held per chart, so
// opening a different one animates it in rather than inheriting the last one's
// spent animation.
const animatedFor = ref(undefined)

// A Record field holds record IDs, so a chart grouped by one labelled its
// legend with raw IDs. Record lists resolve the same references through the
// record store, and this mirrors that: the field's labelField option names the
// field to show, falling back to the referenced module's first field, exactly
// as CFieldRecordViewer picks it.
async function resolveRecordLabels(labels, fieldObj) {
  const moduleID = fieldObj.options?.moduleID
  const namespaceID = props.chart?.namespaceID

  if (!moduleID || !namespaceID) return labels

  // Dimension values can carry a placeholder ('undefined') for records that
  // have no value for the field; only real IDs are worth a lookup.
  const recordIDs = labels.filter(v => /^\d+$/.test(String(v)))
  if (!recordIDs.length) return labels

  // The referenced module is often one the page never loaded — the chart's own
  // module is in the store, its neighbours need not be.
  const refModule =
    moduleStore.getByID(moduleID) ||
    (await moduleStore.findByID({ namespaceID, moduleID }).catch(() => undefined))

  const labelField = fieldObj.options?.labelField || refModule?.fields?.[0]?.name
  if (!labelField) return labels

  await recordStore.resolveRecordLabels({ namespaceID, moduleID, recordIDs }).catch(() => {})

  return labels.map(value => {
    const label = readRecordValue(value, labelField) ?? value
    valueMap.value.set(label, value)
    return label
  })
}

// The store keeps records in two shapes: compose.Record, whose values are an
// object, and the raw API record the label cache holds when the module was not
// loaded, whose values are [{ name, value }]. A label has to read from either.
function readRecordValue(recordID, labelField) {
  const { values } = recordStore.getByID(recordID) || {}
  if (!values) return undefined

  const value = Array.isArray(values)
    ? values.find(v => v?.name === labelField)?.value
    : values[labelField]

  // A multi-value label field resolves to its first value rather than "[object]".
  return Array.isArray(value) ? value[0] : value
}

async function updateChart() {
  error.value = undefined

  const [report = {}] = props.chart.config.reports

  if (!report.moduleID) {
    renderer.value = undefined
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
              const label =
                value === '1'
                  ? trueLabel || t('general.label.yes')
                  : falseLabel || t('general.label.no')
              valueMap.value.set(label, value)
              return label
            })
          } else if (fieldObj.kind === 'Select') {
            data.labels = data.labels.map(value => {
              const found = (fieldObj.options?.options || []).find(o => o.value === value)
              const text = found?.text
              const label = text || value
              valueMap.value.set(label, value)
              return label
            })
          } else if (fieldObj.kind === 'Record') {
            data.labels = await resolveRecordLabels(data.labels, fieldObj)
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

    // The built-in palettes are compiled in; the ones an admin defines live in
    // settings, and a chart naming one resolves to no palette without them.
    data.customColorSchemes = readColorSchemes($Settings)

    const options = chart.makeOptions(data)

    // The configured animation belongs to the chart arriving, not to every
    // refetch after it: an editor that replays the grow-in on each keystroke,
    // or a block that replays it per live filter, reads as the chart reloading
    // rather than updating. Later renders hand echarts the new data and let it
    // transition between states.
    const chartKey = props.chart?.chartID ?? null
    if (animatedFor.value === chartKey) {
      options.animation = false
    }

    renderer.value = options
    animatedFor.value = chartKey
  } catch (e) {
    let msg = e instanceof Error ? e.message : String(e)

    if (msg && msg.startsWith('notification.')) {
      msg = t(msg)
    }

    error.value = msg || t('notification.chart.loadFailed')
    processing.value = false
    return
  }

  setTimeout(() => {
    processing.value = false
    emit('updated')
  }, 300)
}

function handleChartClick(e) {
  const trueName = valueMap.value.get(e.name) ?? e.name
  emit('drill-down', { ...e, trueName })
}

// Labels that no longer fit where echarts put them. A category label wraps
// between words: on a horizontal axis to the width of its band, on a vertical
// one to a third of the chart, and hideOverlap only drops what still collides;
// rotated labels keep their own layout. A radar shrinks until the names beside
// it fit on the canvas.
function fitLabels() {
  const el = chartRef.value?.$el
  const instance = el && getInstanceByDom(el)
  if (!instance || instance.isDisposed()) return

  const patch = { ...fitCategoryAxes(instance), ...fitRadars(instance) }
  if (Object.keys(patch).length) instance.setOption(patch)
}

function fitCategoryAxes(instance) {
  const patch = {}

  for (const mainType of ['xAxis', 'yAxis']) {
    const axes = []
    let changed = false

    instance.getModel().eachComponent(mainType, (axisModel, i) => {
      const axisPatch = {}
      if (axisModel.axis.type === 'category' && !+axisModel.get(['axisLabel', 'rotate'])) {
        const width =
          mainType === 'xAxis'
            ? Math.floor(axisModel.axis.getBandWidth()) - 8
            : Math.floor(instance.getWidth() / 3)
        const fit = fittedLabels.get(`${mainType}${i}`)

        if (fit?.width !== width || fit.formatter !== axisModel.get(['axisLabel', 'formatter'])) {
          const font = axisModel.getModel('axisLabel').getFont()
          const formatter = value => wrapLabel(value, width, text => measureText(text, font))
          fittedLabels.set(`${mainType}${i}`, { width, formatter })
          axisPatch.axisLabel = { formatter }
          changed = true
        }
      }
      axes.push(axisPatch)
    })

    if (changed) patch[mainType] = axes
  }

  return patch
}

function fitRadars(instance) {
  const radars = []
  let changed = false

  instance.getModel().eachComponent('radar', radarModel => {
    const { cx, startAngle } = radarModel.coordinateSystem
    const names = radarModel.get('indicator').map(({ name }) => String(name))
    const font = radarModel.getModel('axisName').getFont()
    const gap = radarModel.get('axisNameGap') + 4
    const width = instance.getWidth()

    // echarts' own radius, half of half the shorter side, as the ceiling
    let radius = Math.min(width, instance.getHeight()) / 4
    names.forEach((name, i) => {
      const cos = Math.cos(startAngle + (i * Math.PI * 2) / names.length)
      if (Math.abs(cos) < 0.1) return
      const room = cos > 0 ? width - cx : cx
      radius = Math.min(radius, (room - measureText(name, font)) / Math.abs(cos) - gap)
    })
    radius = Math.max(Math.floor(radius), 20)

    const radarPatch = {}
    if (radius !== radarModel.get('radius')) {
      radarPatch.radius = radius
      changed = true
    }
    radars.push(radarPatch)
  })

  return changed ? { radar: radars } : {}
}

function measureText(text, font) {
  measureContext ??= document.createElement('canvas').getContext('2d')
  measureContext.font = font
  return measureContext.measureText(text).width
}

function getThemeVariables() {
  const getCssVariable = variableName => {
    return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
  }

  return {
    white: getCssVariable('--p-content-background') || '#ffffff',
    black: getCssVariable('--p-text-color') || '#333333',
    primary: getCssVariable('--p-primary-color') || '#3B82F6',
    secondary: getCssVariable('--p-text-muted-color') || '#6B7280',
    success: getCssVariable('--p-green-500') || '#22C55E',
    warning: getCssVariable('--p-yellow-500') || '#F59E0B',
    danger: getCssVariable('--p-red-500') || '#EF4444',
    light: getCssVariable('--p-content-border-color') || '#E5E7EB',
    'extra-light': getCssVariable('--p-content-hover-background') || '#F3F4F6',
    dark: getCssVariable('--p-text-color') || '#374151',
    'font-regular': getComputedStyle(document.body).fontFamily || 'sans-serif',
  }
}

function setDefaultValues() {
  processing.value = false
  renderer.value = undefined
  valueMap.value.clear()
  animatedFor.value = undefined
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
