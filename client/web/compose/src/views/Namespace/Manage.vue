<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('namespace.manage.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="namespaceID"
      :fields="namespaceFields"
      :items="namespaceList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.label.namespace.single'),
        resourcePlural: $t('general.label.namespace.plural'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <CRouterLinkButton
            :to="{ name: 'namespace.create' }"
            :label="$t('namespace.manage.toolbar.buttons.create')"
            icon="pi pi-plus"
            size="small"
          />
          <NamespaceImporter @imported="onImported" @failed="onFailed" />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.name || '—' }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-changedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #body-actions="{ data }">
        <Button
          icon="pi pi-ellipsis-v"
          text
          severity="secondary"
          size="small"
          class="row-action-btn w-full"
          @click="toggleActionsMenu($event, data)"
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
import NamespaceImporter from '@/components/Namespaces/NamespaceImporter.vue'
import { useNamespaceStore } from '@/stores/namespace'
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const router = useRouter()
const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()
const namespaceStore = useNamespaceStore()

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
    class: 'text-right w-12',
    header: '',
    frozen: true,
    alignFrozen: 'right',
  },
]

const {
  items: namespaceList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $ComposeAPI.namespaceListCancellable(params), {
  filter: { query: '' },
  sorting: { sortBy: 'name', sortDesc: false },
  pagination: { limit: 50 },
})

function handleRowClick({ data }) {
  if (!(data.canUpdateNamespace || data.canDeleteNamespace)) return
  router.push({
    name: 'namespace.edit',
    params: { slug: data.slug || data.namespaceID },
  })
}

// Actions menu methods
const toggleActionsMenu = (event, namespace) => {
  currentNamespace.value = namespace
  actionsMenuItems.value = getActionsMenuItems(namespace)
  actionsMenu.value.toggle(event)
}

const getActionsMenuItems = namespace => {
  const items = []

  if (namespace.canExportNamespace) {
    items.push({
      label: t('namespace.export'),
      icon: 'pi pi-download',
      command: () => exportNamespace(namespace),
    })
  }

  if (namespace.canDeleteNamespace) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => handleDelete(namespace),
    })
  }

  return items
}

const handleDelete = namespace => {
  confirmDelete({
    message: t('namespace.manage.delete.confirm', {
      name: namespace.name || namespace.slug || namespace.namespaceID,
    }),
    header: t('general.label.delete'),
    onConfirm: () => {
      $ComposeAPI
        .namespaceDelete({ namespaceID: namespace.namespaceID })
        .then(() => {
          $toast.toastSuccess(t('namespace.manage.delete.success'))
          filterList()
        })
        .catch(error => {
          $toast.toastDanger(error.message || t('namespace.manage.delete.error'))
        })
    },
  })
}

// Export namespace
function exportNamespace(namespace) {
  const params = {
    namespaceID: namespace.namespaceID,
    filename: encodeURIComponent((namespace.name || 'namespace').replace(/\./g, '-')),
  }

  const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''
  const exportUrl = `${$ComposeAPI.baseURL}${$ComposeAPI.namespaceExportEndpoint(params)}?jwt=${encodeURIComponent(token)}`
  window.open(exportUrl)
}

// Import handlers
function onImported() {
  namespaceStore
    .load({ force: true })
    .then(() => {
      filterList()
      $toast.toastSuccess(t('notification.namespace.imported'))
    })
    .catch(() => {
      $toast.toastDanger(t('notification.namespace.importFailed'))
    })
}

function onFailed(err) {
  $toast.toastDanger(err?.message || t('notification.namespace.importFailed'))
}

// Lifecycle
onMounted(() => {
  document.title = t('general.label.app-name.namespace.list')
})
</script>
