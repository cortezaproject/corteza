import Vue from 'vue'
import Router from 'vue-router'
import { plugins } from '@cortezaproject/corteza-vue'
import DiscoveryAPI from './searcher.js'

Vue.use(Router)
Vue.use(plugins.Auth(), { app: 'discovery' })
Vue.use(plugins.CortezaAPI('system'))
Vue.use(plugins.Settings, { api: Vue.prototype.$SystemAPI })
Vue.use(DiscoveryAPI())
