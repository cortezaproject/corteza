<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.connections.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="connectionID"
      :fields="connectionListFields"
      :items="connectionList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
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
        <div class="flex gap-2">
          <CRouterLinkButton
            v-if="canCreate"
            :to="{ name: 'system.connections.create' }"
            :label="$t('system.connections.list.createLabel')"
            icon="pi pi-plus"
            size="small"
          />
          <CPermissionsButton
            v-if="canGrant"
            resource="corteza::system:dal-connection/*"
            v-tooltip.bottom="$t('general.label.permissions')"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.short || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
      </template>

      <template #body-status="{ data }">
        <Tag v-if="data.status" :value="data.status" />
        <span v-else>-</span>
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #filter>
        <Button
          icon="pi pi-filter"
          severity="secondary"
          size="small"
          text
          @click="toggleFilterMenu"
        />
      </template>
    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.connections.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.connections.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.connections.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.connections.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
  useRBACStore,
  usePermissions,
} from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const canCreate = computed(() => rbac.can('system/', 'dal-connection.create'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const connectionListFields = [
  {
    key: 'name',
    sortable: false,
    header: t('system.connections.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.connections.list.columns.handle'),
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
  filter: { query: '', deleted: '0' },
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
function getActionsMenuItems(connection) {
  const items = []

  if (connection.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::system:dal-connection/${connection.connectionID}`,
          title: connection.meta?.short || connection.handle || connection.connectionID,
        })
      },
    })
  }

  if (connection.canDeleteConnection) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
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
  resourceListRef.value.hideActionsMenu()
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
