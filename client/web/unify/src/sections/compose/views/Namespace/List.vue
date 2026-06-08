<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('namespace.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <div v-if="viewMode === 'cards'" class="h-full flex flex-col min-w-0 overflow-hidden gap-4">
      <Card
        :pt="{
          body: { class: 'p-3' },
          content: { class: 'flex flex-wrap items-center justify-between gap-3' },
        }"
        class="shrink-0"
      >
        <template #content>
          <div class="flex-1 min-w-0 flex items-center gap-2">
            <CRouterLinkButton
              v-if="canCreate"
              :to="{ name: 'namespace.create' }"
              :label="$t('namespace.manage.toolbar.buttons.create')"
              icon="pi pi-plus"
              size="small"
            />
            <NamespaceImporter v-if="canCreate" @imported="onImported" @failed="onFailed" />
            <CPermissionsButton
              v-if="canGrant"
              v-tooltip.bottom="$t('general.label.permissions')"
              resource="corteza::compose:namespace/*"
            />
          </div>
          <div class="flex-1 flex items-center justify-end gap-2">
            <SelectButton
              v-model="viewMode"
              :options="viewModeOptions"
              option-label="label"
              option-value="value"
              :allow-empty="false"
              size="small"
            >
              <template #option="{ option }">
                <i :class="option.icon" v-tooltip.bottom="option.label" />
              </template>
            </SelectButton>
            <CInputSearch
              v-model="query"
              :placeholder="$t('namespace.searchPlaceholder')"
              size="small"
              class="flex-1 min-w-0 max-w-xl"
            />
          </div>
        </template>
      </Card>

      <div v-if="areNamespacesVisible" class="flex-1 overflow-auto min-h-0 min-w-0">
        <div class="py-2">
          <div class="flex flex-wrap justify-center gap-7">
            <RouterLink
              v-for="namespace in sortedNamespaces"
              :key="namespace.namespaceID"
              :to="{
                name: 'namespace.view',
                params: { slug: namespace.slug || namespace.namespaceID },
              }"
              class="block relative group cursor-pointer hover:scale-105 hover:text-primary transition-all duration-100"
              v-show="isNamespaceVisible(namespace)"
            >
              <Card
                :pt="{
                  body: {
                    class: 'grow justify-center gap-0 py-1',
                  },
                  title: {
                    class: 'text-center line-clamp-2 group-hover:line-clamp-none',
                  },
                  subtitle: {
                    class: 'text-center line-clamp-2 group-hover:line-clamp-none',
                  },
                }"
                class="group-hover:shadow-lg w-80 min-h-72 group-hover:h-full overflow-hidden"
              >
                <template #header>
                  <div
                    class="relative flex items-center justify-center w-full h-full pt-7 shrink-0"
                  >
                    <Avatar
                      :label="namespace.meta.logoEnabled ? null : namespace.initials"
                      :image="
                        namespace.meta.logoEnabled
                          ? namespace.meta.logo || $Settings.attachment('ui.mainLogo')
                          : null
                      "
                      :pt="{
                        image: { class: 'object-contain' },
                      }"
                      shape="circle"
                      size="xlarge"
                      :class="{ 'text-muted-color bg-emphasis': !namespace.meta.logoEnabled }"
                      class="w-32 h-32 font-bold"
                    />
                    <Button
                      v-if="getActionsMenuItems(namespace).length"
                      icon="pi pi-ellipsis-v"
                      size="small"
                      text
                      severity="secondary"
                      class="!absolute !top-2 !right-2 z-10 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity bg-(--p-card-background)"
                      @click.stop.prevent="showCardMenu($event, namespace)"
                    />
                  </div>
                </template>

                <template #title>
                  {{ namespace.name }}
                </template>

                <template v-if="namespace.meta.description" #subtitle>
                  {{ namespace.meta.description }}
                </template>
              </Card>
            </RouterLink>
          </div>
        </div>
      </div>
    </div>

    <CResourceList
      v-else
      primary-key="namespaceID"
      :fields="namespaceFields"
      :items="namespaceList"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('general.label.namespace.single'),
        resourcePlural: $t('general.label.namespace.plural'),
        searchPlaceholder: $t('namespace.searchPlaceholder'),
      }"
      clickable
      class="h-full"
      @sort="handleSort"
      @search="filterList"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
      @update:filter="onFilterUpdate"
    >
      <template #filter>
        <SelectButton
          v-model="viewMode"
          :options="viewModeOptions"
          option-label="label"
          option-value="value"
          :allow-empty="false"
          size="small"
        >
          <template #option="{ option }">
            <i :class="option.icon" v-tooltip.bottom="option.label" />
          </template>
        </SelectButton>
      </template>

      <template #header>
        <div class="flex items-center gap-2">
          <CRouterLinkButton
            v-if="canCreate"
            :to="{ name: 'namespace.create' }"
            :label="$t('namespace.manage.toolbar.buttons.create')"
            icon="pi pi-plus"
            size="small"
          />
          <NamespaceImporter v-if="canCreate" @imported="onImported" @failed="onFailed" />
          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::compose:namespace/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex items-center gap-2">
          <Avatar
            :label="data.meta?.logoEnabled ? null : getInitials(data)"
            :image="
              data.meta?.logoEnabled ? data.meta?.logo || $Settings.attachment('ui.mainLogo') : null
            "
            :pt="{ image: { class: 'object-contain' } }"
            shape="circle"
            :class="{ 'text-muted-color bg-emphasis': !data.meta?.logoEnabled }"
            class="font-bold shrink-0 mr-1 !w-10 !h-10 !text-sm"
          />
          <div class="flex flex-col min-w-0">
            <span class="truncate text-sm">{{ data.name || '—' }}</span>
            <span v-if="data.meta?.description" class="text-xs text-muted-color truncate">
              {{ data.meta.description }}
            </span>
          </div>
        </div>
      </template>
    </CResourceList>
  </div>

  <TieredMenu ref="cardMenuRef" :model="cardMenuItems" popup>
    <template #item="{ item, props: menuProps }">
      <a v-ripple v-bind="menuProps.action" :class="item.class">
        <span :class="item.icon" />
        <span class="ml-2">{{ item.label }}</span>
      </a>
    </template>
  </TieredMenu>
