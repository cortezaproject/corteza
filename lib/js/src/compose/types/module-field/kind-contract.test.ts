import { expect } from 'chai'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

import { ModuleFieldRegistry } from './index'

/**
 * Cross-language contract test.
 *
 * server/compose/types/module_field.go holds ModuleFieldKinds, the canonical
 * list, snapshotted (by a Go test) to testdata/module_field_kinds.json. The
 * agentic module tools build their documentation from it, so a kind named
 * there is a kind agents will create.
 *
 * Nothing rejects an unknown kind: the DAL falls through to TypeText and
 * ModuleFieldMaker returns a plain ModuleField, so the webapp renders a String
 * editor and the field looks like it works. That silence is why this test
 * exists — the agentic tools advertised Currency and Duration, which no class
 * here has ever registered, and omitted Geometry, which one does.
 *
 * If the Go side changes, regenerate the snapshot first:
 *   cd server && go test ./compose/types/ -run ModuleFieldKindsSnapshot -update-field-kinds
 */

const SNAPSHOT_PATH = resolve(
  __dirname,
  '../../../../../../server/compose/types/testdata/module_field_kinds.json',
)

const advertised: string[] = JSON.parse(readFileSync(SNAPSHOT_PATH, 'utf8'))

describe('server field kinds match the webapp field classes', () => {
  const registered = [...ModuleFieldRegistry.keys()].sort()

  it('every advertised kind has a webapp class', () => {
    const missing = advertised.filter(kind => !ModuleFieldRegistry.has(kind))

    expect(missing, 'kinds the server advertises with no webapp field class').to.deep.equal([])
  })

  it('every webapp field class is advertised', () => {
    const unadvertised = registered.filter(kind => !advertised.includes(kind))

    expect(unadvertised, 'field kinds the webapp registers that the server never names').to.deep.equal(
      [],
    )
  })

  it('names each kind exactly once', () => {
    expect(advertised).to.have.lengthOf(new Set(advertised).size)
  })
})
