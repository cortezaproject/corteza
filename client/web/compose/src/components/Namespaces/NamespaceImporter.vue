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
        <div
          class="flex flex-col items-center justify-center gap-3 p-8 border-2 border-dashed rounded-border cursor-pointer transition-colors duration-200"
          :class="dragOver ? 'border-primary bg-primary/5' : 'border-surface-300 hover:border-primary'"
          @click="$refs.fileInput.click()"
          @dragover.prevent="dragOver = true"
          @dragleave.prevent="dragOver = false"
          @drop.prevent="onFileDrop"
        >
          <i class="pi pi-cloud-upload text-4xl text-muted-color" />
          <span class="text-muted-color text-sm text-center">
            {{ $t('namespace.import.uploadFilePlaceholder') }}
          </span>
        </div>

        <input
          ref="fileInput"
          type="file"
          accept=".zip"
          class="hidden"
          @change="onFileSelected"
        />

        <Message v-if="uploadError" severity="error" :closable="false">
          {{ uploadError }}
        </Message>
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const emit = defineEmits(['imported', 'failed'])

// State
const showDialog = ref(false)
const step = ref(0)
const importing = ref(false)
const uploadError = ref('')
const session = ref({})
const name = ref('')
const slug = ref('')
const dragOver = ref(false)
const fileInput = ref(null)

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
  uploadError.value = ''
  importing.value = false
  dragOver.value = false
  if (fileInput.value) fileInput.value.value = ''
}

function onBack() {
  step.value = 0
  uploadError.value = ''
}

function onFileSelected(event) {
  const file = event.target?.files?.[0]
  if (file) handleUpload(file)
}

function onFileDrop(event) {
  dragOver.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) handleUpload(file)
}

async function handleUpload(file) {
  uploadError.value = ''
  if (!file) return

  try {
    const formData = new FormData()
    formData.append('upload', file)

    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.namespaceImportInitEndpoint()
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

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
      throw new Error(errData?.error?.message || t('notification.namespace.importFailed'))
    }

    const data = await response.json()
    session.value = data.response || data
    step.value = 1
  } catch (err) {
    uploadError.value = err.message || t('notification.namespace.importFailed')
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
