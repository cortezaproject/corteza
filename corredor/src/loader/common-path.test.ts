import { describe, it } from 'mocha'
import { expect } from 'chai'
import { CommonPath } from './index'

describe(__filename, () => {
  describe('common path finding', () => {
    it('should find the longest common path', () => {
      expect(CommonPath([])).to.equal('/')
      expect(CommonPath(['/foo/bar/baz', '/foo/bar', '/foo'])).to.equal('/foo')
      expect(CommonPath(['/foo/a', '/foo/b', '/bar/c'])).to.equal('/')
      expect(CommonPath(['/bar/a', '/bar/b', '/bar/c'])).to.equal('/bar')
      expect(CommonPath(['/foo/bar/baz', '/foo/bar', '/bar/bar'])).to.equal('/')
      expect(CommonPath(['/foo', '/bar', '/baz'])).to.equal('/')
      expect(CommonPath(['', '', ''])).to.equal('/')
    })
  })
})
