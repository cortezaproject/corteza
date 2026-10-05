import { expect } from 'chai'
import { Apply, HumanID, NoID } from './cast'
class Foo {
  baz = ''
  bool = false
  arr = []
  ff: Foo[] = []
}

describe('cast', () => {
  describe('ensure that applying actually works', () => {
    let foo: Foo

    beforeEach(() => {
      foo = new Foo()
    })

    it('should handle simple assignment', () => {
      Apply(foo, { baz: 'value' }, String, 'baz')
      expect(foo.baz).to.equal('value')
    })

    it('should assign number to string', () => {
      Apply(foo, { baz: 42 }, String, 'baz')
      expect(foo.baz).to.equal('42')
    })

    it('should assign number to string using default (string) caster', () => {
      Apply(foo, { baz: 84 }, 'baz')
      expect(foo.baz).to.equal('84')
    })

    it('should assign boolean (empty string)', () => {
      Apply(foo, { bool: '' }, Boolean, 'bool')
      expect(foo.bool).to.equal(false)
    })

    it('should assign boolean (non-empty string)', () => {
      Apply(foo, { bool: 'true' }, Boolean, 'bool')
      expect(foo.bool).to.equal(true)
    })

    it('should assign boolean (number, 1)', () => {
      Apply(foo, { bool: 1 }, Boolean, 'bool')
      expect(foo.bool).to.equal(true)
    })

    it('should assign boolean (number, 0)', () => {
      Apply(foo, { bool: 0 }, Boolean, 'bool')
      expect(foo.bool).to.equal(false)
    })

    it('should assign simple array', () => {
      Apply(foo, { arr: [42, '42'] }, o => o, 'arr')
      expect(foo.arr).to.deep.equal([42, '42'])
    })

    it('should assign array of Foos', () => {
      // just a primitive check, not actually casting
      // to Foo[] so we're not actually checking for that
      Apply(foo, { ff: [{ baz: 'one' }, { baz: 'two' }] }, o => o, 'ff')
      expect(foo.ff).to.deep.equal([{ baz: 'one' }, { baz: 'two' }])
    })

    it('should cast null to empty string', () => {
      Apply(foo, { baz: null }, String, 'baz')
      expect(foo.baz).to.equal('')
    })
  })
})

describe('HumanID', () => {
  it('keeps a string ID as it is', () => {
    expect(HumanID('516687706987429889')).to.eq('516687706987429889')
  })

  it('reads an empty value as no ID', () => {
    expect(HumanID(undefined)).to.eq(NoID)
    expect(HumanID(0)).to.eq(NoID)
  })

  it('accepts a number that is exact', () => {
    expect(HumanID(123)).to.eq('123')
  })

  it('refuses a number past 2^53, which is already a different ID', () => {
    expect(() => HumanID(516687706987429889)).to.throw('pass IDs as strings')
  })
})
