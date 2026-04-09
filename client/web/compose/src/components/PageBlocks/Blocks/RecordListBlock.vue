<template>
  <PageBlock :block="block" @refreshBlock="fetchRecords(true)">
    <!-- Toolbar -->
    <div v-if="recordListModule" class="flex items-center gap-2 p-3 border-b">
      <!-- Add Record button (inline mode: prepend new row; otherwise: navigate) -->
      <Button
        v-if="!options.hideAddButton && recordListModule?.canCreateRecord && (options.editable || recordPageID)"
        :label="$t('block.recordList.addRecord')"
        icon="pi pi-plus"
        size="small"
        @click="options.editable ? addInlineRecord() : handleAddRecord()"
      />

      <!-- Filter presets dropdown -->
      <Menu
        v-if="visiblePresets.length"
        ref="presetMenuRef"
        :model="presetMenuItems"
        :popup="true"
      />
      <Button
        v-if="visiblePresets.length"
        :label="$t('block.recordList.filter.filters.label')"
        icon="pi pi-sliders-h"
        severity="secondary"
        outlined
        size="small"
        @click="$refs.presetMenuRef?.toggle($event)"
      />

      <!-- Spacer -->
      <div class="flex-1" />

      <!-- Filter button + Search -->
      <div class="flex items-center gap-1 flex-1 max-w-xl">
        <RecordListFilter
          v-if="showFilterButton"
          :module="recordListModule"
          :namespace="namespace"
          :model-value="recordListFilter"
          :allow-preset-save="!!options.customFilterPresets"
          @update:model-value="onFilterChange"
          @reset="onFilterReset"
          @save-preset="onSaveFilterPreset"
        />
        <CInputSearch
          v-if="!options.hideSearch"
          v-model="searchQuery"
          :placeholder="$t('general.label.search')"
          size="small"
          class="flex-1"
        />
      </div>
    </div>

    <!-- Active filters bar -->
    <div
      v-if="activeFilterDisplay.length"
      class="flex items-center flex-wrap gap-2 px-3 py-2 border-b"
    >
      <template v-for="(segment, si) in groupedActiveFilters" :key="si">
        <div class="flex items-center flex-wrap gap-1 border border-surface rounded-border p-1">
          <template v-for="(fg, fgi) in segment.groups" :key="fg.originalIndex">
            <Chip
              v-for="(f, fi) in fg.filter"
              :key="fi"
              removable
              class="text-sm"
              @remove="removeFilter(fg.originalIndex, fi)"
            >
              <span class="font-medium">{{ f.label || f.name }}</span>
              <span class="mx-1 text-muted-color">{{ getOperatorLabel(f.operator) }}</span>
              <span v-if="f.value != null" class="font-semibold text-primary">{{ formatFilterValue(f) }}</span>
              <span v-else class="text-muted-color italic">{{ $t('block.recordList.filter.nil') }}</span>
            </Chip>
            <span
              v-if="fgi < segment.groups.length - 1"
              class="text-xs text-muted-color uppercase font-medium"
            >
              {{ $t('block.recordList.filter.conditions.and') }}
            </span>
          </template>
        </div>
        <span
          v-if="segment.connector"
          class="text-xs text-muted-color uppercase font-medium"
        >
          {{ $t('block.recordList.filter.conditions.or') }}
        </span>
      </template>

      <Button
        :label="$t('block.recordList.filter.reset')"
        severity="secondary"
        text
        size="small"
        class="ml-auto"
        @click="onFilterReset"
      />
    </div>

    <!-- Selection / dirty bar -->
    <div
      v-if="selectedRecords.length || showBulkSave"
      class="flex items-center gap-2 px-3 py-2 bg-highlight border-b"
    >
      <span v-if="selectedRecords.length" class="text-sm font-medium">
        {{
          $t('block.recordList.selected', {
            count: selectedRecords.length,
            total: totalRecords,
          })
        }}
      </span>

      <div class="flex-1" />

      <!-- Save / discard all dirty (or just selected dirty rows) -->
      <template v-if="showBulkSave">
        <Button
          v-tooltip.bottom="$t('block.recordList.tooltip.saveChanges')"
          icon="pi pi-check"
          text
          size="small"
          severity="success"
          :loading="processingDirtyRecords === 'save'"
          :disabled="!!processingDirtyRecords"
          @click="handleSaveDirtyRecords"
        />
        <Button
          v-tooltip.bottom="$t('block.recordList.tooltip.discardChanges')"
          icon="pi pi-times"
          text
          size="small"
          severity="secondary"
          :loading="processingDirtyRecords === 'deny'"
          :disabled="!!processingDirtyRecords"
          @click="handleDenyDirtyRecords"
        />
      </template>

      <Button
        v-if="canDeleteSelected"
        :label="$t('block.recordList.tooltip.deleteSelected')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        @click="deleteSelected"
      />
    </div>

    <!-- No module configured -->
    <div v-if="!recordListModule" class="flex items-center justify-center flex-1 p-4">
      <span class="text-muted italic">
        {{ $t('block.recordList.noModule') }}
      </span>
    </div>

    <!-- Data table -->
    <div v-else class="flex-1 overflow-auto">
      <DataTable
        v-model:selection="selectedRecords"
        :value="records"
        :loading="loading"
        :rows="currentPerPage"
        :total-records="totalRecords"
        :lazy="true"
        :row-class="rowClass"
        scrollable
        scroll-height="flex"
        row-hover
        resizable-columns
        column-resize-mode="expand"
        data-key="recordID"
        :sort-field="sortField"
        :sort-order="sortOrder"
        @sort="onSort"
        @row-click="onRowClick"
        class="record-list-table"
        :pt="{
          headerCell: { class: 'bg-highlight-emphasis' },
        }"
      >
        <template #empty>
          <div class="flex items-center justify-center p-4 text-muted">
            {{ $t('block.recordList.noRecords') }}
          </div>
        </template>

        <!-- Selection checkbox column -->
        <Column
          v-if="canSelectRecords"
          selection-mode="multiple"
          header-style="width: 3rem"
          frozen
        />

        <!-- Data columns -->
        <Column
          v-for="col in columns"
          :key="col.name"
          :field="col.name"
          :header="col.label"
          :sortable="!options.hideSorting && !options.editable && !col.isMulti"
        >
          <template #body="{ data, index }">
            <CFieldEditor
              v-if="shouldShowEditor(data, col)"
              :field="col"
              :namespace="namespace"
              :model-value="data.values[col.name]"
              style="min-width: 200px;"
              @update:model-value="data.values[col.name] = $event; onInlineFieldChange(data)"
              @click.stop
            />
            <div v-else class="group flex items-start gap-1 min-w-0">
              <CFieldViewer
                :field="col"
                :record="data"
                :namespace="namespace"
              />
              <div
                v-if="showInlineActions(col, data)"
                class="flex items-center opacity-0 group-hover:opacity-100 transition-opacity shrink-0"
              >
                <Button
                  v-if="options.inlineRecordEditEnabled && canInlineEdit(data, col)"
                  v-tooltip.top="$t('block.recordList.record.tooltip.edit')"
                  icon="pi pi-pencil"
                  text
                  rounded
                  size="small"
                  severity="secondary"
                  class="w-6 h-6 p-0 mt-0.5"
                  @click.stop="startInlineEdit(data, col)"
                />
                <Button
                  v-if="showInlineFilter(col)"
                  v-tooltip.top="$t('block.recordList.filterByValue')"
                  icon="pi pi-filter"
                  text
                  rounded
                  size="small"
                  severity="secondary"
                  class="w-6 h-6 p-0 mt-0.5"
                  @click.stop="filterByValue(data, col)"
                />
              </div>
            </div>
          </template>
        </Column>

        <!-- Actions column: save/discard when dirty, ellipsis menu otherwise -->
        <Column
          v-if="options.editable || options.inlineRecordEditEnabled || hasRowActions"
          header-style="width: 5rem"
          :pt="{
            headerCell: { class: 'border-l-0' },
            bodyCell: { class: 'px-2 py-1 border-l-0' },
          }"
          frozen
          align-frozen="right"
        >
          <template #body="{ data, index }">
            <div class="flex items-center justify-end gap-1">
              <!-- Inline save/discard (only when row is dirty) -->
              <template v-if="(options.editable || options.inlineRecordEditEnabled) && showSaveAction(data)">
                <Button
                  v-tooltip.top="$t('block.recordList.tooltip.saveChanges')"
                  icon="pi pi-check"
                  text
                  size="small"
                  severity="success"
                  :loading="processingRecords[getRecordKey(data)] === 'save'"
                  :disabled="!!processingRecords[getRecordKey(data)]"
                  @click.stop="handleSaveInline(data, index)"
                />
                <Button
                  v-tooltip.top="$t('block.recordList.tooltip.discardChanges')"
                  icon="pi pi-times"
                  text
                  size="small"
                  severity="secondary"
                  :loading="processingRecords[getRecordKey(data)] === 'deny'"
                  :disabled="!!processingRecords[getRecordKey(data)]"
                  @click.stop="handleDenyInline(data, index)"
                />

                <Divider v-if="hasRowActions" layout="vertical" class="mx-1 h-5" />
              </template>

              <!-- Row action menu button (always visible for dirty rows) -->
              <Button
                v-if="hasRowActions"
                icon="pi pi-ellipsis-v"
                text
                size="small"
                severity="secondary"
                :class="showSaveAction(data) ? '' : 'row-action-btn'"
                @click.stop="openRowMenu($event, data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <template v-if="recordListModule && !options.hidePaging && totalRecords > 0" #footer>
      <div class="flex items-center flex-wrap gap-2 px-3 py-2">
        <div class="flex items-center text-sm">
          <span v-if="options.showTotalCount !== false" class="whitespace-nowrap text-muted">
            {{ paginationRangeText }}
          </span>

          <Divider v-if="options.showRecordPerPageOption && options.showTotalCount !== false" layout="vertical" />

          <!-- Page size selector -->
          <div v-if="options.showRecordPerPageOption" class="flex items-center gap-2 whitespace-nowrap">
            <span class="text-muted">
              {{ $t('block.recordList.pagination.recordsPerPage') }}
            </span>
            <Select
              v-model="currentPerPage"
              :options="pageSizeOptions"
              size="small"
              class="w-20"
              @change="onPageSizeChange"
            />
          </div>
        </div>

        <!-- Prev/Next -->
        <div class="flex items-center ml-auto gap-1">
          <Button
            icon="pi pi-chevron-left"
            text
            severity="secondary"
            size="small"
            :disabled="!hasPrevPage"
            @click="goToPrevPage"
          />
          <Button
            icon="pi pi-chevron-right"
            text
            severity="secondary"
            size="small"
            :disabled="!hasNextPage"
            @click="goToNextPage"
          />
        </div>
      </div>
    </template>

    <!-- Row action menu (teleported) -->
    <Menu :key="rowMenuKey" ref="rowMenuRef" :model="rowMenuItems" :popup="true">
      <template #item="{ item, props }">
        <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
          <a v-ripple :href="href" v-bind="props.action" @click="navigate">
            <span :class="item.icon" />
            <span class="ml-2">{{ item.label }}</span>
          </a>
        </router-link>
        <a v-else v-ripple v-bind="props.action" :class="item.class">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </template>
    </Menu>
  </PageBlock>
