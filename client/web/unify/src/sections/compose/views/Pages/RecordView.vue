<template>
  <!-- Page title in topbar -->
  <Teleport v-if="!inModal" to="#topbar-title" :defer="true">
    <span v-if="page">{{ page.title }}</span>
  </Teleport>

  <!-- Admin tools in topbar -->
  <Teleport v-if="!inModal" to="#topbar-tools" :defer="true">
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
        icon="pi pi-pencil"
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

    <!-- Record Toolbar (only on record pages) -->
    <div v-if="page?.isRecordPage" class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <!-- Left side -->
        <div class="flex gap-2">
          <Button
            v-if="mode === 'view' && layoutButtons.back"
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="handleCancel"
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

        <!-- Center: prev/next navigation -->
        <div
          v-if="mode === 'view' && (recordNavigation.prev || recordNavigation.next)"
          class="flex gap-1"
        >
          <Button
            icon="pi pi-chevron-left"
            severity="secondary"
            :disabled="!recordNavigation.prev || navigating !== null"
            :loading="navigating === 'prev'"
            :title="$t('general.recordNavigation.prev')"
            @click="navigateToRecord(recordNavigation.prev, 'prev')"
          />
          <Button
            icon="pi pi-chevron-right"
            severity="secondary"
            :disabled="!recordNavigation.next || navigating !== null"
            :loading="navigating === 'next'"
            :title="$t('general.recordNavigation.next')"
            @click="navigateToRecord(recordNavigation.next, 'next')"
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
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import { fetchBlockID, usePageVisibility } from '@/sections/compose/composables/usePageVisibility'
import { useModuleStore } from '@planetcrust/human-vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import { compose, validator } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import { computed, inject, nextTick, onBeforeUnmount, provide, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'

const { CInputDelete } = components

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
  inModal: {
    type: Boolean,
    default: false,
  },
  modalPageID: {
    type: String,
    default: null,
  },
  modalRecordID: {
    type: String,
    default: null,
  },
})

const emit = defineEmits(['close'])

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI', null)
const $auth = inject('$auth', {})
const $eventBus = inject('$eventBus', null)

const { buildExpressionVariables, determineLayout, evaluateBlocks } = usePageVisibility($SystemAPI, $auth)

const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const formRef = ref(null)
const serverErrors = ref({})

const loading = ref(false)
const invisibleBlockIDs = ref(new Set())

const recordNavigation = computed(() => {
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (!recordID || recordID === '0') return {}
  // Access the ref directly so Vue tracks it as a reactive dependency
  const ids = recordStore.paginationRecordIDs
  const idx = ids.indexOf(recordID)
  if (idx === -1) return {}
  return {
    prev: idx > 0 ? ids[idx - 1] : undefined,
    next: idx < ids.length - 1 ? ids[idx + 1] : undefined,
  }
})

function navigateToRecord(targetRecordID, direction) {
  navigating.value = direction
  if (!targetRecordID) return
  if (props.inModal) {
    router.push({ query: { ...route.query, recordID: targetRecordID, edit: undefined } })
  } else {
    router.push({
      name: 'page.record',
      params: { slug: route.params.slug, pageID: route.params.pageID, recordID: targetRecordID },
    })
  }
}
const deleting = ref(false)
const isSaving = ref(false)
const page = ref(null)
const layout = ref(null)
const record = ref(null)
const pristineRecord = ref(null)
const navigatingAfterSave = ref(false)

// Mode derived from route or props
const mode = computed(() => {
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
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
  const blocks = (() => {
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
  })()

  // meta.hidden is handled by Grid (tab children must still reach TabsBlock via props.blocks)
  // invisibleBlockIDs are blocks hidden by visibility expressions/roles — remove entirely
  return blocks.filter(b => !invisibleBlockIDs.value.has(fetchBlockID(b)))
})

const navigating = ref(null) // 'prev' | 'next' | null

// Cancels the in-flight record load when we navigate to another record / leave.
let recordLoadAbort = null
function abortRecordLoad() {
  if (recordLoadAbort) {
    recordLoadAbort.abort()
    recordLoadAbort = null
  }
}

