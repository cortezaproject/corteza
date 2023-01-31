import Vue from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import CContentHeader from './CContentHeader'
import CResourceListStatusFilter from './CResourceListStatusFilter'
import { components } from '@cortezaproject/corteza-vue'
import { LMap, LTileLayer, LMarker } from 'vue2-leaflet'
import 'leaflet/dist/leaflet.css'
import { Icon } from 'leaflet'

// import ECharts modules manually to reduce bundle size
import ECharts from 'vue-echarts'
import { use } from 'echarts/core'
import {
  SVGRenderer,
} from 'echarts/renderers'
import {
  LineChart,
} from 'echarts/charts'
import {
  TitleComponent,
  GridComponent,
  TooltipComponent,
} from 'echarts/components'

const { CCorredorManualButtons, CPermissionsButton, CInputConfirm } = components

Vue.use(PortalVue)
Vue.component('c-corredor-manual-buttons', CCorredorManualButtons)
Vue.component('c-permissions-button', CPermissionsButton)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('c-content-header', CContentHeader)
Vue.component('c-resource-list-status-filter', CResourceListStatusFilter)
Vue.component('c-input-confirm', CInputConfirm)

// Map things
Vue.component('l-map', LMap)
Vue.component('l-tile-layer', LTileLayer)
Vue.component('l-marker', LMarker)

use([
  LineChart,
  SVGRenderer,
  TitleComponent,
  GridComponent,
  TooltipComponent,
])

Vue.component('e-charts', ECharts)

delete Icon.Default.prototype._getIconUrl
Icon.Default.mergeOptions({
  iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: require('leaflet/dist/images/marker-icon.png'),
  shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
})