</template>

<script setup>
import NamespaceImporter from '@/sections/compose/components/Namespaces/NamespaceImporter.vue'
import { useNamespaceStore } from '@planetcrust/human-vue'
import {
  components,
  useConfirmDelete,
  usePermissions,
  useRBACStore,
  useResourceList,
} from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
const { CInputSearch, CPermissionsButton, CResourceList, CRouterLinkButton } = components

const { t } = useI18n()
const router = useRouter()
const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const namespaceStore = useNamespaceStore()
const rbac = useRBACStore()

const canGrant = computed(() => rbac.can('compose/', 'grant'))
const canCreate = computed(() => rbac.can('compose/', 'namespace.create'))

const viewMode = ref('cards')
const viewModeOptions = computed(() => [
  { value: 'cards', icon: 'pi pi-th-large', label: t('namespace.view-mode.cards') },
  { value: 'list', icon: 'pi pi-list', label: t('namespace.view-mode.list') },
])

// Cards mode
const query = ref('')

const normalizedQuery = computed(() => (query.value || '').trim().toUpperCase())

const isNamespaceVisible = namespace => {
  const q = normalizedQuery.value
  if (!q) return true
  return (namespace.slug + namespace.name).toUpperCase().indexOf(q) > -1
}

const sortedNamespaces = computed(() =>
  [...namespaceStore.set].sort((a, b) =>
    (a.name || '').localeCompare(b.name || '', undefined, { sensitivity: 'base' }),
  ),
)

const areNamespacesVisible = computed(() => sortedNamespaces.value.some(isNamespaceVisible))

const cardMenuRef = ref()
const cardMenuItems = ref([])

function showCardMenu(event, namespace) {
  cardMenuItems.value = getActionsMenuItems(namespace)
  cardMenuRef.value?.toggle(event)
}

function getInitials(ns) {
  let base = ns?.name || ns?.slug || ''
  if (base.length <= 3) return base
  const initials = base
    .split(/\s+/)
    .map(w => w[0])
    .filter(c => /[a-zA-Z]/.test(c))
    .join('')
  return initials.slice(0, 3) || base.slice(0, 3)
}

// List mode
const namespaceFields = [
  {
    key: 'name',
    sortable: true,
    header: t('namespace.manage.table.columns.name'),
    style: 'max-width: 32rem;',
  },
  {
    key: 'slug',
    sortable: true,
    header: t('namespace.manage.table.columns.slug'),
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
  router.push({
    name: 'namespace.view',
    params: { slug: data.slug || data.namespaceID },
  })
}

function onFilterUpdate(newFilter) {
  Object.assign(filter, newFilter)
}

function goToEdit(namespace) {
  router.push({
    name: 'namespace.edit',
    params: { slug: namespace.slug || namespace.namespaceID },
  })
}

const getActionsMenuItems = namespace => {
  const items = []

  if (namespace.canUpdateNamespace) {
    items.push({
      label: t('general.label.edit'),
      icon: 'pi pi-pencil',
      command: () => goToEdit(namespace),
    })
  }

  if (namespace.canExportNamespace) {
    items.push({
      label: t('namespace.export'),
      icon: 'pi pi-download',
      command: () => exportNamespace(namespace),
    })
  }

  if (namespace.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        openPermissions({
          resource: `corteza::compose:namespace/${namespace.namespaceID}`,
          title: namespace.name || namespace.slug || namespace.namespaceID,
        })
      },
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

function handleDelete(namespace) {
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
          refreshAll()
        })
        .catch(error => {
          $toast.toastDanger(error.message || t('namespace.manage.delete.error'))
        })
    },
  })
}

function exportNamespace(namespace) {
  const params = {
    namespaceID: namespace.namespaceID,
    filename: encodeURIComponent((namespace.name || 'namespace').replace(/\./g, '-')),
  }

  const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''
  const exportUrl = `${$ComposeAPI.baseURL}${$ComposeAPI.namespaceExportEndpoint(params)}?jwt=${encodeURIComponent(token)}`
  window.open(exportUrl)
}

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

function refreshAll() {
  namespaceStore.load()
  filterList()
}
</script>
