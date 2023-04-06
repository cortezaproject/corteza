import Vue from 'vue'

import resourceTranslations from './resource-translations'

import toast from './toast'
import vueSelectPosition from './vue-select-position'

import './eventbus'

Vue.mixin(resourceTranslations)

Vue.mixin(toast)
Vue.mixin(vueSelectPosition)
