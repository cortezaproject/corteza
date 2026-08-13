<template>
  <!-- Gallery mode -->
  <div v-if="isGallery" class="flex items-start justify-around gap-3 flex-wrap h-full">
    <div v-for="att in resolvedAttachments" :key="att.attachmentID" class="item-preview relative">
      <!-- Image preview -->
      <template v-if="att.isImage">
        <a
          v-if="isClickToView && att.originalUrl"
          :href="att.originalUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          <img
            :src="att.previewUrl || att.originalUrl"
            :alt="att.name"
            :style="{ width: 'unset', ...inlineCustomStyles(att) }"
            class="object-contain"
          />
        </a>
        <img
          v-else
          :src="att.previewUrl || att.originalUrl"
          :alt="att.name"
          :style="{ width: 'unset', ...inlineCustomStyles(att) }"
          class="object-contain"
        />
      </template>
      <i v-else class="pi pi-file text-3xl text-primary" />

      <!-- File name -->
      <div
        class="flex items-start justify-center"
        :style="{ width: `calc(${inlineCustomStyles(att).width || '100%'})` }"
      >
        <div
          v-if="!opts.hideFileName"
          class="filename-container text-center"
          style="margin-top: 0.1rem"
        >
          <a
            v-if="isClickToView && att.originalUrl"
            :href="att.originalUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:underline"
          >
            {{ att.name }}
          </a>
          <span v-else>{{ att.name }}</span>
        </div>
      </div>

      <!-- Download button (overlay) -->
      <a
        v-if="isDownloadEnabled && att.downloadUrl"
        :href="att.downloadUrl"
        class="preview-download-button absolute top-0 right-0"
        @click.stop
      >
        <Button icon="pi pi-download" text size="small" severity="secondary" />
      </a>
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
            v-if="isClickToView && att.originalUrl"
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
      <a v-if="isDownloadEnabled && att.downloadUrl" :href="att.downloadUrl" @click.stop>
        <Button icon="pi pi-download" text size="small" severity="secondary" />
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
const isDownloadEnabled = computed(() => opts.value.enableDownload !== false)
const isClickToView = computed(() => opts.value.clickToView !== false)

const attachmentIDs = computed(() => {
  const v = props.record?.values?.[props.field.name]
  if (!v) return []
  if (Array.isArray(v)) return v.filter(Boolean)
  return [v].filter(Boolean)
})

const resolvedAttachments = ref([])

function inlineCustomStyles(att) {
  const o = opts.value
  let { width, height, maxWidth, maxHeight, margin, borderRadius, backgroundColor } = o

  maxWidth = maxWidth || '100%'
  maxHeight = maxHeight || '100%'
  margin = margin || 'auto'

  if (!att.isImage) {
    width = width || '200px'
    height = height || 'auto'
  }

  return {
    width,
    height,
    maxWidth,
    maxHeight,
    borderRadius,
    backgroundColor: backgroundColor ? `#${backgroundColor}` : undefined,
    margin,
  }
}

function formatSize(bytes) {
  if (!bytes) return ''
  const k = 1000
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
        previewUrl: att.previewUrl ? baseURL + att.previewUrl : url,
        originalUrl: url,
        downloadUrl: url ? url + '&download=1' : '',
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

<style scoped>
.item-preview .preview-download-button {
  opacity: 0;
  transition: opacity 0.2s;
}

.item-preview:hover .preview-download-button {
  opacity: 1;
}

.filename-container {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-word;
  max-width: 100%;
}

.filename-container:hover {
  -webkit-line-clamp: unset;
  line-clamp: unset;
  overflow: visible;
}
</style>
