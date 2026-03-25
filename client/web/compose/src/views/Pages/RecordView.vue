<template>
  <!-- Page title in topbar -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="page">{{ page.title }}</span>
  </Teleport>

  <!-- Admin tools in topbar -->
  <Teleport to="#topbar-tools" :defer="true">
    <ButtonGroup v-if="page?.canUpdatePage" class="gap-1">
      <Button
        v-if="page.isRecordPage"
        :label="$t('page.moduleEdit')"
        icon="pi pi-database"
        size="small"
        @click="goToModuleEdit"
      />
      <Button
        :label="$t('page.block.general.label.pageBuilder')"
        icon="pi pi-wrench"
        size="small"
        @click="goToBuilder"
      />
      <Button
        v-tooltip.bottom="$t('navigation.editPage')"
        icon="pi pi-objects-column"
        size="small"
        @click="goToEditPage"
      />
    </ButtonGroup>
  </Teleport>

  <!-- Loading state -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner style="width: 32px; height: 32px" />
  </div>

  <!-- Page content with record context -->
  <Form
    ref="formRef"
    v-else-if="page && positionedBlocks.length"
    :resolver="resolver"
    @submit="handleSave"
    class="flex flex-col h-full"
  >
    <div class="flex-1 overflow-auto">
      <Grid :blocks="positionedBlocks" :namespace="namespace" :page="page" :record="record" />
    </div>

    <!-- Record Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <!-- Left side -->
        <div class="flex gap-2">
          <Button
            v-if="mode === 'view' && layoutButtons.back"
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="$router.back()"
          />
          <Button
            v-else-if="mode !== 'view'"
            :label="$t('general.label.cancel')"
            icon="pi pi-times"
            severity="secondary"
            :disabled="isSaving"
            @click="handleCancel"
          />
        </div>

        <!-- Right side -->
        <div class="flex gap-2">
          <!-- Delete button (view mode or edit mode for existing records) -->
          <CInputDelete
            v-if="showDeleteButton"
            :label="$t('general.label.delete')"
            :message="$t('page.public.record.toolbar.deleteConfirm')"
            :header="page.title"
            :disabled="deleting || isSaving"
            @confirm="handleDelete"
          />

          <!-- Clone / Save as copy (view mode, existing record) -->
          <Button
            v-if="mode === 'view' && !isNew && record && layoutButtons.clone"
            :label="$t('general.label.saveAsCopy')"
            icon="pi pi-copy"
            severity="secondary"
            @click="handleClone"
          />

          <!-- New record button (view mode) -->
          <Button
            v-if="mode === 'view' && layoutButtons.new"
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            severity="secondary"
            @click="handleNew"
          />

          <!-- Edit button (view mode only) -->
          <Button
            v-if="mode === 'view' && record?.canUpdateRecord && layoutButtons.edit"
            :label="$t('general.label.edit')"
            icon="pi pi-pencil"
            severity="primary"
            @click="handleEdit"
          />

          <!-- Save button (edit/create mode) -->
          <Button
            v-if="mode !== 'view' && layoutButtons.submit"
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-check"
            :loading="isSaving"
          />
        </div>
      </div>
    </div>
  </Form>

  <!-- No blocks -->
  <div v-else-if="page" class="flex items-center justify-center h-full">
    <p class="text-muted-color">{{ $t('page.noBlock') }}</p>
  </div>

  <!-- Page not found -->
  <div v-else class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('page.invalid') }}
    </Message>
  </div>
</template>

<script setup>
import Grid from '@/components/PageBlocks/Grid.vue'
import { useModuleStore } from '@/stores/module'
import { usePageLayoutStore } from '@/stores/page-layout'
import { usePageStore } from '@/stores/page'
import { useRecordStore } from '@/stores/record'
import { compose, validator } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, nextTick, provide, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'

const { CInputDelete } = components

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
const $auth = inject('$auth', {})
const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const formRef = ref(null)
const serverErrors = ref({})

const loading = ref(false)
const deleting = ref(false)
const isSaving = ref(false)
const page = ref(null)
const layout = ref(null)
const record = ref(null)
const pristineRecord = ref(null)
const navigatingAfterSave = ref(false)

