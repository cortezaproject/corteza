<template>
  <CFieldEditor
    :field="field"
    :namespace="namespace"
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { provide } from 'vue'
import { components } from '@planetcrust/human-vue'

const { CFieldEditor } = components

defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
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

const emit = defineEmits(['update:modelValue', 'stageFiles'])

// File editors stage their files and leave the upload to whoever saves the
// record, reaching that saver through an injected $fileUploadContext keyed only
// by field name. A record editor holds one record so a single context suffices,
// but a record list holds many rows sharing those field names — one context per
// list would let one row's staged file land on another. Each editing cell is its
// own component instance, so it can provide a context of its own and report
// upward; the row it belongs to is whatever the parent bound this handler with.
provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    emit('stageFiles', { fieldName, files })
  },
})
</script>
