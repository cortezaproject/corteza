<template>
  <div class="flex flex-col h-full">
    <CSidebarNav
      :items="navItems"
      id-key="_id"
      parent-key="_parentId"
      label-key="_label"
      icon-key="_icon"
      divider-key="_divider"
      route-key="_route"
      badge-key="_badge"
      expand-all
    />
  </div>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
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

// A live (published) project opens its dashboard; a draft opens the wizard.
// Live status is `active` (BE never sets `published`), matching ProjectList.
const routeFor = p => ({
  name: ['active', 'published'].includes(p.status) ? 'project.overview' : 'project.wizard',
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

const navItems = computed(() => [
  {
    _id: 'list',
    _parentId: '0',
    _label: t('project.sidebar.allProjects'),
    _icon: 'pi pi-folder',
    _route: { name: 'project.list' },
  },
  {
    _id: 'projects',
    _parentId: '0',
    _label: t('project.sidebar.projects'),
    _icon: 'pi pi-folder-open',
    _divider: true,
  },
  // Archived projects are hidden here (deleted ones never reach the store);
  // the All Projects list still shows everything.
  ...projects.value
    .filter(p => p.status !== 'archived')
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
