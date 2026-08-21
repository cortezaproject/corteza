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

  // The one in-flight or settled fetch. Section access is decided from this
  // list, and the router asks before the shell has mounted, so both have to
  // await the same request rather than race two of them.
  let pending = null

  function fetchApplications() {
    state.loading = true

    pending = $SystemAPI
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

    return pending
  }

  // Resolves once the list has been fetched, fetching it if nobody has yet.
  function ready() {
    return pending || fetchApplications()
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

  async function restoreApplication(applicationID) {
    await $SystemAPI.applicationUndelete({ applicationID })
    const raw = await $SystemAPI.applicationRead({ applicationID })
    const idx = state.apps.findIndex(a => a.applicationID === raw.applicationID)
    if (idx !== -1) state.apps[idx] = raw
    else state.apps.push(raw)
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

  // Trailing and leading slashes are written inconsistently across the
  // registry, so every comparison goes through the bare form.
  const bareUrl = app => (app.unify?.url || '').toLowerCase().replace(/^\/|\/$/g, '')

  /** The registry application serving the given route path, or undefined. */
  function appForPath(path) {
    const p = (path || '').toLowerCase()
    return state.apps.find(app => {
      const url = bareUrl(app)
      if (!url) return false
      return p === '/' + url || p.startsWith('/' + url + '/')
    })
  }

  /** The registry application with the given `unify.url`, or undefined. */
  function appByUrl(url) {
    const wanted = (url || '').toLowerCase().replace(/^\/|\/$/g, '')
    if (!wanted) return undefined
    return state.apps.find(app => bareUrl(app) === wanted)
  }

  /**
   * Whether the application matching the given route path is enabled (usable).
   * Used by the unified shell to show the "disabled" screen per section.
   * Matches the path against each app's `unify.url`; returns true when no app
   * matches (e.g. the home/root section, or unregistered paths).
   */
  function isPathEnabled(path) {
    const match = appForPath(path)
    return match ? match.enabled : true
  }

  /**
   * Whether the current user may open the application with the given
   * `unify.url` — the `access` operation, which is distinct from `read`
   * (`read` only puts the application in the menu).
   *
   * Fails closed, unlike isPathEnabled: an unknown url, an unfetched list and
   * a list that failed to load all answer false. The shell decides section
   * entry on this, and "we don't know yet" must never read as "go ahead".
   */
  function canAccessApp(url) {
    return !!appByUrl(url)?.canAccessApplication
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
    ready,
    findByID,
    reorder,
    create,
    update,
    delete: deleteApplication,
    restore: restoreApplication,
    isCurrentAppEnabled,
    isPathEnabled,
    appForPath,
    appByUrl,
    canAccessApp,
  }
})
