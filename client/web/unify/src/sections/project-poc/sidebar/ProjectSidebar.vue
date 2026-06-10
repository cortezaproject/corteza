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
import { useProjectsStore } from '@/sections/project-poc/stores/projects'
import { components } from '@planetcrust/human-vue'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'

const { CSidebarNav } = components

const store = useProjectsStore()
const { projects } = storeToRefs(store)

store.load()

// A published project opens its read-only overview; a draft opens the wizard.
// Mirrors ProjectList's row navigation.
const routeFor = p => ({
  name: p.status === 'published' ? 'project.overview' : 'project.wizard',
  params: { projectId: p.id },
})

const navItems = computed(() => [
  {
    _id: 'list',
    _parentId: '0',
    _label: 'All Projects',
    _icon: 'pi pi-folder',
    _route: { name: 'project.list' },
  },
  {
    _id: 'projects',
    _parentId: '0',
    _label: 'Projects',
    _icon: 'pi pi-folder-open',
    _divider: true,
  },
  ...[...projects.value]
    .sort((a, b) => (a.name || '').localeCompare(b.name || ''))
    .map(p => ({
      _id: p.id,
      _parentId: 'projects',
      _label: p.name || 'Untitled project',
      _icon: p.mode === 'gated' ? 'pi pi-shield' : 'pi pi-unlock',
      _route: routeFor(p),
    })),
])
</script>
