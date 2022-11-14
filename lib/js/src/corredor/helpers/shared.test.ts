import { describe, it, beforeEach } from 'mocha'
import { expect } from 'chai'
import sinon, { stubObject } from 'ts-sinon'
import { extractID, genericPermissionUpdater } from './shared'
import { NoID } from '../../cast'
import { System as SystemAPI } from '../../api-clients'

describe(__filename, () => {
  beforeEach(() => {
    sinon.restore()
  })

  describe('extractID', () => {
    it('should extract the ID', () => {
      expect(extractID(4200001)).to.equal('4200001')
      expect(extractID('4200002')).to.equal('4200002')
      expect(extractID({ prop: '4200003' }, 'prop')).to.equal('4200003')
    })

    it('should throw error on invalid input', () => {
      expect(() => extractID('abc')).to.throw()
      expect(() => extractID({ id: 'abc' }, 'id')).to.throw()
    })

    it('should extract 0', () => {
      expect(extractID([]), 'array as ID should equal NoID').to.equal(NoID)
      expect(extractID({}), 'empty object as ID should equal NoID').to.equal(NoID)
      expect(extractID(''), 'empty string as ID should equal NoID').to.equal(NoID)
      expect(extractID(), 'nothing as ID should equal NoID').to.equal(NoID)
    })
  })

  describe('generic permission updater', () => {
    it('should do group by role and make multiple calls to the API (one per role', async () => {
      const API = stubObject<SystemAPI>(new SystemAPI({}))

      await genericPermissionUpdater(API, [
        { role: { roleID: '1111' }, resource: { resourceID: 'r.0001' }, operation: 'read', access: '1' },
        { role: { roleID: '1111' }, resource: { resourceID: 'r.0001' }, operation: 'read', access: '1' },
        { role: { roleID: '5555' }, resource: { resourceID: 'r.0001' }, operation: 'read', access: '1' },
        { role: { roleID: '5555' }, resource: { resourceID: 'r.0001' }, operation: 'read', access: '1' },
      ])

      expect(API.permissionsUpdate.calledTwice).true
    })
  })
})
