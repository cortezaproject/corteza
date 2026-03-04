<template>
  <CInputUser
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import CInputUser from '../../input/CInputUser.vue'

defineProps({
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

defineEmits(['update:modelValue'])
</script>
