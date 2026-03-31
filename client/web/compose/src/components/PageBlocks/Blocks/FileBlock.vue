<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full p-3">
      <ProgressSpinner style="width: 2rem; height: 2rem" />
    </div>

    <div v-else-if="resolvedAttachments.length" class="p-2">
      <!-- List mode -->
      <div v-if="viewMode === 'list'" class="flex flex-col gap-2">
        <div
          v-for="att in resolvedAttachments"
          :key="att.attachmentID"
          class="flex items-center gap-3 p-2 border border-surface rounded"
        >
          <i class="pi pi-file text-xl text-primary" />
          <div class="flex-1 min-w-0">
            <div v-if="!hideFileName" class="font-medium truncate">{{ att.name }}</div>
            <div class="text-sm text-muted-color">{{ formatSize(att.size) }}</div>
          </div>
          <a
            v-if="att.download"
            :href="att.download"
            @click.stop
          >
            <Button
              icon="pi pi-download"
              text
              size="small"
              severity="secondary"
            />
          </a>
        </div>
      </div>

      <!-- Gallery mode -->
      <div
        v-else
        class="flex items-start justify-around gap-3 flex-wrap h-full"
      >
        <div
          v-for="att in resolvedAttachments"
          :key="att.attachmentID"
          class="item-preview relative"
        >
          <!-- Image preview -->
          <template v-if="isImage(att)">
            <a
              v-if="att.clickToView && att.url"
              :href="att.url"
              target="_blank"
              rel="noopener noreferrer"
            >
              <img
                :src="att.previewUrl || att.url"
                :alt="att.name"
                :style="{ width: 'unset', ...inlineCustomStyles(att) }"
                class="object-contain"
              />
            </a>
            <img
              v-else
              :src="att.previewUrl || att.url"
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
              v-if="!hideFileName"
              class="filename-container text-center"
              style="margin-top: 0.1rem;"
            >
              <a
                v-if="att.clickToView && att.url"
                :href="att.url"
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
            v-if="att.download"
            :href="att.download"
            class="preview-download-button absolute top-0 right-0"
            @click.stop
          >
            <Button
              icon="pi pi-download"
              text
              size="small"
              severity="secondary"
            />
          </a>
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
const hideFileName = computed(() => !!props.block.options?.hideFileName)
const enableDownload = computed(() => props.block.options?.enableDownload !== false)
const clickToView = computed(() => props.block.options?.clickToView !== false)

const loading = ref(false)
const resolvedAttachments = ref([])

const rawAttachments = computed(() => props.block.options?.attachments || [])

function inlineCustomStyles (att) {
  const o = props.block.options || {}
  let { width, height, maxWidth, maxHeight, margin, borderRadius, backgroundColor } = o

  maxWidth = maxWidth || '100%'
  maxHeight = maxHeight || '100%'
  margin = margin || 'auto'

  if (!isImage(att)) {
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

      const url = att.url ? baseURL + att.url : ''

      results.push({
        attachmentID: att.attachmentID,
        name: att.name,
        meta: att.meta,
        size: att.meta?.original?.size || 0,
        url,
        previewUrl: att.previewUrl ? baseURL + att.previewUrl : url,
        download: enableDownload.value && url ? url + '&download=1' : undefined,
        clickToView: clickToView.value,
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


