<template>
  <ToggleSwitch
    v-if="field.options?.switch"
    :model-value="boolValue"
    :disabled="disabled"
    @update:model-value="onUpdate"
  />
  <Checkbox
    v-else
    :model-value="boolValue"
    :binary="true"
    :disabled="disabled"
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

const boolValue = computed(() => props.modelValue === '1')

function onUpdate(value) {
  emit('update:modelValue', value ? '1' : '')
}
</script>
