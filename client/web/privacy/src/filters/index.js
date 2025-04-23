import Vue from 'vue'
import { filters } from '@cortezaproject/corteza-vue'

for (const n in filters) {
  Vue.filter(n, filters[n])
}
