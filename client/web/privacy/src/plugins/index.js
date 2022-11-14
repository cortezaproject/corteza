import Vue from 'vue'
import Router from 'vue-router'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { plugins } from '@cortezaproject/corteza-vue'

Vue.use(Router)

Vue.use(BootstrapVue, {
  BToast: {
    // see https://bootstrap-vue.org/docs/components/toast#comp-ref-b-toast-props
    autoHideDelay: 7000,
    toaster: 'b-toaster-bottom-right',
  },
})
Vue.use(BootstrapVueIcons)

Vue.use(plugins.Auth(), { app: 'privacy' })

Vue.use(plugins.CortezaAPI('system'))
Vue.use(plugins.CortezaAPI('compose'))

Vue.use(plugins.Settings, { api: Vue.prototype.$SystemAPI })
