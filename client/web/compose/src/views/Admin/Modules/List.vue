<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('module.navigation.module') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="moduleID"
      :fields="moduleFields"
      :items="moduleList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('module.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.label.module.single'),
        resourcePlural: $t('general.label.module.plural'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <CRouterLinkButton
            v-if="namespace?.canCreateModule"
            :to="{ name: 'admin.modules.create' }"
            :label="$t('module.createLabel')"
            icon="pi pi-plus"
            size="small"
          />
          <Button
            v-if="namespace?.canExportModules && moduleList.length"
            :label="$t('general.label.export')"
            icon="pi pi-download"
            size="small"
            severity="secondary"
            @click="exportAllModules"
          />
          <ModuleImporter
            v-if="namespace?.canCreateModule"
            :namespace="namespace"
            @imported="filterList"
          />
          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::compose:module/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <div class="flex items-center gap-2">
            <span>{{ data.name }}</span>
            <Tag v-if="isFederated(data)" severity="info" :value="$t('module.federated')" />
          </div>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
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
  usePermissions,
  useRBACStore,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import ModuleImporter from '@/components/Modules/ModuleImporter.vue'

const { CResourceList, CRouterLinkButton } = components
const { locFullDateTime } = filters

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('compose/', 'grant'))
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')

const resourceListRef = ref()

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
]

// Resource list composable
const {
  items: moduleList,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
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
function getActionsMenuItems(module) {
  const items = []

  if (module.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value.hideActionsMenu()
        openPermissions({
          resource: `corteza::compose:module/${props.namespace.namespaceID}/${module.moduleID}`,
          title: module.name || module.handle || module.moduleID,
          target: module.name || module.handle || module.moduleID,
        })
      },
    })
  }

  if (props.namespace?.canExportModules) {
    items.push({
      label: t('general.label.export'),
      icon: 'pi pi-download',
      command: () => exportModule(module),
    })
  }

  if (module.canDeleteModule) {
    if (items.length > 0) items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(module),
    })
  }

  return items
}

function onConfirmDelete(module) {
  confirmDelete({
    message: t('module.list.delete'),
    header: module.name,
    onConfirm: () => handleDelete(module),
  })
}

async function handleDelete(module) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $ComposeAPI.moduleDelete({
      namespaceID: props.namespace.namespaceID,
      moduleID: module.moduleID,
    })
    $toast.toastSuccess(t('notification.module.deleted'))
    filterList() // Refresh the list
  } catch (e) {
    console.error('Failed to delete module:', e)
    $toast.toastErrorHandler(t('notification.module.deleteFailed'))(e)
  }
}

function exportModule(module) {
  const blob = new Blob(
    [JSON.stringify({ type: 'module', list: [module] }, null, 2)],
    { type: 'application/json' },
  )
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${module.handle || module.name || 'module'}-export.json`
  a.click()
  URL.revokeObjectURL(url)
}

function exportAllModules() {
  const blob = new Blob(
    [JSON.stringify({ type: 'module', list: moduleList.value }, null, 2)],
    { type: 'application/json' },
  )
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'modules-export.json'
  a.click()
  URL.revokeObjectURL(url)
}
</script>
