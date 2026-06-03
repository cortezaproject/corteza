<template>
  <div class="flex flex-col gap-3">
    <CFileDropZone
      accept="*"
      :multiple="true"
      :uploading="uploading"
      :error="uploadError"
      @select="onFilesSelected"
    />

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

    <Fieldset :legend="$t('block.file.view.modeLabel')">
      <div class="flex flex-col gap-3">
        <CFormGroup
          :label="$t('block.file.view.modeLabel')"
          :description="$t('block.file.view.modeFootnote')"
        >
          <SelectButton v-model="mode" :options="modes" option-label="text" option-value="value" />
        </CFormGroup>

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
      </div>
    </Fieldset>

    <template v-if="mode === 'gallery'">
      <Divider />

      <Fieldset :legend="$t('block.file.view.previewStyle')">
        <div class="flex flex-col gap-3">
          <small class="text-muted-color">{{ $t('block.file.view.description') }}</small>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <CFormGroup :label="$t('block.file.view.height')">
              <InputText v-model="height" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.width')">
              <InputText v-model="width" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.maxHeight')">
              <InputText v-model="maxHeight" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.maxWidth')">
              <InputText v-model="maxWidth" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.borderRadius')">
              <InputText v-model="borderRadius" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.margin')">
              <InputText v-model="margin" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('block.file.view.backgroundColor')">
              <CInputColorPicker
                :model-value="backgroundColor"
                show-text
                @update:model-value="backgroundColor = $event"
              />
            </CFormGroup>
          </div>
        </div>
      </Fieldset>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, useFileUpload } from '@planetcrust/human-vue'

const { CInputColorPicker, CFileDropZone } = components

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

// Upload state (from composable)
const { uploading, uploadError, uploadFileRaw } = useFileUpload()
const attachmentMeta = reactive({})

const modes = [
  { value: 'list', text: t('block.file.view.list') },
  { value: 'gallery', text: t('block.file.view.gallery') },
]

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

function formatSize(bytes) {
  if (!bytes) return ''
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// Attachment IDs
const attachmentIDs = computed(() => block.value.options?.attachments || [])

// Resolve attachment details for display
async function resolveAttachments(ids) {
  const namespaceID = props.namespace?.namespaceID
  if (!namespaceID) return

  for (const id of ids) {
    if (attachmentMeta[id]) continue
    try {
      const att = await $ComposeAPI.attachmentRead({ kind: 'page', namespaceID, attachmentID: id })
      attachmentMeta[id] = { name: att.name, size: att.meta?.original?.size || 0 }
    } catch {
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

// Upload handler using composable
async function onFilesSelected(files) {
  const pageID = props.page?.pageID
  const namespaceID = props.namespace?.namespaceID

  if (!pageID || !namespaceID) {
    uploadError.value = t('notification.namespace.importFailed')
    return
  }

  const endpoint = $ComposeAPI.baseURL + $ComposeAPI.pageUploadEndpoint({ namespaceID, pageID })
  const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

  try {
    for (const file of files) {
      const att = await uploadFileRaw(file, { url: endpoint, token })
      if (att.attachmentID) {
        attachmentMeta[att.attachmentID] = {
          name: att.name || file.name,
          size: att.meta?.original?.size || file.size,
        }
        const updated = [...attachmentIDs.value, att.attachmentID]
        updateOptions('attachments', updated)
      }
    }
  } catch {
    // uploadError is already set by the composable
  }
}

const mode = computed({
  get: () => block.value.options?.mode || 'list',
  set: v => updateOptions('mode', v),
})

const hideFileName = computed({
  get: () => !block.value.options?.hideFileName,
  set: v => updateOptions('hideFileName', !v),
})

const clickToView = computed({
  get: () => !!block.value.options?.clickToView,
  set: v => updateOptions('clickToView', v),
})

const enableDownload = computed({
  get: () => block.value.options?.enableDownload !== false,
  set: v => updateOptions('enableDownload', v),
})

const height = computed({
  get: () => block.value.options?.height || '',
  set: v => updateOptions('height', v),
})

const width = computed({
  get: () => block.value.options?.width || '',
  set: v => updateOptions('width', v),
})

const maxHeight = computed({
  get: () => block.value.options?.maxHeight || '',
  set: v => updateOptions('maxHeight', v),
})

const maxWidth = computed({
  get: () => block.value.options?.maxWidth || '',
  set: v => updateOptions('maxWidth', v),
})

const borderRadius = computed({
  get: () => block.value.options?.borderRadius || '',
  set: v => updateOptions('borderRadius', v),
})

const margin = computed({
  get: () => block.value.options?.margin || '',
  set: v => updateOptions('margin', v),
})

const backgroundColor = computed({
  get: () => block.value.options?.backgroundColor || '',
  set: v => updateOptions('backgroundColor', v),
})
</script>
