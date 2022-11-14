import Vue from 'vue'
import Vuex from 'vuex'

import ui from './ui'
import { store as cvStore } from '@cortezaproject/corteza-vue'

Vue.use(Vuex)

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',

  modules: {
    ui,
    rbac: {
      namespaced: true,
      ...cvStore.RBAC(
        Vue.prototype.$SystemAPI,
        Vue.prototype.$ComposeAPI,
        Vue.prototype.$AutomationAPI,
      ),
    },
    wfPrompts: {
      namespaced: true,
      ...cvStore.wfPrompts({
        api: Vue.prototype.$AutomationAPI,
        ws: Vue.prototype.$socket,
        webapp: 'admin',
      }),
    },
  },
})
