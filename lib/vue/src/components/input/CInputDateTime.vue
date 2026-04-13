<template>
  <DatePicker
    :model-value="dateValue"
    :show-time="showTime"
    :time-only="timeOnly"
    :min-date="minDate"
    :max-date="maxDate"
    :disabled="disabled"
    show-icon
    fluid
    :size="size"
    show-button-bar
    @update:model-value="onUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  size: {
    type: String,
    default: '',
  },
  showTime: {
    type: Boolean,
    default: true,
  },
  timeOnly: {
    type: Boolean,
    default: false,
  },
  minDate: {
    type: Date,
    default: undefined,
  },
  maxDate: {
    type: Date,
    default: undefined,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  onlyDate: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const dateValue = computed(() => {
  if (!props.modelValue) return null
  if (props.timeOnly) {
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

  if (props.timeOnly) {
    emit(
      'update:modelValue',
      `${pad(value.getHours())}:${pad(value.getMinutes())}:${pad(value.getSeconds())}`,
    )
    return
  }

  if (props.onlyDate) {
    emit(
      'update:modelValue',
      `${value.getFullYear()}-${pad(value.getMonth() + 1)}-${pad(value.getDate())}`,
    )
    return
  }

  emit('update:modelValue', value.toISOString())
}
</script>
