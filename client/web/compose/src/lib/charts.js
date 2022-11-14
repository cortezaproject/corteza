import { compose } from '@cortezaproject/corteza-js'

/**
 * Helper function to construct the proper chart sub type (if possible)
 * @param {compose.Chart} c Base chart object
 */
export function chartConstructor (c) {
  for (const r of c.config.reports) {
    for (const m of r.metrics) {
      if (m.type === 'funnel') {
        return new compose.FunnelChart(c)
      } else if (m.type === 'gauge') {
        return new compose.GaugeChart(c)
      }
    }
  }

  return new compose.Chart(c)
}
