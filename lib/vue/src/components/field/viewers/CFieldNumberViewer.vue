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
          :class="progressBarClasses"
        />
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

const progressBarClasses = computed(() => {
  const classes = []
  if (props.field.options?.striped) classes.push('striped')
  if (props.field.options?.animated) classes.push('animated')
  return classes
})

const viewerClasses = computed(() => {
  return props.field.isMulti ? ['multiline'] : []
})
</script>

<style scoped>
.multiline {
  white-space: pre-line;
}
</style>
