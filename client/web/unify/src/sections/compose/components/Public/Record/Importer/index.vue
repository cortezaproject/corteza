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
      :header="$t('block.recordList.import.to', { modulename: module.name })"
      modal
      :closable="step !== 2"
      :style="{ width: '60rem' }"
      :breakpoints="{ '1200px': '75vw', '768px': '95vw' }"
      :content-style="{ 'max-height': '75vh', overflow: 'auto' }"
      @hide="onDialogHide"
    >
      <!-- Step 0: File Upload -->
      <div v-if="step === 0" class="flex flex-col gap-5">
        <CFileDropZone
          accept=".csv,.json"
          :uploading="uploading"
          :error="uploadError"
          :drop-label="$t('block.recordList.import.uploadFile')"
          @select="onFilesSelected"
        />

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.import.onError') }}
            </label>
            <Select
              v-model="onError"
              :options="onErrorOptions"
              option-label="label"
              option-value="value"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('block.recordList.import.multiValueDelimiter.label') }}
            </label>
            <Select
              v-model="multiValueDelimiter"
              :options="delimiterOptions"
              option-label="label"
              option-value="value"
              class="w-full"
            />
          </div>
        </div>
      </div>

      <!-- Step 1: Field Mapping -->
      <div v-else-if="step === 1" class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.import.matchFields') }}
          </label>
          <p v-if="hasUnmappedRequired" class="text-xs text-red-500">
            {{ $t('block.recordList.import.hasRequiredFileFields') }}: {{ unmappedRequiredNames }}
          </p>
        </div>

        <DataTable
          :value="rows"
          size="small"
          scrollable
          scroll-height="50vh"
          class="text-sm"
        >
          <Column style="width: 3rem">
            <template #header>
              <Checkbox
                :model-value="allSelected"
                binary
                @update:model-value="toggleAll"
              />
            </template>
            <template #body="{ data }">
              <Checkbox v-model="data.selected" binary />
            </template>
          </Column>
          <Column :header="$t('block.recordList.import.fileColumns')" field="fileColumn" />
          <Column :header="$t('block.recordList.import.moduleFields')" style="width: 55%">
            <template #body="{ data }">
              <Select
                v-model="data.moduleField"
                :options="mappableFields"
                option-label="label"
                option-value="name"
                :placeholder="$t('block.recordList.import.pickModuleField')"
                :invalid="data.selected && !data.moduleField"
                class="w-full"
                @update:model-value="onFieldMapped(data)"
              />
              <small
                v-if="data.fileColumn === 'id'"
                class="text-muted-color block mt-1"
              >
                {{ $t('block.recordList.import.idFieldDescription') }}
              </small>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- Step 2: Progress -->
      <div v-else-if="step === 2" class="flex flex-col gap-4 py-2">
        <ProgressBar
          :value="progressPct"
          :show-value="true"
          style="height: 1.75rem"
        />
        <div class="flex items-center justify-between">
          <div v-if="!progress.finishedAt" class="flex items-center gap-2 text-muted-color text-sm">
            <ProgressSpinner style="width: 1.25rem; height: 1.25rem" stroke-width="6" />
            {{ $t('block.recordList.import.importing') }}
          </div>
          <div v-else-if="!progress.failed" class="text-green-600 text-sm font-medium">
            {{ $t('block.recordList.import.success') }}
          </div>
          <Button
            :label="progress.finishedAt ? $t('general.label.close') : $t('general.label.cancel')"
            severity="secondary"
            size="small"
            @click="closeDialog"
          />
        </div>
      </div>

      <!-- Step 3: Error Report -->
      <div v-else-if="step === 3" class="flex flex-col gap-5">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.import.report.title') }}
          </label>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1 text-sm pl-1">
            <div>
              <span class="text-muted-color">{{ $t('block.recordList.import.report.startedAt') }}:</span>
              <span class="font-medium ml-1">{{ fmtDate(progress.startedAt) }}</span>
            </div>
            <div>
              <span class="text-muted-color">{{ $t('block.recordList.import.report.finishedAt') }}:</span>
              <span class="font-medium ml-1">{{ fmtDate(progress.finishedAt) }}</span>
            </div>
            <div>
              <span class="text-muted-color">{{ $t('block.recordList.import.report.totalRecords') }}:</span>
              <span class="font-medium ml-1">{{ progress.entryCount }}</span>
            </div>
            <div>
              <span class="text-muted-color">{{ $t('block.recordList.import.report.importedRecords') }}:</span>
              <span class="font-medium ml-1 text-green-600">{{ progress.completed }}</span>
            </div>
            <div>
              <span class="text-muted-color">{{ $t('block.recordList.import.report.failedRecords') }}:</span>
              <span class="font-medium ml-1 text-red-500">{{ progress.failed }}</span>
            </div>
          </div>
        </div>

        <div v-if="errorRows.length" class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.import.report.detectedErrors') }}
          </label>
          <DataTable :value="errorRows" size="small" class="text-sm">
            <Column :header="$t('block.recordList.import.report.error')" field="message" />
            <Column :header="$t('block.recordList.import.report.count')" field="count" style="width: 8rem" />
          </DataTable>
        </div>

        <div v-if="failedEntries.length" class="flex flex-col gap-2">
          <label class="font-medium text-sm text-primary">
            {{ $t('block.recordList.import.report.failedEntries') }}
          </label>
          <div class="flex flex-col gap-1 text-sm pl-1">
            <div v-for="(ee, ix) in failedEntries" :key="ix">
              <template v-if="ee.length === 1 || ee[0] === ee[1]">
                <span class="text-muted-color">
                  {{ $t('block.recordList.import.report.failedEntriesLine') }}:
                </span>
                <span class="font-medium ml-1">{{ ee[0] }}</span>
              </template>
              <template v-else>
                <span class="text-muted-color">
                  {{ $t('block.recordList.import.report.failedEntriesLines') }}:
                </span>
                <span class="font-medium ml-1">{{ ee[0] }} → {{ ee[1] }}</span>
              </template>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <!-- Step 0: hidden (upload is the action) -->
        <div v-if="step === 1" class="flex items-center justify-between w-full">
          <Button
            :label="$t('general.label.back')"
            severity="secondary"
            text
            size="small"
            @click="step = 0"
          />
          <Button
            :label="$t('general.label.import')"
            size="small"
            :disabled="!canImport"
            @click="runImport"
          />
        </div>
        <div v-else-if="step === 3" class="flex items-center justify-end w-full">
          <Button
            :label="$t('general.label.close')"
            severity="secondary"
            text
            size="small"
            @click="closeDialog"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import moment from 'moment'
