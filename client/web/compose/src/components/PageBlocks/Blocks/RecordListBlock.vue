<template>
  <PageBlock :block="block">
    <!-- Toolbar -->
    <div v-if="recordListModule" class="flex items-center gap-2 p-3 border-b">
      <!-- Add Record button -->
      <Button
        v-if="!options.hideAddButton && recordPageID && recordListModule?.canCreateRecord"
        :label="$t('block.recordList.addRecord')"
        icon="pi pi-plus"
        size="small"
        @click="handleAddRecord"
      />

      <!-- Spacer -->
      <div class="flex-1" />

      <!-- Search -->
      <CInputSearch
        v-if="!options.hideSearch"
        v-model="searchQuery"
        :placeholder="$t('general.label.search')"
        size="small"
        class="flex-1 max-w-xl"
      />
    </div>

    <!-- Selection bar -->
    <div
      v-if="selectedRecords.length"
      class="flex items-center gap-2 px-3 py-2 bg-highlight border-b"
    >
      <span class="text-sm font-medium">
        {{
          $t('block.recordList.selected', {
            count: selectedRecords.length,
            total: totalRecords,
          })
        }}
      </span>

      <div class="flex-1" />

      <Button
        v-if="canDeleteSelected"
        :label="$t('block.recordList.tooltip.deleteSelected')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        @click="confirmDeleteSelected"
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
        scrollable
        scroll-height="flex"
        row-hover
        show-gridlines
        resizable-columns
        column-resize-mode="expand"
        data-key="recordID"
        :sort-field="sortField"
        :sort-order="sortOrder"
        @sort="onSort"
        @row-click="onRowClick"
        class="record-list-table"
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
          :sortable="!options.hideSorting"
        >
          <template #body="{ data }">
            <CFieldViewer :field="col" :record="data" :namespace="namespace" />
          </template>
        </Column>

        <!-- Actions column -->
        <Column
          v-if="hasRowActions"
          header-style="width: 3rem"
          :pt="{
            headerCell: { class: 'border-l-0' },
            bodyCell: { class: 'p-0 border-l-0' },
          }"
          frozen
          align-frozen="right"
        >
          <template #body="{ data }">
            <Button
              icon="pi pi-ellipsis-v"
              text
              size="small"
              severity="secondary"
              class="row-action-btn w-full mr-2"
              @click="openRowMenu($event, data)"
              @click.stop
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <template v-if="recordListModule && !options.hidePaging && totalRecords > 0" #footer>
      <div class="flex items-center flex-wrap gap-2 px-3 py-2">
        <div class="flex items-center text-sm">
          <span class="whitespace-nowrap text-muted">
            {{ paginationRangeText }}
          </span>

          <Divider layout="vertical" />

          <!-- Page size selector -->
          <div class="flex items-center gap-2 whitespace-nowrap">
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
    <Menu ref="rowMenuRef" :model="rowMenuItems" :popup="true">
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
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Menu from 'primevue/menu'
import Select from 'primevue/select'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer, CInputSearch } = components
import { useModuleStore } from '@/stores/module'
import { useRecordStore } from '@/stores/record'
import PageBlock from './PageBlock.vue'
import { usePageStore } from '@/stores/page'

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
})

const { t } = useI18n()
const router = useRouter()
const { confirmDelete } = useConfirmDelete()
const $ComposeAPI = inject('$ComposeAPI')
const moduleStore = useModuleStore()
const recordStore = useRecordStore()
const pageStore = usePageStore()

const loading = ref(false)
const records = ref([])
const totalRecords = ref(0)
const searchQuery = ref('')
const sortField = ref(null)
const sortOrder = ref(null)
const selectedRecords = ref([])

// Pagination state — cursor based
const pageCursors = ref([]) // stack of previous page cursors
const nextPageCursor = ref(null)
const currentPageIndex = ref(0) // 0-based page index for display

// Page size
const pageSizeOptions = [10, 20, 50, 100]
const currentPerPage = ref(20)

// Row action menu
const rowMenuRef = ref(null)
const activeRowRecord = ref(null)

const options = computed(() => props.block.options || {})

// Initialize perPage from block options
watch(
  () => options.value.perPage,
  val => {
    if (val && pageSizeOptions.includes(val)) {
      currentPerPage.value = val
    }
  },
  { immediate: true },
)

