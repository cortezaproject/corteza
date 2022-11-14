import Vue from 'vue'
import { capitalize } from 'lodash'

import Card from './Card'
import Plain from './Plain'

const Registry = {
  Card,
  Plain,
}

const defaultWrap = 'Card'

/**
 * @param block {compose.PageBLock}
 * @returns component
 */
function GetWrapComponent ({ block, wrap = defaultWrap }) {
  const cmpName = capitalize(wrap)
  if (Object.hasOwnProperty.call(Registry, cmpName)) {
    return Registry[capitalize(cmpName)]
  }

  throw new Error('unknown wrap: ' + wrap)
}

/**
 * Wraps page block with one of the configured (on page block options) components
 */
export default Vue.component('page-block', {
  functional: true,

  render (ce, ctx) {
    return ce(GetWrapComponent(ctx.props), ctx.data, ctx.children)
  },
})
