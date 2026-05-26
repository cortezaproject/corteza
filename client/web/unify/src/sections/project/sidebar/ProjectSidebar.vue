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
import { components } from '@planetcrust/human-vue'
import { computed } from 'vue'

const { CSidebarNav } = components

// Grows as we add screens.
const navItems = computed(() => [
  {
    _id: 'projects',
    _parentId: '0',
    _label: 'Projects',
    _icon: 'pi pi-folder',
    _route: { name: 'project.list' },
  },
])
</script>
