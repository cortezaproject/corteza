import { createRouter, createWebHistory } from 'vue-router'
import { routes } from '../sections'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...routes,

    // Anything unknown falls back to the home section.
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
