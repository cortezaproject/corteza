<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">{{ recordModule.name }} — {{ $t('module.edit.createRecord') }}</span>
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

  <!-- Module not found -->
  <div v-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Module has no fields — records can't be created until fields are added -->
  <div v-else-if="!recordModule.fields?.length" class="flex items-center justify-center h-full">
    <Message severity="info" :closable="false">
      {{ $t('module.edit.createRecordNoFields') }}
    </Message>
  </div>

  <!-- Record create — same layout as public RecordView in create mode -->
  <Form
    ref="formRef"
    v-else-if="record"
    :resolver="resolver"
    @submit="handleSave"
    class="flex flex-col h-full"
  >
    <div
      class="flex-1 overflow-auto"
      @pointerdown.capture="freezeBaseline"
      @keydown.capture="freezeBaseline"
    >
      <Grid :blocks="blocks" :namespace="namespace" :page="syntheticPage" :record="record" />
    </div>

    <!-- Record toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.cancel')"
          icon="pi pi-times"
          severity="secondary"
          :disabled="isSaving"
          @click="handleCancel"
        />
        <Button
          type="submit"
          :label="$t('general.label.save')"
          icon="pi pi-check"
          :loading="isSaving"
        />
      </div>
    </div>
  </Form>
</template>

<script setup>
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import { compose, validator } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'
import { computed, inject, nextTick, provide, reactive, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import {
  mergeAttachmentIDs,
  uploadRecordAttachment,
} from '@/sections/compose/lib/record-attachments'

const { CRouterLinkButton } = components

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
const $Auth = inject('$Auth', {})
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const formRef = ref(null)
const serverErrors = ref({})
const isSaving = ref(false)
const isCancelling = ref(false)
const record = ref(null)

// The form as the user was shown it: what a blank record, a clone source, a
// prefilled reference field and the field editors' own presets add up to. A User
// field set to preset-with-authenticated fills itself in, and that is not an
// unsaved edit — anything changing after the user's first touch is.
const initialValues = ref(null)

const pendingByField = reactive(new Map())
provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    if (files.length > 0) pendingByField.set(fieldName, files)
    else pendingByField.delete(fieldName)
  },
})

const moduleID = computed(() => route.params.moduleID)
const recordModule = computed(() => (moduleID.value ? moduleStore.getByID(moduleID.value) : null))

const isMultiField = fieldName =>
  !!recordModule.value?.fields.find(f => f.name === fieldName)?.isMulti

const mode = ref('create')
const isNew = ref(true)

// Provide recordViewContext so RecordBlock can inject it
provide('recordViewContext', {
  mode,
  record,
  isNew,
  isSaving,
})

const blocks = computed(() => [
  {
    blockID: '_admin_record_create',
    kind: 'Record',
    title: '',
    description: '',
    style: { wrap: { kind: 'card' } },
    options: {
      fields: [],
    },
    xywh: [0, 0, 48, 18],
    meta: { tempID: '_admin_record_create' },
  },
])

const syntheticPage = computed(() => ({
  pageID: '0',
  title: recordModule.value?.name || '',
  moduleID: moduleID.value,
  blocks: [],
}))

function setRecord(rec) {
  record.value = rec
  initialValues.value = null
}

// Taken on the first pointer or key event in the form, before the value that
// event carries lands. Presets resolve on their own schedule — some after a
// round trip — so no fixed moment after init is reliably "the form as shown",
// but everything before the user's first touch of it is.
function freezeBaseline() {
  if (initialValues.value || !record.value) return
  initialValues.value = cloneDeep(toRaw(record.value).values)
}

const isDirty = computed(() => {
  if (isSaving.value || isCancelling.value) return false
  if (pendingByField.size) return true
  // Untouched: the presets a form fills in for itself are not unsaved work.
  if (!record.value || !initialValues.value) return false
  // Read through the reactive record, not toRaw: a raw read registers no
  // dependency and the computed would never see the user's edits.
  return !isEqual(record.value.values, initialValues.value)
})

const { markSaved } = useUnsavedGuard({
  isDirty,
  messageKey: 'general.editor.unsavedChanges',
})

function initRecord() {
  // A record needs a module WITH fields (compose.Record throws otherwise).
  if (!recordModule.value?.fields?.length) return

  const prefillRefField = record => {
    const refField = route.query.refField
    const refValue = route.query.refValue
    if (!record || !refField || !refValue) return

    const field = recordModule.value.fields.find(f => f.name === refField)
    if (!field) return

    const valueArray = Array.isArray(refValue) ? refValue : [refValue]
    if (field.isMulti) {
      record.setValue(refField, valueArray)
    } else {
      record.setValue(refField, valueArray[0])
    }
  }

  // Handle clone
  if (route.query.cloneFromID) {
    recordStore
      .findByID({
        namespaceID: props.namespace.namespaceID,
        moduleID: moduleID.value,
        recordID: route.query.cloneFromID,
      })
      .then(source => {
        const newRec = new compose.Record(recordModule.value)
        for (const field of recordModule.value.fields) {
          newRec.setValue(field.name, source.values[field.name])
        }
        // Prefill ownedBy with current user
        newRec.ownedBy = $Auth?.user?.userID || undefined
        prefillRefField(newRec)
        setRecord(newRec)
      })
      .catch(e => {
        console.error('Failed to load source record for clone:', e)
        const newRec = new compose.Record(recordModule.value, { ownedBy: $Auth?.user?.userID })
        prefillRefField(newRec)
        setRecord(newRec)
      })
  } else {
    const newRec = new compose.Record(recordModule.value, { ownedBy: $Auth?.user?.userID })
    prefillRefField(newRec)
    setRecord(newRec)
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
            recordID: '',
            fieldName,
            file,
          }),
        ),
      )
      record.value.setValue(
        fieldName,
        mergeAttachmentIDs(record.value.values[fieldName], ids, isMultiField(fieldName)),
      )
    }

    const saved = await recordStore.create(record.value)
    $toast.toastSuccess(t('notification.record.createSuccess'))
    markSaved()
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
      $toast.toastErrorHandler(t('notification.record.createFailed'))(e)
    }
  } finally {
    isSaving.value = false
  }
}

function handleCancel() {
  // Cancel is the deliberate discard — the guard has nothing to warn about.
  isCancelling.value = true
  router.replace({
    name: 'admin.modules.record.list',
    params: { moduleID: moduleID.value },
  })
}

watch(
  () => [
    recordModule.value?.moduleID,
    route.query.cloneFromID,
    route.query.refField,
    route.query.refValue,
  ],
  () => {
    if (recordModule.value) initRecord()
  },
  { immediate: true },
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
</script>
