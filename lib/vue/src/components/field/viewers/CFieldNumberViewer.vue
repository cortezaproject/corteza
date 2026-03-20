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
        :class="{ 'mt-2': i }"
        style="min-width: 15rem"
      >
        <ProgressBar
          :value="normalizedProgress(v)"
          :show-value="field.options.showValue !== false"
          :style="{ height: '1.5rem' }"
          :class="progressBarClasses(v)"
          :pt="progressPt(v)"
        >
          <template v-if="field.options.showValue !== false" #default>
            <span :class="progressTextClass(v)">{{ progressLabel(v) }}</span>
          </template>
        </ProgressBar>
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

// Variant color mapping for PrimeVue severity
const variantColorMap = {
  primary: 'var(--p-primary-500, #3b82f6)',
  secondary: 'var(--p-surface-500, #6b7280)',
  success: 'var(--p-green-500, #22c55e)',
  warning: 'var(--p-yellow-500, #f59e0b)',
  danger: 'var(--p-red-500, #ef4444)',
  info: 'var(--p-sky-500, #06b6d4)',
  light: 'var(--p-surface-200, #e5e7eb)',
  dark: 'var(--p-surface-800, #1f2937)',
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

// Pass-through styling for progress bar via PrimeVue's pt
function progressPt(v) {
  const variant = getProgressVariant(v)
  const bgColor = variantColorMap[variant] || variantColorMap.success

  return {
    value: {
      style: {
        backgroundColor: bgColor,
      },
    },
  }
}

// Progress bar label (mirrors Corteza's CProgress.progressLabel)
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
function progressBarClasses(v) {
  const classes = []
  if (props.field.options?.animated) classes.push('animated')
  return classes
}

// Text color for label based on variant darkness
function progressTextClass(v) {
  const variant = getProgressVariant(v)
  return ['dark', 'primary'].includes(variant) ? 'text-white' : 'text-color'
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
