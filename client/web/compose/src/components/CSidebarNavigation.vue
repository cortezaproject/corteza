<template>
  <div class="flex flex-col h-full">
    <!-- Page tree navigation (grows to fill) -->
    <div v-if="pageNavItems.length" class="flex-1 overflow-auto">
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

    <!-- Admin navigation (pinned at bottom) -->
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
</template>

<script setup>
import { useModuleStore } from '@/stores/module'
import { usePageStore } from '@/stores/page'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const moduleStore = useModuleStore()
const pageStore = usePageStore()

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
</script>