// Mode derived from route
const mode = computed(() => {
  const recordID = route.params.recordID
  if (recordID === '0') return 'create'
  if (route.query.edit === '1') return 'edit'
  return 'view'
})

const isNew = computed(() => mode.value === 'create')

// Button visibility from layout config (all enabled by default when no layout)
const layoutButtons = computed(() => {
  const b = layout.value?.config?.buttons ?? {}
  return {
    back: b.back?.enabled ?? true,
    delete: b.delete?.enabled ?? true,
    clone: b.clone?.enabled ?? true,
    new: b.new?.enabled ?? true,
    edit: b.edit?.enabled ?? true,
    submit: b.submit?.enabled ?? true,
  }
})

// Show delete button: in view mode if record has permission, or in edit mode for existing record
const showDeleteButton = computed(() => {
  if (!record.value) return false
  if (!layoutButtons.value.delete) return false
  if (mode.value === 'view') return record.value.canDeleteRecord
  if (mode.value === 'edit') return record.value.canDeleteRecord
  return false
})

const pendingByField = reactive(new Map())

// Provide context so child blocks (RecordBlock) can inject it
provide('recordViewContext', {
  mode,
  record,
  isNew,
  isSaving,
})

provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    if (files.length > 0) pendingByField.set(fieldName, files)
    else pendingByField.delete(fieldName)
  },
})

const positionedBlocks = computed(() => {
  if (!page.value || !layout.value) {
    if (page.value?.blocks?.length) {
      return page.value.blocks
    }
    return []
  }

  return layout.value.blocks
    .map(layoutBlock => {
      const pageBlock = page.value.blocks.find(b => b.blockID === layoutBlock.blockID)
      if (!pageBlock) return null

      return {
        ...pageBlock,
        xywh: layoutBlock.xywh || pageBlock.xywh,
      }
    })
    .filter(Boolean)
})

async function loadPage() {
  const pageID = route.params.pageID
  const recordID = route.params.recordID
  if (!pageID) return

  loading.value = true
  record.value = null
  pristineRecord.value = null

  try {
    page.value = pageStore.getByID(pageID) || null

    if (page.value) {
      const layouts = pageLayoutStore.getByPageID(pageID)
      layout.value = layouts.length > 0 ? layouts[0] : null

      const moduleID = page.value.moduleID
      if (moduleID) {
        const mod = moduleStore.getByID(moduleID)
        if (mod) {
          if (mode.value === 'create') {
            // Handle clone
            if (route.query.cloneFromID) {
              try {
                const source = await recordStore.findByID({
                  namespaceID: mod.namespaceID,
                  moduleID: mod.moduleID,
                  recordID: route.query.cloneFromID,
                })
                const newRec = new compose.Record(mod)
                for (const field of mod.fields) {
                  newRec.setValue(field.name, source.values[field.name])
                }
                // Prefill ownedBy with current user
                newRec.ownedBy = $auth?.user?.userID || undefined
                record.value = newRec
              } catch (e) {
                console.error('Failed to load source record for clone:', e)
                record.value = new compose.Record(mod, { ownedBy: $auth?.user?.userID })
              }
            } else {
              record.value = new compose.Record(mod, { ownedBy: $auth?.user?.userID })
            }
          } else if (recordID && recordID !== '0') {
            try {
              const loaded = await recordStore.findByID({
                namespaceID: mod.namespaceID,
                moduleID: mod.moduleID,
                recordID,
                force: true,
              })
              pristineRecord.value = loaded
              // In edit mode, start with a clone; mode watch will re-clone on later transitions
              record.value = mode.value === 'edit' ? loaded.clone() : loaded
            } catch (e) {
              console.error('Failed to load record:', e)
              record.value = null
            }
          }
        }
      }
    }
  } finally {
    loading.value = false
  }
}

