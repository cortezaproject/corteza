import { useApplicationsStore } from '@planetcrust/human-vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes, sectionById } from '../sections'
import { makeSectionAccessGuard } from './sectionAccess'

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

router.beforeEach(makeSectionAccessGuard({ useApplications: useApplicationsStore, sectionById }))

export default router
