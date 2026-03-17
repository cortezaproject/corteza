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
        <div class="flex flex-col">
          <span class="font-medium">{{ data.meta?.short || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-status="{ data }">
        <Tag
          :value="$t(`agent.list.status.${data.status || 'inactive'}`)"
          :severity="data.status === 'active' ? 'success' : 'secondary'"
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
import Tag from 'primevue/tag'

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

  // Duplicate action
  if (agent.canUpdateAgent !== false) {
    items.push({
      label: t('agent.list.actions.duplicate'),
      icon: 'pi pi-copy',
      command: () => handleDuplicate(agent),
    })
  }

  // Show either delete or undelete depending on state
  if (agent.deletedAt) {
    items.push({
      label: t('agent.list.actions.undelete'),
      icon: 'pi pi-undo',
      command: () => handleUndelete(agent),
    })
  } else if (agent.canDeleteAgent !== false) {
    if (items.length > 0) {
      items.push({ separator: true })
    }
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(agent),
    })
  }

  return items
}

function onConfirmDelete(agent) {
  confirmDelete({
    message: t('agent.list.delete'),
    header: agent.meta?.short || agent.handle || t('general.label.delete'),
    onConfirm: () => handleDelete(agent),
  })
}

async function handleDelete(agent) {
  try {
    await $SystemAPI.agentDelete({
      agentID: agent.agentID,
    })
    $toast.toastSuccess(t('notification.agent.deleted'))
    filterList()
  } catch (e) {
    console.error('Failed to delete agent:', e)
    $toast.toastDanger(t('notification.agent.deleteFailed'))
  }
}

async function handleDuplicate(agent) {
  try {
    // Read the full agent to get all config
    const source = await $SystemAPI.agentRead({ agentID: agent.agentID })

    // Create a copy with modified name
    const copy = {
      handle: source.handle ? `${source.handle}_copy` : '',
      status: 'inactive',
      meta: {
        ...(source.meta || {}),
        short: `${source.meta?.short || ''} (Copy)`,
      },
      behavior: source.behavior || {},
      execution: source.execution || {},
      access: source.access || {},
      invocation: source.invocation || {},
    }

    const created = await $SystemAPI.agentCreate(copy)
    $toast.toastSuccess(t('notification.agent.duplicated'))

    // Navigate to the new agent
    router.push({
      name: 'agent.edit',
      params: { agentID: created.agentID },
    })
  } catch (e) {
    console.error('Failed to duplicate agent:', e)
    $toast.toastDanger(t('notification.agent.duplicateFailed'))
  }
}

async function handleUndelete(agent) {
  try {
    await $SystemAPI.agentUndelete({
      agentID: agent.agentID,
    })
    $toast.toastSuccess(t('notification.agent.restored'))
    filterList()
  } catch (e) {
    console.error('Failed to restore agent:', e)
    $toast.toastDanger(t('notification.agent.restoreFailed'))
  }
}
</script>
