import { expect } from 'chai'
import { Make } from './url'

const pr = 'https'
const hs = 'www.test.tld'
const pt = '/path'
const arr = { a: ['f1', 'f2'] }
const simple = { s: 'src' }
const hash = 'hash'

const ref = `${pr}://www.ref.tld`

describe(__filename, () => {
  describe('make url', () => {
    it('entire url provided', () => {
      const test = Make({ url: `${pr}://${hs}${pt}`, ref })
      expect(test).to.eq(`${pr}://${hs}${pt}`)
    })

    it('no proto provided - with path', () => {
      const test = Make({ url: `//${hs}${pt}`, ref })
      expect(test).to.eq(`${pr}://${hs}${pt}`)
    })

    it('no proto provided - no path', () => {
      const test = Make({ url: `//${hs}`, ref })
      expect(test).to.eq(`${pr}://${hs}/`)
    })

    it('path provided', () => {
      const test = Make({ url: `${pt}`, ref })
      expect(test).to.eq(`${ref}${pt}`)
    })

    it('nothing provided', () => {
      const test = Make({ ref })
      expect(test).to.eq(`${ref}/`)
    })
  })

  describe('make query string', () => {
    it('ignore undefined keys', () => {
      const test = Make({ ref, query: { k: undefined } })
      expect(test).to.eq(`${ref}/`)
    })

    it('keep null keys', () => {
      const test = Make({ ref, query: { k: null } })
      expect(test).to.eq(`${ref}/?k=`)
    })

    it('simple structs', () => {
      const test = Make({ ref, query: simple })
      expect(test).to.eq(`${ref}/?s=src`)
    })

    it('arrays', () => {
      const test = Make({ ref, query: arr })
      expect(test).to.eq(`${ref}/?a[]=f1&a[]=f2`)
    })

    it('mix', () => {
      const test = Make({ ref, query: { ...arr, ...simple } })
      expect(test).to.contain('a[]=f1&a[]=f2')
      expect(test).to.contain('s=src')
    })
  })

  it('add hash', () => {
    const test = Make({ hash, ref })
    expect(test).to.eq(`${ref}/#${hash}`)
  })
})