</template>

<script setup>
import axios from 'axios'
import { computed, inject, nextTick, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete, usePermissions } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer, CFieldEditor, CInputSearch } = components
import { useModuleStore } from '@/stores/module'
import { useRecordStore } from '@/stores/record'
import { useReminderStore } from '@/stores/reminder'
import PageBlock from './PageBlock.vue'
import { usePageStore } from '@/stores/page'
import RecordListFilter from '../../Common/RecordListFilter.vue'
import { evaluatePrefilter, queryToFilter, getFieldFilter, convertRecordListFilter, formatActiveFilterOperator, isBetweenOperator } from '../../../lib/record-filter'

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  page: {
    type: Object,
    default: () => ({}),
  },
  record: {
    type: Object,
    default: undefined,
  },
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const $ComposeAPI = inject('$ComposeAPI')
const $auth = inject('$auth', {})
const $eventBus = inject('$eventBus', null)
const $toast = inject('$toast')
const moduleStore = useModuleStore()
const recordStore = useRecordStore()
const pageStore = usePageStore()
const reminderStore = useReminderStore()

// Injectable route resolver — allows admin views to override routes
// Default: public page routes using recordPageID
const $recordRoutes = inject('$recordRoutes', null)

const loading = ref(false)
const records = ref([])
const totalRecords = ref(0)
const searchQuery = ref('')
const sortField = ref(null)
const sortOrder = ref(null)

