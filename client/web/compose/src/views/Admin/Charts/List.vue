<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chart.navigation.chart') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="chartID"
      :fields="chartFields"
      :items="chartList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('chart.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.label.chart.single'),
        resourcePlural: $t('general.label.chart.plural'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <Button
            v-if="namespace?.canCreateChart"
            :label="$t('chart.createLabel')"
            icon="pi pi-plus"
            size="small"
            @click="showTypeSelector = true"
          />
          <Button
            v-if="namespace?.canExportCharts && chartList.length"
            :label="$t('general.label.export')"
            icon="pi pi-download"
            size="small"
            severity="secondary"
            @click="exportAllCharts"
          />
          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::compose:chart/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.name }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-handle="{ data }">
        <code class="text-sm">{{ data.handle || '-' }}</code>
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

    </CResourceList>

    <!-- Chart Type Selector Dialog -->
    <Dialog
      v-model:visible="showTypeSelector"
      :header="$t('chart.createLabel')"
      modal
      :style="{ width: '600px' }"
    >
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 p-2">
        <div
          v-for="ct in chartTypes"
          :key="ct.category"
          class="flex flex-col items-center gap-2 p-4 border border-surface rounded cursor-pointer hover:bg-highlight transition-colors"
          @click="createChart(ct.category)"
        >
          <i :class="ct.icon" class="text-2xl text-primary" />
          <span class="text-sm font-medium text-center">{{ ct.label }}</span>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import {
  components,
  filters,
  useConfirmDelete,
  usePermissions,
  useRBACStore,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList } = components
const { locFullDateTime } = filters

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('compose/', 'grant'))
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')

const resourceListRef = ref()
const showTypeSelector = ref(false)

const chartTypes = [
  { category: '', label: t('block.chart.addGeneric'), icon: 'pi pi-chart-bar' },
  { category: 'funnel', label: t('block.chart.addFunnel'), icon: 'pi pi-sort-amount-down' },
  { category: 'gauge', label: t('block.chart.addGauge'), icon: 'pi pi-gauge' },
  { category: 'radar', label: t('block.chart.addRadar'), icon: 'pi pi-chart-scatter' },
]

// Column definitions
const chartFields = [
  {
    key: 'name',
    sortable: true,
    header: t('chart.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('chart.list.columns.handle'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('chart.list.columns.changedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items: chartList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(
  params =>
    $ComposeAPI.chartListCancellable({
      namespaceID: props.namespace.namespaceID,
      ...params,
    }),
  {
    filter: { query: '' },
    sorting: { sortBy: 'name', sortDesc: false },
    pagination: { limit: 50 },
  },
)

// Methods
function handleRowClick({ data }) {
  if (!(data.canUpdateChart || data.canDeleteChart)) {
    return
  }
  router.push({
    name: 'admin.charts.edit',
    params: { chartID: data.chartID },
  })
}

function createChart(category) {
  showTypeSelector.value = false
  router.push({
    name: 'admin.charts.create',
    query: category ? { category } : {},
  })
}

// Actions menu methods
function getActionsMenuItems(chart) {
  const items = []

  if (chart.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::compose:chart/${props.namespace.namespaceID}/${chart.chartID}`,
          title: chart.name || chart.handle || chart.chartID,
          target: chart.name || chart.handle || chart.chartID,
        })
      },
    })
  }

  if (props.namespace?.canExportCharts) {
    items.push({
      label: t('general.label.export'),
      icon: 'pi pi-download',
      command: () => exportChart(chart),
    })
  }

  if (chart.canDeleteChart) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(chart),
    })
  }

  return items
}

function onConfirmDelete(chart) {
  confirmDelete({
    message: t('chart.list.delete'),
    header: chart.name,
    onConfirm: () => handleDelete(chart),
  })
}

async function handleDelete(chart) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $ComposeAPI.chartDelete({
      namespaceID: props.namespace.namespaceID,
      chartID: chart.chartID,
    })
    $toast.toastSuccess(t('notification.chart.deleted'))
    filterList()
  } catch (e) {
    console.error('Failed to delete chart:', e)
    $toast.toastDanger(t('notification.chart.deleteFailed'))
  }
}

function exportChart(chart) {
  const blob = new Blob(
    [JSON.stringify({ type: 'chart', list: [chart] }, null, 2)],
    { type: 'application/json' },
  )
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${chart.handle || chart.name || 'chart'}-export.json`
  a.click()
  URL.revokeObjectURL(url)
}

function exportAllCharts() {
  const blob = new Blob(
    [JSON.stringify({ type: 'chart', list: chartList.value }, null, 2)],
    { type: 'application/json' },
  )
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'charts-export.json'
  a.click()
  URL.revokeObjectURL(url)
}
</script>