// Resolve the module for this block
const recordListModule = computed(() => {
  const moduleID = options.value.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Find the record page for this module (page with matching moduleID)
const recordPageID = computed(() => {
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
  return recordPageID.value || recordListModule.value?.canCreateRecord
})

// Row action menu items — filtered by per-record permissions
const rowMenuItems = computed(() => {
  const record = activeRowRecord.value
  const items = []

  if (recordPageID.value && record?.canReadRecord) {
    items.push({
      label: t('block.recordList.record.tooltip.view'),
      icon: 'pi pi-eye',
      route: {
        name: 'page.record',
        params: {
          pageID: recordPageID.value,
          recordID: record.recordID,
        },
      },
    })
  }

  if (recordPageID.value && record?.canUpdateRecord) {
    items.push({
      label: t('block.recordList.record.tooltip.edit'),
      icon: 'pi pi-pencil',
      route: {
        name: 'page.record',
        params: {
          pageID: recordPageID.value,
          recordID: record.recordID,
        },
        query: { edit: '1' },
      },
    })
  }

  if (recordPageID.value && recordListModule.value?.canCreateRecord) {
    items.push({
      label: t('block.recordList.record.tooltip.clone'),
      icon: 'pi pi-copy',
      command: () => handleCloneRecord(record),
    })
  }

  if (record?.canDeleteRecord) {
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

// Request cancellation
let cancelPendingRequest = null

function abortPendingRequest() {
  if (cancelPendingRequest) {
    cancelPendingRequest()
    cancelPendingRequest = null
  }
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
  }

  try {
    // Build sort expression
    let sort = options.value.presort || ''
    if (sortField.value) {
      sort = `${sortField.value} ${sortOrder.value === 1 ? 'ASC' : 'DESC'}`
    }

    // Build query
    let query = options.value.prefilter || ''
    if (searchQuery.value) {
      const searchFilter = searchQuery.value
      query = query ? `(${query}) AND ${searchFilter}` : searchFilter
    }

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
      incTotal: true,
    })

    cancelPendingRequest = cancel
    const result = await response()
    cancelPendingRequest = null

    const mod = recordListModule.value
    const set = (result.set || []).map(r => new compose.Record(mod, r))

    records.value = set
    totalRecords.value = result.filter?.total || set.length

    // Store next page cursor
    nextPageCursor.value = result.filter?.nextPage || null
    loading.value = false
  } catch (error) {
    if (!axios.isCancel(error)) {
      console.error('Failed to fetch records for RecordList block:', error)
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

// Cleanup on unmount
onBeforeUnmount(() => {
  abortPendingRequest()
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
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

  // Don't navigate if clicking on a checkbox or action button
  const target = event.originalEvent?.target
  if (target?.closest('.p-checkbox') || target?.closest('button')) return

  if (!recordPageID.value) return

  const route = {
    name: 'page.record',
    params: {
      pageID: recordPageID.value,
      recordID: record.recordID,
    },
  }

  const displayOption = options.value.recordDisplayOption || 'sameTab'

  if (displayOption === 'newTab') {
    const resolved = router.resolve(route)
    window.open(resolved.href, '_blank')
  } else {
    router.push(route)
  }
}

// Row actions
function openRowMenu(event, record) {
  activeRowRecord.value = record
  rowMenuRef.value.toggle(event)
}

function handleAddRecord() {
  if (!recordPageID.value) return

  router.push({
    name: 'page.record',
    params: {
      pageID: recordPageID.value,
      recordID: '0',
    },
  })
}

function handleCloneRecord(record) {
  if (!record || !recordPageID.value) return

  router.push({
    name: 'page.record',
    params: {
      pageID: recordPageID.value,
      recordID: '0',
    },
    query: { cloneFromID: record.recordID },
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

function confirmDeleteSelected() {
  if (!selectedRecords.value.length) return

  confirmDelete({
    message: t('block.recordList.tooltip.deleteSelected'),
    header: t('block.recordList.record.tooltip.delete'),
    onConfirm: async () => {
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
    },
  })
}

// Reload when module changes
watch(
  () => recordListModule.value?.moduleID,
  () => fetchRecords(true),
  { immediate: true },
)
</script>

<style scoped>
.record-list-table :deep(.p-datatable-tbody > tr) {
  cursor: pointer;
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
