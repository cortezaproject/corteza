<template>
  <div class="flex flex-col min-h-0 border-t surface-border">
    <!-- Search input -->
    <CInputSearch
      id="sidebar-search"
      v-model="searchQuery"
      :placeholder="searchPlaceholder"
      size="small"
      class="my-2"
    />

    <!-- The nav's one scroller: pages and admin panel scroll together inside
         it, while the search box and the namespace switcher above stay put.
         `min-h-full` on the column makes it fill the scroller when there is
         little to show and grow past it when there is more — so with nothing
         to scroll the admin panel is pushed to the foot, and with something to
         scroll it simply travels with the pages. -->
    <div class="flex-1 overflow-auto">
      <div class="flex flex-col min-h-full">
        <div class="flex-1">
          <CSidebarNav
            v-if="filteredPageNavItems.length"
            :items="filteredPageNavItems"
            id-key="pageID"
            parent-key="selfID"
            label-key="title"
            icon-key="_icon"
            weight-key="weight"
            route-key="_route"
            :filter-fn="p => p.visible"
            :expand-all="hasSearch"
          />

          <!-- No results -->
          <div
            v-if="hasSearch && !filteredPageNavItems.length && !filteredAdminNavItems.length"
            class="flex items-center justify-center py-8 text-muted-color text-sm"
          >
            {{ $t('sidebar.noResults') }}
          </div>
        </div>

        <!-- The namespace's administration, named so it reads as a different
             kind of place from the pages above it. -->
        <div v-if="filteredAdminNavItems.length" class="pt-2" data-testid="sidebar-admin-panel">
          <h2 class="px-2 pb-1 text-xs font-semibold uppercase tracking-wide text-muted-color">
            {{ $t('sidebar.adminPanel') }}
          </h2>
          <div class="border-t surface-border pt-1">
            <CSidebarNav
              :items="filteredAdminNavItems"
              id-key="_id"
              parent-key="_parentId"
              label-key="_label"
              icon-key="_icon"
              route-key="_route"
              :expand-all="hasSearch"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useModuleStore } from '@planetcrust/human-vue'
import { useNamespaceStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CSidebarNav, CInputSearch } = components
const { t } = useI18n()
const route = useRoute()
const $ComposeAPI = inject('$ComposeAPI')

const moduleStore = useModuleStore()
const namespaceStore = useNamespaceStore()
const pageStore = usePageStore()

const searchQuery = ref('')

const hasSearch = computed(() => searchQuery.value.trim().length > 0)

const normalizedQuery = computed(() => searchQuery.value.trim().toLowerCase())

const isAdminRoute = computed(() => {
  return route.name?.toString().includes('admin.')
})

const searchPlaceholder = computed(() => {
  return isAdminRoute.value
    ? t('sidebar.searchPlaceholder.admin')
    : t('sidebar.searchPlaceholder.public')
})

// The nav draws a page icon as an <img>, so only the types that carry a URL can
// resolve: 'link' is one already, 'attachment' is a path on this instance. A
// 'library' ("font-awesome://home") or inline-svg icon has nothing to point at —
// left to fall through it becomes an <img> src that 404s and the page shows a
// broken image next to its name, so it renders as no icon instead.
function resolvePageIcon(p) {
  const icon = p?.config?.navItem?.icon
  if (!icon?.src) return undefined
  if (icon.type === 'library' || icon.type === 'inline-svg' || icon.type === 'svg') return undefined
  if (icon.type === 'link') return icon.src
  return `${$ComposeAPI?.baseURL || ''}${icon.src}`
}

// Page tree items with routes to public page view.
// Record pages route to the new-record creator (recordID '0') so a click
// opens the creation form instead of the page (which has no record context).
const pageNavItems = computed(() => {
  return pageStore.set.map(p => ({
    ...p,
    _icon: resolvePageIcon(p),
    _route: p.isRecordPage
      ? {
          name: 'page.record',
          params: { slug: route.params.slug, pageID: p.pageID, recordID: '0' },
        }
      : { name: 'page', params: { slug: route.params.slug, pageID: p.pageID } },
  }))
})

