<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="automationID"
      :fields="listFields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('list.new'),
        resourcePlural: $t('list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('list.new')"
            icon="pi pi-plus"
            size="small"
            @click="openCreateDialog"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::automation:ng-automation/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.short || $t('list.untitled') }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('list.active') : $t('list.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
        />
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
            {{ $t('list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('list.filterForm.disabled.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis0" value="0" />
            <label for="dis0" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis1" value="1" />
            <label for="dis1" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis2" value="2" />
            <label for="dis2" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>

    <!-- Create TAQ Dialog -->
    <TaqConfigModal v-model:visible="showCreateDialog" mode="create" />
  </div>
</template>

<script setup>
import {
  components,
  filters,
  useConfirmDelete,
  useRBACStore,
  useResourceList,
  usePermissions,
} from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import TaqConfigModal from '../components/common/TaqConfigModal.vue'
import { useAutomationStore } from '@planetcrust/human-vue'

const { CResourceList } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const automationStore = useAutomationStore()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $AutomationAPI = inject('$AutomationAPI')
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('automation/', 'grant'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const listFields = [
  {
    key: 'name',
    sortable: true,
    header: t('list.columns.name'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('list.columns.enabled'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('list.columns.updatedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(params => $AutomationAPI.ngAutomationListCancellable(params), {
    filter: { query: '', deleted: '0', disabled: '1' },
    sorting: { sortBy: 'name', sortDesc: false },
    pagination: { limit: 50 },
  })

// Create dialog state
const showCreateDialog = ref(false)

function openCreateDialog() {
  showCreateDialog.value = true
}

// Row click navigation
function handleRowClick({ data }) {
  router.push(`/taq/builder/${data.automationID}`)
}

// Actions menu methods
function getActionsMenuItems(automation) {
  const items = []

  if (automation.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::automation:ng-automation/${automation.automationID}`,
          title: automation.meta?.short || automation.automationID,
        })
      },
    })
  }

  if (automation.canDeleteNgAutomation) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(automation),
    })
  }

  return items
}

function onConfirmDelete(automation) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: automation.meta?.short || t('list.untitled'),
    onConfirm: () => handleDelete(automation),
  })
}

async function handleDelete(automation) {
  resourceListRef.value.hideActionsMenu()
  try {
    await automationStore.remove(automation.automationID)
    $toast.toastSuccess(t('notification.automation.delete.success'))
    filterList()
  } catch (e) {
    console.error('Failed to delete automation:', e)
    $toast.toastErrorHandler(t('notification.automation.delete.error'))(e)
  }
}
</script>
