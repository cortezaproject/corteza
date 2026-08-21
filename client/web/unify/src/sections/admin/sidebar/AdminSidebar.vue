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
import { components, useRBACStore } from '@planetcrust/human-vue'
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()
const $Settings = inject('$Settings')
const rbac = useRBACStore()

const federationEnabled = computed(() => $Settings?.get('federation.enabled', false))

const can = (component, operation) => rbac.can(component, operation)

// Labels and corredor scripts are served by endpoints that carry no operation
// of their own, so there is nothing to ask about them. Reading system settings
// is the nearest true statement — it is held by system administrators and
// nobody else — and erring towards them is the right direction for a screen
// that has no gate at all otherwise.
const canAdministerSystem = computed(() => can('system/', 'settings.read'))

// Every entry that leads somewhere states the permission that reveals it in
// `_can`; a group header states none and survives only while it still has a
// child. Entries the user cannot reach are dropped rather than disabled: the
// admin area is a list of destinations, and one that refuses on arrival is
// worse than one that was never offered.
function reachable(items) {
  const kept = items.filter(item => item._can !== false)
  const parents = new Set(kept.map(item => item._parentId))
  return kept.filter(item => item._route || parents.has(item._id))
}

const navItems = computed(() =>
  reachable([
    {
      _id: 'dashboard',
      _parentId: '0',
      _label: t('navigation.dashboard'),
      _icon: 'pi pi-home',
      _route: { name: 'dashboard' },
      // The landing route: reaching the admin area at all is the permission,
      // and it was already checked before this sidebar rendered.
      _can: true,
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
      _can: can('system/', 'users.search'),
    },
    {
      _id: 'roles',
      _parentId: 'system',
      _label: t('navigation.system.items.roles'),
      _icon: 'pi pi-id-card',
      _route: { name: 'system.roles' },
      _can: can('system/', 'roles.search'),
    },
    {
      _id: 'user-groups',
      _parentId: 'system',
      _label: t('navigation.system.items.usergroups'),
      _icon: 'pi pi-sitemap',
      _route: { name: 'system.userGroups' },
      _can: can('system/', 'user-groups.search'),
    },
    {
      _id: 'labels',
      _parentId: 'system',
      _label: 'Projects (Labels)',
      _icon: 'pi pi-tags',
      _route: { name: 'system.labels' },
      _can: canAdministerSystem.value,
    },
    {
      _id: 'llm-providers',
      _parentId: 'system',
      _label: t('navigation.system.items.llm-providers'),
      _icon: 'pi pi-sparkles',
      _route: { name: 'system.llmProviders' },
      _can: can('system/', 'llm-providers.search'),
    },

    // Infrastructure
    {
      _id: 'connections',
      _parentId: 'system',
      _label: t('navigation.system.items.connections'),
      _icon: 'pi pi-link',
      _route: { name: 'system.connections' },
      _can: can('system/', 'connections.search'),
    },
    {
      _id: 'data-sources',
      _parentId: 'system',
      _label: t('navigation.system.items.data-sources'),
      _icon: 'pi pi-database',
      _route: { name: 'system.dataSources' },
      _can: can('system/', 'dal-connections.search'),
    },

    // Resources
    {
      _id: 'applications',
      _parentId: 'system',
      _label: t('navigation.system.items.applications'),
      _icon: 'pi pi-th-large',
      _route: { name: 'system.applications' },
      _can: can('system/', 'applications.search'),
    },
    {
      _id: 'auth-clients',
      _parentId: 'system',
      _label: t('navigation.system.items.authclients'),
      _icon: 'pi pi-key',
      _route: { name: 'system.authClients' },
      _can: can('system/', 'auth-clients.search'),
    },
    {
      _id: 'templates',
      _parentId: 'system',
      _label: t('navigation.system.items.templates'),
      _icon: 'pi pi-file',
      _route: { name: 'system.templates' },
      _can: can('system/', 'templates.search'),
    },
    {
      _id: 'code-snippets',
      _parentId: 'system',
      _label: t('navigation.system.items.code-snippets'),
      _icon: 'pi pi-code',
      _route: { name: 'system.codeSnippets' },
      _can: can('system/', 'settings.read'),
    },
    {
      _id: 'queues',
      _parentId: 'system',
      _label: t('navigation.system.items.queues'),
      _icon: 'pi pi-inbox',
      _route: { name: 'system.queues' },
      _can: can('system/', 'queues.search'),
    },
    {
      _id: 'api-gateway',
      _parentId: 'system',
      _label: t('navigation.system.items.apigw'),
      _icon: 'pi pi-sliders-h',
      _route: { name: 'system.apiGateway' },
      _can: can('system/', 'apigw-routes.search'),
    },
    {
      _id: 'email',
      _parentId: 'system',
      _label: t('navigation.system.items.email'),
      _icon: 'pi pi-envelope',
      _route: { name: 'system.email' },
      _can: can('system/', 'settings.read'),
    },

    // Policy & Monitoring
    {
      _id: 'action-log',
      _parentId: 'system',
      _label: t('navigation.system.items.actionlog'),
      _icon: 'pi pi-list',
      _route: { name: 'system.actionLog' },
      _can: can('system/', 'action-log.read'),
    },

    // Configuration
    {
      _id: 'settings',
      _parentId: 'system',
      _label: t('navigation.system.items.settings'),
      _icon: 'pi pi-wrench',
      _route: { name: 'system.settings' },
      _can: can('system/', 'settings.read'),
    },
    {
      _id: 'system-permissions',
      _parentId: 'system',
      _label: t('navigation.system.items.permissions'),
      _icon: 'pi pi-lock',
      _route: { name: 'system.permissions' },
      _can: can('system/', 'grant'),
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
      _can: can('compose/', 'settings.read'),
    },
    {
      _id: 'compose-permissions',
      _parentId: 'compose',
      _label: t('navigation.compose.items.permissions'),
      _icon: 'pi pi-lock',
      _route: { name: 'compose.permissions' },
      _can: can('compose/', 'grant'),
    },

    // ── Automation ───────────────────────────────────────────
    {
      _id: 'automation',
      _parentId: '0',
      _label: t('navigation.automation.group'),
      _icon: 'pi pi-bolt',
    },
    {
      _id: 'automation-taq',
      _parentId: 'automation',
      _label: t('navigation.automation.items.taq'),
      _icon: 'pi pi-microchip-ai',
      _route: { name: 'automation.taq' },
      _can: can('automation/', 'ng-automations.search'),
    },
    {
      _id: 'automation-workflows',
      _parentId: 'automation',
      _label: t('navigation.automation.items.workflows'),
      _icon: 'pi pi-share-alt',
      _route: { name: 'automation.workflows' },
      _can: can('automation/', 'workflows.search'),
    },
    {
      _id: 'automation-sessions',
      _parentId: 'automation',
      _label: t('navigation.automation.items.sessions'),
      _icon: 'pi pi-play',
      _route: { name: 'automation.sessions' },
      _can: can('automation/', 'sessions.search'),
    },
    {
      _id: 'automation-scripts',
      _parentId: 'automation',
      _label: t('navigation.automation.items.scripts'),
      _icon: 'pi pi-file-edit',
      _route: { name: 'automation.scripts' },
      _can: canAdministerSystem.value,
    },
    {
      _id: 'automation-permissions',
      _parentId: 'automation',
      _label: t('navigation.automation.items.permissions'),
      _icon: 'pi pi-lock',
      _route: { name: 'automation.permissions' },
      _can: can('automation/', 'grant'),
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
            _can: can('federation/', 'nodes.search'),
          },
          {
            _id: 'federation-permissions',
            _parentId: 'federation',
            _label: t('navigation.federation.items.permissions'),
            _icon: 'pi pi-lock',
            _route: { name: 'federation.permissions' },
            _can: can('federation/', 'grant'),
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
      _can: can('system/', 'settings.read'),
    },
    {
      _id: 'ui-navigation',
      _parentId: 'ui',
      _label: t('navigation.ui.items.navigation'),
      _icon: 'pi pi-bars',
      _route: { name: 'ui.navigation' },
      _can: can('system/', 'settings.read'),
    },
    {
      _id: 'ui-location',
      _parentId: 'ui',
      _label: t('navigation.ui.items.location'),
      _icon: 'pi pi-map-marker',
      _route: { name: 'ui.location' },
      _can: can('system/', 'settings.read'),
    },
  ]),
)
</script>