// Record list filter state
const recordListFilter = ref([])
const presetMenuRef = ref(null)
const selectedRecords = ref([])

// Pagination state — cursor based
const pageCursors = ref([]) // stack of previous page cursors
const nextPageCursor = ref(null)
const currentPageIndex = ref(0) // 0-based page index for display

// Page size — always includes the configured perPage, deduped and sorted
const pageSizeOptions = computed(() => {
  const configured = props.block.options?.perPage
  const base = [10, 20, 50, 100]
  const set = configured ? [...new Set([...base, configured])] : base
  return set.sort((a, b) => a - b)
})
const currentPerPage = ref(props.block.options?.perPage || 20)

// Row action menu
const rowMenuRef = ref(null)
const rowMenuKey = ref(0)
const activeRowRecord = ref(null)

// Inline editing — Corteza-style dirty tracking
const dirtyRecords = reactive({})
const activeInlineEdits = ref(new Map())
const processingRecords = reactive({})
const processingDirtyRecords = ref('')

// Counter for temp IDs on newly added inline rows
let _tempIdCounter = 0
function getRecordKey(record) {
  return record.recordID || record._tempID || ''
}

const options = computed(() => props.block.options || {})

// Sync perPage from block options (e.g. configurator changes)
watch(
  () => options.value.perPage,
  val => {
    if (val) currentPerPage.value = val
  },
)

