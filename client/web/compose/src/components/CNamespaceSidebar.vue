<template>
  <!-- Teleport sidebar content to CSidebar targets -->
  <Teleport to="#sidebar-header-expanded" :defer="true">
    <!-- Namespace Switcher -->
    <FloatLabel variant="on">
      <Select
        id="namespace-selector"
        :model-value="currentNamespaceObject"
        :options="enabledNamespaces"
        option-label="name"
        data-key="namespaceID"
        :placeholder="$t('sidebar.namespaceSelector.placeholder')"
        class="w-full"
        @update:model-value="handleNamespaceChange"
      >
        <template #option="{ option }">
          <div class="flex items-center gap-2">
            <span>{{ option.name }}</span>
          </div>
        </template>
        <template #footer>
          <RouterLink
            :to="{ name: 'namespace.list' }"
            class="block p-2 text-sm text-muted-color hover:text-primary transition-colors text-center border-t"
          >
            {{ $t('sidebar.namespaceSelector.viewAll') }}
          </RouterLink>
        </template>
      </Select>
      <label for="namespace-selector">{{ $t('sidebar.namespaceSelector.label') }}</label>
    </FloatLabel>
  </Teleport>

  <!-- Page tree navigation -->
  <Teleport to="#sidebar-body-expanded" :defer="true">
    <div v-if="namespace && pageNavItems.length">
      <Divider class="my-0 mb-2" />
      <CSidebarNav
        :items="pageNavItems"
        id-key="pageID"
        parent-key="selfID"
        label-key="title"
        weight-key="weight"
        route-key="_route"
        :filter-fn="p => p.visible"
      />
    </div>
  </Teleport>

  <!-- Navigation Items at the bottom of sidebar -->
  <Teleport to="#sidebar-footer-expanded" :defer="true">
    <div v-if="namespace">
      <CSidebarNav
        :items="adminNavItems"
        id-key="_id"
        parent-key="_parentId"
        label-key="_label"
        icon-key="_icon"
        divider-key="_divider"
        route-key="_route"
      />
    </div>
  </Teleport>
</template>

<script setup>
import { useModuleStore } from '@/stores/module'
import { usePageStore } from '@/stores/page'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { CSidebarNav } = components
const { t } = useI18n()

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
  namespaces: {
    type: Array,
    default: () => [],
  },
})

const route = useRoute()
const router = useRouter()
const moduleStore = useModuleStore()
const pageStore = usePageStore()

// Computed
const enabledNamespaces = computed(() => {
  return props.namespaces.filter(ns => ns.enabled)
})

const currentNamespaceSlug = computed(() => {
  return route.params.slug || null
})

// For Select model-value
const currentNamespaceObject = computed(() => {
  const urlPart = currentNamespaceSlug.value
  return (
    enabledNamespaces.value.find(ns => ns.slug === urlPart || ns.namespaceID === urlPart) || null
  )
})

const isAdminRoute = computed(() => {
  return route.name?.toString().includes('admin.')
})

// Page tree items with routes to public page view
const pageNavItems = computed(() => {
  return pageStore.set.map(p => ({
    ...p,
    _route: { name: 'page', params: { pageID: p.pageID } },
  }))
})

// Admin nav items: Modules (with children), Pages (with children), Charts
const adminNavItems = computed(() => [
  {
    _id: 'modules',
    _parentId: '0',
    _label: t('sidebar.modules'),
    _icon: 'pi pi-database',
    _divider: true,
    _route: { name: 'admin.modules' },
  },
  ...moduleStore.set.map(m => ({
    _id: m.moduleID,
    _parentId: 'modules',
    _label: m.name || m.handle || m.moduleID,
    _route: { name: 'admin.modules.edit', params: { moduleID: m.moduleID } },
  })),
  {
    _id: 'pages',
    _parentId: '0',
    _label: t('sidebar.pages'),
    _icon: 'pi pi-file',
    _route: { name: 'admin.pages' },
  },
  ...pageStore.set.map(p => ({
    _id: `page-${p.pageID}`,
    _parentId: p.selfID && p.selfID !== '0' ? `page-${p.selfID}` : 'pages',
    _label: p.title || p.handle || p.pageID,
    _route: { name: 'admin.pages.edit', params: { pageID: p.pageID } },
    weight: p.weight,
  })),
  {
    _id: 'charts',
    _parentId: '0',
    _label: t('sidebar.charts'),
    _icon: 'pi pi-chart-bar',
    _route: { name: 'admin.charts' },
  },
])

// Methods
function handleNamespaceChange(namespace) {
  const newSlug = namespace?.slug || namespace?.namespaceID

  if (newSlug && newSlug !== currentNamespaceSlug.value) {
    const routeName = isAdminRoute.value ? 'admin.modules' : 'pages'
    router.push({ name: routeName, params: { slug: newSlug } })
  }
}
</script>
