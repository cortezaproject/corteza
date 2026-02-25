<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.users.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="userID"
      :fields="userListFields"
      :items="userList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.users.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
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
            :to="{ name: 'system.users.create' }"
            :label="$t('system.users.list.new')"
            icon="pi pi-plus"
            size="small"
          />
          <Button
            icon="pi pi-filter"
            severity="secondary"
            outlined
            size="small"
            @click="toggleFilterMenu"
          />
        </div>
      </template>

      <template #body-email="{ data }">
        {{ data.email }}
      </template>
      <template #body-handle="{ data }">
        <span class="font-medium">{{ data.handle || '-' }}</span>
      </template>

      <template #body-name="{ data }">
        {{ data.name || '-' }}
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
        <a v-else v-ripple v-bind="props.action" :class="item.class" @click="item.command">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </template>
    </TieredMenu>

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

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const userListFields = [
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
    key: 'name',
    sortable: true,
    header: t('system.users.list.columns.name'),
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
function toggleActionsMenu(event, user) {
  actionsMenuItems.value = getActionsMenuItems(user)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(user) {
  const items = []

  if (user.canUpdateUser) {
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
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
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
  try {
    await $SystemAPI.userDelete({ userID: user.userID })
    $toast.toastSuccess(t('notification.user.delete.success'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete user:', e)
    $toast.toastErrorHandler(t('notification.user.delete.error'))(e)
  }
}

async function handleSuspend(user) {
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