// Resolve the module for this block
const recordListModule = computed(() => {
  const moduleID = options.value.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Find the record page for this module (page with matching moduleID)
const recordPageID = computed(() => {
  // When using custom routes, we don't need a public record page
  if ($recordRoutes) return 'admin'

  const moduleID = recordListModule.value?.moduleID
  if (!moduleID) return null

  const pages = pageStore.set
  const page = (pages || []).find(p => p.moduleID === moduleID)
  return page?.pageID || null
})

// Columns derived from block field config or module fields
const columns = computed(() => {
  if (!recordListModule.value) return []

  const configuredFields = options.value.fields || []

  if (configuredFields.length === 0) {
    // Fallback: show first 5 fields
    return (recordListModule.value.fields || []).slice(0, 5)
  }

  // Map configured field names to module field definitions
  if (recordListModule.value.filterFields) {
    return recordListModule.value.filterFields(configuredFields)
  }

  return configuredFields
})

// Pagination helpers
const hasPrevPage = computed(() => pageCursors.value.length > 0)
const hasNextPage = computed(() => !!nextPageCursor.value)

// "Showing X–Y of N" text
const paginationRangeText = computed(() => {
  const start = currentPageIndex.value * currentPerPage.value + 1
  const end = Math.min(start + currentPerPage.value - 1, totalRecords.value)
  return `${start} - ${end} of ${totalRecords.value} records`
})

// Permission-based helpers
const canSelectRecords = computed(() => options.value.selectable !== false)
const canDeleteSelected = computed(() => selectedRecords.value.some(r => r.canDeleteRecord))
const hasRowActions = computed(() => {
  if (!recordPageID.value && !recordListModule.value?.canCreateRecord) return false
  const o = options.value
  return !o.hideRecordViewButton
    || !o.hideRecordEditButton
    || !o.hideRecordCloneButton
    || !o.hideRecordReminderButton
    || !o.hideRecordPermissionsButton
    || !o.hideRecordDeleteButton
})

// Row action menu items — filtered by per-record permissions and configurator hide options
const rowMenuItems = computed(() => {
  const record = activeRowRecord.value
  const items = []

  if (recordPageID.value && record?.canReadRecord && !options.value.hideRecordViewButton) {
    items.push({
      label: t('block.recordList.record.tooltip.view'),
      icon: 'pi pi-eye',
      route: $recordRoutes
        ? $recordRoutes.view(recordListModule.value.moduleID, record.recordID)
        : {
            name: 'page.record',
            params: {
              pageID: recordPageID.value,
              recordID: record.recordID,
            },
          },
    })
  }

  if (recordPageID.value && record?.canUpdateRecord && !options.value.hideRecordEditButton) {
    items.push({
      label: t('block.recordList.record.tooltip.edit'),
      icon: 'pi pi-pencil',
      route: $recordRoutes
        ? $recordRoutes.edit
          ? $recordRoutes.edit(recordListModule.value.moduleID, record.recordID)
          : $recordRoutes.view(recordListModule.value.moduleID, record.recordID)
        : {
            name: 'page.record',
            params: {
              pageID: recordPageID.value,
              recordID: record.recordID,
            },
            query: { edit: '1' },
          },
    })
  }

  if (recordPageID.value && recordListModule.value?.canCreateRecord && !options.value.hideRecordCloneButton) {
    items.push({
      label: t('block.recordList.record.tooltip.clone'),
      icon: 'pi pi-copy',
      command: () => handleCloneRecord(record),
    })
  }

  if (record?.recordID && !options.value.hideRecordReminderButton) {
    items.push({
      label: t('reminder.add'),
      icon: 'pi pi-clock',
      command: () => createReminder(record),
    })
  }

  if (record?.recordID && (record.canGrant ?? recordListModule.value?.canGrant) && !options.value.hideRecordPermissionsButton) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        openPermissions({
          resource: `corteza::compose:record/${props.namespace.namespaceID}/${recordListModule.value.moduleID}/${record.recordID}`,
          title: record.recordID,
        })
      },
    })
  }

  if (record?.canDeleteRecord && !options.value.hideRecordDeleteButton) {
    if (items.length > 0) {
      items.push({ separator: true })
    }

    items.push({
      label: t('block.recordList.record.tooltip.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => confirmDeleteRecord(record),
    })
  }

  return items
})

// Inline editing helpers
const dirtyRecordsCount = computed(() => Object.keys(dirtyRecords).length)

const showBulkSave = computed(() => {
  if (!options.value.editable) return false
  if (selectedRecords.value.length > 0) {
    const toSaveCount = selectedRecords.value.filter(r => showSaveAction(r)).length
    return toSaveCount > 0
  }
  return dirtyRecordsCount.value > 1
})

function rowClass(data) {
  return dirtyRecords[getRecordKey(data)] ? 'bg-yellow-50 dark:bg-yellow-950' : ''
}

function isInlineEditField(col) {
  const editFields = options.value.editFields || []
  if (!editFields.length) return true
  return editFields.some(f => (f.name ?? f) === col.name)
}

function shouldShowEditor(data, col) {
  if (options.value.editable && isInlineEditField(col)) {
    return !data.recordID || data.recordID === '0' || data.canUpdateRecord !== false
  }
  const key = getRecordKey(data)
  const fields = activeInlineEdits.value.get(key)
  return fields && fields.has(col.name)
}

function canInlineEdit(data, col) {
  if (data.canUpdateRecord === false) return false
  if (data.deletedAt) return false
  return isInlineEditField(col)
}

function startInlineEdit(data, col) {
  const key = getRecordKey(data)
  let fields = activeInlineEdits.value.get(key)
  if (!fields) {
    fields = new Set()
    activeInlineEdits.value.set(key, fields)
  }
  fields.add(col.name)
}

function clearInlineEdits(record) {
  const key = getRecordKey(record)
  activeInlineEdits.value.delete(key)
}

function onInlineFieldChange(record) {
  const key = getRecordKey(record)
  if (!dirtyRecords[key]) {
    dirtyRecords[key] = true
  }
}

function showSaveAction(record) {
  if (!recordListModule.value) return false
  const key = getRecordKey(record)
  // New inline record (no recordID) always shows save action if module allows create
  if (!record.recordID || record.recordID === '0') return !!recordListModule.value.canCreateRecord
  if (!dirtyRecords[key]) return false
  return record.canUpdateRecord !== false
}

function addInlineRecord() {
  const r = new compose.Record(recordListModule.value, {})
  // Assign a temp ID so multiple new rows can be tracked independently
  r._tempID = `_new_${++_tempIdCounter}`

  // Prefill refField if present (links to parent record)
  if (options.value.refField && props.record) {
    const refField = recordListModule.value.fields.find(f => f.name === options.value.refField)
    if (refField?.isMulti) {
      r.values[options.value.refField] = [props.record.recordID]
    } else {
      r.values[options.value.refField] = props.record.recordID
    }
  }

  records.value.unshift(r)
  // Mark new row as dirty immediately so save/discard controls appear
  dirtyRecords[r._tempID] = true
}

