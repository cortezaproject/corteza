import { describe, it } from 'mocha'
import { expect } from 'chai'

import pino from 'pino'
import ServerScripts from './server-scripts'
import { Trigger } from '../scripts/trigger'

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
