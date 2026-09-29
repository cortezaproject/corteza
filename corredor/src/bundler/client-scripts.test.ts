import { describe, it } from 'mocha'
import { expect } from 'chai'
import fs from 'fs'
import os from 'os'
import path from 'path'
import pino from 'pino'
import bundler from './client-scripts'

interface Bundled {
  exposed: { scripts?: Array<{ exec: () => unknown }> }
  loadError?: Error
  warnings: string[]
}

// Bundles one client script, loads the bundle and returns what it exposes plus the warnings logged
async function bundle(source: string): Promise<Bundled> {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'corredor-bundle-'))
  const warnings: string[] = []
  const log = pino({ level: 'warn' }, { write: (line: string) => warnings.push(line) })

  try {
    fs.writeFileSync(path.join(dir, 'Script.js'), source)
    fs.writeFileSync(
      path.join(dir, 'entry.js'),
      ["import script from './Script.js'", 'export const scripts = [script]'].join('\n'),
    )

    await bundler.Pack('test', path.join(dir, 'entry.js'), dir, dir, log)

    const exposed: { testClientScripts?: object } = {}
    let loadError: Error | undefined

    try {
      new Function(fs.readFileSync(path.join(dir, 'test.client-scripts.js'), 'utf8')).call(exposed)
    } catch (e) {
      loadError = e as Error
    }

    return { exposed: exposed.testClientScripts ?? {}, loadError, warnings }
  } finally {
    fs.rmSync(dir, { recursive: true, force: true })
  }
}

describe('client scripts bundler', function () {
  this.timeout(30000)

  it('bundles a script that imports a Node built-in', async () => {
    const { exposed, loadError, warnings } = await bundle(
      [
        "import { stringify } from 'querystring'",
        'export default { exec: () => stringify({ a: 1 }) }',
      ].join('\n'),
    )

    expect(loadError).to.equal(undefined)
    expect(warnings).to.have.lengthOf(0)
    expect(exposed.scripts?.[0].exec()).to.equal('a=1')
  })

  it('reports an import it cannot resolve, naming the script', async () => {
    const { loadError, warnings } = await bundle(
      [
        "import missing from 'no-such-module-anywhere'",
        'export default { exec: () => missing }',
      ].join('\n'),
    )

    expect(loadError?.message).to.contain('no-such-module-anywhere')
    expect(warnings.join('\n')).to.contain('no-such-module-anywhere').and.to.contain('Script.js')
  })
})
