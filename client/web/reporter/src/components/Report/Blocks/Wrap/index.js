import Vue from 'vue'
import { capitalize } from 'lodash'

import Card from './Card'

const Registry = {
  Card,
}

const defaultWrap = 'Card'

/**
 * @param block {reporter.Block}
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
export default Vue.component('block', {
  functional: true,

  render (ce, ctx) {
    return ce(GetWrapComponent(ctx.props), ctx.data, ctx.children)
  },
})
