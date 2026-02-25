<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.user-groups.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="userGroupID"
      :fields="userGroupListFields"
      :items="userGroupList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
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
        <CRouterLinkButton
          :to="{ name: 'system.userGroups.create' }"
          :label="$t('system.user-groups.list.new')"
          icon="pi pi-plus"
          size="small"
        />
      </template>

      <template #body-handle="{ data }">
        <span class="font-medium">{{ data.handle || '-' }}</span>
      </template>

      <template #body-name="{ data }">
        {{ data.name || '-' }}
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
const userGroupListFields = [
  {
    key: 'handle',
    sortable: true,
    header: t('system.user-groups.list.columns.handle'),
  },
  {
    key: 'name',
    sortable: true,
    header: t('system.user-groups.list.columns.meta.short', 'Name'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.user-groups.editor.info.updatedAt', 'Updated at'),
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
  items: userGroupList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.userGroupListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
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
function toggleActionsMenu(event, userGroup) {
  actionsMenuItems.value = getActionsMenuItems(userGroup)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(userGroup) {
  const items = []

  if (userGroup.canDeleteUserGroup) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => onConfirmDelete(userGroup),
    })
  }

  return items
}

function onConfirmDelete(userGroup) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: userGroup.name || userGroup.handle || t('system.user-groups.list.new'),
    onConfirm: () => handleDelete(userGroup),
  })
}

async function handleDelete(userGroup) {
  try {
    await $SystemAPI.userGroupDelete({
      userGroupID: userGroup.userGroupID,
    })
    $toast.toastSuccess(t('notification.userGroup.delete.success', 'User Group deleted.'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete user group:', e)
    $toast.toastErrorHandler(
      t('notification.userGroup.delete.error', 'Failed to delete User Group'),
    )(e)
  }
}
</script>
