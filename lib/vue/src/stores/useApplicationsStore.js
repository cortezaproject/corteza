import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'
import { cloneDeep } from 'lodash-es'

export const useApplicationsStore = defineStore('applications', () => {
  const $SystemAPI = inject('$SystemAPI')

  const state = reactive({
    apps: [],
    loading: true,
    error: null,
  })

  // Menu visibility is controlled by `unify.listed`. `enabled` is a separate
  // concern (usability) handled by isPathEnabled — a listed-but-disabled app
  // still appears in the menu but shows the "disabled" screen when opened.
  const unifyOnly = computed(() => {
    return state.apps.filter(app => app.unify && app.unify.listed)
  })

  function fetchApplications() {
    state.loading = true

    return $SystemAPI
      .applicationList({ sort: 'weight ASC' })
      .then(({ set = [] }) => {
        state.apps = set
        state.error = null
      })
      .catch(error => {
        state.error = error
      })
      .finally(() => {
        state.loading = false
      })
  }

  async function findByID(applicationID) {
    const raw = await $SystemAPI.applicationRead({ applicationID })
    const idx = state.apps.findIndex(a => a.applicationID === raw.applicationID)
    if (idx !== -1) state.apps[idx] = raw
    else state.apps.push(raw)
    return cloneDeep(raw)
  }

  async function reorder(reorderedApps) {
    const applicationIDs = reorderedApps.map(a => a.applicationID)
    await $SystemAPI.applicationReorder({ applicationIDs })
    const reorderedSet = new Set(applicationIDs)
    const rest = state.apps.filter(a => !reorderedSet.has(a.applicationID))
    state.apps = [...reorderedApps, ...rest]
  }

  async function create(payload) {
    const created = await $SystemAPI.applicationCreate(payload)
    state.apps.push(created)
    return cloneDeep(created)
  }

  async function update(payload) {
    const updated = await $SystemAPI.applicationUpdate(payload)
    const idx = state.apps.findIndex(a => a.applicationID === updated.applicationID)
    if (idx !== -1) state.apps[idx] = updated
    return cloneDeep(updated)
  }

  async function deleteApplication(applicationID) {
    await $SystemAPI.applicationDelete({ applicationID })
    state.apps = state.apps.filter(a => a.applicationID !== applicationID)
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
      if (appId && url === appId) return true
      const normalizedUrl = '/' + url + '/'
      return path.startsWith(normalizedUrl) || path === '/' + url
    })
    return match ? match.enabled : true
  }

  /**
   * Whether the application matching the given route path is enabled (usable).
   * Used by the unified shell to show the "disabled" screen per section.
   * Matches the path against each app's `unify.url`; returns true when no app
   * matches (e.g. the home/root section, or unregistered paths).
   */
  function isPathEnabled(path) {
    const p = (path || '').toLowerCase()
    const match = state.apps.find(app => {
      let url = app.unify?.url
      if (!url) return false
      url = url.toLowerCase().replace(/^\/|\/$/g, '')
      if (!url) return false
      const base = '/' + url
      return p === base || p.startsWith(base + '/')
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
    findByID,
    reorder,
    create,
    update,
    delete: deleteApplication,
    isCurrentAppEnabled,
    isPathEnabled,
  }
})
