import Vue from 'vue'

import toast from './toast'
import resourceTranslations from './resource-translations'
import uiHelpers from './uiHelpers'

Vue.mixin(toast)
Vue.mixin(resourceTranslations)
Vue.mixin(uiHelpers)
