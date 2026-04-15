import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useApplicationsStore = defineStore('applications', () => {
  const $SystemAPI = inject('$SystemAPI')

  const state = reactive({
    apps: [],
    loading: true,
    error: null,
  })

  const unifyOnly = computed(() => {
    return state.apps.filter(app => app.enabled && app.unify && app.unify.listed)
  })

  function fetchApplications() {
    state.loading = true

    if (!$SystemAPI) {
      state.error = new Error('SystemAPI not available via inject')
      state.loading = false
      return Promise.reject(state.error)
    }

    return $SystemAPI
      .applicationList()
      .then(({ set }) => {
        state.apps = set
        state.error = null
      })
      .catch(error => {
        state.error = error
        console.error('Failed to fetch applications:', error)
      })
      .finally(() => {
        state.loading = false
      })
  }

  async function reorder(reorderedApps) {
    const applicationIDs = reorderedApps.map(a => a.applicationID)
    await $SystemAPI.applicationReorder({ applicationIDs })
    // Reflect new order locally: reordered apps first, then any remaining
    const reorderedSet = new Set(applicationIDs)
    const rest = state.apps.filter(a => !reorderedSet.has(a.applicationID))
    state.apps = [...reorderedApps, ...rest]
  }

  /**
   * Checks if the current app is enabled.
   * Matches using VITE_APP_ID against unify.url, falling back to pathname.
   * Returns true if no matching app is found (app not registered).
   */
  function isCurrentAppEnabled() {
    const appId = (import.meta.env?.VITE_APP_ID || '').toLowerCase()
    const path = window.location.pathname.toLowerCase()

    const match = state.apps.find(app => {
      let url = app.unify?.url
      if (!url) return false
      url = url.toLowerCase().replace(/^\/|\/$/g, '')
      // Match by VITE_APP_ID first, then fallback to pathname
      if (appId && url === appId) return true
      const normalizedUrl = '/' + url + '/'
      return path.startsWith(normalizedUrl) || path === '/' + url
    })
    return match ? match.enabled : true
  }

  return {
    // state
    apps: toRef(state, 'apps'),
    loading: toRef(state, 'loading'),
    error: toRef(state, 'error'),

    // getters
    unifyOnly,

    // actions
    fetchApplications,
    reorder,
    isCurrentAppEnabled,
  }
})
