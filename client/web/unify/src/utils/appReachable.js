import { isAppUrlLocal, useApplicationsStore } from '@planetcrust/human-vue'
import { useRouter } from 'vue-router'
import { sectionById } from '@/sections'
import { sectionAllows } from '@/router/sectionAccess'

// Whether the app menu should offer an application.
//
// A url this shell serves is governed by the section it lands in rather than by
// the application's own `access`, because the registry accepts a url pointing
// inside another section. Such an entry is a deep link, not a webapp of its
// own, and asking its own access would offer a tile that bounces the user
// straight back home.
//
// A url no section serves — an address elsewhere, or an app this shell does not
// host — has no section to speak for it, so it answers for itself.
//
// `sectionFor` is passed in rather than imported so the rule can be exercised
// without a router around it.
export function appReachable(app, { sectionFor, applications }) {
  const raw = String(app?.unify?.url || '')
  const path = raw.replace(/^\/|\/$/g, '')
  if (!path) return false

  const section = isAppUrlLocal(raw) ? sectionFor('/' + path) : null
  if (section) return sectionAllows(section, applications)

  return !!app.canAccessApplication
}

export function useAppReachable() {
  const router = useRouter()
  const applications = useApplicationsStore()

  return app =>
    appReachable(app, {
      sectionFor: path => sectionById(router.resolve(path).meta?.section),
      applications,
    })
}
