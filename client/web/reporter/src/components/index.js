import Vue from 'vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import { components } from '@cortezaproject/corteza-vue'

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