async function handleSaveInline(record, index) {
  if (!recordListModule.value) return
  const key = getRecordKey(record)
  processingRecords[key] = 'save'
  try {
    const isNew = !record.recordID || record.recordID === '0'
    const saved = isNew ? await recordStore.create(record) : await recordStore.update(record)
    delete dirtyRecords[key]
    clearInlineEdits(record)
    const newRecord = new compose.Record(recordListModule.value, saved)
    records.value.splice(index, 1, newRecord)
  } catch (e) {
    console.error('Failed to save inline record:', e)
  } finally {
    delete processingRecords[key]
  }
}

async function handleDenyInline(record, index) {
  if (!recordListModule.value) return
  const key = getRecordKey(record)
  processingRecords[key] = 'deny'
  delete dirtyRecords[key]
  clearInlineEdits(record)
  const isNew = !record.recordID || record.recordID === '0'
  if (isNew) {
    records.value.splice(index, 1)
    delete processingRecords[key]
    return
  }
  try {
    const fresh = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: recordListModule.value.moduleID,
      recordID: record.recordID,
      force: true,
    })
    records.value.splice(index, 1, fresh)
  } catch (e) {
    console.error('Failed to revert inline record:', e)
  } finally {
    delete processingRecords[key]
  }
}

async function handleSaveDirtyRecords() {
  if (!recordListModule.value) return
  // Collect: if rows selected, only those; otherwise all dirty
  const toSave = records.value.filter((r, idx) => {
    if (!showSaveAction(r)) return false
    if (selectedRecords.value.length > 0) {
      return selectedRecords.value.some(s => getRecordKey(s) === getRecordKey(r))
    }
    return true
  })
  if (!toSave.length) return

  processingDirtyRecords.value = 'save'
  let hasError = false

  for (const record of toSave) {
    const key = getRecordKey(record)
    const index = records.value.findIndex(r => getRecordKey(r) === key)
    if (index === -1) continue
    const isNew = !record.recordID || record.recordID === '0'
    try {
      const saved = isNew ? await recordStore.create(record) : await recordStore.update(record)
      delete dirtyRecords[key]
      clearInlineEdits(record)
      const newRecord = new compose.Record(recordListModule.value, saved)
      records.value.splice(index, 1, newRecord)
    } catch (e) {
      hasError = true
      console.error('Failed to save record:', e)
    }
  }

  processingDirtyRecords.value = ''
  if (!hasError) selectedRecords.value = []
}

async function handleDenyDirtyRecords() {
  if (!recordListModule.value) return
  const toDeny = records.value.filter(r => {
    if (!showSaveAction(r)) return false
    if (selectedRecords.value.length > 0) {
      return selectedRecords.value.some(s => getRecordKey(s) === getRecordKey(r))
    }
    return true
  })
  if (!toDeny.length) return

  processingDirtyRecords.value = 'deny'

  for (const record of [...toDeny]) {
    const key = getRecordKey(record)
    const index = records.value.findIndex(r => getRecordKey(r) === key)
    if (index === -1) continue
    delete dirtyRecords[key]
    clearInlineEdits(record)
    const isNew = !record.recordID || record.recordID === '0'
    if (isNew) {
      records.value.splice(index, 1)
      continue
    }
    try {
      const fresh = await recordStore.findByID({
        namespaceID: props.namespace.namespaceID,
        moduleID: recordListModule.value.moduleID,
        recordID: record.recordID,
        force: true,
      })
      records.value.splice(index, 1, fresh)
    } catch (e) {
      console.error('Failed to revert record:', e)
    }
  }

  processingDirtyRecords.value = ''
  selectedRecords.value = []
}

// Request cancellation
let cancelPendingRequest = null

function abortPendingRequest() {
  if (cancelPendingRequest) {
    cancelPendingRequest()
    cancelPendingRequest = null
  }
}

/**
 * Build the evaluated prefilter string from block options,
 * matching Corteza's prepRecordList() logic.
 *
 * - Evaluates template expressions (${recordID}, ${ownerID}, ${userID})
 * - Handles refField for parent record linking
 */
