<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">
      {{ recordModule.name }} —
      {{ isEditMode ? $t('general.label.edit') : $t('general.label.view') }}
    </span>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner style="width: 32px; height: 32px" />
  </div>

  <!-- Module not found -->
  <div v-else-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Record form -->
  <Form v-else :resolver="resolver" @submit="handleSubmit" class="flex flex-col h-full">
    <div class="flex-1 overflow-auto">
      <div class="container mx-auto p-4 max-w-3xl">
        <div class="flex flex-col gap-5">
          <div
            v-for="field in recordModule.fields"
            :key="field.fieldID || field.name"
          >
            <label v-if="field.kind !== 'Bool' || field.options?.switch || !isEditMode" class="text-sm font-semibold text-primary mb-1.5 block">
              {{ field.label || field.name }}
              <span v-if="field.isRequired" class="text-red-500 ml-0.5">*</span>
            </label>

            <!-- Editor -->
            <FormField v-if="isEditMode" :name="field.name" v-slot="{ invalid, error }">
              <CFieldEditor
                :field="field"
                :namespace="namespace"
                :model-value="getFieldValue(field)"
                @update:model-value="setFieldValue(field, $event)"
              />
              <Message v-if="invalid" severity="error" size="small" variant="simple">
                {{ error?.message }}
              </Message>
            </FormField>

            <!-- Viewer -->
            <div v-else class="text-color">
              <CFieldViewer
                v-if="field.canReadRecordValue !== false"
                :field="field"
                :record="record"
                :namespace="namespace"
              />
              <span v-else class="text-muted-color italic text-sm">
                {{ $t('block.field.noPermission') }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <!-- Left: Cancel/Back -->
        <div class="flex gap-2">
          <Button
            v-if="isEditMode"
            :label="$t('general.label.cancel')"
            icon="pi pi-times"
            severity="secondary"
            :disabled="saving"
            @click="handleCancel"
          />
          <Button
            v-else
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="$router.back()"
          />
        </div>

        <!-- Right: actions -->
        <div class="flex gap-2">
          <!-- Delete -->
          <CInputDelete
            v-if="record?.canDeleteRecord"
            :label="$t('general.label.delete')"
            :message="$t('page.public.record.toolbar.deleteConfirm')"
            :header="recordModule.name"
            :disabled="saving || deleting"
            @confirm="handleDelete"
          />

          <!-- Edit button (view mode only) -->
          <Button
            v-if="!isEditMode && record?.canUpdateRecord"
            :label="$t('general.label.edit')"
            icon="pi pi-pencil"
            @click="goToEdit"
          />

          <!-- Save button (edit mode only) -->
          <Button
            v-if="isEditMode"
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-check"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, provide, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { validator } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer, CFieldEditor, CInputDelete } = components
import { useModuleStore } from '@/stores/module'
import { useRecordStore } from '@/stores/record'

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const pendingByField = reactive(new Map())
provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    if (files.length > 0) pendingByField.set(fieldName, files)
    else pendingByField.delete(fieldName)
  },
})

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const record = ref(null)

const moduleID = computed(() => route.params.moduleID)
const recordID = computed(() => route.params.recordID)
const isEditMode = computed(() => route.name === 'admin.modules.record.edit')

const recordModule = computed(() => (moduleID.value ? moduleStore.getByID(moduleID.value) : null))

function getFieldValue(field) {
  if (!record.value) return field.isMulti ? [] : ''
  const val = record.value.values[field.name]
  if (val === undefined || val === null) return field.isMulti ? [] : ''
  return val
}

function setFieldValue(field, value) {
  if (!record.value) return
  record.value.setValue(field.name, value)
}

async function loadRecord() {
  if (!recordModule.value || !recordID.value) return

  loading.value = true
  record.value = null

  try {
    const loaded = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: moduleID.value,
      recordID: recordID.value,
      force: true,
    })
    record.value = loaded.clone()
  } catch (e) {
    console.error('Failed to load record:', e)
    $toast.toastDanger(t('notification.record.loadFailed'))
  } finally {
    loading.value = false
  }
}

function resolver() {
  const errors = {}
  if (!record.value || !recordModule.value) return { errors }
  for (const field of recordModule.value.fields) {
    if (field.isRequired) {
      const val = record.value.values[field.name]
      if (validator.IsEmpty(val)) {
        errors[field.name] = [{ message: t('field.required-field') }]
      }
    }
  }
  return { errors }
}

async function uploadFile({ namespaceID, moduleID, recordID, fieldName, file }) {
  const url = $ComposeAPI.recordUploadEndpoint({ namespaceID, moduleID })
  const formData = new FormData()
  formData.append('recordID', recordID || '')
  formData.append('fieldName', fieldName)
  formData.append('upload', file, file.name)
  const { data } = await $ComposeAPI.api().post(url, formData, { headers: { 'Content-Type': undefined } })
  if (data?.error) throw new Error(data.error)
  const attachment = data?.response ?? data
  if (!attachment?.attachmentID) throw new Error(`Upload failed for "${file.name}": no attachmentID in response`)
  return attachment.attachmentID
}

async function handleSubmit({ valid }) {
  if (!valid) return
  if (!record.value) return

  saving.value = true

  try {
    for (const [fieldName, files] of pendingByField) {
      const ids = await Promise.all(
        files.map(file => uploadFile({
          namespaceID: props.namespace.namespaceID,
          moduleID: moduleID.value,
          recordID: recordID.value,
          fieldName,
          file,
        }))
      )
      const existing = record.value.values[fieldName]
      const existingIDs = Array.isArray(existing) ? existing.filter(Boolean) : (existing ? [existing] : [])
      record.value.setValue(fieldName, [...existingIDs, ...ids])
    }

    const saved = await recordStore.update(record.value)
    $toast.toastSuccess(t('notification.record.updateSuccess'))
    router.replace({
      name: 'admin.modules.record.view',
      params: { moduleID: moduleID.value, recordID: saved.recordID },
    })
  } catch (e) {
    console.error('Failed to save record:', e)
    $toast.toastDanger(t('notification.record.updateFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!record.value) return

  deleting.value = true
  try {
    await recordStore.delete({
      namespaceID: props.namespace.namespaceID,
      moduleID: record.value.moduleID,
      recordID: record.value.recordID,
    })
    $toast.toastSuccess(t('notification.record.deleteSuccess'))
    router.replace({
      name: 'admin.modules.record.list',
      params: { moduleID: moduleID.value },
    })
  } catch (e) {
    console.error('Failed to delete record:', e)
    $toast.toastDanger(t('notification.record.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function handleCancel() {
  router.replace({
    name: 'admin.modules.record.view',
    params: { moduleID: moduleID.value, recordID: recordID.value },
  })
  // Reload to restore pristine record for view mode (route.name no longer in watch)
  loadRecord()
}

function goToEdit() {
  // Re-clone so in-flight edits don't bleed back into view mode on cancel
  if (record.value) {
    record.value = record.value.clone()
  }
  router.push({
    name: 'admin.modules.record.edit',
    params: { moduleID: moduleID.value, recordID: recordID.value },
  })
}

// Only reload when the actual record/module changes, not on view↔edit route-name toggle
watch(
  () => [recordModule.value?.moduleID, recordID.value],
  () => {
    if (recordModule.value) loadRecord()
  },
  { immediate: true },
)
</script>
