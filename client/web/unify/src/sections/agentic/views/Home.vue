<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('agent.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="agentID"
      :fields="agentFields"
      :items="agentList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
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
        <div class="flex gap-2">
          <CRouterLinkButton
            v-if="canCreate"
            :to="{ name: 'agentic.create' }"
            :label="$t('agent.list.create')"
            icon="pi pi-plus"
            size="small"
          />
          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:agent/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.short || '-' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-md">
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

      <template #body-changedAt="{ data }">
        {{ changedAtText(data) }}
      </template>
    </CResourceList>
  </div>
</template>

<script setup>
import {
  changedAtText,
  changedAtField,
  components,
  useConfirmDelete,
  usePermissions,
  useRBACStore,
  useResourceList,
} from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAgentStore } from '@planetcrust/human-vue'

const { CResourceList, CRouterLinkButton } = components

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const agentStore = useAgentStore()
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const canCreate = computed(() => rbac.can('system/', 'agent.create'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()

// Column definitions
const agentFields = [
  {
    key: 'name',
    sortable: true,
    header: t('agent.list.columns.name'),
  },
  {
    key: 'status',
    sortable: true,
    header: t('agent.list.columns.status'),
  },
  changedAtField(t('general.columns.changedAt')),
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
  pagination: { limit: 50 },
})

// Methods
function handleRowClick({ data }) {
  if (!data.canUpdateAgent && !data.canDeleteAgent) {
    return
  }
  router.push({
    name: 'agentic.edit',
    params: { agentID: data.agentID },
  })
}

// Actions menu methods
function getActionsMenuItems(agent) {
  const items = []

  if (agent.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::system:agent/${agent.agentID}`,
          title: agent.meta?.short || agent.handle || agent.agentID,
        })
      },
    })
  }

  if (agent.canUpdateAgent) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      route: {
        name: 'agentic.edit',
        params: { agentID: agent.agentID },
      },
    })
  }

  // Duplicate action
  if (agent.canUpdateAgent) {
    items.push({
      label: t('agent.list.actions.duplicate'),
      icon: 'pi pi-copy',
      command: () => handleDuplicate(agent),
    })
  }

  // Show either delete or restore depending on state. The server authorises an
  // undelete with the delete permission, so both answer to the same flag.
  if (agent.deletedAt) {
    if (agent.canDeleteAgent) {
      items.push({
        label: t('general.label.restore'),
        icon: 'pi pi-replay',
        command: () => handleRestore(agent),
      })
    }
  } else if (agent.canDeleteAgent) {
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
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.agentDelete({ agentID: agent.agentID })
    agentStore.removeFromList(agent.agentID)
    $toast.toastSuccess(t('notification.agent.deleted'))
    filterList()
  } catch (e) {
    console.error('Failed to delete agent:', e)
    $toast.toastErrorHandler(t('notification.agent.deleteFailed'))(e)
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
    agentStore.updateInList(created)
    $toast.toastSuccess(t('notification.agent.duplicated'))

    // Navigate to the new agent
    router.push({
      name: 'agentic.edit',
      params: { agentID: created.agentID },
    })
  } catch (e) {
    console.error('Failed to duplicate agent:', e)
    $toast.toastErrorHandler(t('notification.agent.duplicateFailed'))(e)
  }
}

async function handleRestore(agent) {
  resourceListRef.value.hideActionsMenu()
  try {
    // Undelete answers with a bare OK, so the row that comes back into the
    // sidebar is the one the list is holding, with its deletion cleared.
    await $SystemAPI.agentUndelete({ agentID: agent.agentID })
    agentStore.updateInList({ ...agent, deletedAt: null })
    $toast.toastSuccess(t('notification.agent.restored'))
    filterList()
  } catch (e) {
    console.error('Failed to restore agent:', e)
    $toast.toastErrorHandler(t('notification.agent.restoreFailed'))(e)
  }
}
</script>
