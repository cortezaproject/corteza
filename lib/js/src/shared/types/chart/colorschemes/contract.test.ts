import { expect } from 'chai'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import colorschemes from './index'

/**
 * Cross-language contract test.
 *
 * compose_chart_create refuses a colorScheme it does not recognise, and the list
 * it checks against is a generated copy of these tables, embedded at
 * server/compose/agentic/chart_color_schemes.json. The copy has to stay exact in
 * both directions: a name missing from it makes the server reject a scheme the
 * webapp can draw, and a name only in it lets through one the webapp cannot —
 * which is the silent failure the check exists to stop (an unresolved scheme
 * reaches echarts as an undefined palette, so the chart renders its legend and
 * none of its series).
 *
 * Regenerate the embedded list when a scheme is added or renamed:
 *   node lib/js/src/shared/types/chart/colorschemes/emit-contract.mjs
 */

const EMBEDDED_PATH = resolve(
  __dirname,
  '../../../../../../../server/compose/agentic/chart_color_schemes.json',
)

const embedded: string[] = JSON.parse(readFileSync(EMBEDDED_PATH, 'utf8'))

const actual: string[] = Object.entries(
  colorschemes as Record<string, Record<string, unknown>>,
).flatMap(([family, schemes]) => Object.keys(schemes).map(name => `${family}.${name}`))

describe('server chart colorScheme list matches the webapp colour tables', () => {
  it('carries every scheme the webapp can resolve', () => {
    expect([...embedded].sort()).to.deep.equal([...actual].sort())
  })

  it('is family-qualified and sorted, as the server embeds it', () => {
    expect(embedded).to.deep.equal([...embedded].sort())
    for (const name of embedded) {
      expect(name, `${name} must be <family>.<Name>`).to.match(/^[a-z]+\.[A-Za-z0-9_]+$/)
    }
  })

  // The near-miss that produced an invisible chart on a real dashboard: the
  // number is the swatch count and part of the name, so the wrong one is a
  // different scheme, not a variant of the same one.
  it('distinguishes schemes that differ only by swatch count', () => {
    expect(actual).to.include('tableau.ClassicOrangeBlue13')
    expect(actual).to.not.include('tableau.ClassicOrangeBlue7')
  })
})
