<template>
  <c-input-date-time
    :model-value="modelValue"
    :show-time="showTime"
    :time-only="timeOnly"
    :min-date="minDate"
    :max-date="maxDate"
    :only-date="onlyDate"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import { CInputDateTime } from '../../input'

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

const onlyDate = computed(() => !!props.field.options?.onlyDate)
const timeOnly = computed(() => !!props.field.options?.onlyTime)
const showTime = computed(() => !onlyDate.value)

const minDate = computed(() => {
  if (props.field.options?.onlyFutureValues) return new Date()
  return undefined
})

const maxDate = computed(() => {
  if (props.field.options?.onlyPastValues) return new Date()
  return undefined
})
</script>
