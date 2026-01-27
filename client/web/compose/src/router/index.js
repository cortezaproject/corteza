import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      redirect: '/namespaces',
      children: [
        {
          path: '/namespaces',
          name: 'namespace.list',
          component: () => import('../views/Namespace/List.vue'),
        },
        {
          path: '/namespaces/manage',
          name: 'namespace.manage',
          component: () => import('../views/Namespace/Manage.vue'),
        },
        {
          path: '/namespaces/create',
          name: 'namespace.create',
          component: () => import('../views/Namespace/Edit.vue'),
        },
        {
          path: '/namespaces/edit/:slug',
          name: 'namespace.edit',
          component: () => import('../views/Namespace/Edit.vue'),
        },
        {
          path: '/namespace/:slug',
          name: 'namespace.view',
          component: () => import('../views/Namespace/View.vue'),
        },
      ],
    },

    // Redirect all other routes to root
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