import { components, useFileUpload } from '@planetcrust/human-vue'

const { CFileDropZone } = components

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')

const props = defineProps({
  namespace: { type: Object, required: true },
  module: { type: Object, required: true },
})

const emit = defineEmits(['import-successful'])

// State
const showDialog = ref(false)
const step = ref(0)
const session = ref({})
const rows = ref([])
const progress = ref({})
let pollHandle = null

// Upload
const { uploading, uploadError, uploadFileRaw, reset: resetUpload } = useFileUpload()

// Form values
const onError = ref('SKIP')
const multiValueDelimiter = ref(';')

const onErrorOptions = computed(() => [
  { value: 'SKIP', label: t('block.recordList.import.onErrorSkip') },
  { value: 'FAIL', label: t('block.recordList.import.onErrorFail') },
])

const delimiterOptions = computed(() => [
  { value: ';', label: t('block.recordList.import.multiValueDelimiter.semicolon.label') },
  { value: ',', label: t('block.recordList.import.multiValueDelimiter.comma.label') },
  { value: '|', label: t('block.recordList.import.multiValueDelimiter.pipe.label') },
  { value: '[;]', label: t('block.recordList.import.multiValueDelimiter.semicolonArray.label') },
  { value: '[,]', label: t('block.recordList.import.multiValueDelimiter.commaArray.label') },
  { value: '[|]', label: t('block.recordList.import.multiValueDelimiter.pipeArray.label') },
])

const mappableFields = computed(() => {
  const moduleFields = (props.module.fields || [])
    .filter(f => f.kind !== 'File')
    .map(f => ({
      name: f.name,
      label: f.isRequired ? `${f.label || f.name}*` : (f.label || f.name),
      isRequired: f.isRequired,
    }))

  const systemFields = (props.module.systemFields?.() || []).map(f => ({
    name: f.name,
    label: f.label || f.name,
    isRequired: false,
  }))

  return [...moduleFields, ...systemFields]
})

const requiredFields = computed(() =>
  (props.module.fields || []).filter(f => f.isRequired),
)

const mappedRequiredFields = computed(() =>
  rows.value.filter(r => r.selected && requiredFields.value.some(rf => rf.name === r.moduleField)),
)

