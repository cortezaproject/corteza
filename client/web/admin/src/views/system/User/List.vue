<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.users.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="userID"
      :fields="userListFields"
      :items="userList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('system.users.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.users.list.single'),
        resourcePlural: $t('system.users.list.title'),
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
            :to="{ name: 'system.users.create' }"
            :label="$t('system.users.list.new')"
            icon="pi pi-plus"
            size="small"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:user/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        {{ data.name || '-' }}
      </template>

      <template #body-email="{ data }">
        {{ data.email }}
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
      </template>

      <template #body-state="{ data }">
        <Tag
          v-if="data.suspendedAt"
          :value="$t('system.users.list.rows.filters.suspended')"
          severity="danger"
        />
        <Tag
          v-else-if="data.deletedAt"
          :value="$t('system.users.list.rows.filters.deleted')"
          severity="warning"
        />
        <Tag v-else :value="$t('system.users.list.columns.enabled')" severity="success" />
      </template>

      <template #body-createdAt="{ data }">
        {{ locFullDateTime(data.createdAt) }}
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
            {{ $t('system.users.list.filterForm.suspended.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.suspended" :inputId="'susp0'" value="0" />
            <label for="susp0" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.suspended" :inputId="'susp1'" value="1" />
            <label for="susp1" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.suspended" :inputId="'susp2'" value="2" />
            <label for="susp2" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.users.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" :inputId="'del0'" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" :inputId="'del1'" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" :inputId="'del2'" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.users.list.filterForm.exclusive.label') }}
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
const canCreate = computed(() => rbac.can('system/', 'user.create'))
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const userListFields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.users.list.columns.name'),
  },
  {
    key: 'email',
    sortable: true,
    header: t('system.users.list.columns.email'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.users.list.columns.handle'),
  },
  {
    key: 'state',
    sortable: false,
    header: t('system.users.list.columns.state'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('system.users.list.columns.createdAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items: userList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.userListCancellable(params), {
  filter: { query: '', suspended: '0', deleted: '0' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (!data.canUpdateUser && !data.canDeleteUser) {
    return
  }
  router.push({
    name: 'system.users.edit',
    params: { userID: data.userID },
  })
}

// Actions menu methods
function getActionsMenuItems(user) {
  const items = []

  if (user.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::system:user/${user.userID}`,
          title: user.name || user.handle || user.email || user.userID,
          target: user.name || user.handle || user.email || user.userID,
        })
      },
    })
  }

  if (user.canUpdateUser) {
    if (items.length > 0) items.push({ separator: true })
    if (user.suspendedAt) {
      items.push({
        label: t('system.users.list.unsuspend'),
        icon: 'pi pi-play',
        command: () => handleUnsuspend(user),
      })
    } else {
      items.push({
        label: t('system.users.list.suspend'),
        icon: 'pi pi-pause',
        command: () => handleSuspend(user),
      })
    }
  }

  if (user.canDeleteUser) {
    if (items.length > 0) {
      items.push({ separator: true })
    }
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(user),
    })
  }

  return items
}

function onConfirmDelete(user) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: user.name || user.handle || user.email,
    onConfirm: () => handleDelete(user),
  })
}

async function handleDelete(user) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.userDelete({ userID: user.userID })
    $toast.toastSuccess(t('notification.user.delete.success'))
    filterList()
  } catch (e) {
    console.error('Failed to delete user:', e)
    $toast.toastErrorHandler(t('notification.user.delete.error'))(e)
  }
}

async function handleSuspend(user) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.userSuspend({ userID: user.userID })
    $toast.toastSuccess(t('notification.user.suspend.success'))
    filterList()
  } catch (e) {
    console.error('Failed to suspend user:', e)
    $toast.toastErrorHandler(t('notification.user.suspend.error'))(e)
  }
}

async function handleUnsuspend(user) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.userUnsuspend({ userID: user.userID })
    $toast.toastSuccess(t('notification.user.unsuspend.success'))
    filterList()
  } catch (e) {
    console.error('Failed to unsuspend user:', e)
    $toast.toastErrorHandler(t('notification.user.unsuspend.error'))(e)
  }
}
</script>
