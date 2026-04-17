import { compose } from '@planetcrust/human-js'

/**
 * Helper function to construct the proper chart sub type (if possible).
 * Inspects report metrics to determine chart type.
 * @param {compose.Chart} c Base chart object
 * @returns {compose.Chart|compose.FunnelChart|compose.GaugeChart|compose.RadarChart}
 */
export function chartConstructor(c) {
  for (const r of c.config.reports) {
    for (const m of r.metrics) {
      if (m.type === 'funnel') {
        return new compose.FunnelChart(c)
      } else if (m.type === 'gauge') {
        return new compose.GaugeChart(c)
      } else if (m.type === 'radar') {
        return new compose.RadarChart(c)
      }
    }
  }

  return new compose.Chart(c)
}
