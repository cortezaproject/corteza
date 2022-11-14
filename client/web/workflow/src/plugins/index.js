import Vue from 'vue'
import Router from 'vue-router'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { plugins } from '@cortezaproject/corteza-vue'

Vue.use(Router)

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

Vue.use(plugins.Auth(), { app: 'workflow' })

Vue.use(plugins.CortezaAPI('system'))
Vue.use(plugins.CortezaAPI('compose'))
Vue.use(plugins.CortezaAPI('automation'))
Vue.use(plugins.Settings, { api: Vue.prototype.$SystemAPI })
