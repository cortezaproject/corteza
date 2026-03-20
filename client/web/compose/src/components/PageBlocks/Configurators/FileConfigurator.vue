<template>
  <div class="flex flex-col gap-3">
    <!-- File Upload -->
    <div class="flex flex-col gap-2">
      <div
        class="flex flex-col items-center justify-center gap-3 p-6 border-2 border-dashed rounded-border cursor-pointer transition-colors duration-200"
        :class="
          dragOver ? 'border-primary bg-primary/5' : 'border-surface-300 hover:border-primary'
        "
        @click="$refs.fileInput.click()"
        @dragover.prevent="dragOver = true"
        @dragleave.prevent="dragOver = false"
        @drop.prevent="onDrop"
      >
        <i v-if="!uploading" class="pi pi-cloud-upload text-4xl text-muted-color" />
        <ProgressSpinner v-else style="width: 2rem; height: 2rem" />
        <span class="text-muted-color text-sm text-center">
          {{ uploading ? $t('general.label.uploading') : $t('general.label.dropFiles') }}
        </span>
      </div>

      <input ref="fileInput" type="file" multiple class="hidden" @change="onFileSelected" />

      <Message v-if="uploadError" severity="error" size="small" :closable="false">
        {{ uploadError }}
      </Message>
    </div>

    <!-- Uploaded attachments list -->
    <div v-if="attachmentIDs.length" class="flex flex-col gap-1">
      <div
        v-for="(attID, index) in attachmentIDs"
        :key="attID"
        class="flex items-center gap-2 px-3 py-2 border border-surface rounded"
      >
        <i class="pi pi-file text-primary" />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium truncate">{{ attachmentMeta[attID]?.name || attID }}</div>
          <div v-if="attachmentMeta[attID]?.size" class="text-xs text-muted-color">
            {{ formatSize(attachmentMeta[attID].size) }}
          </div>
        </div>
        <Button
          icon="pi pi-times"
          text
          severity="danger"
          size="small"
          @click="removeAttachment(index)"
        />
      </div>
    </div>

    <Divider />

    <!-- View mode -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.file.view.modeLabel') }}</label>
      <SelectButton v-model="mode" :options="modes" option-label="text" option-value="value" />
      <small class="text-muted-color">{{ $t('block.file.view.modeFootnote') }}</small>
    </div>

    <!-- Display options -->
    <div class="flex flex-col gap-2">
      <div v-if="mode === 'gallery'" class="flex items-center gap-2">
        <Checkbox v-model="hideFileName" :binary="true" input-id="hideFileName" />
        <label for="hideFileName" class="text-sm">{{ $t('block.file.view.showName') }}</label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox v-model="clickToView" :binary="true" input-id="clickToView" />
        <label for="clickToView" class="text-sm">{{ $t('block.file.view.clickToView') }}</label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox v-model="enableDownload" :binary="true" input-id="enableDownload" />
        <label for="enableDownload" class="text-sm">
          {{ $t('block.file.view.enableDownload') }}
        </label>
      </div>
    </div>

    <!-- Gallery preview styling -->
    <template v-if="mode === 'gallery'">
      <Divider />

      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.file.view.previewStyle') }}
      </h5>
      <small class="text-muted-color">{{ $t('block.file.view.description') }}</small>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.height') }}</label>
          <InputText v-model="height" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.width') }}</label>
          <InputText v-model="width" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.file.view.maxHeight') }}
          </label>
          <InputText v-model="maxHeight" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.file.view.maxWidth') }}
          </label>
          <InputText v-model="maxWidth" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.file.view.borderRadius') }}
          </label>
          <InputText v-model="borderRadius" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.margin') }}</label>
          <InputText v-model="margin" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.file.view.backgroundColor') }}</label>
          <CInputColorPicker
            :model-value="backgroundColor"
            show-text
            @update:model-value="backgroundColor = $event"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputColorPicker } = components

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

// Upload state
const dragOver = ref(false)
const uploading = ref(false)
const uploadError = ref('')
const fileInput = ref(null)
const attachmentMeta = reactive({})

