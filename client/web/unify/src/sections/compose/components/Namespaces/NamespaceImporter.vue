<template>
  <div>
    <Button
      :label="$t('namespace.import.buttonLabel')"
      icon="pi pi-upload"
      size="small"
      severity="secondary"
      @click="showDialog = true"
    />

    <Dialog
      v-model:visible="showDialog"
      :header="$t('namespace.import.title')"
      modal
      :closable="!importing"
      :style="{ width: '32rem' }"
      @hide="onDialogHide"
    >
      <!-- Importing spinner -->
      <div v-if="importing" class="flex items-center justify-center py-12">
        <ProgressSpinner />
      </div>

      <!-- Step 0: File Upload -->
      <div v-else-if="step === 0" class="flex flex-col gap-4">
        <CFileDropZone
          accept=".zip"
          :uploading="uploading"
          :error="uploadError"
          :drop-label="$t('namespace.import.uploadFilePlaceholder')"
          @select="onFilesSelected"
        />
      </div>

      <!-- Step 1: Configure name & slug -->
      <div v-else-if="step === 1" class="flex flex-col gap-5">
        <CFormGroup :label="$t('namespace.name.label')" input-id="import-name">
          <InputText
            id="import-name"
            v-model="name"
            :placeholder="$t('namespace.name.placeholder')"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('namespace.import.slug.label')"
          :description="$t('namespace.slug.description')"
          input-id="import-slug"
        >
          <InputText
            id="import-slug"
            v-model="slug"
            :placeholder="$t('namespace.slug.placeholder')"
            :invalid="slug.length > 0 && !slugValid"
          />
          <Message
            v-if="slug.length > 0 && !slugValid"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ $t('namespace.slug.invalid-handle-characters') }}
          </Message>
        </CFormGroup>
      </div>

      <template v-if="!importing && step === 1" #footer>
        <div class="flex items-center justify-between w-full">
          <Button
            :label="$t('namespace.import.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            text
            size="small"
            @click="onBack"
          />
          <Button
            :label="$t('namespace.import.import')"
            icon="pi pi-download"
            size="small"
            :disabled="!canImport"
            @click="onImport"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, useFileUpload } from '@planetcrust/human-vue'

const { CFileDropZone } = components

const $ComposeAPI = inject('$ComposeAPI')
const { t } = useI18n()

const emit = defineEmits(['imported', 'failed'])

// State
const showDialog = ref(false)
const step = ref(0)
const importing = ref(false)
const session = ref({})
const name = ref('')
const slug = ref('')

// Upload state (from composable)
const { uploading, uploadError, uploadFileRaw, reset: resetUpload } = useFileUpload()

// Validation
const slugRegex = /^[a-zA-Z][a-zA-Z0-9_]*$/
const slugValid = computed(() => slug.value.length === 0 || slugRegex.test(slug.value))

const canImport = computed(() => {
  return name.value.trim().length > 0 && slug.value.trim().length > 0 && slugValid.value
})

// Methods
function onDialogHide() {
  if (!importing.value) {
    reset()
  }
}

function reset() {
  step.value = 0
  session.value = {}
  name.value = ''
  slug.value = ''
  importing.value = false
  resetUpload()
}

function onBack() {
  step.value = 0
  uploadError.value = ''
}

// Server-side accepts only zipped Envoy YAML bundles (mimetype application/zip);
// the dropzone's `accept=".zip"` is just a hint that drag-drop bypasses, so
// re-check here so the user gets a clear message instead of a generic 400 →
// "field sessionID is empty" downstream.
const ZIP_MIME_TYPES = ['application/zip', 'application/x-zip-compressed', 'application/x-zip']

function isZipFile(file) {
  if (file.name?.toLowerCase().endsWith('.zip')) return true
  return ZIP_MIME_TYPES.includes(file.type)
}

async function onFilesSelected(files) {
  const file = files[0]
  if (!file) return

  if (!isZipFile(file)) {
    uploadError.value = t('namespace.import.invalidFileFormat')
    return
  }

  try {
    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.namespaceImportInitEndpoint()
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

    const data = await uploadFileRaw(file, { url: endpoint, token })
    if (!data?.sessionID) {
      uploadError.value = t('namespace.import.invalidResponse')
      return
    }
    session.value = data
    step.value = 1
  } catch {
    // uploadError is already set by the composable
  }
}

async function onImport() {
  importing.value = true

  try {
    const out = await $ComposeAPI.namespaceImportRun({
      sessionID: session.value.sessionID,
      name: name.value,
      slug: slug.value,
    })

    emit('imported', out)
    showDialog.value = false
    reset()
  } catch (err) {
    emit('failed', err)
    showDialog.value = false
    reset()
  }
}
</script>
