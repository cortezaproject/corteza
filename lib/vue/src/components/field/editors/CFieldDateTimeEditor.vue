<template>
  <DatePicker
    :model-value="dateValue"
    :show-time="showTime"
    :time-only="timeOnly"
    :min-date="minDate"
    :max-date="maxDate"
    :disabled="disabled"
    date-format="yy-mm-dd"
    show-icon
    fluid
    @update:model-value="onUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const timeOnly = computed(() => !!props.field.options?.onlyTime)
const showTime = computed(() => {
  if (props.field.options?.onlyDate) return false
  return true
})

const minDate = computed(() => {
  if (props.field.options?.onlyFutureValues) return new Date()
  return undefined
})

const maxDate = computed(() => {
  if (props.field.options?.onlyPastValues) return new Date()
  return undefined
})

const dateValue = computed(() => {
  if (!props.modelValue) return null
  if (timeOnly.value) {
    const [h, m, s] = props.modelValue.split(':').map(Number)
    const d = new Date()
    d.setHours(h || 0, m || 0, s || 0, 0)
    return d
  }
  // Append T00:00:00 to date-only strings to ensure local timezone parsing
  const raw = props.modelValue
  const d = /^\d{4}-\d{2}-\d{2}$/.test(raw) ? new Date(raw + 'T00:00:00') : new Date(raw)
  return isNaN(d.getTime()) ? null : d
})

function pad(n) {
  return String(n).padStart(2, '0')
}

function onUpdate(value) {
  if (!value) {
    emit('update:modelValue', '')
    return
  }

  if (timeOnly.value) {
    emit(
      'update:modelValue',
      `${pad(value.getHours())}:${pad(value.getMinutes())}:${pad(value.getSeconds())}`,
    )
    return
  }

  if (props.field.options?.onlyDate) {
    emit(
      'update:modelValue',
      `${value.getFullYear()}-${pad(value.getMonth() + 1)}-${pad(value.getDate())}`,
    )
    return
  }

  emit('update:modelValue', value.toISOString())
}
</script>
