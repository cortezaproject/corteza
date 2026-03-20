<template>
  <!-- Gallery mode -->
  <div v-if="isGallery" class="grid grid-cols-2 md:grid-cols-3 gap-3">
    <div
      v-for="att in resolvedAttachments"
      :key="att.attachmentID"
      class="flex flex-col items-center gap-1 p-3 border border-surface rounded"
    >
      <i v-if="!att.isImage" class="pi pi-file text-3xl text-primary" />
      <a
        v-else-if="opts.clickToView"
        :href="att.originalUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img
          :src="att.previewUrl"
          :alt="att.name"
          :style="galleryImgStyle"
          class="w-full max-h-32 object-contain rounded"
        />
      </a>
      <img
        v-else
        :src="att.previewUrl"
        :alt="att.name"
        :style="galleryImgStyle"
        class="w-full max-h-32 object-contain rounded"
      />
      <span v-if="!opts.hideFileName" class="text-sm truncate max-w-full text-center">
        {{ att.name }}
      </span>
    </div>
  </div>

  <!-- List mode -->
  <div v-else class="flex flex-col gap-1">
    <div
      v-for="att in resolvedAttachments"
      :key="att.attachmentID"
      class="flex items-center gap-2 px-3 py-2 border border-surface rounded"
    >
      <i class="pi pi-file text-primary" />
      <div class="flex-1 min-w-0">
        <div class="text-sm font-medium truncate">
          <a
            v-if="opts.clickToView"
            :href="att.originalUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:underline"
          >
            {{ att.name }}
          </a>
          <span v-else>{{ att.name }}</span>
        </div>
        <div v-if="att.size" class="text-xs text-muted-color">{{ formatSize(att.size) }}</div>
      </div>
      <a
        v-if="opts.enableDownload && att.downloadUrl"
        :href="att.downloadUrl"
        download
        class="text-muted-color hover:text-primary transition-colors"
      >
        <i class="pi pi-download text-sm" />
      </a>
    </div>
    <span v-if="!resolvedAttachments.length" class="text-muted-color text-sm">—</span>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const $ComposeAPI = inject('$ComposeAPI', null)

const opts = computed(() => props.field.options || {})
const isGallery = computed(() => opts.value.mode === 'gallery')
const namespaceID = computed(() => props.namespace?.namespaceID || '')

const attachmentIDs = computed(() => {
  const v = props.record?.values?.[props.field.name]
  if (!v) return []
  if (Array.isArray(v)) return v.filter(Boolean)
  return [v].filter(Boolean)
})

const resolvedAttachments = ref([])

const galleryImgStyle = computed(() => {
  const o = opts.value
  return {
    height: o.height || undefined,
    width: o.width || undefined,
    maxHeight: o.maxHeight || undefined,
    maxWidth: o.maxWidth || undefined,
    borderRadius: o.borderRadius || undefined,
    margin: o.margin || undefined,
    backgroundColor: o.backgroundColor ? `#${o.backgroundColor}` : undefined,
  }
})

function formatSize(bytes) {
  if (!bytes) return ''
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

async function resolveAttachments(ids) {
  if (!ids.length || !$ComposeAPI || !namespaceID.value) {
    resolvedAttachments.value = []
    return
  }

  const baseURL = $ComposeAPI.baseURL || ''
  const results = []

  for (const attachmentID of ids) {
    try {
      const att = await $ComposeAPI.attachmentRead({
        kind: 'record',
        namespaceID: namespaceID.value,
        attachmentID,
      })

      const mime = att.meta?.original?.mimetype || ''
      const url = att.url ? baseURL + att.url : ''

      results.push({
        attachmentID: att.attachmentID,
        name: att.name || attachmentID,
        size: att.meta?.original?.size || 0,
        isImage: mime.startsWith('image/'),
        previewUrl: url,
        originalUrl: url,
        downloadUrl: url,
      })
    } catch {
      results.push({
        attachmentID,
        name: attachmentID,
        size: 0,
        isImage: false,
        previewUrl: '',
        originalUrl: '',
        downloadUrl: '',
      })
    }
  }

  resolvedAttachments.value = results
}

watch(attachmentIDs, ids => resolveAttachments(ids), { immediate: true })
</script>
