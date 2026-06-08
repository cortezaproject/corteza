<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.applications.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="applicationID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('system.applications.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.applications.list.new'),
        resourcePlural: $t('system.applications.list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="
        ({ data }) =>
          $router.push({
            name: 'system.applications.edit',
            params: { applicationID: data.applicationID },
          })
      "
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('system.applications.list.new')"
            icon="pi pi-plus"
            size="small"
            @click="$router.push({ name: 'system.applications.create' })"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:application/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.name || '—' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('general.label.enabled') : $t('general.label.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
        />
      </template>

      <template #body-createdAt="{ data }">
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
            {{ $t('system.applications.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.applications.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.applications.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.applications.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  components,
  filters,
  useConfirmDelete,
  usePermissions,
  useRBACStore,
  useResourceList,
  useApplicationsStore,
} from '@planetcrust/human-vue'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const applicationsStore = useApplicationsStore()

const resourceListRef = ref()
const filterMenu = ref()

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.applications.list.columns.name'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('system.applications.list.columns.enabled'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('system.applications.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(params => $SystemAPI.applicationListCancellable({ ...params }), {
    filter: { query: '', deleted: '0' },
    sorting: { sortBy: 'name', sortDesc: false },
    pagination: { limit: 50 },
  })

function getActionsMenuItems(item) {
  const menuItems = []

  if (item.canGrant) {
    menuItems.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::system:application/${item.applicationID}`,
          title: item.name || item.applicationID,
          target: item.name || item.applicationID,
        })
      },
    })
  }

  if (item.canDeleteApplication) {
    if (menuItems.length > 0) menuItems.push({ separator: true })
    menuItems.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(item),
    })
  }

  return menuItems
}

function onConfirmDelete(item) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: item.name || item.applicationID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  resourceListRef.value.hideActionsMenu()
  try {
    await applicationsStore.delete(item.applicationID)
    $toast.toastSuccess(t('notification.application.delete.success'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.application.delete.error'))(e)
  }
}
</script>
