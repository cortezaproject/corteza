import { expect } from 'chai'
import { Make, MakeAppURL } from './url'

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

  describe('make app url', () => {
    it('uses the configured webapp base', () => {
      expect(MakeAppURL({ app: 'workflow', path: '42/edit', webapp: 'https://host.tld/process', base: 'https://host.tld/process/admin/' }))
        .to.eq('https://host.tld/process/workflow/42/edit')
    })

    it('resolves a relative webapp base against the current location', () => {
      expect(MakeAppURL({ app: 'workflow', path: '42/edit', webapp: '/process/', base: 'https://host.tld/process/admin/' }))
        .to.eq('https://host.tld/process/workflow/42/edit')
    })

    it('falls back to the base tag of the current app under a path prefix', () => {
      expect(MakeAppURL({ app: 'workflow', path: '42/edit', webapp: '', base: 'https://host.tld/process/admin/' }))
        .to.eq('https://host.tld/process/workflow/42/edit')
    })

    it('falls back to the base tag of the current app at the root', () => {
      expect(MakeAppURL({ app: 'workflow', path: '42/edit', webapp: '', base: 'https://host.tld/admin/' }))
        .to.eq('https://host.tld/workflow/42/edit')
    })

    it('links to the app root without a path', () => {
      expect(MakeAppURL({ app: 'compose', webapp: '', base: 'https://host.tld/process/admin/' }))
        .to.eq('https://host.tld/process/compose')
    })
  })
})
