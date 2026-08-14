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
import { useAgentStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const $SystemAPI = inject('$SystemAPI')
const agentStore = useAgentStore()

// The shell does global setup only; each section fetches what its own nav needs.
onMounted(() => {
  if (!agentStore.list.length) {
    agentStore.fetchList()
  }
})

const navItems = computed(() => [
  {
    _id: 'agents',
    _parentId: '0',
    _label: t('navigation.agents'),
    _icon: 'pi pi-android',
    _route: { name: 'agentic' },
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
      _label: a.meta?.short || a.handle || t('agent.list.untitled'),
      _route: { name: 'agentic.edit', params: { agentID: a.agentID } },
    })),
])
</script>
