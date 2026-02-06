<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('page.navigation.page') }}</span>
  </Teleport>

  <div class="container mx-auto p-3 h-full overflow-hidden">
    <CResourceList
      primary-key="pageID"
      :fields="pageFields"
      :items="pageList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('page.searchPlaceholder'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
    >
      <template #header>
        <Button
          v-if="namespace?.canCreatePage"
          :label="$t('page.createLabel')"
          icon="pi pi-plus"
          @click="$router.push({ name: 'admin.pages.create' })"
        />
      </template>

      <template #body-title="{ data }">
        <span class="font-medium">{{ data.title }}</span>
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
const pageFields = [
  {
    key: 'title',
    sortable: true,
    header: t('page.list.columns.title'),
  },
  {
    key: 'handle',
    sortable: true,
    header: t('page.list.columns.handle'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('page.list.columns.changedAt'),
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
  items: pageList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  filterList,
} = useResourceList(
  params =>
    $ComposeAPI.pageListCancellable({
      namespaceID: props.namespace.namespaceID,
      ...params,
    }),
  {
    filter: { query: '' },
    sorting: { sortBy: 'title', sortDesc: false },
    pagination: { limit: 50 },
  },
)

// Methods
function handleRowClick({ data }) {
  if (!(data.canUpdatePage || data.canDeletePage)) {
    return
  }
  router.push({
    name: 'admin.pages.edit',
    params: { pageID: data.pageID },
  })
}

// Actions menu methods
function toggleActionsMenu(event, page) {
  actionsMenuItems.value = getActionsMenuItems(page)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(page) {
  const items = []

  if (page.canUpdatePage) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      command: () =>
        router.push({
          name: 'admin.pages.edit',
          params: { pageID: page.pageID },
        }),
    })
  }

  if (page.canDeletePage) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => confirmDelete(page),
    })
  }

  return items
}

function confirmDelete(page) {
  confirm.require({
    message: t('page.list.delete'),
    header: page.title,
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: () => handleDelete(page),
  })
}

async function handleDelete(page) {
  try {
    await $ComposeAPI.pageDelete({
      namespaceID: props.namespace.namespaceID,
      pageID: page.pageID,
    })
    $toast.toastSuccess(t('notification.page.deleted'))
    filterList()
  } catch (e) {
    console.error('Failed to delete page:', e)
    $toast.toastDanger(t('notification.page.deleteFailed'))
  }
}
</script>
