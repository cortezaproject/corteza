<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.user-groups.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="userGroupID"
      :fields="userGroupListFields"
      :items="userGroupList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('system.user-groups.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.user-groups.list.new'),
        resourcePlural: $t('system.user-groups.list.title'),
      }"
      clickable
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <CRouterLinkButton
            v-if="canCreate"
            :to="{ name: 'system.userGroups.create' }"
            :label="$t('system.user-groups.list.new')"
            icon="pi pi-plus"
            size="small"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:user-group/*"
          />
        </div>
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.short || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
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
            {{ $t('system.user-groups.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.user-groups.list.filterForm.archived.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch0" value="0" />
            <label for="arch0" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch1" value="1" />
            <label for="arch1" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.archived" inputId="arch2" value="2" />
            <label for="arch2" class="text-sm cursor-pointer">
              {{ $t('system.user-groups.list.filterForm.exclusive.label') }}
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
} from '@planetcrust/human-vue'
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
const canCreate = computed(() => rbac.can('system/', 'user-group.create'))
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const userGroupListFields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.user-groups.list.columns.meta.short'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.user-groups.list.columns.handle'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.user-groups.editor.info.updatedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items: userGroupList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.userGroupListCancellable(params), {
  filter: { query: '', deleted: '0', archived: '0' },
  sorting: { sortBy: 'name', sortDesc: false },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (!data.canUpdateUserGroup && !data.canDeleteUserGroup) {
    return
  }
  router.push({
    name: 'system.userGroups.edit',
    params: { userGroupID: data.userGroupID },
  })
}

// Actions menu methods
function getActionsMenuItems(userGroup) {
  const items = []

  if (userGroup.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::system:user-group/${userGroup.userGroupID}`,
          title: userGroup.meta?.short || userGroup.handle || userGroup.userGroupID,
          target: userGroup.meta?.short || userGroup.handle || userGroup.userGroupID,
        })
      },
    })
  }

  if (userGroup.canDeleteUserGroup) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(userGroup),
    })
  }

  return items
}

function onConfirmDelete(userGroup) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: userGroup.meta?.short || userGroup.handle || t('system.user-groups.list.new'),
    onConfirm: () => handleDelete(userGroup),
  })
}

async function handleDelete(userGroup) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.userGroupDelete({
      userGroupID: userGroup.userGroupID,
    })
    $toast.toastSuccess(t('notification.userGroup.delete.success'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete user group:', e)
    $toast.toastErrorHandler(t('notification.userGroup.delete.error'))(e)
  }
}
</script>
