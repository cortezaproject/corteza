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
      match-type="exact"
      expand-all
    />
  </div>
</template>

<script setup>
import { useAgentStore } from '@/stores/agent'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const agentStore = useAgentStore()

const navItems = computed(() => [
  {
    _id: 'home',
    _parentId: '0',
    _label: t('navigation.home', 'Home'),
    _icon: 'pi pi-home',
    _route: { name: 'root' },
  },
  {
    _id: 'agents',
    _parentId: '0',
    _label: t('navigation.agents', 'Agents'),
    _icon: 'pi pi-android',
    _divider: true,
  },
  ...[...agentStore.list]
    .sort((a, b) => {
      const labelA = a.meta?.short || a.handle || ''
      const labelB = b.meta?.short || b.handle || ''
      return labelA.localeCompare(labelB)
    })
    .map(a => ({
      _id: a.agentID,
      _parentId: 'agents',
      _label: a.meta?.short || a.handle || t('agent.list.untitled', 'Untitled'),
      _route: { name: 'agent.edit', params: { agentID: a.agentID } },
    })),
])
</script>
