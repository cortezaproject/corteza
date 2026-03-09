<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.sensitivityLevel.list.title', 'Sensitivity Levels') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="sensitivityLevelID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t(
          'system.sensitivityLevel.list.filterForm.query.placeholder',
          'Filter sensitivity levels',
        ),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.sensitivityLevel.list.new'),
        resourcePlural: $t('system.sensitivityLevel.list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="
        ({ data }) =>
          $router.push({
            name: 'system.sensitivityLevels.edit',
            params: { sensitivityLevelID: data.sensitivityLevelID },
          })
      "
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('system.sensitivityLevel.list.new')"
            icon="pi pi-plus"
            size="small"
            @click="$router.push({ name: 'system.sensitivityLevels.create' })"
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

      <template #body-createdAt="{ data }">
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
        <a v-ripple v-bind="props.action" :class="item.class">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </template>
    </TieredMenu>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.sensitivityLevel.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.sensitivityLevel.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.sensitivityLevel.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.sensitivityLevel.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const actionsMenu = ref()
const actionsMenuItems = ref([])
const filterMenu = ref()

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.sensitivityLevel.list.columns.name', 'Name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.sensitivityLevel.list.columns.handle', 'Handle'),
  },
  {
    key: 'level',
    sortable: true,
    header: t('system.sensitivityLevel.list.columns.level', 'Level'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('system.sensitivityLevel.list.columns.createdAt', 'Created'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
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

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(params => $SystemAPI.dalSensitivityLevelListCancellable({ ...params }), {
    filter: { query: '', deleted: '0' },
    sorting: { sortBy: 'level', sortDesc: false },
    pagination: { limit: 50 },
  })

function toggleActionsMenu(event, item) {
  actionsMenuItems.value = getActionsMenuItems(item)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(item) {
  const items = []

  if (item.canDeleteDalSensitivityLevel) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(item),
    })
  }

  return items
}

function onConfirmDelete(item) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: item.name || item.handle || item.sensitivityLevelID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  try {
    await $SystemAPI.dalSensitivityLevelDelete({ sensitivityLevelID: item.sensitivityLevelID })
    $toast.toastSuccess(
      t('notification.sensitivityLevel.delete.success', 'Sensitivity level deleted'),
    )
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(
      t('notification.sensitivityLevel.delete.error', 'Failed to delete sensitivity level'),
    )(e)
  }
}
</script>