const hasUnmappedRequired = computed(() =>
  mappedRequiredFields.value.length < requiredFields.value.length,
)

const unmappedRequiredNames = computed(() => {
  const mapped = new Set(mappedRequiredFields.value.map(r => r.moduleField))
  return requiredFields.value
    .filter(f => !mapped.has(f.name))
    .map(f => f.label || f.name)
    .join(', ')
})

const canImport = computed(() => {
  const selected = rows.value.filter(r => r.selected)
  return (
    selected.length > 0 &&
    selected.every(r => !!r.moduleField) &&
    !hasUnmappedRequired.value
  )
})

const allSelected = computed(() => rows.value.length > 0 && rows.value.every(r => r.selected))

const progressPct = computed(() => {
  const { completed = 0, entryCount = 0 } = progress.value
  if (!entryCount) return 0
  return Math.round((completed / entryCount) * 100)
})

const errorRows = computed(() => {
  const errors = progress.value?.failLog?.errors || {}
  return Object.entries(errors).map(([message, count]) => ({ message, count }))
})

const failedEntries = computed(() => progress.value?.failLog?.records || [])

// Methods
function fmtDate(dt) {
  if (!dt) return '-'
  const m = moment(dt)
  return m.isValid() ? m.format('YYYY-MM-DD HH:mm:ss') : String(dt)
}

function onDialogHide() {
  if (step.value !== 2) {
    reset()
  }
}

function reset() {
  step.value = 0
  session.value = {}
  rows.value = []
  progress.value = {}
  resetUpload()
  stopPolling()
}

function closeDialog() {
  stopPolling()
  showDialog.value = false
  reset()
}

function toggleAll(val) {
  rows.value.forEach(r => (r.selected = val))
}

function onFieldMapped(row) {
  if (row.moduleField) {
    row.selected = true
  }
}

async function onFilesSelected(files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.recordImportInitEndpoint({
      namespaceID: props.namespace.namespaceID,
      moduleID: props.module.moduleID,
    })
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''

    const data = await uploadFileRaw(file, { url: endpoint, token })
    session.value = {
      ...data,
      namespaceID: props.namespace.namespaceID,
      moduleID: props.module.moduleID,
      onError: onError.value,
      multiValueDelimiter: multiValueDelimiter.value,
    }

    // Build field mapping rows with auto-match
    const fieldMap = {}
    ;(props.module.fields || []).forEach(f => { fieldMap[f.name] = f.name })
    ;(props.module.systemFields?.() || []).forEach(f => { fieldMap[f.name] = f.name })
    fieldMap['id'] = 'recordID'

    rows.value = Object.keys(data.fields || {}).map(fileColumn => {
      const auto = fieldMap[fileColumn] || data.fields[fileColumn] || null
      return { fileColumn, moduleField: auto, selected: !!auto }
    })

    step.value = 1
  } catch {
    // uploadError is set by the composable
  }
}

async function runImport() {
  if (!canImport.value) return

  const fields = {}
  rows.value.forEach(r => {
    if (r.selected && r.moduleField) {
      fields[r.fileColumn] = r.moduleField
    }
  })

  session.value.fields = fields
  session.value.onError = onError.value
  session.value.multiValueDelimiter = multiValueDelimiter.value

  step.value = 2
  progress.value = {}

  try {
    await $ComposeAPI.recordImportRun(session.value)
    startPolling()
  } catch (e) {
    progress.value = {
      startedAt: new Date().toISOString(),
      finishedAt: new Date().toISOString(),
      failed: 1,
      entryCount: 0,
      completed: 0,
      failLog: { errors: { [e.message || 'Unknown error']: 1 }, records: [] },
    }
    step.value = 3
  }
}

function startPolling() {
  stopPolling()
  pollHandle = setInterval(pollProgress, 2000)
  pollProgress()
}

function stopPolling() {
  if (pollHandle !== null) {
    clearInterval(pollHandle)
    pollHandle = null
  }
}

async function pollProgress() {
  try {
    const result = await $ComposeAPI.recordImportProgress(session.value)
    progress.value = result.progress || {}

    if (progress.value.finishedAt) {
      stopPolling()
      if (progress.value.failed) {
        step.value = 3
      } else {
        emit('import-successful')
        setTimeout(closeDialog, 1500)
      }
    }
  } catch {
    stopPolling()
  }
}
</script>
