import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      redirect: '/list',
    },
    {
      path: '/list',
      name: 'workflow.list',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/new',
      name: 'workflow.create',
      component: () => import('../views/Editor.vue'),
    },
    {
      path: '/:workflowID/edit',
      name: 'workflow.edit',
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
