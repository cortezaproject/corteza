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

store.load()

// A live (published) project opens its dashboard; a draft opens the wizard.
// Live status is `active` (BE never sets `published`), matching ProjectList.
const routeFor = p => ({
  name: ['active', 'published'].includes(p.status) ? 'project.overview' : 'project.wizard',
  params: { projectId: p.id },
})

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
  ...[...projects.value]
    .sort((a, b) => (a.name || '').localeCompare(b.name || ''))
    .map(p => ({
      _id: p.id,
      _parentId: 'projects',
      _label: p.name || t('project.list.untitled'),
      _icon: p.mode === 'gated' ? 'pi pi-shield' : 'pi pi-unlock',
      _route: routeFor(p),
    })),
])
</script>
