import Vue from 'vue'
import Vuex from 'vuex'

import discovery from './discovery'
import { store as cvStore } from '@cortezaproject/corteza-vue'

Vue.use(Vuex)

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',

  modules: {
    discovery,
    notifications: {
      namespaced: true,
      ...cvStore.notifications({
        api: Vue.prototype.$SystemAPI,
      }),
    },
  },
})
