<template>
  <div>
    <Button
      :label="$t('general.label.import')"
      icon="pi pi-upload"
      size="small"
      severity="secondary"
      @click="showDialog = true"
    />

    <Dialog
      v-model:visible="showDialog"
      :header="$t('module.import')"
      modal
      :closable="!importing"
      :style="{ width: '32rem' }"
      @hide="onDialogHide"
    >
      <!-- Importing spinner -->
      <div v-if="importing" class="flex flex-col items-center justify-center py-12 gap-3">
        <ProgressSpinner />
        <span class="text-muted-color text-sm">
          {{ $t('general.label.processing') }}
        </span>
      </div>

      <!-- Step 0: File Upload -->
      <div v-else-if="step === 0" class="flex flex-col gap-4">
        <CFileDropZone
          accept=".json"
          :uploading="false"
          :error="uploadError"
          :drop-label="dropLabel"
          @select="onFilesSelected"
        />
      </div>

      <!-- Step 1: Select & confirm -->
      <div v-else-if="step === 1" class="flex flex-col gap-4">
        <div class="flex gap-2 mb-2">
          <Button
            :label="$t('field.selectAll')"
            severity="secondary"
            size="small"
            @click="selectAll(true)"
          />
          <Button
            :label="$t('field.unselectAll')"
            severity="secondary"
            size="small"
            @click="selectAll(false)"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
          <div
            v-for="(item, index) in parsedItems"
            :key="index"
            class="flex items-center gap-2"
          >
            <Checkbox
              v-model="item.import"
              :binary="true"
              :input-id="`import-item-${index}`"
            />
            <label :for="`import-item-${index}`" class="cursor-pointer">
              {{ item.name || item.title }}
            </label>
          </div>
        </div>
      </div>

      <template v-if="!importing && step === 1" #footer>
        <div class="flex items-center justify-between w-full">
          <Button
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            text
            size="small"
            @click="onBack"
          />
          <Button
            :label="$t('general.label.import')"
            icon="pi pi-download"
            size="small"
            :disabled="!hasSelectedItems"
            @click="onImport"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { inject } from 'vue'
import { components } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'

const { CFileDropZone } = components

const { t } = useI18n()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['imported', 'failed'])

// State
const showDialog = ref(false)
const step = ref(0)
const importing = ref(false)
const uploadError = ref('')
const parsedItems = ref([])

const dropLabel = computed(() => t('module.import') + ' (.json)')

const hasSelectedItems = computed(() => {
  return parsedItems.value.some(i => i.import)
})

// Methods
function onDialogHide() {
  if (!importing.value) {
    reset()
  }
}

function reset() {
  step.value = 0
  parsedItems.value = []
  importing.value = false
  uploadError.value = ''
}

function onBack() {
  step.value = 0
  uploadError.value = ''
}

function selectAll(select) {
  parsedItems.value = parsedItems.value.map(i => {
    i.import = select
    return i
  })
}

async function onFilesSelected(files) {
  const file = files[0]
  if (!file) return

  uploadError.value = ''

  try {
    const text = await file.text()
    const data = JSON.parse(text)

    if (!data.list) {
      uploadError.value = t('notification.general.import.readingError')
      return
    }

    parsedItems.value = data.list.map(i => {
      return { import: true, ...i }
    })
    step.value = 1
  } catch {
    uploadError.value = t('notification.general.import.readingError')
  }
}

async function onImport() {
  importing.value = true

  try {
    for (const item of parsedItems.value.filter(i => i.import)) {
      await moduleStore.create({
        ...item,
        handle: item.handle ? `${item.handle}_imported` : '',
        namespaceID: props.namespace.namespaceID,
      })
    }

    $toast.toastSuccess(t('notification.general.import.successful'))
    emit('imported')
    showDialog.value = false
    reset()
  } catch (e) {
    console.error('Module import failed:', e)
    $toast.toastDanger(e.message || t('notification.general.import.failed'))
    emit('failed', e)
    showDialog.value = false
    reset()
  }
}
</script>
