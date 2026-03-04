<template>
  <CInputRecord
    :model-value="modelValue"
    :namespace-id="namespaceID"
    :module-id="field.options.moduleID"
    :label-field="field.options.labelField"
    placeholder=""
    :disabled="disabled || !namespaceID || !field.options.moduleID"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed, inject } from 'vue'
import CInputRecord from '../../input/CInputRecord.vue'

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
  namespace: {
    type: Object,
    default: () => ({}),
  },
})

defineEmits(['update:modelValue'])

// $namespace is provided as a ref by Namespace/View.vue
const $namespace = inject('$namespace', null)

const namespaceID = computed(
  () => props.namespace?.namespaceID || $namespace?.value?.namespaceID || '',
)
</script>
