import Vue from 'vue'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'

import PortalVue from 'portal-vue'
import './faIcons'
import { components } from '@cortezaproject/corteza-vue'

Vue.use(PortalVue)

Vue.use(BootstrapVue, {
  BToast: {
    // see https://bootstrap-vue.org/docs/components/toast#comp-ref-b-toast-props
    autoHideDelay: 7000,
    toaster: 'b-toaster-bottom-right',
  },
})
Vue.use(BootstrapVueIcons)

Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('c-permissions-button', components.CPermissionsButton)
Vue.component('c-input-confirm', components.CInputConfirm)
Vue.component('c-input-processing', components.CInputProcessing)
Vue.component('c-input-search', components.CInputSearch)
