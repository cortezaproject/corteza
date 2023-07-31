import Vue from 'vue'

import resourceTranslations from './resource-translations'

import toast from './toast'

import './eventbus'

Vue.mixin(resourceTranslations)

Vue.mixin(toast)
