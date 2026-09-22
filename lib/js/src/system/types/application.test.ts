import { expect } from 'chai'
import { Application } from './application'

describe('system Application', () => {
  // A custom application is only one while its kind survives a round trip
  // through this class: the editor saves what the class holds.
  it('keeps the kind and source meta of a custom application across a clone', () => {
    const app = new Application({
      applicationID: '514791068074901505',
      unify: {
        name: 'Desk',
        listed: true,
        url: 'app/514791068074901505',
        config: '',
        iconID: '0',
        logoID: '0',
        kind: 'custom',
      },
      sourceMeta: { size: 7520, namespace: 'crm', modules: ['Lead'] },
      canManageSourceOnApplication: true,
    })

    const copy = app.clone()
    expect(copy.unify?.kind).to.equal('custom')
    expect(copy.sourceMeta).to.deep.equal({ size: 7520, namespace: 'crm', modules: ['Lead'] })
    expect(copy.canManageSourceOnApplication).to.equal(true)
  })

  // Owner IDs are 64-bit; as a number the last digits are lost.
  it('holds the owner as an ID string, not a number', () => {
    const app = new Application({ ownerID: '508584512737837057' })
    expect(app.ownerID).to.equal('508584512737837057')
  })

  it('defaults to a section or link with no source and no right to change one', () => {
    const app = new Application()
    expect(app.unify?.kind).to.equal('')
    expect(app.sourceMeta).to.deep.equal({ size: 0 })
    expect(app.canManageSourceOnApplication).to.equal(false)
  })
})
