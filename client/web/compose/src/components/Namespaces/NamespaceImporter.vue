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
        <div class="flex flex-col gap-2">
          <label for="import-name" class="font-medium text-primary">
            {{ $t('namespace.name.label') }}
          </label>
          <InputText
            id="import-name"
            v-model="name"
            :placeholder="$t('namespace.name.placeholder')"
          />
        </div>

        <div class="flex flex-col gap-2">
          <label for="import-slug" class="font-medium text-primary">
            {{ $t('namespace.import.slug.label') }}
          </label>
          <InputText
            id="import-slug"
            v-model="slug"
            :placeholder="$t('namespace.slug.placeholder')"
            :invalid="slug.length > 0 && !slugValid"
          />
          <small class="text-muted-color">
            {{ $t('namespace.slug.description') }}
          </small>
          <Message v-if="slug.length > 0 && !slugValid" severity="error" size="small" variant="simple">
            {{ $t('namespace.slug.invalid-handle-characters') }}
          </Message>
        </div>

        <div class="flex items-center justify-between pt-2">
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
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { components, useFileUpload } from '@planetcrust/human-vue'

const { CFileDropZone } = components

const $ComposeAPI = inject('$ComposeAPI')

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

async function onFilesSelected(files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.namespaceImportInitEndpoint()
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

    const data = await uploadFileRaw(file, { url: endpoint, token })
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
