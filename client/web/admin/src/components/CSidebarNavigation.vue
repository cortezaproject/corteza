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

  // ── System ──────────────────────────────────────────────
  {
    _id: 'system',
    _parentId: '0',
    _label: t('navigation.system.group'),
    _icon: 'pi pi-cog',
  },
  {
    _id: 'users',
    _parentId: 'system',
    _label: t('navigation.system.items.users'),
    _icon: 'pi pi-users',
    _route: { name: 'system.users' },
  },
  {
    _id: 'roles',
    _parentId: 'system',
    _label: t('navigation.system.items.roles', 'Roles'),
    _icon: 'pi pi-id-card',
    _route: { name: 'system.roles' },
  },
  {
    _id: 'user-groups',
    _parentId: 'system',
    _label: t('navigation.system.items.usergroups', 'User Groups'),
    _icon: 'pi pi-sitemap',
    _route: { name: 'system.userGroups' },
  },
  {
    _id: 'applications',
    _parentId: 'system',
    _label: t('navigation.system.items.applications', 'Applications'),
    _icon: 'pi pi-th-large',
    _route: { name: 'system.applications' },
  },
  {
    _id: 'templates',
    _parentId: 'system',
    _label: t('navigation.system.items.templates', 'Templates'),
    _icon: 'pi pi-file-edit',
    _route: { name: 'system.templates' },
  },
  {
    _id: 'settings',
    _parentId: 'system',
    _label: t('navigation.system.items.settings', 'Settings'),
    _icon: 'pi pi-wrench',
    _route: { name: 'system.settings' },
  },
  {
    _id: 'auth-clients',
    _parentId: 'system',
    _label: t('navigation.system.items.authclients', 'Auth Clients'),
    _icon: 'pi pi-key',
    _route: { name: 'system.authClients' },
  },
  {
    _id: 'action-log',
    _parentId: 'system',
    _label: t('navigation.system.items.actionlog', 'Action Log'),
    _icon: 'pi pi-list',
    _route: { name: 'system.actionLog' },
  },
  {
    _id: 'queues',
    _parentId: 'system',
    _label: t('navigation.system.items.queues', 'Queues'),
    _icon: 'pi pi-inbox',
    _route: { name: 'system.queues' },
  },
  {
    _id: 'api-gateway',
    _parentId: 'system',
    _label: t('navigation.system.items.apigw', 'API Gateway'),
    _icon: 'pi pi-sliders-h',
    _route: { name: 'system.apiGateway' },
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
    _route: { name: 'system.dataSources' },
  },
  {
    _id: 'sensitivity-levels',
    _parentId: 'system',
    _label: t('navigation.system.items.sensitivityLevel', 'Sensitivity Levels'),
    _icon: 'pi pi-shield',
    _route: { name: 'system.sensitivityLevels' },
  },

  // ── Compose ─────────────────────────────────────────────
  {
    _id: 'compose',
    _parentId: '0',
    _label: t('navigation.compose.group', 'Compose'),
    _icon: 'pi pi-table',
  },
  {
    _id: 'compose-settings',
    _parentId: 'compose',
    _label: t('navigation.compose.items.settings', 'Settings'),
    _icon: 'pi pi-wrench',
    _route: { name: 'compose.settings' },
  },

  // ── Automation ───────────────────────────────────────────
  {
    _id: 'automation',
    _parentId: '0',
    _label: t('navigation.automation.group', 'Automation'),
    _icon: 'pi pi-bolt',
  },
  {
    _id: 'automation-workflows',
    _parentId: 'automation',
    _label: t('navigation.automation.items.workflows', 'Workflows'),
    _icon: 'pi pi-sitemap',
    _route: { name: 'automation.workflows' },
  },
  {
    _id: 'automation-sessions',
    _parentId: 'automation',
    _label: t('navigation.automation.items.sessions', 'Sessions'),
    _icon: 'pi pi-play',
    _route: { name: 'automation.sessions' },
  },
])
</script>
