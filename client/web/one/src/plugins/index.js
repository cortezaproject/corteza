import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import Router from 'vue-router'
import Vuex from 'vuex'
import VueTour from 'vue-tour'
import VueNativeSock from 'vue-native-websocket'

import { plugins, websocket } from '@cortezaproject/corteza-vue'

Vue.use(BootstrapVue, {
  BToast: {
    // see https://bootstrap-vue.org/docs/components/toast#comp-ref-b-toast-props
    autoHideDelay: 7000,
    toaster: 'b-toaster-bottom-right',
  },
})

Vue.use(plugins.Auth(), {
  app: 'unify',
  rootApp: true,
})

Vue.use(VueTour)
Vue.use(Router)
Vue.use(Vuex)
Vue.use(BootstrapVue)

Vue.use(plugins.CortezaAPI('system'))
Vue.use(plugins.CortezaAPI('compose'))
Vue.use(plugins.CortezaAPI('automation'))

Vue.use(plugins.Settings, { api: Vue.prototype.$SystemAPI })

Vue.use(VueNativeSock, websocket.endpoint(), websocket.config)
