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
import { useWorkflowStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const workflowStore = useWorkflowStore()

// The shell does global setup only; each section fetches what its own nav needs.
onMounted(() => {
  if (!workflowStore.list.length) {
    workflowStore.fetchList()
  }
})

const navItems = computed(() => [
  {
    _id: 'workflows',
    _parentId: '0',
    _label: t('navigation.workflows'),
    _icon: 'pi pi-sitemap',
    _route: { name: 'workflow.list' },
  },
  ...[...workflowStore.list]
    .sort((a, b) => {
      const labelA = a.meta?.name || a.handle || ''
      const labelB = b.meta?.name || b.handle || ''
      return labelA.localeCompare(labelB)
    })
    .map(w => ({
      _id: w.workflowID,
      _parentId: 'workflows',
      _label: w.meta?.name || w.handle || t('workflow.untitled'),
      _route: { name: 'workflow.edit', params: { workflowID: w.workflowID } },
    })),
])
</script>
