import { useApplicationsStore } from '@planetcrust/human-vue'
import { useRouter } from 'vue-router'
import { sectionById } from '@/sections'
import { sectionAllows } from '@/router/sectionAccess'

// Whether the app menu should offer an application — the same question the
// router's section gate answers on arrival.
//
// It resolves the application's url to a section rather than reading the
// application's own `access`, because the registry accepts a url pointing
// inside another section. Such an entry is a deep link, not a webapp of its
// own, and it is governed by the section it lands in; asking its own `access`
// would offer a tile that bounces the user straight back home.
export function useAppReachable() {
  const router = useRouter()
  const applications = useApplicationsStore()

  return app => {
    const url = String(app?.unify?.url || '').replace(/^\/|\/$/g, '')
    if (!url) return false
    return sectionAllows(sectionById(router.resolve('/' + url).meta?.section), applications)
  }
}