function buildPrefilter() {
  const { prefilter, refField } = options.value
  const mod = recordListModule.value
  const filterParts = []

  if (prefilter) {
    const record = props.record
    const user = $auth?.user || {}
    const pf = evaluatePrefilter(prefilter, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
    filterParts.push(`(${pf})`)
  }

  // Handle refField — links to parent record
  if (refField && props.record) {
    const refFieldObj = mod?.fields?.find(f => f.name === refField)
    if (refFieldObj && refFieldObj.isMulti) {
      filterParts.push(getFieldFilter(refField, 'Record', props.record.recordID, 'IN'))
    } else {
      filterParts.push(getFieldFilter(refField, 'Record', props.record.recordID, '='))
    }
  }

  if (drillDownPrefilter.value) {
    filterParts.push(`(${drillDownPrefilter.value})`)
  }

  return filterParts.filter(Boolean).join(' AND ')
}

// Fetch records with cancellation support
async function fetchRecords(resetCursor = false) {
  if (!recordListModule.value) return

  abortPendingRequest()
  loading.value = true

  if (resetCursor) {
    pageCursors.value = []
    nextPageCursor.value = null
    currentPageIndex.value = 0
    selectedRecords.value = []
    // Clear inline editing state on full refresh
    Object.keys(dirtyRecords).forEach(k => delete dirtyRecords[k])
    activeInlineEdits.value.clear()
    Object.keys(processingRecords).forEach(k => delete processingRecords[k])
  }

  try {
    // Build sort expression
    let sort = options.value.presort || ''
    if (sortField.value) {
      sort = `${sortField.value} ${sortOrder.value === 1 ? 'ASC' : 'DESC'}`
    }

    // Build query using queryToFilter — matches Corteza's RecordListBase behavior
    const evaluatedPrefilter = buildPrefilter()

    // Determine search fields: use configured searchableFields or fall back to visible columns
    let searchFields = []
    const searchableFieldConfig = options.value.searchableFields || []
    if (searchableFieldConfig.length > 0 && recordListModule.value.filterFields) {
      searchFields = recordListModule.value.filterFields(searchableFieldConfig)
    } else {
      // Default to visible columns (same as old Corteza behavior)
      searchFields = columns.value
    }

    // Prepare filter groups exactly like old Corteza RecordListBase
    const filterGroups = recordListFilter.value.map(g => {
      const filter = convertRecordListFilter(g.filter || [])
      return { ...g, filter }
    }).filter(g => g.filter?.length)

    const query = queryToFilter(
      searchQuery.value || '',
      evaluatedPrefilter,
      searchFields,
      filterGroups,
    )

    // Determine page cursor for API call
    let pageCursor
    if (resetCursor) {
      pageCursor = undefined
    } else if (pageCursors.value.length > 0) {
      pageCursor = pageCursors.value[pageCursors.value.length - 1]
    } else {
      pageCursor = nextPageCursor.value
    }

    const { response, cancel } = $ComposeAPI.recordListCancellable({
      namespaceID: props.namespace.namespaceID,
      moduleID: recordListModule.value.moduleID,
      query,
      sort,
      limit: currentPerPage.value,
      pageCursor,
      // incTotal is only supported on the first page (no cursor); sending it with a
      // cursor causes a server error. Keep the total from the first page on subsequent pages.
      incTotal: !pageCursor,
    })

    cancelPendingRequest = cancel
    const result = await response()
    cancelPendingRequest = null

    const mod = recordListModule.value
    const set = (result.set || []).map(r => new compose.Record(mod, r))

    records.value = set
    if (result.filter?.total !== undefined) {
      totalRecords.value = result.filter.total
    } else if (!pageCursor) {
      totalRecords.value = set.length
    }

    // Store next page cursor
    nextPageCursor.value = result.filter?.nextPage || null
    loading.value = false
  } catch (error) {
    if (!axios.isCancel(error)) {
      console.error('Failed to fetch records for RecordList block:', error)
      $toast?.toastDanger(error?.message || t('notification.record.loadFailed'))
      records.value = []
      totalRecords.value = 0
      loading.value = false
    }
  }
}

// Debounced search watcher
let searchDebounceTimer = null
watch(searchQuery, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    fetchRecords(true)
  }, 300)
})

// Drill down event listener
const drillDownPrefilter = ref(null)
const offDrillDown = $eventBus?.on(`drill-down-recordList:${props.block.blockID}`, payload => {
  drillDownPrefilter.value = payload?.prefilter
  fetchRecords(true)
})

// Cleanup on unmount
onBeforeUnmount(() => {
  abortPendingRequest()
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  offDrillDown?.()
})

function onSort(event) {
  sortField.value = event.sortField
  sortOrder.value = event.sortOrder
  fetchRecords(true)
}

function onPageSizeChange() {
  fetchRecords(true)
}

function goToNextPage() {
  if (!nextPageCursor.value) return
  pageCursors.value.push(nextPageCursor.value)
  currentPageIndex.value++
  fetchRecords()
}

function goToPrevPage() {
  if (pageCursors.value.length === 0) return
  pageCursors.value.pop()
  currentPageIndex.value = Math.max(0, currentPageIndex.value - 1)
  fetchRecords()
}

function onRowClick(event) {
  const record = event.data
  if (!record?.recordID) return

  // Don't navigate in inline editing mode — fields are edited in-place
  if (options.value.editable) return

  // Don't navigate if clicking on a checkbox or action button
  const target = event.originalEvent?.target
  if (target?.closest('.p-checkbox') || target?.closest('button')) return

  if (!recordPageID.value) return

  const displayOption = options.value.recordDisplayOption || 'sameTab'
  if (displayOption === 'doNothing') return

  const isEditMode = options.value.openRecordInEditMode && record.canUpdateRecord

  if (displayOption === 'modal' && !$recordRoutes) {
    const query = {
      ...route.query,
      recordPageID: recordPageID.value,
      recordID: record.recordID,
    }
    if (isEditMode) {
      query.edit = '1'
    }
    router.push({ query })
    return
  }

  let routeObj
  if ($recordRoutes) {
    if (isEditMode && $recordRoutes.edit) {
      routeObj = $recordRoutes.edit(recordListModule.value.moduleID, record.recordID)
    } else {
      routeObj = $recordRoutes.view(recordListModule.value.moduleID, record.recordID)
      if (isEditMode) {
        routeObj.query = { ...(routeObj.query || {}), edit: '1' }
      }
    }
  } else {
    routeObj = {
      name: 'page.record',
      params: {
        pageID: recordPageID.value,
        recordID: record.recordID,
      },
    }
    if (isEditMode) {
      routeObj.query = { edit: '1' }
    }
  }

  if (displayOption === 'newTab') {
    const resolved = router.resolve(routeObj)
    window.open(resolved.href, '_blank')
  } else {
    router.push(routeObj)
  }
}

