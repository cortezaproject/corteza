import Vue from 'vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import { components } from '@cortezaproject/corteza-vue'
import { LMap, LTileLayer, LMarker } from 'vue2-leaflet'
import 'leaflet/dist/leaflet.css'
import { Icon } from 'leaflet'

Vue.use(PortalVue)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('c-permissions-button', components.CPermissionsButton)
Vue.component('c-input-confirm', components.CInputConfirm)
Vue.component('c-input-processing', components.CInputProcessing)

// Map things
Vue.component('l-map', LMap)
Vue.component('l-tile-layer', LTileLayer)
Vue.component('l-marker', LMarker)

delete Icon.Default.prototype._getIconUrl
Icon.Default.mergeOptions({
  iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: require('leaflet/dist/images/marker-icon.png'),
  shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
})
