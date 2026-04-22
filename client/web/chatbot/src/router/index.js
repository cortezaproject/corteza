import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      component: () => import('../views/List.vue'),
    },
    {
      path: '/create',
      name: 'chatbot.create',
      component: () => import('../views/Editor.vue'),
    },
    {
      path: '/:chatbotID/edit',
      name: 'chatbot.edit',
      component: () => import('../views/Editor.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
