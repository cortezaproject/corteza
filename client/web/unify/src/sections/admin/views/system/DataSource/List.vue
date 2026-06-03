<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.data-sources.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0 flex flex-col gap-4">
    <Panel toggleable :collapsed="false" class="shrink-0">
      <template #header>
        <div class="flex items-center gap-2 flex-1">
          <span class="font-medium">
            {{ $t('system.data-sources.primary.title') }}
          </span>
          <Button
            v-if="primaryConnection"
            v-tooltip.bottom="$t('general.label.edit')"
            icon="pi pi-pencil"
            text
            rounded
            size="small"
            @click="
              $router.push({
                name: 'system.dataSources.edit',
                params: { connectionID: primaryConnection.connectionID },
              })
            "
          />
        </div>
      </template>

      <div v-if="loadingPrimary" class="flex items-center justify-center py-4">
        <ProgressSpinner style="width: 2rem; height: 2rem" />
      </div>

      <div v-else-if="primaryConnection" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col">
          <span class="font-medium text-primary text-sm mb-1">
            {{ $t('system.data-sources.primary.name') }}
          </span>
          <span>{{ primaryConnection.meta?.name || '—' }}</span>
        </div>
        <div class="flex flex-col">
          <span class="font-medium text-primary text-sm mb-1">
            {{ $t('system.data-sources.primary.handle') }}
          </span>
          <span>{{ primaryConnection.handle || '—' }}</span>
        </div>
        <div class="flex flex-col">
          <span class="font-medium text-primary text-sm mb-1">
            {{ $t('system.data-sources.primary.location') }}
          </span>
          <span>{{ primaryLocationLabel || '—' }}</span>
        </div>
        <div class="flex flex-col">
          <span class="font-medium text-primary text-sm mb-1">
            {{ $t('system.data-sources.primary.ownership') }}
          </span>
          <span>{{ primaryConnection.meta?.ownership || '—' }}</span>
        </div>
      </div>

      <div v-else class="text-muted-color text-sm py-2">
        {{ $t('general.resourceList.noItems') }}
      </div>
    </Panel>

    <div class="flex-1 min-h-0">
      <CResourceList
        ref="resourceListRef"
        primary-key="connectionID"
        :fields="fields"
        :items="items"
        :filter="filter"
        :sorting="sorting"
        :pagination="pagination"
        :loading="loading"
        :action-items="getActionsMenuItems"
        :translations="{
          searchPlaceholder: $t('system.data-sources.list.filterForm.query.placeholder'),
          showingPagination: 'general.resourceList.pagination.showing',
          singlePluralPagination: 'general.resourceList.pagination.single',
          prevPagination: $t('general.resourceList.pagination.prev'),
          nextPagination: $t('general.resourceList.pagination.next'),
          recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
          resourceSingle: $t('system.data-sources.list.add-button'),
          resourcePlural: $t('system.data-sources.list.title'),
        }"
        clickable
        class="h-full"
        @update:filter="Object.assign(filter, $event)"
        @sort="handleSort"
        @row-click="
          ({ data }) =>
            $router.push({
              name: 'system.dataSources.edit',
              params: { connectionID: data.connectionID },
            })
        "
        @page-change="handlePageChange"
      >
        <template #header>
          <div class="flex gap-2">
            <Button
              :label="$t('system.data-sources.list.add-button')"
              icon="pi pi-plus"
              size="small"
              @click="$router.push({ name: 'system.dataSources.create' })"
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

        <template #body-type="{ data }">
          <Tag :value="data.type || 'corteza::system:dal-connection'" severity="info" />
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
    </div>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.data-sources.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.data-sources.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.data-sources.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.data-sources.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
} from '@planetcrust/human-vue'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const PRIMARY_TYPE = 'corteza::system:primary-dal-connection'
const EXTERNAL_TYPE = 'corteza::system:dal-connection'

const resourceListRef = ref()
const filterMenu = ref()

const primaryConnection = ref(null)
const loadingPrimary = ref(false)

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.data-sources.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('system.data-sources.list.columns.handle'),
  },
  {
    key: 'type',
    sortable: false,
    header: t('system.data-sources.list.columns.type'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('system.data-sources.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(
    params => $SystemAPI.dalConnectionListCancellable({ ...params, type: EXTERNAL_TYPE }),
    {
      filter: { query: '', deleted: '0' },
      sorting: { sortBy: 'name', sortDesc: false },
      pagination: { limit: 50 },
    },
  )

const primaryLocationLabel = computed(() => {
  if (!primaryConnection.value) return ''
  const loc = primaryConnection.value.meta?.location || {}
  const name = loc.properties?.name
  const coords = loc.geometry?.coordinates || []
  if (name) return name
  if (coords.length === 2) {
    return coords.map(c => c?.toFixed?.(7) || c).join(', ')
  }
  return ''
})

async function loadPrimary() {
  loadingPrimary.value = true
  try {
    const { set = [] } = await $SystemAPI.dalConnectionList({ type: PRIMARY_TYPE })
    primaryConnection.value = set.find(c => c.type === PRIMARY_TYPE) || null
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.fetch.error'))(e)
  } finally {
    loadingPrimary.value = false
  }
}

function getActionsMenuItems(item) {
  const items = []

  if (item.canDeleteConnection) {
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
    header: item.meta?.name || item.handle || item.connectionID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.dalConnectionDelete({ connectionID: item.connectionID })
    $toast.toastSuccess(t('notification.data-source.delete.success'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.delete.error'))(e)
  }
}

onMounted(loadPrimary)
</script>
