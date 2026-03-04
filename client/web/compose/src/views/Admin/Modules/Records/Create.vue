<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">
      {{ recordModule.name }} — {{ $t('module.edit.createRecord') }}
    </span>
  </Teleport>

  <!-- Module not found -->
  <div v-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Create form -->
  <Form ref="formRef" v-else-if="recordModule" :resolver="resolver" @submit="handleSubmit" class="flex flex-col h-full">
    <div class="flex-1 overflow-auto">
      <div class="container mx-auto p-4 max-w-3xl">
        <div class="flex flex-col gap-5">
          <div
            v-for="field in recordModule.fields"
            :key="field.fieldID || field.name"
          >
            <label v-if="field.kind !== 'Bool' || field.options?.switch" class="text-sm font-semibold text-primary mb-1.5 block">
              {{ field.label || field.name }}
              <span v-if="field.isRequired" class="text-red-500 ml-0.5">*</span>
            </label>

            <FormField :name="field.name" v-slot="{ invalid, error }">
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
          </div>
        </div>
      </div>
    </div>

    <!-- Footer toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.cancel')"
          icon="pi pi-times"
          severity="secondary"
          :disabled="saving"
          @click="handleCancel"
        />
        <Button
          type="submit"
          :label="$t('general.label.save')"
          icon="pi pi-check"
          :loading="saving"
        />
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, provide, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { compose, validator } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
const { CFieldEditor } = components
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

const formRef = ref(null)
const serverErrors = ref({})

const saving = ref(false)
const record = ref(null)

const moduleID = computed(() => route.params.moduleID)
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

function initRecord() {
  if (!recordModule.value) return
  record.value = new compose.Record(recordModule.value)
}

function resolver() {
  const errors = {}

  for (const [fieldName, message] of Object.entries(serverErrors.value)) {
    errors[fieldName] = [{ message }]
  }

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
          recordID: '',
          fieldName,
          file,
        }))
      )
      const existing = record.value.values[fieldName]
      const existingIDs = Array.isArray(existing) ? existing.filter(Boolean) : (existing ? [existing] : [])
      record.value.setValue(fieldName, [...existingIDs, ...ids])
    }

    const saved = await recordStore.create(record.value)
    $toast.toastSuccess(t('notification.record.createSuccess'))
    router.replace({
      name: 'admin.modules.record.view',
      params: { moduleID: moduleID.value, recordID: saved.recordID },
    })
  } catch (e) {
    console.error('Failed to create record:', e)
    const details = e?.details ?? []
    const fieldErrors = {}
    for (const detail of details) {
      if (detail.meta?.field) {
        fieldErrors[detail.meta.field] = detail.message
      }
    }
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
      await nextTick()
      formRef.value?.validate()
    } else {
      $toast.toastDanger(t('notification.record.createFailed'))
    }
    saving.value = false
  }
}

function handleCancel() {
  router.replace({
    name: 'admin.modules.record.list',
    params: { moduleID: moduleID.value },
  })
}

watch(
  () => recordModule.value?.moduleID,
  () => {
    if (recordModule.value) initRecord()
  },
  { immediate: true },
)

// Clear server errors as soon as the user edits anything
watch(
  () => record.value?.values,
  () => {
    if (Object.keys(serverErrors.value).length > 0) {
      serverErrors.value = {}
      nextTick(() => formRef.value?.validate())
    }
  },
  { deep: true },
)
</script>
