<template>
  <CInputFile
    :model-value="normalizedValue"
    :disabled="disabled"
    :multiple="field.isMulti"
    :accept="field.options?.mimetypes || ''"
    :max-file-size="field.options?.maxSize || 0"
    @update:model-value="$emit('update:modelValue', $event)"
    @stage-files="onStageFiles"
  />
</template>

<script setup>
import { computed, inject, onUnmounted } from 'vue'
import CInputFile from '../../input/CInputFile.vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
    default: () => [],
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

const normalizedValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue.filter(Boolean)
  if (props.modelValue) return [props.modelValue]
  return []
})

const $fileUploadContext = inject('$fileUploadContext', null)

function onStageFiles(files) {
  if ($fileUploadContext) {
    if (files.length > 0) {
      $fileUploadContext.registerPending(props.field.name, files)
    } else {
      $fileUploadContext.registerPending(props.field.name, [])
    }
  }
}

onUnmounted(() => {
  if ($fileUploadContext) {
    $fileUploadContext.registerPending(props.field.name, [])
  }
})
</script>
