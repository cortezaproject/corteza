<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">
      {{ recordModule.name }} —
      {{ isEditMode ? $t('general.label.edit') : $t('general.label.view') }}
    </span>
  </Teleport>

  <!-- Topbar navigation -->
  <Teleport to="#topbar-tools" :defer="true">
    <ButtonGroup v-if="recordModule" class="gap-1">
      <CRouterLinkButton
        :to="{ name: 'admin.modules.edit', params: { moduleID: recordModule.moduleID } }"
        :label="$t('module.edit.edit')"
        icon="pi pi-pencil"
        size="small"
      />
      <CRouterLinkButton
        :to="{ name: 'admin.modules.record.list', params: { moduleID: recordModule.moduleID } }"
        v-tooltip.bottom="$t('module.allRecords.label')"
        icon="pi pi-table"
        size="small"
      />
    </ButtonGroup>
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

  <!-- Record view/edit — same layout as public RecordView -->
  <Form
    ref="formRef"
    v-else-if="record"
    :resolver="resolver"
    @submit="handleSave"
    class="flex flex-col h-full"
  >
    <div class="flex-1 overflow-auto">
      <Grid :blocks="blocks" :namespace="namespace" :page="syntheticPage" :record="record" />
    </div>

    <!-- Record toolbar — same as public RecordView -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <!-- Left: Back / Cancel -->
        <div class="flex gap-2">
          <Button
            v-if="!isEditMode"
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="$router.back()"
          />
          <Button
            v-else
            :label="$t('general.label.cancel')"
            icon="pi pi-times"
            severity="secondary"
            :disabled="isSaving"
            @click="handleCancel"
          />
        </div>

        <!-- Right: Delete / Edit / Save -->
        <div class="flex gap-2">
          <CInputDelete
            v-if="record.canDeleteRecord"
            :label="$t('general.label.delete')"
            :message="$t('page.public.record.toolbar.deleteConfirm')"
            :header="recordModule.name"
            :disabled="isSaving || deleting"
            @confirm="handleDelete"
          />

          <Button
            v-if="!isEditMode && record.canUpdateRecord"
            :label="$t('general.label.edit')"
            icon="pi pi-pencil"
            @click="goToEdit"
          />

          <Button
            v-if="isEditMode"
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-check"
            :loading="isSaving"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import { validator } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import { computed, inject, nextTick, provide, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import {
  mergeAttachmentIDs,
  uploadRecordAttachment,
} from '@/sections/compose/lib/record-attachments'

const { CInputDelete, CRouterLinkButton } = components

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

const formRef = ref(null)
const serverErrors = ref({})
const loading = ref(false)
const isSaving = ref(false)
const deleting = ref(false)
const record = ref(null)
const pristineRecord = ref(null)
const navigatingAfterSave = ref(false)

const pendingByField = reactive(new Map())
provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    if (files.length > 0) pendingByField.set(fieldName, files)
    else pendingByField.delete(fieldName)
  },
})

const moduleID = computed(() => route.params.moduleID)
const recordID = computed(() => route.params.recordID)
const recordModule = computed(() => (moduleID.value ? moduleStore.getByID(moduleID.value) : null))
const isEditMode = computed(() => route.name === 'admin.modules.record.edit')
const mode = computed(() => (isEditMode.value ? 'edit' : 'view'))
const isNew = computed(() => false)

// Provide recordViewContext so RecordBlock can inject it
provide('recordViewContext', {
  mode,
  record,
  isNew,
  isSaving,
})

// Synthetic blocks: a single Record block showing all fields
const blocks = computed(() => [
  {
    blockID: '_admin_record_view',
    kind: 'Record',
    title: '',
    description: '',
    style: { wrap: { kind: 'card' } },
    options: {
      fields: [],
    },
    xywh: [0, 0, 48, 18],
    meta: { tempID: '_admin_record_view' },
  },
])

const syntheticPage = computed(() => ({
  pageID: '0',
  title: recordModule.value?.name || '',
  moduleID: moduleID.value,
  blocks: [],
}))

async function loadRecord() {
  if (!recordModule.value || !recordID.value) return

  loading.value = true
  record.value = null
  pristineRecord.value = null

  try {
    const loaded = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: moduleID.value,
      recordID: recordID.value,
      force: true,
    })
    pristineRecord.value = loaded
    record.value = isEditMode.value ? loaded.clone() : loaded
  } catch (e) {
    console.error('Failed to load record:', e)
    $toast.toastDanger(t('notification.record.loadFailed'))
  } finally {
    loading.value = false
  }
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

async function handleSave({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }
  if (!record.value) return

  isSaving.value = true

  try {
    for (const [fieldName, files] of pendingByField) {
      const ids = await Promise.all(
        files.map(file =>
          uploadRecordAttachment($ComposeAPI, {
            namespaceID: props.namespace.namespaceID,
            moduleID: moduleID.value,
            recordID: record.value.recordID,
            fieldName,
            file,
          }),
        ),
      )
      record.value.setValue(fieldName, mergeAttachmentIDs(record.value.values[fieldName], ids))
    }

    const saved = await recordStore.update(record.value)

    $toast.toastSuccess(t('notification.record.updateSuccess'))

    pristineRecord.value = saved
    navigatingAfterSave.value = true
    router.replace({
      name: 'admin.modules.record.view',
      params: { moduleID: moduleID.value, recordID: saved.recordID },
    })
  } catch (e) {
    console.error('Failed to save record:', e)
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
      $toast.toastErrorHandler(t('notification.record.updateFailed'))(e)
    }
  } finally {
    isSaving.value = false
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
    $toast.toastErrorHandler(t('notification.record.deleteFailed'))(e)
  } finally {
    deleting.value = false
  }
}

function handleCancel() {
  navigatingAfterSave.value = true
  router.replace({
    name: 'admin.modules.record.view',
    params: { moduleID: moduleID.value, recordID: recordID.value },
  })
}

function goToEdit() {
  router.push({
    name: 'admin.modules.record.edit',
    params: { moduleID: moduleID.value, recordID: recordID.value },
  })
}

// Guard against navigating away with unsaved changes
onBeforeRouteLeave(() => {
  if (navigatingAfterSave.value) {
    navigatingAfterSave.value = false
    return true
  }
  if (isEditMode.value && !isSaving.value) {
    return window.confirm(t('general.editor.unsavedChanges'))
  }
})

// Handle view↔edit in-place without reloading
watch(
  () => mode.value,
  (newMode, oldMode) => {
    if (newMode === 'edit' && oldMode === 'view' && pristineRecord.value) {
      record.value = pristineRecord.value.clone()
    } else if (newMode === 'view' && oldMode === 'edit') {
      record.value = pristineRecord.value
    }
  },
)

// Clear server errors on edit
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

// Only reload when the actual record/module changes
watch(
  () => [recordModule.value?.moduleID, recordID.value],
  () => {
    if (recordModule.value) loadRecord()
  },
  { immediate: true },
)
</script>
