<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('module.navigation.module') }}</span>
  </Teleport>

  <div class="container mx-auto p-5 h-full overflow-hidden">
    <CResourceList
      primary-key="moduleID"
      :fields="moduleFields"
      :items="moduleList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('module.searchPlaceholder'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
    >
      <template #header>
        <Button
          v-if="namespace?.canCreateModule"
          :label="$t('module.createLabel')"
          icon="pi pi-plus"
          @click="$router.push({ name: 'admin.modules.create' })"
        />
      </template>

      <template #body-name="{ data }">
        <div class="flex items-center gap-2">
          <span class="font-medium">{{ data.name }}</span>
          <Tag v-if="isFederated(data)" severity="info" :value="$t('module.federated')" />
        </div>
      </template>

      <template #body-handle="{ data }">
        <code class="text-sm">{{ data.handle || '-' }}</code>
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
          @click.stop="toggleActionsMenu($event, data)"
        />
      </template>
    </CResourceList>

    <TieredMenu ref="actionsMenu" :model="actionsMenuItems" popup />
  </div>
</template>

<script setup>
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'
import { useConfirm } from 'primevue/useconfirm'
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList } = components
const { locFullDateTime } = filters

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const { t } = useI18n()
const confirm = useConfirm()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])

// Column definitions
const moduleFields = [
  {
    key: 'name',
    sortable: true,
    header: t('module.list.columns.name'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('module.list.columns.handle'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('module.list.columns.changedAt'),
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
  },
]

// Resource list composable
const {
  items: moduleList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  filterList,
} = useResourceList(
  params =>
    $ComposeAPI.moduleListCancellable({
      namespaceID: props.namespace.namespaceID,
      ...params,
    }),
  {
    filter: { query: '' },
    sorting: { sortBy: 'name', sortDesc: false },
    pagination: { limit: 50 },
  },
)

// Methods
function isFederated(module) {
  return Object.keys(module.labels || {}).includes('federation')
}

function handleRowClick({ data }) {
  if (!(data.canUpdateModule || data.canDeleteModule)) {
    return
  }
  router.push({
    name: 'admin.modules.edit',
    params: { moduleID: data.moduleID },
  })
}

// Actions menu methods
function toggleActionsMenu(event, module) {
  actionsMenuItems.value = getActionsMenuItems(module)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(module) {
  const items = []

  if (module.canUpdateModule) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      command: () =>
        router.push({
          name: 'admin.modules.edit',
          params: { moduleID: module.moduleID },
        }),
    })
  }

  if (module.canDeleteModule) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => confirmDelete(module),
    })
  }

  return items
}

function confirmDelete(module) {
  confirm.require({
    message: t('module.list.delete'),
    header: module.name,
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: () => handleDelete(module),
  })
}

async function handleDelete(module) {
  try {
    await $ComposeAPI.moduleDelete({
      namespaceID: props.namespace.namespaceID,
      moduleID: module.moduleID,
    })
    $toast.toastSuccess(t('notification.module.deleted'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete module:', e)
    $toast.toastDanger(t('notification.module.deleteFailed'))
  }
}
</script>
