<template>
  <div>
    <!-- Number display mode -->
    <div v-if="isNumberDisplay" :class="viewerClasses">
      {{ formattedNumber }}
    </div>

    <!-- Progress bar display mode -->
    <template v-else>
      <div
        v-for="(v, i) in progressValues"
        :key="i"
        :class="['relative', { 'mt-2': i }]"
        style="min-width: 15rem; height: 1.5rem"
      >
        <ProgressBar
          :value="normalizedProgress(v)"
          :show-value="false"
          :style="progressStyle(v)"
          :class="progressBarClasses(v)"
          class="w-full h-full"
        />
        <div
          v-if="field.options.showValue !== false"
          class="absolute inset-0 flex items-center justify-center pointer-events-none"
        >
          <span class="text-sm font-medium" :style="progressTextStyle(v)">
            {{ progressLabel(v) }}
          </span>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ProgressBar from 'primevue/progressbar'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const value = computed(() => {
  if (props.field.isSystem) {
    return props.record[props.field.name]
  }
  return props.record?.values?.[props.field.name]
})

const isNumberDisplay = computed(() => {
  return props.field.options?.display === 'number' || !props.field.options?.display
})

// Format value using the field's formatValue method if available
const formattedNumber = computed(() => {
  const v = value.value
  if (v === undefined || v === null) return ''

  const vals = props.field.isMulti ? v : [v]
  const delimiter = props.field.options?.multiDelimiter || ', '

  if (props.field.formatValue) {
    return vals.map(val => props.field.formatValue(val)).join(delimiter)
  }

  // Fallback formatting
  const { precision, prefix = '', suffix = '' } = props.field.options || {}

  return vals
    .map(val => {
      const num = Number(val)
      if (isNaN(num)) return val

      let out = precision !== undefined ? num.toFixed(Number(precision)) : num.toLocaleString()

      return `${prefix}${out}${suffix}`
    })
    .join(delimiter)
})

// Progress bar values
const progressValues = computed(() => {
  const v = value.value
  const min = parseFloat(props.field.options?.min || 0)

  if (v === undefined || v === null) {
    return [min]
  }

  const vals = props.field.isMulti ? v : [v]
  return vals.length ? vals : [min]
})

function normalizedProgress(v) {
  const min = parseFloat(props.field.options?.min || 0)
  const max = parseFloat(props.field.options?.max || 100)
  const num = parseFloat(v) || 0

  if (max === min) return 0
  return Math.round(((num - min) / (max - min)) * 100)
}

// Map progress variant → PrimeVue Button severity token suffix
// Reading the same --p-button-{severity}-* tokens the Button component uses guarantees
// the progress fill stays in sync with the configurator's Button tags across themes.
const severityTokenMap = {
  primary: 'primary',
  secondary: 'secondary',
  success: 'success',
  warning: 'warn',
  danger: 'danger',
  info: 'info',
  dark: 'contrast',
}

// Determine variant based on thresholds and current value percentage
function getProgressVariant(v) {
  const thresholds = props.field.options?.thresholds || []
  const baseVariant = props.field.options?.variant || 'success'

  if (!thresholds.length) return baseVariant

  const pct = normalizedProgress(v)
  const sorted = [...thresholds]
    .filter(t => t.value >= 0)
    .sort((a, b) => b.value - a.value)

  for (const t of sorted) {
    if (pct >= t.value) return t.variant || baseVariant
  }

  return baseVariant
}

// Set the PrimeVue ProgressBar's value-background design token so the fill matches
// the configurator's Button tag across themes (same approach as the block ProgressBlock).
function progressStyle(v) {
  const sev = severityTokenMap[getProgressVariant(v)] || 'primary'
  return {
    '--p-progressbar-value-background': `var(--p-button-${sev}-background)`,
  }
}

// Progress bar label (mirrors Human's CProgress.progressLabel)
function progressLabel(v) {
  const opts = props.field.options || {}
  const min = parseFloat(opts.min || 0)
  const max = parseFloat(opts.max || 100)
  const num = parseFloat(v) || 0
  const maxRange = Math.abs(max - min)
  const progressValue = Math.abs(num - min)

  if (opts.showValue === false) return ''

  let label = String(num)

  if (opts.showRelative) {
    label = maxRange === 0 ? '0%' : `${Math.round((progressValue / maxRange) * 10000) / 100}%`
  }

  if (opts.showProgress) {
    const maxLabel = opts.showRelative ? '100%' : String(max)
    label = `${label} / ${maxLabel}`
  }

  return label
}

// CSS classes for the progress bar
function progressBarClasses(_v) {
  const classes = []
  if (props.field.options?.animated) classes.push('animated')
  return classes
}

// Text color reads the matching --p-button-{severity}-color so it stays in sync with progress fill
function progressTextStyle(v) {
  const sev = severityTokenMap[getProgressVariant(v)] || 'primary'
  return { color: `var(--p-button-${sev}-color)` }
}

const viewerClasses = computed(() => {
  return props.field.isMulti ? ['multiline'] : []
})
</script>

<style scoped>
.multiline {
  white-space: pre-line;
}
</style>
