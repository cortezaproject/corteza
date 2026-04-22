<template>
  <InputGroup>
    <InputText
      type="number"
      :model-value="numericValue"
      :min="min"
      :max="max"
      :step="step"
      @update:model-value="onUpdate"
    />
    <InputGroupAddon>px</InputGroupAddon>
  </InputGroup>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  min: { type: Number, default: 0 },
  max: { type: Number, default: 999 },
  step: { type: Number, default: 1 },
})

const emit = defineEmits(['update:modelValue'])

const numericValue = computed(() => {
  const n = parseInt(props.modelValue, 10)
  return Number.isFinite(n) ? n : null
})

function onUpdate(v) {
  if (v === '' || v === null || v === undefined) {
    emit('update:modelValue', '')
    return
  }
  emit('update:modelValue', `${v}px`)
}
</script>
