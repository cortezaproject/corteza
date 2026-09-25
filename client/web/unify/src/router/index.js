import { isAppUrlLocal, resolveAppUrl, useApplicationsStore } from '@planetcrust/human-vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes, sectionById } from '../sections'
import { appReachable } from '../utils/appReachable'
import { makeHomeEntryGuard, resolveHome } from './homeEntry'
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

const openable = app =>
  app.enabled &&
  appReachable(app, {
    sectionFor: path => sectionById(router.resolve(path).meta?.section),
    applications: useApplicationsStore(),
  })

const hrefOf = app => ({
  href: resolveAppUrl(app.unify?.url),
  local: isAppUrlLocal(app.unify?.url),
})

// Where the topbar's home button leads: the home application the user can
// open, or '' for the launcher.
export function homeHref() {
  const applications = useApplicationsStore()
  const { app } = resolveHome({ ownID: applications.ownHomeID, apps: applications.apps, openable })
  return app ? hrefOf(app).href : ''
}

router.beforeEach(
  makeHomeEntryGuard({
    useApplications: useApplicationsStore,
    ownHomeID: () => useApplicationsStore().ownHomeID,
    openable,
    hrefOf,
    leave: href => window.location.replace(href),
  }),
)
router.beforeEach(makeSectionAccessGuard({ useApplications: useApplicationsStore, sectionById }))

export default router
