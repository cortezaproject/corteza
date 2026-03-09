import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      redirect: { name: 'dashboard' },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue'),
    },

    // Connections
    {
      path: '/system/connections',
      name: 'system.connections',
      component: () => import('../views/system/Connection/List.vue'),
    },
    {
      path: '/system/connections/new',
      name: 'system.connections.create',
      component: () => import('../views/system/Connection/Editor.vue'),
    },
    {
      path: '/system/connections/:connectionID',
      name: 'system.connections.edit',
      component: () => import('../views/system/Connection/Editor.vue'),
    },

    // Data Sources (DalConnection)
    {
      path: '/system/data-sources',
      name: 'system.dataSources',
      component: () => import('../views/system/DataSource/List.vue'),
    },
    {
      path: '/system/data-sources/new',
      name: 'system.dataSources.create',
      component: () => import('../views/system/DataSource/Editor.vue'),
    },
    {
      path: '/system/data-sources/:connectionID',
      name: 'system.dataSources.edit',
      component: () => import('../views/system/DataSource/Editor.vue'),
    },

    // Sensitivity Levels
    {
      path: '/system/sensitivity-levels',
      name: 'system.sensitivityLevels',
      component: () => import('../views/system/SensitivityLevel/List.vue'),
    },
    {
      path: '/system/sensitivity-levels/new',
      name: 'system.sensitivityLevels.create',
      component: () => import('../views/system/SensitivityLevel/Editor.vue'),
    },
    {
      path: '/system/sensitivity-levels/:sensitivityLevelID',
      name: 'system.sensitivityLevels.edit',
      component: () => import('../views/system/SensitivityLevel/Editor.vue'),
    },

    // Users
    {
      path: '/system/users',
      name: 'system.users',
      component: () => import('../views/system/User/List.vue'),
    },
    {
      path: '/system/users/new',
      name: 'system.users.create',
      component: () => import('../views/system/User/Editor.vue'),
    },
    {
      path: '/system/users/:userID',
      name: 'system.users.edit',
      component: () => import('../views/system/User/Editor.vue'),
    },

    // User Groups
    {
      path: '/system/user-groups',
      name: 'system.userGroups',
      component: () => import('../views/system/UserGroup/List.vue'),
    },
    {
      path: '/system/user-groups/new',
      name: 'system.userGroups.create',
      component: () => import('../views/system/UserGroup/Editor.vue'),
    },
    {
      path: '/system/user-groups/:userGroupID',
      name: 'system.userGroups.edit',
      component: () => import('../views/system/UserGroup/Editor.vue'),
    },

    // Roles
    {
      path: '/system/roles',
      name: 'system.roles',
      component: () => import('../views/system/Role/List.vue'),
    },
    {
      path: '/system/roles/new',
      name: 'system.roles.create',
      component: () => import('../views/system/Role/Editor.vue'),
    },
    {
      path: '/system/roles/:roleID',
      name: 'system.roles.edit',
      component: () => import('../views/system/Role/Editor.vue'),
    },

    // Applications
    {
      path: '/system/applications',
      name: 'system.applications',
      component: () => import('../views/system/Application/List.vue'),
    },
    {
      path: '/system/applications/new',
      name: 'system.applications.create',
      component: () => import('../views/system/Application/Editor.vue'),
    },
    {
      path: '/system/applications/:applicationID',
      name: 'system.applications.edit',
      component: () => import('../views/system/Application/Editor.vue'),
    },

    // Auth Clients
    {
      path: '/system/auth-clients',
      name: 'system.authClients',
      component: () => import('../views/system/AuthClient/List.vue'),
    },
    {
      path: '/system/auth-clients/new',
      name: 'system.authClients.create',
      component: () => import('../views/system/AuthClient/Editor.vue'),
    },
    {
      path: '/system/auth-clients/:authClientID',
      name: 'system.authClients.edit',
      component: () => import('../views/system/AuthClient/Editor.vue'),
    },

    // Templates
    {
      path: '/system/templates',
      name: 'system.templates',
      component: () => import('../views/system/Template/List.vue'),
    },
    {
      path: '/system/templates/new',
      name: 'system.templates.create',
      component: () => import('../views/system/Template/Editor.vue'),
    },
    {
      path: '/system/templates/:templateID',
      name: 'system.templates.edit',
      component: () => import('../views/system/Template/Editor.vue'),
    },

    // Queues
    {
      path: '/system/queues',
      name: 'system.queues',
      component: () => import('../views/system/Queue/List.vue'),
    },
    {
      path: '/system/queues/new',
      name: 'system.queues.create',
      component: () => import('../views/system/Queue/Editor.vue'),
    },
    {
      path: '/system/queues/:queueID',
      name: 'system.queues.edit',
      component: () => import('../views/system/Queue/Editor.vue'),
    },

    // API Gateway
    {
      path: '/system/api-gateway',
      name: 'system.apiGateway',
      component: () => import('../views/system/ApiGateway/List.vue'),
    },
    {
      path: '/system/api-gateway/new',
      name: 'system.apiGateway.create',
      component: () => import('../views/system/ApiGateway/Editor.vue'),
    },
    {
      path: '/system/api-gateway/:routeID',
      name: 'system.apiGateway.edit',
      component: () => import('../views/system/ApiGateway/Editor.vue'),
    },

    // Action Log
    {
      path: '/system/action-log',
      name: 'system.actionLog',
      component: () => import('../views/system/ActionLog/List.vue'),
    },

    // Settings
    {
      path: '/system/settings',
      name: 'system.settings',
      component: () => import('../views/system/Settings/Index.vue'),
    },

    // Email
    {
      path: '/system/email',
      name: 'system.email',
      component: () => import('../views/system/Email/Index.vue'),
    },

    // Code Snippets
    {
      path: '/system/code-snippets',
      name: 'system.codeSnippets',
      component: () => import('../views/system/CodeSnippets/Index.vue'),
    },

    // ── Compose ──────────────────────────────────────────────
    {
      path: '/compose/settings',
      name: 'compose.settings',
      component: () => import('../views/compose/Settings/Index.vue'),
    },

    // ── Automation ───────────────────────────────────────────
    // Workflows
    {
      path: '/automation/workflows',
      name: 'automation.workflows',
      component: () => import('../views/automation/Workflow/List.vue'),
    },
    {
      path: '/automation/workflows/new',
      name: 'automation.workflows.create',
      component: () => import('../views/automation/Workflow/Editor.vue'),
    },
    {
      path: '/automation/workflows/:workflowID',
      name: 'automation.workflows.edit',
      component: () => import('../views/automation/Workflow/Editor.vue'),
    },

    // Sessions
    {
      path: '/automation/sessions',
      name: 'automation.sessions',
      component: () => import('../views/automation/Session/List.vue'),
    },
    {
      path: '/automation/sessions/:sessionID',
      name: 'automation.sessions.view',
      component: () => import('../views/automation/Session/View.vue'),
    },

    // Scripts
    {
      path: '/automation/scripts',
      name: 'automation.scripts',
      component: () => import('../views/automation/Script/Index.vue'),
    },

    // ── UI ─────────────────────────────────────────────────────
    {
      path: '/ui/settings',
      name: 'ui.settings',
      component: () => import('../views/ui/Settings/Index.vue'),
    },

    // Redirect all other routes to root
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
