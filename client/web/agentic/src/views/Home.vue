<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('agent.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="agentID"
      :fields="agentFields"
      :items="agentList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('agent.list.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('agent.list.title'),
        resourcePlural: $t('agent.list.title'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <CRouterLinkButton
          :to="{ name: 'agent.create' }"
          :label="$t('agent.list.create')"
          icon="pi pi-plus"
          size="small"
        />
      </template>

      <template #body-name="{ data }">
        <span class="font-medium">{{ data.meta?.short || '-' }}</span>
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
const $SystemAPI = inject('$SystemAPI')

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])

// Column definitions
const agentFields = [
  {
    key: 'name',
    sortable: false,
    header: t('agent.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('agent.list.columns.handle'),
  },
  {
    key: 'status',
    sortable: true,
    header: t('agent.list.columns.status'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('agent.list.columns.updatedAt'),
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
  items: agentList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.agentListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (data.canUpdateAgent === false && data.canDeleteAgent === false) {
    return
  }
  router.push({
    name: 'agent.edit',
    params: { agentID: data.agentID },
  })
}

// Actions menu methods
function toggleActionsMenu(event, agent) {
  actionsMenuItems.value = getActionsMenuItems(agent)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(agent) {
  const items = []

  if (agent.canUpdateAgent !== false) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      route: {
        name: 'agent.edit',
        params: { agentID: agent.agentID },
      },
    })
  }

  if (agent.canDeleteAgent !== false) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => onConfirmDelete(agent),
    })
  }

  return items
}

function onConfirmDelete(agent) {
  confirmDelete({
    message: t('agent.list.delete'),
    header: agent.meta?.short || agent.handle,
    onConfirm: () => handleDelete(agent),
  })
}

async function handleDelete(agent) {
  try {
    await $SystemAPI.agentDelete({
      agentID: agent.agentID,
    })
    $toast.toastSuccess(t('notification.agent.deleted'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete agent:', e)
    $toast.toastDanger(t('notification.agent.deleteFailed'))
  }
}
</script>
