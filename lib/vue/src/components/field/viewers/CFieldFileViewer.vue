<template>
  <!-- Gallery mode -->
  <div v-if="isGallery" class="flex flex-wrap">
    <div
      v-for="attID in attachmentIDs"
      :key="attID"
    >
      <a
        v-if="opts.clickToView"
        :href="originalUrl(attID)"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img
          :src="previewUrl(attID)"
          :alt="metaName(attID)"
          :style="galleryImgStyle"
          class="object-cover"
        />
      </a>
      <img
        v-else
        :src="previewUrl(attID)"
        :alt="metaName(attID)"
        :style="galleryImgStyle"
        class="object-cover"
      />
      <div v-if="!opts.hideFileName" class="text-xs text-center text-muted-color mt-1 truncate max-w-[120px]">
        {{ metaName(attID) }}
      </div>
      <a
        v-if="opts.enableDownload"
        :href="downloadUrl(attID)"
        download
        class="block text-xs text-center text-blue-500 hover:underline"
      >
        Download
      </a>
    </div>
  </div>

  <!-- List mode -->
  <div v-else class="flex flex-col gap-1">
    <div
      v-for="attID in attachmentIDs"
      :key="attID"
      class="flex items-center gap-2 text-sm"
    >
      <i class="pi pi-paperclip text-xs text-muted-color" />
      <a
        v-if="opts.clickToView"
        :href="originalUrl(attID)"
        target="_blank"
        rel="noopener noreferrer"
        class="text-blue-500 hover:underline truncate"
      >
        {{ metaName(attID) }}
      </a>
      <span v-else class="truncate">{{ metaName(attID) }}</span>
      <a
        v-if="opts.enableDownload"
        :href="downloadUrl(attID)"
        download
        class="ml-auto text-xs text-blue-500 hover:underline shrink-0"
      >
        <i class="pi pi-download" />
      </a>
    </div>
    <span v-if="!attachmentIDs.length" class="text-muted-color text-sm">—</span>
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

const metaCache = ref(new Map())

const opts = computed(() => props.field.options || {})
const isGallery = computed(() => opts.value.mode === 'gallery')

const namespaceID = computed(() => props.namespace?.namespaceID || '')

const attachmentIDs = computed(() => {
  const v = props.record?.values?.[props.field.name]
  if (!v) return []
  if (Array.isArray(v)) return v.filter(Boolean)
  return [v].filter(Boolean)
})

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

function metaName(attID) {
  return metaCache.value.get(attID)?.name || attID
}

function previewUrl(attID) {
  if (!$ComposeAPI) return ''
  try {
    return $ComposeAPI.attachmentPreviewEndpoint({
      kind: 'record',
      namespaceID: namespaceID.value,
      attachmentID: attID,
      ext: 'png',
    })
  } catch {
    return `/namespace/${namespaceID.value}/attachment/record/${attID}/preview.png`
  }
}

function originalUrl(attID) {
  if (!$ComposeAPI) return ''
  const name = metaName(attID)
  try {
    return $ComposeAPI.attachmentOriginalEndpoint({
      kind: 'record',
      namespaceID: namespaceID.value,
      attachmentID: attID,
      name,
      download: false,
    })
  } catch {
    return `/namespace/${namespaceID.value}/attachment/record/${attID}/original/${name}`
  }
}

function downloadUrl(attID) {
  if (!$ComposeAPI) return ''
  const name = metaName(attID)
  try {
    return $ComposeAPI.attachmentOriginalEndpoint({
      kind: 'record',
      namespaceID: namespaceID.value,
      attachmentID: attID,
      name,
      download: true,
    })
  } catch {
    return `/namespace/${namespaceID.value}/attachment/record/${attID}/original/${name}`
  }
}

async function fetchMeta(attID) {
  if (!$ComposeAPI || metaCache.value.has(attID)) return
  try {
    const meta = await $ComposeAPI.attachmentRead({
      kind: 'record',
      namespaceID: namespaceID.value,
      attachmentID: attID,
    })
    metaCache.value.set(attID, meta)
  } catch {
    // keep ID as display fallback
  }
}

watch(
  attachmentIDs,
  ids => {
    for (const id of ids) fetchMeta(id)
  },
  { immediate: true },
)
</script>
