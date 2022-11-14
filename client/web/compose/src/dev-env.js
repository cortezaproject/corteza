import Vue from 'vue'
import Router from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import c3catalogue from './components/C3'
import { components, i18n } from '@cortezaproject/corteza-vue'
import './mixins'

import './components'
import './themes'

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
    { app: 'corteza-webapp-compose' },
    'block',
    'chart',
    'field',
    'general',
    'module',
    'namespace',
    'navigation',
    'notification',
    'onboarding',
    'page',
    'permissions',
    'preview',
    'sidebar',
    'resource-translator',
  ),
})
