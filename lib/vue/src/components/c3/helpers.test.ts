import { expect } from 'chai'
import { ExtractComponents, extractGroups } from './helpers'

const case1 = [['a'], ['a', 'b'], ['a', 'b'], ['a', 'c', 'd', '123'], ['c'], []]

describe(__filename, () => {
  describe('extractGroups', () => {
    it('return empty array when there are no groups', () => {
      expect(extractGroups([])).to.be.empty
    })
    it('should return unique groups on the 1st level', () => {
      expect(extractGroups(case1)).to.eql(['a', 'c'])
    })
    it('should return unique groups on the 2nd level', () => {
      expect(extractGroups(case1, 'a')).to.eql(['b', 'c'])
    })
    it('should return nothing on nonexisgin paths', () => {
      expect(extractGroups(case1, 'foo')).to.eql([])
    })
    it('should properly handle deep paths', () => {
      expect(extractGroups(case1, 'a', 'c', 'd')).to.eql(['123'])
    })
  })

  describe('ExtractComponents', () => {
    it('nothing in, nothing out', () => {
      expect(ExtractComponents({})).to.be.empty
    })

    it('match all on empty path', () => {
      expect(ExtractComponents({
        a: { name: 'a' },
        b: { name: 'b' },
        c: { name: 'c', group: ['foo'] },
      })).to.be.eql([
        { name: 'a' },
        { name: 'b' },
      ])
    })

    it('root match', () => {
      expect(ExtractComponents({
        a: { name: 'a', group: ['foo'] },
        b: { name: 'b' },
      }, 'foo')).to.be.eql([
        { name: 'a', group: ['foo'] },
      ])
    })
  })
})
