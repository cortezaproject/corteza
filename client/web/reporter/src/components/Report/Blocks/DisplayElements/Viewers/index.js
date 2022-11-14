import Vue from 'vue'
import { capitalize } from 'lodash'
import { components } from '@cortezaproject/corteza-vue'

// Need to import it here since it inits when imported
// Order is important, colorschemes first others after
/* eslint-disable no-unused-vars */
import colorschemes from 'chartjs-plugin-colorschemes'
import Funnel from 'chartjs-plugin-funnel'

const {
  CReportChart,
  CReportMetric,
  CReportTable,
  CReportText,
} = components

/**
 * List of all known display element components
 *
 */
const Registry = {
  Chart: CReportChart,
  Metric: CReportMetric,
  Table: CReportTable,
  Text: CReportText,
}

function GetComponent ({ displayElement }) {
  if (!displayElement) {
    throw new Error('displayElement prop missing')
  }

  const { kind } = displayElement
  if (Object.hasOwnProperty.call(Registry, capitalize(kind))) {
    return Registry[kind]
  }

  throw new Error('unknown displayElement kind: ' + kind)
}

export default Vue.component('display-element', {
  functional: true,

  render (ce, ctx) {
    return ce(GetComponent(ctx.props), ctx.data, ctx.children)
  },
})
