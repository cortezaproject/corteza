import { expect } from 'chai'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import { PageBlockMaker, PageBlockRegistry } from './index'
import Feed from './calendar/feed'
import { PageBlockMetric } from './metric'

/**
 * Cross-language contract test.
 *
 * The server's compose_page_block_schema MCP tool serves the structs in
 * server/compose/types/page_block_options.go, snapshotted (by a Go test) to
 * server/compose/types/testdata/page_block_option_schemas.json. Agents build
 * pages from that schema — so every option key it advertises MUST be a key
 * these page-block classes (the webapp renderers' contract) actually read.
 * A key the webapp ignores fails silently as broken UI (the historical
 * `module` vs `moduleID` bug), which is exactly what this test catches.
 *
 * If the Go side changes, regenerate the snapshot first:
 *   cd server && go test ./compose/types/ -run PageBlockOptionSchemasSnapshot -update-block-schemas
 */

const SNAPSHOT_PATH = resolve(
  __dirname,
  '../../../../../../server/compose/types/testdata/page_block_option_schemas.json',
)

// Block kinds the server advertises but the unify webapp has no renderer
// for. Shrink this list when a renderer lands; never grow it silently.
const UNRENDERED_KINDS = ['SocialFeed']

const schemas: Record<string, Record<string, unknown>> = JSON.parse(
  readFileSync(SNAPSHOT_PATH, 'utf8'),
)

const keysOf = (o: unknown): string[] => Object.keys(o as Record<string, unknown>)

describe('server page-block schemas match webapp page-block contracts', () => {
  for (const [kind, schema] of Object.entries(schemas)) {
    if (UNRENDERED_KINDS.includes(kind)) {
      it(`${kind}: stays unrendered (no webapp class registered)`, () => {
        expect(PageBlockRegistry.get(kind)).to.be.undefined
      })
      continue
    }

    it(`${kind}: advertised option keys exist in the webapp class defaults`, () => {
      expect(PageBlockRegistry.get(kind), `no page-block class for kind ${kind}`).to.not.be
        .undefined

      const feKeys = keysOf(PageBlockMaker({ kind }).options)
      for (const key of keysOf(schema)) {
        expect(feKeys, `${kind}.options.${key} is advertised by the server schema`).to.include(key)
      }
    })
  }

  // The other direction, and the one that used to go unwatched: a key the
  // webapp reads but the server never advertises is a key no agent can discover.
  // It fails as a block built exactly to the published contract that still does
  // not work — a RecordOrganizer with no `group` shows nothing, an Automation
  // button with no `automationID` cannot reach a TAQ.
  it('every rendered block kind has a server schema', () => {
    const unadvertised = [...PageBlockRegistry.keys()].filter(kind => !schemas[kind])

    expect(unadvertised, 'block kinds the webapp renders with no server schema').to.deep.equal([])
  })

  for (const [kind, schema] of Object.entries(schemas)) {
    if (UNRENDERED_KINDS.includes(kind)) continue

    it(`${kind}: every option the webapp reads is advertised`, () => {
      const srvKeys = keysOf(schema)
      const missing = keysOf(PageBlockMaker({ kind }).options).filter(k => !srvKeys.includes(k))

      expect(
        missing,
        `${kind} options the webapp reads but the server never advertises`,
      ).to.deep.equal([])
    })
  }

  it('Metric: advertised metric item keys exist in the metric defaults', () => {
    const feMetricKeys = keysOf(new PageBlockMetric().makeMetric())
    const [item] = schemas.Metric.metrics as Record<string, unknown>[]
    for (const key of keysOf(item)) {
      expect(feMetricKeys, `Metric metrics[].${key}`).to.include(key)
    }
  })

  it('Calendar: advertised feed keys exist in the Feed class', () => {
    const feed = new Feed()
    const feFeedKeys = keysOf(feed)
    const [item] = schemas.Calendar.feeds as Record<string, unknown>[]
    for (const key of keysOf(item)) {
      if (key === 'options') continue
      expect(feFeedKeys, `Calendar feeds[].${key}`).to.include(key)
    }
    const feOptKeys = keysOf(feed.options)
    for (const key of keysOf((item as { options: object }).options)) {
      expect(feOptKeys, `Calendar feeds[].options.${key}`).to.include(key)
    }
  })

  it('Progress: advertised nested option keys exist in the progress defaults', () => {
    const options = PageBlockMaker({ kind: 'Progress' }).options as Record<string, object>
    for (const nested of ['value', 'minValue', 'maxValue', 'display']) {
      const feNested = keysOf(options[nested])
      for (const key of keysOf((schemas.Progress as Record<string, object>)[nested])) {
        if (nested === 'display' && key === 'thresholds') continue // FE default is an empty list
        expect(feNested, `Progress ${nested}.${key}`).to.include(key)
      }
    }
  })
})
