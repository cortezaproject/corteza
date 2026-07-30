<template>
  <div class="flex flex-col h-full">
    <CSidebarNav
      :items="navItems"
      id-key="_id"
      parent-key="_parentId"
      label-key="_label"
      icon-key="_icon"
      route-key="_route"
      expand-all
    >
      <!-- Where each project stands, in the section's own status indicator
           rather than the nav's one-letter Tag: same colour and icon per status
           as the wizard and the revision switcher state, so a project reads the
           same in the sidebar as it does once opened. Icon-only, with the
           wording in the tooltip — the row is already narrow and already named.
           The root "Projects" entry has no status, hence the guard. -->
      <template #badge="{ node }">
        <StatusChip v-if="node._status" :status="node._status" icon-only />
      </template>
    </CSidebarNav>
  </div>
</template>

<script setup>
import StatusChip from '@/sections/project/components/project/StatusChip.vue'
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
    // No per-project icon: every row carried the same one, so it said nothing
    // the indentation doesn't already say, and the status indicator on the
    // right is the thing worth the room (see the #badge slot above).
    .map(p => ({
      _id: p.projectID,
      _parentId: 'projects',
      _label: p.name || t('project.list.untitled'),
      _route: routeFor(p),
      _status: p.status,
    })),
])
</script>