// Filtered page nav items: keep matching items and their ancestors
const filteredPageNavItems = computed(() => {
  if (!hasSearch.value) return pageNavItems.value

  const query = normalizedQuery.value
  const items = pageNavItems.value

  // Find items whose title matches the query
  const matchingIds = new Set()
  for (const item of items) {
    const label = (item.title || item.handle || item.pageID || '').toLowerCase()
    if (label.includes(query)) {
      matchingIds.add(item.pageID)
    }
  }

  // Collect ancestor IDs so the tree stays intact
  const keepIds = new Set(matchingIds)
  const itemMap = new Map(items.map(i => [i.pageID, i]))
  for (const id of matchingIds) {
    let current = itemMap.get(id)
    while (current && current.selfID && current.selfID !== '0') {
      keepIds.add(current.selfID)
      current = itemMap.get(current.selfID)
    }
  }

  return items.filter(i => keepIds.has(i.pageID))
})

// The admin panel opens on the namespace's `manage` op and nothing else. Without
// it the section is absent rather than disabled: every entry in it leads to a
// screen the server refuses, and an offer that refuses on arrival is worse than
// one that was never made.
const canManageNamespace = computed(
  () => !!namespaceStore.getByUrlPart(route.params.slug)?.canManageNamespace,
)

// Admin nav items: Modules (with children), Pages (with children), Charts
const adminNavItems = computed(() => {
  if (!canManageNamespace.value) return []

  return [
    {
      _id: 'modules',
      _parentId: '0',
      _label: t('sidebar.modules'),
      _icon: 'pi pi-database',
      _route: { name: 'admin.modules', params: { slug: route.params.slug } },
    },
    ...moduleStore.set.map(m => ({
      _id: m.moduleID,
      _parentId: 'modules',
      _label: m.name || m.handle || m.moduleID,
      _route: {
        name: 'admin.modules.edit',
        params: { slug: route.params.slug, moduleID: m.moduleID },
      },
    })),
    {
      _id: 'pages',
      _parentId: '0',
      _label: t('sidebar.pages'),
      _icon: 'pi pi-objects-column',
      _route: { name: 'admin.pages', params: { slug: route.params.slug } },
    },
    ...pageStore.set.map(p => ({
      _id: `page-${p.pageID}`,
      _parentId: p.selfID && p.selfID !== '0' ? `page-${p.selfID}` : 'pages',
      _label: p.title || p.handle || p.pageID,
      _route: { name: 'admin.pages.edit', params: { slug: route.params.slug, pageID: p.pageID } },
      weight: p.weight,
    })),
    {
      _id: 'charts',
      _parentId: '0',
      _label: t('sidebar.charts'),
      _icon: 'pi pi-chart-bar',
      _route: { name: 'admin.charts', params: { slug: route.params.slug } },
    },
  ]
})

// Filtered admin nav items: keep matching children + their parent groups
const filteredAdminNavItems = computed(() => {
  if (!hasSearch.value) return adminNavItems.value

  const query = normalizedQuery.value
  const items = adminNavItems.value
  const rootIds = new Set(['modules', 'pages', 'charts'])

  // Find child items matching the query
  const matchingIds = new Set()
  for (const item of items) {
    if (rootIds.has(item._id)) continue // Skip root group items from matching
    const label = (item._label || '').toLowerCase()
    if (label.includes(query)) {
      matchingIds.add(item._id)
    }
  }

  // Also check if root labels match (e.g., searching "chart")
  for (const item of items) {
    if (rootIds.has(item._id)) {
      const label = (item._label || '').toLowerCase()
      if (label.includes(query)) {
        matchingIds.add(item._id)
      }
    }
  }

  // Collect ancestor IDs so the tree stays intact
  const keepIds = new Set(matchingIds)
  const itemMap = new Map(items.map(i => [i._id, i]))
  for (const id of matchingIds) {
    let current = itemMap.get(id)
    while (current && current._parentId && current._parentId !== '0') {
      keepIds.add(current._parentId)
      current = itemMap.get(current._parentId)
    }
  }

  // If a root group matches, keep all its children too
  for (const id of [...matchingIds]) {
    if (rootIds.has(id)) {
      for (const item of items) {
        if (item._parentId === id) {
          keepIds.add(item._id)
        }
      }
    }
  }

  return items.filter(i => keepIds.has(i._id))
})
</script>
