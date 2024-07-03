import Vue from 'vue'
import { createLocalVue, shallowMount as sm, mount as rm } from '@vue/test-utils'
import sinon from 'sinon'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import CContentHeader from 'corteza-webapp-admin/src/components/CContentHeader'
import CResourceListStatusFilter from 'corteza-webapp-admin/src/components/CResourceListStatusFilter'
import store from 'corteza-webapp-admin/src/store'

import resourceTranslations from 'corteza-webapp-admin/src/mixins/resource-translations'
import toast from 'corteza-webapp-admin/src/mixins/toast'
import { components } from '@cortezaproject/corteza-vue'
const { CCorredorManualButtons, CPermissionsButton } = components

// Mixins
Vue.mixin(resourceTranslations)
Vue.mixin(toast)

// Components
Vue.config.ignoredElements = [
  'font-awesome-icon',
]

Vue.use(BootstrapVue)
Vue.use(PortalVue)

Vue.component('c-corredor-manual-buttons', CCorredorManualButtons)
Vue.component('c-permissions-button', CPermissionsButton)
Vue.component('c-content-header', CContentHeader)
Vue.component('c-resource-list-status-filter', CResourceListStatusFilter)

export const writeableWindowLocation = ({ path: value = '/' } = {}) => Object.defineProperty(window, 'location', { writable: true, value })

export const mount = (component, params = {}) => shallowMount(component, { ...params })

export const stdReject = () => sinon.stub().rejects({ message: 'err' })
export const stdResolve = (rtr) => sinon.stub().resolves(rtr)
export const networkReject = () => sinon.stub().rejects({ message: 'Network Error' })

export const stdAuth = (mocks = {}) => ({
  is: sinon.stub().returns(true),
  check: stdResolve(),
  goto: (u) => u,
  open: (u) => u,
  ...mocks,
})

const mounter = (component, { localVue = createLocalVue(), $auth = {}, mocks = {}, stubs = [], ...options } = {}, mount) => {
  return mount(component, {
    localVue,
    stubs: [
      'router-view',
      'router-link',
      'user-roles',
      'c-permissions-button',
      ...stubs,
    ],
    mocks: {
      $t: (e) => e,
      $i18n: {
        i18next: {
          language: 'en',
        },
      },
      $SystemAPI: {},
      $route: { query: { fullPath: '', token: undefined } },
      $auth,
      $store: store,
      ...mocks,
    },
    ...options,
  })
}

export const shallowMount = (...e) => {
  return mounter(...e, sm)
}

export const fullMount = (...e) => {
  return mounter(...e, rm)
}

export default shallowMount
