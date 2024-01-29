import Vue from 'vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import { components } from '@cortezaproject/corteza-vue'

// import ECharts modules manually to reduce bundle size
import ECharts from 'vue-echarts'
import { use } from 'echarts/core'
import {
  CanvasRenderer,
} from 'echarts/renderers'
import {
  LineChart,
  BarChart,
  PieChart,
  GaugeChart,
  RadarChart,
  FunnelChart,
  ScatterChart,
} from 'echarts/charts'
import {
  TitleComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  VisualMapComponent,
  ToolboxComponent,
  DataZoomComponent,
} from 'echarts/components'

use([
  BarChart,
  LineChart,
  PieChart,
  GaugeChart,
  RadarChart,
  FunnelChart,
  ScatterChart,
  CanvasRenderer,
  TitleComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  VisualMapComponent,
  ToolboxComponent,
  DataZoomComponent,
])

Vue.component('e-charts', ECharts)

Vue.use(PortalVue)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('c-permissions-button', components.CPermissionsButton)
Vue.component('c-input-confirm', components.CInputConfirm)
Vue.component('c-input-processing', components.CInputProcessing)
Vue.component('c-resource-list', components.CResourceList)
Vue.component('c-input-checkbox', components.CInputCheckbox)
Vue.component('c-button-submit', components.CButtonSubmit)
Vue.component('c-hint', components.CHint)
Vue.component('c-input-select', components.CInputSelect)
Vue.component('c-input-list', components.CInputList)
