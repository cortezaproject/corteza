import Vue from 'vue'
import Router from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import c3catalogue from './components/C3'
import { components, i18n } from '@cortezaproject/corteza-vue'
import './components'
import './themes'
import './mixins'

const routes = [
  {
    path: '/c3',
    name: 'c3',
    component: components.C3.View,
    props: { catalogue: c3catalogue },
  },
  { path: '*', redirect: { name: 'c3' } },
]

Vue.use(Router)
Vue.use(BootstrapVue)

export default new Vue({
  el: '#app',
  name: 'DevEnv',
  async created () {
    document.body.setAttribute('dir', this.textDirectionality())
  },
  template: '<router-view/>',
  router: new Router({
    mode: 'history',
    routes,
  }),
  i18n: i18n(Vue,
    { app: 'corteza-webapp-admin' },
    'admin',
    'dashboard',
    'navigation',
    'notification',
    'permissions',
    'system.stats',
    'system.applications',
    'system.users',
    'system.roles',
    'system.templates',
    'system.scripts',
    'system.settings',
    'system.authclients',
    'system.permissions',
    'system.actionlog',
    'system.queues',
    'system.apigw',
    'compose.settings',
    'compose.permissions',
    'compose.automation',
    'federation.nodes',
    'federation.permissions',
    'automation.workflows',
    'automation.scripts',
    'automation.sessions',
    'automation.permissions',
    'ui.settings',
  ),
})