const modes = [
  { value: 'list', text: t('block.file.view.list') },
  { value: 'gallery', text: t('block.file.view.gallery') },
]

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

function formatSize(bytes) {
  if (!bytes) return ''
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// Attachment IDs
const attachmentIDs = computed(() => props.block.options?.attachments || [])

// Resolve attachment details for display
async function resolveAttachments(ids) {
  const namespaceID = props.namespace?.namespaceID
  if (!namespaceID) return

  for (const id of ids) {
    if (attachmentMeta[id]) continue
    try {
      const att = await $ComposeAPI.attachmentRead({ kind: 'page', namespaceID, attachmentID: id })
      attachmentMeta[id] = { name: att.name, size: att.meta?.original?.size || 0 }
    } catch (e) {
      // If we can't resolve, just leave it
    }
  }
}

watch(attachmentIDs, ids => resolveAttachments(ids), { immediate: true })

function removeAttachment(index) {
  const updated = [...attachmentIDs.value]
  updated.splice(index, 1)
  updateOptions('attachments', updated)
}

// Upload handlers
function onFileSelected(event) {
  const files = Array.from(event.target?.files || [])
  if (files.length) uploadFiles(files)
  if (fileInput.value) fileInput.value.value = ''
}

function onDrop(event) {
  dragOver.value = false
  const files = Array.from(event.dataTransfer?.files || [])
  if (files.length) uploadFiles(files)
}

async function uploadFiles(files) {
  uploading.value = true
  uploadError.value = ''

  const pageID = props.page?.pageID
  const namespaceID = props.namespace?.namespaceID

  if (!pageID || !namespaceID) {
    uploadError.value = t('notification.namespace.importFailed')
    uploading.value = false
    return
  }

  const endpoint = $ComposeAPI.baseURL + $ComposeAPI.pageUploadEndpoint({ namespaceID, pageID })
  const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

  try {
    for (const file of files) {
      const formData = new FormData()
      formData.append('upload', file)

      const response = await fetch(endpoint, {
        method: 'POST',
        body: formData,
        credentials: 'include',
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      })

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}))
        throw new Error(errData?.error?.message || `Upload failed for ${file.name}`)
      }

      const data = await response.json()
      const att = data.response || data
      if (att.attachmentID) {
        attachmentMeta[att.attachmentID] = {
          name: att.name || file.name,
          size: att.meta?.original?.size || file.size,
        }
        const updated = [...attachmentIDs.value, att.attachmentID]
        updateOptions('attachments', updated)
      }
    }
  } catch (err) {
    uploadError.value = err.message || t('notification.namespace.importFailed')
  } finally {
    uploading.value = false
  }
}

const mode = computed({
  get: () => props.block.options?.mode || 'list',
  set: v => updateOptions('mode', v),
})

const hideFileName = computed({
  get: () => !!props.block.options?.hideFileName,
  set: v => updateOptions('hideFileName', v),
})

const clickToView = computed({
  get: () => !!props.block.options?.clickToView,
  set: v => updateOptions('clickToView', v),
})

const enableDownload = computed({
  get: () => props.block.options?.enableDownload !== false,
  set: v => updateOptions('enableDownload', v),
})

const height = computed({
  get: () => props.block.options?.height || '',
  set: v => updateOptions('height', v),
})

const width = computed({
  get: () => props.block.options?.width || '',
  set: v => updateOptions('width', v),
})

const maxHeight = computed({
  get: () => props.block.options?.maxHeight || '',
  set: v => updateOptions('maxHeight', v),
})

const maxWidth = computed({
  get: () => props.block.options?.maxWidth || '',
  set: v => updateOptions('maxWidth', v),
})

const borderRadius = computed({
  get: () => props.block.options?.borderRadius || '',
  set: v => updateOptions('borderRadius', v),
})

const margin = computed({
  get: () => props.block.options?.margin || '',
  set: v => updateOptions('margin', v),
})

const backgroundColor = computed({
  get: () => props.block.options?.backgroundColor || '',
  set: v => updateOptions('backgroundColor', v),
})
</script>
