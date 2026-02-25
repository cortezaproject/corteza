<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.connections.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="connectionID"
      :fields="connectionListFields"
      :items="connectionList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.connections.list.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.connections.list.resourceSingle'),
        resourcePlural: $t('system.connections.list.resourcePlural'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <CRouterLinkButton
          :to="{ name: 'system.connections.create' }"
          :label="$t('system.connections.list.createLabel')"
          icon="pi pi-plus"
          size="small"
        />
      </template>

      <template #body-handle="{ data }">
        <span class="font-medium">{{ data.handle || '-' }}</span>
      </template>

      <template #body-name="{ data }">
        {{ data.meta?.short || '-' }}
      </template>

      <template #body-status="{ data }">
        <Tag v-if="data.status" :value="data.status" />
        <span v-else>-</span>
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

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])

// Column definitions
const connectionListFields = [
  {
    key: 'handle',
    sortable: true,
    header: t('system.connections.list.columns.handle'),
  },
  {
    key: 'name',
    sortable: false,
    header: t('system.connections.list.columns.name'),
  },
  {
    key: 'status',
    sortable: true,
    header: t('system.connections.list.columns.status'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.connections.list.columns.updatedAt'),
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
      bodyCell: { class: 'p-0 border-l-0' },
    },
  },
]

// Resource list composable
const {
  items: connectionList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.connectionListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (!data.canUpdateConnection && !data.canDeleteConnection) {
    return
  }
  router.push({
    name: 'system.connections.edit',
    params: { connectionID: data.connectionID },
  })
}

// Actions menu methods
function toggleActionsMenu(event, connection) {
  actionsMenuItems.value = getActionsMenuItems(connection)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(connection) {
  const items = []

  if (connection.canDeleteConnection) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => onConfirmDelete(connection),
    })
  }

  return items
}

function onConfirmDelete(connection) {
  confirmDelete({
    message: t('system.connections.list.deleteConfirm'),
    header:
      connection.meta?.short || connection.handle || t('system.connections.list.resourceSingle'),
    onConfirm: () => handleDelete(connection),
  })
}

async function handleDelete(connection) {
  try {
    await $SystemAPI.connectionDelete({
      connectionID: connection.connectionID,
    })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete connection:', e)
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  }
}
</script>
