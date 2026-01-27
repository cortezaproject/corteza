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
    return state.apps.filter(app => app.unify && app.unify.listed)
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
        state.apps = set.filter(app => app.enabled)
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

  return {
    // state
    apps: toRef(state, 'apps'),
    loading: toRef(state, 'loading'),
    error: toRef(state, 'error'),

    // getters
    unifyOnly,

    // actions
    fetchApplications,
  }
})
