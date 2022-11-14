import Vue from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import CContentHeader from './CContentHeader'
import CResourceList from './CResourceList'
import CResourceListStatusFilter from './CResourceListStatusFilter'
import { components } from '@cortezaproject/corteza-vue'

const { CCorredorManualButtons, CPermissionsButton } = components

Vue.use(PortalVue)
Vue.component('c-corredor-manual-buttons', CCorredorManualButtons)
Vue.component('c-permissions-button', CPermissionsButton)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('c-content-header', CContentHeader)
Vue.component('c-resource-list', CResourceList)
Vue.component('c-resource-list-status-filter', CResourceListStatusFilter)
