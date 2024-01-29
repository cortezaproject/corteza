import Vue from 'vue'
import { FontAwesomeIcon, FontAwesomeLayers } from '@fortawesome/vue-fontawesome'
import PortalVue from 'portal-vue'
import './faIcons'
import CContentHeader from './CContentHeader'
import CSystemFields from './CSystemFields'
import CResourceListStatusFilter from './CResourceListStatusFilter'
import { components } from '@cortezaproject/corteza-vue'

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

const { CCorredorManualButtons, CPermissionsButton, CInputConfirm, CButtonSubmit, CInputCheckbox, CInputSelect, CFormTableWrapper } = components

Vue.use(PortalVue)
Vue.component('c-corredor-manual-buttons', CCorredorManualButtons)
Vue.component('c-permissions-button', CPermissionsButton)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('c-content-header', CContentHeader)
Vue.component('c-resource-list-status-filter', CResourceListStatusFilter)
Vue.component('c-input-confirm', CInputConfirm)
Vue.component('c-button-submit', CButtonSubmit)
Vue.component('c-input-checkbox', CInputCheckbox)
Vue.component('c-system-fields', CSystemFields)
Vue.component('c-input-select', CInputSelect)
Vue.component('c-form-table-wrapper', CFormTableWrapper)

use([
  LineChart,
  SVGRenderer,
  TitleComponent,
  GridComponent,
  TooltipComponent,
])

Vue.component('e-charts', ECharts)