function resolver() {
  const errors = {}

  for (const [fieldName, message] of Object.entries(serverErrors.value)) {
    errors[fieldName] = [{ message }]
  }

  const recordModule = page.value ? moduleStore.getByID(page.value.moduleID) : null
  if (!record.value || !recordModule) return { errors }
  for (const field of recordModule.fields) {
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
  const { data } = await $ComposeAPI
    .api()
    .post(url, formData, { headers: { 'Content-Type': undefined } })
  if (data?.error) throw new Error(data.error)
  const attachment = data?.response ?? data
  if (!attachment?.attachmentID)
    throw new Error(`Upload failed for "${file.name}": no attachmentID in response`)
  return attachment.attachmentID
}

async function handleSave({ valid }) {
  if (!valid) return
  if (!record.value || !page.value) return

  isSaving.value = true

  try {
    const namespaceID = props.namespace.namespaceID
    const moduleID = page.value.moduleID
    const recordID = isNew.value ? '' : record.value.recordID

    for (const [fieldName, files] of pendingByField) {
      const ids = await Promise.all(
        files.map(file => uploadFile({ namespaceID, moduleID, recordID, fieldName, file })),
      )
      const existing = record.value.values[fieldName]
      const existingIDs = Array.isArray(existing)
        ? existing.filter(Boolean)
        : existing
          ? [existing]
          : []
      record.value.setValue(fieldName, [...existingIDs, ...ids])
    }

    const saved = isNew.value
      ? await recordStore.create(record.value)
      : await recordStore.update(record.value)

    $toast.toastSuccess(
      t(isNew.value ? 'notification.record.createSuccess' : 'notification.record.updateSuccess'),
    )

    pristineRecord.value = saved
    // Set flag so leave guard allows this programmatic navigation
    navigatingAfterSave.value = true
    router.replace({
      name: 'page.record',
      params: {
        slug: route.params.slug,
        pageID: route.params.pageID,
        recordID: saved.recordID,
      },
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
      $toast.toastDanger(
        t(isNew.value ? 'notification.record.createFailed' : 'notification.record.updateFailed'),
      )
    }
  } finally {
    isSaving.value = false
  }
}

function handleEdit() {
  router.push({ query: { edit: '1' } })
}

function handleClone() {
  router.push({
    name: 'page.record',
    params: { slug: route.params.slug, pageID: route.params.pageID, recordID: '0' },
    query: { cloneFromID: record.value.recordID },
  })
}

function handleNew() {
  router.push({
    name: 'page.record',
    params: { slug: route.params.slug, pageID: route.params.pageID, recordID: '0' },
  })
}

function handleCancel() {
  if (mode.value === 'create') {
    router.back()
  } else {
    // Return to view mode (remove edit query)
    router.replace({ query: {} })
  }
}

async function handleDelete() {
  if (!record.value || !page.value) return

  deleting.value = true
  try {
    await recordStore.delete({
      namespaceID: record.value.namespaceID,
      moduleID: record.value.moduleID,
      recordID: record.value.recordID,
    })
    $toast.toastSuccess(t('notification.record.deleteSuccess'))
    router.back()
  } catch (e) {
    console.error('Failed to delete record:', e)
    $toast.toastDanger(t('notification.record.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function goToBuilder() {
  if (page.value) {
    router.push({
      name: 'admin.pages.builder',
      params: { pageID: page.value.pageID },
    })
  }
}

function goToEditPage() {
  if (page.value) {
    router.push({
      name: 'admin.pages.edit',
      params: { pageID: page.value.pageID },
    })
  }
}

function goToModuleEdit() {
  if (page.value?.moduleID) {
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: page.value.moduleID },
    })
  }
}

// Guard against navigating away with unsaved changes
onBeforeRouteLeave(() => {
  if (navigatingAfterSave.value) {
    navigatingAfterSave.value = false
    return true
  }
  if (mode.value !== 'view' && !isSaving.value) {
    return window.confirm(t('general.record.unsavedChanges'))
  }
})

// Handle view↔edit transitions in-place without reloading
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

// Load on mount and when pageID, recordID, or cloneFromID changes (not on edit-query toggle)
watch(
  () => [route.params.pageID, route.params.recordID, route.query.cloneFromID],
  () => loadPage(),
  { immediate: true },
)
</script>
