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
import { components } from '@cortezaproject/corteza-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const navItems = computed(() => [
  {
    _id: 'dashboard',
    _parentId: '0',
    _label: t('navigation.dashboard'),
    _icon: 'pi pi-home',
    _route: { name: 'dashboard' },
  },
  {
    _id: 'system',
    _parentId: '0',
    _label: t('navigation.system.group'),
    _icon: 'pi pi-cog',
  },
  {
    _id: 'connections',
    _parentId: 'system',
    _label: t('navigation.system.items.connections'),
    _icon: 'pi pi-link',
    _route: { name: 'system.connections' },
  },
  {
    _id: 'data-sources',
    _parentId: 'system',
    _label: t('navigation.system.items.data-sources'),
    _icon: 'pi pi-database',
    _route: { name: 'system.data-sources' },
  },
])
</script>
