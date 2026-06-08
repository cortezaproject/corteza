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
import { useAutomationStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const $AutomationAPI = inject('$AutomationAPI')
const automationStore = useAutomationStore()

// The shell does global setup only; the sidebar loads the automation list it
// renders. (The builder's functions/triggers catalog is loaded by Builder.vue,
// not here — this sidebar is a lazy PrimeVue Drawer and may not be mounted.)
onMounted(() => {
  if (!automationStore.list.length) {
    automationStore.fetchList()
  }
})

const navItems = computed(() => [
  {
    _id: 'automations',
    _parentId: '0',
    _label: t('navigation.automations'),
    _icon: 'pi pi-bolt',
    _route: { name: 'taq' },
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
      _label: a.meta?.short || t('list.untitled'),
      _route: { name: 'taq.builder-edit', params: { id: a.automationID } },
    })),
])
</script>
