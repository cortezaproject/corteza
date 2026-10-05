import { expect } from 'chai'
import { User } from './user'

describe('system User', () => {
  it('keeps the kind the server sends', () => {
    expect(new User({ userID: '1', kind: 'sys' }).kind).to.equal('sys')
    expect(new User({ userID: '1' }).kind).to.equal('')
  })
})
