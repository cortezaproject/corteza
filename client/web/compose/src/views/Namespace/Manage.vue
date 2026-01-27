<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('namespace.manage.title') }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <Button asChild v-slot="slotProps" size="small">
      <RouterLink :to="{ name: 'namespace.list' }" :class="slotProps.class">
        {{ $t('namespace.manage.list-view') }}
      </RouterLink>
    </Button>
  </Teleport>

  <div class="container mx-auto p-3 h-full">
    <CResourceList
      primary-key="namespaceID"
      :fields="namespaceFields"
      :items="namespaceList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
    >
      <template #header>
        <Button asChild v-slot="slotProps" size="large">
          <RouterLink :to="{ name: 'namespace.create' }" :class="slotProps.class">
            {{ $t('namespace.manage.toolbar.buttons.create') }}
          </RouterLink>
        </Button>
      </template>

      <template #body-changedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #body-actions="{ data }">
        <Button
          icon="pi pi-ellipsis-v"
          text
          rounded
          severity="secondary"
          @click="toggleActionsMenu($event, data)"
        />
      </template>
    </CResourceList>

    <TieredMenu ref="actionsMenu" :model="actionsMenuItems" popup />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { components, filters, useResourceList } from '@cortezaproject/corteza-vue-next'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const $ComposeAPI = inject('$ComposeAPI')
const confirm = useConfirm()
const toast = useToast()

// Actions menu
const actionsMenu = ref()
const actionsMenuItems = ref([])
const currentNamespace = ref(null)

const namespaceFields = [
  {
    key: 'name',
    sortable: true,
    header: t('namespace.manage.table.columns.name'),
  },
  {
    key: 'slug',
    sortable: true,
    header: t('namespace.manage.table.columns.slug'),
  },
  {
    key: 'enabled',
    header: t('namespace.manage.table.columns.enabled'),
  },
  {
    key: 'changedAt',
    sortable: true,
    header: t('namespace.manage.table.columns.changedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
  {
    key: 'actions',
    class: 'text-right',
    header: '',
  },
]

const {
  items: namespaceList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  filterList,
  handleRowClick,
} = useResourceList(params => $ComposeAPI.namespaceListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'name', sortDesc: false },
  pagination: { limit: 50 },
})

// Actions menu methods
const toggleActionsMenu = (event, namespace) => {
  currentNamespace.value = namespace
  actionsMenuItems.value = getActionsMenuItems(namespace)
  actionsMenu.value.toggle(event)
}

const getActionsMenuItems = namespace => {
  const items = []

  if (namespace.canDeleteNamespace) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      command: () => handleDelete(namespace),
    })
  }

  return items
}

const handleDelete = namespace => {
  confirm.require({
    message: t('namespace.manage.delete.confirm', {
      name: namespace.name || namespace.slug || namespace.namespaceID,
    }),
    header: t('general.label.delete'),
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: t('general.label.delete'),
    rejectLabel: t('general.label.cancel'),
    accept: () => {
      $ComposeAPI
        .namespaceDelete({ namespaceID: namespace.namespaceID })
        .then(() => {
          toast.add({
            severity: 'success',
            summary: t('general.notification.success'),
            detail: t('namespace.manage.delete.success'),
            life: 3000,
          })
          // Refresh the list
          filterList()
        })
        .catch(error => {
          toast.add({
            severity: 'error',
            summary: t('general.notification.error'),
            detail: error.message || t('namespace.manage.delete.error'),
            life: 5000,
          })
        })
    },
  })
}

// Lifecycle
onMounted(() => {
  document.title = t('general.label.app-name.namespace.list')
})
</script>