async function loadRecord(recordID) {
  if (!page.value || !recordID || recordID === '0') return
  const moduleID = page.value.moduleID
  if (!moduleID) return
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return

  abortRecordLoad()
  const ac = new AbortController()
  recordLoadAbort = ac

  try {
    const loaded = await recordStore.findByID({
      namespaceID: mod.namespaceID,
      moduleID: mod.moduleID,
      recordID,
      force: true,
      signal: ac.signal,
    })
    pristineRecord.value = loaded
    record.value = loaded
    serverErrors.value = {}
  } catch (e) {
    if (ac.signal.aborted) return
    console.error('Failed to load record:', e)
    record.value = null
  } finally {
    if (recordLoadAbort === ac) recordLoadAbort = null
    if (!ac.signal.aborted) navigating.value = null
  }
}

async function loadPage() {
  const pageID = props.inModal ? props.modalPageID : route.params.pageID
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (!pageID) return

  abortRecordLoad()
  const ac = new AbortController()
  recordLoadAbort = ac

  loading.value = true
  record.value = null
  pristineRecord.value = null
  invisibleBlockIDs.value = new Set()

  try {
    page.value = pageStore.getByID(pageID) || null

    if (page.value) {
      if (!page.value.isRecordPage) {
        router.replace({
          name: 'page',
          params: {
            slug: route.params.slug,
            pageID: page.value.pageID,
          },
        })
        return
      }

      const layouts = pageLayoutStore.getByPageID(pageID)
      const vars = buildExpressionVariables({ isRecordPage: true, mode: mode.value })
      layout.value = await determineLayout(layouts, vars)

      const moduleID = page.value.moduleID
      if (moduleID) {
        const mod = moduleStore.getByID(moduleID)
        if (mod) {
          // A record needs a module WITH fields (compose.Record throws otherwise);
          // a fieldless module renders the page without an initialized record.
          if (mode.value === 'create' && mod.fields?.length) {
            // Handle clone
            if (route.query.cloneFromID) {
              try {
                const source = await recordStore.findByID({
                  namespaceID: mod.namespaceID,
                  moduleID: mod.moduleID,
                  recordID: route.query.cloneFromID,
                  signal: ac.signal,
                })
                const newRec = new compose.Record(mod)
                for (const field of mod.fields) {
                  newRec.setValue(field.name, source.values[field.name])
                }
                // Prefill ownedBy with current user
                newRec.ownedBy = $auth?.user?.userID || undefined
                record.value = newRec
              } catch (e) {
                if (ac.signal.aborted) return
                console.error('Failed to load source record for clone:', e)
                record.value = new compose.Record(mod, { ownedBy: $auth?.user?.userID })
              }
            } else {
              record.value = new compose.Record(mod, { ownedBy: $auth?.user?.userID })
            }

            const refField = route.query.refField
            const refValue = route.query.refValue
            if (record.value && refField && refValue) {
              const field = mod.fields.find(f => f.name === refField)
              if (field) {
                const valueArray = Array.isArray(refValue) ? refValue : [refValue]
                if (field.isMulti) {
                  record.value.setValue(refField, valueArray)
                } else {
                  record.value.setValue(refField, valueArray[0])
                }
              }
            }
          } else if (recordID && recordID !== '0') {
            try {
              const loaded = await recordStore.findByID({
                namespaceID: mod.namespaceID,
                moduleID: mod.moduleID,
                recordID,
                force: true,
                signal: ac.signal,
              })
              pristineRecord.value = loaded
              // In edit mode, start with a clone; mode watch will re-clone on later transitions
              record.value = mode.value === 'edit' ? loaded.clone() : loaded
            } catch (e) {
              if (ac.signal.aborted) return
              console.error('Failed to load record:', e)
              record.value = null
            }
          }
        }
      }
    }

    // Evaluate block visibility before revealing content (no flash)
    if (page.value?.blocks?.length) {
      const vars = buildExpressionVariables({
        record: record.value,
        isRecordPage: true,
        mode: mode.value,
      })
      invisibleBlockIDs.value = await evaluateBlocks(page.value.blocks, vars)
    }
  } finally {
    // If this load was superseded/cancelled, leave state to the newer load.
    if (!ac.signal.aborted) {
      loading.value = false
      if (recordLoadAbort === ac) recordLoadAbort = null
    }
  }
}

