<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.workflows.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="workflowID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('automation.workflows.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('automation.workflows.list.new'),
        resourcePlural: $t('automation.workflows.list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="
        ({ data }) =>
          $router.push({
            name: 'automation.workflows.edit',
            params: { workflowID: data.workflowID },
          })
      "
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('automation.workflows.list.new')"
            icon="pi pi-plus"
            size="small"
            @click="$router.push({ name: 'automation.workflows.create' })"
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

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.name || '—' }}</span>
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
        {{ locFullDateTime(data.createdAt) }}
      </template>

    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('automation.workflows.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('automation.workflows.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('automation.workflows.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('automation.workflows.list.filterForm.exclusive.label') }}
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
const $AutomationAPI = inject('$AutomationAPI')

const resourceListRef = ref()
const filterMenu = ref()

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  {
    key: 'name',
    sortable: false,
    header: t('automation.workflows.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('automation.workflows.list.columns.handle'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('automation.workflows.list.columns.enabled'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('automation.workflows.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(params => $AutomationAPI.workflowListCancellable({ ...params }), {
    filter: { query: '', deleted: '0' },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 50 },
  })

function getActionsMenuItems(item) {
  const items = []

  if (item.canDeleteWorkflow) {
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
    header: item.meta?.name || item.handle || item.workflowID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $AutomationAPI.workflowDelete({ workflowID: item.workflowID })
    $toast.toastSuccess(t('notification.workflow.delete.success'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.workflow.delete.error'))(e)
  }
}
</script>
