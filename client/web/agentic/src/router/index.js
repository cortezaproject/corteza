import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/create',
      name: 'agent.create',
      component: () => import('../views/Editor.vue'),
    },
    {
      path: '/:agentID/edit',
      name: 'agent.edit',
      component: () => import('../views/Editor.vue'),
    },
    // Redirect all other routes to root
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
