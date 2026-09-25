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
        resourceSingle: $t('general.label.chatbot.single'),
        resourcePlural: $t('general.label.chatbot.plural'),
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
          <span v-if="data.handle && data.name" class="text-xs text-muted-color truncate max-w-md">
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

      <template #body-changedAt="{ data }">
        {{ changedAtText(data) }}
      </template>
      <template #filter>
        <Button
          :label="$t('general.filter.label')"
          icon="pi pi-filter"
          severity="secondary"
          size="small"
          outlined
          @click="toggleFilterMenu"
        />
      </template>
    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-2 p-2 w-64">
        <span class="font-medium text-sm text-primary">
          {{ $t('chatbot.list.filterForm.deleted.label') }}
        </span>
        <div class="flex items-center gap-2">
          <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
          <label for="del0" class="text-sm cursor-pointer">
            {{ $t('chatbot.list.filterForm.excluded.label') }}
          </label>
        </div>
        <div class="flex items-center gap-2">
          <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
          <label for="del1" class="text-sm cursor-pointer">
            {{ $t('chatbot.list.filterForm.inclusive.label') }}
          </label>
        </div>
        <div class="flex items-center gap-2">
          <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
          <label for="del2" class="text-sm cursor-pointer">
            {{ $t('chatbot.list.filterForm.exclusive.label') }}
          </label>
        </div>
      </div>
    </Popover>
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
import { useChatbotStore } from '@planetcrust/human-vue'

const { CResourceList, CRouterLinkButton } = components

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

const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  { key: 'name', sortable: true, header: t('chatbot.list.columns.name') },
  { key: 'handle', sortable: true, header: t('chatbot.list.columns.handle') },
  { key: 'enabled', sortable: true, header: t('chatbot.list.columns.enabled') },
  changedAtField(t('general.columns.changedAt')),
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
  filter: { query: '', deleted: '0' },
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
    if (cb.canDeleteChatbot) {
      items.push({
        label: t('general.label.restore'),
        icon: 'pi pi-replay',
        command: () => handleRestore(cb),
      })
    }
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

async function handleRestore(cb) {
  resourceListRef.value.hideActionsMenu()
  try {
    // Undelete answers with a bare OK, so the row that comes back into the
    // sidebar is the one the list is holding, with its deletion cleared.
    await $SystemAPI.chatbotUndelete({ chatbotID: cb.chatbotID })
    chatbotStore.updateInList({ ...cb, deletedAt: null })
    $toast.toastSuccess(t('notification.chatbot.restored'))
    filterList()
  } catch (e) {
    console.error(e)
    $toast.toastErrorHandler(t('notification.chatbot.restoreFailed'))(e)
  }
}
</script>
