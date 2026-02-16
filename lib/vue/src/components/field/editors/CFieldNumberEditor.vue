<template>
  <InputNumber
    :model-value="numericValue"
    :disabled="disabled"
    :min-fraction-digits="0"
    :max-fraction-digits="field.options?.precision ?? 3"
    :min="field.options?.min"
    :max="field.options?.max"
    :step="field.options?.step ?? 1"
    :prefix="field.options?.prefix"
    :suffix="field.options?.suffix"
    class="w-full"
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

const numericValue = computed(() => {
  if (props.modelValue === '' || props.modelValue == null) return null
  const n = parseFloat(props.modelValue)
  return isNaN(n) ? null : n
})

function onUpdate(value) {
  emit('update:modelValue', value == null ? '' : String(value))
}
</script>
