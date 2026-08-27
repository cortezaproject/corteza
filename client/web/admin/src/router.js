import Vue from 'vue'
import Router from 'vue-router'
import routes from './views/routes'
import * as CortezaVue from '@cortezaproject/corteza-vue'

const c311Router = CortezaVue.c311Router

const router = new Router({
  mode: 'history',
  routes,
})

if (c311Router?.installC311RouterGuards && Vue.prototype.$C311) {
  c311Router.installC311RouterGuards(router, () => Vue.prototype.$C311)
}

export default router
