<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('module.navigation.module') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      primary-key="moduleID"
      :fields="moduleFields"
      :items="moduleList"
      :filter="filter"
      @update:filter="Object.assign(filter, $event)"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
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
        <CRouterLinkButton
          v-if="namespace?.canCreateModule"
          :to="{ name: 'admin.modules.create' }"
          :label="$t('module.createLabel')"
          icon="pi pi-plus"
          size="small"
        />
      </template>

      <template #body-name="{ data }">
        <div class="flex items-center gap-2">
          <span class="font-medium">{{ data.name }}</span>
          <Tag v-if="isFederated(data)" severity="info" :value="$t('module.federated')" />
        </div>
      </template>

      <template #body-handle="{ data }">
        {{ data.handle || '-' }}
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
          class="row-action-btn w-full"
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

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
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
    pt: {
      headerCell: { class: 'border-l-0' },
      bodyCell: { class: 'p-0 border-l-0' },
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
function toggleActionsMenu(event, module) {
  actionsMenuItems.value = getActionsMenuItems(module)
  actionsMenu.value.toggle(event)
}

function getActionsMenuItems(module) {
  const items = []

  if (module.canDeleteModule) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
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
