import { expect } from 'chai'
import { ConstraintMaker } from './constraints'

describe(__filename, () => {
  describe('maker', () => {
    it('should make simple equal matcher', () => {
      expect(ConstraintMaker({ value: ['foo'] }).Match('foo')).to.be.true
    })
    it('should make simple like matcher', () => {
      expect(ConstraintMaker({ op: 'like', value: ['f*'] }).Match('foo')).to.be.true
    })
    it('should make simple match matcher', () => {
      expect(ConstraintMaker({ op: '~', value: ['^foo$'] }).Match('foo')).to.be.true
    })
    it('should make simple (not) equal matcher', () => {
      expect(ConstraintMaker({ op: '!=', value: ['foo'] }).Match('bar')).to.be.true
    })
    it('should make simple (not) like matcher', () => {
      expect(ConstraintMaker({ op: 'not like', value: ['f*'] }).Match('bar')).to.be.true
    })
    it('should make simple (not) match matcher', () => {
      expect(ConstraintMaker({ op: '!~', value: ['^foo$'] }).Match('bar')).to.be.true
    })
  })
})
