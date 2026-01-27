import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue'),
    },
    {
      path: '/builder',
      name: 'builder',
      component: () => import('../views/Builder.vue'),
    },
    {
      path: '/builder/:id',
      name: 'builder-edit',
      component: () => import('../views/Builder.vue'),
    },
    // Catch-all redirect to dashboard
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
