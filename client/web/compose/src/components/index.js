import Vue from 'vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
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
  BarChart,
  PieChart,
  GaugeChart,
  HeatmapChart,
  FunnelChart,
} from 'echarts/charts'
import {
  TitleComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  VisualMapComponent,
  ToolboxComponent,
} from 'echarts/components'

use([
  BarChart,
  LineChart,
  PieChart,
  GaugeChart,
  HeatmapChart,
  FunnelChart,
  SVGRenderer,
  TitleComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  VisualMapComponent,
  ToolboxComponent,
])

Vue.component('e-charts', ECharts)

Vue.use(PortalVue)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('c-permissions-button', components.CPermissionsButton)
Vue.component('c-input-confirm', components.CInputConfirm)
Vue.component('c-input-processing', components.CInputProcessing)
Vue.component('c-resource-list', components.CResourceList)

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
