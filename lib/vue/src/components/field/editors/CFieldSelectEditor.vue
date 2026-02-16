<template>
  <Select
    :model-value="modelValue"
    :options="selectOptions"
    option-label="text"
    option-value="value"
    :show-clear="!field.isRequired"
    :disabled="disabled"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event ?? '')"
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

defineEmits(['update:modelValue'])

const selectOptions = computed(() => {
  return props.field.options?.options || []
})
</script>
