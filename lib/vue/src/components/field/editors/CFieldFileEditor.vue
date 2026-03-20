<template>
  <CInputFile
    :model-value="normalizedValue"
    :disabled="disabled"
    :multiple="field.isMulti"
    :accept="field.options?.mimetypes || ''"
    :max-file-size="field.options?.maxSize || 0"
    :attachment-info="attachmentInfo"
    @update:model-value="$emit('update:modelValue', $event)"
    @stage-files="onStageFiles"
  />
</template>

<script setup>
import { computed, inject, onUnmounted, reactive, watch } from 'vue'
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

const $ComposeAPI = inject('$ComposeAPI', null)
const $fileUploadContext = inject('$fileUploadContext', null)
const attachmentInfo = reactive({})

// Resolve attachment details for display
async function resolveAttachments(ids) {
  if (!$ComposeAPI || !props.namespace?.namespaceID) return

  for (const id of ids) {
    if (attachmentInfo[id]) continue
    try {
      const att = await $ComposeAPI.attachmentRead({
        kind: 'record',
        namespaceID: props.namespace.namespaceID,
        attachmentID: id,
      })
      const baseURL = $ComposeAPI.baseURL || ''
      attachmentInfo[id] = {
        name: att.name,
        size: att.meta?.original?.size || 0,
        downloadUrl: att.url ? baseURL + att.url : '',
      }
    } catch (e) {
      // If we can't resolve, leave it as ID
    }
  }
}

watch(normalizedValue, (ids) => resolveAttachments(ids), { immediate: true })

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

