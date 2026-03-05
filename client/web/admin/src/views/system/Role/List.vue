<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.roles.list.title', 'Roles') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="roleID"
      :fields="roleListFields"
      :items="roleList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.roles.list.filterForm.query.placeholder', 'Filter roles'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.roles.list.new', 'New Role'),
        resourcePlural: $t('system.roles.list.title', 'Roles'),
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
            :label="$t('system.roles.list.new', 'New Role')"
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
const roleListFields = [
  {
    key: 'handle',
    sortable: true,
    header: t('system.roles.list.columns.handle', 'Handle'),
  },
  {
    key: 'name',
    sortable: true,
    header: t('system.roles.list.columns.name', 'Name'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.roles.editor.info.updatedAt', 'Updated at'),
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
function toggleActionsMenu(event, role) {
  actionsMenuItems.value = getActionsMenuItems(role)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(role) {
  const items = []

  if (role.canDeleteRole) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => onConfirmDelete(role),
    })
  }

  return items
}

function onConfirmDelete(role) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: role.name || role.handle || t('system.roles.list.new', 'New Role'),
    onConfirm: () => handleDelete(role),
  })
}

async function handleDelete(role) {
  try {
    await $SystemAPI.roleDelete({
      roleID: role.roleID,
    })
    $toast.toastSuccess(t('notification.role.delete.success', 'Role deleted.'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete role:', e)
    $toast.toastErrorHandler(t('notification.role.delete.error', 'Failed to delete Role'))(e)
  }
}
</script>
