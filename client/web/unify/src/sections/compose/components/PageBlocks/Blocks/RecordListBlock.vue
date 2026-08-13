<template>
  <PageBlock :block="block" @refreshBlock="fetchRecords(true)">
    <!-- Toolbar -->
    <div v-if="recordListModule" class="flex items-center gap-2 p-3 border-b">
      <!-- Add Record button (inline mode: prepend new row; otherwise: navigate) -->
      <span
        v-if="!options.hideAddButton && recordListModule?.canCreateRecord"
        v-tooltip.bottom="addRecordDisabled ? $t('block.noRecordPage') : ''"
      >
        <Button
          :label="$t('block.recordList.addRecord')"
          icon="pi pi-plus"
          size="small"
          :disabled="addRecordDisabled"
          @click="options.editable ? addInlineRecord() : handleAddRecord()"
        />
      </span>

      <!-- Import button -->
      <RecordImporter
        v-if="!options.hideImportButton && recordListModule?.canCreateRecord"
        :module="recordListModule"
        :namespace="namespace"
        @import-successful="fetchRecords(true)"
      />

      <!-- Export button -->
      <RecordExporter
        v-if="options.allowExport && recordListModule"
        :module="recordListModule"
        :namespace="namespace"
        :filter="currentQuery"
        :selection="selectedRecords"
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

      <!-- Show deleted records toggle -->
      <Button
        v-if="options.showDeletedRecordsOption"
        :label="
          showingDeletedRecords
            ? $t('block.recordList.showRecords.existing')
            : $t('block.recordList.showRecords.deleted')
        "
        :icon="showingDeletedRecords ? 'pi pi-eye' : 'pi pi-trash'"
        :severity="showingDeletedRecords ? 'warn' : 'secondary'"
        :outlined="!showingDeletedRecords"
        size="small"
        @click="toggleDeletedRecords"
      />

      <!-- Configure Fields button -->
      <Button
        v-if="!options.hideConfigureFieldsButton"
        icon="pi pi-table"
        :label="$t('block.recordList.configureFields')"
        severity="secondary"
        outlined
        size="small"
        @click="openFieldPicker"
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
          :presets="allPresets"
          @update:model-value="onFilterChange"
          @reset="onFilterReset"
          @save-preset="onSaveFilterPreset"
          @delete-preset="onDeleteFilterPreset"
          @load-preset="onLoadFilterPreset"
        />
        <CInputSearch
          v-if="!options.hideSearch"
          v-model="searchInput"
          :placeholder="$t('general.label.search')"
          size="small"
          class="flex-1"
          @keydown.enter="commitSearch"
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
              style="border-radius: var(--p-border-radius)"
              @remove="removeFilter(fg.originalIndex, fi)"
            >
              <span class="font-medium">{{ getFieldLabel(f) }}</span>
              <span class="mx-1 text-muted-color">{{ getOperatorLabel(f.operator) }}</span>
              <span
                v-if="f.value != null"
                class="font-semibold text-primary inline-flex gap-1 flex-wrap items-center"
              >
                <template v-if="isBetweenOperator(f.operator)">
                  <CFieldViewer
                    v-if="getField(f.name)"
                    :field="getField(f.name)"
                    :record="getFilterMockRecord(f, f.value?.start)"
                    :namespace="namespace"
                  />
                  <span v-else>{{ f.value?.start || '?' }}</span>
                  <span class="text-muted-color mx-1">-</span>
                  <CFieldViewer
                    v-if="getField(f.name)"
                    :field="getField(f.name)"
                    :record="getFilterMockRecord(f, f.value?.end)"
                    :namespace="namespace"
                  />
                  <span v-else>{{ f.value?.end || '?' }}</span>
                </template>
                <template
                  v-else-if="['IN', 'NOT IN'].includes(f.operator) && !getField(f.name)?.isMulti"
                >
                  <template v-for="(v, i) in Array.isArray(f.value) ? f.value : [f.value]" :key="i">
                    <CFieldViewer
                      v-if="getField(f.name)"
                      :field="getField(f.name)"
                      :record="getFilterMockRecord(f, v)"
                      :namespace="namespace"
                    />
                    <span v-else>{{ v }}</span>
                    <span
                      v-if="i < (Array.isArray(f.value) ? f.value.length : 1) - 1"
                      class="text-muted-color"
                    >
                      ,
                    </span>
                  </template>
                </template>
                <template v-else>
                  <CFieldViewer
                    v-if="getField(f.name)"
                    :field="getField(f.name)"
                    :record="getFilterMockRecord(f, f.value)"
                    :namespace="namespace"
                  />
                  <span v-else>{{ formatFilterValue(f) }}</span>
                </template>
              </span>
              <!-- "is empty" already says what it tests; appending NULL to it
                   would read "Text is empty NULL". -->
              <span v-else-if="!isValuelessOperator(f.operator)" class="text-muted-color italic">
                {{ $t('block.recordList.filter.nil') }}
              </span>
            </Chip>
            <span
              v-if="fgi < segment.groups.length - 1"
              class="text-xs text-muted-color uppercase font-medium"
            >
              {{ $t('block.recordList.filter.conditions.and') }}
            </span>
          </template>
        </div>
        <span v-if="segment.connector" class="text-xs text-muted-color uppercase font-medium">
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
      <span v-if="selectedRecords.length && !selectedAllRecords" class="text-sm font-medium">
        {{
          $t('block.recordList.selected', {
            count: selectedRecords.length,
            total: totalRecords,
          })
        }}
      </span>
      <span v-else-if="selectedAllRecords" class="text-sm font-medium">
        {{
          $t('block.recordList.selected', {
            count: totalRecords,
            total: totalRecords,
          })
        }}
      </span>

      <Button
        v-if="
          selectedRecords.length === records.length &&
          records.length > 0 &&
          totalRecords > records.length &&
          !selectedAllRecords
        "
        :label="$t('block.recordList.selectAllRecords', { count: totalRecords })"
        text
        size="small"
        class="text-primary"
        @click="selectedAllRecords = true"
      />

      <Button
        v-if="selectedAllRecords"
        :label="$t('block.recordList.unselectAllRecords')"
        text
        size="small"
        class="text-primary"
        @click="clearSelection"
      />

      <div class="flex-1" />

      <template
        v-if="canSelectRecords && selectedRecords.length && (options.selectionButtons || []).length"
      >
        <AutomationButtons
          :buttons="options.selectionButtons || []"
          :namespace="namespace"
          :page="page"
          :module="recordListModule"
          :records="selectedRecords"
          :filter="currentQuery"
          @refresh="fetchRecords(true)"
        />
        <Divider layout="vertical" class="!mx-1 !my-0 h-6" />
      </template>

      <!-- Save / discard all dirty (or just selected dirty rows) -->
      <template v-if="showBulkSave">
        <Button
          v-tooltip.bottom="$t('block.recordList.tooltip.saveChanges')"
          icon="pi pi-check"
          text
          size="small"
          severity="primary"
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
        v-if="options.bulkRecordEditEnabled && canUpdateSelected && !showingDeletedRecords"
        v-tooltip.bottom="$t('block.recordList.bulkRecord.title')"
        icon="pi pi-pencil"
        severity="secondary"
        text
        size="small"
        @click="showBulkEditModal = true"
      />

      <Button
        v-if="canDeleteSelected"
        v-tooltip.bottom="$t('block.recordList.tooltip.deleteSelected')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        @click="deleteSelected"
      />

      <CBulkRecordEditModal
        v-model:visible="showBulkEditModal"
        :module="recordListModule"
        :namespace="namespace"
        :query="activeBulkQuery"
        allow-add-field
        @save="fetchRecords(true)"
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
          <div class="flex items-center justify-center p-4 text-muted-color text-center">
            {{
              prefilterUnresolved
                ? $t('block.recordList.noRecordContext')
                : $t('block.recordList.noRecords')
            }}
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
          :sortable="!options.hideSorting && !options.editable && !col.isMulti"
        >
          <!-- The header slot renders in place of PrimeVue's built-in title span,
               so it has to carry that span's class or the labels lose their
               weight. The mark goes inside it: header items are laid out with a
               flex gap that would otherwise push it away from the label. -->
          <template #header>
            <span class="p-datatable-column-title">
              {{ col.label }}
              <span
                v-if="showsRequiredMark(col)"
                v-tooltip.top="$t('field.required-field')"
                class="text-red-500"
                aria-hidden="true"
              >
                *
              </span>
            </span>
          </template>

          <template #body="{ data }">
            <div v-if="shouldShowEditor(data, col)">
              <CInlineFieldEditor
                :key="`${getRecordKey(data)}:${editorGeneration[getRecordKey(data)] || 0}`"
                :field="col"
                :namespace="namespace"
                :model-value="data.values[col.name]"
                :class="{ 'rounded ring-1 ring-red-500': cellError(data, col) }"
                style="min-width: 200px"
                @update:model-value="onInlineFieldUpdate(data, col.name, $event)"
                @stage-files="onInlineFilesStaged(data, $event)"
                @click.stop
              />
              <small v-if="cellError(data, col)" class="block mt-1 text-red-500">
                {{ cellError(data, col) }}
              </small>
            </div>
            <div v-else class="group flex items-start gap-1 min-w-0">
              <CFieldViewer :field="col" :record="data" :namespace="namespace" />
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
                <Button
                  v-if="showCopyFieldButton(data, col)"
                  v-tooltip.top="$t('block.recordList.tooltip.copy')"
                  icon="pi pi-copy"
                  text
                  rounded
                  size="small"
                  severity="secondary"
                  class="w-6 h-6 p-0 mt-0.5"
                  @click.stop="copyFieldValue(data, col)"
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
              <template
                v-if="(options.editable || options.inlineRecordEditEnabled) && showSaveAction(data)"
              >
                <Button
                  v-tooltip.top="$t('block.recordList.tooltip.saveChanges')"
                  icon="pi pi-check"
                  text
                  size="small"
                  severity="primary"
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
          <span v-if="options.showTotalCount !== false" class="whitespace-nowrap text-muted-color">
            {{ paginationRangeText }}
          </span>

          <Divider
            v-if="options.showRecordPerPageOption && options.showTotalCount !== false"
            layout="vertical"
          />

          <!-- Page size selector -->
          <div
            v-if="options.showRecordPerPageOption"
            class="flex items-center gap-2 whitespace-nowrap"
          >
            <span class="text-muted-color">
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

    <!-- Configure Fields modal -->
    <Dialog
      v-model:visible="showFieldPickerModal"
      modal
      :header="$t('block.recordList.configureFields')"
      :style="{ width: '90vw', maxWidth: '56rem' }"
      :contentStyle="{ minHeight: 'min(60vh, 32rem)' }"
      @hide="cancelFieldPicker"
    >
      <CFieldPicker
        :all-fields="allModuleFields"
        :model-value="localFieldNames"
        list-class="max-h-[32rem]"
        :available-label="$t('field.selector.available')"
        :selected-label="$t('field.selector.selected')"
        :select-all-label="$t('field.selector.selectAll')"
        :unselect-all-label="$t('field.selector.unselectAll')"
        :search-placeholder="$t('general.label.search')"
        :no-items-label="$t('field.no-items-found')"
        @update:model-value="localFieldNames = $event"
      />
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="cancelFieldPicker"
          />
          <Button :label="$t('general.label.apply')" size="small" @click="applyFieldPicker" />
        </div>
      </template>
    </Dialog>

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
import { compose, validator } from '@planetcrust/human-js'
import { components, useConfirmDelete, usePermissions } from '@planetcrust/human-vue'
const { CFieldViewer, CInputSearch, CFieldPicker } = components
import { useModuleStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import { useReminderStore } from '@/sections/compose/stores/reminder'
import PageBlock from './PageBlock.vue'
import { usePageStore } from '@planetcrust/human-vue'
import CBulkRecordEditModal from './CBulkRecordEditModal.vue'
import CInlineFieldEditor from './CInlineFieldEditor.vue'
import {
  mergeAttachmentIDs,
  uploadRecordAttachment,
} from '@/sections/compose/lib/record-attachments'
import AutomationButtons from '../Shared/AutomationButtons.vue'
import RecordListFilter from '../../Common/RecordListFilter.vue'
import RecordImporter from '../../Public/Record/Importer/index.vue'
import RecordExporter from '../../Public/Record/Exporter/index.vue'
import {
  evaluatePrefilter,
  usesRecordVariables,
  queryToFilter,
  getFieldFilter,
  convertRecordListFilter,
  formatActiveFilterOperator,
  isBetweenOperator,
  isValuelessOperator,
  recordListFilterStorageKey,
  recordListPresetsStorageKey,
} from '../../../lib/record-filter'

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
const $Auth = inject('$Auth', {})
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
const searchInput = ref('')
const searchQuery = ref('')
const sortField = ref(null)
const sortOrder = ref(null)

// Record list filter state
const recordListFilter = ref([])
const presetMenuRef = ref(null)
const showingDeletedRecords = ref(false)
const selectedRecords = ref([])
const showBulkEditModal = ref(false)

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

// Configure Fields modal
const showFieldPickerModal = ref(false)
const localFieldNames = ref((props.block.options?.fields || []).map(f => f.name ?? f))
let _fieldPickerSnapshot = []
let _fieldPickerApplied = false

function openFieldPicker() {
  _fieldPickerSnapshot = [...localFieldNames.value]
  _fieldPickerApplied = false
  showFieldPickerModal.value = true
}

function cancelFieldPicker() {
  if (!_fieldPickerApplied) {
    localFieldNames.value = _fieldPickerSnapshot
  }
  showFieldPickerModal.value = false
}

function applyFieldPicker() {
  _fieldPickerApplied = true
  showFieldPickerModal.value = false
}

// Inline editing — Human-style dirty tracking
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
// Disabled rather than hidden when there is no record page to open: the button
// used to vanish, which told nobody why. Inline editing needs no destination,
// so it is never disabled in that mode.
const addRecordDisabled = computed(() => !options.value.editable && !recordPageID.value)

const recordPageID = computed(() => {
  // When using custom routes, we don't need a public record page
  if ($recordRoutes) return 'admin'

  const moduleID = recordListModule.value?.moduleID
  if (!moduleID) return null

  const pages = pageStore.set
  const page = (pages || []).find(p => p.moduleID === moduleID)
  return page?.pageID || null
})

// All module fields available for the field picker (regular + system)
const allModuleFields = computed(() => {
  if (!recordListModule.value) return []
  const regular = recordListModule.value.fields || []
  const system = (recordListModule.value.systemFields?.() || []).map(f => ({
    ...f,
    label: t(`field.system.${f.name}`, f.label || f.name),
    isSystem: true,
  }))
  return [...regular, ...system]
})

// Columns derived from block field config or module fields
// localFieldNames (user runtime selection) takes precedence over block options
const columns = computed(() => {
  if (!recordListModule.value) return []

  const activeNames = localFieldNames.value.length
    ? localFieldNames.value
    : (options.value.fields || []).map(f => f.name ?? f)

  if (activeNames.length === 0) {
    // Fallback: show first 5 fields
    return (recordListModule.value.fields || []).slice(0, 5)
  }

  return recordListModule.value.filterFields(activeNames)
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
const canUpdateSelected = computed(() =>
  selectedRecords.value.some(r => r.canUpdateRecordValue !== false),
)
const hasRowActions = computed(() => {
  const o = options.value
  // View/Edit/Clone require a record page to navigate to
  const hasPageActions =
    recordPageID.value &&
    (!o.hideRecordViewButton || !o.hideRecordEditButton || !o.hideRecordCloneButton)
  // These actions don't require a record page
  const hasStandaloneActions =
    !o.hideRecordReminderButton || !o.hideRecordPermissionsButton || !o.hideRecordDeleteButton
  return hasPageActions || hasStandaloneActions
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

  if (
    recordPageID.value &&
    recordListModule.value?.canCreateRecord &&
    !options.value.hideRecordCloneButton
  ) {
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

  if (
    record?.recordID &&
    (record.canGrant ?? recordListModule.value?.canGrant) &&
    !options.value.hideRecordPermissionsButton
  ) {
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
  if (data.deletedAt && !showingDeletedRecords.value) return false
  if (showingDeletedRecords.value) return false
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

// Validation errors per row, as { recordKey: { fieldName: message } }. Shown
// under the offending cell, the way the record editors show them under a field.
const rowErrors = reactive({})

function cellError(record, col) {
  return rowErrors[getRecordKey(record)]?.[col.name]
}

// Only worth marking a column required where the list can actually be edited.
function showsRequiredMark(col) {
  if (!options.value.editable && !options.value.inlineRecordEditEnabled) return false
  return !!col.isRequired && isInlineEditField(col)
}

function clearRowErrors(record, fieldName) {
  const errors = rowErrors[getRecordKey(record)]
  if (!errors) return
  if (fieldName) delete errors[fieldName]
  if (!fieldName || !Object.keys(errors).length) delete rowErrors[getRecordKey(record)]
}

// Required fields the user can actually fill from here. One that is required but
// not shown as a column has no cell to complain in, so it is left to the server —
// its message still reaches the user through the toast below.
function missingRequiredValues(record) {
  const errors = {}
  for (const field of columns.value) {
    if (!field.isRequired || !isInlineEditField(field)) continue
    if (validator.IsEmpty(record.values[field.name])) {
      errors[field.name] = t('field.required-field')
    }
  }
  return errors
}

// Map a rejected save onto the cells it came from. Anything that does not match a
// visible column is reported in the toast instead of vanishing into the console.
function reportSaveFailure(record, e, fallbackKey) {
  const key = getRecordKey(record)
  const shown = new Set(columns.value.map(c => c.name))
  const errors = {}
  const unshown = []

  for (const detail of e?.details ?? []) {
    const field = detail.meta?.field
    if (!field || !detail.message) continue
    if (shown.has(field)) {
      errors[field] = detail.message
    } else {
      const label = recordListModule.value?.fields.find(f => f.name === field)?.label || field
      unshown.push(`${label}: ${detail.message}`)
    }
  }

  if (Object.keys(errors).length) rowErrors[key] = errors

  if (unshown.length) {
    $toast?.toastDanger(unshown.join('\n'))
  } else if (Object.keys(errors).length) {
    $toast?.toastWarning(t('general.notification.formErrors'))
  } else {
    $toast?.toastErrorHandler(t(fallbackKey))(e)
  }
}

function onInlineFieldUpdate(record, fieldName, value) {
  record.values[fieldName] = value
  clearRowErrors(record, fieldName)
  onInlineFieldChange(record)
}

// Files picked in a File cell are only uploaded once the row is saved — a brand
// new row has no recordID to attach them to yet. Hold them per row until then.
// Not reactive: nothing renders from this, and the cells keep their own preview.
const stagedFiles = new Map()

function onInlineFilesStaged(record, { fieldName, files }) {
  const key = getRecordKey(record)
  const forRecord = stagedFiles.get(key)

  if (!files.length) {
    forRecord?.delete(fieldName)
    if (forRecord && !forRecord.size) stagedFiles.delete(key)
    return
  }

  if (forRecord) {
    forRecord.set(fieldName, files)
  } else {
    stagedFiles.set(key, new Map([[fieldName, files]]))
  }
  // A staged file is an unsaved change like any other — without this the row
  // never goes dirty and there is no save control to flush it with.
  onInlineFieldChange(record)
}

// A File cell keeps its staged files in local state that the saved record cannot
// reach, so it would go on showing them as pending next to the attachment they
// became. Bumping this remounts that row's editors on the values that persisted.
const editorGeneration = reactive({})

// Upload whatever the row staged and fold the attachment IDs into its values.
// Runs before the save so a new record claims its attachments on create: the
// upload endpoint accepts an empty recordID and the server binds them there.
// Reports whether anything was uploaded, since only then do the editors go stale.
async function applyStagedFiles(record) {
  const staged = stagedFiles.get(getRecordKey(record))
  if (!staged?.size) return false

  for (const [fieldName, files] of staged) {
    const ids = await Promise.all(
      files.map(file =>
        uploadRecordAttachment($ComposeAPI, {
          namespaceID: props.namespace.namespaceID,
          moduleID: recordListModule.value.moduleID,
          recordID: record.recordID || '',
          fieldName,
          file,
        }),
      ),
    )
    // setValue rather than a plain assignment: it normalises the value to the
    // field's own multiplicity, the same way the record editors write it.
    const field = recordListModule.value.fields.find(f => f.name === fieldName)
    record.setValue(fieldName, mergeAttachmentIDs(record.values[fieldName], ids, !!field?.isMulti))
  }
  return true
}

function refreshRowEditors(key) {
  editorGeneration[key] = (editorGeneration[key] || 0) + 1
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
  const isNew = !record.recordID || record.recordID === '0'

  const missing = missingRequiredValues(record)
  if (Object.keys(missing).length) {
    rowErrors[key] = missing
    $toast?.toastWarning(t('general.notification.formErrors'))
    return
  }

  processingRecords[key] = 'save'
  try {
    const uploaded = await applyStagedFiles(record)
    const saved = isNew ? await recordStore.create(record) : await recordStore.update(record)
    delete dirtyRecords[key]
    stagedFiles.delete(key)
    clearRowErrors(record)
    clearInlineEdits(record)
    const newRecord = new compose.Record(recordListModule.value, saved)
    records.value.splice(index, 1, newRecord)
    if (uploaded) refreshRowEditors(getRecordKey(newRecord))
  } catch (e) {
    console.error('Failed to save inline record:', e)
    reportSaveFailure(
      record,
      e,
      isNew ? 'notification.record.createFailed' : 'notification.record.updateFailed',
    )
  } finally {
    delete processingRecords[key]
  }
}

async function handleDenyInline(record, index) {
  if (!recordListModule.value) return
  const key = getRecordKey(record)
  processingRecords[key] = 'deny'
  delete dirtyRecords[key]
  stagedFiles.delete(key)
  clearRowErrors(record)
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
  const toSave = records.value.filter(r => {
    if (!showSaveAction(r)) return false
    if (selectedRecords.value.length > 0) {
      return selectedRecords.value.some(s => getRecordKey(s) === getRecordKey(r))
    }
    return true
  })
  if (!toSave.length) return

  processingDirtyRecords.value = 'save'
  let hasError = false
  let hasRequiredGap = false

  for (const record of toSave) {
    const key = getRecordKey(record)
    const index = records.value.findIndex(r => getRecordKey(r) === key)
    if (index === -1) continue
    const isNew = !record.recordID || record.recordID === '0'

    // Marked in place and skipped — the other rows still save.
    const missing = missingRequiredValues(record)
    if (Object.keys(missing).length) {
      rowErrors[key] = missing
      hasError = true
      hasRequiredGap = true
      continue
    }

    try {
      const uploaded = await applyStagedFiles(record)
      const saved = isNew ? await recordStore.create(record) : await recordStore.update(record)
      delete dirtyRecords[key]
      stagedFiles.delete(key)
      clearRowErrors(record)
      clearInlineEdits(record)
      const newRecord = new compose.Record(recordListModule.value, saved)
      records.value.splice(index, 1, newRecord)
      if (uploaded) refreshRowEditors(getRecordKey(newRecord))
    } catch (e) {
      hasError = true
      console.error('Failed to save record:', e)
      reportSaveFailure(
        record,
        e,
        isNew ? 'notification.record.createFailed' : 'notification.record.updateFailed',
      )
    }
  }

  processingDirtyRecords.value = ''
  if (hasRequiredGap) $toast?.toastWarning(t('general.notification.formErrors'))
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

// A prefilter reading record variables (${record.values.x}, ${ownerID}) can only be
// resolved on a record page. In the page builder or on a list page there is no record
// to interpolate, so the block lists nothing — evaluating anyway threw on the missing
// record, and dropping the prefilter would list the whole module instead.
const prefilterUnresolved = computed(
  () => !props.record && usesRecordVariables(options.value.prefilter || ''),
)

/**
 * Build the evaluated prefilter string from block options,
 * matching Human's prepRecordList() logic.
 *
 * - Evaluates template expressions (${recordID}, ${ownerID}, ${userID})
 * - Handles refField for parent record linking
 */
function buildPrefilter() {
  const { prefilter, refField } = options.value
  const mod = recordListModule.value
  const filterParts = []

  if (prefilter && !prefilterUnresolved.value) {
    const record = props.record
    const user = $Auth?.user || {}
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

const selectedAllRecords = ref(false)

watch(selectedRecords, newVal => {
  if (newVal.length === 0 || newVal.length < records.value.length) {
    selectedAllRecords.value = false
  }
})

function clearSelection() {
  selectedRecords.value = []
  selectedAllRecords.value = false
}

const currentQuery = computed(() => {
  if (!recordListModule.value) return ''

  // An unresolvable record prefilter must not degrade into "no prefilter": every
  // consumer of this query (export, bulk actions, navigation) has to see an empty
  // set, not the whole module. recordID '0' is NoID, so nothing matches.
  if (prefilterUnresolved.value) return "recordID = '0'"

  const evaluatedPrefilter = buildPrefilter()

  let searchFields = []
  const searchableFieldConfig = options.value.searchableFields || []
  if (searchableFieldConfig.length > 0 && recordListModule.value.filterFields) {
    searchFields = recordListModule.value.filterFields(searchableFieldConfig)
  } else {
    searchFields = columns.value
  }

  const filterGroups = recordListFilter.value
    .map(g => ({ ...g, filter: convertRecordListFilter(g.filter || []) }))
    .filter(g => g.filter?.length)

  return queryToFilter(searchQuery.value || '', evaluatedPrefilter, searchFields, filterGroups)
})

const activeBulkQuery = computed(() => {
  if (selectedAllRecords.value) return currentQuery.value
  const ids = selectedRecords.value.map(r => `'${getRecordKey(r)}'`)
  if (ids.length === 1) return `recordID = ${ids[0]}`
  return `recordID IN (${ids.join(',')})`
})

// Fetch a flat list of record IDs for prev/next navigation (mirrors Human's loadPaginationRecords)
async function loadNavigationIDs() {
  if (!recordListModule.value || prefilterUnresolved.value) return
  try {
    let sort = options.value.presort || ''
    if (sortField.value) {
      sort = `${sortField.value} ${sortOrder.value === 1 ? 'ASC' : 'DESC'}`
    }
    const response = await $ComposeAPI.recordList({
      namespaceID: props.namespace.namespaceID,
      moduleID: recordListModule.value.moduleID,
      query: currentQuery.value,
      sort,
      limit: 50,
    })
    const ids = (response.set || []).map(r => r.recordID)
    recordStore.setNavigationIDs(ids)
  } catch {
    // non-critical, ignore
  }
}

// Fetch records with cancellation support
async function fetchRecords(resetCursor = false) {
  if (!recordListModule.value) return

  if (prefilterUnresolved.value) {
    console.warn(
      'Skipping record list: prefilter uses record variables outside a record page',
      options.value.prefilter,
    )
    abortPendingRequest()
    records.value = []
    totalRecords.value = 0
    nextPageCursor.value = null
    pageCursors.value = []
    currentPageIndex.value = 0
    loading.value = false
    return
  }

  abortPendingRequest()
  loading.value = true

  if (resetCursor) {
    pageCursors.value = []
    nextPageCursor.value = null
    currentPageIndex.value = 0
    selectedRecords.value = []
    // Clear inline editing state on full refresh
    Object.keys(dirtyRecords).forEach(k => delete dirtyRecords[k])
    Object.keys(rowErrors).forEach(k => delete rowErrors[k])
    activeInlineEdits.value.clear()
    stagedFiles.clear()
    Object.keys(processingRecords).forEach(k => delete processingRecords[k])
  }

  try {
    // Build sort expression
    let sort = options.value.presort || ''
    if (sortField.value) {
      sort = `${sortField.value} ${sortOrder.value === 1 ? 'ASC' : 'DESC'}`
    }

    // Build query using computed property
    const query = currentQuery.value

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
      // deleted: 0 = only existing, 2 = only deleted
      deleted: showingDeletedRecords.value ? 2 : 0,
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

const searchSubmittable = computed(
  () => (props.block.options?.searchSubmitMode || 'typing') === 'submit',
)

const commitSearch = () => {
  if (searchQuery.value === searchInput.value) return
  searchQuery.value = searchInput.value
  fetchRecords(true)
}

// Search watcher — debounced commit in typing mode; in submit mode only commit on clear
let searchDebounceTimer = null
watch(searchInput, newVal => {
  if (searchSubmittable.value) {
    if (newVal === '') commitSearch()
    return
  }
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(commitSearch, 300)
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

  if (options.value.enableRecordPageNavigation !== false) {
    loadNavigationIDs()
  }

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

function buildRecordCreatePrefillQuery() {
  if (!options.value.refField || !props.record?.recordID) return {}
  return {
    refField: options.value.refField,
    refValue: props.record.recordID,
  }
}

function handleAddRecord() {
  if (!recordPageID.value) return

  const displayOption = options.value.addRecordDisplayOption || 'sameTab'
  const refQuery = buildRecordCreatePrefillQuery()

  if (displayOption === 'modal' && !$recordRoutes) {
    router.push({
      query: {
        ...route.query,
        recordPageID: recordPageID.value,
        recordID: '0',
        ...refQuery,
      },
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

  if (Object.keys(refQuery).length > 0) {
    routeObj.query = { ...(routeObj.query || {}), ...refQuery }
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
      },
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

function createReminder(record) {
  if (!record?.recordID) return

  const sourceField = (options.value.fields || []).find(({ name }) => {
    const value = record.values?.[name]
    return Array.isArray(value) ? value.length > 0 : !!value
  })

  const fieldValue = sourceField ? record.values?.[sourceField.name] : null
  const title = Array.isArray(fieldValue) ? fieldValue.join(', ') : fieldValue || ''
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
  if (!selectedRecords.value.length && !selectedAllRecords.value) return
  loading.value = true
  if (selectedAllRecords.value) {
    try {
      await $ComposeAPI.recordBulkDelete({
        namespaceID: props.namespace.namespaceID,
        moduleID: recordListModule.value.moduleID,
        query: activeBulkQuery.value,
      })
      selectedRecords.value = []
      selectedAllRecords.value = false
      fetchRecords(true)
    } catch (e) {
      console.error('Failed to mass delete records:', e)
    }
  } else {
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
}

$eventBus?.on('refetch-records', () => fetchRecords(true))

// Reload when module changes
watch(
  () => recordListModule.value?.moduleID,
  () => {
    // Load persisted filters when module changes
    loadStoredFilter()
    fetchRecords(true)
  },
  { immediate: true },
)

// Reload when the page record changes — the block stays mounted across record
// swaps, so a query built from the record (prefilter variables, refField) would
// otherwise keep showing the previous record's rows. This is also what recovers
// the list once a record arrives and prefilterUnresolved clears.
watch(
  () => props.record?.recordID,
  () => {
    if (!options.value.refField && !usesRecordVariables(options.value.prefilter || '')) return
    fetchRecords(true)
  },
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
  return (recordListFilter.value || []).filter(g => g.filter?.some(f => f.name))
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

// Resolves a filtered field for the active-filter chips.
//
// The filter builder offers system fields alongside the module's own
// (RecordListFilter.vue builds its list from both), so looking only at
// module.fields left every createdAt/ownedBy filter unresolved: the chip fell
// back to the raw field name and the raw stored value, showing
// "createdAt >= 2024-06-29T12:30:00.000Z" where a module field showed
// "Status = Open". System labels are translated the same way the builder does.
function getField(name) {
  if (!recordListModule.value) return null

  return (
    recordListModule.value.fields.find(f => f.name === name) ||
    (recordListModule.value.systemFields?.() || []).find(f => f.name === name) ||
    null
  )
}

// The field instance is returned as-is above rather than copied with a
// translated label: a spread would drop its prototype, and with it formatValue()
// and isSystem, which are exactly what the viewer needs to render the value.
// The system-field label is translated here instead, as the builder does for its
// own picker.
function getFieldLabel(f) {
  const field = getField(f.name)
  if (!field) return f.label || f.name

  if (field.isSystem) {
    return t(`field.system.${field.name}`, field.label || field.name)
  }

  return field.label || field.name
}

function getFilterMockRecord(f, val) {
  const base = getField(f.name)
  let value = val
  if (base?.isMulti && !Array.isArray(value)) {
    value = value != null && value !== '' ? [value] : []
  }

  // A system field is read off the record itself, not out of values
  // (CFieldDateTimeViewer, CFieldUserViewer both branch on field.isSystem), so
  // the mock has to carry it in both places for the viewer to find it.
  const mock = { values: { [f.name]: value } }
  if (base?.isSystem) {
    mock[f.name] = value
  }

  return new compose.Record(recordListModule.value, mock)
}

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

// --- Deleted records toggle ---
function toggleDeletedRecords() {
  showingDeletedRecords.value = !showingDeletedRecords.value
  fetchRecords(true)
}

// --- Inline value filtering ---

function showInlineActions(col, data) {
  return (
    (options.value.inlineRecordEditEnabled && canInlineEdit(data, col)) ||
    showInlineFilter(col) ||
    showCopyFieldButton(data, col)
  )
}

function showCopyFieldButton(data, col) {
  if (!options.value.inlineRecordCopyEnabled) return false
  if (!data) return false
  if (col.canReadRecordValue === false) return false
  if (shouldShowEditor(data, col)) return false
  const val = data.values?.[col.name]
  if (val === undefined || val === null) return false
  if (Array.isArray(val) && val.length === 0) return false
  if (val === '') return false
  return true
}

function formatFieldValueForClipboard(record, col) {
  const val = record?.values?.[col.name]
  if (val === undefined || val === null) return ''
  if (Array.isArray(val)) return val.map(v => (v == null ? '' : String(v))).join('\n')
  return String(val)
}

async function copyFieldValue(record, col) {
  const text = formatFieldValueForClipboard(record, col)
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const ta = document.createElement('textarea')
      ta.value = text
      ta.style.position = 'fixed'
      ta.style.opacity = '0'
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    $toast?.toastSuccess(t('block.recordList.tooltip.copySuccess'))
  } catch (e) {
    console.error('Failed to copy field value:', e)
    $toast?.toastErrorHandler(
      t('block.recordList.tooltip.copyError'),
      t('block.recordList.tooltip.copyErrorSummary'),
    )(e)
  }
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

// Scoped to the page as well as the block — see recordListFilterStorageKey.
//
// A function, not a computed: the moduleID watcher below runs with
// `immediate: true` from higher up in this setup block, so a `const` declared
// down here is still in its temporal dead zone when loadStoredFilter() first
// runs. The ReferenceError that produces lands in that function's catch —
// written for "localStorage not available" — and the saved filter silently
// never loads. A function declaration hoists, so the first call works.
function filterStorageKey() {
  return recordListFilterStorageKey(props.page?.pageID, props.block.blockID)
}

function persistFilter() {
  try {
    const key = filterStorageKey()
    if (recordListFilter.value.length) {
      localStorage.setItem(key, JSON.stringify(recordListFilter.value))
    } else {
      localStorage.removeItem(key)
    }
  } catch {
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
  } catch {
    // localStorage not available or invalid data
  }
}

// --- Filter presets ---

// Page-scoped for the same reason the filter key is — see the helper's comment.
const userPresetsKey = computed(() =>
  recordListPresetsStorageKey(props.page?.pageID, props.block.blockID),
)
const userPresets = ref([])

// Load user presets from localStorage on init
try {
  const stored = localStorage.getItem(userPresetsKey.value)
  if (stored) userPresets.value = JSON.parse(stored)
} catch {
  // ignore
}

const allPresets = computed(() => {
  const adminPresets = (options.value.filterPresets || []).filter(p => p.name)
  return [...adminPresets, ...userPresets.value.filter(p => p.name)]
})

const visiblePresets = computed(() => allPresets.value)

const presetMenuItems = computed(() => {
  return visiblePresets.value.map(preset => ({
    label: preset.name,
    command: () => onLoadFilterPreset(preset),
  }))
})

function onLoadFilterPreset(preset) {
  if (!preset.filter) return
  if (Array.isArray(preset.filter)) {
    recordListFilter.value = JSON.parse(JSON.stringify(preset.filter))
  } else {
    recordListFilter.value = []
  }
  persistFilter()
  fetchRecords(true)
}

function onSaveFilterPreset({ name, filter }) {
  if (!name) return
  userPresets.value = [...userPresets.value, { name, filter }]
  persistUserPresets()
}

function onDeleteFilterPreset(index) {
  // Admin presets come first in the combined list, so offset for user presets
  const adminCount = (options.value.filterPresets || []).filter(p => p.name).length
  const userIndex = index - adminCount
  if (userIndex >= 0 && userIndex < userPresets.value.length) {
    const updated = [...userPresets.value]
    updated.splice(userIndex, 1)
    userPresets.value = updated
    persistUserPresets()
  }
}

function persistUserPresets() {
  try {
    if (userPresets.value.length) {
      localStorage.setItem(userPresetsKey.value, JSON.stringify(userPresets.value))
    } else {
      localStorage.removeItem(userPresetsKey.value)
    }
  } catch {
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
