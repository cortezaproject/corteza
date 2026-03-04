<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chart.navigation.chart') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="chartID"
      :fields="chartFields"
      :items="chartList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
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
        <CRouterLinkButton
          v-if="namespace?.canCreateChart"
          :to="{ name: 'admin.charts.create' }"
          :label="$t('chart.createLabel')"
          icon="pi pi-plus"
          size="small"
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
          class="row-action-btn w-full"
          @click.stop="toggleActionsMenu($event, data)"
        />
      </template>
    </CResourceList>

    <TieredMenu ref="actionsMenu" :model="actionsMenuItems" popup>
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
    </TieredMenu>
  </div>
</template>

<script setup>
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList, CRouterLinkButton } = components
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
    pt: {
      headerCell: { class: 'border-l-0' },
      bodyCell: { class: 'px-2 py-1 border-l-0' },
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

// Actions menu methods
function toggleActionsMenu(event, chart) {
  actionsMenuItems.value = getActionsMenuItems(chart)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(chart) {
  const items = []

  if (chart.canDeleteChart) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
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
