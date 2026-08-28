import Vue from 'vue'
import Router from 'vue-router'
import * as CortezaVue from '@cortezaproject/corteza-vue'
import routes from './views/routes'

const c311Router = CortezaVue.c311Router

const router = new Router({
  mode: 'history',
  routes,
})

// Add global error handler for navigation errors
router.onError((error) => {
  console.warn('Navigation error occurred:', error)
  // Silently handle the error without crashing
})

if (c311Router?.installC311RouterGuards && Vue.prototype.$C311) {
  c311Router.installC311RouterGuards(router, () => Vue.prototype.$C311)
}

export default router
