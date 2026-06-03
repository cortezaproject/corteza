<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chatbot.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="chatbotID"
      :fields="fields"
      :items="chatbots"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('chatbot.list.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('chatbot.list.title'),
        resourcePlural: $t('chatbot.list.title'),
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
          <CRouterLinkButton
            v-if="canCreate"
            :to="{ name: 'chatbot.create' }"
            :label="$t('chatbot.list.create')"
            icon="pi pi-plus"
            size="small"
          />
          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::system:chatbot/*"
            :title="$t('chatbot.list.title')"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.name || data.handle || '-' }}</span>
          <span v-if="data.handle && data.name" class="text-xs text-muted-color truncate">
            {{ data.handle }}
          </span>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('chatbot.list.enabled.on') : $t('chatbot.list.enabled.off')"
          :severity="data.enabled ? 'success' : 'secondary'"
          rounded
        />
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>
    </CResourceList>
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
import { useChatbotStore } from '@/sections/chatbot/stores/chatbot'

const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const canCreate = computed(() => rbac.can('system/', 'chatbot.create'))
const { open: openPermissions } = usePermissions()
const chatbotStore = useChatbotStore()

const resourceListRef = ref()

const fields = [
  { key: 'name', sortable: true, header: t('chatbot.list.columns.name') },
  { key: 'handle', sortable: true, header: t('chatbot.list.columns.handle') },
  { key: 'enabled', sortable: true, header: t('chatbot.list.columns.enabled') },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('chatbot.list.columns.updatedAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const {
  items: chatbots,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $SystemAPI.chatbotListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'name', sortDesc: false },
  pagination: { limit: 50 },
})

function handleRowClick({ data }) {
  if (!data.canUpdateChatbot && !data.canDeleteChatbot) return
  router.push({ name: 'chatbot.edit', params: { chatbotID: data.chatbotID } })
}

function getActionsMenuItems(cb) {
  const items = []
  if (cb.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::system:chatbot/${cb.chatbotID}`,
          title: cb.name || cb.handle || cb.chatbotID,
        })
      },
    })
  }
  if (cb.canUpdateChatbot) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      route: { name: 'chatbot.edit', params: { chatbotID: cb.chatbotID } },
    })
    items.push({
      label: t('chatbot.list.actions.duplicate'),
      icon: 'pi pi-copy',
      command: () => handleDuplicate(cb),
    })
  }
  if (cb.deletedAt) {
    items.push({
      label: t('chatbot.list.actions.undelete'),
      icon: 'pi pi-undo',
      command: () => handleUndelete(cb),
    })
  } else if (cb.canDeleteChatbot) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(cb),
    })
  }
  return items
}

function onConfirmDelete(cb) {
  confirmDelete({
    message: t('chatbot.list.delete'),
    header: cb.name || cb.handle || t('general.label.delete'),
    onConfirm: () => handleDelete(cb),
  })
}

async function handleDelete(cb) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.chatbotDelete({ chatbotID: cb.chatbotID })
    chatbotStore.removeFromList(cb.chatbotID)
    $toast.toastSuccess(t('notification.chatbot.deleted'))
    filterList()
  } catch (e) {
    console.error(e)
    $toast.toastErrorHandler(t('notification.chatbot.deleteFailed'))(e)
  }
}

async function handleDuplicate(cb) {
  resourceListRef.value.hideActionsMenu()
  try {
    const source = await $SystemAPI.chatbotRead({ chatbotID: cb.chatbotID })
    const copy = {
      handle: source.handle ? `${source.handle}_copy` : '',
      name: source.name ? t('chatbot.list.duplicateNameSuffix', { name: source.name }) : '',
      enabled: false,
      sessionTTL: source.sessionTTL || '',
      allowedOrigins: source.allowedOrigins || [],
      handoff: source.handoff || {},
      styling: source.styling || {},
      scenarios: source.scenarios || [],
      labels: source.labels || {},
    }
    const created = await $SystemAPI.chatbotCreate(copy)
    chatbotStore.updateInList(created)
    $toast.toastSuccess(t('notification.chatbot.duplicated'))
    router.push({ name: 'chatbot.edit', params: { chatbotID: created.chatbotID } })
  } catch (e) {
    console.error('Failed to duplicate chatbot:', e)
    $toast.toastErrorHandler(t('notification.chatbot.duplicateFailed'))(e)
  }
}

async function handleUndelete(cb) {
  resourceListRef.value.hideActionsMenu()
  try {
    const restored = await $SystemAPI.chatbotUndelete({ chatbotID: cb.chatbotID })
    chatbotStore.updateInList(restored || { ...cb, deletedAt: null })
    $toast.toastSuccess(t('notification.chatbot.restored'))
    filterList()
  } catch (e) {
    console.error(e)
    $toast.toastErrorHandler(t('notification.chatbot.restoreFailed'))(e)
  }
}
</script>
