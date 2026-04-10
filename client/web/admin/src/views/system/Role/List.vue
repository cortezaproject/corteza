<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.roles.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="roleID"
      :fields="roleListFields"
      :items="roleList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('system.roles.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.roles.list.new'),
        resourcePlural: $t('system.roles.list.title'),
      }"
      clickable
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <CRouterLinkButton
            :to="{ name: 'system.roles.create' }"
            :label="$t('system.roles.list.new')"
            icon="pi pi-plus"
            size="small"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:role/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.name || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
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
            {{ $t('system.roles.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.roles.list.filterForm.archived.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch0" value="0" />
            <label for="arch0" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch1" value="1" />
            <label for="arch1" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch2" value="2" />
            <label for="arch2" class="text-sm cursor-pointer">
              {{ $t('system.roles.list.filterForm.exclusive.label') }}
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
  usePermissions,
  useRBACStore,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// Resource list ref (for hideActionsMenu)
const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const roleListFields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.roles.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.roles.list.columns.handle'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.roles.editor.info.updatedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items: roleList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.roleListCancellable(params), {
  filter: { query: '', deleted: '0', archived: '0' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (!data.canUpdateRole && !data.canDeleteRole) {
    return
  }
  router.push({
    name: 'system.roles.edit',
    params: { roleID: data.roleID },
  })
}

// Actions menu methods
function getActionsMenuItems(role) {
  const items = []

  if (role.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::system:role/${role.roleID}`,
          title: role.name || role.handle || role.roleID,
          target: role.name || role.handle || role.roleID,
        })
      },
    })
  }

  if (role.canDeleteRole) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(role),
    })
  }

  return items
}

function onConfirmDelete(role) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: role.name || role.handle || t('system.roles.list.new'),
    onConfirm: () => handleDelete(role),
  })
}

async function handleDelete(role) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.roleDelete({
      roleID: role.roleID,
    })
    $toast.toastSuccess(t('notification.role.delete.success'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete role:', e)
    $toast.toastErrorHandler(t('notification.role.delete.error'))(e)
  }
}
</script>
