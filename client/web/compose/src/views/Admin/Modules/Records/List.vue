<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">{{ recordModule.name }} — {{ $t('module.allRecords.label') }}</span>
  </Teleport>

  <!-- Module not found -->
  <div v-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Records list -->
  <div v-else-if="recordModule" class="flex flex-col h-full">
    <!-- Toolbar -->
    <div class="flex items-center gap-2 p-3 border-b border-surface">
      <Button
        v-if="recordModule.canCreateRecord"
        :label="$t('module.edit.createRecord')"
        icon="pi pi-plus"
        size="small"
        @click="handleCreate"
      />

      <div class="flex-1" />

      <CInputSearch
        v-model="searchQuery"
        :placeholder="$t('general.label.search')"
        size="small"
        class="max-w-sm"
      />
    </div>

    <!-- Selection bar -->
    <div
      v-if="selectedRecords.length"
      class="flex items-center gap-2 px-3 py-2 bg-highlight border-b border-surface"
    >
      <span class="text-sm font-medium">
        {{ $t('block.recordList.selected', { count: selectedRecords.length, total: totalRecords }) }}
      </span>
      <div class="flex-1" />
      <Button
        :label="$t('general.label.delete')"
        icon="pi pi-trash"
        severity="danger"
        text
        size="small"
        @click="confirmDeleteSelected"
      />
    </div>

    <!-- Data table -->
    <div class="flex-1 overflow-auto">
      <DataTable
        v-model:selection="selectedRecords"
        :value="records"
        :loading="loading"
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
        class="records-table"
      >
        <template #empty>
          <div class="flex items-center justify-center p-4 text-muted-color">
            {{ $t('block.recordList.noRecords') }}
          </div>
        </template>

        <!-- Selection column -->
        <Column
          selection-mode="multiple"
          header-style="width: 3rem"
          frozen
        />

        <!-- Data columns -->
        <Column
          v-for="col in columns"
          :key="col.name"
          :field="col.name"
          :header="col.label || col.name"
          sortable
        >
          <template #body="{ data }">
            <CFieldViewer :field="col" :record="data" :namespace="namespace" />
          </template>
        </Column>

        <!-- Actions column -->
        <Column
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
              class="row-action-btn w-full"
              @click.stop="openRowMenu($event, data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Pagination footer -->
    <div
      v-if="totalRecords > 0"
      class="flex items-center flex-wrap gap-2 px-3 py-2 border-t border-surface"
    >
      <span class="text-sm text-muted-color whitespace-nowrap">
        {{ paginationText }}
      </span>

      <div class="flex items-center gap-2 ml-2">
        <span class="text-sm text-muted-color">{{ $t('block.recordList.pagination.recordsPerPage') }}</span>
        <Select
          v-model="perPage"
          :options="pageSizeOptions"
          size="small"
          class="w-20"
          @change="fetchRecords(true)"
        />
      </div>

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
  </div>

  <!-- Row action menu -->
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
</template>

<script setup>
import axios from 'axios'
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Menu from 'primevue/menu'
import Select from 'primevue/select'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer, CInputSearch } = components
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
const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const loading = ref(false)
const records = ref([])
const totalRecords = ref(0)
const searchQuery = ref('')
const sortField = ref(null)
const sortOrder = ref(null)
const selectedRecords = ref([])

// Pagination
const pageSizeOptions = [10, 20, 50, 100]
const perPage = ref(20)
const pageCursors = ref([])
const nextPageCursor = ref(null)
const currentPageIndex = ref(0)

// Row menu
const rowMenuRef = ref(null)
const activeRowRecord = ref(null)

const moduleID = computed(() => route.params.moduleID)
const recordModule = computed(() => (moduleID.value ? moduleStore.getByID(moduleID.value) : null))

// Columns: all module fields (up to 8 for readability)
const columns = computed(() => {
  if (!recordModule.value) return []
  return (recordModule.value.fields || []).slice(0, 8)
})

const hasPrevPage = computed(() => pageCursors.value.length > 0)
const hasNextPage = computed(() => !!nextPageCursor.value)

const paginationText = computed(() => {
  const from = currentPageIndex.value * perPage.value + 1
  const to = Math.min(from + perPage.value - 1, totalRecords.value)
  return t('block.recordList.pagination.showing', { from, to, count: totalRecords.value })
})

const rowMenuItems = computed(() => {
  const rec = activeRowRecord.value
  const items = []

  if (rec?.canReadRecord) {
    items.push({
      label: t('block.recordList.record.tooltip.view'),
      icon: 'pi pi-eye',
      route: {
        name: 'admin.modules.record.view',
        params: { moduleID: moduleID.value, recordID: rec.recordID },
      },
    })
  }

  if (rec?.canUpdateRecord) {
    items.push({
      label: t('block.recordList.record.tooltip.edit'),
      icon: 'pi pi-pencil',
      route: {
        name: 'admin.modules.record.edit',
        params: { moduleID: moduleID.value, recordID: rec.recordID },
      },
    })
  }

  if (rec?.canDeleteRecord) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('block.recordList.record.tooltip.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => confirmDeleteRecord(rec),
    })
  }

  return items
})

