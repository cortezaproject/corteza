import Vue from 'vue'
import Vuex from 'vuex'

import { store as cvStore } from '@cortezaproject/corteza-vue'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    rbac: {
      namespaced: true,
      ...cvStore.RBAC(Vue.prototype.$SystemAPI),
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
