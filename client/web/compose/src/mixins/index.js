import Vue from 'vue'

import toast from './toast'
import resourceTranslations from './resource-translations'
import vueSelectPosition from './vue-select-position'
import uiHelpers from './uiHelpers'

Vue.mixin(toast)
Vue.mixin(resourceTranslations)
Vue.mixin(vueSelectPosition)
Vue.mixin(uiHelpers)
