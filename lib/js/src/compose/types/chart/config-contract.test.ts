import { expect } from 'chai'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import Chart from './chart'
import FunnelChart from './funnel'
import GaugeChart from './gauge'
import RadarChart from './radar'

/**
 * Cross-language contract test.
 *
 * The compose_chart_create MCP tool takes a chart config from an agent and
 * stores it; these classes are what has to draw it. They validate only at draw
 * time, so a config the two sides disagree about is stored without complaint and
 * surfaces as "Metrics aggregate not defined" where the chart should be — which
 * is exactly what happened to agent-built charts.
 *
 * The fixture is shared with server/compose/agentic/chart_handler_test.go: the
 * configs listed there as rejected must fail here too, with the named error, and
 * the ones listed as accepted must validate.
 */

const FIXTURE_PATH = resolve(
  __dirname,
  '../../../../../../server/compose/agentic/testdata/chart_config_cases.json',
)

interface Case {
  name: string
  error?: string
  config: Record<string, any>
}

const cases: { rejected: Array<Case>; accepted: Array<Case> } = JSON.parse(
  readFileSync(FIXTURE_PATH, 'utf8'),
)

// Mirrors the webapp's chartConstructor
// (client/web/unify/src/sections/compose/lib/charts.js): the metric type is what
// picks the class, and each subtype validates differently.
const makeChart = (config: Record<string, any>) => {
  for (const r of config.reports || []) {
    for (const m of r.metrics || []) {
      if (m.type === 'funnel') return new FunnelChart({ config })
      if (m.type === 'gauge') return new GaugeChart({ config })
      if (m.type === 'radar') return new RadarChart({ config })
    }
  }

  return new Chart({ config })
}

describe('MCP chart configs match what the webapp can render', () => {
  for (const c of cases.rejected) {
    it(`rejects: ${c.name}`, () => {
      expect(() => makeChart(c.config).isValid()).to.throw(
        `notification.chart.invalidConfig.${c.error}`,
      )
    })
  }

  for (const c of cases.accepted) {
    it(`renders: ${c.name}`, () => {
      expect(() => makeChart(c.config).isValid()).to.not.throw()
    })
  }
})
