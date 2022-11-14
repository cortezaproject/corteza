import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import { i18n } from '@cortezaproject/corteza-vue'

import './themes'
import './config-check'
import './console-splash'
import './plugins'
import './mixins'
import './components'

new Vue({
  router,
  store,
  i18n: i18n(Vue,
    { app: 'corteza-webapp-discovery' },
    'general',
    'notification',
    'navigation',
    'filters',
    'search',
  ),
  render: h => h(App),
}).$mount('#app')
