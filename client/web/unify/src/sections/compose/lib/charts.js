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

/**
 * Breaks a label onto lines no wider than `width`, at spaces only: a word wider
 * than that keeps a line of its own rather than being split.
 * @param {string|number} text
 * @param {number} width
 * @param {(text: string) => number} measure Rendered width of a piece of text
 * @returns {string}
 */
export function wrapLabel(text, width, measure) {
  const lines = []

  for (const word of String(text).split(' ')) {
    const last = lines.length - 1
    if (last >= 0 && measure(`${lines[last]} ${word}`) <= width) {
      lines[last] += ` ${word}`
    } else {
      lines.push(word)
    }
  }

  return lines.join('\n')
}
