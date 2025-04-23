import Vue from 'vue'
import Vuex from 'vuex'

import applications from './applications'
import { store as cvStore } from '@cortezaproject/corteza-vue'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    applications: applications({ api: Vue.prototype.$SystemAPI }),
    wfPrompts: {
      namespaced: true,
      ...cvStore.wfPrompts({
        api: Vue.prototype.$AutomationAPI,
        ws: Vue.prototype.$socket,
        webapp: 'one',
      }),
    },
    notifications: {
      namespaced: true,
      ...cvStore.notifications({
        api: Vue.prototype.$SystemAPI,
      }),
    },
  },
})

export default store
