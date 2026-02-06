<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chart.navigation.chart') }}</span>
  </Teleport>

  <div class="container mx-auto p-5 h-full overflow-hidden">
    <CResourceList
      primary-key="chartID"
      :fields="chartFields"
      :items="chartList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('chart.searchPlaceholder'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
    >
      <template #header>
        <Button
          v-if="namespace?.canCreateChart"
          :label="$t('chart.createLabel')"
          icon="pi pi-plus"
          @click="$router.push({ name: 'admin.charts.create' })"
        />
      </template>

      <template #body-name="{ data }">
        <span class="font-medium">{{ data.name }}</span>
      </template>

      <template #body-handle="{ data }">
        <code class="text-sm">{{ data.handle || '-' }}</code>
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #body-actions="{ data }">
        <Button
          icon="pi pi-ellipsis-v"
          text
          severity="secondary"
          size="small"
          @click.stop="toggleActionsMenu($event, data)"
        />
      </template>
    </CResourceList>

    <TieredMenu ref="actionsMenu" :model="actionsMenuItems" popup />
  </div>
</template>

<script setup>
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'
import { useConfirm } from 'primevue/useconfirm'
import { inject, ref } from 'vue'
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
const confirm = useConfirm()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])

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
  {
    key: 'actions',
    class: 'text-right w-12',
    header: '',
    frozen: true,
    alignFrozen: 'right',
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

// Actions menu methods
function toggleActionsMenu(event, chart) {
  actionsMenuItems.value = getActionsMenuItems(chart)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(chart) {
  const items = []

  if (chart.canUpdateChart) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      command: () =>
        router.push({
          name: 'admin.charts.edit',
          params: { chartID: chart.chartID },
        }),
    })
  }

  if (chart.canDeleteChart) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => confirmDelete(chart),
    })
  }

  return items
}

function confirmDelete(chart) {
  confirm.require({
    message: t('chart.list.delete'),
    header: chart.name,
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: () => handleDelete(chart),
  })
}

async function handleDelete(chart) {
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
</script>
