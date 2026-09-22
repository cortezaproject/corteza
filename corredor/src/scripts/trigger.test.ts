import { describe, it } from 'mocha'
import { expect } from 'chai'
import Make, { Trigger } from './trigger'

describe('trigger making', () => {
  // sample trigger
  const st = new Trigger()

  describe('from array', () => {
    it('should return empty array on empty array', () => {
      expect(Make([])).is.deep.eq([])
    })

    it('should filter out non-trigger values', () => {
      expect(Make([true, false, 42, '42', st, {}])).is.deep.eq([st])
    })
  })

  describe('from non-array', () => {
    it('should cast into array', () => {
      expect(Make(st)).is.deep.eq([st])
    })
  })

  describe('from function', () => {
    it('should collect all triggers', () => {
      expect(
        Make(function () {
          return [st, st]
        }),
      ).is.deep.eq([st, st])
    })
  })

  describe('from function*', () => {
    it('should collect all triggers', () => {
      expect(
        Make(function* () {
          yield st
          yield st
        }),
      ).is.deep.eq([st, st])
    })
  })

  describe('from nothing', () => {
    it('should handle undefined', () => {
      expect(Make(undefined)).is.deep.eq([])
    })
    it('should handle null', () => {
      expect(Make(null)).is.deep.eq([])
    })
    it('should handle boolean', () => {
      expect(Make(false)).is.deep.eq([])
      expect(Make(true)).is.deep.eq([])
    })
    it('should handle fn w/ void', () => {
      expect(
        Make(() => {
          /* */
        }),
      ).is.deep.eq([])
    })
  })

  describe('deferred', () => {
    it('should make deferred trigger with type & resource', () => {
      expect(Make(({ every }) => every('* * * * *'))).to.deep.eq([
        {
          eventTypes: ['onInterval'],
          resourceTypes: ['system'],
          constraints: [{ name: 'interval', value: ['* * * * *'] }],
          uiProps: [],
        },
      ])
    })
  })
})
