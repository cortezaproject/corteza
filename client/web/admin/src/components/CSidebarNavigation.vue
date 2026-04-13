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
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()
const $Settings = inject('$Settings')

const federationEnabled = computed(() => $Settings?.get('federation.enabled', false))

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
  // Identity & Access
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
    _label: t('navigation.system.items.roles'),
    _icon: 'pi pi-id-card',
    _route: { name: 'system.roles' },
  },
  {
    _id: 'user-groups',
    _parentId: 'system',
    _label: t('navigation.system.items.usergroups'),
    _icon: 'pi pi-sitemap',
    _route: { name: 'system.userGroups' },
  },
  {
    _id: 'auth-clients',
    _parentId: 'system',
    _label: t('navigation.system.items.authclients'),
    _icon: 'pi pi-key',
    _route: { name: 'system.authClients' },
  },

  // Resources
  {
    _id: 'applications',
    _parentId: 'system',
    _label: t('navigation.system.items.applications'),
    _icon: 'pi pi-th-large',
    _route: { name: 'system.applications' },
  },
  {
    _id: 'templates',
    _parentId: 'system',
    _label: t('navigation.system.items.templates'),
    _icon: 'pi pi-file',
    _route: { name: 'system.templates' },
  },
  {
    _id: 'code-snippets',
    _parentId: 'system',
    _label: t('navigation.system.items.code-snippets'),
    _icon: 'pi pi-code',
    _route: { name: 'system.codeSnippets' },
  },
  {
    _id: 'llm-providers',
    _parentId: 'system',
    _label: t('navigation.system.items.llm-providers'),
    _icon: 'pi pi-sparkles',
    _route: { name: 'system.llmProviders' },
  },

  // Infrastructure
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
    _id: 'queues',
    _parentId: 'system',
    _label: t('navigation.system.items.queues'),
    _icon: 'pi pi-inbox',
    _route: { name: 'system.queues' },
  },
  {
    _id: 'api-gateway',
    _parentId: 'system',
    _label: t('navigation.system.items.apigw'),
    _icon: 'pi pi-sliders-h',
    _route: { name: 'system.apiGateway' },
  },
  {
    _id: 'email',
    _parentId: 'system',
    _label: t('navigation.system.items.email'),
    _icon: 'pi pi-envelope',
    _route: { name: 'system.email' },
  },

  // Policy & Monitoring
  {
    _id: 'labels',
    _parentId: 'system',
    _label: t('navigation.system.items.labels'),
    _icon: 'pi pi-tags',
    _route: { name: 'system.labels' },
  },
  {
    _id: 'action-log',
    _parentId: 'system',
    _label: t('navigation.system.items.actionlog'),
    _icon: 'pi pi-list',
    _route: { name: 'system.actionLog' },
  },

  // Configuration
  {
    _id: 'settings',
    _parentId: 'system',
    _label: t('navigation.system.items.settings'),
    _icon: 'pi pi-wrench',
    _route: { name: 'system.settings' },
  },
  {
    _id: 'system-permissions',
    _parentId: 'system',
    _label: t('navigation.system.items.permissions'),
    _icon: 'pi pi-lock',
    _route: { name: 'system.permissions' },
  },

  // ── Compose ─────────────────────────────────────────────
  {
    _id: 'compose',
    _parentId: '0',
    _label: t('navigation.compose.group'),
    _icon: 'pi pi-table',
  },
  {
    _id: 'compose-settings',
    _parentId: 'compose',
    _label: t('navigation.compose.items.settings'),
    _icon: 'pi pi-wrench',
    _route: { name: 'compose.settings' },
  },
  {
    _id: 'compose-permissions',
    _parentId: 'compose',
    _label: t('navigation.compose.items.permissions'),
    _icon: 'pi pi-lock',
    _route: { name: 'compose.permissions' },
  },

  // ── Automation ───────────────────────────────────────────
  {
    _id: 'automation',
    _parentId: '0',
    _label: t('navigation.automation.group'),
    _icon: 'pi pi-bolt',
  },
  {
    _id: 'automation-workflows',
    _parentId: 'automation',
    _label: t('navigation.automation.items.workflows'),
    _icon: 'pi pi-share-alt',
    _route: { name: 'automation.workflows' },
  },
  {
    _id: 'automation-sessions',
    _parentId: 'automation',
    _label: t('navigation.automation.items.sessions'),
    _icon: 'pi pi-play',
    _route: { name: 'automation.sessions' },
  },
  {
    _id: 'automation-scripts',
    _parentId: 'automation',
    _label: t('navigation.automation.items.scripts'),
    _icon: 'pi pi-file-edit',
    _route: { name: 'automation.scripts' },
  },
  {
    _id: 'automation-permissions',
    _parentId: 'automation',
    _label: t('navigation.automation.items.permissions'),
    _icon: 'pi pi-lock',
    _route: { name: 'automation.permissions' },
  },

  // ── Federation ───────────────────────────────────────────
  ...(federationEnabled.value
    ? [
        {
          _id: 'federation',
          _parentId: '0',
          _label: t('navigation.federation.group'),
          _icon: 'pi pi-share-alt',
        },
        {
          _id: 'federation-nodes',
          _parentId: 'federation',
          _label: t('navigation.federation.items.nodes'),
          _icon: 'pi pi-sitemap',
          _route: { name: 'federation.nodes' },
        },
        {
          _id: 'federation-permissions',
          _parentId: 'federation',
          _label: t('navigation.federation.items.permissions'),
          _icon: 'pi pi-lock',
          _route: { name: 'federation.permissions' },
        },
      ]
    : []),

  // ── UI ───────────────────────────────────────────────────
  {
    _id: 'ui',
    _parentId: '0',
    _label: t('navigation.ui.group'),
    _icon: 'pi pi-desktop',
  },
  {
    _id: 'ui-theming',
    _parentId: 'ui',
    _label: t('navigation.ui.items.theming'),
    _icon: 'pi pi-palette',
    _route: { name: 'ui.theming' },
  },
  {
    _id: 'ui-navigation',
    _parentId: 'ui',
    _label: t('navigation.ui.items.navigation'),
    _icon: 'pi pi-bars',
    _route: { name: 'ui.navigation' },
  },
  {
    _id: 'ui-location',
    _parentId: 'ui',
    _label: t('navigation.ui.items.location'),
    _icon: 'pi pi-map-marker',
    _route: { name: 'ui.location' },
  },
])
</script>
