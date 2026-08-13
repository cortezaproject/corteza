<template>
  <CInputFile
    :model-value="normalizedValue"
    :disabled="disabled"
    :multiple="field.isMulti"
    :accept="accept"
    :max-file-size="maxFileSize"
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

const $Settings = inject('$Settings', null)

// Both constraints mirror the server's resolution order in
// compose/service/attachment.go: the field option overrides the system-wide
// record-attachment setting, and an unset constraint means "no limit".
// Staying in step matters — uploads are deferred to record save, so a check the
// client skips surfaces as a failed save rather than a rejected file.
const maxFileSize = computed(() => {
  // Configured and stored in megabytes; the server enforces maxSize * 1_000_000.
  // CInputFile compares raw byte counts, so convert with the same factor.
  const mb = props.field.options?.maxSize || globalSetting('MaxSize') || 0
  return mb * 1_000_000
})

const accept = computed(() => {
  const own = props.field.options?.mimetypes
  if (own) return own
  // The system-wide allow list is a list; CInputFile takes it comma-separated.
  return (globalSetting('Mimetypes') || []).join(',')
})

// The current-settings endpoint returns the server's struct, so these read back
// capital-cased — unlike the kv keys the settings editor writes.
function globalSetting(name) {
  return $Settings?.get(`compose.Record.Attachments.${name}`)
}

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
    } catch {
      // If we can't resolve, leave it as ID
    }
  }
}

watch(normalizedValue, ids => resolveAttachments(ids), { immediate: true })

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
