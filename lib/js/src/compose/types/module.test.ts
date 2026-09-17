import { expect } from 'chai'
import { Module } from './module'
import { AreObjectsOf } from '../../guards'
import { ModuleField } from './module-field/base'

const mod = new Module({
  name: 'modName',

  fields: [
    { name: 'f0' },
    { name: 'f1s', kind: 'String' },
    { name: 'f2b', kind: 'Bool' },
    { name: 'f3n', kind: 'Number' },
  ],
})

describe('module', () => {
  describe('check module casting', () => {
    it('simple assignment', () => {
      expect(mod.name).to.equal('modName')
      expect(mod.fields).to.be.lengthOf(4)
      expect(AreObjectsOf<ModuleField>(mod.fields, 'isSystem')).to.equal(true)
    })
  })

  describe('check field creation', () => {
    it('simple assignment', () => {
      const mod = new Module({
        name: 'modName',
        fields: [
          { name: 'f001', kind: 'Bool' },
          { name: 'f002', kind: 'DateTime' },
          { name: 'f003', kind: 'Email' },
          { name: 'f004', kind: 'File' },
          { name: 'f005', kind: 'Select' },
          { name: 'f006', kind: 'Number' },
          { name: 'f007', kind: 'Record' },
          { name: 'f008', kind: 'String' },
          { name: 'f009', kind: 'Url' },
          { name: 'f010', kind: 'User' },
        ],
      })

      expect(mod.fields).lengthOf(10)
    })
  })

  describe('records flag', () => {
    it('defaults to false and takes the server value', () => {
      expect(new Module({ name: 'm' }).hasRecords).to.equal(false)
      expect(new Module({ name: 'm', hasRecords: true }).hasRecords).to.equal(true)
    })

    it('survives a clone', () => {
      expect(new Module({ name: 'm', hasRecords: true }).clone().hasRecords).to.equal(true)
    })
  })

  describe('field operations', () => {
    it('simple search', () => {
      expect(mod.findField('f1s')).to.not.equal(undefined)
      expect(mod.findField('foo')).to.equal(undefined)
    })
  })
})
