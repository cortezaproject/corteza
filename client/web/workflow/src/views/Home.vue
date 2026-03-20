<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('general.workflow-list') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="workflowID"
      :fields="workflowFields"
      :items="workflowList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('general.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.workflow'),
        resourcePlural: $t('general.workflows'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <CRouterLinkButton
          v-if="canCreate"
          :to="{ name: 'workflow.create' }"
          :label="$t('general.new-workflow')"
          icon="pi pi-plus"
          size="small"
        />
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.name || data.handle || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('general.enabled') : $t('general.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
          rounded
        />
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
          class="row-action-btn w-full mr-2"
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
const $AutomationAPI = inject('$AutomationAPI')

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])

// RBAC
const canCreate = ref(true)

// Column definitions
const workflowFields = [
  {
    key: 'name',
    sortable: false,
    header: t('general.columns.name'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('general.columns.enabled'),
    class: 'text-center w-28',
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('general.columns.changedAt'),
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
  items: workflowList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $AutomationAPI.workflowListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  router.push({
    name: 'workflow.edit',
    params: { workflowID: data.workflowID },
  })
}

function toggleActionsMenu(event, workflow) {
  actionsMenuItems.value = getActionsMenuItems(workflow)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(workflow) {
  const items = []

  items.push({
    label: t('general.label.edit'),
    icon: 'pi pi-pencil',
    route: {
      name: 'workflow.edit',
      params: { workflowID: workflow.workflowID },
    },
  })

  if (workflow.deletedAt) {
    items.push({
      label: t('general.undelete'),
      icon: 'pi pi-undo',
      command: () => handleUndelete(workflow),
    })
  } else {
    if (items.length > 0) {
      items.push({ separator: true })
    }
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(workflow),
    })
  }

  return items
}

function onConfirmDelete(workflow) {
  confirmDelete({
    message: t('notification.delete.confirm'),
    header: workflow.meta?.name || workflow.handle || t('general.label.delete'),
    onConfirm: () => handleDelete(workflow),
  })
}

async function handleDelete(workflow) {
  try {
    await $AutomationAPI.workflowDelete({
      workflowID: workflow.workflowID,
    })
    $toast.toastSuccess(t('notification.delete.success'))
    filterList()
  } catch (e) {
    console.error('Failed to delete workflow:', e)
    $toast.toastDanger(t('notification.delete.failed'))
  }
}

async function handleUndelete(workflow) {
  try {
    await $AutomationAPI.workflowUndelete({
      workflowID: workflow.workflowID,
    })
    $toast.toastSuccess(t('notification.undelete.success'))
    filterList()
  } catch (e) {
    console.error('Failed to restore workflow:', e)
    $toast.toastDanger(t('notification.undelete.failed'))
  }
}
</script>