let _blockVisibilityTimer = null
async function evaluateBlockVisibility() {
  if (!page.value?.blocks?.length) return
  // Debounce rapid changes (e.g. mode switch, record swap)
  clearTimeout(_blockVisibilityTimer)
  _blockVisibilityTimer = setTimeout(async () => {
    const vars = buildExpressionVariables({
      record: record.value,
      isRecordPage: true,
      mode: mode.value,
    })
    invisibleBlockIDs.value = await evaluateBlocks(page.value.blocks, vars)
  }, 300)
}

// Re-evaluate when mode switches (view ↔ edit) or record reference changes
watch([record, mode], evaluateBlockVisibility, { deep: false })

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
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }
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

    if (props.inModal) {
      if (isNew.value) {
        // Update query to new recordID instead of '0'
        router.replace({
          query: {
            ...route.query,
            recordID: saved.recordID,
            edit: undefined,
          },
        })
      } else {
        router.replace({ query: { ...route.query, edit: undefined } })
      }
    } else {
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
    }
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
      $toast.toastErrorHandler(
        t(isNew.value ? 'notification.record.createFailed' : 'notification.record.updateFailed'),
      )(e)
    }
  } finally {
    isSaving.value = false
  }
}

function handleEdit() {
  router.push({ query: { ...route.query, edit: '1' } })
}

function handleClone() {
  if (props.inModal) {
    router.push({
      query: {
        ...route.query,
        recordID: '0',
        cloneFromID: record.value.recordID,
      },
    })
  } else {
    router.push({
      name: 'page.record',
      params: {
        slug: route.params.slug,
        pageID: props.inModal ? props.modalPageID : route.params.pageID,
        recordID: '0',
      },
      query: { cloneFromID: record.value.recordID },
    })
  }
}

function handleNew() {
  if (props.inModal) {
    router.push({
      query: {
        ...route.query,
        recordID: '0',
      },
    })
  } else {
    router.push({
      name: 'page.record',
      params: { slug: route.params.slug, pageID: route.params.pageID, recordID: '0' },
    })
  }
}

function handleCancel() {
  if (props.inModal) {
    if (mode.value === 'view') {
      emit('close')
    } else {
      const q = { ...route.query }
      delete q.edit
      if (isNew.value) {
        emit('close')
      } else {
        router.replace({ query: q })
      }
    }
    return
  }

  if (mode.value === 'edit') {
    // Cancel edit — return to view mode by removing the edit query param
    router.replace({ query: {} })
  } else {
    goBack()
  }
}

function goBack() {
  // If we have history, go back. Otherwise navigate to the namespace pages.
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push({
      name: 'pages',
      params: { slug: route.params.slug },
    })
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
    if (props.inModal) {
      emit('close')
    } else {
      goBack()
    }
  } catch (e) {
    console.error('Failed to delete record:', e)
    $toast.toastErrorHandler(t('notification.record.deleteFailed'))(e)
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
  if (page.value?.isRecordPage) {
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
    return window.confirm(t('general.editor.unsavedChanges'))
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
  () => [
    props.inModal ? props.modalPageID : route.params.pageID,
    props.inModal ? props.modalRecordID : route.params.recordID,
    route.query.cloneFromID,
    route.query.refField,
    route.query.refValue,
  ],
  ([newPageID, newRecordID, newCloneFromID, newRefField, newRefValue], old) => {
    const [oldPageID, , oldCloneFromID, oldRefField, oldRefValue] = old || []
    // If only the recordID changed (same page, no clone transition, no change in prefill params), just swap the record
    if (
      old &&
      newPageID === oldPageID &&
      newCloneFromID === oldCloneFromID &&
      newRefField === oldRefField &&
      newRefValue === oldRefValue &&
      page.value
    ) {
      loadRecord(newRecordID)
    } else {
      loadPage()
    }
  },
  { immediate: true },
)

const offRefetch = $eventBus?.on('refetch-records', () => loadPage())

onBeforeUnmount(() => {
  abortRecordLoad()
  offRefetch?.()
})
</script>
