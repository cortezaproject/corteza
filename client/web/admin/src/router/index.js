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
    {
      path: '/system/data-sources',
      name: 'system.data-sources',
      component: () => import('../views/system/DataSource/List.vue'),
    },

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

    // Redirect all other routes to root
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