// Row actions
function openRowMenu(event, record) {
  activeRowRecord.value = record
  rowMenuKey.value++
  nextTick(() => {
    rowMenuRef.value.show(event)
  })
}

function handleAddRecord() {
  if (!recordPageID.value) return

  const displayOption = options.value.addRecordDisplayOption || 'sameTab'

  if (displayOption === 'modal' && !$recordRoutes) {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPageID.value,
        recordID: '0',
      }
    })
    return
  }

  const routeObj = $recordRoutes
      ? $recordRoutes.create(recordListModule.value.moduleID)
      : {
          name: 'page.record',
          params: {
            pageID: recordPageID.value,
            recordID: '0',
          },
        }

  if (displayOption === 'newTab') {
    const resolved = router.resolve(routeObj)
    window.open(resolved.href, '_blank')
  } else {
    router.push(routeObj)
  }
}

function handleCloneRecord(record) {
  if (!record || !recordPageID.value) return

  const displayOption = options.value.addRecordDisplayOption || 'sameTab'

  if (displayOption === 'modal' && !$recordRoutes) {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPageID.value,
        recordID: '0',
        cloneFromID: record.recordID,
      }
    })
    return
  }

  const routeObj = $recordRoutes
      ? $recordRoutes.create(recordListModule.value.moduleID, record.recordID)
      : {
          name: 'page.record',
          params: {
            pageID: recordPageID.value,
            recordID: '0',
          },
          query: { cloneFromID: record.recordID },
        }

  if (displayOption === 'newTab') {
    const resolved = router.resolve(routeObj)
    window.open(resolved.href, '_blank')
  } else {
    router.push(routeObj)
  }
}

function createReminder (record) {
  if (!record?.recordID) return

  const sourceField = (options.value.fields || []).find(({ name }) => {
    const value = record.values?.[name]
    return Array.isArray(value) ? value.length > 0 : !!value
  })

  const fieldValue = sourceField ? record.values?.[sourceField.name] : null
  const title = Array.isArray(fieldValue) ? fieldValue.join(', ') : (fieldValue || '')
  const payload = {
    title,
  }

  if (recordPageID.value) {
    payload.link = {
      name: 'page.record',
      label: t('reminder.recordPageLink'),
      params: {
        slug: props.namespace.slug || props.namespace.namespaceID,
        pageID: recordPageID.value,
        recordID: record.recordID,
      },
    }
  }

  reminderStore.startCreate({
    resource: `compose:record:${record.recordID}`,
    payload,
  })
}

function confirmDeleteRecord(record) {
  if (!record?.recordID) return

  confirmDelete({
    message: t('block.record.confirmDelete'),
    header: t('block.recordList.record.tooltip.delete'),
    onConfirm: async () => {
      try {
        await recordStore.delete({
          namespaceID: props.namespace.namespaceID,
          moduleID: recordListModule.value.moduleID,
          recordID: record.recordID,
        })
        fetchRecords(true)
      } catch (e) {
        console.error('Failed to delete record:', e)
      }
    },
  })
}

async function deleteSelected() {
  if (!selectedRecords.value.length) return
  try {
    for (const record of selectedRecords.value) {
      await recordStore.delete({
        namespaceID: props.namespace.namespaceID,
        moduleID: recordListModule.value.moduleID,
        recordID: record.recordID,
      })
    }
    selectedRecords.value = []
    fetchRecords(true)
  } catch (e) {
    console.error('Failed to delete selected records:', e)
  }
}

const offRefetch = $eventBus?.on('refetch-records', () => fetchRecords(true))

// Reload when module changes
watch(
  () => recordListModule.value?.moduleID,
  (newID) => {
    // Load persisted filters when module changes
    loadStoredFilter()
    fetchRecords(true)
  },
  { immediate: true },
)

// --- Filter logic ---

const showFilterButton = computed(() => {
  if (options.value.hideFiltering) return false
  if (!recordListModule.value) return false
  return filterableFields.value.length > 0
})

const filterableFields = computed(() => {
  if (!recordListModule.value) return []
  return (recordListModule.value.fields || []).filter(
    f => f.isFilterable && f.canReadRecordValue !== false,
  )
})

// Active filters display
const activeFilterDisplay = computed(() => {
  return (recordListFilter.value || []).filter(
    g => g.filter?.some(f => f.name),
  )
})

const groupedActiveFilters = computed(() => {
  const groups = activeFilterDisplay.value.map((g, idx) => ({
    ...g,
    originalIndex: idx,
  }))

  // Segment groups by OR/AND connectors for display
  const segments = []
  let current = { groups: [], connector: null }
  for (let i = 0; i < groups.length; i++) {
    const g = groups[i]
    if (i > 0 && g.groupCondition === 'OR') {
      current.connector = 'OR'
      segments.push(current)
      current = { groups: [], connector: null }
    }
    current.groups.push(g)
  }
  if (current.groups.length) segments.push(current)
  return segments
})

