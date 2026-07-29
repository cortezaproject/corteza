<template>
  <div class="flex flex-col h-full">
    <CSidebarNav
      :items="navItems"
      id-key="_id"
      parent-key="_parentId"
      label-key="_label"
      icon-key="_icon"
      route-key="_route"
      badge-key="_badge"
      expand-all
    />
  </div>
</template>

<script setup>
import { chainHasPublished } from '@/sections/project/config/publishState'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { NoID } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components

const { t } = useI18n()
const store = useProjectsStore()
const { projects } = storeToRefs(store)

// The list is loaded by the section's `preload` (run by the shell on section
// entry), so it's ready here even before this lazily-mounted drawer opens.
// Mutations keep the store fresh via `absorb`, so the tree stays current.

// One row per chain, not per revision (see ProjectList.intent.md). The
// `projects` cache isn't guaranteed heads-only like the list's own fetch: a
// chain's non-head revisions get absorbed too once its wizard is visited
// (store.listRevisions loads the whole chain), so filter here instead of
// trusting what's in the cache. A project is a chain head when no other
// cached project points at it via parentRevisionID — mirrors the backend's
// `heads_only` subquery (server/store/adapters/rdbms/filter.go f.Project).
const referencedRevisionIds = computed(() => {
  const ids = new Set()
  for (const p of projects.value) {
    if (p.parentRevisionID !== NoID) ids.add(p.parentRevisionID)
  }
  return ids
})
const isChainHead = p => !referencedRevisionIds.value.has(p.projectID)

// Routing favours the dashboard once a chain has ever published (ruled
// 2026-07-28, mirrors ProjectList.vue's onRowClick): it's the project's home,
// and the wizard is entered deliberately from its revision switcher from here
// on. A chain that has never published has no dashboard at all (see
// config/publishState.js), so it links to the wizard.
const routeFor = p => ({
  name: chainHasPublished(p) ? 'project.overview' : 'project.wizard',
  params: { projectId: p.projectID },
})

// Trailing one-letter status tag; severities match ProjectList's mapping,
// extended with the statuses the list view doesn't color yet.
const STATUS_SEVERITY = {
  active: 'success',
  published: 'success',
  draft: 'info',
  suspended: 'warn',
}

const badgeFor = p =>
  STATUS_SEVERITY[p.status] && {
    value: t(`project.statusShort.${p.status}`),
    severity: STATUS_SEVERITY[p.status],
    title: t(`project.status.${p.status}`),
  }

// One root entry: the label routes to the list (CSidebarNavItem auto-expands
// on navigate, never collapses), the right chevron toggles the children
// independently.
const navItems = computed(() => [
  {
    _id: 'projects',
    _parentId: '0',
    _label: t('project.sidebar.projects'),
    _icon: 'pi pi-folder',
    _route: { name: 'project.list' },
  },
  // Archived projects are hidden here (deleted ones never reach the store);
  // the All Projects list still shows everything.
  ...projects.value
    .filter(p => p.status !== 'archived' && isChainHead(p))
    .sort((a, b) => (a.name || '').localeCompare(b.name || ''))
    .map(p => ({
      _id: p.projectID,
      _parentId: 'projects',
      _label: p.name || t('project.list.untitled'),
      _icon: 'pi pi-sitemap',
      _route: routeFor(p),
      _badge: badgeFor(p),
    })),
])
</script>
