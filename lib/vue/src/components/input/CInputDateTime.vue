<template>
  <DatePicker
    :model-value="dateValue"
    :show-time="withTime"
    :time-only="timeOnly"
    :min-date="minDate"
    :max-date="maxDate"
    :disabled="disabled"
    :invalid="invalid"
    :placeholder="placeholder"
    hour-format="24"
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
  // Accepts either a string (the default contract: ISO for a datetime,
  // YYYY-MM-DD when date-only, HH:mm:ss when time-only) or a Date. What comes
  // back out is decided by `valueType`, not by what went in — a caller holding
  // Date objects sets valueType="date" and keeps them.
  modelValue: {
    type: [String, Date],
    default: '',
  },
  // 'string' (default) | 'date'. PrimeVue's DatePicker only understands Date
  // objects: handed a string it renders it through a different code path that
  // drops the time, and with show-time it throws outright. Every caller goes
  // through this component so that conversion happens in exactly one place.
  valueType: {
    type: String,
    default: 'string',
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
  invalid: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: undefined,
  },
  onlyDate: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

// A date-only field has no time to show, whichever way the caller left
// show-time — CFieldDateTimeEditor derives the same thing from field options.
const withTime = computed(() => props.showTime && !props.onlyDate)

const dateValue = computed(() => {
  const raw = props.modelValue
  if (!raw) return null
  if (raw instanceof Date) return isNaN(raw.getTime()) ? null : raw
  if (props.timeOnly) {
    const [h, m, s] = raw.split(':').map(Number)
    const d = new Date()
    d.setHours(h || 0, m || 0, s || 0, 0)
    return d
  }
  // Append T00:00:00 to date-only strings to ensure local timezone parsing
  const d = /^\d{4}-\d{2}-\d{2}$/.test(raw) ? new Date(raw + 'T00:00:00') : new Date(raw)
  return isNaN(d.getTime()) ? null : d
})

function pad(n) {
  return String(n).padStart(2, '0')
}

function onUpdate(value) {
  if (!value) {
    emit('update:modelValue', props.valueType === 'date' ? null : '')
    return
  }

  if (props.valueType === 'date') {
    emit('update:modelValue', value)
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
