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
      component: () => import('../views/system/Connections.vue'),
    },
    {
      path: '/system/data-sources',
      name: 'system.data-sources',
      component: () => import('../views/system/DataSources.vue'),
    },

    // Redirect all other routes to root
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
