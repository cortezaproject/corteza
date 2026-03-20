<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full p-3">
      <ProgressSpinner style="width: 2rem; height: 2rem" />
    </div>

    <div v-else-if="resolvedAttachments.length" class="p-3">
      <!-- List mode -->
      <div v-if="viewMode === 'list'" class="flex flex-col gap-2">
        <div
          v-for="att in resolvedAttachments"
          :key="att.attachmentID"
          class="flex items-center gap-3 p-2 border border-surface rounded"
        >
          <i class="pi pi-file text-xl text-primary" />
          <div class="flex-1 min-w-0">
            <div v-if="showName" class="font-medium truncate">{{ att.name }}</div>
            <div class="text-sm text-muted-color">{{ formatSize(att.meta?.original?.size) }}</div>
          </div>
          <Button
            v-if="enableDownload && att.downloadUrl"
            icon="pi pi-download"
            text
            size="small"
            severity="secondary"
            @click="downloadAttachment(att)"
          />
        </div>
      </div>

      <!-- Grid / Gallery mode -->
      <div v-else class="grid grid-cols-2 md:grid-cols-3 gap-3">
        <div
          v-for="att in resolvedAttachments"
          :key="att.attachmentID"
          class="flex flex-col items-center gap-1 p-3 border border-surface rounded"
        >
          <i v-if="!isImage(att)" class="pi pi-file text-3xl text-primary" />
          <img
            v-else
            :src="att.previewUrl || '#'"
            :alt="att.name"
            class="w-full max-h-32 object-contain rounded"
          />
          <span v-if="showName" class="text-sm truncate max-w-full text-center">
            {{ att.name }}
          </span>
        </div>
      </div>
    </div>
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.file.label') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const $ComposeAPI = inject('$ComposeAPI')

const viewMode = computed(() => props.block.options?.mode || 'list')
const showName = computed(() => props.block.options?.showName ?? true)
const enableDownload = computed(() => props.block.options?.enableDownload ?? true)

const loading = ref(false)
const resolvedAttachments = ref([])

const rawAttachments = computed(() => props.block.options?.attachments || [])

// Resolve attachment IDs into full attachment objects
async function resolveAttachments(ids) {
  if (!ids.length || !$ComposeAPI || !props.namespace?.namespaceID) {
    resolvedAttachments.value = []
    return
  }

  loading.value = true
  const results = []
  const baseURL = $ComposeAPI.baseURL || ''

  for (const entry of ids) {
    // Already a full object
    if (typeof entry === 'object' && entry.attachmentID) {
      results.push(entry)
      continue
    }

    // It's an ID string — resolve it
    const attachmentID = typeof entry === 'string' ? entry : String(entry)
    try {
      const att = await $ComposeAPI.attachmentRead({
        kind: 'page',
        namespaceID: props.namespace.namespaceID,
        attachmentID,
      })

      results.push({
        attachmentID: att.attachmentID,
        name: att.name,
        meta: att.meta,
        previewUrl: att.url ? baseURL + att.url : '',
        downloadUrl: att.url ? baseURL + att.url : '',
      })
    } catch (e) {
      // Skip unresolvable attachments
    }
  }

  resolvedAttachments.value = results
  loading.value = false
}

watch(rawAttachments, ids => resolveAttachments(ids), { immediate: true, deep: true })

function isImage(att) {
  const mime = att.meta?.original?.mimetype || ''
  return mime.startsWith('image/')
}

function formatSize(bytes) {
  if (!bytes) return ''
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function downloadAttachment(att) {
  if (att.downloadUrl) {
    window.open(att.downloadUrl, '_blank')
  }
}
</script>
