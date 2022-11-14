/* eslint-disable @typescript-eslint/ban-ts-ignore */
import { expect } from 'chai'
import { Record, Values } from './record'
import { Module } from './module'

const m = Object.freeze(new Module({
  fields: [
    { name: 'bool', kind: 'Bool' },
    { name: 'simple', kind: 'String' },
    { name: 'required', isRequired: true },
    { name: 'multi', isMulti: true },
    { name: 'multiRequired', isRequired: true, isMulti: true },
  ],
}))

describe(__filename, () => {
  describe('record creation', () => {
    const assertAllUndefined = function (r: Record): void {
      expect(r.module).to.eq(m)
      /**
       * It's extremely important that properties are set
       * even if to undefined. Without this we will have problems
       * with Vue reactivity
       */
      expect(r.values).to.have.property('bool').and.to.be.undefined
      expect(r.values).to.have.property('simple').and.to.be.undefined
      expect(r.values).to.have.property('required').to.be.undefined
      expect(r.values).to.have.property('multi').to.be.deep.eq([])
      expect(r.values).to.have.property('multiRequired').to.be.deep.eq([])
    }

    const assertSimpleSet = function (r: Record): void {
      expect(r.module).to.eq(m)
      expect(r.values.bool).to.eq('0')
      expect(r.values.simple).to.eq('foo')
      expect(r.values.multi).to.deep.eq(['bar'])
    }

    const values: Values = { bool: '0', simple: 'foo', multi: ['bar'] }
    const rawValues = [
      { name: 'bool' },
      { name: 'simple', value: 'foo' },
      { name: 'multi', value: 'bar' },
    ]

    it('module only', () => {
      assertAllUndefined(new Record(m))
    })

    it('module with empty object', () => {
      assertAllUndefined(new Record(m, {}))
    })

    it('empty object with module', () => {
      assertAllUndefined(new Record({}, m))
    })

    it('module with non-empty values object', () => {
      assertSimpleSet(new Record(m, values))
    })

    it('non-empty values object with module', () => {
      assertSimpleSet(new Record(values, m))
    })

    it('module with non-empty raw values', () => {
      assertSimpleSet(new Record(m, rawValues))
    })

    it('raw values with module', () => {
      assertSimpleSet(new Record(rawValues, m))
    })

    it('should not corrupt recordID when there are no values', () => {
      expect(new Record({ recordID: '42' }, m)).to.have.property('recordID').equal('42')
    })

    it('should handle garbage input', () => {
      // @ts-ignore
      expect(() => new Record(false, m)).to.throw
      // @ts-ignore
      expect(() => new Record(null, m)).to.throw
      // @ts-ignore
      expect(() => new Record('some string', m)).to.throw
      // @ts-ignore
      expect(() => new Record(42, m)).to.throw
    })

    it('should handle Record instance', () => {
      const rr = new Record({ recordID: '42' }, m)
      expect(new Record(rr)).to.have.property('recordID').equal('42')
    })
  })

  describe('multi value reading', () => {
    const bars = ['bar1', 'bar2', 'bar3', 'bar4']

    it('should have all values from explicit value setting', () => {
      const r = new Record(m)
      r.values.multi = bars
      expect(r.values.multi).to.be.lengthOf(4)
      expect(r.values.multi).to.deep.eq(bars)
    })

    it('should have all values from setValues (with Values)', () => {
      const r = new Record(m)
      r.setValues({ multi: bars })
      expect(r.values.multi).to.be.lengthOf(4)
      expect(r.values.multi).to.deep.eq(bars)
    })

    it('should have all values from setValues (with RawValue[])', () => {
      const r = new Record(m)
      r.setValues([
        { name: 'multi', value: 'bar1' },
        { name: 'multi', value: 'bar2' },
        { name: 'multi', value: 'bar3' },
        { name: 'multi', value: 'bar4' },
      ])
      expect(r.values.multi).to.be.lengthOf(4)
      expect(r.values.multi).to.deep.eq(bars)
    })

    it('should have all values from record initialization', () => {
      const r = new Record(m, [
        { name: 'multi', value: 'bar1' },
        { name: 'multi', value: 'bar2' },
        { name: 'multi', value: 'bar3' },
        { name: 'multi', value: 'bar4' },
      ])

      expect(r.values.multi).to.be.lengthOf(4)
      expect(r.values.multi).to.deep.eq(bars)
    })
  })

  describe('value setting', () => {
    let r: Record
    beforeEach(() => {
      r = new Record(m)
    })

    it('should properly set values via object', () => {
      r.values.simple = 'foo'
      expect(r.values.simple).to.eq('foo')
    })

    it('should properly set from values object', () => {
      r.setValues({ simple: 'foo' })
      expect(r.values.simple).to.eq('foo')
    })

    it('should properly set from array of raw-values objects', () => {
      r.setValues([{ name: 'simple', value: 'foo' }])
      expect(r.values.simple).to.eq('foo')
    })

    it('should properly set multiple values via values to a non-multi-value field', () => {
      r.setValues({ simple: ['bar', 'baz'] })
      expect(r.values.simple).to.eq('bar')
    })

    it.skip('should properly set value directly to a multi-value field', () => {
      r.values.multi = 'bar'
      expect(r.values.multi).to.deep.eq(['bar'])
    })

    it('should properly set value directly to a multi-value field', () => {
      r.values.multi = ['bar']
      expect(r.values.multi).to.deep.eq(['bar'])
    })

    it('should properly set value via setValues to a multi-value field', () => {
      r.setValues({ multi: 'bar' })
      expect(r.values.multi).to.deep.eq(['bar'])
    })

    it('should properly set multiple values via raw-values to a non-multi-value field', () => {
      r.setValues([{ name: 'simple', value: 'foo' }, { name: 'simple', value: 'foo' }])
      expect(r.values.simple).to.eq('foo')
    })

    it('should properly set value via raw-values to a multi-value field', () => {
      r.setValues([{ name: 'multi', value: 'bar' }])
      expect(r.values.multi).to.deep.eq(['bar'])
    })
  })

  describe('JSON serialization', () => {
    let r: Record
    beforeEach(() => {
      r = new Record(m)
    })

    it('should properly serialize whole record', () => {
      r.setValues([{ simple: 'foo', multi: ['bar', 'baz'] }])
      expect(JSON.stringify(r)).to.equal('{"recordID":"0","moduleID":"0","namespaceID":"0","revision":0,"values":[{"name":"simple","value":"foo"},{"name":"multi","value":"bar"},{"name":"multi","value":"baz"}],"valueErrors":{},"meta":{},"canUpdateRecord":false,"canReadRecord":false,"canDeleteRecord":false,"canManageOwnerOnRecord":false,"canSearchRevision":false,"canGrant":false}')
    })

    it('serialization magic should sustain object manipulation', () => {
      r.setValues([{ simple: 'foo', multi: ['bar', 'baz'] }])
      const { values } = r

      expect(JSON.stringify(values)).to.equal(
        '[{"name":"simple","value":"foo"},{"name":"multi","value":"bar"},{"name":"multi","value":"baz"}]')
    })
  })

  describe('prevent value corruption when when value prop is missing on rawValue and one of the fields is named "name"', () => {
    const m = Object.freeze(new Module({
      fields: [
        { name: 'name' },
        { name: 'alt' },
      ],
    }))

    it('it should not set value for name when constructing an object', () => {
      const r = new Record(m, [{ name: 'alt' }])
      expect(r.values.name).to.eq(undefined)
      expect(r.values.alt).to.eq(undefined)
    })

    it('it should set some values for name when constructing an object', () => {
      const r = new Record(m, [{ name: 'alt' }, { name: 'name', value: 'bar' }])
      expect(r.values.name).to.eq('bar')
      expect(r.values.alt).to.eq(undefined)
    })

    it('it should not set value for name when applying new values as array', () => {
      const r = new Record(m)
      r.apply([{ name: 'alt' }])
      expect(r.values.name).to.eq(undefined)
      expect(r.values.alt).to.eq(undefined)
    })

    it('it should set value for name when applying new values as object', () => {
      const r = new Record(m)
      r.apply({ name: 'alt' })
      expect(r.values.name).to.eq('alt')
      expect(r.values.alt).to.eq(undefined)
    })

    it('it should not set value for name when applying new values', () => {
      const r = new Record(m)
      r.apply([{ name: 'alt' }, { alt: 'foo' }])
      expect(r.values.name).to.eq('alt')
      expect(r.values.alt).to.eq('foo')
    })
  })
})
