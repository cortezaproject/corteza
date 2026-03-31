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
import { useAutomationStore } from '@/stores/automation'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const automationStore = useAutomationStore()

const navItems = computed(() => [
  {
    _id: 'automations',
    _parentId: '0',
    _label: t('navigation.automations', 'Automations'),
    _icon: 'pi pi-bolt',
    _route: { name: 'list' },
  },
  ...[...automationStore.list]
    .sort((a, b) => {
      const labelA = a.meta?.short || ''
      const labelB = b.meta?.short || ''
      return labelA.localeCompare(labelB)
    })
    .map(a => ({
      _id: a.automationID,
      _parentId: 'automations',
      _label: a.meta?.short || t('list.untitled', 'Untitled'),
      _route: { name: 'builder-edit', params: { id: a.automationID } },
    })),
])
</script>
