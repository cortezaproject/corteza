import { describe, it } from 'mocha'
import { expect } from 'chai'

import pino from 'pino'
import ServerScripts from './server-scripts'
import { Trigger } from '../scripts/trigger'
import fs from 'fs'
import os from 'os'
import path from 'path'
import Loader from '../loader'

const baseScript = {
  src: 'path/to/script',
  name: 'scriptname',
  updatedAt: new Date(),
}

const svcCtorArgs = Object.freeze({
  logger: pino({ level: 'silent' }),
  config: {
    cServers: {
      system: { apiBaseURL: 'unit-test' },
      compose: { apiBaseURL: 'unit-test' },
    },
  },
})

describe('scripts list', () => {
  describe('empty', () => {
    it('should be empty', () => {
      const svc = new ServerScripts(svcCtorArgs)
      expect(svc.list()).to.have.lengthOf(0)
    })
  })

  describe('filled', () => {
    const svc = new ServerScripts(svcCtorArgs)
    beforeEach(() => {
      svc.update([
        {
          ...baseScript,
          label: 'label1',
          name: 'Name',
          description: 'deskription',
          errors: [],
          triggers: [
            new Trigger({ eventTypes: ['foo'], resourceTypes: ['res1'], constraints: [] }),
          ],
        },
        {
          ...baseScript,
          label: 'label2',
          name: 'Name',
          description: 'deskription',
          errors: [],
          triggers: [
            new Trigger({ eventTypes: ['afterMyThing'], resourceTypes: ['res2'], constraints: [] }),
          ],
        },
        {
          ...baseScript,
          label: 'label3',
          name: 'Name',
          description: 'deskription',
          errors: [],
          triggers: [
            new Trigger({
              eventTypes: ['beforeMyThing', 'afterMyThing'],
              resourceTypes: ['res2'],
              constraints: [],
            }),
          ],
        },
      ])
    })

    it('should match all 3 with "abel"', () => {
      expect(svc.list({ query: 'abel' })).to.have.lengthOf(3)
    })

    it('should match 1 with labe1', () => {
      expect(svc.list({ query: 'label1' })).to.have.lengthOf(1)
    })

    it('should match all 3 with description', () => {
      expect(svc.list({ query: 'deskr' })).to.have.lengthOf(3)
    })

    it('should match 2 with res2 for resource', () => {
      expect(svc.list({ resourceType: 'res2' })).to.have.lengthOf(2)
    })

    it('should match none with re for resource', () => {
      expect(svc.list({ resourceType: 're' })).to.have.lengthOf(0)
    })

    it('should match 2 with events', () => {
      expect(svc.list({ eventTypes: ['afterMyThing'] })).to.have.lengthOf(2)
    })

    it('should match 1 with event', () => {
      expect(svc.list({ eventTypes: ['beforeMyThing'] })).to.have.lengthOf(1)
    })

    it('should match all with no events', () => {
      expect(svc.list({ eventTypes: [] })).to.have.lengthOf(3)
    })
  })
})

describe('scripts processing', () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'corredor-test-'))
  const valid = path.join(dir, 'valid.js')
  const invalid = path.join(dir, 'invalid.js')

  before(() => {
    fs.writeFileSync(valid, 'exports.default = { exec () { return 42 } }')
    fs.writeFileSync(invalid, 'exports.default = { iterator () { throw new Error() }, exec () {} }')
  })

  after(() => {
    fs.rmSync(dir, { recursive: true, force: true })
  })

  it('gives only valid scripts their exec and keeps raw functions out of invalid ones', async () => {
    const loader = {
      searchPaths: [dir],
      scripts: async () => [
        { src: valid, name: 'valid', updatedAt: new Date(), triggers: [], errors: [] },
        {
          src: invalid,
          name: 'invalid',
          updatedAt: new Date(),
          triggers: [],
          errors: ['iterator not defined'],
        },
      ],
    } as unknown as Loader

    const svc = new ServerScripts({ ...svcCtorArgs, loader })
    await svc.process()
    await new Promise(resolve => setTimeout(resolve, 50))

    const byName = Object.fromEntries(svc.list().map(s => [s.name, s]))
    expect(byName.valid.exec).to.be.a('function')
    expect(byName.invalid).to.not.have.property('exec')
    expect(byName.invalid).to.not.have.property('iterator')
  })
})