let cancelRequest = null

function abortRequest() {
  if (cancelRequest) {
    cancelRequest()
    cancelRequest = null
  }
}

async function fetchRecords(reset = false) {
  if (!recordModule.value) return

  abortRequest()
  loading.value = true

  if (reset) {
    pageCursors.value = []
    nextPageCursor.value = null
    currentPageIndex.value = 0
    selectedRecords.value = []
  }

  try {
    let sort = ''
    if (sortField.value) {
      sort = `${sortField.value} ${sortOrder.value === 1 ? 'ASC' : 'DESC'}`
    }

    let pageCursor
    if (!reset && pageCursors.value.length > 0) {
      pageCursor = pageCursors.value[pageCursors.value.length - 1]
    }

    const { response, cancel } = $ComposeAPI.recordListCancellable({
      namespaceID: props.namespace.namespaceID,
      moduleID: moduleID.value,
      query: searchQuery.value || '',
      sort,
      limit: perPage.value,
      pageCursor,
      incTotal: true,
    })

    cancelRequest = cancel
    const result = await response()
    cancelRequest = null

    const mod = recordModule.value
    records.value = (result.set || []).map(r => new compose.Record(mod, r))
    totalRecords.value = result.filter?.total || records.value.length
    nextPageCursor.value = result.filter?.nextPage || null
  } catch (e) {
    if (!axios.isCancel(e)) {
      console.error('Failed to fetch records:', e)
      cancelRequest = null
      records.value = []
      totalRecords.value = 0
    }
  } finally {
    if (cancelRequest === null) {
      loading.value = false
    }
  }
}

// Debounced search
let searchTimer = null
watch(searchQuery, () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchRecords(true), 300)
})

onBeforeUnmount(() => {
  abortRequest()
  if (searchTimer) clearTimeout(searchTimer)
})

function onSort(event) {
  sortField.value = event.sortField
  sortOrder.value = event.sortOrder
  fetchRecords(true)
}

function onRowClick(event) {
  const rec = event.data
  if (!rec?.recordID) return

  const target = event.originalEvent?.target
  if (target?.closest('.p-checkbox') || target?.closest('button')) return

  if (rec.canReadRecord) {
    router.push({
      name: 'admin.modules.record.view',
      params: { moduleID: moduleID.value, recordID: rec.recordID },
    })
  }
}

function openRowMenu(event, rec) {
  activeRowRecord.value = rec
  rowMenuRef.value.toggle(event)
}

function handleCreate() {
  router.push({
    name: 'admin.modules.record.create',
    params: { moduleID: moduleID.value },
  })
}

function goToNextPage() {
  if (!nextPageCursor.value) return
  pageCursors.value.push(nextPageCursor.value)
  currentPageIndex.value++
  fetchRecords()
}

function goToPrevPage() {
  if (!pageCursors.value.length) return
  pageCursors.value.pop()
  currentPageIndex.value = Math.max(0, currentPageIndex.value - 1)
  fetchRecords()
}

function confirmDeleteRecord(rec) {
  confirmDelete({
    message: t('block.record.confirmDelete'),
    header: t('block.recordList.record.tooltip.delete'),
    onConfirm: async () => {
      try {
        await recordStore.delete({
          namespaceID: props.namespace.namespaceID,
          moduleID: moduleID.value,
          recordID: rec.recordID,
        })
        fetchRecords(true)
      } catch (e) {
        console.error('Failed to delete record:', e)
        $toast.toastDanger(t('notification.record.deleteFailed'))
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
      const results = await Promise.allSettled(
        selectedRecords.value.map(rec =>
          recordStore.delete({
            namespaceID: props.namespace.namespaceID,
            moduleID: moduleID.value,
            recordID: rec.recordID,
          }),
        ),
      )
      const failed = results.filter(r => r.status === 'rejected')
      if (failed.length) {
        console.error('Failed to delete some records:', failed)
        $toast.toastDanger(t('notification.record.deleteBulkFailed'))
      } else {
        $toast.toastSuccess(t('notification.record.deleteBulkSuccess'))
      }
      selectedRecords.value = []
      fetchRecords(true)
    },
  })
}

// Load records when module is ready
watch(
  () => recordModule.value?.moduleID,
  () => {
    if (recordModule.value) fetchRecords(true)
  },
  { immediate: true },
)
</script>

<style scoped>
.records-table :deep(.p-datatable-tbody > tr) {
  cursor: pointer;
}

.records-table :deep(.row-action-btn) {
  opacity: 0;
  transition: opacity 0.15s ease;
}

.records-table :deep(.p-datatable-tbody > tr:hover .row-action-btn) {
  opacity: 1;
}

:deep(.p-datatable-gridlines .p-datatable-thead > tr > th) {
  border-top: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :first-child) {
  border-left: 0;
}

:deep(.p-datatable-gridlines :is(.p-datatable-thead, .p-datatable-tbody) tr > :last-child) {
  border-right: 0;
}
</style>