function getOperatorLabel(op) {
  const key = formatActiveFilterOperator(op)
  return t(`block.recordList.filter.operatorLabels.${key}`)
}

function formatFilterValue(f) {
  if (isBetweenOperator(f.operator) && f.value) {
    return `${f.value.start || '?'} – ${f.value.end || '?'}`
  }
  if (Array.isArray(f.value)) return f.value.join(', ')
  return String(f.value ?? '')
}

function removeFilter(groupIndex, filterIndex) {
  const updated = JSON.parse(JSON.stringify(recordListFilter.value))
  updated[groupIndex].filter.splice(filterIndex, 1)
  if (!updated[groupIndex].filter.length) {
    updated.splice(groupIndex, 1)
  }
  recordListFilter.value = updated
  persistFilter()
  fetchRecords(true)
}

function onFilterChange(newFilter) {
  recordListFilter.value = newFilter
  persistFilter()
  fetchRecords(true)
}

function onFilterReset() {
  recordListFilter.value = []
  persistFilter()
  fetchRecords(true)
}

// --- Inline value filtering ---

function showInlineActions(col, data) {
  return (options.value.inlineRecordEditEnabled && canInlineEdit(data, col)) || showInlineFilter(col)
}

function showInlineFilter(col) {
  if (options.value.hideFiltering) return false
  if (!options.value.inlineValueFiltering) return false
  return col.isFilterable !== false
}

function filterByValue(record, col) {
  const value = record.values?.[col.name]
  const newFilter = {
    name: col.name,
    kind: col.kind,
    operator: col.isMulti ? 'IN' : '=',
    value: Array.isArray(value) ? value[0] : value,
    condition: '',
  }

  // Add to existing filter or create a new group
  const updated = JSON.parse(JSON.stringify(recordListFilter.value))
  if (updated.length) {
    updated[0].filter.push(newFilter)
  } else {
    updated.push({ filter: [newFilter], groupCondition: 'OR' })
  }
  recordListFilter.value = updated
  persistFilter()
  fetchRecords(true)
}

// --- Filter persistence (localStorage) ---

function filterStorageKey() {
  return `recordListFilter-${props.block.blockID}`
}

function persistFilter() {
  try {
    const key = filterStorageKey()
    if (recordListFilter.value.length) {
      localStorage.setItem(key, JSON.stringify(recordListFilter.value))
    } else {
      localStorage.removeItem(key)
    }
  } catch (e) {
    // localStorage not available
  }
}

function loadStoredFilter() {
  try {
    const key = filterStorageKey()
    const stored = localStorage.getItem(key)
    if (stored) {
      recordListFilter.value = JSON.parse(stored)
    }
  } catch (e) {
    // localStorage not available or invalid data
  }
}

// --- Filter presets ---

const visiblePresets = computed(() => {
  const presets = options.value.filterPresets || []
  // Also merge user-saved presets from localStorage
  try {
    const stored = localStorage.getItem(`recordListFilterPresets-${props.block.blockID}`)
    if (stored) {
      const userPresets = JSON.parse(stored)
      return [...presets, ...userPresets].filter(p => p.name)
    }
  } catch (e) {
    // ignore
  }
  return presets.filter(p => p.name)
})

const presetMenuItems = computed(() => {
  return visiblePresets.value.map(preset => ({
    label: preset.name,
    command: () => applyPreset(preset),
  }))
})

function applyPreset(preset) {
  if (!preset.filter) return
  // Preset filter can be a structured array (user-saved) or empty
  if (Array.isArray(preset.filter)) {
    recordListFilter.value = preset.filter
  } else {
    recordListFilter.value = []
  }
  persistFilter()
  fetchRecords(true)
}

function onSaveFilterPreset(filterData) {
  const name = prompt(t('block.recordList.filterPresets.filterName'))
  if (!name) return

  try {
    const key = `recordListFilterPresets-${props.block.blockID}`
    const stored = localStorage.getItem(key)
    const presets = stored ? JSON.parse(stored) : []
    presets.push({ name, filter: filterData })
    localStorage.setItem(key, JSON.stringify(presets))
  } catch (e) {
    // localStorage not available
  }
}
</script>

<style scoped>
.record-list-table :deep(.p-datatable-tbody > tr) {
  cursor: pointer;
}

.record-list-table :deep(.p-datatable-tbody > tr > td) {
  vertical-align: top;
}

.record-list-table :deep(.row-action-btn) {
  opacity: 0;
  transition: opacity 0.15s ease;
}

.record-list-table :deep(.p-datatable-tbody > tr:hover .row-action-btn) {
  opacity: 1;
}

:deep(.p-datatable-gridlines .p-datatable-paginator-bottom) {
  border-width: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :first-child) {
  border-left: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :last-child) {
  border-right: 0;
}

:deep(.p-datatable-gridlines .p-datatable-thead > tr > th) {
  border-top: 0;
}

:deep(.p-datatable-mask) {
  background: color-mix(in srgb, var(--p-content-background) 80%, transparent) !important;
}

:deep(.p-datatable-mask .p-icon-spin) {
  color: var(--p-primary-color);
}
</style>
